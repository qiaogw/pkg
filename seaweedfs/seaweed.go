package seaweedfs

//文件系统操作
import (
	"github.com/astaxie/beego/httplib"

	"github.com/astaxie/beego"
	//	"io/ioutil"
)

var (
	msvr = "http://" + beego.AppConfig.String("seaweedfsAddress") + ":9333/dir/assign"
	vsvr = "http://" + beego.AppConfig.String("seaweedfsAddress") + ":8080/"
)

//SetFile 文件写入分布式系统，
//输入文件路径，输出文件fid，文件大小，错误
func SetFile(file string) (fid string, size int, err error) {

	req := httplib.Get(msvr)
	_, err = req.Response()
	if err != nil {
		beego.Error(err)
		return
	}
	result := new(map[string]interface{})
	req.ToJSON(&result)
	fid = (*result)["fid"].(string)

	url := vsvr + fid
	req = httplib.Post(url)
	req.PostFile("uploadfile1", file)
	//	req.Body(bt)
	_, err = req.Response()
	if err != nil {
		beego.Error(err)
		return
	}
	req.ToJSON(&result)
	siz, _ := (*result)["size"].(float64)
	size = int(siz)
	return

}

//DeleteFile 文件删除分布式系统，
//输入文件路径，输出错误
func DeleteFile(fid string) (err error) {

	url := vsvr + fid
	req := httplib.Delete(url)
	_, err = req.Response()
	if err != nil {
		beego.Error(err)
		return
	}
	return

}

//UbdateFile 文件更新分布式系统，
//输入文件路径，输出错误
func UbdateFile(fid, file string) (err error) {

	url := vsvr + fid
	req := httplib.Post(url)

	req.PostFile("uploadfile1", file)
	//	req.Body(bt)
	_, err = req.Response()
	if err != nil {
		beego.Error(err)
		return
	}
	return

}

//InitJpg 初始化加载初始化页面
func InitJpg() {

	url := vsvr + "times.jpg"
	req := httplib.Post(url)
	req.PostFile("uploadfile1", "imgs/times.jpg")
	_, err := req.Response()
	if err != nil {
		beego.Error(err)
		return
	}
	url = vsvr + "default.jpg"
	req = httplib.Post(url)
	req.PostFile("uploadfile1", "imgs/default.jpg")
	_, err = req.Response()
	if err != nil {
		beego.Error(err)
		return
	}
	return
}
