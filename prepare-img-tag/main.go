package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/snapADDY/actions/internal/actions"
	"github.com/snapADDY/actions/internal/msgflag"
)

// Set of branches which should always create and publish artifacts.
// Overriding the environment-name is also not allowed for commits
// to these branches to prevent issues with squash merges.
var alwaysPublishBranch = map[string]struct{}{
	"master":  {},
	"main":    {},
	"staging": {},
}

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {
	action := actions.NewContext()
	log.Printf("Actions Context: %#v\n", action)

	commitMsg := os.Getenv("INPUT_COMMITMESSAGE")
	flags := msgflag.Parse(commitMsg)

	log.Printf("Parsed commit message: %s\nFlags: %#v\n", commitMsg, flags)

	imgFlag := parseImgFlag(flags["img"])
	tag, envName, publish := makeImgTag(action, imgFlag)

	log.Printf("Prepared Tag: %s; Should publish %v\n", tag, publish)

	actions.SetOutput("tag-name", tag)
	actions.SetOutput("publish", publish)

	actions.SetOutput("has-img-flag", imgFlag.ShouldPublish)
	actions.SetOutput("env-name", envName)
	actions.SetOutput("env-override", imgFlag.EnvNameOverride)

	return nil
}

// makeImgTag creates the tag for our container registry and decides if the image should be published
func makeImgTag(action actions.Context, imgFlag ImgFlag) (tag, envName string, publish bool) {
	tagRegex := regexp.MustCompile(`^\w[\w.-]{0,127}$`)
	invalidChar := regexp.MustCompile(`[^\w.-]+`)

	if action.EventName == "push" {
		if action.Ref.Type == actions.Tag && tagRegex.MatchString(action.Ref.Name) {
			// Tags should always be published.
			return action.Ref.Name, action.Ref.Name, true
		} else if action.Ref.Type == actions.Branch {
			// Clean the branch name to be compatible with container registry tags.
			br := action.Ref.Name
			br = strings.ReplaceAll(br, "/", "-")
			br = invalidChar.ReplaceAllString(br, "")

			// Some branches should always publish images
			_, alwaysPublish := alwaysPublishBranch[br]

			pub := alwaysPublish || imgFlag.ShouldPublish

			// Override the environment name if requested.
			// Branches commits from branches that always publish tags
			// can not override the Tag to prevent issues with squash merges.
			ov := imgFlag.EnvNameOverride
			if ov != "" && tagRegex.MatchString(ov) && !alwaysPublish {
				br = ov
			}

			hash := action.SHA[:8]
			tag := fmt.Sprintf("%s-%s-%d", br, hash, action.RunNumber)

			return tag, br, pub
		}
	}

	// We don't want to publish but still need a placeholder tag for CI.
	return fmt.Sprintf("run-%d", action.RunNumber), "", false
}
