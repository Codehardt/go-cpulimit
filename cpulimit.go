package cpulimit

import (
	"sync"
	"time"

	"github.com/shirou/gopsutil/cpu"
)

const defaultLimit float64 = 80.0

type Limiter struct {
	MaxCPUUsage float64
	stop        bool
	wg          *sync.WaitGroup
	mutex       *sync.RWMutex
}

func (l *Limiter) Start() {
	if l.MaxCPUUsage == 0.0 {
		l.MaxCPUUsage = defaultLimit
	}
	l.wg = &sync.WaitGroup{}
	l.mutex = &sync.RWMutex{}
	l.wg.Add(1)
	go l.run()
}

func (l *Limiter) Stop() {
	l.stop = true
	l.wg.Wait()
}

func (l *Limiter) Wait() {
	l.mutex.RLock()
	l.mutex.RUnlock()
}

func (l *Limiter) run() {
	defer l.wg.Done()
	var (
		busy1    float64
		busy2    float64
		all1     float64
		all2     float64
		cpuUsage float64
		locked   bool
		m        = make([]float64, 4) // measurements
	)
	tk := time.NewTicker(time.Millisecond * 500)
	defer tk.Stop()
	busy2, all2 = getBusy()
	var counter int
	for range tk.C {
		if l.stop {
			if locked {
				l.mutex.Unlock()
				locked = false
			}
			break
		}
		busy1, all1 = busy2, all2
		busy2, all2 = getBusy()
		cpuUsage = getCPUUsage(busy1, all1, busy2, all2)
		m[counter] = cpuUsage
		if (m[0]+m[1]+m[2]+m[3])/4.0 > l.MaxCPUUsage {
			if !locked {
				l.mutex.Lock()
				locked = true
			}
		} else {
			if locked {
				l.mutex.Unlock()
				locked = false
			}
		}
		counter++
		if counter > 3 {
			counter = 0
		}
	}
}

func getBusy() (busy, all float64) {
	ts, err := cpu.Times(false)
	if err != nil {
		panic(err)
	}
	t := ts[0]
	busy = t.User + t.System + t.Nice + t.Iowait + t.Irq + t.Softirq + t.Steal + t.Guest + t.GuestNice + t.Stolen
	all = busy + t.Idle
	return
}

func getCPUUsage(busy1, all1, busy2, all2 float64) float64 {
	if all1 == all2 {
		return 0.0
	}
	usage := ((busy2 - busy1) / (all2 - all1)) * 100.0
	return usage
}
