package main

import (
	"errors"
	"flag"
	"os"
	"strings"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/nalabelle/miniflux-sidekick/filter"
	"github.com/nalabelle/miniflux-sidekick/rules"
	ff "github.com/peterbourgon/ff/v3"
	pkg_cron "github.com/robfig/cron/v3"
	miniflux "miniflux.app/v2/client"
)

func main() {
	fs := flag.NewFlagSet("mf", flag.ExitOnError)
	var (
		killfilePath        = fs.String("killfile-path", "", "the path to the local killfile")
		logLevel            = fs.String("log-level", "", "the level to filter logs at eg. debug, info, warn, error")
		minifluxAPIEndpoint = fs.String("api-endpoint", "https://rss.notmyhostna.me", "the api of your miniflux instance")
		minifluxAPIKey      = fs.String("api-key", "", "api key used for authentication")
		refreshInterval     = fs.String("refresh-interval", "", "interval defining how often we check for new entries in miniflux")
	)

	ff.Parse(fs, os.Args[1:],
		ff.WithConfigFileFlag("config"),
		ff.WithConfigFileParser(ff.PlainParser),
		ff.WithEnvVarPrefix("MF"),
	)

	l := log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	switch strings.ToLower(*logLevel) {
	case "debug":
		l = level.NewFilter(l, level.AllowDebug())
	case "info":
		l = level.NewFilter(l, level.AllowInfo())
	case "warn":
		l = level.NewFilter(l, level.AllowWarn())
	case "error":
		l = level.NewFilter(l, level.AllowError())
	default:
		l = level.NewFilter(l, level.AllowError())
	}
	l = log.With(l, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)

	var client *miniflux.Client
	if *minifluxAPIKey != "" {
		client = miniflux.New(*minifluxAPIEndpoint, *minifluxAPIKey)
	} else {
		level.Error(l).Log("err", errors.New("api endpoint and api key need to be provided"))
		return
	}
	u, err := client.Me()
	if err != nil {
		level.Error(l).Log("err", err)
		return
	}
	level.Info(l).Log("msg", "user successfully logged in", "username", u.Username, "user_id", u.ID, "is_admin", u.IsAdmin)

	var rr rules.Repository
	if *killfilePath != "" {
		level.Info(l).Log("msg", "using a local killfile", "path", *killfilePath)
		localRepo := rules.NewRepository()
		if err != nil {
			level.Error(l).Log("err", err)
			return
		}
		localRepo.FetchRules(*killfilePath, l)
		rr = localRepo
	}

	filterService := filter.NewService(l, client, rr)
	cron := pkg_cron.New()
	level.Info(l).Log("msg", "running filter job in destructive mode", "interval_cron", *refreshInterval)
	_, err = cron.AddJob(*refreshInterval, filterService)
	if err != nil {
		level.Error(l).Log("msg", "error adding cron job to scheduler", "err", err)
	}
	cron.Run()
}
