package main

import (
	"fmt"
	"testing"
)

func TestParse(t *testing.T) {
	c := make(chan string)
	o := make(chan *CollectionStat, 100)
	go func() {
		defer close(c)
		c <- mkStart(1)
		c <- mkStats(1000, 20000, 300000)
		c <- mkFullEnd(5, 30, 0.002)
		c <- mkStart(2)
		c <- mkStats(232, 293093, 290)
		c <- mkEnd(0.05)
	}()

	p := NewParser()
	p.Parse(c, o)
}

func mkStart(gen int) string {
	return fmt.Sprintf("gc: collecting generation %d...\n", gen)
}

func mkStats(gen0 int, gen1 int, gen2 int) string {
	return fmt.Sprintf("gc: objects in each generation: %d %d %d\n", gen0, gen1, gen2)
}

func mkEnd(duration float32) string {
	return fmt.Sprintf("gc: done, %fs elapsed.\n", duration)
}

func mkFullEnd(unreachable int, uncollectable int, duration float32) string {
	return fmt.Sprintf("gc: done, %d unreachable, %d uncollectable, %fs elapsed.\n", unreachable, uncollectable, duration)
}

