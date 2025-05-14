// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 58.
//!+

// Surface computes an SVG rendering of a 3-D surface function.
package main

import (
	"fmt"
	"io"
	"math"
)

const (
	width, height = 600, 320            // canvas size in pixels
	cells         = 100                 // number of grid cells
	xyrange       = 30.0                // axis ranges (-xyrange..+xyrange)
	xyscale       = width / 2 / xyrange // pixels per x or y unit
	zscale        = height * 0.4        // pixels per z unit
	angle         = math.Pi / 6         // angle of x, y axes (=30°)
)

var sin30, cos30 = math.Sin(angle), math.Cos(angle) // sin(30°), cos(30°)

func surface(w io.Writer, cor int) {
	fmt.Fprintf(w, "<svg xmlns='http://www.w3.org/2000/svg' "+
		"style='stroke: grey; stroke-width: 0.7' "+
		"width='%d' height='%d'>", width, height)
	for i := 0; i < cells; i++ {
		for j := 0; j < cells; j++ {
			ax, ay, z := corner(i+1, j)
			bx, by, _ := corner(i, j)
			cx, cy, _ := corner(i, j+1)
			dx, dy, _ := corner(i+1, j+1)

			// Calculate color based on z value
			color := getColor(z, cor)

			fmt.Fprintf(w, "<polygon points='%g,%g %g,%g %g,%g %g,%g' fill='%s'/>\n",
				ax, ay, bx, by, cx, cy, dx, dy, color)
		}
	}
	fmt.Fprintln(w, "</svg>")

}

func getColor(z float64, cor int) string {

	fmt.Printf("Debug: z=%v, cor=%v\n", z, cor) // Debug

	switch cor {
	case 1:
		if z > 0.5 {
			return "#ff0000"
		} else if z < 0.5 && z > 0 {
			return "#0000ff"
		} else {
			return "#00ff00"
		}

	case 2:

		if z > 0.5 {
			return "#fff200"
		} else if z < 0.5 && z > 0 {
			return "#00ff00"
		} else {
			return "#00ff00"
		}

	case 3:

		if z > 0.5 {
			return "#00fff2"
		} else if z < 0.5 && z > 0 {
			return "#00ffff"
		} else {
			return "#00ff00"
		}

	default:
		if z > 0.5 {
			// Peaks are light yellow
			return "#000000"
		} else if z < 0.5 && z > 0 {

			return "#0000ff"
		} else {
			// Valleys are green
			return "#00ff00"
		}
	}
	// Normalize z value to determine color

}

func corner(i, j int) (float64, float64, float64) {
	// Find point (x,y) at corner of cell (i,j).
	x := xyrange * (float64(i)/cells - 0.5)
	y := xyrange * (float64(j)/cells - 0.5)

	// Compute surface height z.
	z := f(x, y)

	// Project (x,y,z) isometrically onto 2-D SVG canvas (sx,sy).
	sx := width/2 + (x-y)*cos30*xyscale
	sy := height/2 + (x+y)*sin30*xyscale - z*zscale
	return sx, sy, z
}

func f(x, y float64) float64 {
	r := math.Hypot(x, y) // distance from (0,0)
	return math.Sin(r) / r
}

//!-
