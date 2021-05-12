package structs

// Statistics will collect how many files/directories
// were there and how many new files/directories were found
type Statistics struct {
	OverallFiles       int
	OverallDirectories int
	NewDirectories     int
	NewFiles           int
}
