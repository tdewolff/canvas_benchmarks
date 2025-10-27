package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/tdewolff/canvas"
	"github.com/tdewolff/canvas/renderers"
)

func loadPath(name string) (*canvas.Path, error) {
	f, err := os.Open(fmt.Sprintf("../data/%s_8.json", name))
	if err != nil {
		return nil, err
	}

	polygons := [][][2]float64{}
	if err := json.NewDecoder(f).Decode(&polygons); err != nil {
		return nil, err
	}

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
	return p, nil
}

func main() {
	europe, err := loadPath("europe")
	if err != nil {
		panic(err)
	}

	bounds := canvas.Rect{9.5, 59.5, 11.0, 61.0}

	//europe = europe.Settle(canvas.NonZero)
	europe = europe.Or(canvas.Rect{0.0, 0.0, 1.0, 1.0}.ToPath())

	c := canvas.New(bounds.W(), bounds.H())
	ctx := canvas.NewContext(c)

	// background
	ctx.SetFillColor(canvas.White)
	ctx.DrawPath(0, 0, canvas.Rectangle(c.W, c.H))

	// result
	ctx.SetStrokeWidth(0.015)
	ctx.SetStrokeColor(canvas.Black)
	ctx.SetFillColor(canvas.Hex("#0001"))

	europe.Translate(-bounds.X0, -bounds.Y0)
	ctx.DrawPath(0.0, 0.0, europe)

	renderers.Write("out.png", c, canvas.DPMM(200.0))
}
