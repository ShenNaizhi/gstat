package main

import (
	"fmt"
	"strings"
	"log"
	"os"
	"bufio"
	"path/filepath"
)

// when given a path, scan() crawls it and its subfolders
// searching for Git repositories
func scan(folder string) {
	fmt.Printf("Found folders:\n\n")

	repositories := recursiveScanFolder(folder) // get a slice of strings
	filePath := dotFilePath() // get the path of the dot file for output
	addNewRepoPathsToFile(filePath, repositories) // output the slice contents to dot file

	fmt.Printf("\n\nSuccessfully added.\n\n")
}

// starts the recursive search of git repo living in the `folder` subtree
func recursiveScanFolder(folder string) []string {
	return scanGitFolders([]string{}, folder)
}

// returns a list of subfolders of `folder` ending with `.git`.
// Returns the base folder of the repo, the `.git` folder parent.
// Recursively searches in the subfolders by passing an existing `folders` slice.
func scanGitFolders(folders []string, folder string) []string {
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
		if toSkip[fileName] { // skip unrelated folders
			continue
		}
		
		path = filepath.Join(folder, fileName)
		if fileName == ".git" { // end with `.git`
			fmt.Println(path)
			folders = append(folders, path)
			continue
		}

		folders = scanGitFolders(folders, path)
	}

	return folders
}

// Returns the dot file path of the repos list.
// Create the dot file and the enclosing folder if they do not exist.
// Adapted from `getDotFilePath()` from the original tutorial.
func dotFilePath() string {
	/*
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	res := usr.HomeDir + "/.gogitlocalstats"
	*/

	// Adapt from the original tutorial,
	// make the dot file and the whole program in the same path.
	res, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}

	res, err = filepath.EvalSymlinks(res) // to prevent res being a symlink
	if err != nil {
		log.Fatal(err)
	}

	res = filepath.Dir(res) // discard the executable's name
	res = filepath.Join(res, ".gogitlocalstats") // dot file path

	return res
}

// Give a slice of strings representing repo paths, store them in the file system.
// Adapted from `addNewSliceElementsToFile()` in original tutorial.
func addNewRepoPathsToFile(filePath string, repos []string) {
	existingRepos := parseFileLinesToSlice(filePath)
	toPersist := joinSlices(existingRepos, repos)
	persistRepoPathsToFile(toPersist, filePath)
}

// Give a file path string, it gets the content of each line, returns a slice of string.
func parseFileLinesToSlice(filePath string) []string {
	f := openFile(filePath) // no need to check, it will panic if necessary
	defer f.Close()

	lines := []string{}
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}

	if err := sc.Err(); err != nil {
		panic(err)
	}

	return lines
}

// Opens the target file at `filePath`. Creates it if not existing.
func openFile(filePath string) *os.File {
	res, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644) // change 0755 (original in tutorial) by 0644
	if err != nil {
		panic(err)
	}

	return res
}

// Append elems of `from` to `to`, without duplicates.
func joinSlices(to []string, from []string) []string {
	set := map[string]bool{}
	for _, el := range to {
		set[el] = true
	}

	for _, el := range from {
		if !set[el] {
			to = append(to, el)
			set[el] = true
		}
	}

	return to
}

// Writes contents of `repos` in file at `filePath` (totally overwrite).
// Adapted from `dumpStringsSliceToFile()` in original tutorial.
func persistRepoPathsToFile(repos []string, filePath string) {
	data := strings.Join(repos, "\n")
	os.WriteFile(filePath, []byte(data), 0644) // Change 0755 (original in tutorial) by 0644
}