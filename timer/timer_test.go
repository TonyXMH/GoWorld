package timer

import (
	"github.com/TonyXMH/GoWorld/gwlog"
	//"net/http/pprof"
	"os"
	"testing"
	"time"
	"math/rand"
	"runtime/pprof"
)

func init() {
	StartTicks(time.Millisecond)
}

func TestCallback(t *testing.T) {
	INTERVAL := 100 * time.Millisecond
	for i := 0; i < 10; i++ {
		x := false
		AddCallback(INTERVAL, func() {
			gwlog.Info("callback!")
			x = true
		})
		time.Sleep(INTERVAL * 2)
		if !x {
			t.Fatalf("x should be true,but it's false")
		}
	}
}

func TestTimer(t *testing.T) {
	INTERVAL := 100 * time.Millisecond
	x := 0
	px := x
	now := time.Now()
	nextTime := now.Add(INTERVAL)
	gwlog.Info("now is %s, next time should be%s\n", time.Now(), nextTime)
	AddTimer(INTERVAL, func() {
		x += 1
		gwlog.Info("timer %s x %v px %v\n", time.Now(), x, px)
	})

	for i := 0; i < 10; i++ {
		time.Sleep(nextTime.Add(INTERVAL / 2).Sub(time.Now()))
		gwlog.Info("Check x %v px %v @ %s\n", x, px, time.Now())
		if x != px+1 {
			t.Fatalf("x should be %d, but it's %d", px+1, x)
		}
		px = x
		nextTime = nextTime.Add(INTERVAL)
		gwlog.Info("now is %s, next time should be %s\n", time.Now(), nextTime)

	}
}

func TestCallbackSeq(t *testing.T) {
	a := 0
	d := time.Second
	for i := 0; i < 100; i++ {
		j := i
		AddCallback(d, func() {
			if a != j {
				t.Error(j, a)
			}
			a += 1
		})
	}
	time.Sleep(d + time.Second*1)
}

func TestCancelCallback(t *testing.T) {
	INTERVAL := 20 * time.Millisecond
	x := 0
	timer := AddCallback(INTERVAL, func() {
		x = 1
	})
	if !timer.IsActive() {
		t.Fatalf("time should be active")
	}
	timer.Cancel()
	if timer.IsActive() {
		t.Fatalf("timer should be inactive")
	}
	time.Sleep(INTERVAL * 2)
	if x != 0 {
		t.Fatalf("x should be 0,but is %v", x)
	}
}

func TestCancelTimer(t *testing.T) {
	INTERVAL := 20 * time.Millisecond
	x := 0
	timer := AddTimer(INTERVAL, func() {
		x += 1
	})
	if !timer.IsActive() {
		t.Fatalf("timer should be active")
	}
	timer.Cancel()
	if timer.IsActive() {
		t.Fatalf("time should be inactive")
	}
	time.Sleep(INTERVAL * 2)
	if x != 0 {
		t.Fatalf("x should be 0,but is %v", x)
	}
}

func TestTimerPerformance(t *testing.T) {
	f, err := os.Create("TestTimerPerformance.cpuprof")
	if err != nil {
		panic(err)
	}
	pprof.StartCPUProfile(f)
	duration := 10 * time.Second
	for i:=0;i<400000 ;i++  {
		if rand.Float32()<0.5{
			d:=time.Duration(rand.Int63n(int64(duration)))
			AddCallback(d, func() {

			})
		}else{
			d:=time.Duration(rand.Int63n(int64(time.Second)))
			AddTimer(d, func() {})
		}
	}
	gwlog.Info("Waiting for",duration,"...")
	time.Sleep(duration)
	pprof.StopCPUProfile()
}
