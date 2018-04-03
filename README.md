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
go get -d -v ./...

go run cmd/alien-invasion/main.go --map examples/testmap1.txt --aliens 3
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

An example run on a 25x25 full map:

```
$ go run cmd/alien-invasion/main.go -map examples/testmap25x25.txt -aliens 51
2018/03/29 22:45:56 [iter     0] City_2_18 has been destroyed by alien 33 and alien 51
2018/03/29 22:45:56 [iter     1] City_6_10 has been destroyed by alien 4 and alien 50
2018/03/29 22:45:56 [iter     1] City_9_8 has been destroyed by alien 34 and alien 46
2018/03/29 22:45:56 [iter     1] City_18_3 has been destroyed by alien 3 and alien 16
2018/03/29 22:45:56 [iter     1] City_24_13 has been destroyed by alien 13 and alien 25
2018/03/29 22:45:56 [iter     1] City_22_19 has been destroyed by alien 38 and alien 48
2018/03/29 22:45:56 [iter     1] City_23_21 has been destroyed by alien 10 and alien 44
2018/03/29 22:45:56 [iter     3] City_17_14 has been destroyed by alien 40 and alien 24
2018/03/29 22:45:56 [iter     4] City_3_4 has been destroyed by alien 41 and alien 20
2018/03/29 22:45:56 [iter     4] City_7_11 has been destroyed by alien 18 and alien 22
2018/03/29 22:45:56 [iter     5] City_13_24 has been destroyed by alien 32 and alien 17
2018/03/29 22:45:56 [iter     5] City_9_6 has been destroyed by alien 15 and alien 9
2018/03/29 22:45:56 [iter     6] City_17_25 has been destroyed by alien 26 and alien 8
2018/03/29 22:45:56 [iter     6] City_21_14 has been destroyed by alien 19 and alien 42
2018/03/29 22:45:56 [iter     9] City_11_19 has been destroyed by alien 45 and alien 43
2018/03/29 22:45:56 [iter    14] City_5_17 has been destroyed by alien 27 and alien 7
2018/03/29 22:45:56 [iter    15] City_1_21 has been destroyed by alien 31 and alien 49
2018/03/29 22:45:56 [iter    21] City_24_9 has been destroyed by alien 23 and alien 35
2018/03/29 22:45:56 [iter    24] City_24_5 has been destroyed by alien 47 and alien 11
2018/03/29 22:45:56 [iter    26] City_20_22 has been destroyed by alien 14 and alien 6
2018/03/29 22:45:56 [iter    36] City_12_9 has been destroyed by alien 39 and alien 21
2018/03/29 22:45:56 [iter    48] City_11_8 has been destroyed by alien 36 and alien 1
2018/03/29 22:45:56 [iter    54] City_19_13 has been destroyed by alien 2 and alien 5
2018/03/29 22:45:56 [iter   107] City_16_4 has been destroyed by alien 30 and alien 12
2018/03/29 22:45:56 [iter   244] City_16_21 has been destroyed by alien 37 and alien 28
2018/03/29 22:45:57 [iter 10000] simulation terminated after 10000 iterations
2018/03/29 22:45:57 [iter 10000] Statistics:
2018/03/29 22:45:57 [iter 10000]    number of moves executed: 19159
2018/03/29 22:45:57 [iter 10000]    25 cities were destroyed out of 625, 600 remained
```

## Using a Docker container

The provided Dockerfile helps you run the simulation in a dedicated container holding all libraries needed.
An example Docker build and run:

```
$ docker build .
Sending build context to Docker daemon  233.5kB
Step 1/7 : FROM golang:1.10
 ---> d632bbfe5767
Step 2/7 : WORKDIR /go/src/app
 ---> Using cache
 ---> bbd7861fe64b
Step 3/7 : COPY . .
 ---> c5ea32f34268
Removing intermediate container 2b59875e044f
Step 4/7 : ENV GOPATH /go/src/app
 ---> Running in 941daed8540f
 ---> 1a8319738046
Removing intermediate container 941daed8540f
Step 5/7 : ENV GOBIN /go/src/app/bin
 ---> Running in cc7b136ac184
 ---> 757a8f3f795c
Removing intermediate container cc7b136ac184
Step 6/7 : RUN go get -d -v ./... && go install -v cmd/alien-invasion/main.go
 ---> Running in 528145003c6a
github.com/davecgh/go-spew (download)
github.com/kr/pretty (download)
github.com/kr/text (download)
github.com/davecgh/go-spew/spew
invasion
command-line-arguments
 ---> b7aad02d67c4
Removing intermediate container 528145003c6a
Step 7/7 : CMD bin/main
 ---> Running in 390a2db044c0
 ---> 262837db3f94
Removing intermediate container 390a2db044c0
Successfully built 262837db3f94

$ docker run 262837db3f94
2018/04/03 09:06:55 [iter     0] Qu-ux has been destroyed by alien 1 and alien 2
2018/04/03 09:06:55 [iter 10000] simulation terminated after 10000 iterations
2018/04/03 09:06:55 [iter 10000] Statistics:
2018/04/03 09:06:55 [iter 10000]        number of moves executed: 10002
2018/04/03 09:06:55 [iter 10000]        1 city was destroyed out of 5, 4 remained
```

## Utilities

The util/ directory contains some utility scripts: a full-mesh **map generator** (you
can provide the width and height and it will generate a map full of cities, all
iter-connected):

```
$ go run util/map_generator.go -width 20 -height 10 -out my_10x20_map.txt
Written 200 lines into file my_10x20_map.txt
```

And a **map visualizer**, which is not yet functional in the current release, but will be fixed soon:

```
$ go run util/map_visualizer.go -map examples/testmap1.txt

[][]string{
    {"B", "Q", "B"},
    {" ", " ", "F"},
}
```

## Running the tests

Simple Go tests are provided for invasion package:

```
go test ./cmd/...
```

You can check the status of latest Travis builds at https://travis-ci.org/bert42/alien-invasion


## TODO

* code coverage report
* **dep** for library version dependency tracking
* concurrency for Alien moves
* map visualizer, simulation visualizer
* map generator option to _drop_ number of cities
* map generator option to leave out roads between cities (percent of roads missing?)

## Authors

* **Norbert Csongradi** - [Bert42](https://github.com/bert42)

## License

This project is licensed under the Apache License - see the [LICENSE.md](LICENSE.md) file for details
