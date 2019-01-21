package main

import (
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	gcStart   = regexp.MustCompile("collecting generation (\\d+)...")
	gcStats   = regexp.MustCompile("objects in each generation: (\\d+) (\\d+) (\\d+)")
	gcEnd     = regexp.MustCompile("done, ([\\d\\.]+)s elapsed.")
	gcEndFull = regexp.MustCompile("done, (\\d+) unreachable, (\\d+) uncollectable, ([\\d\\.]+)s elapsed.")
)

type Parser struct {
	currentCollection *CollectionStat
}

func NewParser() *Parser {
	p := &Parser{}
	return p
}

func (p *Parser) Parse(c chan string, o chan *CollectionStat) {
	for line := range c {
		if strings.HasPrefix(line, "gc:") {
			p.processGcLine(line, o)
		}
	}
}

func (p *Parser) processGcLine(line string, o chan *CollectionStat) {
	ts := time.Now()

	if matches := gcStart.FindStringSubmatch(line); len(matches) > 0 {
		if p.currentCollection != nil {
			log.Print("found start of gc while previous is active")
			return
		}

		genNum, err := strconv.Atoi(matches[1])
		if err != nil {
			log.Printf("unable to parse collection number %v", err)
			return
		}

		p.currentCollection = &CollectionStat{startTime: ts, gen: genNum}
	} else if matches := gcStats.FindStringSubmatch(line); len(matches) > 0 {
		if p.currentCollection == nil {
			log.Print("found gc collection sizes when no gc is active")
			return
		}

		sizes := make([]int, 3)
		for i, sizeStr := range matches[1:4] {
			size, err := strconv.Atoi(sizeStr)
			if err != nil {
				log.Printf("unable to parse generation size %v", err)
			}
			sizes[i] = size
		}
		p.currentCollection.genSizes = sizes
	} else if matches := gcEnd.FindStringSubmatch(line); len(matches) > 0 {
		if p.currentCollection == nil {
			log.Print("found gc end while no gc is active")
			return
		}

		elapsed, err := strconv.ParseFloat(matches[1], 32)
		if err != nil {
			log.Printf("unable to parse elapsed time %v", err)
			return
		}

		p.currentCollection.elapsedMicros = int(elapsed*10000)
		p.currentCollection.endTime = ts

		o <- p.currentCollection

		p.currentCollection = nil
	} else if matches := gcEndFull.FindStringSubmatch(line); len(matches) > 0 {
		if p.currentCollection == nil {
			log.Print("found gc end while no gc is active")
			return
		}

		unreachable, err := strconv.Atoi(matches[1])
		if err != nil {
			log.Printf("unable to parse unreachable count %v", err)
			return
		}

		uncollectable, err := strconv.Atoi(matches[2])
		if err != nil {
			log.Printf("unable to parse uncollectable count %v", err)
			return
		}

		elapsed, err := strconv.ParseFloat(matches[3], 32)
		if err != nil {
			log.Printf("unable to parse elapsed time %v", err)
			return
		}

		p.currentCollection.unreachableCount = unreachable
		p.currentCollection.uncollectableCount = uncollectable
		p.currentCollection.elapsedMicros = int(elapsed*10000)
		p.currentCollection.endTime = ts

		o <- p.currentCollection

		p.currentCollection = nil
	} else {
		log.Print("Unknown gc log line")
		log.Print(line)
	}
}
