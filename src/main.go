package main

import (
	"fmt"
	"strings"
	"log"
	"os"
	"os/user"
	"flag"
)

// when given a path, scan() crawls it and its subfolders
// searching for Git repositories
func scan(folder string) {
	fmt.Printf("Found folders:\n\n")

	repositories := recursiveScanFolder(folder) // get a slice of strings
	filePath := getDotFilePath() // get the path of the dot file for output
	addNewSliceElementsToFile(filePath, repositories) // output the slice contents to dot file

	fmt.Printf("\n\nSuccessfully added\n\n")
}

// starts the recursive search of git repo living in the `folder` subtree
func recursiveScanFolders(folder string) []string {
	return scanGitFolders([]string{}, folder)
}

// returns a list of subfolders of `folder` ending with `.git`.
// Returns the base folder of the repo, the `.git` folder parent.
// Recursively searches in the subfolders by passing an existing `folders` slice.
func scanGitFolders(folders []string, folder string) []string {
	// trim the last `/`
	folder = strings.TrimSuffix(folder, "/")

	f, err := os.Open(folder)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	files, err := f.ReadDir(-1)
	if err != nil {
		log.Fatal(err)
	}

	path := ""
	toSkip := map[string]bool {
		"vendor": true,
		"node_modules": true,
	}
	for _, file := range files {
		if !file.IsDir() { // only check directories
			continue
		}

		fileName := file.Name()

		if fileName == ".git" { // end with `.git`
			fmt.Println(path)
			folders = append(folders, path)
			continue
		}

		if toSkip[fileName] { // skip unrelated folders
			continue
		}

		path = path + "/" + fileName
		folders = scanGitFolders(folders, path)
	}

	return folders
}

// Returns the dot file path of the repos list.
// Create the dot file and the enclosing folder if they do not exist.
func getDotFilePath() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	res := usr.HomeDir + "/.gogitlocalstats"

	return res
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