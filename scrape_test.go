package openhpibadge

import (
	"fmt"
	"time"

	"testing"
)

// TestScrapeMOOCByName ...
func TestHasClasses(t *testing.T) {
	if !hasClasses([]string{}, []string{}) {
		t.Error("Should be")
	}
	if !hasClasses([]string{"test"}, []string{}) {
		t.Error("Should be")
	}
	if hasClasses([]string{}, []string{"test"}) {
		t.Error("Should be")
	}
	if hasClasses([]string{"a"}, []string{"b"}) {
		t.Error("Should be")
	}
	if !hasClasses([]string{"a", "b", "c"}, []string{"a", "c", "b"}) {
		t.Error("Should be")
	}
}

// TestScrapeMOOCByName ...
func TestScrapeMOOCByName(t *testing.T) {
	name := "neuralnets2020"
	var err error
	var urlMooc, nameMooc *MOOC
	urlMooc, err = ScrapeMOOCByURL(fmt.Sprintf("https://open.hpi.de/courses/%s", name))
	nameMooc, err = ScrapeMOOCByName(name)
	if err != nil {
		t.Fatal("Failed to scrape by URL or name")
	}
	if _, equal := urlMooc.Equal(nameMooc); !equal {
		t.Fatal("ScrapeMOOCByName and ScrapeMOOCByURL get different results")
	}
}

// TestScrapeNeuralnets2020 ...
func TestScrapeNeuralnets2020(t *testing.T) {
	URL := "https://open.hpi.de/courses/neuralnets2020"
	course, err := ScrapeMOOCByURL(URL)
	if err != nil {
		t.Fatal(err)
	}
	if course.URL != URL {
		t.Errorf("Expected URL %s for neuralnets2020 but got %s", URL, course.URL)
	}
	title := "Praktische Einführung in Deep Learning für Computer Vision"
	if course.Title != title {
		t.Errorf("Expected title %s for neuralnets2020 but got %s", title, course.Title)
	}
	language := "Deutsch"
	if course.Language != language {
		t.Errorf("Expected language %s for neuralnets2020 but got %s", language, course.Language)
	}
	start := time.Date(2020, 3, 11, 0, 0, 0, 0, time.UTC)
	if course.Start != start {
		t.Errorf("Expected start date %v for neuralnets2020 but got %v", start, course.Start)
	}
	end := time.Date(2020, 4, 14, 0, 0, 0, 0, time.UTC)
	if course.End != end {
		t.Errorf("Expected end date %v for neuralnets2020 but got %v", end, course.End)
	}
	partsStart := 6334
	partsEnd := 8707
	minPartsCurrent := 9496
	if course.Participants.Start != partsStart {
		t.Errorf("Expected %d participants at the start for neuralnets2020 but got %d", partsStart, course.Participants.Start)
	}
	if course.Participants.End != partsEnd {
		t.Errorf("Expected %d participants at the start for neuralnets2020 but got %d", partsEnd, course.Participants.End)
	}
	if course.Participants.Current < minPartsCurrent {
		t.Errorf("Expected at least %d participants at the start for neuralnets2020 but got %d", minPartsCurrent, course.Participants.Current)
	}
}

// TestScrapePythonjunior2020 ...
func TestScrapePythonjunior2020(t *testing.T) {
	URL := "https://open.hpi.de/courses/pythonjunior2020"
	course, err := ScrapeMOOCByURL(URL)
	if err != nil {
		t.Fatal(err)
	}
	if course.URL != URL {
		t.Errorf("Expected URL %s for pythonjunior2020 but got %s", URL, course.URL)
	}
	title := "Programmieren lernen mit Python"
	if course.Title != title {
		t.Errorf("Expected title %s for pythonjunior2020 but got %s", title, course.Title)
	}
	language := "Deutsch"
	if course.Language != language {
		t.Errorf("Expected language %s for pythonjunior2020 but got %s", language, course.Language)
	}
	start := time.Date(2020, 10, 27, 0, 0, 0, 0, time.UTC)
	if course.Start != start {
		t.Errorf("Expected start date %v for pythonjunior2020 but got %v", start, course.Start)
	}
	end := time.Date(2020, 11, 23, 0, 0, 0, 0, time.UTC)
	if course.End != end {
		t.Errorf("Expected end date %v for pythonjunior2020 but got %v", end, course.End)
	}
	partsStart := 0
	partsEnd := 0
	minPartsCurrent := 500
	if course.Participants.Start != partsStart {
		t.Errorf("Expected %d participants at the start for pythonjunior2020 but got %d", partsStart, course.Participants.Start)
	}
	if course.Participants.End != partsEnd {
		t.Errorf("Expected %d participants at the start for pythonjunior2020 but got %d", partsEnd, course.Participants.End)
	}
	if course.Participants.Current < minPartsCurrent {
		t.Errorf("Expected at least %d participants at the start for pythonjunior2020 but got %d", minPartsCurrent, course.Participants.Current)
	}
}