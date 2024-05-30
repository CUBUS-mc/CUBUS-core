package gui

import (
	"CUBUS-core/shared/types"
	"context"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"golang.org/x/sync/semaphore"
	"math"
	"sync"
	"time"
)

type CubeContainer struct {
	Container        *fyne.Container
	X                float32
	Y                float32
	NCubes           int
	Selected         *Cube
	UnselectCallback func()
	IsoDistance      float32
	Mu               sync.Mutex
}

type Cube struct {
	widget.Icon
	X, Y           float32
	CubeSize       float32
	SelectCallback func(c *Cube)
	Id             string
	Config         types.CubeConfig
}

func NewCubeContainer(unselectCallback func(), x float32, y float32) *CubeContainer {
	return &CubeContainer{
		Container:        container.NewWithoutLayout(),
		UnselectCallback: unselectCallback,
		X:                x,
		Y:                y,
		IsoDistance:      float32(70),
	}
}

func (cc *CubeContainer) ChangeSelected(c *Cube, selectCallback func(c *Cube)) {
	cc.Mu.Lock()
	defer cc.Mu.Unlock()
	if cc.Selected == c {
		cc.Selected = nil
		cc.UnselectCallback()
	} else {
		cc.Selected = c
		selectCallback(cc.Selected)
	}
}

func (cc *CubeContainer) AddCube(textureUrl string, selectCallback func(c *Cube), id string, cubeConfig types.CubeConfig) { // TODO: Change method so it adds the cubes in a square not a rectangle
	if math.IsNaN(float64(cc.X)) { // TODO: Fix that there could be nan (this is a problem in the setupDialog method)
		cc.X = 0
	}
	if math.IsNaN(float64(cc.Y)) {
		cc.Y = 0
	}
	xNew := cc.X + (float32(cc.NCubes%5)*cc.IsoDistance - float32(cc.NCubes/5)*cc.IsoDistance)
	yNew := cc.Y + (float32(cc.NCubes/5)*cc.IsoDistance/2 + float32(cc.NCubes%5)*cc.IsoDistance/2) + float32(cc.NCubes)
	c := newCube(textureUrl, func(c *Cube) { cc.ChangeSelected(c, selectCallback) }, id, xNew, yNew, cubeConfig)
	cc.Container.Add(c)
	cc.NCubes++
	// cc.Container.Objects[cc.NCubes-1].Move(fyne.NewPos(xNew, yNew))
}

func (cc *CubeContainer) MoveContainer(x float32, y float32) {
	deltaX := x - cc.X
	deltaY := y - cc.Y
	cc.X = x
	cc.Y = y
	for i := 0; i < cc.NCubes; i++ {
		xNew := cc.Container.Objects[i].Position().X + deltaX
		yNew := cc.Container.Objects[i].Position().Y + deltaY
		if cube, ok := cc.Container.Objects[i].(*Cube); ok {
			cube.MoveSmoothlyTo(xNew, yNew)
		}
	}
}

func (cc *CubeContainer) CenterCubes() {
	var sumX, sumY float32
	for i := 0; i < cc.NCubes; i++ {
		if cube, ok := cc.Container.Objects[i].(*Cube); ok {
			sumX += cube.Position().X + cube.CubeSize/2
			sumY += cube.Position().Y + cube.CubeSize/2
		}
	}
	meanX := sumX / float32(cc.NCubes)
	meanY := sumY / float32(cc.NCubes)

	deltaX := 700 - meanX
	deltaY := 450 - meanY

	cc.MoveContainer(cc.X+deltaX, cc.Y+deltaY)
}

var sem = semaphore.NewWeighted(1)

func newCube(textureUrl string, selectCallback func(c *Cube), id string, x float32, y float32, cubeConfig types.CubeConfig) *Cube {
	texture, err := fyne.LoadResourceFromURLString(textureUrl) // TODO: Change texture based on cube type
	if err != nil {
		panic(err)
	}
	cubeImage := &Cube{
		X:              x,
		Y:              y,
		CubeSize:       150,
		SelectCallback: selectCallback,
		Id:             id,
		Config:         cubeConfig,
	}
	cubeImage.Resource = texture
	cubeImage.Move(fyne.NewPos(cubeImage.X, cubeImage.Y))
	cubeImage.Resize(fyne.NewSize(cubeImage.CubeSize, cubeImage.CubeSize))
	return cubeImage
}

func (d *Cube) isPointInHitbox(_, _ float32) bool {
	// TODO: Create a better hitbox for the cubes
	return true
}

func (d *Cube) Dragged(e *fyne.DragEvent) {
	if !d.isPointInHitbox(e.Position.X, e.Position.Y) {
		return
	}
	if err := sem.Acquire(context.Background(), 1); err != nil {
		return
	}
	go func() {
		defer sem.Release(1)
		dx := d.X - (d.Position().X + e.Dragged.DX)
		dy := d.Y - (d.Position().Y + e.Dragged.DY)
		distance := math.Sqrt(float64(dx*dx + dy*dy))
		scale := float32(1 / (1 + distance/50))
		d.Move(fyne.NewPos(d.Position().X+e.Dragged.DX*scale, d.Position().Y+e.Dragged.DY*scale))
		d.Refresh()
	}()
}

func (d *Cube) DragEnd() {
	d.AnimateTo(d.X, d.Y)
}

func (d *Cube) Tapped(e *fyne.PointEvent) {
	if !d.isPointInHitbox(e.Position.X, e.Position.Y) {
		return
	}
	d.SelectCallback(d)
}

func (d *Cube) MoveSmoothlyTo(x float32, y float32) {
	d.X = x
	d.Y = y
	d.AnimateTo(x, y)
}

func (d *Cube) AnimateTo(x float32, y float32) {
	if err := sem.Acquire(context.Background(), 1); err != nil {
		return
	}
	go func() {
		defer sem.Release(1)
		canvas.NewPositionAnimation(
			d.Position(),
			fyne.NewPos(x, y),
			time.Second/2,
			func(pos fyne.Position) {
				d.Move(pos)
				d.Refresh()
			}).Start()
	}()
}
