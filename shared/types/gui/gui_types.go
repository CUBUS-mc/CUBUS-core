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
	window           fyne.Window
}

type Cube struct {
	widget.Icon
	X, Y           float32
	CubeSize       float32
	SelectCallback func(c *Cube)
	Id             string
	Config         types.CubeConfig
}

func NewCubeContainer(unselectCallback func(), x float32, y float32, window fyne.Window) *CubeContainer {
	return &CubeContainer{
		Container:        container.NewWithoutLayout(),
		UnselectCallback: unselectCallback,
		X:                x,
		Y:                y,
		IsoDistance:      float32(70),
		window:           window,
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

func (cc *CubeContainer) AddCube(textureUrl string, selectCallback func(c *Cube), id string, cubeConfig types.CubeConfig) {
	c := newCube(textureUrl, func(c *Cube) { cc.ChangeSelected(c, selectCallback) }, id, 0, 0, cubeConfig)
	cc.Container.Add(c)
	cc.NCubes++
	cc.CenterCubes()
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
			cube.Move(fyne.NewPos(xNew, yNew))
			cube.X = xNew
			cube.Y = yNew
		}
	}
}

func (cc *CubeContainer) rePlaceCubes() {
	cubesPerRow := int(math.Ceil(math.Sqrt(float64(cc.NCubes))))
	for i := 0; i < cc.NCubes; i++ {
		row := i / cubesPerRow
		col := i % cubesPerRow

		threeDimX := float32(row*95) + cc.IsoDistance
		threeDimY := float32(0)
		threeDimZ := float32(col*95) + cc.IsoDistance

		xNew := cc.X + ((threeDimX - threeDimZ) / 1.4142135624)
		yNew := cc.Y + ((threeDimX + 2*threeDimY + threeDimZ) / 2.4494897428)
		if cube, ok := cc.Container.Objects[i].(*Cube); ok {
			cube.Move(fyne.NewPos(xNew, yNew))
			cube.X = xNew
			cube.Y = yNew
		}
	}
	cc.Container.Refresh()
}

func (cc *CubeContainer) CenterCubes() {
	cc.rePlaceCubes()

	var sumX, sumY float32
	for i := 0; i < cc.NCubes; i++ {
		if cube, ok := cc.Container.Objects[i].(*Cube); ok {
			sumX += cube.Position().X + cube.CubeSize/2
			sumY += cube.Position().Y + cube.CubeSize/2
		}
	}
	meanX := sumX / float32(cc.NCubes)
	meanY := sumY / float32(cc.NCubes)

	deltaX := cc.window.Canvas().Size().Width*0.5 - meanX
	deltaY := cc.window.Canvas().Size().Height*0.5 - meanY

	cc.MoveContainer(cc.X+deltaX, cc.Y+deltaY)
}

func (cc *CubeContainer) ClearCubes() {
	cc.Container.Objects = []fyne.CanvasObject{}
	cc.NCubes = 0
	cc.Selected = nil
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
