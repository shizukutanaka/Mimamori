package i18n

var esphomeMessages = map[string]map[string]string{
	"en": {
		"esphome_noauth":     "The ESPHome dashboard (port 6052) on %s has no password. Anyone on your network can edit and reflash your ESPHome devices.",
		"esphome_noauth_fix": "Set a username and password for ESPHome (in Docker, the ESPHOME_USERNAME and ESPHOME_PASSWORD environment variables; older versions use USERNAME and PASSWORD), then restart it.",
	},
	"ja": {
		"esphome_noauth":     "「%s」の ESPHome の管理画面（6052番）に、パスワードがありません。家のネットワークの誰でも、ESPHome の機器の設定を書き換えたり、ファームウェアを入れ直したりできます。",
		"esphome_noauth_fix": "ESPHome にユーザー名とパスワードを設定して（Docker なら環境変数 ESPHOME_USERNAME と ESPHOME_PASSWORD。古い版は USERNAME と PASSWORD）、再起動してください。",
	},
	"zh-Hans": {
		"esphome_noauth":     "%s 上的 ESPHome 管理界面（端口 6052）没有密码。您网络中的任何人都可以修改并重新刷写您的 ESPHome 设备。",
		"esphome_noauth_fix": "为 ESPHome 设置用户名和密码（Docker 中使用 ESPHOME_USERNAME 和 ESPHOME_PASSWORD 环境变量；旧版本使用 USERNAME 和 PASSWORD），然后重新启动。",
	},
	"zh-Hant": {
		"esphome_noauth":     "%s 上的 ESPHome 管理介面（連接埠 6052）沒有密碼。您網路中的任何人都可以修改並重新燒錄您的 ESPHome 裝置。",
		"esphome_noauth_fix": "為 ESPHome 設定使用者名稱和密碼（Docker 中使用 ESPHOME_USERNAME 和 ESPHOME_PASSWORD 環境變數；舊版本使用 USERNAME 和 PASSWORD），然後重新啟動。",
	},
	"es": {
		"esphome_noauth":     "El panel de ESPHome en %s (puerto 6052) no tiene contraseña. Cualquier persona de su red puede modificar y volver a flashear sus dispositivos ESPHome.",
		"esphome_noauth_fix": "Configure un nombre de usuario y una contraseña para ESPHome (en Docker, las variables ESPHOME_USERNAME y ESPHOME_PASSWORD; las versiones antiguas usan USERNAME y PASSWORD) y reinícielo.",
	},
	"pt": {
		"esphome_noauth":     "O painel do ESPHome em %s (porta 6052) não tem senha. Qualquer pessoa na sua rede pode alterar e regravar os seus dispositivos ESPHome.",
		"esphome_noauth_fix": "Defina um nome de usuário e uma senha para o ESPHome (no Docker, as variáveis ESPHOME_USERNAME e ESPHOME_PASSWORD; versões antigas usam USERNAME e PASSWORD) e reinicie-o.",
	},
	"fr": {
		"esphome_noauth":     "Le tableau de bord ESPHome sur %s (port 6052) n’a pas de mot de passe. Toute personne sur votre réseau peut modifier et reflasher vos appareils ESPHome.",
		"esphome_noauth_fix": "Définissez un nom d’utilisateur et un mot de passe pour ESPHome (dans Docker, les variables ESPHOME_USERNAME et ESPHOME_PASSWORD ; les anciennes versions utilisent USERNAME et PASSWORD), puis redémarrez-le.",
	},
	"de": {
		"esphome_noauth":     "Das ESPHome-Dashboard auf %s (Port 6052) hat kein Passwort. Jede Person in Ihrem Netzwerk kann Ihre ESPHome-Geräte umkonfigurieren und neu flashen.",
		"esphome_noauth_fix": "Legen Sie für ESPHome einen Benutzernamen und ein Passwort fest (in Docker die Umgebungsvariablen ESPHOME_USERNAME und ESPHOME_PASSWORD; ältere Versionen verwenden USERNAME und PASSWORD) und starten Sie es neu.",
	},
	"ru": {
		"esphome_noauth":     "Панель ESPHome на устройстве «%s» (порт 6052) не защищена паролем. Любой пользователь вашей сети может изменить и перепрошить ваши устройства ESPHome.",
		"esphome_noauth_fix": "Задайте имя пользователя и пароль для ESPHome (в Docker — переменные ESPHOME_USERNAME и ESPHOME_PASSWORD; старые версии используют USERNAME и PASSWORD), затем перезапустите его.",
	},
	"ko": {
		"esphome_noauth":     "%s의 ESPHome 대시보드(포트 6052)에 비밀번호가 없습니다. 네트워크의 누구나 ESPHome 기기의 설정을 바꾸고 펌웨어를 다시 쓸 수 있습니다.",
		"esphome_noauth_fix": "ESPHome에 사용자 이름과 비밀번호를 설정하세요(Docker에서는 ESPHOME_USERNAME 및 ESPHOME_PASSWORD 환경 변수, 이전 버전은 USERNAME 및 PASSWORD 사용). 그런 다음 다시 시작하세요.",
	},
	"id": {
		"esphome_noauth":     "Dasbor ESPHome di %s (port 6052) tidak memiliki kata sandi. Siapa pun di jaringan Anda dapat mengubah dan mem-flash ulang perangkat ESPHome Anda.",
		"esphome_noauth_fix": "Atur nama pengguna dan kata sandi untuk ESPHome (di Docker, variabel lingkungan ESPHOME_USERNAME dan ESPHOME_PASSWORD; versi lama menggunakan USERNAME dan PASSWORD), lalu mulai ulang.",
	},
}

func init() {
	add(esphomeMessages)
}
