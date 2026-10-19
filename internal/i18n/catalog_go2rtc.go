package i18n

var go2rtcMessages = map[string]map[string]string{
	"en": {
		"go2rtc_noauth":     "The go2rtc web interface on %s (port 1984) has no password. Anyone on your network can change its settings and run programs on that device.",
		"go2rtc_noauth_fix": "In go2rtc's settings (go2rtc.yaml), set a username and password under api, or set api listen to \"127.0.0.1:1984\" if you only use it on that device, then restart go2rtc.",
	},
	"ja": {
		"go2rtc_noauth":     "「%s」の go2rtc の管理画面（1984番）に、パスワードがありません。家のネットワークの誰でも、設定を変えたり、その機器で命令を実行させたりできます。",
		"go2rtc_noauth_fix": "go2rtc の設定（go2rtc.yaml）の api に username と password を設定するか、その機器の中だけで使うなら listen を \"127.0.0.1:1984\" にして、go2rtc を再起動してください。",
	},
	"zh-Hans": {
		"go2rtc_noauth":     "go2rtc 在 %s 的 Web 管理界面（端口 1984）没有密码。您网络中的任何人都可以更改设置并在该设备上运行程序。",
		"go2rtc_noauth_fix": "在 go2rtc 的设置（go2rtc.yaml）中，在 api 下设置 username 和 password；如果只在该设备上使用，请将 api listen 设为 \"127.0.0.1:1984\"，然后重启 go2rtc。",
	},
	"zh-Hant": {
		"go2rtc_noauth":     "go2rtc 在 %s 的網頁管理介面（連接埠 1984）沒有密碼。您網路中的任何人都可以更改設定，並在該裝置上執行程式。",
		"go2rtc_noauth_fix": "在 go2rtc 的設定（go2rtc.yaml）中，於 api 下設定 username 和 password；如果只在該裝置上使用，請將 api listen 設為 \"127.0.0.1:1984\"，然後重新啟動 go2rtc。",
	},
	"es": {
		"go2rtc_noauth":     "La interfaz web de go2rtc en %s (puerto 1984) no tiene contraseña. Cualquier persona de su red puede cambiar su configuración y ejecutar programas en ese dispositivo.",
		"go2rtc_noauth_fix": "En la configuración de go2rtc (go2rtc.yaml), establezca username y password en api o configure api listen como \"127.0.0.1:1984\" si solo lo usa en ese dispositivo; después, reinicie go2rtc.",
	},
	"pt": {
		"go2rtc_noauth":     "A interface web do go2rtc em %s (porta 1984) não tem senha. Qualquer pessoa na sua rede pode alterar as configurações e executar programas nesse dispositivo.",
		"go2rtc_noauth_fix": "Nas configurações do go2rtc (go2rtc.yaml), defina username e password em api ou configure api listen como \"127.0.0.1:1984\" se você só o usa nesse dispositivo; depois, reinicie o go2rtc.",
	},
	"fr": {
		"go2rtc_noauth":     "L’interface web de go2rtc sur %s (port 1984) n’a pas de mot de passe. Toute personne sur votre réseau peut modifier ses paramètres et exécuter des programmes sur cet appareil.",
		"go2rtc_noauth_fix": "Dans les paramètres de go2rtc (go2rtc.yaml), définissez username et password sous api, ou définissez api listen sur \"127.0.0.1:1984\" si vous l’utilisez uniquement sur cet appareil, puis redémarrez go2rtc.",
	},
	"de": {
		"go2rtc_noauth":     "Die go2rtc-Weboberfläche auf %s (Port 1984) hat kein Passwort. Jede Person in Ihrem Netzwerk kann die Einstellungen ändern und Programme auf diesem Gerät ausführen.",
		"go2rtc_noauth_fix": "Legen Sie in den go2rtc-Einstellungen (go2rtc.yaml) unter api einen username und ein password fest oder setzen Sie api listen auf \"127.0.0.1:1984\", wenn Sie go2rtc nur auf diesem Gerät verwenden, und starten Sie go2rtc neu.",
	},
	"ru": {
		"go2rtc_noauth":     "Веб-интерфейс go2rtc на устройстве «%s» (порт 1984) не защищён паролем. Любой пользователь вашей сети может изменить настройки и запускать программы на этом устройстве.",
		"go2rtc_noauth_fix": "В настройках go2rtc (go2rtc.yaml) задайте username и password в разделе api или установите api listen в значение \"127.0.0.1:1984\", если используете его только на этом устройстве, а затем перезапустите go2rtc.",
	},
	"ko": {
		"go2rtc_noauth":     "%s의 go2rtc 웹 인터페이스(포트 1984)에 비밀번호가 없습니다. 네트워크의 누구나 설정을 변경하고 이 기기에서 프로그램을 실행할 수 있습니다.",
		"go2rtc_noauth_fix": "go2rtc 설정(go2rtc.yaml)의 api 아래에 username과 password를 설정하거나, 이 기기에서만 사용한다면 api listen을 \"127.0.0.1:1984\"로 설정한 뒤 go2rtc를 다시 시작하세요.",
	},
	"id": {
		"go2rtc_noauth":     "Antarmuka web go2rtc di %s (port 1984) tidak memiliki kata sandi. Siapa pun di jaringan Anda dapat mengubah pengaturannya dan menjalankan program pada perangkat tersebut.",
		"go2rtc_noauth_fix": "Di pengaturan go2rtc (go2rtc.yaml), tetapkan username dan password di bawah api, atau atur api listen ke \"127.0.0.1:1984\" jika Anda hanya menggunakannya di perangkat tersebut, lalu mulai ulang go2rtc.",
	},
}

var go2rtcCatalogTargets = map[string]map[string]string{
	"en":      messages["en"],
	"ja":      messages["ja"],
	"zh-Hans": catalogZhHans,
	"zh-Hant": catalogZhHant,
	"es":      catalogEs,
	"pt":      catalogPt,
	"fr":      catalogFr,
	"de":      catalogDe,
	"ru":      catalogRu,
	"ko":      catalogKo,
	"id":      catalogId,
}

func init() {
	for lang, catalog := range go2rtcMessages {
		target := go2rtcCatalogTargets[lang]
		for key, value := range catalog {
			target[key] = value
		}
		messages[lang] = target
	}
}
