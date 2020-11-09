package iosocket

import (
	"container/list"
	"encoding/json"
	"time"

	"github.com/astaxie/beego"
	"github.com/gorilla/websocket"
	"github.com/qiaogw/models"
)

// Subscription 连接者消息
type Subscription struct {
	Archive []models.WsEvent      // All the events from the archive.
	New     <-chan models.WsEvent // New events coming in.
}

// Message 消息获得结构
type Message struct {
	Ws  *websocket.Conn
	Evt models.WsEvent
}

//Subscriber 连接者
type Subscriber struct {
	Id     int             `json:"id,omitempty"`
	Name   string          `json:"name,omitempty"`
	Avatar string          `json:"avatar,omitempty"` //头像
	RoomId string          `json:"roomId,omitempty"`
	Conn   *websocket.Conn // Only for WebSocket users; otherwise nil.
}

var (
	//新建立用户频道chan.
	Subscribe = make(chan Subscriber, 10)
	//退出用户频道chan.
	UnSubscribe = make(chan Subscriber, 10)
	//发送事件发布。chan
	Publish = make(chan Message, 10)
	//长轮询候补名单。队列chan
	WaitingList = list.New()
	//用户队列链表
	Subscribers = list.New()
)

func init() {
	go chatroom()
}

//NewEvent 创建新websocket用户
func NewEvent(ep, username string, userid int, targetid, roomid, msg string, userlist interface{}) models.WsEvent {
	return models.WsEvent{ep, username, userid, targetid, roomid, int(time.Now().Unix()), msg, userlist}
}

//Join 建立连接
func Join(userid int, username, roomid string, ws *websocket.Conn) {
	Subscribe <- Subscriber{Id: userid, Name: username, RoomId: roomid, Conn: ws}
}

//Leave 断开连接
func Leave(userid int, username, roomid string, ws *websocket.Conn) {
	UnSubscribe <- Subscriber{Id: userid, Name: username, RoomId: roomid, Conn: ws}
}

//chatroom 此函数处理所有传入的chan消息。
func chatroom() {
	for {
		select {
		//监听建立连接chan有数据
		case sub := <-Subscribe:
			if !isUserExist(Subscribers, sub.Id) {
				Subscribers.PushBack(sub) //将用户添加到列表的末尾。.
				// //发布JOIN事件。
				userlist := GetUserlist(Subscribers, sub.RoomId)
				evt := NewEvent(models.EVENT_JOIN, sub.Name, sub.Id, sub.Avatar, sub.RoomId, "", userlist)
				BroadcastWebSocket(evt)
				welcomevt := NewEvent(models.WELCOM_MESSAGE, sub.Name, sub.Id, sub.Avatar, sub.RoomId, "欢迎"+sub.Name+"加入聊天，请文明发言！", userlist)
				//beego.Debug("欢迎" + sub.Name + "加入聊天，请文明发言！")
				BroadcastWebSocket(welcomevt)
				models.NewArchive(evt)
				//beego.Info("New user:", sub.Id, ";WebSocket:", sub.Conn != nil)
			} else {
				for unsub := Subscribers.Front(); unsub != nil; unsub = unsub.Next() {
					if unsub.Value.(Subscriber).Id == sub.Id {
						Subscribers.Remove(unsub)
						userlist := GetUserlist(Subscribers, sub.RoomId)
						evt := NewEvent(models.EVENT_LEAVE, sub.Name, sub.Id, sub.Avatar, sub.RoomId, "", userlist)
						models.NewArchive(evt)
					}
				}
				Subscribers.PushBack(sub)
				userlist := GetUserlist(Subscribers, sub.RoomId)
				//beego.Info("userlist userlist:", userlist)
				evt := NewEvent(models.EVENT_JOIN, sub.Name, sub.Id, sub.Avatar, sub.RoomId, "", userlist)
				BroadcastWebSocket(evt)
				models.NewArchive(evt)
				//beego.Info("Old user:", sub.Id, ";WebSocket:", sub.Conn != nil)
			}
			//监听发布信息chan有数据
		case msg := <-Publish:
			// 通知等待队列。
			for ch := WaitingList.Back(); ch != nil; ch = ch.Prev() {
				ch.Value.(chan bool) <- true
				WaitingList.Remove(ch)
			}
			event := msg.Evt
			switch event.Type {
			case "EVENT_JOIN":
				beego.Debug("Join(event.UserId, event.UserName, event.RoomId, msg.Ws)", event.UserId, event.UserName, event.RoomId, msg.Evt)
				Join(event.UserId, event.UserName, event.RoomId, msg.Ws)

			case "EVENT_MESSAGE":
				// Send it out to every client that is currently connected
				BroadcastWebSocket(event)
				models.NewArchive(event)
			case "EVENT_LEAVE":
				beego.Debug("Join(event.UserId, event.UserName, event.RoomId, msg.Ws)", event.UserId, event.UserName, event.RoomId, msg.Evt)
				Leave(event.UserId, event.UserName, event.RoomId, msg.Ws)

			default:
			}

			//监听退出用户频道chan有数据
		case unsub := <-UnSubscribe:
			for sub := Subscribers.Front(); sub != nil; sub = sub.Next() {
				if sub.Value.(Subscriber).Id == unsub.Id {
					beego.Debug("退出:", unsub)
					Subscribers.Remove(sub)
					userlist := GetUserlist(Subscribers, unsub.RoomId)
					evt := NewEvent(models.EVENT_LEAVE, unsub.Name, unsub.Id, unsub.Avatar, unsub.RoomId, "", userlist)
					//					Publish <- Message{ws, }
					BroadcastWebSocket(evt) // Publish a LEAVE event.
					models.NewArchive(evt)
					break
				}
			}
		}
	}
}

