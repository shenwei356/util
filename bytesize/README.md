bytesize
========

package for provides a way to show readable values of byte sizes by reediting the code from http://golang.org/doc/effective_go.html

Install
-------
This package is "go-gettable", just:

    go get github.com/shenwei356/util/bytesize

Example
-------
    
    fmt.Sprintlf("1024 bytes = %v\n", ByteSize(float64(1024)))

Result:

    1024 bytes = 1.00 KB

Copyright (c) 2013, Wei Shen (shenwei356@gmail.com)

[MIT License](https://github.com/shenwei356/util/bytesize/master/LICENSE)