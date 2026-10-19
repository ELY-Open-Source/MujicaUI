package main

import (
	"errors"
	"fmt"
	"image"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/data"
	"github.com/ELY-Open-Source/MujicaUI/input"
	"github.com/ELY-Open-Source/MujicaUI/media"
	"github.com/ELY-Open-Source/MujicaUI/navigation"
)

func imageCropperDemo(c *ui.Context, log func(string)) {
	type state struct {
		crop     media.ImageCrop
		aspect   int
		disabled bool
	}
	s := ui.Local(c.Root(), "341", func() state { return state{} })
	aspects := []float64{0, 1, 16.0 / 9, 4.0 / 3}
	mediaSection(c, "Parameters", func() {
		ui.Row(c).Gap(16).AlignItems(ui.Center).Children(func() {
			before := s.aspect
			navigation.SegmentedControl(c, &s.aspect, []navigation.Segment{{Label: "Free"}, {Label: "1:1"}, {Label: "16:9"}, {Label: "4:3"}}, navigation.SegmentedControlOptions{Label: "Aspect"})
			if s.aspect != before {
				s.crop.Rect = image.Rectangle{} // refit to the new ratio
				log(fmt.Sprintf("aspect = %v", aspects[s.aspect]))
			}
			mediaCheck(c, log, &s.disabled, "Disabled")
		})
	})
	mediaSection(c, "Live", func() {
		if media.ImageCropper(c, mediaPhotos[0], &s.crop, media.ImageCropperOptions{Aspect: aspects[s.aspect], Disabled: s.disabled}).Changed() {
			log(fmt.Sprintf("crop %v rotated %d°", s.crop.Rect, s.crop.Rotation))
		}
	})
}

func imageAnnotatorDemo(c *ui.Context, log func(string)) {
	type state struct {
		notes    []media.Annotation
		disabled bool
	}
	s := ui.Local(c.Root(), "342", func() state {
		return state{notes: []media.Annotation{{Kind: media.AnnotateRect, From: image.Pt(150, 60), To: image.Pt(330, 230)}}}
	})
	mediaSection(c, "Parameters", func() {
		ui.Row(c).Gap(16).AlignItems(ui.Center).Children(func() {
			mediaCheck(c, log, &s.disabled, "Disabled")
			if input.Button(c, "Export JSON", input.ButtonOptions{Variant: input.Secondary}).Clicked() {
				data, err := media.ExportAnnotations(s.notes)
				if err != nil {
					log("export failed: " + err.Error())
				} else {
					c.WriteClipboard(string(data))
					log(fmt.Sprintf("exported %d bytes to the clipboard", len(data)))
				}
			}
		})
	})
	mediaSection(c, "Live", func() {
		if media.ImageAnnotator(c, mediaBitmaps[1], &s.notes, media.ImageAnnotatorOptions{Disabled: s.disabled}).Changed() {
			log(fmt.Sprintf("annotations: %d", len(s.notes)))
		}
	})
}

func imageThumbnailDemo(c *ui.Context, log func(string)) {
	type state struct {
		selected                  int
		loading, failed, disabled bool
	}
	s := ui.Local(c.Root(), "343", func() state { return state{} })
	mediaSection(c, "Parameters", func() {
		ui.Row(c).Gap(16).Children(func() {
			mediaCheck(c, log, &s.loading, "Loading")
			mediaCheck(c, log, &s.failed, "Failed")
			mediaCheck(c, log, &s.disabled, "Disabled")
		})
	})
	mediaSection(c, "Live", func() {
		ui.Row(c).Gap(12).Wrap().Children(func() {
			for i, b := range mediaBitmaps {
				var err error
				if s.failed && i == 2 {
					err = errors.New("the file could not be read")
				}
				o := media.ImageThumbnailOptions{Alt: fmt.Sprintf("Scene %d", i+1), Aspect: 4.0 / 3, Width: 120,
					Selected: s.selected == i, Loading: s.loading && i == 1, Err: err, Disabled: s.disabled}
				if media.ImageThumbnail(c, b, o).Changed() {
					s.selected = i
					log(fmt.Sprintf("selected scene %d", i+1))
				}
			}
		})
	})
}

func lightboxDemo(c *ui.Context, log func(string)) {
	type state struct {
		open  bool
		index int
	}
	s := ui.Local(c.Root(), "344", func() state { return state{} })
	images := []data.ViewerImage{
		{Source: mediaBitmaps[0], Alt: "The grand hall at dusk"},
		{Source: mediaBitmaps[1], Alt: "The ivory gallery"},
		{Source: mediaBitmaps[2], Alt: "The blue chapel"},
	}
	mediaSection(c, "Live", func() {
		ui.Row(c).Gap(12).Children(func() {
			for i, im := range images {
				if media.ImageThumbnail(c, im.Source, media.ImageThumbnailOptions{Alt: im.Alt, Aspect: 4.0 / 3}).Changed() {
					s.open, s.index = true, i
					log(fmt.Sprintf("open at %d", i+1))
				}
			}
		})
		was := s.open
		if media.Lightbox(c, &s.open, &s.index, images).Changed() {
			log(fmt.Sprintf("image %d", s.index+1))
		}
		if was && !s.open {
			log("closed")
		}
	})
}

func videoStripDemo(c *ui.Context, log func(string)) {
	type state struct {
		current  int
		disabled bool
	}
	s := ui.Local(c.Root(), "347", func() state { return state{} })
	frames := make([]media.StripFrame, 240)
	for i := range frames {
		frames[i] = media.StripFrame{Image: mediaBitmaps[i%3], At: time.Duration(i) * 5 * time.Second}
	}
	mediaSection(c, "Parameters", func() { mediaCheck(c, log, &s.disabled, "Disabled") })
	mediaSection(c, "Live", func() {
		if media.VideoThumbnailStrip(c, &s.current, frames, media.VideoThumbnailStripOptions{Disabled: s.disabled}).Changed() {
			log("jump to " + mediaClockLabel(frames[s.current].At))
		}
	})
}

func compareDemo(c *ui.Context, log func(string)) {
	type state struct {
		split    float64
		disabled bool
	}
	s := ui.Local(c.Root(), "360", func() state { return state{split: 0.5} })
	mediaSection(c, "Parameters", func() { mediaCheck(c, log, &s.disabled, "Disabled") })
	mediaSection(c, "Live", func() {
		ui.Column(c).MaxWidth(560).Children(func() {
			if media.ImageCompare(c, mediaBitmaps[0], mediaBitmaps[1], &s.split, media.ImageCompareOptions{Disabled: s.disabled}).Changed() {
				log(fmt.Sprintf("split %.0f%%", s.split*100))
			}
		})
	})
}
