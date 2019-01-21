package main

import (
	"bufio"
	flag "github.com/spf13/pflag"
	"io"
	"log"
	"os"
	"time"
)

const defaultTimerDuration = "5s"
const channelBufferSize = 2048

func main() {
	timerDuration := flag.String("interval", defaultTimerDuration, "reporting interval in duration format")
	flag.Parse()

	parser := NewParser()
	stats := NewStats()

	readerToParser := make(chan string, channelBufferSize)
	parserToStats := make(chan *CollectionStat, channelBufferSize)
	go reader(os.Stdin, readerToParser)
	go parser.Parse(readerToParser, parserToStats)
	go timer(stats, *timerDuration)
	stats.Process(parserToStats)
}

func timer(s *Stats, duration string) {
	d, err := time.ParseDuration(duration)
	if err != nil {
		log.Panic(err)
	}

	t := time.Tick(d)
	for range t {
		log.Print("gc stats for previous " + duration)
		s.PrintAndFlush()
	}
}

func reader(r io.Reader, c chan string) {
	s := bufio.NewScanner(r)
	for s.Scan() {
		line := s.Text()
		c <- line
	}

	if s.Err() != nil {
		log.Print("stdin closed.")
		close(c)
	}
}

