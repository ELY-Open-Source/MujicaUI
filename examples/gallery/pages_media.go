package main

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/input"
	"github.com/ELY-Open-Source/MujicaUI/theme"
)

const mediaCategory = "Media"

var mediaPages = []page{
	{ID: "341", Name: "ImageCropper", ZhName: "图片裁剪", Category: mediaCategory, Demo: imageCropperDemo, Code: imageCropperCode,
		Keys: "方向键移动裁剪框；Shift+方向键缩放；R / Shift+R 向右 / 向左旋转。"},
	{ID: "342", Name: "ImageAnnotator", ZhName: "图片标注", Category: mediaCategory, Demo: imageAnnotatorDemo, Code: imageAnnotatorCode,
		Keys: "[ / ] 切换选中标注；方向键移动；Delete 删除；Escape 取消选择。"},
	{ID: "343", Name: "ImageThumbnail", ZhName: "缩略图", Category: mediaCategory, Demo: imageThumbnailDemo, Code: imageThumbnailCode,
		Keys: "Tab 聚焦；Enter / Space 激活。"},
	{ID: "344", Name: "Lightbox", ZhName: "灯箱", Category: mediaCategory, Demo: lightboxDemo, Code: lightboxCode,
		Keys: "← / → 切换；Home / End 首尾；Escape 关闭并恢复焦点。"},
	{ID: "345", Name: "VideoPlayer", ZhName: "视频播放器", Category: mediaCategory, Demo: videoPlayerDemo, Code: videoPlayerCode,
		Keys: "画面聚焦时 Space / K 播放暂停，← / → 跳转 5 秒，F 全屏，C 字幕。"},
	{ID: "346", Name: "VideoScrubber", ZhName: "视频时间轴", Category: mediaCategory, Demo: videoScrubberDemo, Code: videoScrubberCode,
		Keys: "← / → 跳转一个步长；Home / End 开头结尾。"},
	{ID: "347", Name: "VideoThumbnailStrip", ZhName: "缩略帧条", Category: mediaCategory, Demo: videoStripDemo, Code: videoStripCode,
		Keys: "Tab 聚焦帧条；← / → 上下帧；Home / End 首尾帧。"},
	{ID: "348", Name: "AudioPlayer", ZhName: "音频播放器", Category: mediaCategory, Demo: audioPlayerDemo, Code: audioPlayerCode,
		Keys: "Space 播放或暂停；进度条聚焦时 ← / → 跳转 5 秒。"},
	{ID: "349", Name: "AudioWaveform", ZhName: "音频波形", Category: mediaCategory, Demo: audioWaveformDemo, Code: audioWaveformCode,
		Keys: "← / → 跳转 5 秒；Shift+← / → 调整选区；Escape 清除选区。"},
	{ID: "350", Name: "AudioSpectrum", ZhName: "音频频谱", Category: mediaCategory, Demo: audioSpectrumDemo, Code: audioSpectrumCode,
		Keys: "Tab 聚焦；← / → 逐个读出频带；Home / End 首尾。"},
	{ID: "351", Name: "VolumeControl", ZhName: "音量控制", Category: mediaCategory, Demo: volumeDemo, Code: volumeCode,
		Keys: "← / → 音量 ±5%；Home / End；M 静音。"},
	{ID: "352", Name: "PlaybackSpeedControl", ZhName: "倍速选择", Category: mediaCategory, Demo: speedDemo, Code: speedCode,
		Keys: "Enter / Space 打开；↑ / ↓ 移动；Enter 确认；Escape 关闭。"},
	{ID: "353", Name: "MediaControls", ZhName: "媒体控制", Category: mediaCategory, Demo: mediaControlsDemo, Code: mediaControlsCode,
		Keys: "Tab 进入；← / → 在按钮间移动；Enter / Space 按下。"},
	{ID: "354", Name: "Playlist", ZhName: "播放列表", Category: mediaCategory, Demo: playlistDemo, Code: playlistCode,
		Keys: "↑ / ↓ 选择；Enter 播放；Alt+↑ / ↓ 移动；Delete 移除；可拖动排序。"},
	{ID: "355", Name: "SubtitleEditor", ZhName: "字幕编辑器", Category: mediaCategory, Demo: subtitleDemo, Code: subtitleCode,
		Keys: "Tab 在一行的开始、结束、文本、删除间移动；时间格式 m:ss.mmm。"},
	{ID: "356", Name: "ScreenRecorderControls", ZhName: "录屏控制", Category: mediaCategory, Demo: recorderDemo, Code: recorderCode,
		Keys: "Tab 在按钮间移动；Enter / Space 按下。"},
	{ID: "357", Name: "CameraPreview", ZhName: "摄像头预览", Category: mediaCategory, Demo: cameraDemo, Code: cameraCode,
		Keys: "Tab 聚焦预览；M 切换镜像。"},
	{ID: "358", Name: "MicLevelMeter", ZhName: "麦克风电平", Category: mediaCategory, Demo: micDemo, Code: micCode,
		Keys: "不获取焦点；读数随描述朗读。"},
	{ID: "359", Name: "DeviceSelector", ZhName: "设备选择", Category: mediaCategory, Demo: deviceDemo, Code: deviceCode,
		Keys: "Tab 在选择框间移动；Enter 打开；输入字母跳转。"},
	{ID: "360", Name: "ImageCompare", ZhName: "前后对比", Category: mediaCategory, Demo: compareDemo, Code: compareCode,
		Keys: "← / → 移动 1%；Home / End；可拖动。"},
}

