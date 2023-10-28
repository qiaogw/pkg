// Package fsm
// fsm 是一个简单的，功能强大的FSM实现，与其他FSM实现具有一些不同的功能。
// fsm 的一个特点是它不会保留/保留对象的状态。 当它处理转换时，您必须将当前状态传递给id，
// 因此您可以将 fsm 视为“无状态”状态机。 这个好处是一个 fsm 实例可以用来处理很多对象实例的转换，而不是创建大量的FSM实例。
// 对象实例自身保持状态。
// 另一个特点是它为Moore和Mealy FSM提供了一个通用接口。 您可以为这两个FSM实现相应的方法（OnExit，Action，OnEnter）。
// 第三个有趣的特性是您可以将配置的转换导出到状态图。 一张图片胜过千言万语。
// fsm 的样式执行
package fsm

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

// Transition 是状态转换，所有数据都是id，简化了FSM的使用，并使其成为通用的。
type Transition struct {
	FromName   string
	EventName  string
	ToName     string
	ActionName string
	From       string
	Event      string
	To         string
	Action     string
}

// Delegate 用于处理操作。 因为gofsm使用文字值作为事件，状态和动作，所以您需要使用相应的功能处理它们。
// DefaultDelegate是将处理分为三个操作的默认委托实现：OnExit Action，Action和OnEnter Action。 您可以实现不同的代理。
type Delegate interface {
	// HandleEvent handles transitions
	HandleEvent(action string, fromState string, toState string, args []interface{})
}

// StateMachine 可以处理很多对象的转换的FSM。 委托和转换在使用它们之前进行配置。
type StateMachine struct {
	delegate    Delegate
	transitions []Transition
}

// Error 是处理事件和状态改变时出错。
type Error interface {
	error
	BadEvent() string
	CurrentState() string
}

type smError struct {
	badEvent     string
	currentState string
}

// Error 获取
func (e smError) Error() string {
	return fmt.Sprintf("状态机错误:  事件 [%s] 不能发现转换规则当其处于状态 [%s]\n", e.badEvent, e.currentState)
}

// BadEvent 返回错误事件
func (e smError) BadEvent() string {
	return e.badEvent
}

// CurrentState 当前状态
func (e smError) CurrentState() string {
	return e.currentState
}

// NewStateMachine 创建一个新的状态机
func NewStateMachine(delegate Delegate, transitions ...Transition) *StateMachine {
	return &StateMachine{delegate: delegate, transitions: transitions}
}

// Trigger 发射一个事件 您必须传递处理对象的当前状态，关于此对象的其他信息可以通过args传递。
func (m *StateMachine) Trigger(currentState string, event string, args ...interface{}) Error {
	trans := m.findTransMatching(currentState, event)
	if trans == nil || len(trans) < 1 {
		return smError{event, currentState}
	}
	for _, t := range trans {
		if t.Action != "" {
			m.delegate.HandleEvent(t.Action, currentState, t.To, args)
		}
	}

	return nil
}

// findTransMatching 根据当前状态和事件获得相应的转换。
func (m *StateMachine) findTransMatching(fromState string, event string) (ts []Transition) {
	for _, v := range m.transitions {
		if v.From == fromState && v.Event == event {
			ts = append(ts, v)
		}
	}
	return
}

// Export 将状态图导出到文件中。
func (m *StateMachine) Export(outfile string) error {
	return m.ExportWithDetails(outfile, "png", "dot", "72", "-Gsize=10,5 -Gdpi=200")
}

// ExportWithDetails  使用更多graphviz选项导出状态图。
func (m *StateMachine) ExportWithDetails(outfile string, format string, layout string, scale string, more string) error {
	dot := `digraph StateMachine {
	rankdir=LR
	node[width=1 fixedsize=true shape=circle style=filled fillcolor="darkorchid1" ]
	
	`

	for _, t := range m.transitions {
		link := fmt.Sprintf(`%s -> %s [label="%s | %s"]`, t.FromName, t.ToName, t.EventName, t.ActionName)
		dot = dot + "\r\n" + link
	}

	dot = dot + "\r\n}"
	cmd := fmt.Sprintf("dot -o%s -T%s -K%s -s%s %s", outfile, format, layout, scale, more)

	return system(cmd, dot)
}

// system 系统命令
func system(c string, dot string) error {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command(`cmd`, `/C`, c)
	} else {
		cmd = exec.Command(`/bin/sh`, `-c`, c)
	}
	cmd.Stdin = strings.NewReader(dot)
	return cmd.Run()
}
