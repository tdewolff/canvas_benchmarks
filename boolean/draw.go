package main

import (
	"encoding/json"
	"fmt"
	"image/color"
	"io"
	"math"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/tdewolff/canvas"
	"github.com/tdewolff/canvas/renderers"
)

type Result struct {
	Name   string
	Z      int
	T      time.Duration
	Status int
	Result [][][][2]float64
}

func loadLengths(name string) ([]int, error) {
	var ns []int
	for z := 0; ; z++ {
		f, err := os.Open(fmt.Sprintf("data/%s_%d.json", name, z))
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

		n := 0
		for _, polygon := range polygons {
			n += len(polygon)
		}
		ns = append(ns, n)
	}
	return ns, nil
}

func main() {
	f, err := os.Open("data/results.json")
	if err != nil {
		panic(err)
	}

	var results []Result
	r := json.NewDecoder(f)
	for {
		var result Result
		if err := r.Decode(&result); err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		} else if result.Z != 8 {
			results = append(results, result)
			continue
		}

		bounds := canvas.Rect{-17.507566112883165, 37.01773388676248, 17.904232351182134, 81.64329738477107}

		c := canvas.New(bounds.W(), bounds.H())
		ctx := canvas.NewContext(c)

		// background
		ctx.SetFillColor(canvas.White)
		ctx.DrawPath(0, 0, canvas.Rectangle(c.W, c.H))

		// result
		ctx.SetStrokeWidth(0.015)
		ctx.SetStrokeColor(canvas.Black)
		ctx.SetFillColor(canvas.Hex("#0001"))
		for _, polygons := range result.Result {
			p := &canvas.Path{}
			for _, polygon := range polygons {
				if len(polygon) == 0 {
					continue
				}

				p.MoveTo(polygon[0][0], polygon[0][1])
				for _, coord := range polygon[1:] {
					p.LineTo(coord[0], coord[1])
				}
				p.Close()
			}
			p.Translate(-bounds.X0, -bounds.Y0)
			ctx.DrawPath(0.0, 0.0, p)
		}

		renderers.Write(result.Name+".png", c, canvas.DPMM(20.0))
		results = append(results, result)
	}

	slices.SortFunc(results, func(a, b Result) int {
		cmp := strings.Compare(a.Name, b.Name)
		if cmp != 0 {
			return cmp
		}
		if a.Z < b.Z {
			return -1
		} else if b.Z < a.Z {
			return 1
		}
		return 0
	})

	maxDur := time.Duration(0)
	names := []string{}
	experiments := map[string][]Result{}
	for _, result := range results {
		rows, ok := experiments[result.Name]
		if !ok {
			names = append(names, result.Name)
		}
		rows = append(rows, result)
		experiments[result.Name] = rows

		if maxDur < result.T {
			maxDur = result.T
		}
	}
	slices.Sort(names)

	europeSegs, err := loadLengths("europe")
	chileSegs, err := loadLengths("chile")
	if len(europeSegs) != len(chileSegs) {
		panic("europe/chile Z levels must match")
	} else if len(europeSegs) == 0 {
		panic("no data")
	}

	Z := len(europeSegs)
	Width := 160.0
	Height := 80.0

	xs := []float64{}
	maxSegs := float64(europeSegs[len(europeSegs)-1] + chileSegs[len(chileSegs)-1])
	for z := 0; z < Z; z++ {
		xs = append(xs, Width*float64(europeSegs[z]+chileSegs[z])/maxSegs)
	}

	font, err := canvas.LoadSystemFont("Atkinson Hyperlegible,sans", canvas.FontRegular)
	if err != nil {
		panic(err)
	}
	face := font.Face(12.0, canvas.Black)
	metrics := face.Metrics()

	colors := []color.RGBA{
		canvas.Hex("#001f3f"),
		canvas.Hex("#f012be"),
		canvas.Hex("#ff851b"),
		canvas.Hex("#0074d9"),
		canvas.Hex("#7fdbff"),
		canvas.Hex("#ffdc00"),
		canvas.Hex("#39cccc"),
		canvas.Hex("#b10dc9"),
		canvas.Hex("#85144b"),
		canvas.Hex("#ff4136"),
		canvas.Hex("#3d9970"),
		canvas.Hex("#2ecc40"),
		canvas.Hex("#01ff70"),
	}

	shapes := []*canvas.Path{
		canvas.Circle(1.0),
		canvas.Triangle(1.4),
		canvas.Triangle(1.4).Rotate(180.0),
		canvas.MustParseSVGPath("M-1 -1L1 1M-1 1L1 -1").Stroke(0.4, canvas.ButtCap, canvas.MiterJoin, 0.01),
		canvas.MustParseSVGPath("M-1.4 0L1.4 0M0 -1.4L0 1.4").Stroke(0.4, canvas.ButtCap, canvas.MiterJoin, 0.01),
		canvas.RegularPolygon(4, 1.4, true),
		canvas.RegularPolygon(6, 1.4, false),
	}

	c := canvas.New(200.0, 100.0)
	ctx := canvas.NewContext(c)

	// background
	ctx.SetFillColor(canvas.White)
	ctx.DrawPath(0, 0, canvas.Rectangle(c.W, c.H))

	ctx.SetFillColor(canvas.Transparent)
	ctx.SetStrokeColor(canvas.Hex("#bbb"))
	ctx.SetStrokeWidth(0.5)
	ctx.DrawPath(35.0, 15.0, canvas.Line(0.0, Height))
	ctx.DrawPath(35.0, 15.0, canvas.Line(Width, 0.0))

	dw := math.Pow10(int(math.Round(math.Log10(maxSegs)))) / 4.0
	if float64(maxSegs) < dw*3 {
		dw /= 2.0
	}
	htick := canvas.Line(0.0, -2.0)
	for i := 0; i < 5; i++ {
		x := dw * float64(i) / maxSegs * Width
		if Width < x {
			break
		}
		ctx.DrawPath(35.0+x, 15.0, htick)
		ctx.DrawText(35.0+x, 7.5, canvas.NewTextLine(face, fmt.Sprintf("%d", i*int(dw)), canvas.Center))
	}

	dh := math.Pow10(int(math.Round(math.Log10(float64(maxDur))))) / 4.0
	if float64(maxDur) < dh*3 {
		dh /= 2.0
	}
	vtick := canvas.Line(-2.0, 0.0)
	for i := 0; i <= 5; i++ {
		y := dh * float64(i) / float64(maxDur) * Height
		if Height < y {
			break
		}
		ctx.DrawPath(35.0, 15.0+y, vtick)
		ctx.DrawText(30.0, 15.0+y-metrics.XHeight/2.0, canvas.NewTextLine(face, fmt.Sprintf("%v", time.Duration(float64(i)*dh)), canvas.Right))
	}

	for j, name := range names {
		p := &canvas.Polyline{}
		for i, run := range experiments[name] {
			p.Add(xs[i], Height*float64(run.T)/float64(maxDur))
		}
		ctx.SetFillColor(canvas.Transparent)
		ctx.SetStrokeColor(colors[j])
		ctx.SetStrokeWidth(0.5)
		ctx.DrawPath(35.0, 15.0, p.ToPath())

		ctx.SetFillColor(colors[j])
		ctx.SetStrokeColor(canvas.Transparent)
		for i, run := range experiments[name] {
			ctx.DrawPath(35.0+xs[i], 15.0+Height*float64(run.T)/float64(maxDur), shapes[j])
		}
		ctx.DrawPath(40.0, 90.0-metrics.LineHeight*float64(j)+metrics.XHeight/2.0, shapes[j])
	}

	ctx.DrawText(44.0, 90.0, canvas.NewTextLine(face, strings.Join(names, "\n"), canvas.Left))

	renderers.Write("results.png", c, canvas.DPMM(10.0))

	links := map[string]string{
		"ajohnson1": "http://www.angusj.com/delphi/clipper/documentation/Docs/Overview/_Body.htm",
		"ajohnson2": "https://github.com/AngusJohnson/Clipper2",
		"ioverlay":  "https://github.com/iShape-Rust/iOverlay",
		"tdewolff":  "https://github.com/tdewolff/canvas",
	}

	fmt.Println("Input:", maxSegs)
	fmt.Println()
	fmt.Println("| Library | Time | Polygon size |")
	fmt.Println("| --- | --- | --- |")
	for _, name := range names {
		fmt.Printf("| [%v](%v) | %v | %v |\n", name, links[name], experiments[name][Z-1].T, countSegments(experiments[name][Z-1].Result))
	}
}

func countSegments(polygons [][][][2]float64) int {
	n := 0
	for _, polygon := range polygons {
		for _, ring := range polygon {
			n += len(ring)
		}
	}
	return n
}
