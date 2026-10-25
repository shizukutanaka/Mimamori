package i18n

var nodeREDMessages = map[string]map[string]string{
	"en": {
		"nodered_noauth":     "The Node-RED editor (port 1880) on %s has no password. Anyone on your network can change your flows and run commands on that device.",
		"nodered_noauth_fix": "Set adminAuth in Node-RED's settings.js to require a password, then restart Node-RED.",
	},
	"ja": {
		"nodered_noauth":     "「%s」の Node-RED の編集画面（1880番）に、パスワードがありません。家のネットワークの誰でも、フローを書き換えたり、その機器で命令を実行させたりできます。",
		"nodered_noauth_fix": "Node-RED の settings.js で adminAuth を設定してパスワードを付け、Node-RED を再起動してください。",
	},
	"zh-Hans": {
		"nodered_noauth":     "%s 上的 Node-RED 编辑器（端口 1880）没有密码。您网络中的任何人都可以更改您的流程并在该设备上运行命令。",
		"nodered_noauth_fix": "在 Node-RED 的 settings.js 中设置 adminAuth 以启用密码，然后重启 Node-RED。",
	},
	"zh-Hant": {
		"nodered_noauth":     "%s 上的 Node-RED 編輯器（連接埠 1880）沒有密碼。您網路中的任何人都可以更改您的流程，並在該裝置上執行命令。",
		"nodered_noauth_fix": "在 Node-RED 的 settings.js 中設定 adminAuth 以啟用密碼，然後重新啟動 Node-RED。",
	},
	"es": {
		"nodered_noauth":     "El editor de Node-RED en %s (puerto 1880) no tiene contraseña. Cualquier persona de su red puede cambiar sus flujos y ejecutar comandos en ese dispositivo.",
		"nodered_noauth_fix": "Configure adminAuth en settings.js de Node-RED para exigir una contraseña y reinicie Node-RED.",
	},
	"pt": {
		"nodered_noauth":     "O editor do Node-RED em %s (porta 1880) não tem senha. Qualquer pessoa na sua rede pode alterar os seus fluxos e executar comandos nesse dispositivo.",
		"nodered_noauth_fix": "Configure adminAuth no settings.js do Node-RED para exigir uma senha e reinicie o Node-RED.",
	},
	"fr": {
		"nodered_noauth":     "L’éditeur Node-RED sur %s (port 1880) n’a pas de mot de passe. Toute personne sur votre réseau peut modifier vos flux et exécuter des commandes sur cet appareil.",
		"nodered_noauth_fix": "Configurez adminAuth dans settings.js de Node-RED pour exiger un mot de passe, puis redémarrez Node-RED.",
	},
	"de": {
		"nodered_noauth":     "Der Node-RED-Editor auf %s (Port 1880) hat kein Passwort. Jede Person in Ihrem Netzwerk kann Ihre Flows ändern und Befehle auf diesem Gerät ausführen.",
		"nodered_noauth_fix": "Legen Sie adminAuth in Node-REDs settings.js fest, um ein Passwort zu verlangen, und starten Sie Node-RED neu.",
	},
	"ru": {
		"nodered_noauth":     "Редактор Node-RED на устройстве «%s» (порт 1880) не защищён паролем. Любой пользователь вашей сети может изменить ваши потоки и выполнять команды на этом устройстве.",
		"nodered_noauth_fix": "Настройте adminAuth в settings.js Node-RED, чтобы требовать пароль, затем перезапустите Node-RED.",
	},
	"ko": {
		"nodered_noauth":     "%s의 Node-RED 편집기(포트 1880)에 비밀번호가 없습니다. 네트워크의 누구나 플로우를 바꾸고 이 기기에서 명령을 실행할 수 있습니다.",
		"nodered_noauth_fix": "Node-RED의 settings.js에서 adminAuth를 설정해 비밀번호를 요구하고 Node-RED를 다시 시작하세요.",
	},
	"id": {
		"nodered_noauth":     "Editor Node-RED di %s (port 1880) tidak memiliki kata sandi. Siapa pun di jaringan Anda dapat mengubah flow Anda dan menjalankan perintah pada perangkat tersebut.",
		"nodered_noauth_fix": "Atur adminAuth di settings.js Node-RED agar meminta kata sandi, lalu mulai ulang Node-RED.",
	},
}

func init() {
	add(nodeREDMessages)
}
