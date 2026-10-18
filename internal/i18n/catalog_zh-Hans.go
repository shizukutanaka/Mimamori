package i18n

var catalogZhHans = map[string]string{
	"title":                  "Mimamori — 家庭网络安全检查",
	"safe":                   "安全",
	"caution":                "注意",
	"danger":                 "危险",
	"score":                  "家庭网络安全：%s  %d 分",
	"device_counts":          "设备：%d 台（新增 %d 台）",
	"findings_heading":       "需要处理的问题",
	"devices_heading":        "设备",
	"no_problems":            "未发现问题",
	"first_run":              "本次已记录家庭网络中的设备。从下次运行起会报告新设备。",
	"new_badge":              "新",
	"router":                 "路由器",
	"randomized_device":      "手机或电脑（使用 Wi-Fi 私有地址）",
	"telnet_open":            "“%s”上的 Telnet 已开启。这种旧式远程控制服务容易被用来控制设备。",
	"telnet_open_fix":        "请在设备设置中关闭 Telnet 并安装可用的更新。如果仍在使用初始密码，请一并更改。",
	"ftp_open":               "“%s”上的 FTP 文件访问已开启，文件或账户信息可能因此暴露。",
	"ftp_open_fix":           "如果不需要，请关闭 FTP。如果需要，请改用更安全的方式并设置独一无二的密码。",
	"vnc_open":               "“%s”上的 VNC 远程操作已开启，能够连接的人可能会控制屏幕。",
	"vnc_open_fix":           "如果不需要远程操作，请关闭 VNC。需要使用时，请限制可连接的设备并设置强密码。",
	"router_upnp":            "路由器的 UPnP 已开启。家中设备可能会自动打开通往互联网的入口。",
	"router_upnp_fix":        "如果游戏机等设备不需要 UPnP，请在路由器设置中将其关闭。",
	"new_device":             "发现新设备：%s",
	"new_device_fix":         "如果不认识这台设备，请更改 Wi-Fi 密码，并只重新连接自己的设备。",
	"port_forward":           "%s 可通过 UPnP 打开的端口 %d 从互联网访问。",
	"port_forward_fix":       "如果您不认识这项设置，请在路由器设置中删除它或关闭 UPnP。游戏主机和部分应用可能会有意打开这些端口。",
	"port_forward_risky":     "%s 的远程访问或文件共享通过端口 %d 向整个互联网开放。攻击者会持续扫描这类入口。",
	"port_forward_risky_fix": "请在路由器设置中删除此端口转发并关闭 UPnP。如果需要远程访问，请改用路由器的 VPN 功能。",
	"points_recovered":       "（+%d 分）",
	"press_enter":            "按 Enter 键关闭...",
}

func init() {
	messages["zh-Hans"] = catalogZhHans
}
