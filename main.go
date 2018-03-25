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
    "math/rand"

    "github.com/davecgh/go-spew/spew"
)

const DEBUG = true

// City represents a City with Name and 4 possible roads to neighbour cities
// North, East, South, West is either nil or a string to neighbour cities
type City struct {
    Name    string
    Roads   map[int]string
    Alien   int
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

    runSimulation(fullMap, int(numAliens), 10000)
    _debug(fullMap)

    fmt.Println("Done.")
}

// Function runSimulation: main simulation loop
// Input: mapType map data, int number of aliens, int iterations to run
func runSimulation(fullMap mapType, numAliens int, iterations int) {
    _debug(fmt.Sprintf("Deploying %d aliens into cities...", numAliens))
    deployAliens(&fullMap, numAliens)

    // destroyCity(&fullMap, "Foo", 1, 2)
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
    city  := (*mapData)[cityName]

    for _, direction := range allRoads(city) {
        delete((*mapData)[city.Roads[direction]].Roads, oppositeDirection(direction))
    }

    delete(*mapData, cityName)

    fmt.Printf("%s has been destroyed by alien %d and alien %d\n", cityName, alien1, alien2)
}

func deployAliens(mapData *mapType, numAliens int)  {
    for i:=0; i<numAliens; i++ {
        allCities := allCities(mapData)  // need to re-read city keys as they could be destroyed during deployment
                                         // FIXME: could be more effective with a caching slice here

        randIndex := rand.Intn(len(allCities))
        _debug(fmt.Sprintf("Moved alien %d to %s", i, allCities[randIndex]))
        moveAlienTo(mapData, allCities[randIndex], i)
    }
}

func moveAlienTo(mapData *mapType, cityName string, alien int) {
    city := (*mapData)[cityName]
    if city.Alien == 0 {  // no alien in this city yet, move him in
        city.Alien = alien
    } else {              // already an alien here, so they fight and destroy this city
        destroyCity(mapData, cityName, city.Alien, alien)
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
func allRoads(city City) []int {
    var roads []int;

    for direction := 0; direction < 4; direction++ {
        if _, toOk := city.Roads[direction]; toOk {
            roads = append(roads, direction)
        }
    }

    return roads;
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

func _debug(obj ...interface{}) {
    if (!DEBUG) { return }

    spew.Dump(obj)
}
