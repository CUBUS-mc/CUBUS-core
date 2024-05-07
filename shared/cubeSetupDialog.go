package shared

import "crypto"

type CubeType struct {
	Label       string
	Description string
}

var CubeTypes = map[string]CubeType{
	"queen": {
		Label:       "Queen",
		Description: "The Queen is the central cube in the CUBUS network. It is the C2 server. There can only be one Queen in a CUBUS subnet.",
	},
	"security": {
		Label:       "Security",
		Description: "The Security cube is responsible for the security of the CUBUS network.",
	},
	"database": {
		Label:       "Database",
		Description: "The Database cube is responsible for storing data for the CUBUS network.",
	},
	"api": {
		Label:       "API",
		Description: "The API cube is responsible for providing an API for the CUBUS network.",
	},
	"cubus-mod": {
		Label:       "CUBUS Mod",
		Description: "This is a special cube that runs in a Minecraft mod.",
	},
	"discord-bot": {
		Label:       "Discord Bot",
		Description: "The Discord Bot cube is responsible for providing a Discord bot for the CUBUS network.",
	},
	"web": {
		Label:       "Web",
		Description: "The Web cube is responsible for providing a web interface for the CUBUS network.",
	},
	"drone": {
		Label:       "Drone",
		Description: "The Drone cubes are responsible for scanning, modifying and probiding a proxy for the CUBUS network.",
	},
	"generic-worker": {
		Label:       "Generic Worker",
		Description: "The Generic Worker cube is a generic cube that can be used for any purpose by any other cube.",
	},
}

type CubeConfigType struct {
	Id        string
	CubeType  CubeType
	PublicKey crypto.PublicKey
}

// TODO: Define a method to read a dialoge tree from a file
// TODO: Create a dialog tree in a file for the setup dialog
