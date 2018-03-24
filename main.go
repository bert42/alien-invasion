package main

import (
	"fmt"
	"log"
	"os"
	"io"
	"bufio"
)


type City struct {
	Name	string
	North	string
	East	string
	South	string
	West	string
}

func main() {

	mapData := readMap()
	fullMap := buildMap(mapData)

	fmt.Println(fullMap)
	fmt.Println("Done.")
}


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

		result = append(result, line)
	}

	return result
}

func buildMap(mapData []string) []City {
	var cities []City;

	for _, line := range mapData {
		city := City{
			Name: line,
		}

		cities = append(cities, city)
	}

	return cities
}