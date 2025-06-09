package main

import (
	"math"
)

var (
	A      float64   = 1
	B      float64   = 1
	output []float64 = make([]float64, width*height)
)

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

func preProcessFrame(A float64, B float64) {
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

			if L > 0 {
				index := yp*width + xp
				if index >= 0 && index < width*height && ooz > zbuffer[index] {
					zbuffer[index] = ooz
					output[index] = L
				}
			}
		}
	}
}
