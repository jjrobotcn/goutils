package observer

import (
	"sync"
)

type IObserver interface {
	Register(c chan<- interface{})
	Unregister(c chan<- interface{})
	Update(data interface{})
}

type observerImpl struct {
	observers *sync.Map
}

func NewObserver() IObserver {
	return &observerImpl{
		observers: new(sync.Map),
	}
}

func (o *observerImpl) Register(c chan<- interface{}) {
	o.observers.Store(c, c)
}

func (o *observerImpl) Unregister(c chan<- interface{}) {
	defer func() {
		recover()
	}()
	o.observers.Delete(c)
	close(c)
}

func (o *observerImpl) Update(data interface{}) {
	o.observers.Range(func(key, value interface{}) bool {
		c := value.(chan<- interface{})
		select {
		case c <- data:
		default:
			o.Unregister(c)
		}
		return true
	})
}
