package julia

import (
	"math/cmplx"
)

type ComplexPoint struct {
	Z        complex128
	ZProc    complex128
	N        int
	ScaledN  float64
	diverged bool
}

type IteratedFunctionSystem interface {
	Init() error
	Generate() error
	SetPower() error
	GetPower() float64
	GetCVal() complex128
	GetGrid() *[]*ComplexPoint
}

type JuliaFractal struct {
	grid     []*ComplexPoint
	power    float64
	cval     complex128
	RowCount int
	maxiter  int
}

func (f *JuliaFractal) GetPower() float64 {
	return f.power
}

func (f *JuliaFractal) GetCVal() complex128 {
	return f.cval
}

func (f *JuliaFractal) GetGrid() *[]*ComplexPoint {
	return &f.grid
}

func (f *JuliaFractal) ScaleNValues() {
	maxN := 0
	for i := range f.grid {
		if f.grid[i].N > maxN {
			maxN = f.grid[i].N
		}
		if maxN == f.maxiter {
			break
		}
	}
	for i := range f.grid {
		f.grid[i].ScaledN = float64(f.grid[i].N) / float64(maxN)
	}
}

func (fractal *JuliaFractal) Init(gridWidth int, cval complex128, power float64, x1 float64, x2 float64, y1 float64, y2 float64) {
	dx := float64(x2-x1) / float64(gridWidth)
	dy := float64(y2-y1) / float64(gridWidth)

	ycurrent := float64(y1)
	xcurrent := float64(x1)
	cc := 0
	for i := range fractal.grid {
		fractal.grid[i] = &ComplexPoint{
			Z:        complex(xcurrent, ycurrent),
			ZProc:    complex(0, 0),
			N:        0,
			diverged: false,
		}

		if xcurrent >= x2 {
			xcurrent = x1
			fractal.RowCount = cc
			cc = 0
			ycurrent += dy
		} else {
			cc += 1
			xcurrent += dx
		}

		if ycurrent >= y2 {
			break
		}
	}
}

func NewJuliaFractal(gridWidth int, cval complex128, power float64, x1 float64, x2 float64, y1 float64, y2 float64) (fractal JuliaFractal) {
	nTotal := gridWidth * gridWidth
	fractal = JuliaFractal{
		grid:    make([]*ComplexPoint, nTotal),
		power:   power,
		cval:    cval,
		maxiter: 500,
	}

	fractal.Init(gridWidth, cval, power, x1, x2, y1, y2)

	return
}

func (f *JuliaFractal) StepRealC(val float64) {
	f.cval += complex(val, 0)
}

func (f *JuliaFractal) StepPower(val float64) {
	f.power += val
}

func (f *JuliaFractal) Generate() {
	doComputePoint := func(point *ComplexPoint, cval complex128, maxiter int) {
		point.ZProc = point.Z
		point.diverged = false
		point.N = 0
		for !point.diverged {
			point.ZProc = cmplx.Pow(point.ZProc, complex(f.power, 0)) + cval
			point.N += 1
			if point.N >= (maxiter) {
				break
			}
			if cmplx.Abs(point.ZProc) >= 2 {
				point.diverged = true
			}
		}
	}

	for i := range f.grid {
		doComputePoint(f.grid[i], f.cval, f.maxiter)
	}

}

func (f *JuliaFractal) PGenerate() {
	jobQueue := make(chan int, 100)

	doComputePoint := func(point *ComplexPoint, cval complex128, maxiter int) {
		for !point.diverged {
			point.Z = point.Z*point.Z + cval
			point.N += 1
			if point.N >= (maxiter) {
				break
			}
			if cmplx.Abs(point.Z) >= 2 {
				point.diverged = true
			}
		}
		<-jobQueue
	}

	for i := range f.grid {
		jobQueue <- 1
		go doComputePoint(f.grid[i], f.cval, f.maxiter)
	}

}
