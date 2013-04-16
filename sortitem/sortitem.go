// Copyright 2013 Wei Shen (shenwei356@gmail.com). All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

// Package sortitem provides sort methods to key-value list.
package sortitem

import (
	"sort"
)

// Imitate code from http://golang.org/pkg/sort/#example_Interface .
type Item struct {
	Key   string
	Value float64
}

type ItemList []Item

func (list ItemList) Swap(i, j int) {
	list[i], list[j] = list[j], list[i]
}
func (list ItemList) Len() int {
	return len(list)
}

// Sort alphabetically by Key 
type ByKey struct {
	ItemList
}

func (list ByKey) Less(i, j int) bool {
	return list.ItemList[i].Key < list.ItemList[j].Key
}

// Ascending sort by Value
type ByValue struct {
	ItemList
}

func (list ByValue) Less(i, j int) bool {
	return list.ItemList[i].Value < list.ItemList[j].Value
}

// Reverse embeds a sort.Interface value and implements a reverse sort over
// that value.
type Reverse struct {
	sort.Interface
}

func (r Reverse) Less(i, j int) bool {
	return r.Interface.Less(j, i)
}
