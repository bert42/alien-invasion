package main

import (
	"fmt"
	"log"
	"os"
	"io"
	"bufio"
	"strings"
	"flag"
	"strconv"

	"github.com/davecgh/go-spew/spew"
)

// City represents a City with Name and 4 possible roads to neighbour cities
// North, East, South, West is either nil or a string to neighbour cities
type City struct {
	Name	string
	Roads 	map[int]string
}

type mapType map[string]City



func main() {
	flag.Parse()

	numAliensString := flag.Arg(0)
	numAliens, err  := strconv.ParseInt(numAliensString, 10, 32)
	if numAliensString == "" || err != nil {
		Usage(1)
	}

	mapData := readMapData()
	fullMap := buildMap(mapData)

	spew.Dump(fullMap)

	runSimulation(fullMap, int(numAliens), 10000)
	fmt.Println("Done.")
}

// Function runSimulation: main simulation loop
// Input: mapType map data, int number of aliens, int iterations to run
func runSimulation(fullMap mapType, numAliens int, iterations int) {
	fmt.Println(fullMap, numAliens, iterations)
}


// Function readMap: reads in data from Stdin, removes line endings
// Input: -
// Returns: lines, slice of strings
func readMapData() []string {
	var result []string

	reader := bufio.NewReader(os.Stdin)
	for {
		line, err := reader.ReadString('\n')

		if (err != nil) { 
			if err == io.EOF {
				break
			} else {
				log.Fatal("Can't read stdin: "+err.Error())
			}
		}

		result = append(result, strings.TrimSuffix(line, "\r\n"))
	}

	return result
}

// Function buildMap: builds map hash from lines to Cities
// Input: stringslice lines
// Returns: map of City (key: name of City, value: struct City)
func buildMap(mapData []string) mapType {
	cities := make(mapType);

	for _, line := range mapData {
		s := strings.Split(line, " ")
		cityName := s[0]
		city := City{
			Name: cityName,
			Roads: make(map[int]string),
		}

		for _, directionData := range s[1:] {
			directions := strings.Split(directionData, "=") // splits <direction>=<cityName>
			city.Roads[dirNameToInt(directions[0])] = directions[1]
		}

		cities[cityName] = city
	}

	//TODO: needs validation for roads between cities point both directions

	return cities
}

// Function dirNameToInt: convert a direction name (north, ...) to integer value
// Input: string direction
// Returns: integer representation of direction
func dirNameToInt(direction string) int {
	dirHash := map[string]int{"north": 0, "east": 1, "south": 2, "west": 3}

	return dirHash[direction]
}


// Function Usage: prints usage info to terminal and exits
// Input: int exitCode
// Returns: -
func Usage(exitCode int) {
	fmt.Println(`Alien-Invasion

		Reads in pre-defined map data from given map file (txt format), then
		runs simulation of N aliens wandering in cities fighting when met.

		Usage: alien-invasion <number of aliens>
		`)

	os.Exit(exitCode)
}

