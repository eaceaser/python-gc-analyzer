# python-gc-analyzer

A simple tool to parse Python (at least Python 3.4)'s gc debug logging.

## Building
`go build`

## Usage

Ensure that your python application has GC debugging enabled by adding `gc.set_debug(gc.DEBUG_STATS)`.

```
sh> invoke your python application 2>&1 | ./python-gc-analyzer
..time elapses..
2019/01/21 00:00:19 gc stats for previous 5s
  gen|count|total seconds|longest pause
    0|  113|     0.009400|0.000600
    1|   10|     0.018000|0.003500
    2|    1|     0.194100|0.194100
total|  124|     0.221500|-
```
> NB: python outputs gc debug logging to stderr so make sure you redirect it to stdout with 2>&1

See `python-gc-analyzer --help` for additional options.

## TODO

* additional in-process summary statistics as needed
* writing stats timeseries to data files for external analysis, eg numpy

## Background

Python doesn't output timestamps for individual GC runs, so there's no easy way to 
compute cumulative time spent in GC for a given interval from the debug output alone. 
This tool attempts to solve that by consuming stdin as quickly as possible to assign 
somewhat accurate timestamps to GC runs.
