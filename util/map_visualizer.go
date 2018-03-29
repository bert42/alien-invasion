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

//	fmt.Println(visualData)
}

var cache map[string]bool

func visualize(data *invasion.Invasion) (result [][]string) {
	cache = make(map[string]bool, 0)

	result = [][]string{{""}}
	walkCities(data, firstCityName(data), 0, 0, 0, 0, &result)

	return result
}

func walkCities(data *invasion.Invasion, cityName string, x int, y int, xD int, yD int, result *[][]string) {
	fmt.Println("called: ", cityName, x, y, xD, yD)
	if _, ok := cache[cityName]; ok { // city already processed
		return
	}

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
//	x -= xD
//	y -= yD

	fmt.Println("processing: ", *result, x, y)

	// store City at current coordinates, then walk around in directions
	(*result)[x][y] = cityName[0:1]

	for direction := 0; direction < 4; direction++ {
        if nextCityName, toOk := city.Roads[direction]; toOk {
            if direction == 0 {
            	xD =  0; yD = -1
            }
            if direction == 1 {
            	xD =  1; yD = 0
            }
            if direction == 2 {
            	xD =  0; yD = 1
            }
            if direction == 3 {
            	xD =  -1; yD = 0
            }
        	walkCities(data, nextCityName, x, y, xD, yD, result)
        }
    }
}

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

func appendSpaces(slice *[]string, num int) {
	for i:=0; i<num; i++ {
		*slice = append(*slice, " ")
	}
}

func prependSpaces(slice *[]string, num int) {
	for i:=0; i<num; i++ {
		*slice = append([]string{" "}, *slice...)
	}
}

func firstCityName(data *invasion.Invasion) string {
   	for cityName, _ := range data.Map {
   		return cityName
    }

    return ""
}
