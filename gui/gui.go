package gui

import (
	"CUBUS-core/orchestrator/client"
	"CUBUS-core/shared"
	"CUBUS-core/shared/translation"
	"CUBUS-core/shared/types"
	"CUBUS-core/shared/types/gui"
	"context"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"log"
	"sync"
	"time"
)

var cubusWindow fyne.Window
var lastWindowSize fyne.Size
var T = translation.T

func selectCube(c *gui.Cube, infoContainerShape *canvas.Rectangle, pointerLine *canvas.Line, pointerTip *canvas.Circle, infoContainerText *widget.RichText) {
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

func startResizeListener(infoContainerShape *canvas.Rectangle, infoContainerText *widget.RichText, cubeContainerObject *gui.CubeContainer, ctx context.Context) {
	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			default:
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
		}
	}(ctx)
}

type Gui struct {
	app                 fyne.App
	defaults            *shared.Defaults
	ctx                 context.Context
	cancel              context.CancelFunc
	mu                  sync.Mutex
	cubeConfigs         []types.CubeConfig
	cubeContainerObject *gui.CubeContainer
	remoteAddresses     []string
}

func NewGui(app fyne.App, defaults *shared.Defaults) *Gui {
	ctx, cancel := context.WithCancel(context.Background())
	newGui := &Gui{app: app, defaults: defaults, ctx: ctx, cancel: cancel}
	translation.AddLanguageChangeListener(newGui.RecreateComponents)
	return newGui
}

func (g *Gui) RecreateComponents() {
	g.mu.Lock()
	g.cancel()
	g.ctx, g.cancel = context.WithCancel(context.Background())
	g.CreateGui()
	defer g.mu.Unlock()
}

func (g *Gui) getCubeConfigs() {
	g.remoteAddresses = g.app.Preferences().StringListWithFallback("remoteAddresses", []string{"http://localhost:25560"})
	g.app.Preferences().SetStringList("remoteAddresses", g.remoteAddresses)
	orchestratorClient := client.NewClient()
	g.cubeConfigs = make([]types.CubeConfig, 0)
	for _, remoteAddress := range g.remoteAddresses {
		cubes, err := orchestratorClient.GetAllCubes(remoteAddress)
		if err != nil {
			log.Println("Error fetching cubes:", err)
			continue
		}
		g.cubeConfigs = append(g.cubeConfigs, cubes...)
	}
}

func (g *Gui) addServerUrlIfNotExists(serverUrl string) {
	for _, remoteAddress := range g.remoteAddresses {
		if remoteAddress == serverUrl {
			return
		}
	}
	g.remoteAddresses = append(g.remoteAddresses, serverUrl)
	g.app.Preferences().SetStringList("remoteAddresses", g.remoteAddresses)
}

func (g *Gui) SetupGui() {
	g.getCubeConfigs()

	cubusWindow = g.app.NewWindow("CUBUS core")
	cubusWindow.Resize(WindowSize())
	cubusWindow.CenterOnScreen()
	cubusWindow.SetIcon(g.app.Icon())

	lastWindowSize = cubusWindow.Canvas().Size()

	g.cubeContainerObject = gui.NewCubeContainer(WindowWidth()*0.5, WindowHeight()*0.5, cubusWindow)
	for _, cubeConfig := range g.cubeConfigs {
		g.cubeContainerObject.AddCube(g.defaults.CubeAssetURL, cubeConfig.Id, cubeConfig)
	}
	g.cubeContainerObject.CenterCubes()
}