// mediaSection titles a part of a demo.
func mediaSection(c *ui.Context, title string, fn func()) {
	ui.Column(c).Gap(8).FillWidth().Children(func() {
		ui.Text(c, title).FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(core.Tokens(c).TextMuted)
		fn()
	})
}

// mediaCheck is a parameter checkbox that logs its changes.
func mediaCheck(c *ui.Context, log func(string), v *bool, label string) {
	before := *v
	input.Checkbox(c, v, label, input.CheckboxOptions{})
	if *v != before {
		log(fmt.Sprintf("%s = %v", label, *v))
	}
}

// mediaPicture draws a w×h court scene in variant hues.
func mediaPicture(w, h, variant int) *image.RGBA {
	walls := []color.RGBA{{36, 28, 32, 255}, {242, 235, 221, 255}, {73, 42, 54, 255}}
	arch := []color.RGBA{{125, 32, 52, 255}, {173, 145, 100, 255}, {53, 94, 128, 255}}
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	cx := float64(w) / 2
	for y := range h {
		for x := range w {
			c := walls[variant%3]
			dx, dy := (float64(x)-cx)/(float64(w)*0.28), (float64(y)-float64(h)*0.55)/(float64(h)*0.5)
			if dy > 0 && math.Abs(dx) < 1 || dx*dx+dy*dy < 1 {
				c = arch[variant%3]
			}
			if x%(w/6) < w/40 || y > h*7/8 {
				c = color.RGBA{189, 161, 114, 255}
			}
			if math.Abs(float64(x)-cx)/float64(w)+math.Abs(float64(y)-float64(h)*0.45)/float64(h) < 0.06 {
				c = color.RGBA{251, 247, 239, 255}
			}
			img.Set(x, y, c)
		}
	}
	return img
}

var (
	mediaPhotos  = []*image.RGBA{mediaPicture(480, 270, 0), mediaPicture(480, 270, 1), mediaPicture(360, 360, 2)}
	mediaBitmaps = func() []*ui.Bitmap {
		out := make([]*ui.Bitmap, len(mediaPhotos))
		for i, p := range mediaPhotos {
			out[i] = ui.NewBitmap(p)
		}
		return out
	}()
)

// mediaClockLabel writes d as m:ss for event logs.
func mediaClockLabel(d time.Duration) string {
	s := int(d / time.Second)
	return fmt.Sprintf("%d:%02d", s/60, s%60)
}
