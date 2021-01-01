package module

import (
	"context"
	"sort"
)

func NewHook() *Hook {
	return &Hook{
		listeners: []Listener{},
		orders:    map[int][]int{},
	}
}

type (
	Hook struct {
		listeners []Listener
		orders    map[int][]int
	}

	Listener func(ctx context.Context, event interface{}) error
)

func (hook *Hook) Listen(listener Listener, weight int) {
	hook.listeners = append(hook.listeners, listener)
	index := len(hook.listeners) - 1

	if _, found := hook.orders[weight]; !found {
		hook.orders[weight] = []int{}
	}

	hook.orders[weight] = append(hook.orders[weight], index)
}

func (hook Hook) Trigger(ctx context.Context, event interface{}) error {
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
