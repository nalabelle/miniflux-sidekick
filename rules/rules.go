package rules

import (
	"bufio"
	"errors"
	"os"
	"regexp"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

var (
	reRuleSplitter = regexp.MustCompile(`(.+?)\s\"?(.+?)\"?\s\"(\w+?) (\S+?) (.+)\"`)
)

// Repository defines the interface for the rules repository
type Repository interface {
	// FetchRules fetches the list of rules from a file or remote location
	FetchRules(location string, l log.Logger)
	Rules() []Rule
	SetRules(rules []Rule)
}

type repository struct {
	rules []Rule
}

// NewLocalRepository returns a newly initialized rules repository
func NewRepository() Repository {
	return &repository{}
}

func (r *repository) Rules() []Rule {
	return r.rules
}

func (r *repository) SetRules(rules []Rule) {
	r.rules = rules
}

// FetchRules parses a local killfile to get all rules
func (r *repository) FetchRules(location string, l log.Logger) {
	file, err := os.Open(location)
	if err != nil {
		level.Error(l).Log("err", err)
		return
	}
	defer file.Close()

	var rules []Rule
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		rule, err := Parse(line)
		if err != nil {
			level.Error(l).Log("err", err, "expression", line)
			continue
		}

		rules = append(rules, rule)
	}
	r.rules = rules
}

// Rule contains a killfile rule. There's no official standard so we implement these rules https://newsboat.org/releases/2.15/docs/newsboat.html#_killfiles
type Rule struct {
	Command   string
	URL       string
	Attribute string
	Operator  string
	Match     string
}

func Parse(line string) (Rule, error) {
	matches := reRuleSplitter.FindStringSubmatch(line)
	if len(matches) == 6 {
		return Rule{
			Command:   matches[1],
			URL:       matches[2],
			Attribute: matches[3],
			Operator:  matches[4],
			Match:     matches[5],
		}, nil
	} else {
		return Rule{}, errors.New("invalid filter expression")
	}
}
