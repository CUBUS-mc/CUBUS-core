package gui

import (
	"CUBUS-core/shared"
	"CUBUS-core/shared/translation"
	"CUBUS-core/shared/types"
	"CUBUS-core/shared/types/gui"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"time"
)

func selectCube(c *gui.Cube, infoContainerShape *canvas.Rectangle, pointerLine *canvas.Line, pointerTip *canvas.Circle, infoContainerText *widget.RichText) {
	println(c.X)
	println(c.Y)
	go func() {
		infoContainerText.Segments = []widget.RichTextSegment{ // TODO: Show all info about the cube and add a button to delete or edit the cube
			&widget.TextSegment{
				Text: "Cube info\n",
				Style: widget.RichTextStyle{
					TextStyle: fyne.TextStyle{Bold: true},
					ColorName: theme.ColorNameBackground,
				},
			},
			&widget.TextSegment{
				Text: "Cube ID: " + c.Id + "\n",
				Style: widget.RichTextStyle{
					ColorName: theme.ColorNameBackground,
				},
			},
			&widget.TextSegment{
				Text: "Cube Name: " + c.Config.CubeName + "\n",
				Style: widget.RichTextStyle{
					ColorName: theme.ColorNameBackground,
				},
			},
		}
		infoContainerText.Refresh()
		canvas.NewPositionAnimation(
			infoContainerShape.Position(),
			fyne.NewPos(1400-325, 25),
			time.Second/4,
			func(pos fyne.Position) {
				infoContainerShape.Move(pos)
				infoContainerShape.Refresh()
				infoContainerText.Move(fyne.NewPos(pos.X+5, pos.Y+10))
				infoContainerText.Refresh()
			}).Start()
		pointerTip.Move(fyne.NewPos(c.Position().X+c.CubeSize/2-5, c.Position().Y+c.CubeSize/2-5))
		pointerTip.Show()
		pointerLine.Resize(fyne.NewSize(0, 0))
		pointerLine.Move(fyne.NewPos(c.Position().X+c.CubeSize/2, c.Position().Y+c.CubeSize/2))
		pointerLine.Show()
		canvas.NewSizeAnimation(
			pointerLine.Size(),
			fyne.NewSize(1400-c.Position().X-325, 2),
			time.Second/4,
			func(size fyne.Size) {
				pointerLine.Resize(size)
				pointerLine.Refresh()
			}).Start()
	}()
}

func unselectCube(infoContainerShape *canvas.Rectangle, pointerLine *canvas.Line, pointerTip *canvas.Circle, infoContainerText *widget.RichText) {
	go func() {
		canvas.NewPositionAnimation(
			infoContainerShape.Position(),
			fyne.NewPos(1400, 25),
			time.Second/4,
			func(pos fyne.Position) {
				infoContainerShape.Move(pos)
				infoContainerShape.Refresh()
				infoContainerText.Move(fyne.NewPos(pos.X+5, pos.Y+10))
				infoContainerText.Refresh()
			}).Start()
		pointerTip.Hide()
		pointerLine.Hide()
	}()
}

func Gui(cubusApp fyne.App, defaults *shared.Defaults) {
	T := translation.T
	cubeStrings := cubusApp.Preferences().StringListWithFallback("cubes", []string{})
	cubusApp.Preferences().SetStringList("cubes", cubeStrings)
	cubeConfigs := make([]map[string]interface{}, len(cubeStrings))
	for i, cubeString := range cubeStrings {
		cubeConfigs[i] = shared.JsonStringToObject(cubeString)
	}

	cubusWindow := cubusApp.NewWindow("CUBUS core")
	cubusWindow.Resize(fyne.NewSize(1400, 900))
	cubusWindow.CenterOnScreen()
	cubusWindow.SetIcon(cubusApp.Icon())
	cubusWindow.SetFixedSize(true)

	infoContainerShape := canvas.NewRectangle(color.White)
	infoContainerShape.Resize(fyne.NewSize(300, 825))
	infoContainerShape.Move(fyne.NewPos(1400, 25))
	infoContainerShape.CornerRadius = 12

	infoContainerText := widget.NewRichText()
	infoContainerText.Resize(fyne.NewSize(290, 800))
	infoContainerText.Wrapping = fyne.TextWrapBreak

	pointerLine := canvas.NewLine(color.White)
	pointerLine.StrokeWidth = 2
	pointerLine.Hide()
	pointerTip := canvas.NewCircle(color.White)
	pointerTip.Resize(fyne.NewSize(10, 10))
	pointerTip.Hide()

	cubeContainerObject := gui.NewCubeContainer(func() { unselectCube(infoContainerShape, pointerLine, pointerTip, infoContainerText) }, 500, 300)
	for _, cubeConfig := range cubeConfigs {
		cubeConfigAsCorrectType := types.CubeConfig{
			Id:        cubeConfig["id"].(string),
			CubeType:  types.CubeType{Value: cubeConfig["type"].(string)},
			CubeName:  cubeConfig["name"].(string),
			PublicKey: nil,
		}
		cubeContainerObject.AddCube(defaults.CubeAssetURL, func(c *gui.Cube) { selectCube(c, infoContainerShape, pointerLine, pointerTip, infoContainerText) }, cubeConfig["id"].(string), cubeConfigAsCorrectType)
	}
	cubeContainerObject.CenterCubes()

	windowMenu := fyne.NewMainMenu(
		fyne.NewMenu(T("File"),
			fyne.NewMenuItem("Create a new Cube", func() {
				setupDialog(
					cubusWindow,
					&cubeConfigs,
					&cubeStrings,
					cubusApp,
					defaults,
					cubeContainerObject,
					infoContainerShape,
					pointerLine,
					pointerTip,
					infoContainerText,
				)
			}),
			fyne.NewMenuItem(T("Export config"), func() {
				println("Export config") // TODO: implement config export
			}),
			fyne.NewMenuItem(T("Import config"), func() {
				println("Import config") // TODO: implement config import
			}),
			fyne.NewMenuItem(T("Settings"), func() {
				println("Settings") // TODO: implement settings (e.g. change the language)
			}),
		),
	)
	cubusWindow.SetMainMenu(windowMenu)

	mainContainer := container.NewWithoutLayout(cubeContainerObject.Container, pointerLine, pointerTip, infoContainerShape, infoContainerText)
	cubusWindow.SetContent(mainContainer)

	go func() {
		for {
			infoContainerShape.Refresh()
			pointerLine.Refresh()
			pointerTip.Refresh()
			cubeContainerObject.Mu.Lock()
			selected := cubeContainerObject.Selected
			cubeContainerObject.Mu.Unlock()
			if selected != nil {
				pointerTip.Move(fyne.NewPos(selected.Position().X+selected.CubeSize/2-5, selected.Position().Y+selected.CubeSize/2-40))
				pointerLine.Move(fyne.NewPos(selected.Position().X+selected.CubeSize/2, selected.Position().Y+selected.CubeSize/2-35))
				time.AfterFunc(time.Second/4, func() {
					cubeContainerObject.Mu.Lock()
					selected := cubeContainerObject.Selected
					cubeContainerObject.Mu.Unlock()
					if selected == nil {
						return
					}
					pointerLine.Resize(fyne.NewSize(1400-selected.Position().X-325, 2))
				})
			}
		}
	}()

	cubusWindow.ShowAndRun()
}
