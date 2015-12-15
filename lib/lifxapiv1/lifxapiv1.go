package lifxapiv1

import (
	wd "github.com/joshproehl/lifx-wrangler/lib/watchdog"
)

type v1 struct {
	watchdog *wd.Watchdog
}

func NewLifxCloudV1APICloneHandlers(w *wd.Watchdog) *v1 {
	return &v1{watchdog: w}
}
