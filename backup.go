package main

import (
	"log"
	"os"

	"backup/backup"
	"backup/config"
)

func main() {
	log.Println("Let's try to backup what you want!")

	// this will be read from the command line or from the config later...
	// parse the command line to see what we should do
	if config.UseConfigFile() {
		log.Println("We will use the config file.")
	} else {
		log.Println("We will use the args you entered.")
		log.Printf("Source dir: %s\n", config.GetSourceDirFromCommandLine())
		log.Printf("Target dir: %s\n", config.GetTargetDirFromCommandLine())
	}

	// configSourceDir, configTargetDir := config.GetConfigFromArgs()
	//configSourceDir := "D:\\test\\source"
	//configTargetDir := "D:\\test\\target"
	/*
		if configSourceDir == "" || configTargetDir == "" {
			log.Fatalln("You have to give a source and a target directory")
			os.Exit(1)
		}

		previousTargetDirs, err := backup.GetPreviousBackupDirs(configTargetDir)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Previous backups:", previousTargetDirs)

		sourceDir := configSourceDir
		targetDir := configTargetDir + "/" + backup.GenerateDate()

		previousTargetDir := ""
		if len(previousTargetDirs) > 0 {
			previousTargetDir = configTargetDir + "/" + previousTargetDirs[len(previousTargetDirs)-1]
		}

		os.MkdirAll(filepath.FromSlash(targetDir), os.ModePerm)

		backup.WalkDir(sourceDir, targetDir, previousTargetDir)
	*/

	statistics := backup.Run(config.GetSourceDirFromCommandLine(), config.GetTargetDirFromCommandLine())

	log.Printf(
		"Statistics:\nOverall directories: %d\nOverall files: %d\nNew directories: %d\nNew files: %d",
		statistics.OverallDirectories,
		statistics.OverallFiles,
		statistics.NewDirectories,
		statistics.NewFiles)

	os.Exit(0)
}
