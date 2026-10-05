package mujica

import "testing"

func TestWaveLoaderSnapshot(t *testing.T) { loaderSnapshot(t, "wave_loader", WaveLoader) }

func TestWaveLoaderBehavior(t *testing.T) { testLoaderBehavior(t, WaveLoader) }
