/*
 *  Copyright (C) 2020-2021  AnySwap Ltd. All rights reserved.
 *  Copyright (C) 2020-2021  haijun.cai@anyswap.exchange
 *
 *  This library is free software; you can redistribute it and/or
 *  modify it under the Apache License, Version 2.0.
 *
 *  This library is distributed in the hope that it will be useful,
 *  but WITHOUT ANY WARRANTY; without even the implied warranty of
 *  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 *
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 *
 */

// Package common  Self encapsulated map structure supporting concurrent operation 
package common

import (
	"sync"
)

// SafeMap map + sync mutex
type SafeMap struct {
	sync.RWMutex
	Map map[string]interface{}
}

// NewSafeMap new SafeMap
func NewSafeMap(size int) *SafeMap {
	sm := new(SafeMap)
	sm.Map = make(map[string]interface{})
	return sm
}

// ReadMap get value by key
func (sm *SafeMap) ReadMap(key string) (interface{}, bool) {
	sm.RLock()
	value, ok := sm.Map[key]
	sm.RUnlock()
	return value, ok
}

// WriteMap write value by key
func (sm *SafeMap) WriteMap(key string, value interface{}) {
	sm.Lock()
	sm.Map[key] = value
	sm.Unlock()
}

// DeleteMap delete value by key
func (sm *SafeMap) DeleteMap(key string) {
	sm.Lock()
	delete(sm.Map, key)
	sm.Unlock()
}

// ListMap get all (key,value)
func (sm *SafeMap) ListMap() ([]string, []interface{}) {
	sm.RLock()
	l := len(sm.Map)
	key := make([]string, l)
	value := make([]interface{}, l)
	i := 0
	for k, v := range sm.Map {
		key[i] = k
		value[i] = v
		i++
	}
	sm.RUnlock()

	return key, value
}

// MapLength get len of map
func (sm *SafeMap) MapLength() int {
	sm.RLock()
	l := len(sm.Map)
	sm.RUnlock()
	return l
}
