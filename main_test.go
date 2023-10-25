package main

import (
	"testing"
	"time"
)

func makeItem(t time.Time, d, group int) item {
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

	ret := makeFilterKeepMap(intervals)

	d := time.Now()
	keys := []item{
		makeItem(d, 0, 1),
		makeItem(d, 1, 2),
		makeItem(d, 2, 3),
		makeItem(d, 3, 4),
		makeItem(d, 4, 5),
		makeItem(d, 5, 6),
		makeItem(d, 6, 6),
		makeItem(d, 7, 6),
		makeItem(d, 8, 7),
		makeItem(d, 9, 7),
		makeItem(d, 10, 7),
		makeItem(d, 11, 8),
		makeItem(d, 12, 8),
		makeItem(d, 13, 8),
		makeItem(d, 14, 9),
		makeItem(d, 20, 9),
		makeItem(d, 21, 10),
		makeItem(d, 27, 10),
		makeItem(d, 28, 11),
		makeItem(d, 34, 11),
	}

	for i, ck := range keys {
		group := ret[ck.key]
		if group != ck.group {
			t.Fatalf("invalid group %d %s %d %d", i, ck.key, group, ck.group)
		}
	}
}

func makeLine(t time.Time, d int) string {
	return t.AddDate(0, 0, -d).Format(layoutDateTime)
}

func TestGetResultDelete(t *testing.T) {
	intervals := []interval{
		{1, 5},
		{3, 3},
		{7, 3},
	}

	start := time.Now()
	lines := []string{
		makeLine(start, 0),
		makeLine(start, 4),
		makeLine(start, 6),  // remove, g 6
		makeLine(start, 7),  // g 6
		makeLine(start, 10), // g 7
		makeLine(start, 17), // remove g 9
		makeLine(start, 18), // g 9
		makeLine(start, 21), // remove g 10
		makeLine(start, 24), // remove g 10
		makeLine(start, 27), // g 10
		makeLine(start, 30), // g 11
	}
	// 0-4, 5-7: 6, 8-10: 7, 11-13: 8, 14-20: 9, 21-27: 10

	deleted := getResult(intervals, lines, true)

	if len(deleted) != 4 {
		t.Fatalf("deleted %v lines %v", deleted, lines)
	}

	if deleted[0] != lines[2] {
		t.Fatalf("deleted 0: %s %s", deleted[0], lines[2])
	}
	if deleted[1] != lines[5] {
		t.Fatalf("deleted 1: %s %s", deleted[1], lines[5])
	}
	if deleted[2] != lines[7] {
		t.Fatalf("deleted 2: %s %s", deleted[2], lines[7])
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

	start := time.Now()
	lines := []string{
		makeLine(start, 0),
		makeLine(start, 4),
		makeLine(start, 6),  // remove, g 6
		makeLine(start, 7),  // g 6
		makeLine(start, 10), // g 7
		makeLine(start, 17), // remove g 9
		makeLine(start, 18), // g 9
		makeLine(start, 21), // remove g 10
		makeLine(start, 24), // remove g 10
		makeLine(start, 27), // g 10
		makeLine(start, 30), // g 11
	}
	// 0-4, 5-7: 6, 8-10: 7, 11-13: 8, 14-20: 9, 21-27: 10

	keep := getResult(intervals, lines, false)

	if len(keep) != 7 {
		t.Fatalf("keep %v lines %v", keep, lines)
	}

	if keep[0] != lines[0] {
		t.Fatalf("keep 0: %s %s", keep[0], lines[0])
	}
	if keep[1] != lines[1] {
		t.Fatalf("keep 1: %s %s", keep[1], lines[1])
	}
	if keep[2] != lines[3] {
		t.Fatalf("keep 2: %s %s", keep[2], lines[3])
	}
	if keep[3] != lines[4] {
		t.Fatalf("keep 3: %s %s", keep[3], lines[4])
	}
	if keep[4] != lines[6] {
		t.Fatalf("keep 3: %s %s", keep[4], lines[6])
	}
	if keep[5] != lines[9] {
		t.Fatalf("keep 3: %s %s", keep[5], lines[9])
	}
	if keep[6] != lines[10] {
		t.Fatalf("keep 3: %s %s", keep[6], lines[10])
	}
}

func TestGetResultKeepSameDate(t *testing.T) {
	intervals := []interval{
		{1, 3},
		{3, 3},
	}

	start := time.Now()
	lines := []string{
		makeLine(start, 0),
		makeLine(start, 0),
		makeLine(start, 0),
		makeLine(start, 1),
		makeLine(start, 2),
		makeLine(start, 2),
		makeLine(start, 3), // remove g 4
		makeLine(start, 5), // g 4
		makeLine(start, 6), // remove g 5
		makeLine(start, 6), // remove g 5
		makeLine(start, 7), // g 5
	}
	// 0-2, [3 4 5]: 4, [6 7 8]: 5

	keep := getResult(intervals, lines, false)

	if len(keep) != 8 {
		t.Fatalf("keep %v lines %v", keep, lines)
	}

	for i := 0; i < 6; i++ {
		if keep[i] != lines[i] {
			t.Fatalf("keep %d: %s %s", i, keep[i], lines[i])
		}
	}
	if keep[6] != lines[7] {
		t.Fatalf("keep 7: %s %s", keep[6], lines[7])
	}
	if keep[7] != lines[10] {
		t.Fatalf("keep 10: %s %s", keep[7], lines[10])
	}
}

func TestGetResultDeleteSameDate(t *testing.T) {
	intervals := []interval{
		{1, 3},
		{3, 3},
	}

	start := time.Now()
	lines := []string{
		makeLine(start, 0),
		makeLine(start, 0),
		makeLine(start, 0),
		makeLine(start, 1),
		makeLine(start, 2),
		makeLine(start, 2),
		makeLine(start, 3), // remove g 4
		makeLine(start, 5), // g 4
		makeLine(start, 6), // remove g 5
		makeLine(start, 6), // remove g 5
		makeLine(start, 7), // g 5
	}
	// 0-2, [3 4 5]: 4, [6 7 8]: 5

	deleted := getResult(intervals, lines, true)

	if len(deleted) != 3 {
		t.Fatalf("deleted %v lines %v", deleted, lines)
	}

	if deleted[0] != lines[6] {
		t.Fatalf("deleted 7: %s %s", deleted[0], lines[6])
	}
	if deleted[1] != lines[8] {
		t.Fatalf("deleted 7: %s %s", deleted[1], lines[8])
	}
	if deleted[2] != lines[8] {
		t.Fatalf("deleted 10: %s %s", deleted[2], lines[9])
	}
}
