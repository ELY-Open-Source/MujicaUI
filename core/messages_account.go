package core

func init() {
	addMessages(map[string][2]string{
		"settings_layout.title":       {"Settings", "设置"},
		"settings_search.placeholder": {"Search settings", "搜索设置"},
		"settings_search.none":        {"No settings match", "没有匹配的设置"},

		"theme_selector.label":  {"Appearance", "外观"},
		"theme_selector.System": {"System", "跟随系统"},
		"theme_selector.Light":  {"Light", "浅色"},
		"theme_selector.Dark":   {"Dark", "深色"},

		"accent_color_picker.label":      {"Accent color", "主色"},
		"accent_color_picker.low":        {"Low contrast", "对比度不足"},
		"accent_color_picker.low_note":   {"Marked colors are too faint for text on them.", "带标记的颜色对比度不足，不宜承载文字。"},
		"accent_color_picker.chosen_low": {"%s has low contrast; text on it may be hard to read.", "%s 对比度不足，其上的文字可能难以阅读。"},

		"shortcut_recorder.press":          {"Press keys…", "请按下快捷键…"},
		"shortcut_recorder.none":           {"none", "未设置"},
		"shortcut_recorder.conflict":       {"%s is already used by %s", "%s 已被“%s”使用"},
		"shortcut_recorder.needs_modifier": {"Add Ctrl, Alt or Command; F-keys may stand alone.", "请加上 Ctrl、Alt 或 Command；功能键可单独使用。"},

		"keyboard_shortcuts_list.search": {"Search shortcuts", "搜索快捷键"},

		"login_form.account":  {"Email or username", "邮箱或用户名"},
		"login_form.password": {"Password", "密码"},
		"login_form.remember": {"Remember me", "记住我"},
		"login_form.submit":   {"Sign in", "登录"},

		"oauth_buttons.continue": {"Continue with %s", "使用 %s 继续"},

		"two_factor_input.label":     {"Verification code", "验证码"},
		"two_factor_input.resend":    {"Resend code", "重新发送"},
		"two_factor_input.resend_in": {"Resend in %d:%02d", "%d:%02d 后可重新发送"},

		"profile_editor.avatar":  {"Change picture", "更换头像"},
		"profile_editor.name":    {"Name", "名称"},
		"profile_editor.bio":     {"Bio", "简介"},
		"profile_editor.unsaved": {"Unsaved changes", "有未保存的更改"},
		"profile_editor.discard": {"Discard", "放弃更改"},
		"profile_editor.save":    {"Save", "保存"},

		"account_switcher.accounts": {"Accounts", "账户"},
		"account_switcher.add":      {"Add account…", "添加账户…"},

		"workspace_switcher.label":  {"Workspaces", "工作区"},
		"workspace_switcher.search": {"Find a workspace", "查找工作区"},
		"workspace_switcher.none":   {"No workspace matches", "没有匹配的工作区"},

		"user_menu.label":    {"%s, %s", "%s，%s"},
		"user_menu.status0":  {"Online", "在线"},
		"user_menu.status1":  {"Away", "离开"},
		"user_menu.status2":  {"Do not disturb", "勿扰"},
		"user_menu.status3":  {"Invisible", "隐身"},
		"user_menu.settings": {"Settings", "设置"},
		"user_menu.sign_out": {"Sign out", "退出登录"},

		"session_list.label":           {"Signed-in devices", "已登录设备"},
		"session_list.current":         {"This device", "本设备"},
		"session_list.sign_out":        {"Sign out", "注销"},
		"session_list.sign_out_of":     {"Sign out of %s", "注销 %s"},
		"session_list.sign_out_others": {"Sign out of other sessions", "注销其他会话"},
		"session_list.confirm_title":   {"Sign out everywhere else?", "注销其他所有会话？"},
		"session_list.confirm_text":    {"%d other sessions will need to sign in again.", "其他 %d 个会话需要重新登录。"},
		"session_list.now":             {"Active now", "当前活跃"},
		"session_list.unknown":         {"Last active unknown", "最近活跃时间未知"},

		"api_key_manager.create":       {"Create key", "创建密钥"},
		"api_key_manager.create_title": {"New API key", "新建 API 密钥"},
		"api_key_manager.secret_title": {"Copy your new key", "复制新密钥"},
		"api_key_manager.name":         {"Key name", "密钥名称"},
		"api_key_manager.cancel":       {"Cancel", "取消"},
		"api_key_manager.done":         {"I've saved it", "我已保存"},
		"api_key_manager.once":         {"This key shows only once. Store it somewhere safe now.", "此密钥只显示一次，请立即妥善保存。"},
		"api_key_manager.copy":         {"Copy key", "复制密钥"},
		"api_key_manager.copied":       {"Key copied", "密钥已复制"},
		"api_key_manager.empty":        {"No API keys", "暂无 API 密钥"},
		"api_key_manager.revoke":       {"Revoke", "撤销"},
		"api_key_manager.revoke_named": {"Revoke %s", "撤销 %s"},
		"api_key_manager.revoke_title": {"Revoke %s?", "撤销 %s？"},
		"api_key_manager.revoke_text":  {"Apps using this key stop working at once.", "使用此密钥的应用将立即失效。"},
		"api_key_manager.dates":        {"Created %s · last used %s", "创建于 %s · 最近使用 %s"},
		"api_key_manager.never":        {"never", "从未"},

		"usage_quota.amount":      {"%s of %s %s", "%s / %s %s"},
		"usage_quota.near":        {"Near limit", "接近上限"},
		"usage_quota.over":        {"Limit reached", "已达上限"},
		"usage_quota.reset_due":   {"Resetting now", "正在重置"},
		"usage_quota.reset_hours": {"Resets in %d h", "%d 小时后重置"},
		"usage_quota.reset_on":    {"Resets on %s", "%s 重置"},

		"onboarding_wizard.progress": {"Setup progress", "引导进度"},
		"onboarding_wizard.count":    {"Step %d of %d", "第 %d 步，共 %d 步"},
		"onboarding_wizard.skip":     {"Skip", "跳过"},
		"onboarding_wizard.back":     {"Back", "上一步"},
		"onboarding_wizard.next":     {"Next", "下一步"},
		"onboarding_wizard.finish":   {"Finish", "完成"},

		"coachmark.count": {"%d of %d", "%d / %d"},
		"coachmark.skip":  {"Skip tour", "跳过引导"},
		"coachmark.back":  {"Back", "上一步"},
		"coachmark.next":  {"Next", "下一步"},
		"coachmark.done":  {"Done", "完成"},

		"whats_new_dialog.title": {"What's new", "更新内容"},
		"whats_new_dialog.hide":  {"Don't show again", "不再显示"},
		"whats_new_dialog.close": {"Got it", "知道了"},
		"whats_new_dialog.kind0": {"New", "新增"},
		"whats_new_dialog.kind1": {"Improved", "改进"},
		"whats_new_dialog.kind2": {"Fixed", "修复"},

		"feedback_widget.kind":    {"Feedback type", "反馈类型"},
		"feedback_widget.kind0":   {"Bug", "问题"},
		"feedback_widget.kind1":   {"Idea", "建议"},
		"feedback_widget.kind2":   {"Praise", "称赞"},
		"feedback_widget.kind3":   {"Other", "其他"},
		"feedback_widget.message": {"Message", "内容"},
		"feedback_widget.short":   {"Write at least %d characters.", "请至少写 %d 个字。"},
		"feedback_widget.send":    {"Send feedback", "发送反馈"},
		"feedback_widget.clear":   {"Clear", "清空"},
	})
}
