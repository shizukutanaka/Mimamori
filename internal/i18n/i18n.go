package i18n

import (
	"fmt"
	"strings"
)

var messages = map[string]map[string]string{
	"en": {
		"title":                  "Mimamori — Home Safety Check",
		"safe":                   "Safe",
		"caution":                "Caution",
		"danger":                 "Danger",
		"score":                  "Home safety: %s  %d points",
		"device_counts":          "Devices: %d (%d new)",
		"findings_heading":       "Things to address",
		"devices_heading":        "Devices",
		"no_problems":            "No problems found",
		"first_run":              "This run learned the devices on your network. New devices will be reported from the next run.",
		"new_badge":              "new",
		"router":                 "Router",
		"randomized_device":      "Phone or PC (private Wi-Fi address)",
		"telnet_open":            "Telnet is open on %s. This old remote-control service is often abused to take over devices.",
		"telnet_open_fix":        "Turn off Telnet in the device settings, update its firmware, and change any default password.",
		"ftp_open":               "FTP file access is open on %s. It may expose files or account details.",
		"ftp_open_fix":           "Turn off FTP if you do not need it. If it is needed, use a secure alternative and a unique password.",
		"vnc_open":               "VNC remote access is open on %s. Anyone who can reach it may be able to control the screen.",
		"vnc_open_fix":           "Turn off VNC if you do not need remote access. Otherwise restrict access and use a strong unique password.",
		"router_upnp":            "The router's UPnP is enabled. Devices in your home may open doors to the Internet automatically.",
		"router_upnp_fix":        "If devices such as game consoles do not need it, turn off UPnP in the router settings.",
		"new_device":             "A new device connected: %s",
		"new_device_fix":         "If you do not recognize it, change your Wi-Fi password and reconnect only your own devices.",
		"port_forward":           "%s can be reached from the Internet through port %d (opened by UPnP).",
		"port_forward_fix":       "If you don't recognize this, remove it in the router settings or turn off UPnP. Game consoles and some apps open these on purpose.",
		"port_forward_risky":     "Remote access on %s is open to the whole Internet through port %d. Attackers scan for this constantly.",
		"port_forward_risky_fix": "Remove this port forwarding in the router settings and turn off UPnP. If you need remote access, use the router's VPN feature instead.",
		"points_recovered":       "(+%d points)",
		"press_enter":            "Press Enter to close...",
	},
	"ja": {
		"title":                  "Mimamori — 家の安全診断",
		"safe":                   "安全",
		"caution":                "注意",
		"danger":                 "危険",
		"score":                  "家の安全状態: %s  %d点",
		"device_counts":          "機器: %d台（新しい機器 %d台）",
		"findings_heading":       "対処が必要なこと",
		"devices_heading":        "機器一覧",
		"no_problems":            "危険はありません",
		"first_run":              "今回は家の機器を記録しました。新しい機器は次回からお知らせします。",
		"new_badge":              "新規",
		"router":                 "ルーター",
		"randomized_device":      "スマホ・パソコン（Wi-Fi のプライベートアドレス）",
		"telnet_open":            "「%s」で Telnet（古い遠隔操作の仕組み）が開いています。乗っ取りに悪用されやすい状態です。",
		"telnet_open_fix":        "機器の設定で Telnet を無効にし、ファームウェアを更新してください。初期パスワードのままなら変更してください。",
		"ftp_open":               "「%s」で FTP（ファイル転送の仕組み）が開いています。ファイルやアカウント情報が見られるおそれがあります。",
		"ftp_open_fix":           "使っていなければ FTP を無効にしてください。必要な場合は安全な方法に切り替え、推測されにくいパスワードを使ってください。",
		"vnc_open":               "「%s」で VNC（遠隔操作の仕組み）が開いています。画面を操作されるおそれがあります。",
		"vnc_open_fix":           "遠隔操作が不要なら VNC を無効にしてください。必要な場合は接続元を制限し、強いパスワードを設定してください。",
		"router_upnp":            "ルーターの UPnP が有効です。家の中の機器が勝手にインターネット側へ入口を開けられます。",
		"router_upnp_fix":        "ゲーム機などで不要なら、ルーターの設定画面で UPnP を無効にしてください。",
		"new_device":             "新しい機器が接続されました: %s",
		"new_device_fix":         "心当たりがなければ Wi-Fi のパスワードを変更し、自分の機器だけをつなぎ直してください。",
		"port_forward":           "「%s」は UPnP によってインターネットからポート %d で接続できる状態です。",
		"port_forward_fix":       "心当たりがなければ、ルーターの設定でこの転送を削除するか、UPnP を無効にしてください。ゲーム機や一部のアプリが意図的に設定する場合もあります。",
		"port_forward_risky":     "「%s」の遠隔操作・ファイル共有が、ポート %d を通じてインターネット全体に公開されています。攻撃者は常にこの状態を探しています。",
		"port_forward_risky_fix": "ルーターの設定でこのポート転送を削除し、UPnP を無効にしてください。遠隔アクセスが必要なら、ルーターの VPN 機能を使ってください。",
		"points_recovered":       "（+%d点）",
		"press_enter":            "Enter キーを押すと閉じます...",
	},
}

func Detect(getenv func(string) string) string {
	for _, value := range strings.Split(getenv("LANGUAGE"), ":") {
		if lang, ok := Match(value); ok {
			return lang
		}
	}
	for _, key := range []string{"LC_ALL", "LC_MESSAGES", "LANG"} {
		if value := getenv(key); value != "" {
			if lang, ok := Match(value); ok {
				return lang
			}
			return "en"
		}
	}
	if lang, ok := Match(systemLocale()); ok {
		return lang
	}
	return "en"
}

func Match(tag string) (string, bool) {
	tag = strings.TrimSpace(strings.ToLower(tag))
	if index := strings.IndexAny(tag, ".@"); index >= 0 {
		tag = tag[:index]
	}
	tag = strings.ReplaceAll(tag, "-", "_")
	if tag == "" || tag == "c" || tag == "posix" {
		return "", false
	}
	parts := strings.Split(tag, "_")
	primary := parts[0]
	if primary == "zh" {
		for _, part := range parts[1:] {
			switch part {
			case "hant":
				return "zh-Hant", true
			case "hans":
				return "zh-Hans", true
			}
		}
		for _, part := range parts[1:] {
			switch part {
			case "tw", "hk", "mo":
				return "zh-Hant", true
			}
		}
		return "zh-Hans", true
	}
	for _, lang := range Supported() {
		if strings.SplitN(strings.ToLower(lang), "-", 2)[0] == primary {
			return lang, true
		}
	}
	return "", false
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
	return []string{"en", "ja", "zh-Hans", "zh-Hant", "es", "pt", "fr", "de", "ru", "ko", "id"}
}
