package i18n

var catalogKo = map[string]string{
	"title":             "Mimamori — 가정용 네트워크 안전 점검",
	"safe":              "안전",
	"caution":           "주의",
	"danger":            "위험",
	"score":             "가정용 네트워크 안전: %s  %d점",
	"device_counts":     "기기: %d대 (%d대 신규)",
	"findings_heading":  "확인하고 조치할 내용",
	"devices_heading":   "기기",
	"no_problems":       "문제가 발견되지 않았습니다",
	"first_run":         "이번 실행에서 네트워크의 기기를 기록했습니다. 새 기기는 다음 실행부터 알려드립니다.",
	"new_badge":         "신규",
	"router":            "공유기",
	"randomized_device": "휴대폰 또는 PC (비공개 Wi-Fi 주소)",
	"telnet_open":       "%s에서 Telnet이 열려 있습니다. 오래된 원격 제어 서비스는 기기를 탈취하는 데 악용될 수 있습니다.",
	"telnet_open_fix":   "기기 설정에서 Telnet을 끄고 사용 가능한 업데이트를 설치하세요. 기본 비밀번호도 변경하세요.",
	"ftp_open":          "%s에서 FTP 파일 접근이 열려 있어 파일이나 계정 정보가 노출될 수 있습니다.",
	"ftp_open_fix":      "필요하지 않다면 FTP를 끄세요. 필요하다면 더 안전한 방법과 고유한 비밀번호를 사용하세요.",
	"vnc_open":          "%s에서 VNC 원격 접속이 열려 있습니다. 접속 가능한 사람이 화면을 조작할 수 있습니다.",
	"vnc_open_fix":      "원격 접속이 필요하지 않다면 VNC를 끄세요. 필요하다면 접속 가능한 기기를 제한하고 강력한 비밀번호를 설정하세요.",
	"router_upnp":       "공유기에서 UPnP가 켜져 있습니다. 집 안의 기기가 인터넷에 연결되는 통로를 자동으로 열 수 있습니다.",
	"router_upnp_fix":   "게임기 등에서 필요하지 않다면 공유기 설정에서 UPnP를 끄세요.",
	"new_device":        "새 기기가 연결되었습니다: %s",
	"new_device_fix":    "모르는 기기라면 Wi-Fi 비밀번호를 바꾸고 본인 기기만 다시 연결하세요.",
	"points_recovered":  "(+%d점)",
	"press_enter":       "Enter 키를 누르면 닫힙니다...",
}

func init() {
	messages["ko"] = catalogKo
}
