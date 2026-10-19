package i18n

var scryptedMessages = map[string]map[string]string{
	"en": {
		"scrypted_onboarding":     "Scrypted (port 11080) on %s hasn't had its first account created yet. Anyone on your network can create the admin account first, then install plugins and run commands.",
		"scrypted_onboarding_fix": "Open Scrypted (https://<this device>:10443) and create the admin account now. If you don't use it, turn it off.",
	},
	"ja": {
		"scrypted_onboarding":     "「%s」の Scrypted（11080番）は、最初のアカウントがまだ作られていません。家のネットワークの誰でも、先に管理者のアカウントを作り、プラグインを入れて命令を実行できる状態です。",
		"scrypted_onboarding_fix": "Scrypted（https://この機器:10443）を開いて、いますぐ管理者のアカウントを作ってください。使っていなければ、止めてください。",
	},
	"zh-Hans": {
		"scrypted_onboarding":     "%s 上的 Scrypted（端口 11080）尚未创建第一个账户。您网络中的任何人都可以抢先创建管理员账户，然后安装插件并运行命令。",
		"scrypted_onboarding_fix": "请立即打开 Scrypted（https://<此设备>:10443）并创建管理员账户。如果您不使用它，请将其关闭。",
	},
	"zh-Hant": {
		"scrypted_onboarding":     "%s 上的 Scrypted（連接埠 11080）尚未建立第一個帳戶。您網路中的任何人都可以搶先建立管理員帳戶，然後安裝外掛程式並執行命令。",
		"scrypted_onboarding_fix": "請立即開啟 Scrypted（https://<此裝置>:10443）並建立管理員帳戶。如果您不使用它，請將它關閉。",
	},
	"es": {
		"scrypted_onboarding":     "Aún no se ha creado la primera cuenta de Scrypted (puerto 11080) en %s. Cualquiera de su red puede crear primero la cuenta de administrador, instalar complementos y ejecutar comandos.",
		"scrypted_onboarding_fix": "Abra Scrypted (https://<este dispositivo>:10443) y cree ahora la cuenta de administrador. Si no lo usa, apáguelo.",
	},
	"pt": {
		"scrypted_onboarding":     "A primeira conta do Scrypted (porta 11080) ainda não foi criada em %s. Qualquer pessoa da sua rede pode criar primeiro a conta de administrador, instalar plugins e executar comandos.",
		"scrypted_onboarding_fix": "Abra o Scrypted (https://<este dispositivo>:10443) e crie a conta de administrador agora. Se não o utiliza, desligue-o.",
	},
	"fr": {
		"scrypted_onboarding":     "Le premier compte Scrypted (port 11080) n’a pas encore été créé sur %s. Toute personne de votre réseau peut créer le compte administrateur en premier, puis installer des extensions et exécuter des commandes.",
		"scrypted_onboarding_fix": "Ouvrez Scrypted (https://<cet appareil>:10443) et créez le compte administrateur dès maintenant. Si vous ne l’utilisez pas, désactivez-le.",
	},
	"de": {
		"scrypted_onboarding":     "Auf %s wurde noch kein erstes Konto für Scrypted (Port 11080) erstellt. Jede Person in Ihrem Netzwerk kann zuerst das Administratorkonto anlegen, Plugins installieren und Befehle ausführen.",
		"scrypted_onboarding_fix": "Öffnen Sie Scrypted (https://<dieses Gerät>:10443) und erstellen Sie jetzt das Administratorkonto. Wenn Sie es nicht nutzen, schalten Sie es aus.",
	},
	"ru": {
		"scrypted_onboarding":     "На устройстве %s ещё не создана первая учётная запись Scrypted (порт 11080). Любой пользователь вашей сети может первым создать учётную запись администратора, затем устанавливать плагины и выполнять команды.",
		"scrypted_onboarding_fix": "Откройте Scrypted (https://<это устройство>:10443) и сейчас же создайте учётную запись администратора. Если вы не пользуетесь им, выключите его.",
	},
	"ko": {
		"scrypted_onboarding":     "%s의 Scrypted(포트 11080)에 첫 계정이 아직 만들어지지 않았습니다. 네트워크의 누구나 먼저 관리자 계정을 만들고 플러그인을 설치하거나 명령을 실행할 수 있습니다.",
		"scrypted_onboarding_fix": "Scrypted(https://<이 기기>:10443)를 열어 지금 관리자 계정을 만드세요. 사용하지 않는다면 끄세요.",
	},
	"id": {
		"scrypted_onboarding":     "Akun pertama Scrypted (port 11080) di %s belum dibuat. Siapa pun di jaringan Anda dapat lebih dahulu membuat akun administrator, lalu memasang plugin dan menjalankan perintah.",
		"scrypted_onboarding_fix": "Buka Scrypted (https://<perangkat ini>:10443) dan buat akun administrator sekarang. Jika tidak digunakan, matikan.",
	},
}

var scryptedCatalogTargets = map[string]map[string]string{
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
	for lang, catalog := range scryptedMessages {
		target := scryptedCatalogTargets[lang]
		for key, value := range catalog {
			target[key] = value
		}
		messages[lang] = target
	}
}
