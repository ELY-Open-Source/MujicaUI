package feedback

import (
	"testing"
)

func TestPulseLoaderSnapshot(t *testing.T) { loaderSnapshot(t, "pulse_loader", PulseLoader) }

func TestPulseLoaderBehavior(t *testing.T) { testLoaderBehavior(t, PulseLoader) }
