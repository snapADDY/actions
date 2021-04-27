package msgflag

// ImgFlag describes all options for docker images. The Flag uses `img` as prefix and has multiple optional settings.
// Examles of valid flags:
//	[img]         // Shortest valid flag;
//	[img:testenv] // Override the branch name with "testenv"
//	[img::dev]    // Dev build enabled
//	[img:testenv:dev]
type ImgFlag struct {
	// EnvNameOverride is used to override the branch name of a docker tag.
	EnvNameOverride string
	// ShouldPublish determines if the built package should be published. True when a valid img flag is found.
	ShouldPublish bool
	// BuildDev enables a dev build mode for images. May not be implemented in every project.
	BuildDev bool
}

var blockedEnvOverrides = map[string]bool{
	"master":  true,
	"main":    true,
	"staging": true,
}

func parseImgFlag(flag []string) ImgFlag {
	if len(flag) == 0 || flag[0] != "img" {
		return ImgFlag{}
	}

	fl := ImgFlag{
		ShouldPublish: true,
	}

	if len(flag) >= 2 {
		env := flag[1]
		_, blocked := blockedEnvOverrides[env]
		if !blocked {
			fl.EnvNameOverride = env
		}
	}

	if len(flag) >= 3 && flag[2] == "dev" {
		fl.BuildDev = true
	}

	return fl
}
