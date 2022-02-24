package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/snapADDY/actions/pkg/actions"
	"github.com/snapADDY/actions/pkg/msgflag"
)

func main() {
	ctx := actions.NewContext()
	fmt.Printf("Actions Context: %#v\n", ctx)

	commitMsg := os.Getenv("INPUT_COMMITMESSAGE")
	flags := msgflag.Parse(commitMsg)

	fmt.Printf("Parsed commit message: %s\nFlags: %#v\n", commitMsg, flags)

	imgFlag := parseImgFlag(flags["img"])
	tag, publish := makeImgTag(ctx, imgFlag)

	fmt.Printf("Prepared Tag: %s; Should publish %v\n", tag, publish)

	actions.SetOutput("tag-name", tag)
	actions.SetOutput("publish", publish)
}

var alwaysPublishBranch = map[string]bool{
	"master":  true,
	"main":    true,
	"staging": true,
}

func makeImgTag(action actions.Context, imgFlag ImgFlag) (string, bool) {
	tagRegex := regexp.MustCompile(`^\w[\w.-]{0,127}$`)
	invalidChar := regexp.MustCompile(`[^\w.-]+`)

	if action.EventName == "push" {
		if action.Ref.Type == actions.Tag && tagRegex.MatchString(action.Ref.Name) {
			return action.Ref.Name, true
		} else if action.Ref.Type == actions.Branch {
			br := action.Ref.Name
			br = strings.ReplaceAll(br, "/", "-")
			br = invalidChar.ReplaceAllString(br, "")

			_, pub := alwaysPublishBranch[br]
			pub = pub || imgFlag.ShouldPublish

			ov := imgFlag.EnvNameOverride
			if ov != "" && tagRegex.MatchString(ov) {
				br = ov
			}

			hash := action.SHA[:8]
			tag := fmt.Sprintf("%s-%s-%d", br, hash, action.RunNumber)

			return tag, pub
		}
	}

	return fmt.Sprintf("run-%d", action.RunNumber), false
}
