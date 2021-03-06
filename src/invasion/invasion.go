// package invasion represents a simulation of an alien invasion on a pre-defined
// map; it deploys aliens, let them wander using the defined roads between cities
// and handle their (and the city's) destruction whenever they meet
package invasion

import (
    "fmt"
    "io/ioutil"
    "log"
    "math/rand"
    "time"
    "strings"
    "errors"

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
type Invasion struct {
    Map map[string]*City
    VerboseLog []string
    DebugLog []string
    Iteration int
    Statistics struct {
        MaxCities int
        NumberOfMoves int
    }
}

// New constructs a new Invasion, sets up default values
func New() *Invasion {
    rand.Seed(time.Now().UTC().UnixNano())

    return &Invasion{
        VerboseLog: make([]string, 0),
        DebugLog: make([]string, 0),
        Iteration: 0,
    }
}

// Run starts the main simulation loop: deploys and moves aliens, prints events
func (data *Invasion) Run(numAliens int, iterations int) {
    data.verbose(fmt.Sprintf("Deploying %d aliens into cities...", numAliens))
    data.Deploy(numAliens)

    defer data.PrintStatistics()

    for data.Iteration = 1; data.Iteration < iterations; data.Iteration++ {
        data.Move()
        if isAny := data.AnyCitiesLeft(); !isAny {
            data.Print(fmt.Sprintf("all cities (%d) have been destroyed", data.Statistics.MaxCities))
            return
        }
    }
    data.Print(fmt.Sprintf("simulation terminated after %d iterations", data.Iteration))
}

func (data *Invasion) PrintStatistics() {
    cities := data.AllCities()
    var citiesStr string // just to make it beautiful for singular/plural
    if data.Statistics.MaxCities-len(cities) == 1 {
        citiesStr = "city was"
    } else {
        citiesStr = "cities were"
    }


    data.Print("Statistics:\n")
    data.Print(fmt.Sprintf("\tnumber of moves executed: %d", data.Statistics.NumberOfMoves))
    data.Print(fmt.Sprintf("\t%d %s destroyed out of %d, %d remained",
        data.Statistics.MaxCities-len(cities), citiesStr, data.Statistics.MaxCities, len(cities)))
}

// BuildMap reads map file and builds map struct from lines into Cities, stores map data in Invasion.Map
func (data *Invasion) BuildMap(fileName string) {
    contents, err := ioutil.ReadFile(fileName)
    if err != nil {
        log.Fatal(err)
    }

    mapLines := strings.Split(string(contents[:]), "\n")
    for _, line := range mapLines {
        line = strings.TrimSuffix(line, "\r");
    }

    cities := make(map[string]*City)

    for _, line := range mapLines {
        s := strings.Split(line, " ")
        cityName := s[0]

        if len(cityName) == 0 || cityName[0:1] == "#" {
            continue
        }

        if _, cityExists := cities[cityName]; cityExists {
            log.Fatalf("City '%s' is redefined in map, it's not supported, yet", cityName)
        }
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

    data.Map = cities
    data.Statistics.MaxCities = len(data.AllCities())
    if err := data.ValidateRoads(); err != nil {
        log.Fatal(err.Error())
    }
}

// ValidateRoads walks all defined roads and validates the source and destination
// points, returns error if missing cities are found, describing the problem
func (data *Invasion) ValidateRoads() error {
    for _, city := range data.Map {
        for _, direction := range AllRoads(city) {
            toCityName, toExists := city.Roads[direction]
            if !toExists {
                continue
            }

            toCity, toCityOk := data.Map[toCityName];
            if !toCityOk {
                return errors.New(fmt.Sprintf("Map validation error: road to %s from %s, but %s not found on map", toCityName, city.Name, toCityName))
            }

            if toCity.Roads[oppositeDirection(direction)] != city.Name {
                return errors.New(fmt.Sprintf("Map validation error: no back-road to %s from %s, but should be", toCity.Roads[oppositeDirection(direction)], city.Name))
            }
        }
    }

    return nil
}

// DestroyCity removes a city from map, and removes it from neighbour cities
// prints destruction fact
func (data *Invasion) DestroyCity(cityName string, alien1 int, alien2 int) {
    city := data.Map[cityName]

    for _, direction := range AllRoads(city) {
        delete(data.Map[city.Roads[direction]].Roads, oppositeDirection(direction))
    }

    delete(data.Map, cityName)

    data.Print(fmt.Sprintf("%s has been destroyed by alien %d and alien %d\n", cityName, alien1, alien2))
}

// Deploy initially deploys aliens into cities randomly, takes care of 2 aliens in the same city destroys the city
func (data *Invasion) Deploy(numAliens int) {
    for i := 1; i <= numAliens; i++ {
        allCities := data.AllCities() // need to re-read city keys as they could be destroyed during deployment
        // FIXME: could be more effective with a caching slice here
        numCities := len(allCities)
        randIndex := rand.Intn(numCities)

        data.MoveAlienTo(allCities[randIndex], i)
    }
    data.debug(data.Map)
}

// Move iterates over all cities and moves aliens if any road is still present in that city
func (data *Invasion) Move() {
    allCities := data.AllCities()
    moveCache := make(map[string]bool, 0) // prevent back and forth move between neighbouring cities

    for _, cityName := range allCities {
        city, exists := data.Map[cityName];
        if !exists { // this city has been already destroyed in movements phase
            continue
        }
        if _, exists := moveCache[cityName]; exists { // already moved into this one
            continue
        }
        if city.Alien == 0 { // no alien here to move
            continue
        }

        targets := data.TargetCitiesFrom(city)
        if len(targets) == 0 { // no roads out of this city
            continue
        }

        directionIndex := rand.Intn(len(targets))
        cityTo := targets[directionIndex]
        data.MoveAlienTo(cityTo, city.Alien)
        city.Alien = 0 // moved out from this city
        moveCache[cityTo] = true
    }
}

// MoveAlienTo moves alien into a City, calls destroy if another alien is already there
func (data *Invasion) MoveAlienTo(cityName string, alien int) {
    city := data.Map[cityName]
    data.Statistics.NumberOfMoves++

    if city.Alien == 0 { // no alien in this city yet, move him in
        city.Alien = alien
    } else { // already an alien here, so they fight and destroy this city
        data.DestroyCity(cityName, city.Alien, alien)
    }
}

// TargetCitiesFrom returns all cities reachable from current one
func (data *Invasion) TargetCitiesFrom(city *City) (cities []string) {
    for direction := 0; direction < 4; direction++ {
        if cityName, exists := city.Roads[direction]; exists {
            cities = append(cities, cityName)
        }
    }

    return cities
}

// AnyCitiesLeft returns false if all of the cities have been destroyed
func (data *Invasion) AnyCitiesLeft() bool {
    return len(data.AllCities()) != 0
}

// dirNameToInt convert a direction name (north, ...) to integer value
func dirNameToInt(direction string) int {
    dirHash := map[string]int{"north": 0, "east": 1, "south": 2, "west": 3}

    return dirHash[direction]
}

// oppositeDirection returns the int representation of opposite direction from a city
func oppositeDirection(direction int) int {
    return (direction + 2) % 4
}

// AllCities returns all city names present on the full map
func (data *Invasion) AllCities() []string {
    cities := make([]string, len(data.Map))

    i := 0
    for k := range data.Map {
        cities[i] = k
        i++
    }

    return cities
}

// ALlRoads returns a slice of ints with roads from that city, simplifies other loops
func AllRoads(city *City) (roads []int) {
    for direction := 0; direction < 4; direction++ {
        if _, toOk := city.Roads[direction]; toOk {
            roads = append(roads, direction)
        }
    }

    return roads
}

// Print outputs string in log format with iterations prefix
func (data *Invasion) Print(s string) {
    log.Printf("[iter %5d] %s", data.Iteration, s)
}

// Dump dumps full map data for debugging purposes
func (data *Invasion) Dump() {
    spew.Dump(data)
}

// verbose is an internal verbose printer, collects output in VerboseLog
func (data *Invasion) verbose(str string) {
    data.VerboseLog = append(data.VerboseLog, str)
}

// debug is an internal debug printer, collects output in DebugLog
func (data *Invasion) debug(obj ...interface{}) {
    data.DebugLog = append(data.DebugLog, spew.Sdump(obj))
}
