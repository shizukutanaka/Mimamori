package i18n

import (
	"fmt"
	"strings"
)

var messages = map[string]map[string]string{
	"en": {
		"title":             "Mimamori — Home Safety Check",
		"safe":              "Safe",
		"caution":           "Caution",
		"danger":            "Danger",
		"score":             "Home safety: %s  %d points",
		"device_counts":     "Devices: %d (%d new)",
		"findings_heading":  "Things to address",
		"devices_heading":   "Devices",
		"no_problems":       "No problems found",
		"first_run":         "This run learned the devices on your network. New devices will be reported from the next run.",
		"new_badge":         "new",
		"router":            "Router",
		"randomized_device": "Phone or PC (private Wi-Fi address)",
		"telnet_open":       "Telnet is open on %s. This old remote-control service is often abused to take over devices.",
		"telnet_fix":        "Turn off Telnet in the device settings, update its firmware, and change any default password.",
		"ftp_open":          "FTP file access is open on %s. It may expose files or account details.",
		"ftp_fix":           "Turn off FTP if you do not need it. If it is needed, use a secure alternative and a unique password.",
		"vnc_open":          "VNC remote access is open on %s. Anyone who can reach it may be able to control the screen.",
		"vnc_fix":           "Turn off VNC if you do not need remote access. Otherwise restrict access and use a strong unique password.",
		"router_upnp":       "The router's UPnP is enabled. Devices in your home may open doors to the Internet automatically.",
		"router_upnp_fix":   "If devices such as game consoles do not need it, turn off UPnP in the router settings.",
		"new_device":        "A new device connected: %s",
		"new_device_fix":    "If you do not recognize it, change your Wi-Fi password and reconnect only your own devices.",
		"first_run_json":    "first run",
		"score_label":       "Score",
		"points_recovered":  "+%d points",
	},
	"ja": {
		"title":             "Mimamori — 家の安全診断",
		"safe":              "安全",
		"caution":           "注意",
		"danger":            "危険",
		"score":             "家の安全状態: %s  %d点",
		"device_counts":     "機器: %d台（新しい機器 %d台）",
		"findings_heading":  "対処が必要なこと",
		"devices_heading":   "機器一覧",
		"no_problems":       "危険はありません",
		"first_run":         "今回は家の機器を記録しました。新しい機器は次回からお知らせします。",
		"new_badge":         "新規",
		"router":            "ルーター",
		"randomized_device": "スマホ・パソコン（Wi-Fi のプライベートアドレス）",
		"telnet_open":       "「%s」で Telnet（古い遠隔操作の仕組み）が開いています。乗っ取りに悪用されやすい状態です。",
		"telnet_fix":        "機器の設定で Telnet を無効にし、ファームウェアを更新してください。初期パスワードのままなら変更してください。",
		"ftp_open":          "「%s」で FTP（ファイル転送の仕組み）が開いています。ファイルやアカウント情報が見られるおそれがあります。",
		"ftp_fix":           "使っていなければ FTP を無効にしてください。必要な場合は安全な方法に切り替え、推測されにくいパスワードを使ってください。",
		"vnc_open":          "「%s」で VNC（遠隔操作の仕組み）が開いています。画面を操作されるおそれがあります。",
		"vnc_fix":           "遠隔操作が不要なら VNC を無効にしてください。必要な場合は接続元を制限し、強いパスワードを設定してください。",
		"router_upnp":       "ルーターの UPnP が有効です。家の中の機器が勝手にインターネット側へ入口を開けられます。",
		"router_upnp_fix":   "ゲーム機などで不要なら、ルーターの設定画面で UPnP を無効にしてください。",
		"new_device":        "新しい機器が接続されました: %s",
		"new_device_fix":    "心当たりがなければ Wi-Fi のパスワードを変更し、自分の機器だけをつなぎ直してください。",
		"first_run_json":    "初回",
		"score_label":       "点数",
		"points_recovered":  "+%d点",
	},
}

func Detect(getenv func(string) string) string {
	for _, key := range []string{"LC_ALL", "LC_MESSAGES", "LANG"} {
		if value := getenv(key); value != "" {
			return language(value)
		}
	}
	return language(systemLocale())
}

func language(locale string) string {
	if strings.HasPrefix(strings.ToLower(locale), "ja") {
		return "ja"
	}
	return "en"
}

func T(lang, key string, args ...any) string {
	text, ok := messages[lang][key]
	if !ok {
		text, ok = messages["en"][key]
	}
	if !ok {
		text = key
	}
	if len(args) == 0 {
		return text
	}
	return fmt.Sprintf(text, args...)
}

func Supported() []string {
	return []string{"en", "ja"}
}
