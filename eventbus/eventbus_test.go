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
	eBus.RegisterHandler(reflect.TypeOf(testEvent{}), func(event any) {
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
	eBus.RegisterHandler(reflect.TypeOf(testEvent{}), func(any) {
		et.Increment()
	})
	eBus.RegisterHandler(reflect.TypeOf(testEvent{}), func(any) {
		et.Increment()
	})
	eBus.RegisterHandler(reflect.TypeOf(testEvent{}), func(any) {
		et.Increment()
	})
	eBus.RegisterHandler(reflect.TypeOf(testEvent{}), func(any) {
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

func (et *testEventTarget) Add(value int) {
	et.count += value
}

type testEventWithData struct {
	value1 int
	value2 int
}

func TestEventBus_Handle_WithData(t *testing.T) {
	et := testEventTarget{count: 99}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	e := testEventWithData{
		value1: 5,
		value2: 10,
	}
	var wg sync.WaitGroup
	wg.Add(1)
	eBus := New(log.Default())
	go eBus.Run(ctx, &wg)

	eBus.RegisterHandler(reflect.TypeOf(testEventWithData{}), func(e any) {
		eventData, ok := e.(testEventWithData)
		assert.True(t, ok)
		et.Add(eventData.value1)
		et.Add(eventData.value2)
	})
	eBus.Push(e)
	time.Sleep(1 * time.Second)
	assert.Equal(t, 114, et.count)
	cancel()
	wg.Wait()
}
