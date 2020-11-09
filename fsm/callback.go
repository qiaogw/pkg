package fsm

// EventProcessor 定义OnExit，Action和OnEnter操作。
type EventProcessor interface {
	// OnExit 动作处理退出状态
	OnExit(fromState string, args []interface{})
	// Action 用于处理转换
	Action(action string, fromState string, toState string, args []interface{})
	// OnActionFailure 执行aciton 失败时 的处理
	OnActionFailure(action string, fromState string, toState string, args []interface{}, err error)
	// OnEnter Action 进入状态
	OnEnter(toState string, args []interface{})
}

// DefaultDelegate 是默认代理。
//将操作分为三个操作：OnExit，Action和OnEnter。
type DefaultDelegate struct {
	P EventProcessor
}

// HandleEvent 实现Delegate接口并将HandleEvent分为三个操作。
func (dd *DefaultDelegate) HandleEvent(action string, fromState string, toState string, args []interface{}) {
	// fmt.Println("action is ::", action)
	// fmt.Println("fromState is ::", fromState)
	// fmt.Println("toState is ::", toState)
	// fmt.Println("args is ::", args)
	if fromState != toState {
		dd.P.OnExit(fromState, args)
	}

	dd.P.Action(action, fromState, toState, args)

	if fromState != toState {
		dd.P.OnEnter(toState, args)
	}
}
