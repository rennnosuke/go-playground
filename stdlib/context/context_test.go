package context

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestContext_WithValue(t *testing.T) {
	c := context.Background()

	// with value - set value
	key := "key_hoge"
	val := "value_hoge"
	ctx := context.WithValue(c, key, val)
	fmt.Println(ctx.Value(key)) // get value
}

func TestContext_WithCancel(t *testing.T) {
	c := context.Background()

	// with cancel - get child context and cancelFunc.
	ctx, cancelFunc := context.WithCancel(c)

	// async
	fmt.Println("start async function call.")
	go func(ctx context.Context) {
		select {
		case <-ctx.Done():
			fmt.Println("context was canceled.")
		}
		fmt.Println("end of async function call.")
	}(ctx)

	time.Sleep(5 * time.Second)
	cancelFunc()
}

func TestContext_WithTimeout(t *testing.T) {
	c := context.Background()

	ctx, _ := context.WithTimeout(c, 5*time.Second)

	fmt.Println("start async function call.")
	go func(ctx context.Context) {
		select {
		case <-ctx.Done():
			fmt.Println("context was canceled(timeout).")
		}
		fmt.Println("end of async function call.")
	}(ctx)

	// check deadline..
	if dl, ok := ctx.Deadline(); ok {
		fmt.Printf("deadline: %v\n", dl)
	}

	time.Sleep(10 * time.Second)
}

func TestContext_WithDeadline(t *testing.T) {
	c := context.Background()

	deadline := time.Now().Add(5 * time.Second) //  time is specified.
	// with deadline - set deadline
	ctx, _ := context.WithDeadline(c, deadline)

	fmt.Println("start async function call.")
	go func(ctx context.Context) {
		select {
		case <-ctx.Done():
			fmt.Println("context was canceled(deadline).")
		}
		fmt.Println("end of async function call.")
	}(ctx)

	// check deadline..
	if dl, ok := ctx.Deadline(); ok {
		fmt.Printf("deadline: %v\n", dl)
	}

	time.Sleep(10 * time.Second)
}
