package main

import "flag"

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