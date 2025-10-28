package main

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/jonas-p/go-shp"
	"github.com/wroge/wgs84/v2"

	"github.com/tdewolff/canvas"
)

func main() {
	r, err := shp.OpenZip("ne_10m_admin_0_countries.zip")
	if err != nil {
		panic(err)
	}
	defer r.Close()

	europeCountries := []string{"Norway", "Sweden", "Finland", "Estonia", "Latvia", "Lithuania", "Denmark", "United Kingdom", "Ireland", "Netherlands", "Belgium", "Luxembourg", "Germany", "France", "Poland", "Austria", "Switzerland", "Spain", "Andorra", "Portugal", "Slovenia", "Czechia", "Slovakia", "Croatia", "Bosnia and Herzegovina", "Montenegro", "Kosovo", "Republic of Serbia", "North Macedonia", "Albania", "Hungary", "Bulgaria", "Romania", "Greece", "Italy", "Malta"}

	var chile, europe *canvas.Path
	countries := []string{}
	for r.Next() {
		if err := r.Err(); err != nil {
			panic(err)
		}

		name := r.Attribute(3)
		if null := strings.IndexByte(name, 0); null != -1 {
			name = name[:null]
		}
		countries = append(countries, name)
		if name == "Chile" {
			_, ishape := r.Shape()
			shape, err := ParseShapePath(ishape)
			if err != nil {
				panic(err)
			}
			chile = chile.Append(shape)
		} else if slices.Contains(europeCountries, name) {
			_, ishape := r.Shape()
			shape, err := ParseShapePath(ishape)
			if err != nil {
				panic(err)
			}
			europe = europe.Append(shape)
		}
	}

	// remove islands (coordinates in WGS84)
	chile = chile.Clip(-78, -60, -62, -16)
	europe = europe.Clip(-12, 30, 32, 72)

	// transform Chile to UTM 19 south, this has the least distortion for Chile
	utm19S := wgs84.Transform(wgs84.EPSG(4326), wgs84.EPSG(32719))
	chile = chile.TransformFunc(func(x, y float64) (float64, float64) {
		x, y, _ = utm19S(x, y, 0.0)
		return x / 1e5, y / 1e5
	})

	// transform Europe to UTM 33 north, this has the least distortion for Norway/Italy
	utm33N := wgs84.Transform(wgs84.EPSG(4326), wgs84.EPSG(32633))
	europe = europe.TransformFunc(func(x, y float64) (float64, float64) {
		x, y, _ = utm33N(x, y, 0.0)
		return x / 1e5, y / 1e5
	})

	if err := os.MkdirAll("data", 0755); err != nil {
		panic(err)
	} else if err := writeFiles("data/europe", europe); err != nil {
		panic(err)
	} else if err := writeFiles("data/chile", chile); err != nil {
		panic(err)
	}
}

func writeFiles(name string, p *canvas.Path) error {
	zs := []float64{1.0, 0.35, 0.1, 0.03, 0.01, 0.0035, 0.001, 0.0002, 0.00003, 0.0}
	for z, vwSize := range zs {
		p2 := p
		if vwSize != 0.0 {
			p2 = p.SimplifyVisvalingamWhyatt(vwSize)
		}
		name2 := fmt.Sprintf("%s_%d", name, z)
		fmt.Printf("%s: z=%d vw-size=%g len=%d\n", name2, z, vwSize, p2.Len())

		//if w, err := os.Create(name2 + ".path.gz"); err != nil {
		//	return err
		//} else {
		//	wGzip := gzip.NewWriter(w)
		//	if err := gob.NewEncoder(wGzip).Encode(p2); err != nil {
		//		return err
		//	} else if err := wGzip.Close(); err != nil {
		//		return err
		//	} else if err := w.Close(); err != nil {
		//		return err
		//	}
		//}

		if w, err := os.Create(name2 + ".json"); err != nil {
			return err
		} else {
			fs := [][][2]float64{}
			for _, pi := range p2.Split() {
				coords := pi.Coords()
				f := make([][2]float64, len(coords))
				for i, c := range coords {
					f[i][0], f[i][1] = c.X, c.Y
				}
				fs = append(fs, f)
			}
			if err := json.NewEncoder(w).Encode(fs); err != nil {
				return err
			} else if err := w.Close(); err != nil {
				return err
			}
		}
	}
	return nil
}

type Projection func(float64, float64) (float64, float64)

func ParseShapePath(ishape shp.Shape) (*canvas.Path, error) {
	var d []float64
	switch shape := ishape.(type) {
	case *shp.Polygon:
		d = make([]float64, 0, (len(shape.Parts)+1)*4)
		for i := 0; i < len(shape.Parts); i++ {
			start := shape.Parts[i]
			end := int32(len(shape.Points))
			if i+1 < len(shape.Parts) {
				end = shape.Parts[i+1]
			}
			if end-start < 2 {
				continue
			}

			d = append(d, canvas.MoveToCmd, shape.Points[start].X, shape.Points[start].Y, canvas.MoveToCmd)
			for _, pt := range shape.Points[start+1 : end] {
				d = append(d, canvas.LineToCmd, pt.X, pt.Y, canvas.LineToCmd)
			}
			d = append(d, canvas.CloseCmd, shape.Points[start].X, shape.Points[start].Y, canvas.CloseCmd)
		}
	case *shp.PolyLine:
		d = make([]float64, 0, (len(shape.Parts)+1)*4)
		for i := 0; i < len(shape.Parts); i++ {
			start := shape.Parts[i]
			end := int32(len(shape.Points))
			if i+1 < len(shape.Parts) {
				end = shape.Parts[i+1]
			}
			if end-start < 2 {
				continue
			}

			d = append(d, canvas.MoveToCmd, shape.Points[start].X, shape.Points[start].Y, canvas.MoveToCmd)
			for _, pt := range shape.Points[start+1 : end] {
				d = append(d, canvas.LineToCmd, pt.X, pt.Y, canvas.LineToCmd)
			}
		}
	default:
		return nil, fmt.Errorf("unknown shape type: %T", ishape)
	}
	return canvas.NewPathFromData(d), nil
}
