package filter

import (
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/nalabelle/miniflux-sidekick/rules"
	miniflux "miniflux.app/client"
)

func TestEvaluateRules(t *testing.T) {
	type mockService struct {
		rules []rules.Rule
		l     log.Logger
	}

	tests := []struct {
		name       string
		expression string
		args       *miniflux.Entry
		want       bool
	}{
		{
			name:       "Entry contains string",
			expression: "\"ignore-article\" \"http://example.com/feed.xml\" \"title # Moon\"",
			args: &miniflux.Entry{
				Title: "Moon entry",
			},
			want: true,
		},
		{
			name:       "Entry contains string",
			expression: "\"ignore-article\" \"http://example.com/feed.xml\" \"title # Moon\"",
			args: &miniflux.Entry{
				Title: "Sun entry",
			},
			want: false,
		},
		{
			name:       "Entry contains string, matched with Regexp",
			expression: "\"ignore-article\" \"http://example.com/feed.xml\" \"title =~ [Sponsor]\"",
			args: &miniflux.Entry{
				Title: "[Sponsor] Sun entry",
			},
			want: true,
		},
		{
			name:       "Entry doesn't string, matched with Regexp",
			expression: "\"ignore-article\" \"http://example.com/feed.xml\" `title =~ \\[Sponsor\\]`",
			args: &miniflux.Entry{
				Title: "[SponSomethingElsesor] Sun entry",
			},
			want: false,
		},
		{
			name:       "Entry doesn't string, matched with Regexp, ignore case",
			expression: "\"ignore-article\" \"http://example.com/feed.xml\" \"title =~ (?i)(Podcast|scooter)\"",
			args: &miniflux.Entry{
				Title: "podcast",
			},
			want: true,
		},
		{
			name:       "Entry doesn't string, matched with Regexp, ignore case",
			expression: "\"ignore-article\" \"http://example.com/feed.xml\" \"title =~ (?i)(Podcast|scooter)\"",
			args: &miniflux.Entry{
				Title: "SCOOTER",
			},
			want: true,
		},
		{
			name:       "Entry doesn't string, matched with Regexp, respect case",
			expression: "\"ignore-article\" \"http://example.com/feed.xml\" \"title =~ (Podcast)\"",
			args: &miniflux.Entry{
				Title: "podcast",
			},
			want: false,
		},
		{
			name:       "Entry doesn't tag, case insensitive",
			expression: "\"ignore-article\" \"http://example.com/feed.xml\" \"tag # (?i)podcast\"",
			args: &miniflux.Entry{
				Tags: []string{"this", "Podcast", "testing"},
			},
			want: true,
		},
		{
			name:       "Entry doesn't tag, respect case",
			expression: "\"ignore-article\" \"http://example.com/feed.xml\" \"tag # podcast\"",
			args: &miniflux.Entry{
				Tags: []string{"this", "Podcast", "testing"},
			},
			want: false,
		},
		{
			name:       "Entry doesn't tag, miss",
			expression: "\"ignore-article\" \"http://example.com/feed.xml\" \"tag # test\"",
			args: &miniflux.Entry{
				Tags: []string{"this", "podcast", "testing"},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepository := rules.NewRepository()
			rule, _ := rules.Parse(tt.expression)
			mockRepository.SetRules([]rules.Rule{rule})

			s := service{
				rulesRepository: mockRepository,
			}

			if got := s.evaluateRules(tt.args); got != tt.want {
				t.Errorf("evaluateRules() = %v, want %v", got, tt.want)
			}
		})
	}
}
