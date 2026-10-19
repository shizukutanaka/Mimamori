package i18n

var homeAssistantMessages = map[string]map[string]string{
	"en": {
		"homeassistant_onboarding":     "Home Assistant (port 8123) on %s hasn't finished its first-time setup. Anyone on your network can create the owner account first and control all your smart-home devices.",
		"homeassistant_onboarding_fix": "Open Home Assistant and finish the first-time setup by creating the owner account. If you don't use it, turn it off.",
	},
	"ja": {
		"homeassistant_onboarding":     "「%s」の Home Assistant（8123番）は、最初の設定が終わっていません。家のネットワークの誰でも、先に持ち主のアカウントを作り、家の機器をすべて操作できる状態です。",
		"homeassistant_onboarding_fix": "Home Assistant を開いて、最初の設定（持ち主のアカウントの作成）を済ませてください。使っていないなら、止めてください。",
	},
	"zh-Hans": {
		"homeassistant_onboarding":     "%s 上的 Home Assistant（端口 8123）尚未完成首次设置。网络中的任何人都可以先创建所有者账户并控制家中的所有智能设备。",
		"homeassistant_onboarding_fix": "请打开 Home Assistant 并创建所有者账户以完成首次设置。如果您不使用它，请将其关闭。",
	},
	"zh-Hant": {
		"homeassistant_onboarding":     "%s 上的 Home Assistant（連接埠 8123）尚未完成首次設定。網路中的任何人都可以先建立擁有者帳戶並控制家中所有智慧家庭裝置。",
		"homeassistant_onboarding_fix": "請開啟 Home Assistant 並建立擁有者帳戶以完成首次設定。如果您不使用它，請將它關閉。",
	},
	"es": {
		"homeassistant_onboarding":     "Home Assistant (puerto 8123) en %s no ha completado la configuración inicial. Cualquier persona de su red puede crear primero la cuenta de propietario y controlar todos los dispositivos inteligentes de su hogar.",
		"homeassistant_onboarding_fix": "Abra Home Assistant y complete la configuración inicial creando la cuenta de propietario. Si no lo utiliza, apáguelo.",
	},
	"pt": {
		"homeassistant_onboarding":     "O Home Assistant (porta 8123) em %s não concluiu a configuração inicial. Qualquer pessoa na sua rede pode criar primeiro a conta do proprietário e controlar todos os dispositivos inteligentes da sua casa.",
		"homeassistant_onboarding_fix": "Abra o Home Assistant e conclua a configuração inicial criando a conta do proprietário. Se não o utiliza, desligue-o.",
	},
	"fr": {
		"homeassistant_onboarding":     "Home Assistant (port 8123) sur %s n’a pas terminé sa configuration initiale. Toute personne sur votre réseau peut créer le compte propriétaire en premier et contrôler tous les appareils domotiques de votre maison.",
		"homeassistant_onboarding_fix": "Ouvrez Home Assistant et terminez la configuration initiale en créant le compte propriétaire. Si vous ne l’utilisez pas, arrêtez-le.",
	},
	"de": {
		"homeassistant_onboarding":     "Home Assistant (Port 8123) auf %s hat die Ersteinrichtung nicht abgeschlossen. Jede Person in Ihrem Netzwerk kann zuerst das Besitzerkonto erstellen und alle Smart-Home-Geräte steuern.",
		"homeassistant_onboarding_fix": "Öffnen Sie Home Assistant und schließen Sie die Ersteinrichtung ab, indem Sie das Besitzerkonto erstellen. Wenn Sie Home Assistant nicht verwenden, schalten Sie es aus.",
	},
	"ru": {
		"homeassistant_onboarding":     "Home Assistant (порт 8123) на устройстве %s не завершил первоначальную настройку. Любой пользователь вашей сети может первым создать учётную запись владельца и управлять всеми устройствами умного дома.",
		"homeassistant_onboarding_fix": "Откройте Home Assistant и завершите первоначальную настройку, создав учётную запись владельца. Если вы им не пользуетесь, выключите его.",
	},
	"ko": {
		"homeassistant_onboarding":     "%s의 Home Assistant(포트 8123)가 최초 설정을 완료하지 않았습니다. 네트워크의 누구나 먼저 소유자 계정을 만들고 집의 모든 스마트홈 기기를 제어할 수 있습니다.",
		"homeassistant_onboarding_fix": "Home Assistant를 열어 소유자 계정을 만들고 최초 설정을 완료하세요. 사용하지 않는다면 종료하세요.",
	},
	"id": {
		"homeassistant_onboarding":     "Home Assistant (port 8123) pada %s belum menyelesaikan penyiapan awal. Siapa pun di jaringan Anda dapat membuat akun pemilik lebih dulu dan mengendalikan semua perangkat rumah pintar Anda.",
		"homeassistant_onboarding_fix": "Buka Home Assistant dan selesaikan penyiapan awal dengan membuat akun pemilik. Jika tidak digunakan, matikan.",
	},
}

var homeAssistantCatalogTargets = map[string]map[string]string{
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
	for lang, catalog := range homeAssistantMessages {
		target := homeAssistantCatalogTargets[lang]
		for key, value := range catalog {
			target[key] = value
		}
		messages[lang] = target
	}
}
