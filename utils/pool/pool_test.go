package pool

import (
	"fmt"
	"testing"
)

func TestNewPoolWithSize(t *testing.T) {
	pool := NewPoolWithSize(1)
	pool.RunGo(func() {
		panic("this should never be executed")
	})
	pool.RunGo(func() {
		println("Hello World2")
		panic("hello world 2")
	})
	pool.Wait()
	fmt.Println("done all")
}
