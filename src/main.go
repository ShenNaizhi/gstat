package main

import "flag"

// when given a path, scan crawls it and its subfolders
// searching for Git repositories
func scan(path string) {
	print("scan")
}

// stats generates a graph of the local Git contributions
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