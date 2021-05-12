package backup

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"backup/structs"

	"github.com/udhos/equalfile"
)

var statistics structs.Statistics

// Run will start the backup process
func Run(sourceDir string, targetDir string) structs.Statistics {
	// init the statistics
	statistics = structs.Statistics{OverallFiles: 0, OverallDirectories: 0, NewDirectories: 0, NewFiles: 0}

	previousTargetDirs, err := getPreviousBackupDirs(targetDir)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Previous backups:", previousTargetDirs)

	// directory where the new backup will be placed
	newTargetDir := targetDir + "/" + generateDate()

	// directory of the previous backup, if there is any
	previousTargetDir := ""
	if len(previousTargetDirs) > 0 {
		previousTargetDir = targetDir + "/" + previousTargetDirs[len(previousTargetDirs)-1]
	}

	// create the target dir (if not already existing)
	os.MkdirAll(filepath.FromSlash(targetDir), os.ModePerm)

	// start walking through the source directory
	walkDir(sourceDir, newTargetDir, previousTargetDir)

	return statistics
}

func walkDir(sourceDir string, targetDir string, previousTargetDir string) {
	sourceDir = filepath.FromSlash(sourceDir)
	targetDir = filepath.FromSlash(targetDir)
	previousTargetDir = filepath.FromSlash(previousTargetDir)
	// println("sourceDir: " + sourceDir)
	// println("targetDir: " + targetDir)
	// println("previousTargetDir: " + previousTargetDir)

	dirContent, errReadDir := ioutil.ReadDir(sourceDir)
	if errReadDir != nil {
		log.Fatal(errReadDir)
	}

	for _, c := range dirContent {
		if c.IsDir() {
			// increase counter for overall count of directories
			statistics.OverallDirectories++

			newSourceDir := filepath.FromSlash(sourceDir + "/" + c.Name())
			newTargetDir := filepath.FromSlash(targetDir + "/" + c.Name())
			newPreviousTargetDir := ""
			if previousTargetDir != "" {
				newPreviousTargetDir = filepath.FromSlash(previousTargetDir + "/" + c.Name())
			}

			// check if the directory already existed in the previous backup
			previousTargetName := filepath.FromSlash(previousTargetDir + "/" + c.Name())
			if _, err := os.Stat(previousTargetName); os.IsNotExist(err) {
				log.Printf("New directory detected: %s", newSourceDir)
				statistics.NewDirectories++
			}

			// create dir
			os.MkdirAll(newTargetDir, os.ModePerm)

			walkDir(newSourceDir, newTargetDir, newPreviousTargetDir)
		} else {
			// increase counter for overall count of files
			statistics.OverallFiles++

			sourceName := filepath.FromSlash(sourceDir + "/" + c.Name())
			targetName := filepath.FromSlash(targetDir + "/" + c.Name())

			if previousTargetDir != "" {
				previousTargetName := filepath.FromSlash(previousTargetDir + "/" + c.Name())

				// check if the file is new => look into the last backup
				_, errStat := os.Stat(previousTargetName)
				if errStat != nil {
					// copy the file
					if os.IsNotExist(errStat) {
						copyFile(sourceName, targetName)
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
						linkFile(previousTargetName, targetName)
					}
				}
			} else {
				// we don't have the file anyway, let's copy it
				copyFile(sourceName, targetName)
			}
		}
	}
}

// GenerateDate will generate the current date which later will be used as directory name for the new backup
func generateDate() string {
	time := time.Now()
	//Format is <year>-<month>-<day> <hour>:<minute>:<second> (all 2 digits except year which has 4)
	formattedTime := time.Format("20060102030405")

	return formattedTime
}

// GetPreviousBackupDirs will return an array of all previous backup directories
func getPreviousBackupDirs(targetDir string) ([]string, error) {
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

func copyFile(source string, target string) {
	log.Printf("Copy file %s to %s", source, target)

	// increase counter for new files
	statistics.NewFiles++

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

func linkFile(destination string, new string) {
	log.Printf("Creating hardlink %s to destination %s", new, destination)
	errLink := os.Link(destination, new)
	if errLink != nil {
		log.Fatal(errLink)
	}
}
