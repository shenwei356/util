package main

import (
    . "github.com/shenwei356/util/bytesize"
    "fmt"
)

func main() {
    fmt.Printf("1024 bytes = %v\n", ByteSize(float64(1024)))
    fmt.Printf("13146111 bytes = %v\n", ByteSize(float64(13146111)))
}
