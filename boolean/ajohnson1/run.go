package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	clipper "github.com/ctessum/go.clipper"
)

const Factor = 1000000.0 // float64 to int factor

func loadPaths(name string) ([]clipper.Paths, error) {
	var zpaths []clipper.Paths
	for z := 0; ; z++ {
		f, err := os.Open(fmt.Sprintf("../data/%s_%d.json", name, z))
		if err != nil {
			if os.IsNotExist(err) {
				break
			}
			return nil, err
		}

		polygons := [][][2]float64{}
		if err := json.NewDecoder(f).Decode(&polygons); err != nil {
			return nil, err
		}

		paths := clipper.Paths{}
		for _, polygon := range polygons {
			path := clipper.Path{}
			for _, coord := range polygon {
				path = append(path, &clipper.IntPoint{
					clipper.CInt(coord[0] * Factor),
					clipper.CInt(coord[1] * Factor),
				})
			}
			paths = append(paths, path)
		}
		zpaths = append(zpaths, paths)
	}
	return zpaths, nil
}

func main() {
	europes, err := loadPaths("europe")
	if err != nil {
		panic(err)
	}
	chiles, err := loadPaths("chile")
	if err != nil {
		panic(err)
	}

	f, err := os.OpenFile("../data/results.json", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	w := json.NewEncoder(f)
	for z := 0; z < len(europes); z++ {
		paths, status, dur := exec(europes[z], chiles[z])

		polygons := [][][2]float64{}
		for _, path := range paths {
			polygon := [][2]float64{}
			for _, coord := range path {
				polygon = append(polygon, [2]float64{
					float64(coord.X) / Factor,
					float64(coord.Y) / Factor,
				})
			}
			polygons = append(polygons, polygon)
		}

		w.Encode(struct {
			Name   string
			Z      int
			T      time.Duration
			Status int
			Result [][][][2]float64
		}{"ajohnson1", z, dur, status, [][][][2]float64{polygons}})
		fmt.Println(z, dur)

		if status != 0 {
			break
		}
	}
}

func exec(europe, chile clipper.Paths) (clipper.Paths, int, time.Duration) {
	var paths clipper.Paths
	var status int
	var dur time.Duration
	defer func() {
		if r := recover(); r != nil {
			paths = nil
			status = 3
			dur = 0
		}
	}()

	var n int
	var succeeded bool

	t := time.Now()
	for ; n < 5; n++ {
		if 10*time.Second < time.Since(t) {
			return nil, 2, 0
		}

		c := clipper.NewClipper(clipper.IoStrictlySimple | clipper.IoReverseSolution)
		c.AddPaths(europe, clipper.PtSubject, true)
		c.AddPaths(chile, clipper.PtClip, true)
		paths, succeeded = c.Execute1(clipper.CtUnion, clipper.PftNonZero, clipper.PftNonZero)
		if !succeeded {
			return nil, 1, 0
		}
	}
	dur = time.Since(t) / time.Duration(n)
	return paths, status, dur
}
