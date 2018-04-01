package invasion

import (
    "testing"
    "os"
    "sort"
)

const mapfile1 = "../../../examples/testmap1.txt"

var simulation *Invasion


func TestBuildMap(t *testing.T) {
    if _, err := os.Stat(mapfile1); os.IsNotExist(err) {
        t.Fatalf("Unable to open test map %s", mapfile1)
    }

    simulation = New()
    simulation.BuildMap(mapfile1)

    cities := simulation.AllCities()
    if len(cities) != 5 {
        t.Fatalf("Number of loaded cities differ, got: %d, want: %d", len(cities), 5)
    }
}

func TestDirections(t *testing.T) {
    city := City{
        Name: "Foobar",
        Roads: map[int]string{1: "Blah", 3: "Huhh"},
        Alien: 0,
    }

    roads := AllRoads(&city)
    if len(roads) != 2 {
        t.Fatalf("Number of roads differ, got: %d, want: %d", len(roads), 2)
    }

}

func TestValidateRoads(t *testing.T) {
    simulation.ValidateRoads()
    t.Log("ValidateRoads ran ok") // well, not the best candidate for testing
}

func TestDeploy(t *testing.T) {
    simulation.Deploy(1) // should not destroy any cities and move is possible

    aliens := citiesWithAliens(simulation)

    if len(aliens) != 1 {
        t.Fatalf("Deploy failed, got: %d, want: %d", len(aliens), 1)
    }
}

func TestMove(t *testing.T) {
    aliensBefore := citiesWithAliens(simulation)
t.Log(aliensBefore)
    simulation.Move()
    t.Log("Move ran ok")

    aliensAfter := citiesWithAliens(simulation)
t.Log(aliensAfter)
    if compareSlices(aliensBefore, aliensAfter) {
        t.Fatal("Simulation Move() left aliens in the same cities, should not happen")
    }
}

func TestDestroyCity(t *testing.T) {
    simulation.DestroyCity("Bee", 9999, 9998)
    cities := simulation.AllCities()
    if len(cities) != 4 {
        t.Fatalf("Number of cities after Destroy differ, got: %d, want: %d", len(cities), 4)
    }
}


// utility functions

// citiesWithAliens returns all city names with Aliens in it
func citiesWithAliens(simulation *Invasion) (result []string) {
    for _, cityName := range simulation.AllCities() {
        if alienID := simulation.Map[cityName].Alien; alienID != 0 {
            result = append(result, cityName)
        }
    }

    return result
}

// compareSlices returns true for two string-slices with same data, false for different data
func compareSlices(a, b []string) bool {
    if a == nil && b == nil {
        return true;
    }

    if a == nil || b == nil {
        return false;
    }

    if len(a) != len(b) {
        return false
    }

    // needs sorting to compare if all fit yet
    sort.Strings(a)
    sort.Strings(b)

    for i := range a {
        if a[i] != b[i] {
            return false
        }
    }

    return true
}