func (g *Gui) CreateGui() {
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

	selectCallback := func(c *gui.Cube) { selectCube(c, infoContainerShape, pointerLine, pointerTip, infoContainerText) }
	unselectCallback := func() { unselectCube(infoContainerShape, pointerLine, pointerTip, infoContainerText) }
	g.cubeContainerObject.SetSelectCallback(selectCallback)
	g.cubeContainerObject.SetUnselectCallback(unselectCallback)

	startResizeListener(infoContainerShape, infoContainerText, g.cubeContainerObject, g.ctx)

	windowMenu := fyne.NewMainMenu( // TODO: add section to manage remote addresses
		fyne.NewMenu(T("File"),
			fyne.NewMenuItem(T("Create a new Cube"), func() {
				setupDialog(
					cubusWindow,
					&g.cubeConfigs,
					g.defaults,
					g.cubeContainerObject,
					g.addServerUrlIfNotExists,
				)
			}),
			/*
					fyne.NewMenuItem(T("Export cube configs"), func() { // TODO: save server address of each cube in the config when exporting
						saveDialog := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
							if writer == nil {
								return
							}
							if err != nil {
								log.Println(T("Error exporting config: "), err)
								return
							}
							jsonData, err := json.Marshal(g.cubeConfigs)
							if err != nil {
								log.Println(T("Error exporting config: "), err)
								return
							}
							_, err = writer.Write(jsonData)
							if err != nil {
								log.Println(T("Error writing config file: "), err)
							}
							err = writer.Close()
							if err != nil {
								return
							}
						}, cubusWindow)
						saveDialog.SetFileName("config.json")
						saveDialog.Show()
					}),
					fyne.NewMenuItem(T("Import cube configs"), func() { // TODO: fix this.
				Let the user decide on what server which cube should be created
						openDialog := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
							if reader == nil {
								return
							}
							if err != nil {
								log.Println(T("Error importing config: "), err)
								return
							}
							jsonData, err := io.ReadAll(reader)
							if err != nil {
								log.Println(T("Error reading config file: "), err)
								return
							}
							var importedConfigs []map[string]interface{}
							err = json.Unmarshal(jsonData, &importedConfigs)
							if err != nil {
								log.Println(T("Error parsing config file: "), err)
								return
							}
							cubeStrings := make([]string, len(importedConfigs))
							for i, config := range importedConfigs {
								cubeStrings[i] = shared.ObjectToJsonString(config)
							}
							g.app.Preferences().SetStringList("cubes", cubeStrings)
							g.cubeContainerObject.ClearCubes()
							for _, cubeConfig := range importedConfigs {
								cubeConfigAsCorrectType := types.CubeConfig{
									Id:        cubeConfig["id"].(string),
									CubeType:  types.CubeType{Value: cubeConfig["type"].(string)},
									CubeName:  cubeConfig["name"].(string),
									PublicKey: nil,
								}
								g.cubeContainerObject.AddCube(g.defaults.CubeAssetURL, cubeConfig["id"].(string), cubeConfigAsCorrectType)
							}
							g.cubeContainerObject.CenterCubes()

							err = reader.Close()
							if err != nil {
								log.Println(T("Error closing config file: "), err)
								return
							}
						}, cubusWindow)
						openDialog.Show()
					}),
			*/
			fyne.NewMenuItem(T("Settings"), func() {
				settingsDialog(g.app)
			}),
		),
	)
	cubusWindow.SetMainMenu(windowMenu)

	mainContainer := container.NewWithoutLayout(g.cubeContainerObject.Container, pointerLine, pointerTip, infoContainerShape, infoContainerText)
	cubusWindow.SetContent(mainContainer)

	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				infoContainerShape.Refresh()
				pointerLine.Refresh()
				pointerTip.Refresh()
				g.cubeContainerObject.Mu.Lock()
				selected := g.cubeContainerObject.Selected
				g.cubeContainerObject.Mu.Unlock()
				if selected != nil {
					pointerTip.Move(fyne.NewPos(selected.Position().X+selected.CubeSize/2-5, selected.Position().Y+selected.CubeSize/2-40))
					pointerLine.Move(fyne.NewPos(selected.Position().X+selected.CubeSize/2, selected.Position().Y+selected.CubeSize/2-35))
					time.AfterFunc(time.Second/4, func() {
						g.cubeContainerObject.Mu.Lock()
						selected := g.cubeContainerObject.Selected
						g.cubeContainerObject.Mu.Unlock()
						if selected == nil {
							return
						}
						pointerLine.Resize(fyne.NewSize(WindowWidth()-selected.Position().X-WindowWidth()*0.3-25, 2))
					})
				}
			}
		}
	}(g.ctx)
}

func (g *Gui) Run() {
	g.SetupGui()
	g.CreateGui()
	cubusWindow.ShowAndRun()
}
