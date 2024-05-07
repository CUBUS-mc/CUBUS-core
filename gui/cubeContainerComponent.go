package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"sync"
)

type cubeContainer struct {
	Container        *fyne.Container
	x                float32
	y                float32
	nCubes           int
	selected         *cube
	unselectCallback func()
	isoDistance      float32
	mu               sync.Mutex
}

func newCubeContainer(unselectCallback func(), x float32, y float32) *cubeContainer {
	return &cubeContainer{
		Container:        container.NewWithoutLayout(),
		unselectCallback: unselectCallback,
		x:                x,
		y:                y,
		isoDistance:      float32(70),
	}
}

func (cc *cubeContainer) changeSelected(c *cube, selectCallback func(c *cube)) {
	cc.mu.Lock()
	defer cc.mu.Unlock()
	if cc.selected == c {
		cc.selected = nil
		cc.unselectCallback()
	} else {
		cc.selected = c
		selectCallback(cc.selected)
	}
}

func (cc *cubeContainer) AddCube(textureUrl string, selectCallback func(c *cube), id string) { // TODO: Change method so it adds the cubes in a square not a rectangle
	xNew := cc.x + (float32(cc.nCubes%5)*cc.isoDistance - float32(cc.nCubes/5)*cc.isoDistance)
	yNew := cc.y + (float32(cc.nCubes/5)*cc.isoDistance/2 + float32(cc.nCubes%5)*cc.isoDistance/2) + float32(cc.nCubes)
	c := newCube(textureUrl, func(c *cube) { cc.changeSelected(c, selectCallback) }, id, xNew, yNew)
	cc.Container.Add(c)
	cc.nCubes++
	cc.Container.Objects[cc.nCubes-1].Move(fyne.NewPos(xNew, yNew))
}

func (cc *cubeContainer) MoveContainer(x float32, y float32) {
	deltaX := x - cc.x
	deltaY := y - cc.y
	cc.x = x
	cc.y = y
	for i := 0; i < cc.nCubes; i++ {
		xNew := cc.Container.Objects[i].Position().X + deltaX
		yNew := cc.Container.Objects[i].Position().Y + deltaY
		if cube, ok := cc.Container.Objects[i].(*cube); ok {
			cube.MoveSmoothlyTo(xNew, yNew)
		}
	}
}

func (cc *cubeContainer) CenterCubes() {
	var sumX, sumY float32
	for i := 0; i < cc.nCubes; i++ {
		if cube, ok := cc.Container.Objects[i].(*cube); ok {
			sumX += cube.Position().X + cube.size/2
			sumY += cube.Position().Y + cube.size/2
		}
	}
	meanX := sumX / float32(cc.nCubes)
	meanY := sumY / float32(cc.nCubes)

	deltaX := 700 - meanX
	deltaY := 450 - meanY

	cc.MoveContainer(cc.x+deltaX, cc.y+deltaY)
}
