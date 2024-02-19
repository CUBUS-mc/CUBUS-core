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
	x, y           float32
	size           float32
	selectCallback func(c *cube)
	id             string
}

func newCube(textureUrl string, selectCallback func(c *cube), id string, x float32, y float32) *cube {
	texture, err := fyne.LoadResourceFromURLString(textureUrl)
	if err != nil {
		panic(err)
	}
	cubeImage := &cube{
		x:              x,
		y:              y,
		size:           150,
		selectCallback: selectCallback,
		id:             id,
	}
	cubeImage.Resource = texture
	cubeImage.Move(fyne.NewPos(cubeImage.x, cubeImage.y))
	cubeImage.Resize(fyne.NewSize(cubeImage.size, cubeImage.size))
	return cubeImage
}

func (d *cube) isPointInHitbox(_, _ float32) bool {
	// TODO: Implement this
	return true
}

func (d *cube) Dragged(e *fyne.DragEvent) {
	if !d.isPointInHitbox(e.Position.X, e.Position.Y) {
		return
	}
	go func() {
		dx := d.x - (d.Position().X + e.Dragged.DX)
		dy := d.y - (d.Position().Y + e.Dragged.DY)
		distance := math.Sqrt(float64(dx*dx + dy*dy))
		scale := float32(1 / (1 + distance/50))
		d.Move(fyne.NewPos(d.Position().X+e.Dragged.DX*scale, d.Position().Y+e.Dragged.DY*scale))
		d.Refresh()
	}()
}

func (d *cube) DragEnd() {
	go func() {
		canvas.NewPositionAnimation(
			d.Position(),
			fyne.NewPos(d.x, d.y),
			time.Second/2,
			func(pos fyne.Position) {
				d.Move(pos)
				d.Refresh()
			}).Start()
	}()
}

func (d *cube) Tapped(e *fyne.PointEvent) {
	if !d.isPointInHitbox(e.Position.X, e.Position.Y) {
		return
	}
	d.selectCallback(d)
}

func (d *cube) MoveSmoothlyTo(x float32, y float32) {
	d.x = x
	d.y = y
	go func() {
		canvas.NewPositionAnimation(
			d.Position(),
			fyne.NewPos(d.x, d.y),
			time.Second/2,
			func(pos fyne.Position) {
				d.Move(pos)
				d.Refresh()
			}).Start()
	}()
}
