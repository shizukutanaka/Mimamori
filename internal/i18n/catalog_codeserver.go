package i18n

var codeserverMessages = map[string]map[string]string{
	"en": {
		"codeserver_noauth":     "code-server (VS Code in the browser, port 8443) on %s is open to your network with no password. Anyone at home can open its terminal and run any command.",
		"codeserver_noauth_fix": "Set the PASSWORD or HASHED_PASSWORD environment variable for the container and restart it (don't start it with --auth none).",
	},
	"ja": {
		"codeserver_noauth":     "「%s」の code-server（ブラウザで使う VS Code、8443番）は、合言葉なしで家のネットワークに開いています。家の誰でも、ブラウザから端末を開き、好きな命令を実行できる状態です。",
		"codeserver_noauth_fix": "Docker の環境変数 PASSWORD か HASHED_PASSWORD を設定して、起動し直してください（--auth none で起動しない）。",
	},
	"zh-Hans": {
		"codeserver_noauth":     "%s 上的 code-server（浏览器中的 VS Code，端口 8443）在没有密码的情况下对您的网络开放。家中任何人都可以通过浏览器打开其终端并执行任意命令。",
		"codeserver_noauth_fix": "请为容器设置 PASSWORD 或 HASHED_PASSWORD 环境变量并重新启动（不要使用 --auth none 启动）。",
	},
	"zh-Hant": {
		"codeserver_noauth":     "%s 上的 code-server（瀏覽器中的 VS Code，連接埠 8443）未設定密碼便對您的網路開放。家中任何人都可以透過瀏覽器開啟終端機並執行任何命令。",
		"codeserver_noauth_fix": "請為容器設定 PASSWORD 或 HASHED_PASSWORD 環境變數並重新啟動（不要使用 --auth none 啟動）。",
	},
	"es": {
		"codeserver_noauth":     "code-server (VS Code en el navegador, puerto 8443) en %s está abierto a su red sin contraseña. Cualquiera en casa puede abrir su terminal desde el navegador y ejecutar cualquier comando.",
		"codeserver_noauth_fix": "Configure la variable de entorno PASSWORD o HASHED_PASSWORD del contenedor y reinícielo (no lo inicie con --auth none).",
	},
	"pt": {
		"codeserver_noauth":     "O code-server (VS Code no navegador, porta 8443) em %s está aberto à sua rede sem senha. Qualquer pessoa em casa pode abrir o terminal no navegador e executar qualquer comando.",
		"codeserver_noauth_fix": "Defina a variável de ambiente PASSWORD ou HASHED_PASSWORD do contêiner e reinicie-o (não o inicie com --auth none).",
	},
	"fr": {
		"codeserver_noauth":     "code-server (VS Code dans le navigateur, port 8443) sur %s est ouvert à votre réseau sans mot de passe. Toute personne à la maison peut ouvrir son terminal dans le navigateur et exécuter n’importe quelle commande.",
		"codeserver_noauth_fix": "Définissez la variable d’environnement PASSWORD ou HASHED_PASSWORD du conteneur et redémarrez-le (ne le démarrez pas avec --auth none).",
	},
	"de": {
		"codeserver_noauth":     "code-server (VS Code im Browser, Port 8443) auf %s ist ohne Passwort für Ihr Netzwerk offen. Jede Person zu Hause kann das Terminal im Browser öffnen und beliebige Befehle ausführen.",
		"codeserver_noauth_fix": "Setzen Sie die Umgebungsvariable PASSWORD oder HASHED_PASSWORD für den Container und starten Sie ihn neu (nicht mit --auth none starten).",
	},
	"ru": {
		"codeserver_noauth":     "code-server (VS Code в браузере, порт 8443) на устройстве %s открыт в вашей сети без пароля. Любой человек дома может открыть терминал в браузере и выполнять любые команды.",
		"codeserver_noauth_fix": "Задайте для контейнера переменную окружения PASSWORD или HASHED_PASSWORD и перезапустите его (не запускайте с --auth none).",
	},
	"ko": {
		"codeserver_noauth":     "%s의 code-server(브라우저에서 사용하는 VS Code, 포트 8443)가 비밀번호 없이 네트워크에 열려 있습니다. 집 안의 누구나 브라우저에서 터미널을 열고 어떤 명령이든 실행할 수 있습니다.",
		"codeserver_noauth_fix": "컨테이너의 환경 변수 PASSWORD 또는 HASHED_PASSWORD를 설정하고 다시 시작하세요(--auth none으로 시작하지 마세요).",
	},
	"id": {
		"codeserver_noauth":     "code-server (VS Code di browser, port 8443) di %s terbuka ke jaringan Anda tanpa kata sandi. Siapa pun di rumah dapat membuka terminalnya melalui browser dan menjalankan perintah apa pun.",
		"codeserver_noauth_fix": "Atur variabel lingkungan PASSWORD atau HASHED_PASSWORD untuk kontainer, lalu mulai ulang (jangan jalankan dengan --auth none).",
	},
}

func init() {
	add(codeserverMessages)
}
