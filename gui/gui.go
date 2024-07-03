package gui

import (
	"CUBUS-core/shared"
	"CUBUS-core/shared/translation"
	"CUBUS-core/shared/types"
	"CUBUS-core/shared/types/gui"
	"encoding/json"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"io"
	"log"
	"time"
)

var cubusWindow fyne.Window
var lastWindowSize fyne.Size

func selectCube(c *gui.Cube, infoContainerShape *canvas.Rectangle, pointerLine *canvas.Line, pointerTip *canvas.Circle, infoContainerText *widget.RichText) {
	T := translation.T
	go func() {
		infoContainerShape.Resize(fyne.NewSize(WindowWidth()*0.3, WindowHeight()-80))
		infoContainerText.Resize(fyne.NewSize(WindowWidth()*0.3-10, WindowHeight()-105))
		infoContainerText.Segments = []widget.RichTextSegment{ // TODO: add a button to delete or edit the cube
			&widget.TextSegment{
				Text: T("Cube info") + "\n",
				Style: widget.RichTextStyle{
					TextStyle: fyne.TextStyle{Bold: true},
					ColorName: theme.ColorNameBackground,
				},
			},
			&widget.TextSegment{
				Text: T("Cube ID: ") + c.Id + "\n",
				Style: widget.RichTextStyle{
					ColorName: theme.ColorNameBackground,
				},
			},
			&widget.TextSegment{
				Text: T("Cube Name: ") + c.Config.CubeName + "\n",
				Style: widget.RichTextStyle{
					ColorName: theme.ColorNameBackground,
				},
			},
			&widget.TextSegment{
				Text: T("Cube Type: ") + c.Config.CubeType.Value + "\n",
				Style: widget.RichTextStyle{
					ColorName: theme.ColorNameBackground,
				},
			},
		}
		infoContainerText.Refresh()
		canvas.NewPositionAnimation(
			infoContainerShape.Position(),
			fyne.NewPos(WindowWidth()-WindowWidth()*0.3-25, 25),
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
			fyne.NewSize(WindowWidth()-c.Position().X-WindowWidth()*0.3-25, 2),
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
			fyne.NewPos(WindowWidth(), 25),
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

func startResizeListener(infoContainerShape *canvas.Rectangle, infoContainerText *widget.RichText, cubeContainerObject *gui.CubeContainer) {
	go func() {
		for {
			currentSize := cubusWindow.Canvas().Size()
			if !(currentSize == lastWindowSize) {
				lastWindowSize = currentSize
				if cubeContainerObject.Selected != nil {
					infoContainerShape.Move(fyne.NewPos(WindowWidth()-WindowWidth()*0.3-25, 25))
					infoContainerText.Move(fyne.NewPos(WindowWidth()-WindowWidth()*0.3-20, 35))

					infoContainerShape.Resize(fyne.NewSize(WindowWidth()*0.3, WindowHeight()-80))
					infoContainerText.Resize(fyne.NewSize(WindowWidth()*0.3-10, WindowHeight()-105))
				} else {
					infoContainerShape.Move(fyne.NewPos(WindowWidth(), 25))
					infoContainerText.Move(fyne.NewPos(WindowWidth(), 35))
					infoContainerShape.Resize(fyne.NewSize(0, 0))
					infoContainerText.Resize(fyne.NewSize(0, 0))
				}
				infoContainerShape.Refresh()
				infoContainerText.Refresh()

				cubeContainerObject.CenterCubes()
			}
		}
	}()
}

func Gui(cubusApp fyne.App, defaults *shared.Defaults) { // TODO: make this responsive
	T := translation.T
	cubeStrings := cubusApp.Preferences().StringListWithFallback("cubes", []string{})
	cubusApp.Preferences().SetStringList("cubes", cubeStrings)
	cubeConfigs := make([]map[string]interface{}, len(cubeStrings))
	for i, cubeString := range cubeStrings {
		cubeConfigs[i] = shared.JsonStringToObject(cubeString)
	}

	cubusWindow = cubusApp.NewWindow("CUBUS core")
	cubusWindow.Resize(WindowSize()) // 1400x900 is the default size of the unresponsive GUI
	cubusWindow.CenterOnScreen()
	cubusWindow.SetIcon(cubusApp.Icon())

	lastWindowSize = cubusWindow.Canvas().Size()

	infoContainerShape := canvas.NewRectangle(color.White)
	infoContainerShape.Resize(fyne.NewSize(WindowWidth()*0.3, WindowHeight()-80))
	infoContainerShape.Move(fyne.NewPos(WindowWidth(), 25))
	infoContainerShape.CornerRadius = 12

	infoContainerText := widget.NewRichText()
	infoContainerText.Resize(fyne.NewSize(WindowWidth()*0.3-10, WindowHeight()-105))
	infoContainerText.Wrapping = fyne.TextWrapBreak

	pointerLine := canvas.NewLine(color.White)
	pointerLine.StrokeWidth = 2
	pointerLine.Hide()
	pointerTip := canvas.NewCircle(color.White)
	pointerTip.Resize(fyne.NewSize(10, 10))
	pointerTip.Hide()

	cubeContainerObject := gui.NewCubeContainer(func() { unselectCube(infoContainerShape, pointerLine, pointerTip, infoContainerText) }, WindowWidth()*0.5, WindowHeight()*0.5, cubusWindow)
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

	startResizeListener(infoContainerShape, infoContainerText, cubeContainerObject)

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
			fyne.NewMenuItem(T("Export cube configs"), func() {
				saveDialog := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
					if writer == nil {
						return
					}
					if err != nil {
						log.Println("Error exporting config:", err)
						return
					}
					jsonData, err := json.Marshal(cubeConfigs)
					if err != nil {
						log.Println("Error exporting config:", err)
						return
					}
					_, err = writer.Write(jsonData)
					if err != nil {
						log.Println("Error writing config file:", err)
					}
					err = writer.Close()
					if err != nil {
						return
					}
				}, cubusWindow)
				saveDialog.SetFileName("config.json")
				saveDialog.Show()
			}),
			fyne.NewMenuItem(T("Import cube configs"), func() {
				openDialog := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
					if reader == nil {
						return
					}
					if err != nil {
						log.Println("Error importing config:", err)
						return
					}
					jsonData, err := io.ReadAll(reader)
					if err != nil {
						log.Println("Error reading config file:", err)
						return
					}
					var importedConfigs []map[string]interface{}
					err = json.Unmarshal(jsonData, &importedConfigs)
					if err != nil {
						log.Println("Error parsing config file:", err)
						return
					}
					cubeStrings := make([]string, len(importedConfigs))
					for i, config := range importedConfigs {
						cubeStrings[i] = shared.ObjectToJsonString(config)
					}
					cubusApp.Preferences().SetStringList("cubes", cubeStrings)
					cubeContainerObject.ClearCubes()
					for _, cubeConfig := range importedConfigs {
						cubeConfigAsCorrectType := types.CubeConfig{
							Id:        cubeConfig["id"].(string),
							CubeType:  types.CubeType{Value: cubeConfig["type"].(string)},
							CubeName:  cubeConfig["name"].(string),
							PublicKey: nil,
						}
						cubeContainerObject.AddCube(defaults.CubeAssetURL, func(c *gui.Cube) { selectCube(c, infoContainerShape, pointerLine, pointerTip, infoContainerText) }, cubeConfig["id"].(string), cubeConfigAsCorrectType)
					}
					cubeContainerObject.CenterCubes()

					err = reader.Close()
					if err != nil {
						return
					}
				}, cubusWindow)
				openDialog.Show()
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
					pointerLine.Resize(fyne.NewSize(WindowWidth()-selected.Position().X-WindowWidth()*0.3-25, 2))
				})
			}
		}
	}()

	cubusWindow.ShowAndRun()
}
