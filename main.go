package main

func main() {
	configuration := loadConfig()
	switch configuration["ui_type"].(string) {
	case "c":
		cli(configuration)
	case "g":
		gui(configuration)
	}
}
