package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"strings"
)

var configDialogueTree = questionSequence{
	Name: "",
	SimpleNodes: []any{
		multiSelectNode{
			Name:     "cube_type",
			Question: "What type of cube do you want to setup?",
			Prompt:   "Type of cube: ",
			Options: map[string]interface{}{
				"q": questionSequence{
					Name:        "Queen",
					SimpleNodes: []any{},
				},
				"s": questionSequence{
					Name:        "Security",
					SimpleNodes: []any{},
				},
				"d": questionSequence{
					Name:        "Drone",
					SimpleNodes: []any{},
				},
				"db": questionSequence{
					Name:        "Database",
					SimpleNodes: []any{},
				},
				"a": questionSequence{
					Name:        "API",
					SimpleNodes: []any{},
				},
				"w": questionSequence{
					Name:        "Website",
					SimpleNodes: []any{},
				},
				"b": questionSequence{
					Name:        "Discord bot",
					SimpleNodes: []any{},
				},
			},
			Default: "d",
		},
	},
	AdvancedNodes: []any{
		multiSelectNode{
			Name:     "ui_type",
			Question: "What type of UI do you want to use?",
			Prompt:   "Type of UI: ",
			Options: map[string]interface{}{
				"c": valueNode{
					Name: "CLI",
				},
				"g": valueNode{
					Name: "GUI",
				},
			},
			Default: "g",
		},
	},
}

type AdvancedSetting int

const (
	None AdvancedSetting = iota
	True
	False
)

type multiSelectNode struct {
	Name     string
	Question string
	Prompt   string
	Options  map[string]interface{}
	Default  string
}

type booleanNode struct {
	Name     string
	Question string
	Yes      questionSequence
	No       questionSequence
	Default  bool
}

type valueNode struct {
	Name string
}

type questionSequence struct {
	Name          string
	SimpleNodes   []any
	AdvancedNodes []any
}

func handleQuestionSequenceCLI(node questionSequence, config map[string]interface{}, advanced AdvancedSetting) {
	for _, n := range node.SimpleNodes {
		walkTreeCLI(n, config, advanced)
	}
	if advanced == None && len(node.AdvancedNodes) > 0 {
		advancedBool := askForBoolean("Do you want to configure advanced settings? (y/n) ", false)
		if advancedBool {
			advanced = True
		} else {
			advanced = False
		}
	}
	if advanced == True {
		for _, n := range node.AdvancedNodes {
			walkTreeCLI(n, config, advanced)
		}
	} else if advanced == False {
		for _, n := range node.AdvancedNodes {
			config[n.(multiSelectNode).Name] = n.(multiSelectNode).Default
		}
	}
}

func handeleSelectNodeCLI(node multiSelectNode, config map[string]interface{}) {
	println(node.Question)
	var options []string
	for key, value := range node.Options {
		switch v := value.(type) {
		case questionSequence:
			options = append(options, fmt.Sprintf("%s - %s", key, v.Name))
		case valueNode:
			options = append(options, fmt.Sprintf("%s - %s", key, v.Name))
		}
	}
	println(fmt.Sprintf("%s (default: %s)", strings.Join(options, ", "), node.Default))
	config[node.Name] = strings.ToLower(askForString(node.Prompt, node.Default, func(input string) bool {
		_, ok := node.Options[input]
		return ok
	}))
}

func handeleBooleanNodeCLI(node booleanNode, config map[string]interface{}) {
	config[node.Name] = askForBoolean(node.Question, node.Default)
}

func walkTreeCLI(node any, config map[string]interface{}, advanced AdvancedSetting) {
	switch node.(type) {
	case questionSequence:
		handleQuestionSequenceCLI(node.(questionSequence), config, advanced)
	case multiSelectNode:
		handeleSelectNodeCLI(node.(multiSelectNode), config)
	case booleanNode:
		handeleBooleanNodeCLI(node.(booleanNode), config)
	}
}

func handleQuestionSequenceGUI(node questionSequence, config map[string]interface{}, parent *fyne.Container) {
	for _, n := range node.SimpleNodes {
		walkTreeGUI(n, config, parent)
	}
	if len(node.AdvancedNodes) > 0 {
		for _, n := range node.AdvancedNodes {
			config[n.(multiSelectNode).Name] = n.(multiSelectNode).Default
		}
		advancedContainer := container.NewVBox()
		advancedContainer.Hide()
		advancedBool := widget.NewCheck("Do you want to configure advanced settings?", func(b bool) {
			advancedContainer.Objects = nil
			for _, n := range node.AdvancedNodes {
				config[n.(multiSelectNode).Name] = n.(multiSelectNode).Default
			}
			if b {
				advancedContainer.Show()
				for _, n := range node.AdvancedNodes {
					walkTreeGUI(n, config, advancedContainer)
				}
			} else {
				advancedContainer.Hide()
				for _, n := range node.AdvancedNodes {
					config[n.(multiSelectNode).Name] = n.(multiSelectNode).Default
				}
			}
		})
		parent.Add(advancedBool)
		parent.Add(advancedContainer)
	}
}

func handeleSelectNodeGUI(node multiSelectNode, config map[string]interface{}, parent *fyne.Container) {
	config[node.Name] = node.Default
	options := make([]string, 0, len(node.Options))
	reverseMap := make(map[string]string)
	for key, value := range node.Options {
		switch v := value.(type) {
		case questionSequence:
			optionName := v.Name
			if key == node.Default {
				optionName += " (default)"
				options = append([]string{optionName}, options...)
			} else {
				options = append(options, optionName)
			}
			reverseMap[optionName] = key
		case valueNode:
			optionName := v.Name
			if key == node.Default {
				optionName += " (default)"
				options = append([]string{optionName}, options...)
			} else {
				options = append(options, optionName)
			}
			reverseMap[optionName] = key
		}
	}
	selectLabel := widget.NewLabel(node.Prompt)
	selectWidget := widget.NewSelect(options, func(input string) {
		key := reverseMap[input]
		config[node.Name] = key
		walkTreeGUI(node.Options[key], config, parent)
	})
	selectWidget.SetSelected(options[0])
	parent.Add(container.NewHBox(selectLabel, selectWidget))
}

func handeleBooleanNodeGUI(node booleanNode, config map[string]interface{}, parent *fyne.Container) {
	config[node.Name] = node.Default
	checkbox := widget.NewCheck(node.Question, func(b bool) {
		config[node.Name] = b
		if b {
			walkTreeGUI(node.Yes, config, parent)
		} else {
			walkTreeGUI(node.No, config, parent)
		}
	})
	parent.Add(checkbox)
}

func walkTreeGUI(node any, config map[string]interface{}, parent *fyne.Container) {
	switch node.(type) {
	case questionSequence:
		handleQuestionSequenceGUI(node.(questionSequence), config, parent)
	case multiSelectNode:
		handeleSelectNodeGUI(node.(multiSelectNode), config, parent)
	case booleanNode:
		handeleBooleanNodeGUI(node.(booleanNode), config, parent)
	}
}
