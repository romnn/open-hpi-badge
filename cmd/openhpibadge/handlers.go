package main

import (
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/romnnn/openhpibadge"
)

func courseEnrolledHandler(bundle *i18n.Bundle) func (c *gin.Context) {
	return func (c *gin.Context) {
		_, loc, course := parseLanguageAndCourseParameters(bundle, c)
		mooc, err := openhpibadge.ScrapeMOOCByName(course)
		if err != nil {
			failWithError(c, err)
			return
		}

		var participants int
		switch strings.ToLower(c.Params.ByName("time")) {
		case "start":
			participants = mooc.Participants.Start
		case "end":
			participants = mooc.Participants.End
		default:
			participants = mooc.Participants.Current
		}

		translation, err := loc.Localize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "enrolled",
				One:   "{{.Count}} enrolled",
				Other: "{{.Count}} enrolled",
			},
			TemplateData: map[string]interface{}{
				"Count": participants,
			},
		})
		if err != nil {
			failWithError(c, err)
			return
		}
		c.JSON(http.StatusOK, schemaResponse{
			SchemaVersion: 1,
			Label:         "openHPI",
			Message:       translation,
		})
	}
}

func courseEndDateHandler(bundle *i18n.Bundle) func (c *gin.Context) {
	return func (c *gin.Context) {
		lang, loc, course := parseLanguageAndCourseParameters(bundle, c)
		mooc, err := openhpibadge.ScrapeMOOCByName(course)
		if err != nil {
			failWithError(c, err)
			return
		}
		end := formatTime(mooc.End, lang)
		ended := mooc.End.Before(time.Now())

		translation, err := loc.Localize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "end-date",
				One:   "{{if .Ended}}ended{{else}}ends{{end}} {{.Date}}",
				Other: "{{if .Ended}}ended{{else}}ends{{end}} {{.Date}}",
			},
			TemplateData: map[string]interface{}{
				"Date":  end,
				"Ended": ended,
			},
		})
		if err != nil {
			failWithError(c, err)
			return
		}
		c.JSON(http.StatusOK, schemaResponse{
			SchemaVersion: 1,
			Label:         "openHPI",
			Message:       translation,
		})
	}
}

func courseStartDateHandler(bundle *i18n.Bundle) func (c *gin.Context) {
	return func (c *gin.Context) {
		lang, loc, course := parseLanguageAndCourseParameters(bundle, c)
		mooc, err := openhpibadge.ScrapeMOOCByName(course)
		if err != nil {
			failWithError(c, err)
			return
		}
		start := formatTime(mooc.Start, lang)
		ended := mooc.Start.Before(time.Now())

		translation, err := loc.Localize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "start-date",
				One:   "{{if .Ended}}started{{else}}starts{{end}} {{.Date}}",
				Other: "{{if .Ended}}started{{else}}starts{{end}} {{.Date}}",
			},
			TemplateData: map[string]interface{}{
				"Date":  start,
				"Ended": ended,
			},
		})
		if err != nil {
			failWithError(c, err)
			return
		}
		c.JSON(http.StatusOK, schemaResponse{
			SchemaVersion: 1,
			Label:         "openHPI",
			Message:       translation,
		})
	}
}

func courseCountdownHandler(bundle *i18n.Bundle) func (c *gin.Context) {
	return func (c *gin.Context) {
		lang, loc, course := parseLanguageAndCourseParameters(bundle, c)
		mooc, err := openhpibadge.ScrapeMOOCByName(course)
		if err != nil {
			failWithError(c, err)
			return
		}
		end := formatTime(mooc.End, lang)
		delta := mooc.Start.Sub(time.Now())
		daysLeft := int(math.Floor(delta.Hours() / 24))
		daysRunning := int(math.Abs(float64(daysLeft)))
		started := daysLeft <= 0
		ended := mooc.End.Before(time.Now())
		running := started && !ended

		pluralCount := daysRunning
		if !ended && !running {
			pluralCount = daysLeft
		}

		translation, err := loc.Localize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "countdown",
				One:   "{{if .Ended}}ended {{.EndDate}}{{else if .Running}}running for {{.DaysRunning}} day{{else}}starts in {{.DaysLeft}} day{{end}}",
				Other: "{{if .Ended}}ended {{.EndDate}}{{else if .Running}}running for {{.DaysRunning}} days{{else}}starts in {{.DaysLeft}} days{{end}}",
			},
			TemplateData: map[string]interface{}{
				"EndDate":     end,
				"Running":     running,
				"DaysRunning": daysRunning,
				"DaysLeft":    daysLeft,
				"Ended":       ended,
			},
			PluralCount: pluralCount,
		})
		if err != nil {
			failWithError(c, err)
			return
		}
		c.JSON(http.StatusOK, schemaResponse{
			SchemaVersion: 1,
			Label:         "openHPI",
			Message:       translation,
		})
	}
}

