package future

import "time"

type Future[R any] struct {
	err     error
	result  R
	finsh   chan struct{}
	timeout time.Duration

	dofunc       func() (R, error)
	completeFunc func(R) error
}

func NewFuture[R any]() *Future[R] {
	return &Future[R]{
		finsh: make(chan struct{}),
	}
}

func (fu *Future[R]) OnComplete(c func(R) error) *Future[R] {
	fu.completeFunc = c
	return fu
}

func (fu *Future[R]) Timeout(t time.Duration) *Future[R] {
	fu.timeout = t
	return fu
}

// func (fu *Future[R]) Cancel() {

// }

func (fu *Future[R]) Exec(f func() (R, error)) (R, error) {
	go func() {
		defer func() {
			if e := recover(); e != nil {
				fu.err = e.(error)
			}
		}()
		r, err := f()
		if err != nil {
			fu.err = err
			return
		}
		fu.result = r

		if fu.completeFunc != nil {
			err = fu.completeFunc(r)
			if err != nil {
				fu.err = err
				return
			}
		}
		fu.err = nil
		fu.finsh <- struct{}{}
	}()

	<-fu.finsh
	return fu.result, fu.err
}
