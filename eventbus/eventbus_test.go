package eventbus

import (
	"context"
	"log"
	"reflect"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vidarandrebo/oncetree/consts"
)

func TestEventBus_Handle(t *testing.T) {
	et := testEventTarget{count: 99}
	ctx, cancel := context.WithTimeout(context.Background(), consts.RPCContextTimeout)
	e := testEvent{}
	var wg sync.WaitGroup
	wg.Add(1)
	eBus := New(log.Default())
	go eBus.Run(ctx, &wg)
	eBus.RegisterHandler(reflect.TypeOf(testEvent{}), func(event any) {
		et.Increment()
	})
	eBus.PushEvent(e)

	eBus.closeChannels()
	wg.Wait()
	cancel()
	assert.Equal(t, 100, et.count)
}

func TestEventBus_Handle_ManyHandlers(t *testing.T) {
	et := testEventTarget{count: 99}
	ctx, cancel := context.WithTimeout(context.Background(), consts.RPCContextTimeout)
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
	eBus.PushEvent(e)

	eBus.closeChannels()
	wg.Wait()
	cancel()
	assert.Equal(t, 101, et.count)
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
	ctx, cancel := context.WithTimeout(context.Background(), consts.RPCContextTimeout)
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
	eBus.PushEvent(e)

	eBus.closeChannels()
	wg.Wait()
	cancel()
	assert.Equal(t, 114, et.count)
}

func TestEventBus_HandleOneTask(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	wg.Add(1)
	eBus := New(log.Default())
	go eBus.Run(ctx, &wg)

	et := testEventTarget{count: 10}
	closure := func() {
		et.count++
	}
	eBus.PushTask(closure)

	eBus.closeChannels()
	wg.Wait()
	cancel()
	assert.Equal(t, 11, et.count)
}

func TestEventBus_HandleMultipleTasks(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	wg.Add(1)
	eBus := New(log.Default())
	go eBus.Run(ctx, &wg)

	var mut sync.Mutex
	et := testEventTarget{count: 10}
	closure := func() {
		// need mutex for this test, since there are multiple task handlers
		mut.Lock()
		et.count++
		mut.Unlock()
	}
	eBus.PushTask(closure)
	eBus.PushTask(closure)
	eBus.PushTask(closure)
	eBus.PushTask(closure)
	eBus.closeChannels()

	wg.Wait()
	cancel()
	assert.Equal(t, 14, et.count)
}

// closeChannels testing only fn, used to flush all events at test end
func (eb *EventBus) closeChannels() {
	eb.mut.Lock()
	close(eb.pendingEvents)
	close(eb.pendingTasks)
	eb.mut.Unlock()
}
