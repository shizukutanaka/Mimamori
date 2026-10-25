package i18n

var jupyterMessages = map[string]map[string]string{
	"en": {
		"jupyter_noauth":     "Jupyter (port 8888) on %s is open to your network with no token or password. Anyone at home can run any program on this computer.",
		"jupyter_noauth_fix": "Start Jupyter with a token or password (don't launch it with --IdentityProvider.token='' or --NotebookApp.token='').",
	},
	"ja": {
		"jupyter_noauth":     "「%s」の Jupyter（8888番）は、トークンもパスワードもなしで家のネットワークに開いています。家の誰でも、この PC で好きなプログラムを動かせる状態です。",
		"jupyter_noauth_fix": "Jupyter を起動するときに token かパスワードを設定してください（`--IdentityProvider.token=''` や `--NotebookApp.token=''` を付けて起動しない）。",
	},
	"zh-Hans": {
		"jupyter_noauth":     "%s 上的 Jupyter（端口 8888）在没有 token 或密码的情况下向您的网络开放。家中任何人都可以在这台电脑上运行任意程序。",
		"jupyter_noauth_fix": "启动 Jupyter 时请设置 token 或密码（不要使用 --IdentityProvider.token='' 或 --NotebookApp.token='' 启动）。",
	},
	"zh-Hant": {
		"jupyter_noauth":     "%s 上的 Jupyter（連接埠 8888）在沒有 token 或密碼的情況下向您的網路開放。家中任何人都可以在這台電腦上執行任意程式。",
		"jupyter_noauth_fix": "啟動 Jupyter 時請設定 token 或密碼（不要使用 --IdentityProvider.token='' 或 --NotebookApp.token='' 啟動）。",
	},
	"es": {
		"jupyter_noauth":     "Jupyter (puerto 8888) en %s está abierto a su red sin token ni contraseña. Cualquiera en casa puede ejecutar cualquier programa en este equipo.",
		"jupyter_noauth_fix": "Inicie Jupyter con un token o una contraseña (no lo inicie con --IdentityProvider.token='' o --NotebookApp.token='').",
	},
	"pt": {
		"jupyter_noauth":     "O Jupyter (porta 8888) em %s está aberto à sua rede sem token nem senha. Qualquer pessoa em casa pode executar qualquer programa neste computador.",
		"jupyter_noauth_fix": "Inicie o Jupyter com um token ou senha (não o inicie com --IdentityProvider.token='' ou --NotebookApp.token='').",
	},
	"fr": {
		"jupyter_noauth":     "Jupyter (port 8888) sur %s est ouvert à votre réseau sans token ni mot de passe. Toute personne à la maison peut exécuter n’importe quel programme sur cet ordinateur.",
		"jupyter_noauth_fix": "Lancez Jupyter avec un token ou un mot de passe (ne le lancez pas avec --IdentityProvider.token='' ou --NotebookApp.token='').",
	},
	"de": {
		"jupyter_noauth":     "Jupyter (Port 8888) auf %s ist ohne token oder Passwort in Ihrem Netzwerk erreichbar. Jede Person zu Hause kann beliebige Programme auf diesem Computer ausführen.",
		"jupyter_noauth_fix": "Starten Sie Jupyter mit einem token oder Passwort (nicht mit --IdentityProvider.token='' oder --NotebookApp.token='').",
	},
	"ru": {
		"jupyter_noauth":     "Jupyter (порт 8888) на устройстве %s открыт в вашей сети без token и пароля. Любой человек дома может запускать любые программы на этом компьютере.",
		"jupyter_noauth_fix": "Запускайте Jupyter с token или паролем (не используйте --IdentityProvider.token='' или --NotebookApp.token='').",
	},
	"ko": {
		"jupyter_noauth":     "%s의 Jupyter(포트 8888)가 token이나 비밀번호 없이 네트워크에 열려 있습니다. 집 안의 누구나 이 컴퓨터에서 원하는 프로그램을 실행할 수 있습니다.",
		"jupyter_noauth_fix": "token 또는 비밀번호를 설정해 Jupyter를 실행하세요(--IdentityProvider.token='' 또는 --NotebookApp.token='' 옵션으로 실행하지 마세요).",
	},
	"id": {
		"jupyter_noauth":     "Jupyter (port 8888) di %s terbuka ke jaringan Anda tanpa token atau kata sandi. Siapa pun di rumah dapat menjalankan program apa pun di komputer ini.",
		"jupyter_noauth_fix": "Jalankan Jupyter dengan token atau kata sandi (jangan gunakan --IdentityProvider.token='' atau --NotebookApp.token='').",
	},
}

func init() {
	add(jupyterMessages)
}
