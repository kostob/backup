package config

import (
	"flag"
	"log"
)

var useConfigFile bool
var configFilePath string
var useCommandLine bool
var sourcePath string
var targetPath string

func init() {
	defineCommandLine()

	if configFilePath != "" {
		useConfigFile = true
	}

	if sourcePath != "" || targetPath != "" {
		useCommandLine = true
	}

	parseCommandLine()
}

func defineCommandLine() {
	flag.StringVar(&configFilePath, "c", "", "Use the given config file")
	flag.StringVar(&sourcePath, "s", "", "The source directory to use")
	flag.StringVar(&targetPath, "t", "", "The target directory to use")

	flag.Parse()
}

func parseCommandLine() {
	// check that we have a valid command line
	if useConfigFile == true && useCommandLine == true {
		log.Fatal("Invalid command line. You cannot use a config file and set source and target over the command line.")
	}

	if useCommandLine == true && (sourcePath == "" || targetPath == "") {
		log.Fatal("Invalid command line. If you want to set the source and target path over the command line, you have to add both.")
	}

	// if we have "c" in in the command line we should use this file

	// if there is "s" and "t" we should use the values from the command line

	// if there is nothing, we will use the config file which is placed in the user dir

	if useCommandLine == false && useConfigFile == false {
		log.Fatal("Invalid command line. You have to select either to use a config file or you have to set the source and target directory.")
	}
}

// UseConfigFile tells if a config file will be used (true) or if the source and target paths for the command line will be used
func UseConfigFile() bool {
	return useConfigFile
}

// GetSourceDirFromCommandLine will return the source dir given on the command line
func GetSourceDirFromCommandLine() string {
	if useCommandLine == true {
		return sourcePath
	}

	return ""
}

// GetTargetDirFromCommandLine will return the target dir given on the command line
func GetTargetDirFromCommandLine() string {
	if useCommandLine == true {
		return targetPath
	}

	return ""
}
