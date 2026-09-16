package main

import (
	"fmt"
	"strings"
	"log"
	"os"
	"os/user"
)

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

// Give a slice of strings representing repo paths, store them in the file system.
// Adapted from `addNewSliceElementsToFile()`.
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
	res, err := os.Open(filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644) // change 0755 (original in tutorial) by 0644
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
		}
	}

	return to
}

// Writes contents of `repos` in file at `filePath` (totally overwrite).
func persistRepoPathsToFile(repos []string, filePath string) {
	data := strings.Join(repos, "\n")
	os.WriteFile(filePath, []byte(data), 0644) // Change 0755 (original in tutorial) by 0644
}