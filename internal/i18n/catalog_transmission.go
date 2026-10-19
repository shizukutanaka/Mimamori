package i18n

var transmissionMessages = map[string]map[string]string{
	"en": {
		"transmission_noauth":     "Transmission's remote control (port 9091) on %s accepts commands from your network without a password. Anyone at home can change where files are saved and set a program to run when a download finishes.",
		"transmission_noauth_fix": "Turn on RPC authentication with a username and password in Transmission (for the linuxserver.io Docker image, set USER and PASS).",
	},
	"ja": {
		"transmission_noauth":     "「%s」の Transmission の操作窓口（9091番）は、合言葉なしで家のネットワークに開いています。家の誰でも、ファイルの保存先を変えたり、取り込みが終わったときに動かすプログラムを設定したりできる状態です。",
		"transmission_noauth_fix": "Transmission の設定で、操作窓口のユーザー名と合言葉を設定してください（linuxserver.io の Docker なら USER と PASS）。",
	},
	"zh-Hans": {
		"transmission_noauth":     "%s 上的 Transmission 远程控制（端口 9091）在没有密码的情况下接受来自您网络的命令。家中任何人都可以更改文件保存位置，并设置下载完成后运行的程序。",
		"transmission_noauth_fix": "请在 Transmission 中启用 RPC 身份验证并设置用户名和密码（linuxserver.io Docker 镜像请设置 USER 和 PASS）。",
	},
	"zh-Hant": {
		"transmission_noauth":     "%s 上的 Transmission 遠端控制（連接埠 9091）在沒有密碼的情況下接受來自您網路的命令。家中任何人都可以更改檔案儲存位置，並設定下載完成後執行的程式。",
		"transmission_noauth_fix": "請在 Transmission 中啟用 RPC 驗證並設定使用者名稱和密碼（linuxserver.io Docker 映像請設定 USER 和 PASS）。",
	},
	"es": {
		"transmission_noauth":     "El control remoto de Transmission (puerto 9091) en %s acepta comandos de su red sin contraseña. Cualquiera en casa puede cambiar dónde se guardan los archivos y configurar un programa para que se ejecute al terminar una descarga.",
		"transmission_noauth_fix": "Active la autenticación RPC con usuario y contraseña en Transmission (para la imagen Docker de linuxserver.io, configure USER y PASS).",
	},
	"pt": {
		"transmission_noauth":     "O controle remoto do Transmission (porta 9091) em %s aceita comandos da sua rede sem senha. Qualquer pessoa em casa pode alterar onde os arquivos são salvos e definir um programa para executar quando um download terminar.",
		"transmission_noauth_fix": "Ative a autenticação RPC com usuário e senha no Transmission (na imagem Docker da linuxserver.io, defina USER e PASS).",
	},
	"fr": {
		"transmission_noauth":     "La commande à distance de Transmission (port 9091) sur %s accepte les commandes de votre réseau sans mot de passe. Toute personne à la maison peut modifier l’emplacement d’enregistrement des fichiers et définir un programme à exécuter à la fin d’un téléchargement.",
		"transmission_noauth_fix": "Activez l’authentification RPC avec un nom d’utilisateur et un mot de passe dans Transmission (pour l’image Docker linuxserver.io, définissez USER et PASS).",
	},
	"de": {
		"transmission_noauth":     "Die Fernsteuerung von Transmission (Port 9091) auf %s nimmt Befehle aus Ihrem Netzwerk ohne Passwort an. Jede Person zu Hause kann den Speicherort für Dateien ändern und ein Programm festlegen, das nach einem Download ausgeführt wird.",
		"transmission_noauth_fix": "Aktivieren Sie die RPC-Authentifizierung mit Benutzername und Passwort in Transmission (beim linuxserver.io-Docker-Image USER und PASS setzen).",
	},
	"ru": {
		"transmission_noauth":     "Удалённое управление Transmission (порт 9091) на устройстве %s принимает команды из вашей сети без пароля. Любой человек дома может изменить папку для сохранения файлов и задать программу, запускаемую после завершения загрузки.",
		"transmission_noauth_fix": "Включите RPC-аутентификацию с именем пользователя и паролем в Transmission (для Docker-образа linuxserver.io задайте USER и PASS).",
	},
	"ko": {
		"transmission_noauth":     "%s의 Transmission 원격 제어(포트 9091)가 네트워크에서 비밀번호 없이 명령을 받습니다. 집 안의 누구나 파일 저장 위치를 바꾸고 다운로드가 끝날 때 실행할 프로그램을 설정할 수 있습니다.",
		"transmission_noauth_fix": "Transmission에서 사용자 이름과 비밀번호로 RPC 인증을 켜세요(linuxserver.io Docker 이미지에서는 USER와 PASS를 설정하세요).",
	},
	"id": {
		"transmission_noauth":     "Kendali jarak jauh Transmission (port 9091) di %s menerima perintah dari jaringan Anda tanpa kata sandi. Siapa pun di rumah dapat mengubah lokasi penyimpanan berkas dan mengatur program yang berjalan setelah unduhan selesai.",
		"transmission_noauth_fix": "Aktifkan autentikasi RPC dengan nama pengguna dan kata sandi di Transmission (untuk image Docker linuxserver.io, atur USER dan PASS).",
	},
}

var transmissionCatalogTargets = map[string]map[string]string{
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
	for lang, catalog := range transmissionMessages {
		target := transmissionCatalogTargets[lang]
		for key, value := range catalog {
			target[key] = value
		}
		messages[lang] = target
	}
}
