package i18n

var catalogZhHant = map[string]string{
	"title":                  "Mimamori — 家庭網路安全檢查",
	"safe":                   "安全",
	"caution":                "注意",
	"danger":                 "危險",
	"score":                  "家庭網路安全：%s  %d 分",
	"device_counts":          "裝置：%d 台（新增 %d 台）",
	"findings_heading":       "需要處理的問題",
	"devices_heading":        "裝置",
	"no_problems":            "沒有發現問題",
	"first_run":              "本次已記錄家中網路的裝置。從下次執行起會提醒您有新裝置。",
	"new_badge":              "新增",
	"router":                 "路由器",
	"randomized_device":      "手機或電腦（使用 Wi-Fi 私人位址）",
	"telnet_open":            "「%s」上的 Telnet 已開啟。這種舊式遠端控制服務容易被用來控制裝置。",
	"telnet_open_fix":        "請在裝置設定中關閉 Telnet 並安裝可用的更新。如果仍使用預設密碼，也請一併更改。",
	"ftp_open":               "「%s」上的 FTP 檔案存取已開啟，檔案或帳戶資料可能因此外洩。",
	"ftp_open_fix":           "如果不需要，請關閉 FTP。如果需要使用，請改用更安全的方式並設定專用密碼。",
	"vnc_open":               "「%s」上的 VNC 遠端操作已開啟，能連線的人可能會控制畫面。",
	"vnc_open_fix":           "如果不需要遠端操作，請關閉 VNC。需要使用時，請限制可連線的裝置並設定高強度密碼。",
	"router_upnp":            "路由器的 UPnP 已開啟。家中裝置可能會自動開啟通往網際網路的入口。",
	"router_upnp_fix":        "如果遊戲機等裝置不需要 UPnP，請在路由器設定中將它關閉。",
	"new_device":             "發現新裝置：%s",
	"new_device_fix":         "如果不認得這台裝置，請更改 Wi-Fi 密碼，並只重新連線自己的裝置。",
	"port_forward":           "%s 可透過 UPnP 開啟的連接埠 %d 從網際網路存取。",
	"port_forward_fix":       "如果不認得這項設定，請在路由器設定中刪除它或關閉 UPnP。遊戲主機和部分應用程式可能會刻意開啟這些連接埠。",
	"port_forward_risky":     "%s 的遠端存取或檔案分享透過連接埠 %d 對整個網際網路開放。攻擊者會不斷掃描這類入口。",
	"port_forward_risky_fix": "請在路由器設定中刪除此連接埠轉送並關閉 UPnP。如果需要遠端存取，請改用路由器的 VPN 功能。",
	"points_recovered":       "（+%d 分）",
	"press_enter":            "按 Enter 鍵即可關閉...",
}

func init() {
	messages["zh-Hant"] = catalogZhHant
}
