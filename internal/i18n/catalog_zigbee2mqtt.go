package i18n

var zigbee2mqttMessages = map[string]map[string]string{
	"en": {
		"zigbee2mqtt_noauth":     "The Zigbee2MQTT web interface on %s lets anyone on your network control your Zigbee devices and run programs without a password.",
		"zigbee2mqtt_noauth_fix": "In Zigbee2MQTT's configuration.yaml, set frontend.auth_token to a long random value, then restart Zigbee2MQTT.",
	},
	"ja": {
		"zigbee2mqtt_noauth":     "「%s」の Zigbee2MQTT の管理画面が、パスワードなしで家の中の誰でも Zigbee 機器を操作したり、プログラムを実行させたりできる状態です。",
		"zigbee2mqtt_noauth_fix": "Zigbee2MQTT の configuration.yaml で frontend.auth_token に長いランダムな値を設定し、Zigbee2MQTT を再起動してください。",
	},
	"zh-Hans": {
		"zigbee2mqtt_noauth":     "Zigbee2MQTT 在 %s 上的网页界面没有密码，您网络中的任何人都可以控制 Zigbee 设备并运行程序。",
		"zigbee2mqtt_noauth_fix": "请在 Zigbee2MQTT 的 configuration.yaml 中将 frontend.auth_token 设置为较长的随机值，然后重新启动 Zigbee2MQTT。",
	},
	"zh-Hant": {
		"zigbee2mqtt_noauth":     "Zigbee2MQTT 在 %s 上的網頁介面沒有密碼，您網路中的任何人都可以控制 Zigbee 裝置並執行程式。",
		"zigbee2mqtt_noauth_fix": "請在 Zigbee2MQTT 的 configuration.yaml 中將 frontend.auth_token 設定為較長的隨機值，然後重新啟動 Zigbee2MQTT。",
	},
	"es": {
		"zigbee2mqtt_noauth":     "La interfaz web de Zigbee2MQTT en %s permite que cualquier persona de su red controle sus dispositivos Zigbee y ejecute programas sin contraseña.",
		"zigbee2mqtt_noauth_fix": "En el archivo configuration.yaml de Zigbee2MQTT, establezca frontend.auth_token con un valor aleatorio largo y reinicie Zigbee2MQTT.",
	},
	"pt": {
		"zigbee2mqtt_noauth":     "A interface web do Zigbee2MQTT em %s permite que qualquer pessoa na sua rede controle os seus dispositivos Zigbee e execute programas sem senha.",
		"zigbee2mqtt_noauth_fix": "No arquivo configuration.yaml do Zigbee2MQTT, defina frontend.auth_token com um valor aleatório longo e reinicie o Zigbee2MQTT.",
	},
	"fr": {
		"zigbee2mqtt_noauth":     "L’interface web de Zigbee2MQTT sur %s permet à toute personne de votre réseau de contrôler vos appareils Zigbee et d’exécuter des programmes sans mot de passe.",
		"zigbee2mqtt_noauth_fix": "Dans le fichier configuration.yaml de Zigbee2MQTT, définissez frontend.auth_token sur une longue valeur aléatoire, puis redémarrez Zigbee2MQTT.",
	},
	"de": {
		"zigbee2mqtt_noauth":     "Die Zigbee2MQTT-Weboberfläche auf %s ist nicht durch ein Passwort geschützt. Jede Person in Ihrem Netzwerk kann Ihre Zigbee-Geräte steuern und Programme ausführen.",
		"zigbee2mqtt_noauth_fix": "Legen Sie in der configuration.yaml von Zigbee2MQTT für frontend.auth_token einen langen Zufallswert fest und starten Sie Zigbee2MQTT anschließend neu.",
	},
	"ru": {
		"zigbee2mqtt_noauth":     "Веб-интерфейс Zigbee2MQTT на %s не защищён паролем. Любой пользователь вашей сети может управлять устройствами Zigbee и запускать программы.",
		"zigbee2mqtt_noauth_fix": "В configuration.yaml Zigbee2MQTT задайте для frontend.auth_token длинное случайное значение, затем перезапустите Zigbee2MQTT.",
	},
	"ko": {
		"zigbee2mqtt_noauth":     "%s의 Zigbee2MQTT 웹 인터페이스에 비밀번호가 설정되어 있지 않습니다. 네트워크의 누구나 Zigbee 기기를 제어하고 프로그램을 실행할 수 있습니다.",
		"zigbee2mqtt_noauth_fix": "Zigbee2MQTT의 configuration.yaml에서 frontend.auth_token을 길고 무작위인 값으로 설정한 후 Zigbee2MQTT를 다시 시작하세요.",
	},
	"id": {
		"zigbee2mqtt_noauth":     "Antarmuka web Zigbee2MQTT di %s tidak dilindungi kata sandi. Siapa pun di jaringan Anda dapat mengendalikan perangkat Zigbee dan menjalankan program.",
		"zigbee2mqtt_noauth_fix": "Di configuration.yaml Zigbee2MQTT, tetapkan frontend.auth_token ke nilai acak yang panjang, lalu mulai ulang Zigbee2MQTT.",
	},
}

var zigbee2mqttCatalogTargets = map[string]map[string]string{
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
	for lang, catalog := range zigbee2mqttMessages {
		target := zigbee2mqttCatalogTargets[lang]
		for key, value := range catalog {
			target[key] = value
		}
		messages[lang] = target
	}
}