func courseStatusHandler(bundle *i18n.Bundle) func (c *gin.Context) {
	return func (c *gin.Context) {
		_, loc, course := parseLanguageAndCourseParameters(bundle, c)
		mooc, err := openhpibadge.ScrapeMOOCByName(course)
		if err != nil {
			failWithError(c, err)
			return
		}
		started := mooc.Start.Before(time.Now())
		ended := mooc.End.Before(time.Now())
		running := started && !ended

		translation, err := loc.Localize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "status",
				One:   "{{if .Ended}}ended{{else if .Running}}running{{else}}coming up{{end}}",
				Other: "{{if .Ended}}ended{{else if .Running}}running{{else}}coming up{{end}}",
			},
			TemplateData: map[string]interface{}{
				"Running": running,
				"Ended":   ended,
			},
		})
		if err != nil {
			failWithError(c, err)
			return
		}
		c.JSON(http.StatusOK, schemaResponse{
			SchemaVersion: 1,
			Label:         "status",
			Message:       translation,
		})
	}
}

func courseDurationHandler(bundle *i18n.Bundle) func (c *gin.Context) {
	return func (c *gin.Context) {
		_, loc, course := parseLanguageAndCourseParameters(bundle, c)
		mooc, err := openhpibadge.ScrapeMOOCByName(course)
		if err != nil {
			failWithError(c, err)
			return
		}
		delta := mooc.End.Sub(mooc.Start)
		deltaWeeks := int(math.Ceil(delta.Hours() / 24 / 7))
		translation, err := loc.Localize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "duration",
				One:   "{{.Weeks}} week",
				Other: "{{.Weeks}} weeks",
			},
			TemplateData: map[string]interface{}{
				"Weeks": deltaWeeks,
			},
			PluralCount: deltaWeeks,
		})
		if err != nil {
			failWithError(c, err)
			return
		}
		c.JSON(http.StatusOK, schemaResponse{
			SchemaVersion: 1,
			Label:         "openHPI",
			Message:       translation,
		})
	}
}

func courseTitleHandler(bundle *i18n.Bundle) func (c *gin.Context) {
	return func (c *gin.Context) {
		_, loc, course := parseLanguageAndCourseParameters(bundle, c)
		mooc, err := openhpibadge.ScrapeMOOCByName(course)
		if err != nil {
			failWithError(c, err)
			return
		}
		translation, err := loc.Localize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "title",
				One:   "course",
				Other: "course",
			},
		})
		if err != nil {
			failWithError(c, err)
			return
		}
		c.JSON(http.StatusOK, schemaResponse{
			SchemaVersion: 1,
			Label:         translation,
			Message:       mooc.Title,
		})
	}
}

func courseLanguageHandler(bundle *i18n.Bundle) func (c *gin.Context) {
	return func (c *gin.Context) {
		_, loc, course := parseLanguageAndCourseParameters(bundle, c)
		mooc, err := openhpibadge.ScrapeMOOCByName(course)
		if err != nil {
			failWithError(c, err)
			return
		}
		translation, err := loc.Localize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "language",
				One:   "language",
				Other: "language",
			},
		})
		if err != nil {
			failWithError(c, err)
			return
		}
		c.JSON(http.StatusOK, schemaResponse{
			SchemaVersion: 1,
			Label:         translation,
			Message:       mooc.Language,
		})
	}
}