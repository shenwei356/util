ParallelMap
===========

A lock-free parallel map in go.

ParallelMap uses a backend goroutine for the sequential excution of 
Get and Set or custom function, which was inspired by section 14.17
in book *The Way to Go*.

Install
-------
This package is "go-gettable", just:

    go get github.com/shenwei356/util/ParallelMap

Example
-------
    
    import (
        "fmt"
        pmap "github.com/shenwei356/util/ParallelMap"
        "runtime"
    )

    func main() {
        runtime.GOMAXPROCS(runtime.NumCPU())

        m := pmap.NewParallelMap()
        m.Set("year", 2014)

        if v, ok := m.Get("year"); ok {
            fmt.Println(v)
        }
    }

Documentation
-------------

[See documentation on gowalker for more detail](http://gowalker.org/github.com/shenwei356/util/ParallelMap).

[MIT License](https://github.com/shenwei356/util/blob/master/ParallelMap/LICENSE)