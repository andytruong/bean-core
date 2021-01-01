package module

import (
	"context"
	"sort"
)

func NewHook() *Hook {
	return &Hook{
		listeners: []Callback{},
		orders:    map[int][]int{},
	}
}

type (
	Hook struct {
		listeners []Callback
		orders    map[int][]int
	}

	Callback func(ctx context.Context, event interface{}) error
)

func (hook *Hook) Hook(callback Callback, weight int) {
	hook.listeners = append(hook.listeners, callback)
	index := len(hook.listeners) - 1

	if _, found := hook.orders[weight]; !found {
		hook.orders[weight] = []int{}
	}

	hook.orders[weight] = append(hook.orders[weight], index)
}

func (hook Hook) Invoke(ctx context.Context, event interface{}) error {
	keys := []int{}
	for key := range hook.orders {
		keys = append(keys, key)
	}

	sort.Ints(keys)

	for _, key := range keys {
		for _, index := range hook.orders[key] {
			listener := hook.listeners[index]
			if err := listener(ctx, event); nil != err {
				return err
			}
		}
	}

	return nil
}
