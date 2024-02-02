package main

import (
	"fmt"
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
	Yes      map[string]interface{}
	No       map[string]interface{}
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

func walkTreeCLI(node any, config map[string]interface{}, advanced AdvancedSetting) {
	switch node.(type) {
	case questionSequence:
		handleQuestionSequenceCLI(node.(questionSequence), config, advanced)
	case multiSelectNode:
		handeleSelectNodeCLI(node.(multiSelectNode), config)
	}
}
