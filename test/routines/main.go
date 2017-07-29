package main

import (
	"fmt"
	"math"
	"os"
	"runtime/pprof"
	"runtime/trace"
	"sync"
	"time"
)

func main() {
	f, err := os.Create("./cpu.prof")
	if err != nil {
		fmt.Printf("could not create CPU profile: %s\n", err.Error())
		return
	}

	f2, err := os.Create("./trace.out")
	if err != nil {
		fmt.Printf("could not create trace: %s\n", err.Error())
		return
	}

	if err := pprof.StartCPUProfile(f); err != nil {
		fmt.Printf("could not start CPU profile: %s\n", err.Error())
		return
	}

	trace.Start(f2)

	count := 10
	var wg sync.WaitGroup
	wg.Add(count)

	results := make([]chan float64, count)
	for i := 0; i < count; i++ {
		go calculationLoop(&wg, results[i])
	}

	wg.Wait()

	trace.Stop()

	pprof.StopCPUProfile()

	f.Close()
	f2.Close()

	fmt.Printf("Complete\n")
}

func calculationLoop(wg *sync.WaitGroup, result chan float64) {
	last := 0.0
	for i := 0; i < 5000; i++ {
		// Calculate!
		x := float64(i + 1000000000)
		z := squareRoot(x)
		last = z

		if i%10 == 0 {
			time.Sleep(2 * time.Millisecond)
		}
	}

	wg.Done()

	result <- last
}

func squareRoot(x float64) float64 {
	// Code copied from https://gist.github.com/pstoll/4106979
	z := 1.0
	minDelta := 0.00000000001
	delta := 1.0
	i := 0
	for ; math.Abs(delta) > minDelta; i++ {
		delta = (z*z - x) / (2 * z)
		z = z - delta
	}
	return z
}
