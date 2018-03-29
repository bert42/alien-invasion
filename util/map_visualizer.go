package main

import (
    "flag"
    "fmt"
    "os"
    "github.com/kr/pretty"

    "../cmd/alien-invasion/invasion"
)

func main() {
    flag.Usage = func() {
        fmt.Println(`Map Visualizer

    Reads in pre-defined map data and displays it on screen
    `)

        flag.PrintDefaults()
    }

    mapFile := flag.String("map", "", "map file path (required)")

    flag.Parse()

    if *mapFile == "" {
        fmt.Println("-map is required")
        os.Exit(2)
    }

    Simulation := invasion.New()
    Simulation.BuildMap(*mapFile)

    visualData := visualize(Simulation)
    pretty.Println(visualData)

//  fmt.Println(visualData)
}

// cache holds already processed city names to help walk routine back-tracking
var cache map[string]bool

// visualize returns the 2D string array representation of the map
func visualize(data *invasion.Invasion) (result [][]string) {
    cache = make(map[string]bool, 0)

    result = [][]string{{""}}
    walkCities(data, firstCityName(data), 0, 0, 0, 0, &result)

    return result
}

// walkCities is a recursive function to walk all roads defined between cities
// it expands the 2D array holding the map automatically in all directions
func walkCities(data *invasion.Invasion, cityName string, x int, y int, xD int, yD int, result *[][]string) {
    fmt.Println("-> called: ", cityName, x, y, xD, yD) //debug
    if _, ok := cache[cityName]; ok { // city already processed
        return
    }

    if x < 0 { x = 0 }
    if y < 0 { y = 0 }

    city := data.Map[cityName]
    cache[cityName] = true

    if x+xD < 0 {
        expand(result, -1, 0)
    }
    if x+xD > len(*result) {
        expand(result, 1, 0)
    }
    if y+yD < 0 {
        expand(result, 0, -1)
    }
    if y+yD > len((*result)[0]) {
        expand(result, 0, 1)
    }

    fmt.Println("---> processing: ", *result, x, y) //debug

    // store City at current coordinates, then walk around in all directions
    (*result)[x][y] = cityName[0:1]

    for direction := 0; direction < 4; direction++ {
        if nextCityName, toOk := city.Roads[direction]; toOk {
            var newxD int
            var newyD int
            if direction == 0 { // north
                newxD =  0; newyD = -1
            }
            if direction == 1 { // east
                newxD =  1; newyD =  0
            }
            if direction == 2 { // south
                newxD =  0; newyD =  1
            }
            if direction == 3 { // west
                newxD = -1; newyD =  0
            }
            walkCities(data, nextCityName, x+xD, y+yD, newxD, newyD, result)
        }
    }
}

// expand expands the 2D array in the desired direction
// only one direction supported at a time
func expand(result *[][]string, x int, y int) {
    if y == 1 || y == -1 {
        width := len((*result)[0])
        emptyLine := []string{" "}
        appendSpaces(&emptyLine, width-1)
        if y == 1 {
            *result = append(*result, emptyLine)
        } else {
            emptyArray := [][]string{emptyLine}
            *result = append(emptyArray, *result...)
        }
    } else if x == 1 {
        for i, line := range *result {
            line = append(line, " ")
            (*result)[i] = line
        }
    } else if x == -1 {
        for i, line := range *result {
            line = append([]string{" "}, line...)
            (*result)[i] = line
        }
    }
}

// appendSpaces extends the string array with the number of one-space string provided on the right
func appendSpaces(slice *[]string, num int) {
    for i:=0; i<num; i++ {
        *slice = append(*slice, " ")
    }
}

// prependSpaces extends the string array with the number of one-space string provided at left
func prependSpaces(slice *[]string, num int) {
    for i:=0; i<num; i++ {
        *slice = append([]string{" "}, *slice...)
    }
}

// firstCityName returns a random city name from data.Map hash, or an empty string
func firstCityName(data *invasion.Invasion) string {
    for cityName, _ := range data.Map {
        return cityName
    }

    return ""
}
