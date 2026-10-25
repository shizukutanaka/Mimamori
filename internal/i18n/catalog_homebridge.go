package i18n

var homebridgeMessages = map[string]map[string]string{
	"en": {
		"homebridge_noauth":     "The Homebridge UI on %s has login turned off, so anyone on your network can sign in as an administrator, install plugins, and run commands.",
		"homebridge_noauth_fix": "In Homebridge's config.json, remove \"auth\": \"none\" from the \"platform\": \"config\" section (or change it to \"form\"), then restart Homebridge.",
	},
	"ja": {
		"homebridge_noauth":     "「%s」の Homebridge の管理画面でログインが切られていて、家の中の誰でも管理者として入り、プラグインを入れたり命令を実行したりできる状態です。",
		"homebridge_noauth_fix": "Homebridge の config.json で、\"platform\": \"config\" の部分にある \"auth\": \"none\" を消すか \"form\" に変えて、Homebridge を再起動してください。",
	},
	"zh-Hans": {
		"homebridge_noauth":     "%s 上的 Homebridge 管理界面已关闭登录，因此您网络中的任何人都可以以管理员身份登录、安装插件并运行命令。",
		"homebridge_noauth_fix": "请您在 Homebridge 的 config.json 中，从 \"platform\": \"config\" 部分移除 \"auth\": \"none\"（或将其改为 \"form\"），然后重启 Homebridge。",
	},
	"zh-Hant": {
		"homebridge_noauth":     "%s 上的 Homebridge 管理介面已關閉登入，因此您網路中的任何人都可以以管理員身分登入、安裝外掛程式並執行命令。",
		"homebridge_noauth_fix": "請您在 Homebridge 的 config.json 中，從 \"platform\": \"config\" 部分移除 \"auth\": \"none\"（或將其改為 \"form\"），然後重新啟動 Homebridge。",
	},
	"es": {
		"homebridge_noauth":     "La interfaz de Homebridge en %s tiene el inicio de sesión desactivado, por lo que cualquiera de su red puede entrar como administrador, instalar complementos y ejecutar comandos.",
		"homebridge_noauth_fix": "En el config.json de Homebridge, elimine \"auth\": \"none\" de la sección \"platform\": \"config\" (o cámbielo por \"form\") y reinicie Homebridge.",
	},
	"pt": {
		"homebridge_noauth":     "A interface do Homebridge em %s está com o início de sessão desativado, por isso qualquer pessoa da sua rede pode entrar como administrador, instalar plugins e executar comandos.",
		"homebridge_noauth_fix": "No config.json do Homebridge, remova \"auth\": \"none\" da secção \"platform\": \"config\" (ou altere para \"form\") e reinicie o Homebridge.",
	},
	"fr": {
		"homebridge_noauth":     "L’interface Homebridge sur %s a la connexion désactivée ; toute personne de votre réseau peut donc se connecter comme administrateur, installer des extensions et exécuter des commandes.",
		"homebridge_noauth_fix": "Dans le config.json de Homebridge, supprimez \"auth\": \"none\" de la section \"platform\": \"config\" (ou remplacez-le par \"form\"), puis redémarrez Homebridge.",
	},
	"de": {
		"homebridge_noauth":     "Bei der Homebridge-Oberfläche auf %s ist die Anmeldung deaktiviert. Daher kann sich jede Person in Ihrem Netzwerk als Administrator anmelden, Plugins installieren und Befehle ausführen.",
		"homebridge_noauth_fix": "Entfernen Sie in der config.json von Homebridge \"auth\": \"none\" aus dem Abschnitt \"platform\": \"config\" (oder ändern Sie es in \"form\"), und starten Sie Homebridge neu.",
	},
	"ru": {
		"homebridge_noauth":     "На интерфейсе Homebridge на %s отключён вход, поэтому любой пользователь вашей сети может войти как администратор, установить плагины и выполнять команды.",
		"homebridge_noauth_fix": "В config.json Homebridge удалите \"auth\": \"none\" из раздела \"platform\": \"config\" (или замените на \"form\"), затем перезапустите Homebridge.",
	},
	"ko": {
		"homebridge_noauth":     "%s의 Homebridge 관리 화면에서 로그인이 꺼져 있어 네트워크의 누구나 관리자로 로그인하고 플러그인을 설치하거나 명령을 실행할 수 있습니다.",
		"homebridge_noauth_fix": "Homebridge의 config.json에서 \"platform\": \"config\" 부분의 \"auth\": \"none\"을 삭제하거나 \"form\"으로 바꾼 뒤 Homebridge를 다시 시작하세요.",
	},
	"id": {
		"homebridge_noauth":     "Login pada antarmuka Homebridge di %s dinonaktifkan, sehingga siapa pun di jaringan Anda dapat masuk sebagai administrator, memasang plugin, dan menjalankan perintah.",
		"homebridge_noauth_fix": "Silakan hapus \"auth\": \"none\" dari bagian \"platform\": \"config\" di config.json Homebridge (atau ubah menjadi \"form\"), lalu mulai ulang Homebridge.",
	},
}

func init() {
	add(homebridgeMessages)
}
