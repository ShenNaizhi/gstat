package main

import (
	"time"
	"sort"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5"
)

type column []int

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
func fillCommits(email string, path string, commits map[int]int) map[int]int {
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
	iter, err := repo.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		panic(err)
	}

	// iterate on commits
	offset := weekdayOffset()
	err = iter.ForEach(func (c *object.Commit) error {
		if c.Author.Email != email { // not target user
			return nil
		}

		days := daysBetween(c.Author.When, time.Now())
		if days <= lookbackDays {
			commits[days + offset]++
		}

		return nil
	})
	if err != nil {
		panic(err)
	}

	return commits
}

// Counts how many days between 2 dates.
// Adapted from `countDaysSinceDate()` in the original tutorial.
func daysBetween(from, to time.Time) int {
	if from.After(to) { // ensure `from` is before `to`
		from, to = to, from
	}

	// unify into UTC 00:00 before calculation
	y1, m1, d1 := from.Date()
	y2, m2, d2 := to.Date()
	from = time.Date(y1, m1, d1, 0, 0, 0, 0, time.UTC)
	to = time.Date(y2, m2, d2, 0, 0, 0, 0, time.UTC)

	return int(to.Sub(from).Hours() / 24)
}

// Returns the amount of days missing to fill the last column of the stats graph.
// Adapted from `calcOffset()` in the original tutorial.
func weekdayOffset() int {
	res := int(time.Now().Weekday()) // American standard
	res = res % 7 + 1 // ISO 8601: Monday is the first day of a week

	return res
}

// Render the git commit stats in terminal.
func printCommitsStats(commits map[int]int) {
	keys := sortedMapKeys(commits)
	cols := buildCols(keys, commits)
	printCells(cols)
}

// Returns a slice of keys of `m` in ascending order.
// Adapted from 'sortMapIntoSlice()' in the original tutorial.
func sortedMapKeys(m map[int]int) []int {
	res := []int{}
	for k := range m {
		res = append(res, k)
	}
	sort.Ints(res)

	return res
}

func buildCols(keys []int, commits map[int]int) map[int]column {
	res := map[int]column{}
	col := column{}

	for _, k := range keys {
		week := k / 7
		weekday := k % 7

		if weekday == 0 { // the first day of a week
			col = column{} // reset
		}

		col = append(col, commits[k])

		if week == 6 { // the last day of a week
			res[week] = col // complete a column
		}
	}

	return res
}

// Render the whole stats graph
func printCells(cols map[int]column) {
	printMonths() // head of column

	totalWeeks := (lookbackDays - 1) / 7 + 1 // ceiling
	offset := weekdayOffset()
	for wd := 0; wd < 7; wd++ { // wd: weekday
		for w := totalWeeks; w >= 0; w-- { // w: current week index [0, totalWeeks)
			if w == totalWeeks { // head of row
				printWeekday(wd)
			}

			if c, ok := cols[w]; ok { // c: current column
				if w == 0 && wd == offset { // today
					printOneCell(c[wd], true)
				} else if len(c) > wd { // other dates
					printOneCell(c[wd], false)
				}
			} else {
				printOneCell(0, false)
			}
		}
		fmt.Println() // a row is rendered
	}
}