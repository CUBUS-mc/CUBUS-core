package main

import (
	"math"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

type cube struct {
	widget.Icon
	x, y float32
	size float32
}

func newCube(textureUrl string) *cube {
	texture, err := fyne.LoadResourceFromURLString(textureUrl)
	if err != nil {
		panic(err)
	}
	cubeImage := &cube{
		x:    400,
		y:    300,
		size: 150,
	}
	cubeImage.Resource = texture
	cubeImage.Move(fyne.NewPos(cubeImage.x, cubeImage.y))
	cubeImage.Resize(fyne.NewSize(cubeImage.size, cubeImage.size))
	return cubeImage
}

func (d *cube) Dragged(e *fyne.DragEvent) {
	dx := d.x - (d.Position().X + e.Dragged.DX)
	dy := d.y - (d.Position().Y + e.Dragged.DY)
	distance := math.Sqrt(float64(dx*dx + dy*dy))
	scale := float32(1 / (1 + distance/50))
	d.Move(fyne.NewPos(d.Position().X+e.Dragged.DX*scale, d.Position().Y+e.Dragged.DY*scale))
	d.Refresh()
}

func (d *cube) DragEnd() {
	canvas.NewPositionAnimation(
		d.Position(),
		fyne.NewPos(d.x, d.y),
		time.Second/2,
		func(pos fyne.Position) {
			d.Move(pos)
			d.Refresh()
		}).Start()
}
