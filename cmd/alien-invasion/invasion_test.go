package invasion_test

import (
    "testing"
    "os"

    "./invasion"
)

const mapfile1 = "../../examples/testmap1.txt"

func TestBuildMap(t *testing.T) {
    if _, err := os.Stat(mapfile1); os.IsNotExist(err) {
        t.Fatalf("Unable to open test map %s", mapfile1)
    }

    simulation := &invasion.Invasion{}
    simulation.BuildMap(mapfile1)

    cities := simulation.AllCities()
    if len(cities) != 5 {
        t.Fatalf("Number of loaded cities differ, got: %d, want: %d", len(cities), 5)
    }
}

func TestDirections(t *testing.T) {
    city := invasion.City{
        Name: "Foobar",
        Roads: map[int]string{1: "Blah", 3: "Huhh"},
        Alien: 0,
    }

    roads := invasion.AllRoads(&city)
    if len(roads) != 2 {
        t.Fatalf("Number of roads differ, got: %d, want: %d", len(roads), 2)
    }

}