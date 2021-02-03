package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
)

// TestFormatTime ...
func TestFormatTime(t *testing.T) {
	testDate := time.Date(2019, 2, 5, 0, 0, 0, 0, time.UTC)
	cases := []struct {
		date     time.Time
		lang     language.Tag
		expected string
	}{
		{testDate, language.German, "5. Februar 2019"},
		{testDate, language.English, "February 5, 2019"},
	}
	for _, c := range cases {
		result := formatTime(c.date, c.lang)
		if result != c.expected {
			t.Errorf("formatTime(%v, %v) = %s but expected %s", c.date, c.lang, result, c.expected)
		}
	}
}

func GET(uri string) *http.Request {
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		panic(err)
	}
	return req
}

// TestRoutes ...
func TestRoutes(t *testing.T) {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	data, err := Asset("intn/active.de.toml")
	if err != nil {
		t.Fatal(err)
	}
	bundle.ParseMessageFileBytes(data, "active.de.toml")
	gin.SetMode(gin.ReleaseMode)
	router := setupRouter(bundle)

	type result struct {
		SchemaVersion int    `json:"schemaVersion"`
		Label         string `json:"label"`
		Message       string `json:"message"`
	}

	casesDE := []struct {
		req      *http.Request
		expected result
	}{
		{GET("/openhpi/de/course/neuralnets2020/language"), result{
			SchemaVersion: 1,
			Label:         "Sprache",
			Message:       "Deutsch",
		}},
		{GET("/openhpi/de/course/neuralnets2020/title"), result{
			SchemaVersion: 1,
			Label:         "Kurs",
			Message:       "Praktische Einführung in Deep Learning für Computer Vision",
		}},
		{GET("/openhpi/de/course/neuralnets2020/duration"), result{
			SchemaVersion: 1,
			Label:         "openHPI",
			Message:       "5 Wochen",
		}},
		{GET("/openhpi/de/course/neuralnets2020/status"), result{
			SchemaVersion: 1,
			Label:         "status",
			Message:       "abgeschlossen",
		}},
		{GET("/openhpi/de/course/neuralnets2020/countdown"), result{
			SchemaVersion: 1,
			Label:         "openHPI",
			Message:       "endete am 14. April 2020",
		}},
		{GET("/openhpi/de/course/neuralnets2020/date/end"), result{
			SchemaVersion: 1,
			Label:         "openHPI",
			Message:       "endete am 14. April 2020",
		}},
		{GET("/openhpi/de/course/neuralnets2020/date/start"), result{
			SchemaVersion: 1,
			Label:         "openHPI",
			Message:       "startete am 11. März 2020",
		}},
		{GET("/openhpi/de/course/neuralnets2020/enrolled/start"), result{
			SchemaVersion: 1,
			Label:         "openHPI",
			Message:       "6334 eingeschrieben",
		}},
		{GET("/openhpi/de/course/neuralnets2020/enrolled/end"), result{
			SchemaVersion: 1,
			Label:         "openHPI",
			Message:       "8707 eingeschrieben",
		}},
	}

	casesEN := []struct {
		req      *http.Request
		expected result
	}{
		{GET("/openhpi/en/course/neuralnets2020/language"), result{
			SchemaVersion: 1,
			Label:         "language",
			Message:       "Deutsch",
		}},
		{GET("/openhpi/en/course/neuralnets2020/title"), result{
			SchemaVersion: 1,
			Label:         "course",
			Message:       "Praktische Einführung in Deep Learning für Computer Vision",
		}},
		{GET("/openhpi/en/course/neuralnets2020/duration"), result{
			SchemaVersion: 1,
			Label:         "openHPI",
			Message:       "5 weeks",
		}},
		{GET("/openhpi/en/course/neuralnets2020/status"), result{
			SchemaVersion: 1,
			Label:         "status",
			Message:       "ended",
		}},
		{GET("/openhpi/en/course/neuralnets2020/countdown"), result{
			SchemaVersion: 1,
			Label:         "openHPI",
			Message:       "ended April 14, 2020",
		}},
		{GET("/openhpi/en/course/neuralnets2020/date/end"), result{
			SchemaVersion: 1,
			Label:         "openHPI",
			Message:       "ended April 14, 2020",
		}},
		{GET("/openhpi/en/course/neuralnets2020/date/start"), result{
			SchemaVersion: 1,
			Label:         "openHPI",
			Message:       "started March 11, 2020",
		}},
		{GET("/openhpi/en/course/neuralnets2020/enrolled/start"), result{
			SchemaVersion: 1,
			Label:         "openHPI",
			Message:       "6334 enrolled",
		}},
		{GET("/openhpi/en/course/neuralnets2020/enrolled/end"), result{
			SchemaVersion: 1,
			Label:         "openHPI",
			Message:       "8707 enrolled",
		}},
	}

	for _, c := range append(casesDE, casesEN...) {
		w := httptest.NewRecorder()
		router.ServeHTTP(w, c.req)
		assert.Equal(t, 200, w.Code)
		var out result
		if err := json.Unmarshal(w.Body.Bytes(), &out); err != nil {
			t.Error(err)
		}
		assert.Equal(t, c.expected, out)
	}
}
