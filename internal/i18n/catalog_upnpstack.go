package i18n

var upnpStackMessages = map[string]map[string]string{
	"en": {
		"upnp_stack_vulnerable":     "%s uses an old UPnP component (%s) with known remotely exploitable vulnerabilities.",
		"upnp_stack_vulnerable_fix": "Update the device's firmware. If no update exists, turn off UPnP on it (and on the router), and consider replacing it.",
	},
	"ja": {
		"upnp_stack_vulnerable":     "%s は、遠隔から悪用できる既知の脆弱性がある古い UPnP 部品（%s）を使っています。",
		"upnp_stack_vulnerable_fix": "ファームウェアを更新してください。更新がない場合は、その機器とルーターの UPnP を無効にし、買い替えも検討してください。",
	},
	"zh-Hans": {
		"upnp_stack_vulnerable":     "设备 %s 使用的旧版 UPnP 组件（%s）存在已知的可被远程利用的漏洞。",
		"upnp_stack_vulnerable_fix": "请您更新设备的固件。若没有可用更新，请关闭该设备和路由器上的 UPnP，并考虑更换设备。",
	},
	"zh-Hant": {
		"upnp_stack_vulnerable":     "裝置 %s 使用的舊版 UPnP 元件（%s）有已知的遠端可利用漏洞。",
		"upnp_stack_vulnerable_fix": "請您更新裝置的韌體。若沒有可用更新，請關閉該裝置和路由器上的 UPnP，並考慮更換裝置。",
	},
	"es": {
		"upnp_stack_vulnerable":     "El dispositivo %s utiliza un componente UPnP antiguo (%s) con vulnerabilidades conocidas que pueden explotarse de forma remota.",
		"upnp_stack_vulnerable_fix": "Actualice el firmware del dispositivo. Si no hay ninguna actualización, desactive UPnP en el dispositivo y en el router, y considere sustituirlo.",
	},
	"pt": {
		"upnp_stack_vulnerable":     "O dispositivo %s usa um componente UPnP antigo (%s) com vulnerabilidades conhecidas que podem ser exploradas remotamente.",
		"upnp_stack_vulnerable_fix": "Atualize o firmware do dispositivo. Se não houver atualização disponível, desative o UPnP no dispositivo e no roteador e considere substituí-lo.",
	},
	"fr": {
		"upnp_stack_vulnerable":     "L’appareil %s utilise un ancien composant UPnP (%s) présentant des vulnérabilités connues exploitables à distance.",
		"upnp_stack_vulnerable_fix": "Mettez à jour le micrologiciel de l’appareil. Si aucune mise à jour n’est disponible, désactivez UPnP sur l’appareil et le routeur, et envisagez de le remplacer.",
	},
	"de": {
		"upnp_stack_vulnerable":     "Das Gerät %s verwendet eine veraltete UPnP-Komponente (%s) mit bekannten, aus der Ferne ausnutzbaren Sicherheitslücken.",
		"upnp_stack_vulnerable_fix": "Aktualisieren Sie die Firmware des Geräts. Wenn kein Update verfügbar ist, schalten Sie UPnP am Gerät und am Router aus und erwägen Sie, das Gerät zu ersetzen.",
	},
	"ru": {
		"upnp_stack_vulnerable":     "Устройство %s использует устаревший компонент UPnP (%s) с известными уязвимостями, которые можно использовать удалённо.",
		"upnp_stack_vulnerable_fix": "Обновите прошивку устройства. Если обновления нет, отключите UPnP на устройстве и роутере и подумайте о замене устройства.",
	},
	"ko": {
		"upnp_stack_vulnerable":     "기기 %s에서 원격으로 악용할 수 있는 알려진 취약점이 있는 오래된 UPnP 구성 요소(%s)를 사용하고 있습니다.",
		"upnp_stack_vulnerable_fix": "기기의 펌웨어를 업데이트하세요. 업데이트가 없다면 해당 기기와 라우터에서 UPnP를 끄고 기기 교체를 고려하세요.",
	},
	"id": {
		"upnp_stack_vulnerable":     "Perangkat %s menggunakan komponen UPnP lama (%s) dengan kerentanan yang diketahui dapat dieksploitasi dari jarak jauh.",
		"upnp_stack_vulnerable_fix": "Perbarui firmware perangkat. Jika Anda tidak menemukan pembaruan, nonaktifkan UPnP pada perangkat dan router, lalu pertimbangkan untuk menggantinya.",
	},
}

func init() {
	for lang, catalog := range upnpStackMessages {
		target := wifiCatalogTargets[lang]
		for key, value := range catalog {
			target[key] = value
		}
		messages[lang] = target
	}
}
