package openhpibadge

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/araddon/dateparse"
	"github.com/gocolly/colly/v2"
)

func hasClasses(have []string, want []string) bool {
	found := make([]bool, len(want))
	for i, f := range want {
		for _, c := range have {
			if c == f {
				found[i] = true
			}
		}
	}
	for _, f := range found {
		if !f {
			return false
		}
	}
	return true
}

func parseDate(s string) (t time.Time, err error) {
	raw := strings.TrimSpace(s)
	var layout string
	layout, err = dateparse.ParseFormat(raw)
	if err != nil {
		return
	}
	t, err = time.Parse(layout, raw)
	if err != nil {
		return
	}
	return
}

func parseDateRange(parts []string) (start time.Time, end time.Time, err error) {
	if len(parts) != 2 {
		err = errors.New("Need exactly two daterange components")
		return
	}
	start, err = parseDate(parts[0])
	end, err = parseDate(parts[1])
	return
}

// ScrapeMOOCByName ...
func ScrapeMOOCByName(name string) (*MOOC, error) {
	URL := fmt.Sprintf("https://open.hpi.de/courses/%s", name)
	return ScrapeMOOCByURL(URL)
}

// ScrapeMOOCByURL ...
func ScrapeMOOCByURL(URL string) (*MOOC, error) {
	mooc := MOOC{URL: URL}
	var scrapeErr error
	var scrapeStatusCode int

	// Validate URL
	parsedURL, err := url.ParseRequestURI(URL)
	if err != nil {
		return nil, err
	}
	if parsedURL.Host != "open.hpi.de" {
		return nil, fmt.Errorf("%s is not a valid URL for open.hpi.de", URL)
	}

	c := colly.NewCollector(
		colly.AllowedDomains("open.hpi.de"),
		colly.Async(true),
	)

	c.Limit(&colly.LimitRule{
		DomainGlob:  "open.hpi.de/*",
		Delay:       1 * time.Second,
		RandomDelay: 1 * time.Second,
	})

	// Set error handler
	c.OnError(func(r *colly.Response, err error) {
		scrapeErr = err
		scrapeStatusCode = r.StatusCode
	})

	c.OnHTML("div#schema-course-title", func(e *colly.HTMLElement) {
		mooc.Title = e.Text
	})

	c.OnHTML("div.enrollment-statistics", func(e *colly.HTMLElement) {
		e.ForEach("div[data-type]", func(_ int, elem *colly.HTMLElement) {
			rawCount := elem.ChildText("div.enrollment-statistics__count")
			participants, partErr := strconv.Atoi(strings.ReplaceAll(rawCount, ",", ""))
			date, dateErr := parseDate(elem.ChildText("div.enrollment-statistics__date-value"))
			switch elem.Attr("data-type") {
			case "current":
				if partErr == nil {
					mooc.Participants.Current = participants
				}
			case "course_start":
				if dateErr == nil {
					mooc.Start = date
				}
				if partErr == nil {
					mooc.Participants.Start = participants
				}
			case "course_end":
				if dateErr == nil {
					mooc.End = date
				}
				if partErr == nil {
					mooc.Participants.End = participants
				}
			}
		})
	})

	c.OnHTML("div.course-enrollment-count", func(e *colly.HTMLElement) {
		rawCount := e.ChildText("span.badge")
		participants, err := strconv.Atoi(strings.ReplaceAll(rawCount, ",", ""))
		if err == nil {
			mooc.Participants.Current = participants
		}
	})

	c.OnHTML("span.shortinfo", func(e *colly.HTMLElement) {
		infoClasses := strings.Split(e.ChildAttr("span:nth-child(1)", "class"), " ")
		isLanguage := hasClasses(infoClasses, []string{"icon-language"})
		isCalendar := hasClasses(infoClasses, []string{"icon-calendar"})

		if isLanguage {
			languageRaw := e.ChildText("span:nth-child(2)")
			if re, err := regexp.Compile(`Language: (\w+)`); err == nil {
				result := re.FindStringSubmatch(languageRaw)
				if len(result) == 2 {
					mooc.Language = result[1]
				}
			}
		} else if isCalendar {
			if start, end, err := parseDateRange(strings.Split(e.ChildText("span:nth-child(2)"), "-")); err == nil {
				mooc.Start = start
				mooc.End = end
			}
		}
	})

	c.Visit(URL)
	c.Wait()

	if scrapeErr != nil {
		return nil, scrapeErr
	}

	if scrapeStatusCode > 0 && scrapeStatusCode != 200 {
		return nil, fmt.Errorf("Scraping %s returned bad status code %d", URL, scrapeStatusCode)
	}
	return &mooc, nil
}
