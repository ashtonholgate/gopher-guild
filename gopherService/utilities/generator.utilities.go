package utilities

import (
	"fmt"
	"math"
	"math/rand"
	"sync"
)

func random(c chan int) {
	randomNumber := rand.Intn(1000)
	c <- randomNumber
}

func GenerateRandomNumber() int {
	c := make(chan int)
	const goroutineCount = 1000

	var wg sync.WaitGroup
	wg.Add(goroutineCount)

	for i := 0; i < goroutineCount; i++ {
		go func() {
			defer wg.Done()
			random(c)
		}()
	}

	go func() {
		wg.Wait()
		close(c)
	}()

	var largestNumber int
	for elem := range c {
		largestNumber = int(math.Max(float64(largestNumber), float64(elem)))
		fmt.Println(elem)
	}
	return largestNumber
}
