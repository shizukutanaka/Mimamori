package i18n

import "fmt"

var checkedMessages = map[string]map[string]string{
	"en": {
		"checked_devices":  "Checked: %s",
		"checked_ports":    "risky ports",
		"checked_forwards": "router UPnP port forwards",
		"checked_wifi":     "Wi-Fi encryption",
		"checked_sep":      " · ",
	},
	"ja": {
		"checked_devices":  "確認したこと: 機器 %s",
		"checked_ports":    "危険なポート",
		"checked_forwards": "ルーターの UPnP ポート転送",
		"checked_wifi":     "Wi-Fi の暗号化",
		"checked_sep":      "・",
	},
	"zh-Hans": {
		"checked_devices":  "已检查: %s",
		"checked_ports":    "危险端口",
		"checked_forwards": "路由器 UPnP 端口转发",
		"checked_wifi":     "Wi-Fi 加密",
		"checked_sep":      " · ",
	},
	"zh-Hant": {
		"checked_devices":  "已檢查: %s",
		"checked_ports":    "危險連接埠",
		"checked_forwards": "路由器 UPnP 連接埠轉送",
		"checked_wifi":     "Wi-Fi 加密",
		"checked_sep":      " · ",
	},
	"es": {
		"checked_devices":  "Revisado: %s",
		"checked_ports":    "puertos de riesgo",
		"checked_forwards": "redirecciones de puertos UPnP del router",
		"checked_wifi":     "cifrado de Wi-Fi",
		"checked_sep":      " · ",
	},
	"pt": {
		"checked_devices":  "Verificado: %s",
		"checked_ports":    "portas de risco",
		"checked_forwards": "redirecionamentos de portas UPnP do roteador",
		"checked_wifi":     "criptografia do Wi-Fi",
		"checked_sep":      " · ",
	},
	"fr": {
		"checked_devices":  "Vérifié : %s",
		"checked_ports":    "ports à risque",
		"checked_forwards": "redirections de ports UPnP du routeur",
		"checked_wifi":     "chiffrement du Wi-Fi",
		"checked_sep":      " · ",
	},
	"de": {
		"checked_devices":  "Geprüft: %s",
		"checked_ports":    "riskante Ports",
		"checked_forwards": "UPnP-Portweiterleitungen des Routers",
		"checked_wifi":     "WLAN-Verschlüsselung",
		"checked_sep":      " · ",
	},
	"ru": {
		"checked_devices":  "Проверено: %s",
		"checked_ports":    "опасные порты",
		"checked_forwards": "перенаправление портов UPnP на роутере",
		"checked_wifi":     "шифрование Wi-Fi",
		"checked_sep":      " · ",
	},
	"ko": {
		"checked_devices":  "확인한 항목: %s",
		"checked_ports":    "위험 포트",
		"checked_forwards": "라우터 UPnP 포트 포워딩",
		"checked_wifi":     "Wi-Fi 암호화",
		"checked_sep":      " · ",
	},
	"id": {
		"checked_devices":  "Diperiksa: %s",
		"checked_ports":    "port berisiko",
		"checked_forwards": "penerusan port UPnP router",
		"checked_wifi":     "enkripsi Wi-Fi",
		"checked_sep":      " · ",
	},
}

func CheckedDeviceCount(lang string, count int) string {
	switch lang {
	case "ja":
		return fmt.Sprintf("%d台", count)
	case "zh-Hans":
		return fmt.Sprintf("%d 台设备", count)
	case "zh-Hant":
		return fmt.Sprintf("%d 台裝置", count)
	case "es":
		if count == 1 {
			return "1 dispositivo"
		}
		return fmt.Sprintf("%d dispositivos", count)
	case "pt":
		if count == 1 {
			return "1 dispositivo"
		}
		return fmt.Sprintf("%d dispositivos", count)
	case "fr":
		if count == 1 {
			return "1 appareil"
		}
		return fmt.Sprintf("%d appareils", count)
	case "de":
		if count == 1 {
			return "1 Gerät"
		}
		return fmt.Sprintf("%d Geräte", count)
	case "ru":
		lastTwo := count % 100
		switch {
		case count%10 == 1 && lastTwo != 11:
			return fmt.Sprintf("%d устройство", count)
		case count%10 >= 2 && count%10 <= 4 && (lastTwo < 12 || lastTwo > 14):
			return fmt.Sprintf("%d устройства", count)
		default:
			return fmt.Sprintf("%d устройств", count)
		}
	case "ko":
		return fmt.Sprintf("기기 %d대", count)
	case "id":
		return fmt.Sprintf("%d perangkat", count)
	default:
		if count == 1 {
			return "1 device"
		}
		return fmt.Sprintf("%d devices", count)
	}
}

func init() {
	for lang, catalog := range checkedMessages {
		target := wifiCatalogTargets[lang]
		for key, value := range catalog {
			target[key] = value
		}
		messages[lang] = target
	}
}
