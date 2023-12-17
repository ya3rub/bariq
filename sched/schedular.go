package sched

import (
	"context"
	"fmt"
	"time"
)

var pl = fmt.Println

type (
	TaskFunc[T any] func(context.Context) (T, error)
	response[T any] struct {
		value T
		err   error
	}
)

type Task[T any] struct {
	respch       chan response[T]
	ctx          context.Context
	cancel       context.CancelFunc
	parentCancel context.CancelFunc
}

func (t *Task[T]) Await() (T, error) {
	select {
	case <-t.ctx.Done():
		var val T
		return val, t.ctx.Err()
	case resp := <-t.respch:
		return resp.value, resp.err
	}
}

type TaskOpts struct {
	Timeout time.Duration
}

func (t *Task[T]) Cancel() {
	t.cancel()
	t.parentCancel()
}

func SpawnWithTimeout[T any](
	t TaskFunc[T],
	d time.Duration,
) *Task[T] {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		d,
	)
	return spawn(ctx, cancel, t)
}

func Spawn[T any](t TaskFunc[T]) *Task[T] {
	ctx := context.Background()
	return spawn(ctx, func() {}, t)
}

func spawn[T any](
	ctx context.Context,
	parentCancel context.CancelFunc,
	t TaskFunc[T],
) *Task[T] {
	respch := make(chan response[T])
	// INFO: ctx here is the prarent context
	c, cancel := context.WithCancel(ctx)

	go func() {
		val, err := t(c)
		respch <- response[T]{
			value: val,
			err:   err,
		}
		close(respch)
	}()

	return &Task[T]{
		respch:       respch,
		ctx:          c,
		cancel:       cancel,
		parentCancel: parentCancel,
	}
}
