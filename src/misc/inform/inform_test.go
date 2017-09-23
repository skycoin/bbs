package inform

import (
	"os"
	"sync"
	"testing"
)

func TestNewLogger(t *testing.T) {
	const count = 10
	logger := NewLogger(true, os.Stdout, "TEST")
	var wg sync.WaitGroup
	wg.Add(count)
	for i := 0; i < count; i++ {
		go func(i int) {
			switch i % 2 {
			case 0:
				logger.Printf("Number: %d", i)
			case 1:
				logger.Println("Number:", i)
			}

			wg.Done()
		}(i)
	}
	wg.Wait()
}
