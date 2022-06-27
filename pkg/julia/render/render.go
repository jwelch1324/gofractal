package render

import "gocv.io/x/gocv"

type RGBCanvas struct {
	raw_img_R gocv.Mat
	raw_img_G gocv.Mat
	raw_img_B gocv.Mat
}

func NewRGBCanvas(width int, height int) (r *RGBCanvas) {
	r = &RGBCanvas{}
	r.raw_img_R = gocv.NewMatWithSize(width, height, gocv.MatTypeCV8UC1)
	r.raw_img_G = gocv.NewMatWithSize(width, height, gocv.MatTypeCV8UC1)
	r.raw_img_B = gocv.NewMatWithSize(width, height, gocv.MatTypeCV8UC1)
	return
}

func (r *RGBCanvas) Close() {
	r.raw_img_B.Close()
	r.raw_img_G.Close()
	r.raw_img_R.Close()
}

func (r *RGBCanvas) SetPixelValues(x int, y int, R uint8, G uint8, B uint8) {
	r.raw_img_R.SetUCharAt(x, y, R)
	r.raw_img_G.SetUCharAt(x, y, G)
	r.raw_img_B.SetUCharAt(x, y, B)
}

func (r *RGBCanvas) Draw() (img gocv.Mat) {
	img = gocv.NewMatWithSize(r.raw_img_R.Rows(), r.raw_img_R.Cols(), gocv.MatTypeCV8UC3)
	gocv.Merge([]gocv.Mat{r.raw_img_B, r.raw_img_G, r.raw_img_R}, &img)
	return
}
