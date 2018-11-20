package svm

import (
	"math"
)

type luaTable struct {
	array   []luaValue
	hashmap map[luaValue]luaValue
}

func newLuaTable(arrayLen, mapCap int) *luaTable {
	t := &luaTable{}
	if arrayLen > 0 {
		t.array = make([]luaValue, arrayLen)
	}
	if mapCap > 0 {
		t.hashmap = make(map[luaValue]luaValue, mapCap)
	} else {
		t.hashmap = make(map[luaValue]luaValue)
	}
	return t
}

func (t *luaTable) get(key luaValue) luaValue {
	key = keyFloatToInteger(key)
	if idx, ok := key.(int64); ok {
		if idx >= 1 && idx <= int64(len(t.array)) {
			return t.array[idx-1]
		}
	}
	return t.hashmap[key]
}

func (t *luaTable) set(key, value luaValue) {
	switch idx := key.(type) {
	case nil:
		return
	case float64:
		if math.IsNaN(idx) {
			return
		}
		if i, ok := floatToInteger(idx); ok && float64(i) == idx {
			//set int key with value
			t.setInt(i, value)
		} else if value == nil {
			delete(t.hashmap, key)
		} else {
			t.hashmap[key] = value
		}
	case int64:
		//set int kty with value
		t.setInt(idx, value)
	default:
		if value == nil {
			delete(t.hashmap, key)
		} else {
			t.hashmap[key] = value
		}
	}
}

func (t *luaTable) len() int {
	return len(t.array)
}

func (t *luaTable) setInt(key int64, value luaValue) {
	k := int(key) - 1
	kok := int64(int(k)) == (key - 1)
	if kok && k >= 0 && k < len(t.array) {
		t.array[k] = value
		if k+1 == len(t.array) && value == nil {
			t.removeNilItems()
		}
	} else if kok && k >= 0 && k == len(t.array) {
		delete(t.hashmap, key)
		if value != nil {
			t.array = append(t.array, value)
			t.expandArray()
		}
	} else {
		t.hashmap[key] = value
	}
}

func (t *luaTable) removeNilItems() {
	for i := len(t.array) - 1; i >= 0; i-- {
		if t.array[i] == nil {
			t.array = t.array[0:i]
		}
	}
}

func (t *luaTable) expandArray() {
	for i := len(t.array); i >= 0; i++ {
		if val, find := t.hashmap[i]; find {
			delete(t.hashmap, i)
			t.array = append(t.array, val)
		} else {
			break
		}
	}
}
