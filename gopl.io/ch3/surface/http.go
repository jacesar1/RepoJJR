// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// Run with "web" command-line argument for web server.
// See page 13.
//!+main

// Lissajous generates GIF animations of random Lissajous figures.
package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
)

//!-main

func main() {
	if len(os.Args) > 1 && os.Args[1] == "web" {
		//!+http
		handler := func(w http.ResponseWriter, r *http.Request) {
			if err := r.ParseForm(); err != nil {
				log.Print(err)
			}
			cor := 3
			for k, v := range r.Form {
				if k == "cor" {
					c, err := strconv.Atoi(v[0])
					if err != nil {
						log.Print(err)
					} else {
						cor = c

					}
				}
			}
			w.Header().Set("Content-Type", "image/svg+xml")
			surface(w, cor)
		}
		http.HandleFunc("/", handler)
		//!-http
		log.Fatal(http.ListenAndServe("localhost:8000", nil))
	} else {
		//!+main
		f, err := os.Create("surface.html")
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		surface(f, 3)
	}
}

//!-main
