package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

func failWithError(c *gin.Context, err error) {
	fmt.Printf("%s", err.Error())
	errorMessage := "error"
	responseStatus := http.StatusInternalServerError
	if err.Error() == "Not Found" {
		errorMessage = "unknown"
		responseStatus = http.StatusNotFound
	}
	c.JSON(responseStatus, schemaResponse{
		SchemaVersion: 1,
		Label:         "openHPI",
		Message:       errorMessage,
		IsError:       true,
	})
}

func parseLanguageAndCourseParameters(bundle *i18n.Bundle, c *gin.Context) (language.Tag, *i18n.Localizer, string) {
	lang := c.Params.ByName("lang")
	if lang == "de" {
		return language.German, i18n.NewLocalizer(bundle, language.German.String()), c.Params.ByName("course")
	}
	return language.English, i18n.NewLocalizer(bundle, language.English.String()), c.Params.ByName("course")
}

func formatTime(t time.Time, lang language.Tag) string {
	monthsDE := map[time.Month]string{
		time.January:   "Januar",
		time.February:  "Februar",
		time.March:     "März",
		time.April:     "April",
		time.May:       "Mai",
		time.June:      "Juni",
		time.July:      "Juli",
		time.August:    "August",
		time.September: "September",
		time.October:   "Oktober",
		time.November:  "November",
		time.December:  "Dezember",
	}
	if lang == language.German {
		return fmt.Sprintf("%d. %s %d", t.Day(), monthsDE[t.Month()], t.Year())
	}
	return fmt.Sprintf("%s %d, %d", t.Month().String(), t.Day(), t.Year())
}