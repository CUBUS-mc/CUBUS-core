package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

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
	xNew := cc.x + (float32(cc.nCubes%5)*isoDistance - float32(cc.nCubes/5)*isoDistance)
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
