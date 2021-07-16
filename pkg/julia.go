package julia

import "image"

type ComplexPoint struct {
	z        complex128
	n        int32
	diverged bool
}

type IteratedFunctionSystem interface {
	Init() error
	Step() error
	SetPower() error
}

type JuliaFractal struct {
	grid  []*ComplexPoint
	power float64
}

func NewJuliaFractal(gridWidth int, cval complex128, power float64, coordinateBounds image.Rectangle) (fractal JuliaFractal) {
	nTotal := gridWidth * gridWidth
	fractal = JuliaFractal{
		grid:  make([]*ComplexPoint, nTotal),
		power: power,
	}

	x1, x2, y1, y2 := coordinateBounds.Min.X, coordinateBounds.Max.X, coordinateBounds.Min.Y, coordinateBounds.Max.X
	dx := float64(x2-x1) / float64(gridWidth)
	dy := float64(y2-y1) / float64(gridWidth)

	ycurrent := float64(y1)
	xcurrent := float64(x1)

	for i := range fractal.grid {
		fractal.grid[i] = &ComplexPoint{
			z:        complex(xcurrent, ycurrent),
			n:        0,
			diverged: false,
		}

		xcurrent += dx
		ycurrent += dy
	}

}
