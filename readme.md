# werk [![Build Status](https://travis-ci.org/bendrucker/werk.svg?branch=master)](https://travis-ci.org/bendrucker/werk) [![GoDoc](https://godoc.org/github.com/bendrucker/werk?status.svg)](https://godoc.org/github.com/bendrucker/werk)

Werk provides a worker pool that can concurrently process work up to a specified pool size.

## Usage

```go
import "github.com/bendrucker/werk"

func main() {
	pool := werk.NewPool(10)

	result, _ := pool.Do(context.TODO(), Work{Value: "beep boop"}, func(ctx context.Context, v interface{}) (interface{}, error) {
		fmt.Println("value:", v)
		fmt.Println("workers:", pool.Available())
		return "borp", nil
	})

	fmt.Println("result:", result)

	_, err := pool.Do(context.TODO(), Work{Value: "beep boop"}, func(ctx context.Context, v interface{}) (interface{}, error) {
		return nil, errors.New("oops")
	})

	fmt.Println("err:", err)

	// Output:
	// value: beep boop
	// workers: 9
	// result: borp
	// err: oops
}
```
