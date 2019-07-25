package main

import (
	"context"
	"fmt"
	"github.com/jjrobotcn/goutils/observer"
	"time"
)

func main() {
	c := make(chan interface{}, 100)
	o := observer.NewObserver()
	o.Register(c)
	defer o.Unregister(c)

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
	defer cancel()

	type Test struct {
		YesOrNo bool
	}

	go func() {
		o.Update(123)

		o.Update("text")

		t1 := Test{YesOrNo: true}
		o.Update(t1)

		t2 := &Test{YesOrNo: false}
		o.Update(t2)
	}()

	for {
		select {
		case v, ok := <-c:
			if !ok {
				return
			}

			switch vv := v.(type) {
			case int:
				fmt.Println(vv)
				// 123
			case string:
				fmt.Println(vv)
				// text
			case Test:
				fmt.Println(vv)
				// {true}
			case *Test:
				fmt.Println(*vv)
				// {false}
			default:
				_ = vv
			}
		case <-ctx.Done():
			return
		}
	}
}
