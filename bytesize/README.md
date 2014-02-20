bytesize
========

Package for providing a way to show readable values of byte sizes by reediting
the code from http://golang.org/doc/effective_go.html. It could also parsing
byte size text to ByteSize object.

Install
-------
This package is "go-gettable", just:

    go get github.com/shenwei356/util/bytesize

Example
-------
    
	fmt.Printf("1024 bytes\t%v\n", bytesize.ByteSize(float64(1024)))
	fmt.Printf("13146111 bytes\t%v\n", bytesize.ByteSize(float64(13146111)))

    // parsing
	size, err := bytesize.Parse([]byte("1.5 KB"))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%.0f bytes\n", size)


Result:

    1024 bytes         1.00 KB
    13146111 bytes    12.54 MB
    1536 bytes



Copyright (c) 2013, Wei Shen (shenwei356@gmail.com)

Documentation
-------------

[See documentation on gowalker for more detail](http://gowalker.org/github.com/shenwei356/util/bytesize).

[MIT License](https://github.com/shenwei356/util/blob/master/bytesize/LICENSE)