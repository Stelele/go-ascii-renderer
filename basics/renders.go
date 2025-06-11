package main

import (
	"math"
)

var (
	A      float64   = 1
	B      float64   = 1
	output []float64 = make([]float64, width*height)
)

type Point struct {
	X  float64
	Y  float64
	Z  float64
	Nx float64
	Ny float64
	Nz float64
}

func test(x int, y int) float64 {
	u0 := 2*(float64(x)/float64(width)) - 1
	v0 := 2*(float64(y)/float64(height)) - 1

	u := 2*(u0-math.Floor(u0)) - 1
	v := 2*(v0-math.Floor(v0)) - 1

	d := math.Sqrt(u*u + v*v)

	d = math.Sin(d*8 + t*0.8)
	d = math.Abs(d)

	return d
}

func donut(x int, y int) float64 {
	return output[y*width+x]
}

func preProcessFrame() {
	R1 := 1
	R2 := 2
	K2 := 20
	K1 := width * K2 * 2 / (16 * (R1 + R2))

	cosA := math.Cos(A)
	sinA := math.Sin(A)
	cosB := math.Cos(B)
	sinB := math.Sin(B)

	output = make([]float64, width*height)
	zbuffer := make([]float64, width*height)

	for theta := 0.; theta < 2*math.Pi; theta += 0.3 {
		costheta := math.Cos(theta)
		sintheta := math.Sin(theta)

		for phi := 0.; phi < 2*math.Pi; phi += 0.1 {
			cosphi := math.Cos(phi)
			sinphi := math.Sin(phi)

			circlex := float64(R2) + float64(R1)*costheta
			circley := float64(R1) * sintheta

			x := circlex*(cosB*cosphi+sinA*sinB*sinphi) - circley*cosA*sinB
			y := circlex*(sinB*cosphi-sinA*cosB*sinphi) + circley*cosA*cosB
			z := float64(K2) + cosA*circlex*sinphi + circley*sinA
			ooz := 1 / z

			xp := width/2 + int(float64(K1)*ooz*x)
			yp := height/2 - int(float64(K1)*ooz*y)

			L := 1 / math.Sqrt(2) * (cosphi*costheta*sinB - cosA*costheta*sinphi - sinA*sintheta + cosB*(cosA*sintheta-costheta*sinA*sinphi))

			if L >= 0 && xp >= 0 && yp >= 0 {
				index := yp*width + xp
				if index >= 0 && index < width*height && ooz > zbuffer[index] {
					zbuffer[index] = ooz
					output[index] = L
				}
			}
		}
	}
}

func preProcessAnyFrame() {
	K2 := 20.
	K1 := float64(width) * K2 * 2 / (16 * (1 + 2))

	Rx := 1.
	Ry := 1.
	Rz := -5.
	ood := 1 / math.Sqrt(Rx*Rx+Ry*Ry+Rz*Rz)

	output = make([]float64, width*height)
	zbuffer := make([]float64, width*height)

	points := getShape()
	for _, point := range points {
		rP := rotate(point)

		xP := width/2 + int((K1*rP.X)/(K2+rP.Z))
		yP := height/2 - int((K1*rP.Y)/(K2+rP.Z))
		ooz := 1 / (K2 + rP.Z)

		L := ood * (Rx*rP.Nx + Ry*rP.Ny + Rz*rP.Nz)

		if L > 0 {
			index := yP*width + xP
			if index >= 0 && index < width*height && ooz > zbuffer[index] {
				output[index] = L
				zbuffer[index] = ooz
			}
		}
	}
}

func rotate(p Point) Point {
	cA := math.Cos(A)
	sA := math.Sin(A)
	cB := math.Cos(B)
	sB := math.Sin(B)

	a11 := cB
	a12 := sB
	a13 := 0.

	a21 := cA * -1 * sB
	a22 := cA * cB
	a23 := sA

	a31 := sA * sB
	a32 := -1 * sA * cB
	a33 := cA

	rP := Point{}

	rP.X = p.X*a11 + p.Y*a21 + p.Z*a31
	rP.Y = p.X*a12 + p.Y*a22 + p.Z*a32
	rP.Z = p.X*a13 + p.Y*a23 + p.Z*a33

	rP.Nx = p.Nx*a11 + p.Ny*a21 + p.Nz*a31
	rP.Ny = p.Nx*a12 + p.Ny*a22 + p.Nz*a32
	rP.Nz = p.Nx*a13 + p.Ny*a23 + p.Nz*a33

	return rP
}
