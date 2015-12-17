package lifxapiv1

import (
	"errors"
	"fmt"
	wd "github.com/joshproehl/lifx-wrangler/lib/watchdog"
	"regexp"
)

// Set up the regexes needed for handling lifxAPI selector strings
var (
	regexAll        = regexp.MustCompile(`^all$`)
	regexLabel      = regexp.MustCompile(`^label:`)
	regexId         = regexp.MustCompile(`^id:`)
	regexGroupId    = regexp.MustCompile(`^group_id:`)
	regexLocationId = regexp.MustCompile(`^location_id:`)
	regexLocation   = regexp.MustCompile(`^location:`)
	regexSceneId    = regexp.MustCompile(`^scene_id:`)
)

// GetForSelector takes a selector string used by the LIFX HTTP API and returns the lights managed by the watchdog that are found by that selector
func (v *v1) GetForSelector(s string) ([]*wd.Light, error) {
	switch {
	case regexAll.MatchString(s):
		return v.watchdog.LightCollection.All(), nil

	case regexLabel.MatchString(s):
		label := regexLabel.ReplaceAllLiteralString(s, "") // Get the actual label value by erasing the prefix that we matched.
		lights := v.watchdog.LightCollection.GetForLabel(label)
		if len(lights) == 0 {
			return nil, errors.New(fmt.Sprintf("Could not find light with label: %s", label))
		}

		return lights, nil

	case regexId.MatchString(s):
		lights := v.watchdog.LightCollection.All()
		return lights, nil

	case regexGroupId.MatchString(s):

		lights := v.watchdog.LightCollection.All()
		if len(lights) == 0 {
			return nil, errors.New(fmt.Sprintf("Could not find group_id: %s", s))
		}
		return lights, nil

	case regexLocationId.MatchString(s):
		lights := v.watchdog.LightCollection.All()
		return lights, nil

	case regexLocation.MatchString(s):
		lights := v.watchdog.LightCollection.All()
		return lights, nil

	case regexSceneId.MatchString(s):
		lights := v.watchdog.LightCollection.All()
		return lights, nil

	default:
		return nil, errors.New(fmt.Sprintf("Unknown selector type: '%s'", s))
	}
}
