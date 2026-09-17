package main

import "fmt"

const DefaultLookbackMonths = 6

func stats(email string) {
	fmt.Print("stats")
}

/*
// stats() generates a graph of the local Git contributions
func stats(email string) {
	commits := processRepos(email) // get the list of commits
	printCommitStats(commits) // given commits, render the stats graph
}

// Given user's email, returns the commits made in last `DefaultLookbackMonths` months.
// Adapted from `processRepositories()` in original tutorial.
func processRepos(email string) map[int]int {
	filePath := getDotFilePath() // get the dot file path
	repos := parseFileLinesToSlice(filePath) // parse all lines of the file in a slice

	res := make(map[int]int, DefaultLookbackMonths)
	for i := DefaultLookbackMonths; i > 0; i-- { // initialize the commits map
		res[i] = 0
	}

	for _, path := range repos { // iterate over all repos, fill the commits map
		res = fillCommits(email, path, res)
	}

	return res
}
*/