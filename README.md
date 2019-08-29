## tokenbucket

A elegant token bucket implementation in Go, only math calculations, no goroutine and channel.

## Usage

```go
package main

import (
	"fmt"
	"time"

	"github.com/lifenod/tokenbucket"
)

func main() {
	lim := tokenbucket.NewLimiter(5, 5)

	now := time.Now()
	for i := 0; i < 10; i++ {
		lim.Wait()
		fmt.Println(time.Now().Sub(now))
	}
}
```

Got output like the following:

```
1.722µs
42.783µs
47.06µs
49.333µs
51.114µs
204.77297ms
401.448733ms
604.790073ms
804.361192ms
1.003988204s
```

### Allow function

This function is used for manual control sleeping flow,
for example, you need to force cancel when sleeping:

```go
package main

import (
	"fmt"
	"time"

	"github.com/lifenod/tokenbucket"
)

func main() {
	lim := tokenbucket.NewLimiter(5, 5)

	cancelC := make(chan struct{})
	go func() {
		time.Sleep(time.Second * 2)
		close(cancelC)
	}()

	now := time.Now()
	for i := 0; i < 100; i++ {
		isAllowed, sleep := lim.Allow(time.Now())
		if !isAllowed {
			select {
			case <-cancelC:
				fmt.Println("canceled.")
				return
			case <-time.NewTimer(sleep).C:
			}
		}

		fmt.Println(time.Now().Sub(now))
	}
}
```

Got output like the following:

```
2.058µs
45.852µs
48.722µs
50.58µs
52.429µs
202.293412ms
404.180016ms
602.524459ms
801.456306ms
1.005150736s
canceled.
```
