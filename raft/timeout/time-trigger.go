package timeout

import (
	"time"
)

// TimerTrigger 结构体封装了定时器和控制逻辑
type TimerTrigger struct {
	timer     *time.Timer
	resetChan chan bool
	stopChan  chan bool
	duration  time.Duration
	CallBack  func()
}

// NewTimerControl 创建并返回一个新的 TimerControl 实例
func NewTimerControl(duration time.Duration) *TimerTrigger {
	return &TimerTrigger{
		timer:     time.NewTimer(duration),
		resetChan: make(chan bool),
		stopChan:  make(chan bool),
		duration:  duration,
	}
}

// Start 启动定时器监听循环
func (tc *TimerTrigger) Start(callback func()) {
	go func() {
		for {
			select {
			case <-tc.timer.C:
				go callback()
				return
			case <-tc.resetChan:
				if !tc.timer.Stop() {
					<-tc.timer.C
				}
				tc.timer.Reset(tc.duration)
			case <-tc.stopChan:
				tc.timer.Stop()
				return
			}
		}
	}()
}

// Start 启动定时器监听循环
func (tc *TimerTrigger) StartIntervalTask(callback func()) {
	go func() {
		for {
			select {
			case <-tc.timer.C:
				go callback()
				tc.timer.Reset(tc.duration)
			case <-tc.resetChan:
				if !tc.timer.Stop() {
					<-tc.timer.C
				}
				tc.timer.Reset(tc.duration)
			case <-tc.stopChan:
				tc.timer.Stop()
				return
			}
		}
	}()
}

// Reset 重置定时器
func (tc *TimerTrigger) Reset() {
	tc.resetChan <- true
}

// Stop 停止定时器
func (tc *TimerTrigger) Stop() {
	select {
	case tc.stopChan <- true:
	default:
	}
}
