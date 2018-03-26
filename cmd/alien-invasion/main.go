package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strings"

	"github.com/davecgh/go-spew/spew"
)

// City represents a City with Name and 4 possible roads to neighbour cities
// North, East, South, West is either nil or a string to neighbour cities
// Alien stores the ID of alien currently in the city
type City struct {
	Name  string
	Roads map[int]string
	Alien int
}

type mapType map[string]*City

var iteration int
var maxCities int
var DEBUG bool

func init() {
	iteration = 0
}

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

	DEBUG = *debug
	mapRawData := readMapData(*mapFile)
	mapData := buildMap(mapRawData)
	maxCities := len(allCities(&mapData))
	_ = maxCities // I hate this...

	runSimulation(&mapData, int(*numAliens), 10000)
	_debug(mapData)
}

// Function runSimulation: main simulation loop
// Input: mapType map data, int number of aliens, int iterations to run
func runSimulation(mapData *mapType, numAliens int, iterations int) {
	_debug(fmt.Sprintf("Deploying %d aliens into cities...", numAliens))
	deployAliens(mapData, numAliens)

	for iteration = 1; iteration <= iterations; iteration++ {
		moveAliens(mapData)
	}

	// destroyCity(&fullMap, "Foo", 1, 2)
}

// Function readMap: reads in data from file, removes line endings
// Input: fileName to read from
// Returns: lines, slice of strings
func readMapData(fileName string) []string {
	var result []string

	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')

		if err != nil {
			if err == io.EOF {
				break
			} else {
				log.Fatalf("Can't read %s: %s", fileName, err.Error())
			}
		}

		line = strings.TrimSuffix(line, "\n");
		line = strings.TrimSuffix(line, "\r");
		result = append(result, line)
	}

	return result
}

// Function buildMap: builds map hash from lines to Cities
// Input: stringslice lines
// Returns: map of City (key: name of City, value: struct City)
func buildMap(mapData []string) mapType {
	cities := make(mapType)

	for _, line := range mapData {
		s := strings.Split(line, " ")
		cityName := s[0]
		city := &City{
			Name:  cityName,
			Roads: make(map[int]string),
		}

		for _, directionData := range s[1:] {
			directions := strings.Split(directionData, "=") // splits <direction>=<cityName>
			city.Roads[dirNameToInt(directions[0])] = directions[1]
		}

		cities[cityName] = city
	}

	validateRoads(cities)

	return cities
}

// Function validateRoads: walks all defined roads and validates the source and destination points
// Input: mapType mapData with full map
// Returns: void, but raises exception for validation errors
func validateRoads(mapData mapType) {
	for _, city := range mapData {
		for _, direction := range allRoads(city) {
			if toCityName, toOk := city.Roads[direction]; toOk {
				if toCity, toCityOk := mapData[toCityName]; toCityOk {
					if toCity.Roads[oppositeDirection(direction)] != city.Name {
						log.Fatalf("Map validation error: no back-road to %s from %s, but should be", toCity.Roads[oppositeDirection(direction)], city.Name)
					}
				} else {
					log.Fatalf("Map validation error: road to %s from %s, but %s not found on map", toCityName, city.Name, toCityName)
				}
			}
		}
	}
}

// Function destroyCity: removes a city from mapData, and removes it from neighbour cities as well
// Input: *mapType mapData, string cityToBeRemoved
// Returns: -
func destroyCity(mapData *mapType, cityName string, alien1 int, alien2 int) {
	city := (*mapData)[cityName]

	for _, direction := range allRoads(city) {
		delete((*mapData)[city.Roads[direction]].Roads, oppositeDirection(direction))
	}

	delete(*mapData, cityName)

	log.Printf("[iter %5d] %s has been destroyed by alien %d and alien %d\n", iteration, cityName, alien1, alien2)
	assertAnyCitiesLeft(mapData)
}

func deployAliens(mapData *mapType, numAliens int) {
	for i := 1; i <= numAliens; i++ {
		allCities := allCities(mapData) // need to re-read city keys as they could be destroyed during deployment
		// FIXME: could be more effective with a caching slice here
		numCities := len(allCities)

		randIndex := rand.Intn(numCities)
		// _debug(fmt.Sprintf("Moved alien %d to %s", i, allCities[randIndex]))
		moveAlienTo(mapData, allCities[randIndex], i)
	}
}

func moveAliens(mapData *mapType) {
	allCities := allCities(mapData)

	for _, cityName := range allCities {
		if city, ok := (*mapData)[cityName]; ok {
			if roads := allRoads(city); len(roads) > 0 { // there are still roads out of this city
				directionIndex := rand.Intn(len(roads))
				cityTo := city.Roads[roads[directionIndex]]
				moveAlienTo(mapData, cityTo, city.Alien)
				city.Alien = 0 // moved out from this city
			}
		} // else this city has been already destroyed in movements phase
	}
}

func moveAlienTo(mapData *mapType, cityName string, alien int) {
	city := (*mapData)[cityName]

	if city.Alien == 0 { // no alien in this city yet, move him in
		city.Alien = alien
	} else { // already an alien here, so they fight and destroy this city
		destroyCity(mapData, cityName, city.Alien, alien)
	}
}

func assertAnyCitiesLeft(mapData *mapType) {
	if len(allCities(mapData)) == 0 {
		log.Printf("[iter %5d] all cities (%d) have been destroyed", iteration, maxCities)
		os.Exit(0)
	}

}

// Function dirNameToInt: convert a direction name (north, ...) to integer value
// Input: string direction
// Returns: integer representation of direction
func dirNameToInt(direction string) int {
	dirHash := map[string]int{"north": 0, "east": 1, "south": 2, "west": 3}

	return dirHash[direction]
}

// Function oppositeDirection: returns the int representation of oppsite direction
// Input: direction int (eg. north is 0, etc.)
// Returns: direction int (eg. 2 for south, etc.)
func oppositeDirection(direction int) int {
	return (direction + 2) % 4
}

func allCities(mapData *mapType) []string {
	cities := make([]string, len(*mapData))

	i := 0
	for k := range *mapData {
		cities[i] = k
		i++
	}

	return cities
}

// Function allRoads: returns a slice of ints with roads from that city, simplifies other loops
// Input: City city
// Returns: []int roads
func allRoads(city *City) []int {
	var roads []int

	for direction := 0; direction < 4; direction++ {
		if _, toOk := city.Roads[direction]; toOk {
			roads = append(roads, direction)
		}
	}

	return roads
}

func _debug(obj ...interface{}) {
	if !DEBUG {
		return
	}

	spew.Dump(obj)
}
