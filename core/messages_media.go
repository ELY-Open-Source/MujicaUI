package core

func init() {
	addMessages(map[string][2]string{
		"image_cropper.label":       {"Crop area", "裁剪区域"},
		"image_cropper.rotateLeft":  {"Rotate left", "向左旋转"},
		"image_cropper.rotateRight": {"Rotate right", "向右旋转"},

		"image_annotator.label":     {"Annotations", "标注"},
		"image_annotator.tool":      {"Tool", "工具"},
		"image_annotator.select":    {"Select", "选择"},
		"image_annotator.rect":      {"Rectangle", "矩形"},
		"image_annotator.arrow":     {"Arrow", "箭头"},
		"image_annotator.text":      {"Text", "文字"},
		"image_annotator.textLabel": {"Annotation text", "标注文字"},
		"image_annotator.note":      {"Note", "注释"},
		"image_annotator.delete":    {"Delete annotation", "删除标注"},
		"image_annotator.count":     {"%d annotations", "%d 个标注"},
		"image_annotator.selected":  {"Annotation %d of %d selected: %s", "已选中第 %d 个标注（共 %d 个）：%s"},

		"image_thumbnail.failed":  {"Failed to load", "加载失败"},
		"image_thumbnail.loading": {"Loading", "加载中"},

		"lightbox.title":    {"Image", "图片"},
		"lightbox.close":    {"Close", "关闭"},
		"lightbox.position": {"%d of %d", "第 %d 张，共 %d 张"},
		"lightbox.previous": {"Previous image", "上一张"},
		"lightbox.next":     {"Next image", "下一张"},

		"video_player.label":      {"Video", "视频"},
		"video_player.empty":      {"No video", "暂无视频"},
		"video_player.subtitles":  {"Subtitles", "字幕"},
		"video_player.pip":        {"Picture in picture", "画中画"},
		"video_player.fullscreen": {"Full screen", "全屏"},

		"video_scrubber.label": {"Timeline", "时间轴"},

		"video_thumbnail_strip.label":    {"Frames", "缩略帧"},
		"video_thumbnail_strip.position": {"Frame %d of %d, %s", "第 %d 帧，共 %d 帧，%s"},
		"video_thumbnail_strip.previous": {"Previous frame", "上一帧"},
		"video_thumbnail_strip.next":     {"Next frame", "下一帧"},

		"audio_waveform.label": {"Waveform", "波形"},
		"audio_waveform.empty": {"No audio", "暂无音频"},

		"audio_spectrum.label": {"Spectrum", "频谱"},
		"audio_spectrum.empty": {"No signal", "无信号"},
		"audio_spectrum.band":  {"Band %d", "第 %d 段"},

		"volume_control.label":  {"Volume", "音量"},
		"volume_control.mute":   {"Mute", "静音"},
		"volume_control.unmute": {"Unmute", "取消静音"},

		"playback_speed_control.label": {"Playback speed", "播放速度"},

		"media_controls.label":    {"Playback", "播放控制"},
		"media_controls.play":     {"Play", "播放"},
		"media_controls.pause":    {"Pause", "暂停"},
		"media_controls.previous": {"Previous", "上一首"},
		"media_controls.next":     {"Next", "下一首"},
		"media_controls.shuffle":  {"Shuffle", "随机播放"},
		"media_controls.loopOff":  {"Repeat off", "不循环"},
		"media_controls.loopAll":  {"Repeat all", "列表循环"},
		"media_controls.loopOne":  {"Repeat one", "单曲循环"},

		"playlist.empty":   {"The playlist is empty", "播放列表为空"},
		"playlist.playing": {"Playing", "正在播放"},
		"playlist.remove":  {"Remove", "移除"},

		"subtitle_editor.start":   {"Start", "开始"},
		"subtitle_editor.end":     {"End", "结束"},
		"subtitle_editor.text":    {"Text", "文本"},
		"subtitle_editor.add":     {"Add cue", "添加字幕"},
		"subtitle_editor.remove":  {"Remove cue", "删除字幕"},
		"subtitle_editor.empty":   {"No cues yet", "还没有字幕"},
		"subtitle_editor.badTime": {"Use m:ss.mmm, as 1:05.250", "请使用 m:ss.mmm 格式，例如 1:05.250"},
		"subtitle_editor.order":   {"End must be after start", "结束时间须晚于开始时间"},
		"subtitle_editor.overlap": {"Overlaps cue %d", "与第 %d 条重叠"},

		"screen_recorder_controls.label":     {"Screen recording", "屏幕录制"},
		"screen_recorder_controls.idle":      {"Ready", "就绪"},
		"screen_recorder_controls.recording": {"Recording", "录制中"},
		"screen_recorder_controls.paused":    {"Paused", "已暂停"},
		"screen_recorder_controls.start":     {"Record", "开始录制"},
		"screen_recorder_controls.pause":     {"Pause", "暂停"},
		"screen_recorder_controls.resume":    {"Resume", "继续"},
		"screen_recorder_controls.stop":      {"Stop", "停止"},

		"camera_preview.label":    {"Camera", "摄像头"},
		"camera_preview.noSignal": {"No signal", "无信号"},
		"camera_preview.mirror":   {"Mirror", "镜像"},

		"mic_level_meter.label": {"Microphone level", "麦克风电平"},
		"mic_level_meter.peak":  {"peak", "峰值"},

		"device_selector.label":      {"Devices", "设备"},
		"device_selector.microphone": {"Microphone", "麦克风"},
		"device_selector.camera":     {"Camera", "摄像头"},
		"device_selector.speaker":    {"Speaker", "扬声器"},
		"device_selector.choose":     {"Choose a device", "选择设备"},
		"device_selector.none":       {"No devices found", "未找到设备"},

		"image_compare.before": {"Before", "之前"},
		"image_compare.after":  {"After", "之后"},
	})
}
