package main

import (
	"fmt"
	"github.com/montanaflynn/stats"
	"os"
	"sync"
	"text/tabwriter"
	"time"
)

type CollectionStat struct {
	gen int
	genSizes []int
	unreachableCount int
	uncollectableCount int
	elapsedMicros int
	startTime time.Time
	endTime time.Time
}

type Stats struct {
	lock sync.Mutex
	gen0Times stats.Float64Data
	gen0Sum int
	gen1Times stats.Float64Data
	gen1Sum int
	gen2Times stats.Float64Data
	gen2Sum int
}

func NewStats() *Stats {
	rv := &Stats{}
	rv.resetStats()
	return rv
}

func (s *Stats) Process(c chan *CollectionStat) {
	for stat := range c {
		s.lock.Lock()
		switch stat.gen {
		case 0:
			s.gen0Times = append(s.gen0Times, float64(stat.elapsedMicros))
			s.gen0Sum += stat.elapsedMicros
			break
		case 1:
			s.gen1Times = append(s.gen1Times, float64(stat.elapsedMicros))
			s.gen1Sum += stat.elapsedMicros
			break
		case 2:
			s.gen2Times = append(s.gen2Times, float64(stat.elapsedMicros))
			s.gen2Sum += stat.elapsedMicros
			break
		}
		s.lock.Unlock()
	}
}

func (s *Stats) PrintAndFlush() {
	s.lock.Lock()
	defer s.lock.Unlock()
	defer s.resetStats()

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 0, ' ', tabwriter.Debug|tabwriter.AlignRight)

	gen0Max, err := stats.Max(s.gen0Times)
	if err != nil {
		gen0Max = 0
	}

	gen1Max, err := stats.Max(s.gen1Times)
	if err != nil {
		gen1Max = 0
	}
	gen2Max, err := stats.Max(s.gen2Times)
	if err != nil {
		gen2Max = 0
	}

	totalCount := len(s.gen0Times) + len(s.gen1Times) + len(s.gen2Times)

	fmt.Fprintln(w, "gen\tcount\ttotal seconds\tlongest pause")
	fmt.Fprintf(w, "0\t%d\t%f\t%f\n", len(s.gen0Times), float64(s.gen0Sum) / 10000, gen0Max / 10000)
	fmt.Fprintf(w, "1\t%d\t%f\t%f\n", len(s.gen1Times), float64(s.gen1Sum) / 10000, gen1Max / 10000)
	fmt.Fprintf(w, "2\t%d\t%f\t%f\n", len(s.gen2Times), float64(s.gen2Sum) / 10000, gen2Max / 10000)
	fmt.Fprintf(w, "total\t%d\t%f\t-\n", totalCount, float64(s.gen0Sum+s.gen1Sum+s.gen2Sum) / 10000)
	w.Flush()
}

func (s *Stats) resetStats() {
	s.gen0Sum = 0
	s.gen0Times = make(stats.Float64Data, 0)
	s.gen1Sum = 0
	s.gen1Times = make(stats.Float64Data, 0)
	s.gen2Sum = 0
	s.gen2Times = make(stats.Float64Data, 0)
}
