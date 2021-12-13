package main

import (
	"fmt"
	"math"
	"os"
	"strings"
	"flag"
	"bufio"
	"strconv"
)

type Airport struct {
	identity string
	iata string
	lat float64
	lon float64
}

func main() {

	fmt.Println("Started");

	distFile := "./dist/distances.csv"
	deleteIfFileExists(distFile)

	writeToFile(distFile, fmt.Sprintf("%s;%s;%s", "from", "to", "distanceInKM"))

	fmt.Println("Init end");

	ports := getAllAirports()
	
	fmt.Printf("Running loop for %d airpots\n", len(ports))

	totalExpected := ((len(ports) * (len(ports) - 1)) / 2)

	fmt.Printf("Will make %d calculations\n", totalExpected)

	totalDone := 0

	for x := 0; x<len(ports); x++ {
		for y := x+1; y<len(ports); y++ {
			distance := measureDistanceInKm(ports[x].lat, ports[x].lon, ports[y].lat, ports[y].lon)
			writeToFile(distFile, fmt.Sprintf("%s;%s;%f", ports[x].identity, ports[y].identity, distance))
			totalDone ++;
		} 
	}

	fmt.Printf("Total Done : %d\n", totalDone)
}

func getAllAirports() []Airport {
	fptr := flag.String("fpath", "./source/airport-codes_csv.csv", "Read from")
	flag.Parse()

	f, err := os.Open(*fptr)
	if err != nil {
		panic(err)
	}

	defer func() {
		f.Close()
	}()

	s := bufio.NewScanner(f)
	i := 0

	ports := []Airport{}
	
	for (s.Scan()) {
		i++
		if i == 1 {
			continue
		}
		place := strings.Split(s.Text(), ";")

		coord1str := strings.Replace(place[11], ",", ".", -1)
		coord1, err := strconv.ParseFloat(coord1str, 64);
		if err != nil {
			panic(err)
		}

		coord2str := strings.Replace(place[12], ",", ".", -1)
		coord2, err := strconv.ParseFloat(coord2str, 64);
		if err != nil {
			panic(err)
		}
		
		ports = append(ports, Airport{place[0], place[9], coord2, coord1})

	}

	return ports;
}

func measureDistanceInKm(lat1 float64, lon1 float64, lat2 float64, lon2 float64) float64 {
	R := 6371
	dLat := deg2rad(lat2-lat1)
	dLon := deg2rad(lon2-lon1)

	a := math.Sin(dLat / 2) * math.Sin(dLat/2) +
    math.Cos(deg2rad(lat1)) * math.Cos(deg2rad(lat2)) * 
    math.Sin(dLon/2) * math.Sin(dLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return c * float64(R)
}

func deg2rad(deg float64) float64 {
	return deg * (math.Pi / 180)
}

func deleteIfFileExists(dest string) {
	err := os.Remove(dest)
	if err != nil {
		fmt.Println(err)
	}
}

func writeToFile(dest string, text string) {
	// If the file doesn't exist, create it, or append to the file
	f, err := os.OpenFile(dest, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
	}
	if _, err := f.Write([]byte(text + "\n")); err != nil {
		f.Close() // ignore error; Write error takes precedence
		fmt.Println(err)
	}
	if err := f.Close(); err != nil {
		fmt.Println(err)
	}
}