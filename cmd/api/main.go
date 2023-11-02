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
	"github.com/peterbourgon/ff"
	miniflux "miniflux.app/client"
)

func main() {
	fs := flag.NewFlagSet("mf", flag.ExitOnError)
	var (
		minifluxAPIKey      = fs.String("api-key", "", "api key used for authentication")
		minifluxAPIEndpoint = fs.String("api-endpoint", "https://rss.notmyhostna.me", "the api of your miniflux instance")
		killfilePath        = fs.String("killfile-path", "", "the path to the local killfile")
		logLevel            = fs.String("log-level", "", "the level to filter logs at eg. debug, info, warn, error")
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
	filterService.RunFilterJob(false)
}
