package i18n

var esphomeMessages = map[string]map[string]string{
	"en": {
		"esphome_noauth":     "Anyone can use the ESPHome dashboard (%s) without a password (port 6052)",
		"esphome_noauth_fix": "Set a username and password for ESPHome (in Docker, the ESPHOME_USERNAME and ESPHOME_PASSWORD environment variables; older versions use USERNAME and PASSWORD), then restart it.",
	},
	"ja": {
		"esphome_noauth":     "ESPHome の管理画面（%s）を、パスワードなしで誰でも操作できます（6052番）",
		"esphome_noauth_fix": "ESPHome にユーザー名とパスワードを設定して（Docker なら環境変数 ESPHOME_USERNAME と ESPHOME_PASSWORD。古い版は USERNAME と PASSWORD）、再起動してください。",
	},
	"zh-Hans": {
		"esphome_noauth":     "任何人都可以在没有密码的情况下使用 ESPHome 管理界面（%s）（端口 6052）。",
		"esphome_noauth_fix": "为 ESPHome 设置用户名和密码（Docker 中使用 ESPHOME_USERNAME 和 ESPHOME_PASSWORD 环境变量；旧版本使用 USERNAME 和 PASSWORD），然后重新启动。",
	},
	"zh-Hant": {
		"esphome_noauth":     "任何人都可以在沒有密碼的情況下使用 ESPHome 管理介面（%s）（連接埠 6052）。",
		"esphome_noauth_fix": "為 ESPHome 設定使用者名稱和密碼（Docker 中使用 ESPHOME_USERNAME 和 ESPHOME_PASSWORD 環境變數；舊版本使用 USERNAME 和 PASSWORD），然後重新啟動。",
	},
	"es": {
		"esphome_noauth":     "Cualquiera puede usar el panel de ESPHome (%s) sin contraseña (puerto 6052).",
		"esphome_noauth_fix": "Configure un nombre de usuario y una contraseña para ESPHome (en Docker, las variables ESPHOME_USERNAME y ESPHOME_PASSWORD; las versiones antiguas usan USERNAME y PASSWORD) y reinícielo.",
	},
	"pt": {
		"esphome_noauth":     "Qualquer pessoa pode usar o painel do ESPHome (%s) sem senha (porta 6052).",
		"esphome_noauth_fix": "Defina um nome de usuário e uma senha para o ESPHome (no Docker, as variáveis ESPHOME_USERNAME e ESPHOME_PASSWORD; versões antigas usam USERNAME e PASSWORD) e reinicie-o.",
	},
	"fr": {
		"esphome_noauth":     "N’importe qui peut utiliser le tableau de bord ESPHome (%s) sans mot de passe (port 6052).",
		"esphome_noauth_fix": "Définissez un nom d’utilisateur et un mot de passe pour ESPHome (dans Docker, les variables ESPHOME_USERNAME et ESPHOME_PASSWORD ; les anciennes versions utilisent USERNAME et PASSWORD), puis redémarrez-le.",
	},
	"de": {
		"esphome_noauth":     "Jede Person kann das ESPHome-Dashboard (%s) ohne Passwort verwenden (Port 6052).",
		"esphome_noauth_fix": "Legen Sie für ESPHome einen Benutzernamen und ein Passwort fest (in Docker die Umgebungsvariablen ESPHOME_USERNAME und ESPHOME_PASSWORD; ältere Versionen verwenden USERNAME und PASSWORD) und starten Sie es neu.",
	},
	"ru": {
		"esphome_noauth":     "Любой может использовать панель ESPHome (%s) без пароля (порт 6052).",
		"esphome_noauth_fix": "Задайте имя пользователя и пароль для ESPHome (в Docker — переменные ESPHOME_USERNAME и ESPHOME_PASSWORD; старые версии используют USERNAME и PASSWORD), затем перезапустите его.",
	},
	"ko": {
		"esphome_noauth":     "누구나 비밀번호 없이 ESPHome 대시보드(%s)를 사용할 수 있습니다(포트 6052).",
		"esphome_noauth_fix": "ESPHome에 사용자 이름과 비밀번호를 설정하세요(Docker에서는 ESPHOME_USERNAME 및 ESPHOME_PASSWORD 환경 변수, 이전 버전은 USERNAME 및 PASSWORD 사용). 그런 다음 다시 시작하세요.",
	},
	"id": {
		"esphome_noauth":     "Siapa pun dapat menggunakan dasbor ESPHome (%s) tanpa kata sandi (port 6052).",
		"esphome_noauth_fix": "Atur nama pengguna dan kata sandi untuk ESPHome (di Docker, variabel lingkungan ESPHOME_USERNAME dan ESPHOME_PASSWORD; versi lama menggunakan USERNAME dan PASSWORD), lalu mulai ulang.",
	},
}

func init() {
	add(esphomeMessages)
}
