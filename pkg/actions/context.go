package actions

import (
	"os"
	"strconv"
	"strings"
)

type Context struct {
	EventName  string
	SHA        string
	Ref        GitRef
	Workflow   string
	Action     string
	Actor      string
	Repository string
	Job        string
	RunNumber  int64
	RunID      int64
}

type GitRefType int

const (
	Unknown GitRefType = iota
	Branch
	PullRequest
	Tag
)

type GitRef struct {
	Type GitRefType
	Name string
}

func NewContext() Context {
	c := Context{
		EventName:  os.Getenv("GITHUB_EVENT_NAME"),
		SHA:        os.Getenv("GITHUB_SHA"),
		Ref:        parseRef(os.Getenv("GITHUB_REF")),
		Workflow:   os.Getenv("GITHUB_WORKFLOW"),
		Action:     os.Getenv("GITHUB_ACTION"),
		Actor:      os.Getenv("GITHUB_ACTOR"),
		Repository: os.Getenv("GITHUB_REPOSITORY"),
		Job:        os.Getenv("GITHUB_JOB"),
	}

	c.RunNumber, _ = strconv.ParseInt(os.Getenv("GITHUB_RUN_NUMBER"), 10, 64)
	c.RunID, _ = strconv.ParseInt(os.Getenv("GITHUB_RUN_ID"), 10, 64)

	return c
}

func parseRef(refStr string) GitRef {
	ref := strings.SplitN(refStr, "/", 3)

	if len(ref) != 3 || ref[0] != "refs" {
		return GitRef{Unknown, refStr}
	}

	var r GitRef

	switch ref[1] {
	case "heads":
		r = GitRef{Branch, ref[2]}
	case "pull":
		r = GitRef{PullRequest, ref[2]}
	case "tags":
		r = GitRef{Tag, ref[2]}
	default:
		r = GitRef{Unknown, refStr}
	}

	return r
}
