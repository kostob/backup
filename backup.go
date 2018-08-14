package main

import (
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/udhos/equalfile"
)

func copy(source string, target string) {
	log.Printf("Copy file %s to %s", source, target)

	from, err := os.Open(source)
	if err != nil {
		log.Fatal(err)
	}
	defer from.Close()

	to, err := os.OpenFile(target, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer to.Close()

	_, err = io.Copy(to, from)
	if err != nil {
		log.Fatal(err)
	}
}

func link(destination string, new string) {
	log.Printf("Creating hardlink %s to destination %s", new, destination)
	errLink := os.Link(destination, new)
	if errLink != nil {
		log.Fatal(errLink)
	}
}

func walkDir(sourceDir string, targetDir string, previousTargetDir string) {
	dirContent, errReadDir := ioutil.ReadDir(filepath.FromSlash(sourceDir))
	if errReadDir != nil {
		log.Fatal(errReadDir)
	}

	for _, c := range dirContent {
		if c.IsDir() {
			newSourceDir := filepath.FromSlash(sourceDir + "/" + c.Name())
			newTargetDir := filepath.FromSlash(targetDir + "/" + c.Name())
			newPreviousTargetDir := ""
			if previousTargetDir != "" {
				newPreviousTargetDir = filepath.FromSlash(previousTargetDir + "/" + c.Name())
			}

			// create dir
			os.MkdirAll(newTargetDir, os.ModePerm)

			walkDir(newSourceDir, newTargetDir, newPreviousTargetDir)
		} else {
			sourceName := filepath.FromSlash(sourceDir + "/" + c.Name())
			targetName := filepath.FromSlash(targetDir + "/" + c.Name())
			previousTargetName := ""
			if previousTargetDir != "" {
				previousTargetName = filepath.FromSlash(previousTargetDir + "/" + c.Name())

				// check if the file is new => look into the last backup
				_, errStat := os.Stat(previousTargetName)
				if errStat != nil {
					// copy the file
					if os.IsNotExist(errStat) {
						copy(sourceName, targetName)
					} else {
						log.Fatal(errStat)
					}
				} else {
					// check if the file is the same as in the previous backup => make a link to previous backup
					cmp := equalfile.New(nil, equalfile.Options{})
					equal, errCompare := cmp.CompareFile(sourceName, previousTargetName)
					if errCompare != nil {
						log.Fatal(errCompare)
					}

					if equal {
						link(previousTargetName, targetName)
					}
				}
			} else {
				// we dont have the file anyway, let's copy it
				copy(sourceName, targetName)
			}
		}
	}
}

func generateDate() string {
	time := time.Now()
	//Format is <year><month><day><hour><minute><second> (all 2 digits except year which has 4)
	formattedTime := time.Format("20060102030405")

	return formattedTime
}

func getPrevoiusBackupDirs(targetDir string) ([]string, error) {
	var files []string
	fileInfo, err := ioutil.ReadDir(filepath.FromSlash(targetDir))
	if err != nil {
		return files, err
	}

	for _, file := range fileInfo {
		files = append(files, file.Name())
	}
	return files, nil
}

func getConfigFromArgs() (string, string) {
	flagSource := flag.String("s", "", "The source directory")
	flagTarget := flag.String("t", "", "The target directory")

	flag.Parse()

	return *flagSource, *flagTarget
}

func main() {
	log.Println("Let's backup what you want!")

	// this will be read from the command line or from the config later...
	configSourceDir, configTargetDir := getConfigFromArgs()
	//configSourceDir := "D:\\test\\source"
	//configTargetDir := "D:\\test\\target"

	if configSourceDir == "" || configTargetDir == "" {
		log.Fatalln("You have to give a source and a target directory")
		os.Exit(1)
	}

	previousTargetDirs, err := getPrevoiusBackupDirs(configTargetDir)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Previous backups:", previousTargetDirs)

	sourceDir := configSourceDir
	targetDir := configTargetDir + "/" + generateDate()

	previousTargetDir := ""
	if len(previousTargetDirs) > 0 {
		previousTargetDir = configTargetDir + "/" + previousTargetDirs[len(previousTargetDirs)-1]
	}

	os.MkdirAll(filepath.FromSlash(targetDir), os.ModePerm)

	walkDir(sourceDir, targetDir, previousTargetDir)

	os.Exit(0)
}
