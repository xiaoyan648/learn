package future

import (
	"fmt"
	"time"

	"github.com/fanliao/go-promise"
)

type Resp struct {
	code int
	name string
}

func ExampleMyFuture() {
	f := NewFuture[*Resp]()
	f = f.OnComplete(func(r *Resp) error {
		fmt.Printf("OnComplete result, code:%d, name:%s \n", r.code, r.name)
		return nil
	})
	r, err := f.Exec(func() (*Resp, error) {
		return &Resp{code: 200, name: "wuqinqiang"}, nil
	})
	fmt.Printf("r: %+v, err: %v", r, err)

	// Output:
	// OnComplete result, code:200, name:wuqinqiang
	// r: &{code:200 name:wuqinqiang}, err: <nil>
}

func ExampleGoPromise() {
	task := func() (r interface{}, err error) {
		time.Sleep(100 * time.Millisecond)
		return &Resp{code: 200, name: "wuqinqiang"}, nil
	}

	f := promise.Start(task).OnSuccess(func(v interface{}) {
		fmt.Printf("success: %v", v)
	}).OnFailure(func(v interface{}) {
		fmt.Printf("failed: %v", v)
	}).OnComplete(func(v interface{}) {
		fmt.Printf("complete: %v", v)
	})
	r, err := f.Get()
	fmt.Printf("r: %+v, err: %v", r, err)
	// Output:
}

func ExampleGoPromise2() {
	p := promise.NewPromise()
	p.OnSuccess(func(v interface{}) {
		fmt.Printf("success: %v", v)
	}).OnFailure(func(v interface{}) {
		fmt.Printf("failed: %v", v)
	}).OnComplete(func(v interface{}) {
		fmt.Printf("complete: %v", v)
	})

	go func() {
		time.Sleep(1000 * time.Millisecond)
		p.Resolve(&Resp{code: 200, name: "wuqinqiang"})
	}()
	r, err, isTimeout := p.GetOrTimeout(500)
	fmt.Printf("r: %+v, err: %v, isTimeout:%v", r, err, isTimeout)
	// Output:
}
