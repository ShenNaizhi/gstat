package main

import (
	"fmt"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5"
	"gopkg.in/src-d/go-git.v5"
)

const lookbackDays = 183

// stats() generates a graph of the local Git contributions
func stats(email string) {
	commits := processRepos(email) // get the list of commits
	printCommitStats(commits) // given commits, render the stats graph
}

// Given user's email, returns the commits made in last `lookbackDays` day(s).
// Adapted from `processRepositories()` in original tutorial.
func processRepos(email string) map[int]int {
	filePath := dotFilePath() // get the dot file path
	repos := parseFileLinesToSlice(filePath) // parse all lines of the file in a slice

	res := make(map[int]int, lookbackDays)
	for i := lookbackDays; i > 0; i-- { // initialize the commits map
		res[i] = 0
	}

	for _, path := range repos { // iterate over all repos, fill the commits map
		res = fillCommits(email, path, res)
	}

	return res
}

// Given a repo found in `path`, counts the frequency of commits and
// puts them in `commits`, returns it when completes.
func fillCommits(email string, path string, commits map[int][int]) map[int]int {
	// instantiate a git repo object from `path`
	repo, err := git.PlainOpen(path)
	if err != nil {
		panic(err)
	}

	// get the HEAD ref
	ref, err := repo.Head()
	if err != nil {
		panic(err)
	}

	// get the commits history starting from HEAD
	iter, err := repo.Log(&git.LogOptions({From: ref.Hash()}))
	if err != nil {
		panic(err)
	}

	// iterate on commits
	offset := calcWeekdayOffset()
	err = iter.ForEach(func (c *object.Commit) error {
		if c.Author.Email != email { // not target user
			return nil
		}

		days := countDaysSinceDate(c.Author.When) + offset
		if days <= lookbackDays {
			commits[days]++
		}

		return nil
	})
	if err != nil {
		panic(err)
	}

	return commits
}