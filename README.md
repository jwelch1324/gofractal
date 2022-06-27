A 2D Fractal Renderer written in Go

A small project to procedurally generate Julia sets of the mandelbrot fractal 

# Running 
```bash
go build ./cmd/julia
./julia
```

Keyboard Controls
- `[` and `]` keys will change the power
- `o` will toggle display of the currently rendered fractal
- `i` will invert the color space
- `c` will change the color pallete used to render
- 'f' will freeze the parameter change from being rendered (useful if you want to make a large parameter change without rendering each step)
- `esc` or `q` will close the program

Examples of outputs

![image](https://user-images.githubusercontent.com/4706333/176051581-b453bb4f-bd83-4018-a7a2-4303032e7724.png)


![image](https://user-images.githubusercontent.com/4706333/176051679-9e729437-af3f-4cbb-8e6d-5c9dbe330a50.png)


![image](https://user-images.githubusercontent.com/4706333/176051740-2b853f0b-f548-443a-a839-bc160c54de03.png)
