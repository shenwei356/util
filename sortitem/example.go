package main

import (
    . "github.com/shenwei356/util/sortitem"
    "fmt"
    "sort"
)

func main() {
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
}
