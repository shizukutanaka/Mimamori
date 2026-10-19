package i18n

var shellyMessages = map[string]map[string]string{
	"en": {
		"shelly_noauth":     "The Shelly smart device on %s has no password, so anyone on your network can switch it on or off and replace its firmware.",
		"shelly_noauth_fix": "In the Shelly app or the device's web page, open the authentication settings, turn on password protection, and set a password.",
	},
	"ja": {
		"shelly_noauth":     "「%s」の Shelly の機器にパスワードがなく、家の中の誰でも電源を入れたり切ったり、ファームウェアを書き換えたりできる状態です。",
		"shelly_noauth_fix": "Shelly のアプリか機器の画面で、認証（パスワードの保護）の設定を開いてオンにし、パスワードを設定してください。",
	},
	"zh-Hans": {
		"shelly_noauth":     "%s 上的 Shelly 智能设备没有密码，因此网络中的任何人都可以开关设备并替换固件。",
		"shelly_noauth_fix": "请您在 Shelly 应用或设备网页中打开身份验证设置，启用密码保护并设置密码。",
	},
	"zh-Hant": {
		"shelly_noauth":     "%s 上的 Shelly 智慧裝置沒有密碼，因此網路中的任何人都可以開關裝置並替換韌體。",
		"shelly_noauth_fix": "請您在 Shelly 應用程式或裝置網頁中開啟驗證設定，啟用密碼保護並設定密碼。",
	},
	"es": {
		"shelly_noauth":     "El dispositivo inteligente Shelly en %s no tiene contraseña, por lo que cualquier persona de su red puede encenderlo o apagarlo y reemplazar su firmware.",
		"shelly_noauth_fix": "En la aplicación Shelly o en la página web del dispositivo, abra los ajustes de autenticación, active la protección con contraseña y establezca una contraseña.",
	},
	"pt": {
		"shelly_noauth":     "O dispositivo inteligente Shelly em %s está sem senha, então qualquer pessoa da sua rede pode ligá-lo ou desligá-lo e substituir o firmware.",
		"shelly_noauth_fix": "No aplicativo Shelly ou na página web do dispositivo, abra as configurações de autenticação, ative a proteção por senha e defina uma senha.",
	},
	"fr": {
		"shelly_noauth":     "L’appareil connecté Shelly sur %s n’a pas de mot de passe ; toute personne de votre réseau peut donc l’allumer ou l’éteindre et remplacer son micrologiciel.",
		"shelly_noauth_fix": "Dans l’application Shelly ou sur la page Web de l’appareil, ouvrez les paramètres d’authentification, activez la protection par mot de passe et définissez un mot de passe.",
	},
	"de": {
		"shelly_noauth":     "Das Shelly-Smartgerät auf %s hat kein Passwort. Daher kann jede Person in Ihrem Netzwerk es ein- oder ausschalten und die Firmware ersetzen.",
		"shelly_noauth_fix": "Öffnen Sie in der Shelly-App oder auf der Webseite des Geräts die Authentifizierungseinstellungen, aktivieren Sie den Passwortschutz und legen Sie ein Passwort fest.",
	},
	"ru": {
		"shelly_noauth":     "У устройства Shelly на %s нет пароля, поэтому любой пользователь вашей сети может включать и выключать его и заменять прошивку.",
		"shelly_noauth_fix": "В приложении Shelly или на веб-странице устройства откройте настройки аутентификации, включите защиту паролем и задайте пароль.",
	},
	"ko": {
		"shelly_noauth":     "%s의 Shelly 스마트 기기에는 비밀번호가 없어 네트워크의 누구나 전원을 켜고 끄거나 펌웨어를 바꿀 수 있습니다.",
		"shelly_noauth_fix": "Shelly 앱 또는 기기의 웹 페이지에서 인증 설정을 열고 비밀번호 보호를 켠 다음 비밀번호를 설정하세요.",
	},
	"id": {
		"shelly_noauth":     "Perangkat pintar Shelly di %s tidak memiliki kata sandi, sehingga siapa pun di jaringan Anda dapat menyalakan atau mematikannya dan mengganti firmware-nya.",
		"shelly_noauth_fix": "Di aplikasi Shelly atau halaman web perangkat, buka pengaturan autentikasi, aktifkan perlindungan kata sandi, lalu tetapkan kata sandi.",
	},
}

var shellyCatalogTargets = map[string]map[string]string{
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
	for lang, catalog := range shellyMessages {
		target := shellyCatalogTargets[lang]
		for key, value := range catalog {
			target[key] = value
		}
		messages[lang] = target
	}
}