func isUserExist(Subscribers *list.List, userid int) bool {
	for sub := Subscribers.Front(); sub != nil; sub = sub.Next() {
		if sub.Value.(Subscriber).Id == userid {
			return true
		}
	}
	return false
}

//GetUserlist 获取该房间用户列表
func GetUserlist(Subscribers *list.List, roomid string) []interface{} {
	userlist := make([]interface{}, 0)
	for sub := Subscribers.Front(); sub != nil; sub = sub.Next() {
		if sub.Value.(Subscriber).RoomId == roomid {
			userlist = append(userlist, sub.Value.(Subscriber))
		}
	}
	return userlist
}

// BroadcastWebSocket 向WebSocket用户广播消息,只针对本房间.
func BroadcastWebSocket(event models.WsEvent) {
	data, err := json.Marshal(event)
	if err != nil {
		beego.Error("Fail to marshal event:", err)
		return
	}
	for sub := Subscribers.Front(); sub != nil; sub = sub.Next() {
		//立即将事件发送到WebSocket用户。
		ws := sub.Value.(Subscriber).Conn
		if ws != nil && event.RoomId == sub.Value.(Subscriber).RoomId {
			if ws.WriteMessage(websocket.TextMessage, data) != nil {
				//用户断开连接。
				UnSubscribe <- sub.Value.(Subscriber)
			}
		}
	}
}

// BroadcastUser 向WebSocket特定用户userid广播消息.
func BroadcastUser(event models.WsEvent) {
	data, err := json.Marshal(event)
	if err != nil {
		beego.Error("Fail to marshal event:", err)
		return
	}
	for sub := Subscribers.Front(); sub != nil; sub = sub.Next() {
		//立即将事件发送到WebSocket用户。
		ws := sub.Value.(Subscriber).Conn
		if ws != nil && event.UserId == sub.Value.(Subscriber).Id {
			if ws.WriteMessage(websocket.TextMessage, data) != nil {
				//用户断开连接。
				UnSubscribe <- sub.Value.(Subscriber)
			}
		}
	}
}

// BroadcastAll 向WebSocket所有连接用户广播消息.
func BroadcastAll(event models.WsEvent) {
	data, err := json.Marshal(event)
	if err != nil {
		beego.Error("Fail to marshal event:", err)
		return
	}
	for sub := Subscribers.Front(); sub != nil; sub = sub.Next() {
		//立即将事件发送到WebSocket用户。
		ws := sub.Value.(Subscriber).Conn
		if ws != nil {
			if ws.WriteMessage(websocket.TextMessage, data) != nil {
				//清空用户。
				UnSubscribe <- sub.Value.(Subscriber)
			}
		}
	}
}
