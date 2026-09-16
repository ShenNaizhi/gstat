package main

import (
	"fmt"
	"flag"
)

// when given a path, scan() crawls it and its subfolders
// searching for Git repositories
func scan(folder string) {
	fmt.Printf("Found folders:\n\n")

	repositories := recursiveScanFolder(folder) // get a slice of strings
	filePath := getDotFilePath() // get the path of the dot file for output
	addNewRepoPathToFile(filePath, repositories) // output the slice contents to dot file

	fmt.Printf("\n\nSuccessfully added\n\n")
}

// stats() generates a graph of the local Git contributions
func stats(email string) {
	print("stats")
}

func main() {
	var folder, email string

	// define flags
	flag.StringVar(&folder, "add", "", "add a new folder to scan for Git repositories")
	flag.StringVar(&email, "email", "your@email.com", "the email to scan")

	// parse arguments
	flag.Parse()

	if folder != "" {
		scan(folder)
		return
	}

	stats(email)
}