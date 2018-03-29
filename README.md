# Alien-invasion

[![license](https://img.shields.io/github/license/bert42/alien-invasion.svg)](https://github.com/bert42/alien-invasion/blob/master/LICENSE)
[![Build Status](https://travis-ci.org/bert42/alien-invasion.svg?branch=master)](https://travis-ci.org/bert42/alien-invasion)

A simple iterative simulation of a fictional alien invasion in Go

## Prerequisites

You need to install git (https://git-scm.com/downloads) and Go (https://golang.org/dl/) first.
Also, set the current GOPATH to the current dir and let Go install required libraries by running (assuming a Linux/Unix environment):

```
git clone https://github.com/bert42/alien-invasion
cd alien-invasion
export GOPATH=`pwd`
go get ./...
```

## Getting started

This project simulates an alien invasion on a given map of cities.
You can find some example map files under examples/, feel free to create your own.
Map files are simple text files with format:
 * one city per line
 * city name is first, followed by 1-4 directions (either north, east, south of west)
 * each direction represents a road to another city that lies in that direction

Example map file content:

```
Foo_city north=Bar_city
Bar_city south=Foo_city
```

Map files are validated to be complete and will be refused for missing roads.


You can start the simulation by specifying a map file to be used and the number of aliens to be deployed:

```
go run cmd/alien-invasion/main.go --map examples/testmap1.txt --aliens 3
```

It will randomly drop 3 aliens in the cities. They start wandering around using the roads between cities. Whenever two aliens enter
the same city, they start a fight and in the process kill each other and destroy the city. When a city is destroyed, it is removed
from the map along with any roads to and from it.
Simulation runs until 10.000 iterations, or until all aliens are destroyed.

The program will print a line when a city is destroyed (please note the iteration number), like:

```
2018/03/25 20:26:55 [iter     0] Bar has been destroyed by alien 1 and alien 2
```

## Running the tests

Simple Go tests are provided:

```
go test -v cmd/alien-invasion/*test.go
```

## TODO

* **dep** for library version dependency tracking
* concurrency for Alien moves
* map visualizer, simulation visualizer
* map generator option to _drop_ number of cities
* map generator option to leave out roads between cities (percent of roads missing?)

## Authors

* **Norbert Csongradi** - [Bert42](https://github.com/bert42)

## License

This project is licensed under the Apache License - see the [LICENSE.md](LICENSE.md) file for details
