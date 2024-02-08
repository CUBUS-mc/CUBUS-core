package main

import (
	"math"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
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

func (d *cube) Tapped(_ *fyne.PointEvent) {
	d.selectCallback(d)
}

type cubeContainer struct {
	Container        *fyne.Container
	x                float32
	y                float32
	nCubes           int
	selected         *cube
	unselectCallback func()
}

func newCubeContainer(unselectCallback func(), x float32, y float32) *cubeContainer {
	return &cubeContainer{
		Container:        container.NewWithoutLayout(),
		unselectCallback: unselectCallback,
		x:                x,
		y:                y,
	}
}

func (cc *cubeContainer) changeSelected(c *cube, selectCallback func(c *cube)) {
	if cc.selected == c {
		cc.selected = nil
		cc.unselectCallback()
	} else {
		cc.selected = c
		selectCallback(cc.selected)
	}
}

func (cc *cubeContainer) AddCube(textureUrl string, selectCallback func(c *cube), id string) {
	isoDistance := float32(70)
	xNew := cc.x + float32(cc.nCubes%5)*isoDistance - float32(cc.nCubes/5)*isoDistance
	yNew := cc.y + (float32(cc.nCubes/5)*isoDistance/2 + float32(cc.nCubes%5)*isoDistance/2) + float32(cc.nCubes)
	c := newCube(textureUrl, func(c *cube) { cc.changeSelected(c, selectCallback) }, id, xNew, yNew)
	cc.Container.Add(c)
	cc.nCubes++
	cc.Container.Objects[cc.nCubes-1].Move(fyne.NewPos(xNew, yNew))
	cc.CenterCubes()
}

func (cc *cubeContainer) CenterCubes() {
	// TODO: Center cubes with a nice animation
}
