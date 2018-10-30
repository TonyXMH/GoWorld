package timer

import (
	"container/heap"
	"runtime/debug"
	"sync"
	"time"
	"github.com/TonyXMH/GoWorld/gwlog"
)

const MIN_TIMER_INTERVAL = 1 * time.Millisecond

var nextAddSeq uint = 1

type CallbackFunc func()

type Timer struct {
	fireTime time.Time
	interval time.Duration
	callback CallbackFunc
	repeat   bool
	addseq   uint
}

func (t *Timer) Cancel() {
	t.callback = nil
}

func (t *Timer) IsActive() bool {
	return t.callback != nil
}

type _TimerHeap struct {
	timers []*Timer
}

func (h *_TimerHeap) Len() int {
	return len(h.timers)
}

func (h *_TimerHeap) Less(i, j int) bool {
	t1, t2 := h.timers[i].fireTime, h.timers[j].fireTime
	if t1.Before(t2) {
		return true
	}
	if t1.After(t2) {
		return false
	}
	return h.timers[i].addseq < h.timers[j].addseq
}

func (h *_TimerHeap) Swap(i, j int) {
	h.timers[i], h.timers[j] = h.timers[j], h.timers[i]
}

func (h *_TimerHeap) Push(x interface{}) {
	h.timers = append(h.timers, x.(*Timer))
}

func (h *_TimerHeap) Pop() (ret interface{}) {
	l := len(h.timers)
	ret, h.timers = h.timers[l-1], h.timers[:l-1]
	return ret
}

var (
	timerHeap     _TimerHeap
	timerHeapLock sync.Mutex
)

func init() {
	heap.Init(&timerHeap)
}

func AddCallback(d time.Duration, callback CallbackFunc) *Timer {
	t := &Timer{
		fireTime: time.Now().Add(d),
		interval: d,
		callback: callback,
		repeat:   false,
	}
	timerHeapLock.Lock()
	t.addseq = nextAddSeq
	nextAddSeq += 1
	heap.Push(&timerHeap, t)
	timerHeapLock.Unlock()
	return t
}

func AddTimer(d time.Duration, callback CallbackFunc) *Timer {
	if d < MIN_TIMER_INTERVAL {
		d = MIN_TIMER_INTERVAL
	}
	t := &Timer{
		fireTime: time.Now().Add(d),
		interval: d,
		callback: callback,
		repeat:   true,
	}
	timerHeapLock.Lock()
	t.addseq = nextAddSeq
	nextAddSeq++
	heap.Push(&timerHeap, t)
	timerHeapLock.Unlock()
	return t
}

func Tick() {
	now := time.Now()
	timerHeapLock.Lock()
	for {
		if timerHeap.Len() <= 0 {
			break
		}
		nextFireTime := timerHeap.timers[0].fireTime
		if nextFireTime.After(now) {
			break
		}
		t := heap.Pop(&timerHeap).(*Timer)
		callback := t.callback
		if callback == nil {
			continue
		}
		if !t.repeat {
			t.callback = nil
		}
		timerHeapLock.Unlock()
		runCallback(callback)
		timerHeapLock.Lock()
		if t.repeat {
			t.fireTime = t.fireTime.Add(t.interval)
			if !t.fireTime.After(now) {
				t.fireTime = now.Add(t.interval)
			}
			t.addseq = nextAddSeq
			nextAddSeq += 1
			heap.Push(&timerHeap, t)
		}
	}
	timerHeapLock.Unlock()
}

func StartTicks(tickInterval time.Duration) {
	go selfTickRoutine(tickInterval)
}

func selfTickRoutine(tickInterval time.Duration) {
	for {
		time.Sleep(tickInterval)
		Tick()
	}
}

func runCallback(callback CallbackFunc) {
	defer func() {
		err := recover()
		if err != nil {
			//fmt.Fprintln(os.Stderr, "Callback%v paniced:%v\n", callback, err)
			gwlog.Error("Callback%v paniced:%v\n", callback, err)
			debug.PrintStack()
		}
	}()
	callback()
}
