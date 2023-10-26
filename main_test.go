package main

import (
	"testing"
	"time"
)

func makeItemAgo(t time.Time, d, group int) item {
	return item{t.AddDate(0, 0, -d).Format(layoutDate), group}
}

type item struct {
	key   string
	group int
}

func TestMakeFilterKeepMap(t *testing.T) {
	intervals := []interval{
		{1, 5},
		{3, 3},
		{7, 3},
	}

	ret := makeFilterKeepMap(intervals, time.Now())

	d := time.Now()
	keys := []item{
		makeItemAgo(d, 0, 1),
		makeItemAgo(d, 1, 2),
		makeItemAgo(d, 2, 3),
		makeItemAgo(d, 3, 4),
		makeItemAgo(d, 4, 5),
		makeItemAgo(d, 5, 6),
		makeItemAgo(d, 6, 6),
		makeItemAgo(d, 7, 6),
		makeItemAgo(d, 8, 7),
		makeItemAgo(d, 9, 7),
		makeItemAgo(d, 10, 7),
		makeItemAgo(d, 11, 8),
		makeItemAgo(d, 12, 8),
		makeItemAgo(d, 13, 8),
		makeItemAgo(d, 14, 9),
		makeItemAgo(d, 20, 9),
		makeItemAgo(d, 21, 10),
		makeItemAgo(d, 27, 10),
		makeItemAgo(d, 28, 11),
		makeItemAgo(d, 34, 11),
	}

	for i, ck := range keys {
		group := ret[ck.key]
		if group != ck.group {
			t.Fatalf("invalid group %d %s %d %d", i, ck.key, group, ck.group)
		}
	}
}

func makeLineAgo(t time.Time, d int) string {
	return t.AddDate(0, 0, -d).Format(layoutDateTime)
}

func TestGetResultDelete(t *testing.T) {
	intervals := []interval{
		{1, 5},
		{3, 3},
		{7, 3},
	}

	start, err := time.Parse(layoutDate, "20231026")
	if err != nil {
		t.Fatalf("time: %v", err)
	}

	lines := []string{
		makeLineAgo(start, 30), // 0926 g 11
		makeLineAgo(start, 27), // 0929 g 10
		makeLineAgo(start, 24), // 1002 g 10 remove
		makeLineAgo(start, 21), // 1005 g 10 remove
		makeLineAgo(start, 18), // 1008 g 9
		makeLineAgo(start, 17), // 1009 g 9 remove
		makeLineAgo(start, 10), // 1016 g 7
		makeLineAgo(start, 7),  // 1019 g 6
		makeLineAgo(start, 6),  // 1020 g 6 remove
		makeLineAgo(start, 4),  // 1022 g 1
		makeLineAgo(start, 0),  // 1026 g 1
	}
	// 0-4, 5-7: 6, 8-10: 7, 11-13: 8, 14-20: 9, 21-27: 10

	deleted := getResult(intervals, lines, start, true)

	if len(deleted) != 4 {
		t.Fatalf("deleted %v lines %v", deleted, lines)
	}

	if deleted[0] != lines[2] {
		t.Fatalf("deleted 0: %s %s", deleted[0], lines[2])
	}
	if deleted[1] != lines[3] {
		t.Fatalf("deleted 1: %s %s", deleted[1], lines[3])
	}
	if deleted[2] != lines[5] {
		t.Fatalf("deleted 2: %s %s", deleted[2], lines[5])
	}
	if deleted[3] != lines[8] {
		t.Fatalf("deleted 3: %s %s", deleted[3], lines[8])
	}
}

func TestGetResultKeep(t *testing.T) {
	intervals := []interval{
		{1, 5},
		{3, 3},
		{7, 3},
	}

	start, err := time.Parse(layoutDate, "20231026")
	if err != nil {
		t.Fatalf("time: %v", err)
	}

	lines := []string{
		makeLineAgo(start, 30), // 0926 g 11
		makeLineAgo(start, 27), // 0929 g 10
		makeLineAgo(start, 24), // 1002 g 10 remove
		makeLineAgo(start, 21), // 1005 g 10 remove
		makeLineAgo(start, 18), // 1008 g 9
		makeLineAgo(start, 17), // 1009 g 9 remove
		makeLineAgo(start, 10), // 1016 g 7
		makeLineAgo(start, 7),  // 1019 g 6
		makeLineAgo(start, 6),  // 1020 g 6 remove
		makeLineAgo(start, 4),  // 1022 g 1
		makeLineAgo(start, 0),  // 1026 g 1
	}
	// 0-4, 5-7: 6, 8-10: 7, 11-13: 8, 14-20: 9, 21-27: 10

	keep := getResult(intervals, lines, start, false)

	if len(keep) != 7 {
		t.Fatalf("keep %v lines %v", keep, lines)
	}

	if keep[0] != lines[0] {
		t.Fatalf("keep 0: %s %s", keep[0], lines[0])
	}
	if keep[1] != lines[1] {
		t.Fatalf("keep 1: %s %s", keep[1], lines[1])
	}
	if keep[2] != lines[4] {
		t.Fatalf("keep 2: %s %s", keep[2], lines[4])
	}
	if keep[3] != lines[6] {
		t.Fatalf("keep 3: %s %s", keep[3], lines[6])
	}
	if keep[4] != lines[7] {
		t.Fatalf("keep 4: %s %s", keep[4], lines[7])
	}
	if keep[5] != lines[9] {
		t.Fatalf("keep 5: %s %s", keep[5], lines[9])
	}
	if keep[6] != lines[10] {
		t.Fatalf("keep 6: %s %s", keep[6], lines[10])
	}
}

