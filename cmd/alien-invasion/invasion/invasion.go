package invasion

import (
    "bufio"
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

// Invasion holds a full map and implements methods to run simulation
type Invasion map[string]*City

var iteration int
var maxCities int
var DEBUG bool

func init() {
    iteration = 0
}

// Runs the main simulation loop: deploys and moves aliens, prints events
func (mapData *Invasion) Run(numAliens int, iterations int) {
    _debug(fmt.Sprintf("Deploying %d aliens into cities...", numAliens))
    mapData.Deploy(numAliens)

    for iteration = 1; iteration <= iterations; iteration++ {
        mapData.Move()
    }
}

// Reads in data from map file, removes line endings, returns all lines
func (mapData *Invasion) ReadMap(fileName string) []string {
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

// builds map hash from lines into Cities, stores map data in caller struct
func (mapData *Invasion) BuildMap(fileName string) {
    mapLines := mapData.ReadMap(fileName)

    cities := make(Invasion)

    for _, line := range mapLines {
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

    cities.validateRoads()

    maxCities := len(cities.AllCities())
    _ = maxCities // I hate this...

    *mapData = cities
}

// walks all defined roads and validates the source and destination points, bails if missing cities found
func (mapData *Invasion) validateRoads() {
    for _, city := range *mapData {
        for _, direction := range AllRoads(city) {
            if toCityName, toOk := city.Roads[direction]; toOk {
                if toCity, toCityOk := (*mapData)[toCityName]; toCityOk {
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

// removes a city from map, and removes it from neighbour cities as well, prints destruction fact
func (mapData *Invasion) DestroyCity(cityName string, alien1 int, alien2 int) {
    city := (*mapData)[cityName]

    for _, direction := range AllRoads(city) {
        delete((*mapData)[city.Roads[direction]].Roads, oppositeDirection(direction))
    }

    delete(*mapData, cityName)

    log.Printf("[iter %5d] %s has been destroyed by alien %d and alien %d\n", iteration, cityName, alien1, alien2)
    mapData.AssertAnyCitiesLeft()
}

// initially deploys aliens into cities randomly, takes care of 2 aliens in the same city destroys the city
func (mapData *Invasion) Deploy(numAliens int) {
    for i := 1; i <= numAliens; i++ {
        allCities := mapData.AllCities() // need to re-read city keys as they could be destroyed during deployment
        // FIXME: could be more effective with a caching slice here
        numCities := len(allCities)
        randIndex := rand.Intn(numCities)
        // _debug(fmt.Sprintf("Moved alien %d to %s", i, allCities[randIndex]))
        mapData.MoveAlienTo(allCities[randIndex], i)
    }
}

// iterates over all cities and moves aliens if roads are still present
func (mapData *Invasion) Move() {
    allCities := mapData.AllCities()

    for _, cityName := range allCities {
        if city, ok := (*mapData)[cityName]; ok {
            if city.Alien == 0 { //no alien here to move
                continue
            }
            if roads := AllRoads(city); len(roads) > 0 { // there are still roads out of this city
                directionIndex := rand.Intn(len(roads))
                cityTo := city.Roads[roads[directionIndex]]
                mapData.MoveAlienTo(cityTo, city.Alien)
                city.Alien = 0 // moved out from this city
            }
        } // else this city has been already destroyed in movements phase
    }
}

// moves alien into a City, calls destroy if another alien is already there
func (mapData *Invasion) MoveAlienTo(cityName string, alien int) {
    city := (*mapData)[cityName]

    if city.Alien == 0 { // no alien in this city yet, move him in
        city.Alien = alien
    } else { // already an alien here, so they fight and destroy this city
        mapData.DestroyCity(cityName, city.Alien, alien)
    }
}

// terminates the simulation if all of the cities have been destroyed
func (mapData *Invasion) AssertAnyCitiesLeft() {
    if len(mapData.AllCities()) == 0 {
        log.Printf("[iter %5d] all cities (%d) have been destroyed", iteration, maxCities)
        os.Exit(0)
    }

}

// convert a direction name (north, ...) to integer value
func dirNameToInt(direction string) int {
    dirHash := map[string]int{"north": 0, "east": 1, "south": 2, "west": 3}

    return dirHash[direction]
}

// returns the int representation of oppsite direction
func oppositeDirection(direction int) int {
    return (direction + 2) % 4
}

// returns all city names of full map
func (mapData *Invasion) AllCities() []string {
    cities := make([]string, len(*mapData))

    i := 0
    for k := range *mapData {
        cities[i] = k
        i++
    }

    return cities
}

// returns a slice of ints with roads from that city, simplifies other loops
func AllRoads(city *City) []int {
    var roads []int

    for direction := 0; direction < 4; direction++ {
        if _, toOk := city.Roads[direction]; toOk {
            roads = append(roads, direction)
        }
    }

    return roads
}

// debug, dumps full map data
func (mapData *Invasion) Dump() {
    spew.Dump(mapData)
}

// internal debug printer, controlled by DEBUG bool
func _debug(obj ...interface{}) {
    if !DEBUG {
        return
    }

    spew.Dump(obj)
}
