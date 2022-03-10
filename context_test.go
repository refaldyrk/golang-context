package gocontext

import (
	"context"
	"fmt"
	"runtime"
	"testing"
	"time"
)

func TestCreateContext(t *testing.T) {
	bg := context.Background()
	fmt.Println(bg)

	td := context.TODO()
	fmt.Println(td)
}

func TestContextWithValue(t *testing.T) {
	bg := context.Background()

	childA := context.WithValue(bg, "ChildA", "A")
	childB := context.WithValue(bg, "ChildB", "B")

	childC := context.WithValue(childA, "ChildC", "C")
	childD := context.WithValue(childA, "ChildD", "D")

	childF := context.WithValue(childB, "ChildF", "F")
	childG := context.WithValue(childF, "ChildG", "G")

	fmt.Println(childA)
	fmt.Println(childB)
	fmt.Println(childC)
	fmt.Println(childD)
	fmt.Println(childF)
	fmt.Println(childG)

	fmt.Println(childA.Value("ChildA"))
	fmt.Println(childB.Value("ChildC"))
	fmt.Println(childC.Value("ChildD"))
	fmt.Println(childD.Value("ChildG"))
	fmt.Println(childF.Value("ChildB"))
}

func CreateCounter(ctx context.Context) chan int {
	dest := make(chan int)
	go func() {
		defer close(dest)
		i := 0
		for {
			select {
			case <-ctx.Done():
				return
			default:
				dest <- i
				i++
				time.Sleep(time.Second * 1)
			}
		}
	}()
	return dest
}

func TestGoLeak(t *testing.T) {
	fmt.Println("Num Goroutine:", runtime.NumGoroutine())

	bg := context.Background()
	ctx, cancel := context.WithCancel(bg)

	c := CreateCounter(ctx)
	for s := range c {
		fmt.Println(s)
		if s == 10 {
			break
		}
	}

	cancel()

	time.Sleep(time.Second * 5)

	fmt.Println("Num Goroutine:", runtime.NumGoroutine())
}

func TestContextWithTimeout(t *testing.T) {
	fmt.Println("Num Goroutine:", runtime.NumGoroutine())

	bg := context.Background()
	ctx, cancel := context.WithTimeout(bg, time.Second*5)
	defer cancel()

	c := CreateCounter(ctx)
	for s := range c {
		fmt.Println(s)
	}

	fmt.Println("Num Goroutine:", runtime.NumGoroutine())
}

func TestContextWithDeadline(t *testing.T) {
	fmt.Println("Num Goroutine:", runtime.NumGoroutine())

	bg := context.Background()
	ctx, cancel := context.WithDeadline(bg, time.Now().Add(time.Second*5))
	defer cancel()

	c := CreateCounter(ctx)
	for s := range c {
		fmt.Println(s)
	}

	fmt.Println("Num Goroutine:", runtime.NumGoroutine())

}
