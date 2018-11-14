# go-cpulimit

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/Codehardt/go-cpulimit)](https://goreportcard.com/report/github.com/Codehardt/go-cpulimit)
[![GoDoc](https://godoc.org/github.com/Codehardt/go-cpulimit?status.svg)](https://godoc.org/github.com/Codehardt/go-cpulimit)

With **go-cpulimit** you can limit the CPU usage of your go program. It provides a function called `Wait()` that holds your program, until the CPU usage is below max CPU usage.

You should use this `Wait()` function in as many situations as possible, e.g. on every iteration in a `for loop`.

**Warning:** This limiter does not work, if there are insufficient `Wait()` calls.

This package requires [github.com/shirou/gopsutil](https://github.com/shirou/gopsutil) *(Copyright (c) 2014, WAKAYAMA Shirou)* to be installed.

## Example

```golang
limiter := &cpulimit.Limiter{
    MaxCPUUsage:     50.0,                   // throttle CPU usage to 50%
    MeasureInterval: time.Millisecond * 333, // measure cpu usage in an interval of 333 ms
    Measurements:    3,                      // use the avg of the last 3 measurements
}
limiter.Start()
defer limiter.Stop()
for {
    limiter.Wait() // wait until cpu usage is below 50%
    /*
     * do some work
     */
}
```
