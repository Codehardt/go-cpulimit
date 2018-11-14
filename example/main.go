package main

import (
	"time"

	"github.com/Codehardt/go-cpulimit"
)

func main() {
	limiter := &cpulimit.Limiter{
		MaxCPUUsage:     50.0,                   // throttle if current cpu usage is over 50%
		MeasureInterval: time.Millisecond * 333, // measure cpu usage in an interval of 333 milliseconds
		Measurements:    3,                      // use the average of the last 3 measurements for cpu usage calculation
	}
	limiter.Start()
	defer limiter.Stop()
	for {
		limiter.Wait() // wait until cpu usage is below 50%
		/*
		 * do some work
		 */
	}
}
