package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/tdewolff/canvas"
)

func loadPaths(name string) ([]canvas.Paths, error) {
	var ps []canvas.Paths
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

		p := canvas.Paths{}
		for _, polygon := range polygons {
			if len(polygon) == 0 {
				continue
			}

			pi := &canvas.Path{}
			pi.MoveTo(polygon[0][0], polygon[0][1])
			for _, coord := range polygon[1:] {
				pi.LineTo(coord[0], coord[1])
			}
			pi.Close()
			p = append(p, pi)
		}
		ps = append(ps, p)
	}
	return ps, nil
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
		ps, status, dur := exec(europes[z], chiles[z])

		fss := [][][][2]float64{}
		for _, p := range ps {
			fs := [][][2]float64{}
			for _, pi := range p.Split() {
				coords := pi.Coords()
				f := make([][2]float64, len(coords))
				for i, c := range coords {
					f[i][0], f[i][1] = c.X, c.Y
				}
				fs = append(fs, f)
			}
			fss = append(fss, fs)
		}

		w.Encode(struct {
			Name   string
			Z      int
			T      time.Duration
			Status int
			Result [][][][2]float64
		}{"tdewolff", z, dur, status, fss})
		fmt.Println(z, dur)

		if status != 0 {
			break
		}
	}
}

func exec(europe, chile canvas.Paths) (canvas.Paths, int, time.Duration) {
	var p canvas.Paths
	var status int
	var dur time.Duration
	defer func() {
		if r := recover(); r != nil {
			p = nil
			status = 3
			dur = 0
		}
	}()

	var n int
	t := time.Now()
	for ; n < 10; n++ {
		if 10*time.Second < time.Since(t) {
			return nil, 2, 0
		}

		p = europe.Or(chile)
	}
	dur = time.Since(t) / time.Duration(n)
	return p, status, dur
}
