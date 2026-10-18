package i18n

var mqttMessages = map[string]map[string]string{
	"en": {
		"mqtt_anonymous":     "The smart-home hub (MQTT) on %s lets anyone on your network control devices without a password.",
		"mqtt_anonymous_fix": "In the MQTT broker's settings, turn off anonymous access (for Mosquitto: allow_anonymous false) and set a password, then update it.",
	},
	"ja": {
		"mqtt_anonymous":     "「%s」のスマートホーム用の MQTT が、パスワードなしで家の中の誰でも機器を操作できる状態です。",
		"mqtt_anonymous_fix": "MQTT ブローカーの設定で匿名接続を無効にし（Mosquitto なら allow_anonymous false）、パスワードを設定して、ソフトを更新してください。",
	},
	"zh-Hans": {
		"mqtt_anonymous":     "%s 上的智能家居中枢（MQTT）允许网络中的任何人无密码控制设备。",
		"mqtt_anonymous_fix": "请您在 MQTT 代理设置中关闭匿名访问（对于 Mosquitto：allow_anonymous false），设置密码，然后更新软件。",
	},
	"zh-Hant": {
		"mqtt_anonymous":     "%s 上的智慧家庭中樞（MQTT）允許網路中的任何人無密碼控制裝置。",
		"mqtt_anonymous_fix": "請您在 MQTT 代理程式設定中關閉匿名存取（Mosquitto：allow_anonymous false），設定密碼，然後更新軟體。",
	},
	"es": {
		"mqtt_anonymous":     "El concentrador de hogar inteligente (MQTT) en %s permite que cualquiera en su red controle dispositivos sin contraseña.",
		"mqtt_anonymous_fix": "Desactive el acceso anónimo en la configuración del bróker MQTT (en Mosquitto: allow_anonymous false), establezca una contraseña y actualícelo.",
	},
	"pt": {
		"mqtt_anonymous":     "O hub de casa inteligente (MQTT) em %s permite que qualquer pessoa na sua rede controle dispositivos sem senha.",
		"mqtt_anonymous_fix": "Desative o acesso anônimo nas configurações do broker MQTT (no Mosquitto: allow_anonymous false), defina uma senha e atualize-o.",
	},
	"fr": {
		"mqtt_anonymous":     "Le concentrateur domotique (MQTT) sur %s permet à toute personne de votre réseau de contrôler des appareils sans mot de passe.",
		"mqtt_anonymous_fix": "Désactivez l'accès anonyme dans les paramètres du broker MQTT (pour Mosquitto : allow_anonymous false), définissez un mot de passe, puis mettez-le à jour.",
	},
	"de": {
		"mqtt_anonymous":     "Der Smart-Home-Hub (MQTT) auf %s ermöglicht es jedem in Ihrem Netzwerk, Geräte ohne Passwort zu steuern.",
		"mqtt_anonymous_fix": "Deaktivieren Sie den anonymen Zugriff in den Einstellungen des MQTT-Brokers (bei Mosquitto: allow_anonymous false), legen Sie ein Passwort fest und aktualisieren Sie ihn.",
	},
	"ru": {
		"mqtt_anonymous":     "Центр умного дома (MQTT) на %s позволяет любому в вашей сети управлять устройствами без пароля.",
		"mqtt_anonymous_fix": "Отключите анонимный доступ в настройках MQTT-брокера (для Mosquitto: allow_anonymous false), задайте пароль и обновите его.",
	},
	"ko": {
		"mqtt_anonymous":     "%s의 스마트홈 허브(MQTT)는 네트워크의 누구나 비밀번호 없이 기기를 제어할 수 있게 합니다.",
		"mqtt_anonymous_fix": "MQTT 브로커 설정에서 익명 접속을 끄고(Mosquitto의 경우 allow_anonymous false), 비밀번호를 설정한 다음 업데이트하세요.",
	},
	"id": {
		"mqtt_anonymous":     "Hub rumah pintar (MQTT) di %s memungkinkan siapa pun di jaringan Anda mengendalikan perangkat tanpa kata sandi.",
		"mqtt_anonymous_fix": "Di pengaturan broker MQTT, nonaktifkan akses anonim (untuk Mosquitto: allow_anonymous false), tetapkan kata sandi, lalu perbarui perangkat lunaknya.",
	},
}

var mqttCatalogTargets = map[string]map[string]string{
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
	for lang, catalog := range mqttMessages {
		target := mqttCatalogTargets[lang]
		for key, value := range catalog {
			target[key] = value
		}
		messages[lang] = target
	}
}
