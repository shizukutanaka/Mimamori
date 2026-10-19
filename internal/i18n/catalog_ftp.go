package i18n

var ftpMessages = map[string]map[string]string{
	"en": {
		"ftp_anonymous":     "The file server (FTP) on %s lets anyone on your network log in without a password and see its files.",
		"ftp_anonymous_fix": "In the device's settings, turn off anonymous (guest) FTP access, or turn FTP off if you don't use it. For USB storage plugged into your router, check the router's USB sharing settings.",
	},
	"ja": {
		"ftp_anonymous":     "「%s」のファイル共有（FTP）に、家の中の誰でもパスワードなしでログインし、ファイルを見られる状態です。",
		"ftp_anonymous_fix": "機器の設定で FTP の匿名（ゲスト）アクセスを無効にしてください。使っていない場合は FTP 自体をオフにしてください。ルーターに挿した USB メモリの場合は、ルーターの USB 共有の設定を確認してください。",
	},
	"zh-Hans": {
		"ftp_anonymous":     "“%s”上的文件服务器（FTP）允许您家中网络上的任何人无需密码登录并查看其文件。",
		"ftp_anonymous_fix": "请在设备设置中关闭匿名（访客）FTP 访问；如果您不使用 FTP，请关闭 FTP。对于连接到路由器的 USB 存储设备，请检查路由器的 USB 共享设置。",
	},
	"zh-Hant": {
		"ftp_anonymous":     "「%s」上的檔案伺服器（FTP）允許您家中網路上的任何人無需密碼登入並查看其中的檔案。",
		"ftp_anonymous_fix": "請在裝置設定中關閉匿名（訪客）FTP 存取；如果您不使用 FTP，請將 FTP 關閉。若是插在路由器上的 USB 儲存裝置，請檢查路由器的 USB 分享設定。",
	},
	"es": {
		"ftp_anonymous":     "El servidor de archivos (FTP) en %s permite que cualquier persona de su red inicie sesión sin contraseña y vea sus archivos.",
		"ftp_anonymous_fix": "En la configuración del dispositivo, desactive el acceso FTP anónimo (de invitado) o desactive FTP si no lo utiliza. Para el almacenamiento USB conectado al router, revise la configuración de uso compartido USB del router.",
	},
	"pt": {
		"ftp_anonymous":     "O servidor de arquivos (FTP) em %s permite que qualquer pessoa na sua rede faça login sem senha e veja seus arquivos.",
		"ftp_anonymous_fix": "Nas configurações do dispositivo, desative o acesso FTP anônimo (de convidado) ou desative o FTP se você não o usa. Para armazenamento USB conectado ao roteador, verifique as configurações de compartilhamento USB do roteador.",
	},
	"fr": {
		"ftp_anonymous":     "Le serveur de fichiers (FTP) sur %s permet à toute personne de votre réseau de se connecter sans mot de passe et de voir ses fichiers.",
		"ftp_anonymous_fix": "Dans les paramètres de l’appareil, désactivez l’accès FTP anonyme (invité), ou désactivez FTP si vous ne l’utilisez pas. Pour un stockage USB branché sur votre routeur, vérifiez les paramètres de partage USB du routeur.",
	},
	"de": {
		"ftp_anonymous":     "Der Dateiserver (FTP) auf %s ermöglicht allen in Ihrem Netzwerk, sich ohne Passwort anzumelden und seine Dateien einzusehen.",
		"ftp_anonymous_fix": "Deaktivieren Sie in den Geräteeinstellungen den anonymen (Gast-)FTP-Zugriff oder schalten Sie FTP aus, wenn Sie ihn nicht nutzen. Prüfen Sie bei einem am Router angeschlossenen USB-Speicher die USB-Freigabeeinstellungen des Routers.",
	},
	"ru": {
		"ftp_anonymous":     "Файловый сервер (FTP) на %s позволяет любому пользователю в вашей сети входить без пароля и просматривать его файлы.",
		"ftp_anonymous_fix": "В настройках устройства отключите анонимный (гостевой) доступ по FTP или отключите FTP, если вы им не пользуетесь. Для USB-накопителя, подключённого к роутеру, проверьте настройки USB-общего доступа на роутере.",
	},
	"ko": {
		"ftp_anonymous":     "%s의 파일 서버(FTP)는 네트워크의 누구나 비밀번호 없이 로그인하여 파일을 볼 수 있습니다.",
		"ftp_anonymous_fix": "기기 설정에서 익명(게스트) FTP 접근을 끄세요. 사용하지 않는다면 FTP를 끄세요. 라우터에 연결한 USB 저장장치라면 라우터의 USB 공유 설정을 확인하세요.",
	},
	"id": {
		"ftp_anonymous":     "Server file (FTP) di %s memungkinkan siapa pun di jaringan Anda masuk tanpa kata sandi dan melihat file-nya.",
		"ftp_anonymous_fix": "Di pengaturan perangkat, nonaktifkan akses FTP anonim (tamu), atau matikan FTP jika Anda tidak menggunakannya. Untuk penyimpanan USB yang terhubung ke router, periksa pengaturan berbagi USB router.",
	},
}

var ftpCatalogTargets = map[string]map[string]string{
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
	for lang, catalog := range ftpCatalogTargets {
		for key, value := range ftpMessages[lang] {
			catalog[key] = value
		}
		messages[lang] = catalog
	}
}
