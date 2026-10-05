package feedback

import (
	"testing"
)

func TestBarsLoaderSnapshot(t *testing.T) { loaderSnapshot(t, "bars_loader", BarsLoader) }

func TestBarsLoaderBehavior(t *testing.T) { testLoaderBehavior(t, BarsLoader) }
