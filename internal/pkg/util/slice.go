// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package util

func SliceFilter[T interface{}](data []T, filterFunc func(T) bool) []T {
	var filterData []T

	for _, v := range data {
		if filterFunc(v) {
			filterData = append(filterData, v)
		}
	}
	return filterData
}

func SliceFind[T interface{}](data []*T, someFunc func(*T) bool) *T {

	for _, v := range data {
		if someFunc(v) {
			return v
		}
	}

	return nil
}
