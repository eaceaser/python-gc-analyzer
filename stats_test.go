package main

import (
	"testing"
	"time"
)

func TestStats(t *testing.T) {
	c := make(chan *CollectionStat)
	go func() {
		defer close(c)
		c <- &CollectionStat{
			0,
			[]int{10, 20, 200},
			0,
			0,
			100,
			time.Now(),
			time.Now(),
		}
		c <- &CollectionStat{
			1,
			[]int{10, 20, 200},
			0,
			0,
			100,
			time.Now(),
			time.Now(),
		}
		c <- &CollectionStat{
			2,
			[]int{10, 20, 200},
			0,
			0,
			100,
			time.Now(),
			time.Now(),
		}
	}()

	s := NewStats()
	s.Process(c)
	s.PrintAndFlush()
}
