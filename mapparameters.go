// Copyright 2015 Mathieu MAST. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
package goloc

type MapParameters map[string]interface{}

type KeyValue struct {
	Key   string
	Value interface{}
}

func NewMapParameters(keyValues ...KeyValue) *MapParameters {
	mp := make(MapParameters)
	for _, keyValue := range keyValues {
		mp[keyValue.Key] = keyValue.Value
	}
	return &mp
}

func (mp *MapParameters) Get(key string) interface{} {
	return (*mp)[key]
}

func (mp *MapParameters) Set(key string, value interface{}) {
	(*mp)[key] = value
}
