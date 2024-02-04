package main

import (
	"io/fs"
	"log"
	"os"

	"github.com/go-git/go-git/v5"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func getCube(cubeType string) {
	switch cubeType {
	case "q":
		if !checkPathExists("queen") {
			log.Println("Queen cube not found, installing...")
			downloadCube("queen")
		} else {
			checkForUpdates("queen")
		}
	}
}

func downloadCube(cubeName string) {
	log.Println("Downloading " + cubeName + " cube")
	log.Println("[1/2] Creating " + cubeName + " directory")
	err := os.Mkdir(cubeName, fs.ModeDir)
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Println("[2/2] Cloning " + cubeName + " repository")
	_, err = git.PlainClone(cubeName, false, &git.CloneOptions{
		URL:      "https://github.com/CUBUS-mc/CUBUS-" + cubeName + ".git",
		Progress: os.Stdout,
	})
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Println(cases.Title(language.English).String(cubeName) + " downloaded successfully")
}

func checkForUpdates(cubeName string) {
	log.Println("Checking for updates for " + cubeName + " cube")
	repo, err := git.PlainOpen(cubeName)
	if err != nil {
		log.Fatal(err)
		return
	}
	wt, err := repo.Worktree()
	if err != nil {
		log.Fatal(err)
		return
	}
	err = wt.Pull(&git.PullOptions{})
	if err != nil {
		if err.Error() == "already up-to-date" {
			log.Println(cases.Title(language.English).String(cubeName) + " is up to date")
			return
		}
		log.Fatal(err)
		return
	}
	log.Println(cases.Title(language.English).String(cubeName) + " updated successfully")
}
