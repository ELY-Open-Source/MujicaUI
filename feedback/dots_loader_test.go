package feedback

import (
	"testing"
)

func TestDotsLoaderSnapshot(t *testing.T) { loaderSnapshot(t, "dots_loader", DotsLoader) }

func TestDotsLoaderBehavior(t *testing.T) { testLoaderBehavior(t, DotsLoader) }
