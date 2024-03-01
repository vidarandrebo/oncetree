package eventbus

import (
	"context"
	"log"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEventBus_Handle(t *testing.T) {
	et := testEventTarget{count: 99}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	e := testEvent{}
	var wg sync.WaitGroup
	wg.Add(1)
	eBus := New(log.Default())
	go eBus.Run(ctx, &wg)
	eBus.RegisterHandler(reflect.TypeOf(testEvent{}), func() {
		et.Increment()
	})
	eBus.Push(e)
	time.Sleep(1 * time.Second)
	assert.Equal(t, 100, et.count)
	cancel()
	wg.Wait()
}

func TestEventBus_Handle_ManyHandlers(t *testing.T) {
	et := testEventTarget{count: 99}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	e := testEvent{}
	var wg sync.WaitGroup
	wg.Add(1)
	eBus := New(log.Default())
	go eBus.Run(ctx, &wg)
	eBus.RegisterHandler(reflect.TypeOf(testEvent{}), func() {
		et.Increment()
	})
	eBus.RegisterHandler(reflect.TypeOf(testEvent{}), func() {
		et.Increment()
	})
	eBus.RegisterHandler(reflect.TypeOf(testEvent{}), func() {
		et.Increment()
	})
	eBus.RegisterHandler(reflect.TypeOf(testEvent{}), func() {
		et.Decrement()
	})
	eBus.Push(e)
	time.Sleep(1 * time.Second)
	assert.Equal(t, 101, et.count)
	cancel()
	wg.Wait()
}

type testEvent struct{}

type testEventTarget struct {
	count int
}

func (et *testEventTarget) Increment() {
	et.count++
}

func (et *testEventTarget) Decrement() {
	et.count--
}
