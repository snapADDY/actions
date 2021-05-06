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

	tag, publish := makeImgTag(ctx, flags)

	fmt.Printf("Prepared Tag: %s; Should publish %v\n", tag, publish)

	actions.SetOutput("tag-name", tag)
	actions.SetOutput("publish", publish)

}

var alwaysPublishBranch = map[string]bool{
	"master":  true,
	"main":    true,
	"staging": true,
}

func makeImgTag(ctx actions.Context, flags msgflag.Flags) (string, bool) {
	tagRegex := regexp.MustCompile(`^\w[\w.-]{0,127}$`)
	invalidChar := regexp.MustCompile(`[^\w.-]+`)

	if ctx.EventName == "push" {
		if ctx.Ref.Type == actions.Tag && tagRegex.MatchString(ctx.Ref.Name) {
			return ctx.Ref.Name, true

		} else if ctx.Ref.Type == actions.Branch {
			br := ctx.Ref.Name
			br = strings.ReplaceAll(br, "/", "-")
			br = invalidChar.ReplaceAllString(br, "")

			_, pub := alwaysPublishBranch[br]
			pub = pub || flags.Img.ShouldPublish

			ov := flags.Img.EnvNameOverride
			if ov != "" && tagRegex.MatchString(ov) {
				br = ov
			}

			hash := ctx.SHA[:8]
			tag := fmt.Sprintf("%s-%s-%d", br, hash, ctx.RunNumber)

			return tag, pub
		}
	}

	return fmt.Sprintf("run-%d", ctx.RunNumber), false
}
