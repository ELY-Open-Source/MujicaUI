package mujica

import "testing"

func TestOrbitLoaderSnapshot(t *testing.T) { loaderSnapshot(t, "orbit_loader", OrbitLoader) }

func TestOrbitLoaderBehavior(t *testing.T) { testLoaderBehavior(t, OrbitLoader) }
