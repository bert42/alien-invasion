package main

import (
    "fmt"
    "flag"
    "bufio"
    "os"
    "strconv"
    "strings"
)

const Prefix = "City_"

func main() {
    flag.Usage = func() {
        fmt.Println(`Generates full-connected mash map of given size`)

        flag.PrintDefaults()
    }

    outFile := flag.String("out", "", "out filename (required)")
    width   := flag.Int("width", 5, "map width")
    height  := flag.Int("height", 5, "map height")

    flag.Parse()

    if *outFile == "" {
        fmt.Println("--out is required")
        os.Exit(2)
    }

    file, err := os.OpenFile(*outFile, os.O_RDWR|os.O_CREATE, 0755)
    if err != nil {
        fmt.Printf("Unable to open file %s: %s\n", outFile, err.Error())
    }
    defer file.Close()

    writer := bufio.NewWriter(file)
    defer writer.Flush()

    lines := 0
    for i:=1; i<=*width; i++ {
        for j:=1; j<=*height; j++ {
            elems := []string{CityName(i, j)}
            if i < *width {
                elems = append(elems, "east="+CityName(i+1, j))
            }
            if i > 1 {
                elems = append(elems, "west="+CityName(i-1, j))
            }
            if j > 1 {
                elems = append(elems, "north="+CityName(i, j-1))
            }
            if j < *height {
                elems = append(elems, "south="+CityName(i, j+1))
            }

            line := strings.Join(elems, " ")+"\n"
            lines++
            _, err := writer.WriteString(line)
            if err != nil {
                fmt.Printf("Unable to write into file: %s", err.Error())
            }
        }
    }

    fmt.Printf("Written %d lines into file %s\n", lines, *outFile)
}

func CityName(i int, j int) string {
    return Prefix+strconv.FormatInt(int64(i), 10)+"_"+strconv.FormatInt(int64(j), 10)
}