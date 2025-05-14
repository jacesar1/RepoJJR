// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// Run with "web" command-line argument for web server.
// See page 13.
//!+main

// Lissajous generates GIF animations of random Lissajous figures.
package main

import (
	"image"
	"image/color"
	"image/gif"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

//!-main
// Packages not needed by version in book.

//!+main

var palette = []color.Color{color.White, color.Black}

//const (
//	whiteIndex = 0 // first color in palette
//	blackIndex = 1 // next color in palette
//)

var r *rand.Rand = rand.New(rand.NewSource(time.Now().UTC().UnixNano()))

func main() {
	//!-main
	// The sequence of images is deterministic unless we seed
	// the pseudo-random number generator using the current time.
	// Thanks to Randall McPherson for pointing out the omission.

	// Then inside the animation loop where colors are picked:
	for i := 0; i < len(palette); i++ {
		palette[i] = color.RGBA{uint8(r.Uint32() & 0xFF), uint8(r.Uint32() & 0xFF), uint8(r.Uint32() & 0xFF), 255}
	}

	if len(os.Args) > 1 && os.Args[1] == "web" {
		//!+http
		handler := func(w http.ResponseWriter, r *http.Request) {
			if err := r.ParseForm(); err != nil {
				log.Print(err)
			}
			cycles := 5
			for k, v := range r.Form {
				if k == "cycles" {
					c, err := strconv.Atoi(v[0])
					if err != nil {
						log.Print(err)
					} else {
						cycles = c
					}
				}
			}
			lissajous(w, cycles)
		}
		http.HandleFunc("/", handler)
		//!-http
		log.Fatal(http.ListenAndServe("localhost:8000", nil))
	} else {
		//!+main
		f, err := os.Create("output.gif")
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		lissajous(f, 5)
	}
}

func lissajous(out io.Writer, cycles int) {
	const (
		res     = 0.001 // angular resolution
		size    = 200   // image canvas covers [-size..+size]
		nframes = 64    // number of animation frames
		delay   = 8     // delay between frames in 10ms units
	)
	freq := rand.Float64() * 3.0 // relative frequency of y oscillator
	anim := gif.GIF{LoopCount: nframes}
	phase := 0.0 // phase difference
	for i := 0; i < nframes; i++ {
		rect := image.Rect(0, 0, 2*size+1, 2*size+1)
		img := image.NewPaletted(rect, palette)
		for t := 0.0; t < float64(cycles)*2*math.Pi; t += res {
			x := math.Sin(t)
			y := math.Sin(t*freq + phase)
			img.SetColorIndex(size+int(x*size+0.5), size+int(y*size+0.5),
				uint8(i%len(palette)))
		}
		phase += 0.1
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}
	gif.EncodeAll(out, &anim) // NOTE: ignoring encoding errors
}

//!-main
