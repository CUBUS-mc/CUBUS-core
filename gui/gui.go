package gui

import (
	"CUBUS-core/shared"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"time"
)

func selectCube(c *cube, infoContainerShape *canvas.Rectangle, pointerLine *canvas.Line, pointerTip *canvas.Circle, infoContainerText *widget.RichText) {
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
				Text: "Cube ID: " + c.id + "\n",
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
		pointerTip.Move(fyne.NewPos(c.Position().X+c.size/2-5, c.Position().Y+c.size/2-5))
		pointerTip.Show()
		pointerLine.Resize(fyne.NewSize(0, 0))
		pointerLine.Move(fyne.NewPos(c.Position().X+c.size/2, c.Position().Y+c.size/2))
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
	cubeStrings := cubusApp.Preferences().StringListWithFallback("cubes", []string{})
	cubusApp.Preferences().SetStringList("cubes", cubeStrings)
	cubeConfigs := make([]map[string]interface{}, len(cubeStrings))
	for i, cubeString := range cubeStrings {
		cubeConfigs[i] = shared.JsonStringToObject(cubeString)
	}

	cubusWindow := cubusApp.NewWindow("QUBUS core")
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

	cubeContainerObject := newCubeContainer(func() { unselectCube(infoContainerShape, pointerLine, pointerTip, infoContainerText) }, 500, 300)
	for _, cubeConfig := range cubeConfigs {
		cubeContainerObject.AddCube(defaults.CubeAssetURL, func(c *cube) { selectCube(c, infoContainerShape, pointerLine, pointerTip, infoContainerText) }, cubeConfig["id"].(string))
	}
	cubeContainerObject.CenterCubes()

	windowMenu := fyne.NewMainMenu(
		fyne.NewMenu("File",
			fyne.NewMenuItem("Create a new Cube", func() {
				setupDialog(cubusWindow) // TODO: move the preference update to the setupDialog function and only update the cubeContainerObject with the new cube
				// cubeConfigs = append(cubeConfigs, map[string]interface{}{"id": cubeConfig.id})
				// cubeStrings = []string{}
				// for _, cubeConfig := range cubeConfigs {
				// 	cubeStrings = append(cubeStrings, shared.ObjectToJsonString(cubeConfig))
				// }
				// cubusApp.Preferences().SetStringList("cubes", cubeStrings)
				// cubeContainerObject.AddCube(defaults.CubeAssetURL, func(c *cube) { selectCube(c, infoContainerShape, pointerLine, pointerTip, infoContainerText) }, cubeConfig.id)
				// cubeContainerObject.CenterCubes()
			}),
			fyne.NewMenuItem("Export config", func() {
				println("Export config") // TODO: implement config export
			}),
			fyne.NewMenuItem("Import config", func() {
				println("Import config") // TODO: implement config import
			}),
			fyne.NewMenuItem("Settings", func() {
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
			cubeContainerObject.mu.Lock()
			selected := cubeContainerObject.selected
			cubeContainerObject.mu.Unlock()
			if selected != nil {
				pointerTip.Move(fyne.NewPos(selected.Position().X+selected.size/2-5, selected.Position().Y+selected.size/2-40))
				pointerLine.Move(fyne.NewPos(selected.Position().X+selected.size/2, selected.Position().Y+selected.size/2-35))
				time.AfterFunc(time.Second/4, func() {
					cubeContainerObject.mu.Lock()
					selected := cubeContainerObject.selected
					cubeContainerObject.mu.Unlock()
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
