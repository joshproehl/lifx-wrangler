package lifxapiv1

import (
	"errors"
	"fmt"
	"github.com/pdf/golifx/common"
	"regexp"
	"strconv"
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
func (v *v1) GetForSelector(s string) ([]common.Light, error) {
	switch {
	case regexAll.MatchString(s):
		return v.watchdog.Client.GetLights()

	case regexLabel.MatchString(s):
		label := regexLabel.ReplaceAllLiteralString(s, "") // Get the actual label value by erasing the prefix that we matched.
		return wrapSingleResult(v.watchdog.Client.GetLightByLabel(label))

	case regexId.MatchString(s):
		idStr := regexId.ReplaceAllLiteralString(s, "") // Get the actual label value by erasing the prefix that we matched.
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Unable to parse ID to int. Error was:", err))
		}
		return wrapSingleResult(v.watchdog.Client.GetLightByID(id))

	case regexGroupId.MatchString(s):
		// TODO: Return results matched be selector
		return v.watchdog.Client.GetLights()

	case regexLocationId.MatchString(s):
		// TODO: Return results matched be selector
		return v.watchdog.Client.GetLights()

	case regexLocation.MatchString(s):
		// TODO: Return results matched be selector
		return v.watchdog.Client.GetLights()

	case regexSceneId.MatchString(s):
		// TODO: Return results matched be selector
		return v.watchdog.Client.GetLights()

	default:
		return nil, errors.New(fmt.Sprintf("Unknown selector type: '%s'", s))
	}
}

func wrapSingleResult(l common.Light, e error) ([]common.Light, error) {
	return []common.Light{l}, e
}
