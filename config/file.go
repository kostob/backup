package config

// GetSourceDirsFromConfigFile will return an array of source directories which should be backed up.
func GetSourceDirsFromConfigFile() []string {
	return []string{"foo", "bar", "foobar"}
}

// GetTargetDirFromConfigFile will return the target directory to which the backup should be written.
func GetTargetDirFromConfigFile() string {
	return ""
}
