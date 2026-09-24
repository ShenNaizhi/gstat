package main

import (
	"fmt"
	"time"
	"sort"
	"slices"
	"strings"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type column []int

const lookbackDays = 183
var currentTime time.Time = time.Now()

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
	for i := 0; i < lookbackDays; i++ { // initialize the commits map
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
	offset := weekdayOffset(currentTime)
	err = iter.ForEach(func (c *object.Commit) error {
		if c.Author.Email != email { // not target user
			return nil
		}

		days := daysBetween(c.Author.When, currentTime)
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

	// unify into UTC midnight before calculation
	from = normalizeToUTCMidnight(from)
	to = normalizeToUTCMidnight(to)

	return int(to.Sub(from).Hours() / 24)
}

// Normalize a time to UTC midnight.
func normalizeToUTCMidnight(t time.Time) time.Time {
	y, m, d := t.Date()

	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}

// Make up the last week.
// Adapted from `calcOffset()` in the original tutorial.
func weekdayOffset(t time.Time) int {
	res := int(t.Weekday()) // American standard
	res = (res - 1 + 7) % 7 // ISO 8601: Monday is the first day of a week

	return 7 - (res + 1)
}

// Render the git commit stats in terminal.
func printCommitStats(commits map[int]int) {
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
		weekday := 6 - k % 7

		if weekday == 6 { // the last day of a week
			col = column{} // reset
		}

		col = append(col, commits[k])

		if weekday == 0 { // the first day of a week
			slices.Reverse(col)
			res[week] = col // complete a column
		}
	}

	return res
}

// Render the whole stats graph
func printCells(cols map[int]column) {
	printMonths() // head of column

	offset := weekdayOffset(currentTime)
	totalWeeks := (lookbackDays + offset - 1) / 7 + 1 // ceiling
	for wd := 0; wd < 7; wd++ { // wd: weekday
		for w := totalWeeks; w >= 0; w-- { // w: current week index [0, totalWeeks)
			if w == totalWeeks { // head of row
				printWeekday(wd)
				continue
			}

			if c, ok := cols[w]; ok { // c: current column
				if w == 0 && wd == 6 - offset { // today
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

// Renders the month names in the first line of the whole stat graph.
func printMonths() {
	week := normalizeToUTCMidnight(currentTime).AddDate(0, 0, -lookbackDays) // during a certain week
	month := week.Month()
	cur := currentTime

	fmt.Print(strings.Repeat(" ", 9))
	// fmt.Printf("%s ", week.Month().String()[:3]) // abbreviation of weekdays
	for {
		if week.Month() != month { // next month
			fmt.Printf("%s ", week.Month().String()[:3]) // abbreviation of weekdays
			month = week.Month()
		} else { // still in current month
			fmt.Print(strings.Repeat(" ", 4))
		}

		week = week.AddDate(0, 0, 7) // move to next week

		if week.After(cur) && !inSameWeek(week, cur) { // done
			break
		}
	}
	fmt.Println()
}

// Checks whether 2 given time is in the same week
func inSameWeek(a, b time.Time) bool {
	a = normalizeToUTCMidnight(a)
	b = normalizeToUTCMidnight(b)

	if a.After(b) {
		a, b = b, a // ensure a is before b
	}

	gap := int(b.Sub(a).Hours() / 24)
	if gap > 6 {
		return false // not in the same week
	}

	return weekdayOffset(a) >= weekdayOffset(b)
}

// Prints the abbreviation of the weekday.
func printWeekday(wd int) {
	res := strings.Repeat(" ", 5)
	switch wd {
	case 0: res = " Mon "
	case 2: res = " Wed "
	case 4: res = " Fri "
	case 6: res = " Sun "
	}

	fmt.Print(res)
}

// Prints `val` in different formats depending on
// it's value and `isToday` flag.
func printOneCell(val int, isToday bool) {
	esc := "\033[0;30m" // black background

	cntLvl := []int{0, 4, 9} // commit count level
	colorLvl := []string{"47", "43", "42"} // color level
	for i, cnt := range cntLvl {
		if val > cnt {
			esc = "\033[1;30;" + colorLvl[i] + "m"
		} else {
			break
		}
	}

	if isToday {
		esc = "\033[1;37;45m"
	}

	escReset := "\033[0m"

	if val == 0 {
		fmt.Print(esc + "  - " + escReset)
		return
	}

	fmt.Printf(esc + "%3d " + escReset, val)
}