package i18n

var syncthingMessages = map[string]map[string]string{
	"en": {
		"syncthing_noauth":     "Syncthing's web GUI (port 8384) on %s has no password. Anyone on your network can add their own device and copy or overwrite files on this machine.",
		"syncthing_noauth_fix": "Set a GUI Authentication User and Password in Syncthing's Settings → GUI.",
	},
	"ja": {
		"syncthing_noauth":     "「%s」の Syncthing の管理画面（8384番）は、合言葉が設定されていません。家のネットワークの誰でも同期の相手を足して、この機器のファイルを抜き取ったり書き換えたりできる状態です。",
		"syncthing_noauth_fix": "Syncthing の管理画面の「設定」→「GUI」で、GUI 認証のユーザー名と合言葉を設定してください。",
	},
	"zh-Hans": {
		"syncthing_noauth":     "%s 上的 Syncthing 网页管理界面（端口 8384）没有密码。您网络中的任何人都可以添加自己的设备，并复制或覆盖此机器上的文件。",
		"syncthing_noauth_fix": "请在 Syncthing 的“设置”→“GUI”中设置 GUI 身份验证用户名和密码。",
	},
	"zh-Hant": {
		"syncthing_noauth":     "%s 上的 Syncthing 網頁管理介面（連接埠 8384）沒有密碼。您網路中的任何人都可以新增自己的裝置，並複製或覆寫這台機器上的檔案。",
		"syncthing_noauth_fix": "請在 Syncthing 的「設定」→「GUI」中設定 GUI 驗證使用者名稱和密碼。",
	},
	"es": {
		"syncthing_noauth":     "La interfaz web de Syncthing (puerto 8384) en %s no tiene contraseña. Cualquiera de su red puede añadir su propio dispositivo y copiar o sobrescribir archivos de este equipo.",
		"syncthing_noauth_fix": "Configure un usuario y una contraseña de autenticación GUI en Ajustes → GUI de Syncthing.",
	},
	"pt": {
		"syncthing_noauth":     "A interface web do Syncthing (porta 8384) em %s está sem senha. Qualquer pessoa na sua rede pode adicionar o próprio dispositivo e copiar ou sobrescrever arquivos deste computador.",
		"syncthing_noauth_fix": "Defina um usuário e uma senha de autenticação da GUI em Configurações → GUI do Syncthing.",
	},
	"fr": {
		"syncthing_noauth":     "L’interface Web de Syncthing (port 8384) sur %s n’a pas de mot de passe. Toute personne de votre réseau peut ajouter son propre appareil et copier ou écraser des fichiers sur cette machine.",
		"syncthing_noauth_fix": "Définissez un nom d’utilisateur et un mot de passe d’authentification GUI dans Paramètres → GUI de Syncthing.",
	},
	"de": {
		"syncthing_noauth":     "Die Syncthing-Weboberfläche (Port 8384) auf %s hat kein Passwort. Jede Person in Ihrem Netzwerk kann ein eigenes Gerät hinzufügen und Dateien auf diesem Gerät kopieren oder überschreiben.",
		"syncthing_noauth_fix": "Legen Sie unter Einstellungen → GUI von Syncthing einen GUI-Authentifizierungsbenutzer und ein Passwort fest.",
	},
	"ru": {
		"syncthing_noauth":     "Веб-интерфейс Syncthing (порт 8384) на устройстве %s не защищён паролем. Любой пользователь вашей сети может добавить своё устройство и скопировать или перезаписать файлы на этом компьютере.",
		"syncthing_noauth_fix": "Задайте имя пользователя и пароль для GUI в разделе «Настройки» → «GUI» Syncthing.",
	},
	"ko": {
		"syncthing_noauth":     "%s의 Syncthing 웹 GUI(포트 8384)에 비밀번호가 설정되어 있지 않습니다. 네트워크의 누구나 자신의 기기를 추가하고 이 기기의 파일을 복사하거나 덮어쓸 수 있습니다.",
		"syncthing_noauth_fix": "Syncthing의 설정 → GUI에서 GUI 인증 사용자와 비밀번호를 설정하세요.",
	},
	"id": {
		"syncthing_noauth":     "GUI web Syncthing (port 8384) di %s tidak memiliki kata sandi. Siapa pun di jaringan Anda dapat menambahkan perangkat sendiri dan menyalin atau menimpa berkas di mesin ini.",
		"syncthing_noauth_fix": "Tetapkan pengguna dan kata sandi autentikasi GUI di Pengaturan → GUI Syncthing.",
	},
}

func init() {
	add(syncthingMessages)
}
