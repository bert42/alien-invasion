package main

import (
	"fmt"
	"log"
	"os"
	"io"
	"bufio"
	"strings"

	"github.com/davecgh/go-spew/spew"
)

// City represents a City with Name and 4 possible roads to neighbour cities
// North, East, South, West is either nil or a string to neighbour cities
type City struct {
	Name	string
	North	string
	East	string
	South	string
	West	string
}

type mapType map[string]City


func main() {
	mapData := readMap()
	fullMap := buildMap(mapData)

	spew.Dump(fullMap)

	fmt.Println("Done.")
}


// Function readMap: reads in data from Stdin, removes line endings
// Input: -
// Returns: lines, slice of strings
func readMap() []string {
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
// Input: lines
// Returns: map of City (key: name of City, value: struct City)
func buildMap(mapData []string) mapType {
	cities := make(mapType);

	for _, line := range mapData {
		s := strings.Split(line, " ")
		cityName := s[0]
		city := City{
			Name: cityName,
		}
		for _, directionData := range s[1:] {
			directions := strings.Split(directionData, "=") // splits <direction>=<cityName>
			switch directions[0] {
			case "north":
				city.North = directions[1]
			case "east":
				city.East = directions[1]
			case "south":
				city.South = directions[1]
			case "west":
				city.West = directions[1]
			}
		}

		cities[cityName] = city
	}

	return cities
}
