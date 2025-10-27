package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/epit3d/goclipper2/goclipper2"
)

func loadPaths(name string) ([]*goclipper2.PathsD, error) {
	var zpaths []*goclipper2.PathsD
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

		paths := goclipper2.NewPathsd()
		for _, polygon := range polygons {
			path := *goclipper2.NewPathd()
			for _, coord := range polygon {
				path.AddPoint(*goclipper2.NewPointD(coord[0], coord[1]))
			}
			paths.AddPath(path)
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
		for j := int64(0); j < paths.Length(); j++ {
			polygon := [][2]float64{}
			coords := paths.GetPath(j).ToPoints()
			for _, coord := range coords {
				polygon = append(polygon, [2]float64{coord.X(), coord.Y()})
			}
			polygons = append(polygons, polygon)
		}

		w.Encode(struct {
			Name   string
			Z      int
			T      time.Duration
			Status int
			Result [][][][2]float64
		}{"ajohnson2", z, dur, status, [][][][2]float64{polygons}})
		fmt.Println(z, dur)

		if status != 0 {
			break
		}
	}
}

func exec(europe, chile *goclipper2.PathsD) (*goclipper2.PathsD, int, time.Duration) {
	var paths *goclipper2.PathsD
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
	t := time.Now()
	for ; n < 5; n++ {
		if 10*time.Second < time.Since(t) {
			return nil, 2, 0
		}

		paths = europe.BooleanOp(goclipper2.Union, goclipper2.NonZero, chile, 8)
	}
	dur = time.Since(t) / time.Duration(n)
	return paths, status, dur
}
