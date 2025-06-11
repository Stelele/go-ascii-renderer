package main

import (
	"math"
)

func getShape() []Point {
	framesMark := 600

	if f < framesMark {
		return getTorus()
	}
	if f < 2*framesMark {
		return getSphere()
	}

	if f > 4*framesMark {
		f = 0
	}
	return getCube()
}

func getTorus() []Point {
	R1 := 1.
	R2 := 2.

	points := make([]Point, 0)

	for t := 0.; t < 2*math.Pi; t += 0.1 {
		ct := math.Cos(t)
		st := math.Sin(t)

		for p := 0.; p < 2*math.Pi; p += 0.01 {
			cp := math.Cos(p)
			sp := math.Sin(p)

			point := Point{}

			point.X = (R2 + R1*ct) * cp
			point.Y = R1 * st
			point.Z = (R2 + R1*ct) * sp

			point.Nx = ct * cp
			point.Ny = st
			point.Nz = ct * sp

			points = append(points, point)
		}
	}

	return points
}

func getSphere() []Point {
	R := 5.
	points := make([]Point, 0)

	for t := 0.; t < 2*math.Pi; t += 0.1 {
		ct := math.Cos(t)
		st := math.Sin(t)

		for p := 0.; p < 2*math.Pi; p += 0.01 {
			cp := math.Cos(p)
			sp := math.Sin(p)

			point := Point{}

			point.X = R * ct * cp
			point.Y = R * st
			point.Z = R * ct * sp

			point.Nx = ct * cp
			point.Ny = st
			point.Nz = ct * sp

			points = append(points, point)
		}
	}

	return points
}

func getCube() []Point {
	d := 2.
	points := make([]Point, 0)

	for i := -1 * d; i < d; i += 0.1 {
		for j := -1 * d; j < d; j += 0.1 {
			// left
			l := Point{}

			l.X = -1 * d
			l.Y = i //(rand.Float64()*2 - 1) * d
			l.Z = j //(rand.Float64()*2 - 1) * d

			l.Nx = -1
			l.Ny = 0
			l.Nz = 0

			points = append(points, l)

			// right
			r := Point{}

			r.X = d
			r.Y = i //(rand.Float64()*2 - 1) * d
			r.Z = j //(rand.Float64()*2 - 1) * d

			r.Nx = 1
			r.Ny = 0
			r.Nz = 0

			points = append(points, r)

			// up
			u := Point{}

			u.X = i //(rand.Float64()*2 - 1) * d
			u.Y = d
			u.Z = j //(rand.Float64()*2 - 1) * d

			u.Nx = 0
			u.Ny = 1
			u.Nz = 0

			points = append(points, u)

			// down
			dwn := Point{}

			dwn.X = i //(rand.Float64()*2 - 1) * d
			dwn.Y = -1 * d
			dwn.Z = j //(rand.Float64()*2 - 1) * d

			dwn.Nx = 0
			dwn.Ny = -1
			dwn.Nz = 0

			points = append(points, dwn)

			// front
			f := Point{}

			f.X = i //(rand.Float64()*2 - 1) * d
			f.Y = j //(rand.Float64()*2 - 1) * d
			f.Z = d

			f.Nx = 0
			f.Ny = 0
			f.Nz = 1

			points = append(points, f)

			// back
			b := Point{}

			b.X = i //(rand.Float64()*2 - 1) * d
			b.Y = j //(rand.Float64()*2 - 1) * d
			b.Z = -1 * d

			b.Nx = 0
			b.Ny = 0
			b.Nz = -1

			points = append(points, b)
		}
	}

	return points
}