func TestGetResultKeepSameDate(t *testing.T) {
	intervals := []interval{
		{1, 3},
		{3, 3},
	}

	start, err := time.Parse(layoutDate, "20231026")
	if err != nil {
		t.Fatalf("time: %v", err)
	}

	lines := []string{
		makeLineAgo(start, 7), // g 5
		makeLineAgo(start, 6), // remove g 5
		makeLineAgo(start, 6), // remove g 5
		makeLineAgo(start, 5), // g 4
		makeLineAgo(start, 3), // remove g 4
		makeLineAgo(start, 2),
		makeLineAgo(start, 2),
		makeLineAgo(start, 1),
		makeLineAgo(start, 0),
		makeLineAgo(start, 0),
		makeLineAgo(start, 0),
	}
	// 0-2, [3 4 5]: 4, [6 7 8]: 5

	keep := getResult(intervals, lines, start, false)

	if len(keep) != 8 {
		t.Fatalf("keep %v lines %v", keep, lines)
	}

	if keep[0] != lines[0] {
		t.Fatalf("keep 0: %s %s", keep[0], lines[0])
	}
	if keep[1] != lines[3] {
		t.Fatalf("keep 1: %s %s", keep[1], lines[3])
	}

	for i := 2; i < 8; i++ {
		if keep[i] != lines[i+3] {
			t.Fatalf("keep %d: %s %s", i, keep[i], lines[i+3])
		}
	}
}

func TestGetResultDeleteSameDate(t *testing.T) {
	intervals := []interval{
		{1, 3},
		{3, 3},
	}

	start, err := time.Parse(layoutDate, "20231026")
	if err != nil {
		t.Fatalf("time: %v", err)
	}

	lines := []string{
		makeLineAgo(start, 7), // g 5
		makeLineAgo(start, 6), // remove g 5
		makeLineAgo(start, 6), // remove g 5
		makeLineAgo(start, 5), // g 4
		makeLineAgo(start, 3), // remove g 4
		makeLineAgo(start, 2),
		makeLineAgo(start, 2),
		makeLineAgo(start, 1),
		makeLineAgo(start, 0),
		makeLineAgo(start, 0),
		makeLineAgo(start, 0),
	}
	// 0-2, [3 4 5]: 4, [6 7 8]: 5

	deleted := getResult(intervals, lines, start, true)

	if len(deleted) != 3 {
		t.Fatalf("deleted %v lines %v", deleted, lines)
	}

	if deleted[0] != lines[1] {
		t.Fatalf("deleted 0: %s %s", deleted[0], lines[1])
	}
	if deleted[1] != lines[2] {
		t.Fatalf("deleted 1: %s %s", deleted[1], lines[2])
	}
	if deleted[2] != lines[4] {
		t.Fatalf("deleted 2: %s %s", deleted[2], lines[4])
	}
}

func TestWithText(t *testing.T) {
	lines := []string{
		"20230831_043501",
		"20230907_043501",
		"20230913_043501",
		"20230920_043501",
		"20230927_040501",
		"20231004_040501",
		"20231011_040501",
		"20231012_040501",
		"20231013_040501",
		"20231014_040501",
		"20231015_040501",
		"20231016_040501",
		"20231017_040501",
		"20231018_040501",
		"20231019_040501",
		"20231020_040501",
		"20231021_040501",
		"20231022_040501",
		"20231023_040501",
		"20231024_040501",
		"20231025_040501",
	}

	tm, err := time.Parse(layoutDate, "20231025")
	if err != nil {
		t.Fatalf("time: %v", err)
	}

	intervals := []interval{
		{1, 14},
		{7, 6},
		{30, 10},
	}

	deleted := getResult(intervals, lines, tm, true)
	if len(deleted) != 1 || deleted[0] != "20230913_043501" {
		t.Fatalf("deleted: %v", deleted)
	}
}
