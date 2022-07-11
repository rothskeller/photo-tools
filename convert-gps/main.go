package main

import (
	"fmt"
	"math"
)

func main() {
	for {
		var lat, long float64
		fmt.Scanf("%f, %f\n", &lat, &long)
		var d, m float64
		d = math.Floor(lat)
		lat = lat - d
		lat = lat * 60
		m = math.Floor(lat)
		lat = lat - m
		lat = lat * 60
		fmt.Printf("%.0f°%.0f'%f\", ", d, m, lat)
		d = math.Floor(-long)
		long = -long - d
		long = long * 60
		m = math.Floor(long)
		long = long - m
		long = long * 60
		fmt.Printf("%.0f°%.0f'%f\"\n", d, m, long)
	}
}
