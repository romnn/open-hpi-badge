package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"golang.org/x/text/language"
)

type openHPIColor string

const (
	openHPIColorRed    openHPIColor = "b1063a"
	openHPIColorOrange openHPIColor = "e2681d"
	openHPIColorYellow openHPIColor = "f9a61a"
)

// schemaResponse ...
type schemaResponse struct {
	SchemaVersion int    `json:"schemaVersion,omitempty"`
	Label         string `json:"label,omitempty"`
	Message       string `json:"message,omitempty"`
	Color         string `json:"color,omitempty"`
	LabelColor    string `json:"labelColor,omitempty"`
	IsError       bool   `json:"isError,omitempty"`
	NamedLogo     string `json:"namedLogo,omitempty"`
	LogoSvg       string `json:"logoSvg,omitempty"`
	LogoColor     string `json:"logoColor,omitempty"`
	LogoWidth     string `json:"logoWidth,omitempty"`
	LogoPosition  string `json:"logoPosition,omitempty"`
	Style         string `json:"style,omitempty"`
	CacheSeconds  int    `json:"cacheSeconds,omitempty"`
}

func setupRouter(bundle *i18n.Bundle) *gin.Engine {
	r := gin.Default()

	// Health test
	r.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "up")
	})

	// https://img.shields.io/endpoint?url=...&style=...

	// Course Participants
	r.GET("/openhpi/:lang/course/:course/enrolled", courseEnrolledHandler(bundle))
	r.GET("/openhpi/:lang/course/:course/enrolled/:time", courseEnrolledHandler(bundle))

	// Course end date
	r.GET("/openhpi/:lang/course/:course/date/end", courseEndDateHandler(bundle))

	// Course start date
	r.GET("/openhpi/:lang/course/:course/date/start", courseStartDateHandler(bundle))

	// Course start date countdown
	r.GET("/openhpi/:lang/course/:course/countdown", courseCountdownHandler(bundle))

	// Course status
	r.GET("/openhpi/:lang/course/:course/status", courseStatusHandler(bundle))

	// Course duration
	r.GET("/openhpi/:lang/course/:course/duration", courseDurationHandler(bundle))

	// Course title
	r.GET("/openhpi/:lang/course/:course/title", courseTitleHandler(bundle))

	// Course language
	r.GET("/openhpi/:lang/course/:course/language", courseLanguageHandler(bundle))

	return r
}

func main() {
	app := &cli.App{
		Name:  "openhpibadge",
		Usage: "Serve an api endpoint for custom shields.io badges",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "production",
				Value:   false,
				Aliases: []string{"prod"},
				EnvVars: []string{"PRODUCTION", "RELEASE"},
				Usage:   "enable production mode",
			},
			&cli.IntFlag{
				Name:    "port",
				Value:   80,
				EnvVars: []string{"PORT"},
				Usage:   "service port",
			},
		},
		Action: func(c *cli.Context) error {
			bundle := i18n.NewBundle(language.English)
			bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
			data, err := Asset("intn/active.de.toml")
			if err != nil {
				return err
			}
			bundle.ParseMessageFileBytes(data, "active.de.toml")

			if c.Bool("production") {
				gin.SetMode(gin.ReleaseMode)
			}

			r := setupRouter(bundle)
			r.Run(fmt.Sprintf(":%d", c.Int("port")))
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
