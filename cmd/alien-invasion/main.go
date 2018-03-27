package main

import (
	"flag"
	"fmt"

	"./invasion"
)

func main() {
	flag.Usage = func() {
		fmt.Println(`Alien-Invasion

    Reads in pre-defined map data from given map file (txt format), then
    runs simulation of N aliens wandering in cities fighting when met.

    Usage: alien-invasion [flags]
    `)

		flag.PrintDefaults()
	}

	mapFile := flag.String("map", "examples/testmap1.txt", "map file path")
	numAliens := flag.Int("aliens", 3, "number of aliens to be deployed")
	debug := flag.Bool("debug", false, "debug mode dumps map definitions")

	flag.Parse()

	invasion.DEBUG = *debug

	Simulation := &invasion.Invasion{}
	Simulation.BuildMap(*mapFile)

	Simulation.Run(int(*numAliens), 10000)
//	Simulation.Dump()
}

