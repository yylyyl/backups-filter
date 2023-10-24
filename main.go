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

	if !*fDescending {
		reverseStrings(lines)
	}

	ret := getResult(intervals, lines, !*fKeep)

	if !*fDescending {
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
func makeFilterKeepMap(intervals []interval) map[string]int {
	ret := make(map[string]int)
	intvIndex := 0
	intvCount := 0
	intvDays := 0

	t := time.Now()
	group := 1
	for {
		if intvIndex >= len(intervals) {
			break
		}

		intv := intervals[intvIndex]
		ret[t.Format(layoutDate)] = group
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

func getResult(intervals []interval, lines []string, del bool) []string {
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
	filterKeepMap := makeFilterKeepMap(intervals)

	groupItems := []string{}
	for i := 0; i < len(times); i++ {
		t := times[i]
		date := t.Format(layoutDate)
		group := filterKeepMap[date]

		// find all items of the current group
		for j := i + 1; j < len(times); j++ {
			t := times[j]
			date := t.Format(layoutDate)
			if ng := filterKeepMap[date]; group == ng {
				groupItems = append(groupItems, lines[j])
			} else {
				break
			}
		}

		// keep the farthest item in the same group
		if len(groupItems) > 0 {
			if del {
				ret = append(ret, lines[i])
				if len(groupItems) > 1 {
					ret = append(ret, groupItems[0:len(groupItems)-1]...)
				}
			} else {
				ret = append(ret, groupItems[len(groupItems)-1])
			}
			i += len(groupItems)
			groupItems = []string{}
		} else {
			// the only item in the group
			if !del {
				ret = append(ret, lines[i])
			}
		}
	}

	return ret
}

func reverseStrings(s []string) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}
