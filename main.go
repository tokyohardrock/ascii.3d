package main

import (
	"log"
	"math"
	"time"
)

const (
	d     = 3  // distance to camera
	gridW = 40 // grid width
	gridH = 40 // grid height
	fps   = 60
	delta = math.Pi * 0.01
)

type vertices struct {
	x []float64
	y []float64
	z []float64
}

type scrnCoords struct {
	x []int
	y []int
}

var vrtcs vertices = vertices{
	[]float64{-2, 1, 2, -2, -2.5, 2, 2, -2},
	[]float64{-2, -1, 2, 2, -2.5, -2, 2, 2},
	[]float64{-2, -1, -2, -2, 2.5, 2, 2, 2},
}

var edges []int = []int{
	1, 2,
	2, 3,
	3, 4,
	4, 1,

	5, 6,
	6, 7,
	7, 8,
	8, 5,

	1, 5,
	2, 6,
	3, 7,
	4, 8,
}

var xyzRotations = []float64{0, 0, 0}

var ascii []byte = []byte{' ', '#'}

func normalizeVertices(vrt *vertices) {
	var maxElem float64
	for i := range vrt.x {
		val := math.Max(math.Abs(vrt.x[i]), math.Max(math.Abs(vrt.y[i]), math.Abs(vrt.z[i])))
		maxElem = math.Max(maxElem, val)
	}
	if maxElem == 0 {
		log.Fatal("Singularity detected — normalization aborted")
		return
	}
	scale := 1 / maxElem
	for i := range vrt.x {
		vrt.x[i] *= scale
		vrt.y[i] *= scale
		vrt.z[i] *= scale
	}
}

func getTransformedCoords(x, y, z float64) (float64, float64, float64) {
	var coordsMatrix = [][]float64{
		{x},
		{y},
		{z},
	}
	var xRot = xRotationMatrix(xyzRotations[0])
	var yRot = yRotationMatrix(xyzRotations[1])
	var zRot = zRotationMatrix(xyzRotations[2])
	var resultTransformMatrix = mulitplyMatrices(mulitplyMatrices(xRot, yRot), zRot)
	var result = mulitplyMatrices(resultTransformMatrix, coordsMatrix)
	return result[0][0], result[1][0], result[2][0]
}

func getScreenCoord(x, y, z float64) (float64, float64) {
	var x_proj = x / (z + d)
	var y_proj = y / (z + d)
	var scrnX = (x_proj + 1) * (gridW - 1) * 0.5
	var scrnY = (1 - y_proj) * (gridH - 1) * 0.5
	return scrnX, scrnY
}

func Bresenham(x1, y1, x2, y2 int) scrnCoords {
	var line scrnCoords = scrnCoords{
		make([]int, 0),
		make([]int, 0),
	}
	var dx = int(math.Abs(float64(x2 - x1)))
	var dy = -int(math.Abs(float64(y2 - y1)))
	var err = dx + dy
	var sx, sy = 1, 1
	if x1 >= x2 {
		sx = -1
	}
	if y1 >= y2 {
		sy = -1
	}
	for {
		line.x = append(line.x, x1)
		line.y = append(line.y, y1)
		if x1 == x2 && y1 == y2 {
			break
		}
		e2 := 2 * err
		if e2 >= dy {
			err += dy
			x1 += sx
		}
		if e2 <= dx {
			err += dx
			y1 += sy
		}
	}
	return line
}

func draw() {
	var grid = make([][]int, gridH)
	var projCoords scrnCoords = scrnCoords{
		make([]int, len(vrtcs.x)),
		make([]int, len(vrtcs.x)),
	}
	for i := range grid {
		grid[i] = make([]int, gridW)
	}
	for i := range vrtcs.x {
		scrnX, scrnY := getScreenCoord(getTransformedCoords(vrtcs.x[i], vrtcs.y[i], vrtcs.z[i]))
		projCoords.x[i] = int(scrnX)
		projCoords.y[i] = int(scrnY)
		grid[projCoords.y[i]][projCoords.x[i]] = 1
	}
	for i := range edges[:len(edges)-1] {
		if (i+1)%2 == 0 {
			continue
		}
		var line = Bresenham(
			projCoords.x[edges[i]-1],
			projCoords.y[edges[i]-1],
			projCoords.x[edges[i+1]-1],
			projCoords.y[edges[i+1]-1],
		)
		for j := range line.x {
			grid[line.y[j]][line.x[j]] = 1
		}
	}
	for i := range grid {
		for j := range grid[i] {
			print(" " + string(ascii[grid[i][j]]) + " ")
		}
		print("\n")
	}
}

func main() {
	normalizeVertices(&vrtcs)
	for {
		time.Sleep(time.Second / fps)
		ClearScreen()
		for i := range xyzRotations {
			xyzRotations[i] = math.Mod(xyzRotations[i], math.Pi*2) + delta
		}
		draw()
	}
}
