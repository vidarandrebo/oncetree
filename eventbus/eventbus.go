package eventbus

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"sync"
)

type EventBus struct {
	pendingEvents chan any
	handlers      map[reflect.Type][]func(any)
	mut           sync.RWMutex
	logger        *log.Logger
}

func New(logger *log.Logger) *EventBus {
	return &EventBus{
		handlers:      make(map[reflect.Type][]func(any)),
		pendingEvents: make(chan any, 100),
		logger:        logger,
	}
}

func (eb *EventBus) Run(cxt context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
loop:
	for {
		select {
		case event := <-eb.pendingEvents:
			eb.handle(event)
		case <-cxt.Done():
			break loop
		}
	}
}

func (eb *EventBus) handle(event any) {
	eventType := reflect.TypeOf(event)
	handlers, err := eb.eventHandlers(eventType)
	if err != nil {
		eb.logger.Println(err)
		return
	}
	for _, handler := range handlers {
		handler(event)
	}
}

func (eb *EventBus) Push(event any) {
	eb.pendingEvents <- event
}

func (eb *EventBus) eventHandlers(eventType reflect.Type) ([]func(any), error) {
	eb.mut.RLock()
	handlers, ok := eb.handlers[eventType]
	if !ok {
		return nil, fmt.Errorf("no handlers for event of type %v", eventType)
	}
	return handlers, nil
}

func (eb *EventBus) RegisterHandler(eventType reflect.Type, fn func(any)) {
	eb.mut.Lock()
	defer eb.mut.Unlock()
	if _, ok := eb.handlers[eventType]; !ok {
		eb.handlers[eventType] = make([]func(any), 0)
	}
	eb.handlers[eventType] = append(eb.handlers[eventType], fn)
}
