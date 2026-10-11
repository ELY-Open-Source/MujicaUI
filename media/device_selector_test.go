package media

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
)

var deviceTestList = MediaDevices{
	Microphones: []MediaDevice{{"m1", "Built-in Microphone"}, {"m2", "Studio Mic"}},
	Cameras:     []MediaDevice{},
	Speakers:    []MediaDevice{{"s1", "Headphones"}},
}

func deviceTester(sel *DeviceSelection, devices MediaDevices, disabled bool, n *int) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Padding(20).Width(420).Children(func() {
			if DeviceSelector(c, sel, devices, DeviceSelectorOptions{Disabled: disabled}).Changed() {
				*n++
			}
		})
	}, 500, 400)
}

func TestDeviceSelectorChoose(t *testing.T) {
	sel, n := DeviceSelection{Microphone: "m1"}, 0
	tt := deviceTester(&sel, deviceTestList, false, &n)
	if !tt.HasText("Built-in Microphone") || !tt.HasText("No devices found") {
		t.Fatalf("texts %v", tt.Texts())
	}
	if err := tt.Click("Built-in Microphone"); err != nil {
		t.Fatal(err)
	}
	if err := tt.Click("Studio Mic"); err != nil {
		t.Fatal(err)
	}
	if sel.Microphone != "m2" || n != 1 {
		t.Fatalf("selection %+v changes %d", sel, n)
	}
}

func TestDeviceSelectorHidesNilListsAndDisabled(t *testing.T) {
	sel, n := DeviceSelection{}, 0
	tt := deviceTester(&sel, MediaDevices{Speakers: deviceTestList.Speakers}, true, &n)
	if tt.HasText("Camera") || tt.HasText("Microphone") {
		t.Fatalf("nil lists shown: %v", tt.Texts())
	}
	_ = tt.Click("Choose a device")
	_ = tt.Click("Headphones")
	if sel.Speaker != "" || n != 0 {
		t.Fatal("disabled selector chose")
	}
}

func TestDeviceSelectorBadLists(t *testing.T) {
	sel := DeviceSelection{}
	mediaPanics(t, "duplicate", func(c *ui.Context) {
		DeviceSelector(c, &sel, MediaDevices{Cameras: []MediaDevice{{"a", "A"}, {"a", "B"}}}, DeviceSelectorOptions{})
	})
	gone := DeviceSelection{Camera: "x"}
	mediaPanics(t, "unknown selection", func(c *ui.Context) {
		DeviceSelector(c, &gone, MediaDevices{Cameras: []MediaDevice{{"a", "A"}}}, DeviceSelectorOptions{})
	})
}

func TestDeviceSelectorSnapshot(t *testing.T) {
	sel := DeviceSelection{Microphone: "m2", Speaker: "s1"}
	snapshot(t, "device_selector", 460, 190, func(c *ui.Context) {
		ui.Column(c).Padding(20).Width(420).Children(func() {
			DeviceSelector(c, &sel, deviceTestList, DeviceSelectorOptions{})
		})
	})
}
