sortitem
========

Package sortitem provides sort methods to key-value list.
The code is imitated from http://golang.org/pkg/sort/#example_Interface .

Install
-------
This package is "go-gettable", just:

    go get github.com/shenwei356/util/sortitem

Example
-------
    
    list := make([]Item, 0)
    list = append(list, Item{"a", float64(12)})
    list = append(list, Item{"b", float64(3)})
    list = append(list, Item{"c", float64(45)})
    
    // Sort alphabetically by Key 
    sort.Sort(ByKey{list})
    fmt.Println(list)
    
    // Ascending sort by Value
    sort.Sort(ByValue{list})
    fmt.Println(list)
    
    // Reverse sort
    sort.Sort(Reverse{ByKey{list}})
    fmt.Println(list)

Result:

    [{a 12} {b 3} {c 45}]
    [{b 3} {a 12} {c 45}]
    [{c 45} {b 3} {a 12}]


Copyright (c) 2013, Wei Shen (shenwei356@gmail.com)

[MIT License](https://github.com/shenwei356/util/blob/master/sortitem/LICENSE)