package main

import (
	"fmt"
	"image"
	"image/color"
	"sync"
	"time"

	"github.com/jwelch1324/gofractal/pkg/julia"
	"github.com/jwelch1324/gofractal/pkg/julia/render"
	"github.com/mazznoer/colorgrad"
	"github.com/panjf2000/ants/v2"
	"gocv.io/x/gocv"
)

func main() {
	defer ants.Release()
	var window *gocv.Window = nil
	window = gocv.NewWindow("JuliaTest")
	defer window.Close()

	gridWidth := 500
	renderer := render.NewRGBCanvas(gridWidth, gridWidth)

	fractal := julia.NewJuliaFractal(gridWidth, complex(-0.6772, -0.72), 2, -2.5, 2.5, -2.5, 2.5)

	grads := []colorgrad.Gradient{colorgrad.Turbo(), colorgrad.Blues(), colorgrad.Rainbow(), colorgrad.Magma(), colorgrad.Sinebow()}
	gc := 0
	grad := grads[gc]
	paramChange := true
	invertFlag := false
	showParams := true
	freezeUpdates := false

	for {

		if paramChange && !freezeUpdates {
			fractal.Generate()
			fractal.ScaleNValues() // Puts all N values on the 0...1 scale
			renderer.Close()
			renderer = render.NewRGBCanvas(gridWidth, gridWidth)
			paramChange = false
		}

		grid := fractal.GetGrid()

		xc := 0
		yc := 0

		s := time.Now()

		var wg sync.WaitGroup

		// loadPixel := func(xc int, yc int, r *render.Render, nval float64, grad colorgrad.Gradient) {
		// 	color := grad.At(nval)
		// 	r.SetPixelValues(yc, xc, uint8(color.R*255.0), uint8(color.G*255.0), uint8(color.B*255.0))
		// 	wg.Done()
		// }

		for i := range *grid {
			nval := (*grid)[i].ScaledN
			if invertFlag {
				nval = 1 - nval
			}
			color := grad.At(nval)
			renderer.SetPixelValues(yc, xc, uint8(color.R*255.0), uint8(color.G*255.0), uint8(color.B*255.0))
			// jtask := func() {
			// 	renderer.SetPixelValues(yc, xc, uint8(color.R*255.0), uint8(color.G*255.0), uint8(color.B*255.0))
			// 	wg.Done()
			// }
			// wg.Add(1)
			//_ = ants.Submit(jtask)

			if xc == fractal.RowCount {
				xc = 0
				yc += 1
			} else {
				xc += 1
			}
			if yc == gridWidth {
				break
			}
		}
		wg.Wait()
		fmt.Println(time.Since(s))

		newImg := renderer.Draw()
		if gridWidth < 1000 {
			gocv.Resize(newImg, &newImg, image.Point{}, 2, 2, gocv.InterpolationCubic)
		}

		if showParams {
			cv := fractal.GetCVal()
			gocv.PutText(&newImg, fmt.Sprintf("Power: %.5f", fractal.GetPower()), image.Pt(int(float64(gridWidth)*0.02), int(float64(gridWidth)*0.06)), gocv.FontHersheyPlain, 2, color.RGBA{R: 255, G: 255, B: 255, A: 0}, 2)
			gocv.PutText(&newImg, fmt.Sprintf("CVal: %.5f + i%.5f", real(cv), imag(cv)), image.Pt(int(float64(gridWidth)*0.02), int(float64(gridWidth)*0.11)), gocv.FontHersheyPlain, 2, color.RGBA{R: 255, G: 255, B: 255, A: 0}, 2)
			//gocv.PutText(&newImg, fmt.Sprintf("Power: %.5f", fractal.GetPower()), image.Pt(int(float64(gridWidth)*0.02), int(float64(gridWidth)*0.02)), gocv.FontHersheyPlain, 2, color.RGBA{R: 255, G: 255, B: 255, A: 0}, 2)
		}

		window.IMShow(newImg)
		newImg.Close()

		if k := window.WaitKey(1); k >= 0 {
			switch k {
			case 113: // q
				return
			case 27: //esc
				return
			case 99: //c - change color scheme
				gc += 1
				if gc == len(grads) {
					gc = 0
				}
				grad = grads[gc]
			case 111:
				showParams = !showParams
			case 105:
				invertFlag = !invertFlag
			case 108:
				fractal.StepRealC(0.01)
				paramChange = true
			case 106:
				fractal.StepRealC(-0.01)
				paramChange = true
			case 93: //]
				fractal.StepPower(0.01)
				paramChange = true
			case 91: //[
				fractal.StepPower(-0.01)
				paramChange = true
			case 102: // f
				freezeUpdates = !freezeUpdates
			default:
				fmt.Printf("Key Pressed: %d\n", k)
			}

		}
	}
}
