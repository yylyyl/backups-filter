package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"time"
)

const (
	layoutDateTime = "20060102_150405"
	layoutDate     = "20060102"
)

var (
	fLayout = flag.String("layout", layoutDateTime, "layout of datetime, but this tool "+
		"only takes the dates into consideration")
	fBackupDaily   = flag.Uint("backup-daily", 7, "number of daily backups")
	fBackupWeekly  = flag.Uint("backup-weekly", 3, "number of weekly backups")
	fBackupMonthly = flag.Uint("backup-monthly", 5, "number of monthly backups")
	fDescending    = flag.Bool("descending", false, "descending input and output")
	fKeep          = flag.Bool("keep", false, "print the items which should be kept, "+
		"instead of which should be deleted")

	debug = false
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		usage()
	}
}

func usage() {
	fmt.Println("")
	fmt.Println("This tool helps you remove old daily backups. For example:")
	fmt.Println("  - Keep 7 daily backups (-backup-daily 7), and")
	fmt.Println("  - Keep 3 weekly backups (-backup-weekly 3), and")
	fmt.Println("  - Keep 5 monthly backups (-backup-monthly 5).")
	fmt.Println("  There will be 15 backups to be kept in total.")
	fmt.Println("  Other backups should be deleted.")
	fmt.Println("")
	fmt.Println("To use this tool correctly, you should do these things in your back up scripts:")
	fmt.Println("1. Do the ordinary backup procedures, normally uploading your files")
	fmt.Printf("2. %s\n", "Write the current date into a file, e.g. `date +%Y%m%d_%H%M%S >> backups.txt`")
	fmt.Println("3. Invoke this tool, e.g. `TO_DELETE=$(cat backups.txt | backups-filter)`")
	fmt.Println("4. Delete your old backups, and remove the date line from backups.txt respectively")
}

func main() {
	flag.Parse()

	if *fBackupDaily == 0 {
		fmt.Fprintf(os.Stderr, "Error: at least 1 daily backup")
		os.Exit(1)
	}

	intervals := []interval{
		{1, int(*fBackupDaily)},
	}
	if *fBackupWeekly > 0 {
		intervals = append(intervals, interval{7, int(*fBackupWeekly)})
	}
	if *fBackupMonthly > 0 {
		intervals = append(intervals, interval{30, int(*fBackupMonthly)})
	}

	scanner := bufio.NewScanner(os.Stdin)
	var lines []string
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v", err)
		os.Exit(1)
	}

	if *fDescending {
		reverseStrings(lines)
	}

	t := time.Now()
	ret := getResult(intervals, lines, t, !*fKeep)

	if *fDescending {
		reverseStrings(ret)
	}

	for _, line := range ret {
		fmt.Println(line)
	}
}

type interval struct {
	days  int
	count int
}

// makeFilterKeepMap creates a map, key is the date, value is the group id.
// For the same group id, only one date should be kept.
func makeFilterKeepMap(intervals []interval, t time.Time) map[string]int {
	ret := make(map[string]int)
	intvIndex := 0
	intvCount := 0
	intvDays := 0

	group := 1
	for {
		if intvIndex >= len(intervals) {
			break
		}

		intv := intervals[intvIndex]
		ret[t.Format(layoutDate)] = group
		if debug {
			fmt.Fprintf(os.Stderr, "filterMap: %s %d\n", t.Format(layoutDate), group)
		}
		if intvCount < intv.count && intvDays%intv.days == 0 {
			intvCount++
			group++
		}

		t = t.AddDate(0, 0, -1)
		intvDays++
		if intvCount >= intv.count || intvDays > intv.days*intv.count {
			intvCount = 0
			intvIndex++
			intvDays = 1
		}
	}
	return ret
}

func getResult(intervals []interval, lines []string, t time.Time, del bool) []string {
	times := make([]time.Time, 0, len(lines))
	for _, line := range lines {
		t, err := time.Parse(*fLayout, line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v", err)
			os.Exit(1)
		}
		times = append(times, t)
	}

	var ret []string
	filterKeepMap := makeFilterKeepMap(intervals, t)

	for i := 0; i < len(times); i++ {
		t := times[i]
		date := t.Format(layoutDate)
		group := filterKeepMap[date]

		// find all items with the same date, and maybe keep them
		sameDates := []string{}
		for j := i + 1; j < len(times); j++ {
			t := times[j]
			ndate := t.Format(layoutDate)
			if ndate != date {
				break
			}
			sameDates = append(sameDates, lines[j])
		}

		// find all items of the current group
		groupItems := []string{}
		for j := i + len(sameDates) + 1; j < len(times); j++ {
			t := times[j]
			date := t.Format(layoutDate)
			if ng := filterKeepMap[date]; group == ng {
				groupItems = append(groupItems, lines[j])
			} else {
				break
			}
		}

		// e.g. 18, 18, 19, 19, 20 in the same group
		// line[i] = 1st 18
		// sameDates = [2nd 18]
		// groupItems = [19, 19, 20]
		// keep the oldest, and the same dates
		// remove all the groupItems
		if del {
			if len(groupItems) > 0 {
				ret = append(ret, groupItems...)
			}
		} else {
			ret = append(ret, lines[i])
			if len(sameDates) > 0 {
				ret = append(ret, sameDates...)
			}
		}

		i += len(groupItems) + len(sameDates)
	}

	return ret
}

func reverseStrings(s []string) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}
