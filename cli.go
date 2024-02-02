package main

import (
	"fmt"
	"strings"
)

const banner = `
 ██████╗██╗   ██╗██████╗ ██╗   ██╗███████╗
██╔════╝██║   ██║██╔══██╗██║   ██║██╔════╝
██║     ██║   ██║██████╔╝██║   ██║███████╗
██║     ██║   ██║██╔══██╗██║   ██║╚════██║
╚██████╗╚██████╔╝██████╔╝╚██████╔╝███████║
 ╚═════╝ ╚═════╝ ╚═════╝  ╚═════╝ ╚══════╝
`

func cli(config map[string]interface{}) {
	fmt.Print(banner)
	cubeTypeString := map[string]string{"q": "QUEEN", "s": "SECURITY", "d": "DRONE", "db": "DATABASE", "a": "API", "w": "WEBSITE", "b": "DISCORD BOT"}[strings.ToUpper(config["cube_type"].(string))]
	fmt.Printf("%*s\n\n", 42, cubeTypeString)
	fmt.Println("> ")
}
