// Code generated by goderive DO NOT EDIT.

package pipeline

import (
	"sync"
)

// derivePipeline composes f and g into a concurrent pipeline.
func derivePipeline(f func(lines []string) <-chan string, g func(line string) <-chan int) func([]string) <-chan int {
	return func(a []string) <-chan int {
		b := f(a)
		return deriveJoin(deriveFmap(g, b))
	}
}

// deriveJoin listens on all channels resulting from the input channel and sends all their results on the output channel.
func deriveJoin(in <-chan (<-chan int)) <-chan int {
	out := make(chan int)
	go func() {
		wait := sync.WaitGroup{}
		for c := range in {
			wait.Add(1)
			res := c
			go func() {
				for r := range res {
					out <- r
				}
				wait.Done()
			}()
		}
		wait.Wait()
		close(out)
	}()
	return out
}

// deriveFmap returns an output channel where the items are the result of the input function being applied to the items on the input channel.
func deriveFmap(f func(string) <-chan int, in <-chan string) <-chan (<-chan int) {
	out := make(chan (<-chan int), cap(in))
	go func() {
		for a := range in {
			b := f(a)
			out <- b
		}
		close(out)
	}()
	return out
}
