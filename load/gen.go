package main

import (
	"errors"
	"demo/load/lib"
	"time"
	"ucenter/utils"
)

type myGenerator struct {
	caller      lib.Caller           // 调用器。
	timeoutNs   time.Duration        // 处理超时时间，单位：纳秒。
	lps         uint32               // 每秒载荷量。
	durationNs  time.Duration        // 负载持续时间，单位：纳秒。
	concurrency uint32               // 并发量。
	tickets     lib.GoTickets        // Goroutine票池。
	stopSign    chan byte            // 停止信号的传递通道。
	cancelSign  byte                 // 取消发送后续结果的信号。
	endSign     chan uint64          // 完结信号的传递通道，同时被用于传递调用执行计数。
	callCount   uint64               // 调用执行计数。
	status      lib.GenStatus        // 状态。
	resultCh    chan *lib.CallResult // 调用结果通道。
}

func NewGenerator(
	caller lib.Caller,
	timeoutNs time.Duration,
	lps uint32,
	durationNs time.Duration,
	resultCh chan *lib.CallResult) (lib.Generator, error) {
	utils.CLog("[INFO] 初始化载荷发生器...")
	utils.CLog("[INFO] 检查参数...")

	var errMsg string
	if caller == nil {
		errMsg = "Invalid caller!"
	}
	if timeoutNs == 0 {
		errMsg = "Invalid timeoutNs!"
	}
	if lps == 0 {
		errMsg = "Invalid lps(load per second)!"
	}
	if durationNs == 0 {
		errMsg = "Invalid durationNs!"
	}
	if resultCh == nil {
		errMsg = "Invalid result channel!"
	}
	if errMsg != "" {
		utils.CLog("[TRAC] 发现参数错误...")
		utils.CLog("[TRAC] 发现参数错误...")
		return nil, errors.New(errMsg)
	}
	gen := &myGenerator{
		caller:     caller,
		timeoutNs:  timeoutNs,
		lps:        lps,
		durationNs: durationNs,
		stopSign:   make(chan byte, 1),
		cancelSign: 0,
		status:     lib.STATUS_ORIGINAL,
		resultCh:   resultCh,
	}
	utils.CLog("[INFO] Passed. (timeoutNs=%v, lps=%d, durationNs=%v)", timeoutNs, lps, durationNs)
	err := gen.init()
	if err != nil {
		return nil, err
	}
	return gen, nil
}

func (gen *myGenerator) init() error {
	return nil
}

func (gen *myGenerator) genLoad(throttle <-chan time.Time, endSign chan<- uint64) {
	callCount := uint64(0)
Loop:
	for ; ; callCount++ {
		select {
		case <-gen.stopSign:
			gen.handleStopSign()
			endSign <- callCount
			break Loop
		default:
		}
		gen.asyncCall()
		if gen.lps > 0 {
			select {
			case <-throttle:
			case <-gen.stopSign:
				gen.handleStopSign()
				endSign <- callCount
				break Loop
			}
		}
	}
}

func (gen *myGenerator) handleStopSign() {
	gen.cancelSign = 1
	utils.CLog("[INFO] 取消结果通道...")
	close(gen.resultCh)
}

func (gen *myGenerator) asyncCall() {
	gen.tickets.Take()
	go func() {

		rawReq := gen.caller.BuildReq()
		gen.tickets.Return()
	}()
}

func (gen *myGenerator) Start() {
	utils.CLog("[INFO] 启动荷载发生器...")
	// 设置节流阀
	var throttle <-chan time.Time
	if gen.lps > 0 {
		interval := time.Duration(1e9 / gen.lps)
		utils.CLog("[INFO] 设置发送间隔为...")
		throttle = time.Tick(interval)
	}

	// 初始化停止信号
	go func() {
		time.AfterFunc(gen.durationNs, func() {
			utils.CLog("[INFO] 停止荷载发生器...")
			gen.stopSign <- 0
		})
	}()

	//设置已启动状态
	gen.status = lib.STATUS_STARTED
}

func (gen *myGenerator) Stop() (uint64, bool) {
	return 0, true
}

func (gen *myGenerator) Status() lib.GenStatus {
	return gen.status
}
