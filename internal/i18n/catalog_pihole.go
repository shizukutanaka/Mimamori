package i18n

var piholeMessages = map[string]map[string]string{
	"en": {
		"pihole_noauth":     "Pi-hole (port 80) on %s has an empty admin password. Anyone on your network can change its DNS settings and send every device in your home to fake websites.",
		"pihole_noauth_fix": "Set a password in the Pi-hole admin page or with `pihole setpassword`.",
	},
	"ja": {
		"pihole_noauth":     "「%s」の Pi-hole（80番）は、管理画面のパスワードが空です。家のネットワークの誰でも DNS の設定を書き換えて、家じゅうの機器を偽のサイトへ向けられる状態です。",
		"pihole_noauth_fix": "Pi-hole の管理画面か `pihole setpassword` で、パスワードを設定してください。",
	},
	"zh-Hans": {
		"pihole_noauth":     "%s 上的 Pi-hole（端口 80）管理员密码为空。您网络中的任何人都可以更改其 DNS 设置，将家中所有设备引导至虚假网站。",
		"pihole_noauth_fix": "请在 Pi-hole 管理页面或使用 `pihole setpassword` 设置密码。",
	},
	"zh-Hant": {
		"pihole_noauth":     "%s 上的 Pi-hole（連接埠 80）管理員密碼為空。您網路中的任何人都可以更改其 DNS 設定，將家中所有裝置導向虛假網站。",
		"pihole_noauth_fix": "請在 Pi-hole 管理頁面或使用 `pihole setpassword` 設定密碼。",
	},
	"es": {
		"pihole_noauth":     "Pi-hole (puerto 80) en %s tiene vacía la contraseña de administrador. Cualquiera de su red puede cambiar su configuración DNS y enviar todos los dispositivos de su hogar a sitios web falsos.",
		"pihole_noauth_fix": "Establezca una contraseña en la página de administración de Pi-hole o con `pihole setpassword`.",
	},
	"pt": {
		"pihole_noauth":     "O Pi-hole (porta 80) em %s está sem senha de administrador. Qualquer pessoa da sua rede pode alterar as configurações DNS e enviar todos os dispositivos da sua casa para sites falsos.",
		"pihole_noauth_fix": "Defina uma senha na página de administração do Pi-hole ou com `pihole setpassword`.",
	},
	"fr": {
		"pihole_noauth":     "Pi-hole (port 80) sur %s n’a pas de mot de passe administrateur. Toute personne de votre réseau peut modifier ses paramètres DNS et envoyer tous les appareils de votre domicile vers de faux sites Web.",
		"pihole_noauth_fix": "Définissez un mot de passe dans la page d’administration de Pi-hole ou avec `pihole setpassword`.",
	},
	"de": {
		"pihole_noauth":     "Pi-hole (Port 80) auf %s hat kein Administratorkennwort. Jede Person in Ihrem Netzwerk kann die DNS-Einstellungen ändern und alle Geräte in Ihrem Zuhause auf gefälschte Webseiten umleiten.",
		"pihole_noauth_fix": "Legen Sie auf der Pi-hole-Verwaltungsseite oder mit `pihole setpassword` ein Passwort fest.",
	},
	"ru": {
		"pihole_noauth":     "У Pi-hole (порт 80) на устройстве %s пустой пароль администратора. Любой пользователь вашей сети может изменить настройки DNS и направить все домашние устройства на поддельные сайты.",
		"pihole_noauth_fix": "Задайте пароль на странице администратора Pi-hole или с помощью `pihole setpassword`.",
	},
	"ko": {
		"pihole_noauth":     "%s의 Pi-hole(포트 80)에 관리자 비밀번호가 설정되어 있지 않습니다. 네트워크의 누구나 DNS 설정을 바꾸어 집 안의 모든 기기를 가짜 웹사이트로 보낼 수 있습니다.",
		"pihole_noauth_fix": "Pi-hole 관리자 페이지나 `pihole setpassword` 명령으로 비밀번호를 설정하세요.",
	},
	"id": {
		"pihole_noauth":     "Pi-hole (port 80) di %s tidak memiliki kata sandi administrator. Siapa pun di jaringan Anda dapat mengubah pengaturan DNS dan mengarahkan semua perangkat di rumah Anda ke situs web palsu.",
		"pihole_noauth_fix": "Tetapkan kata sandi di halaman admin Pi-hole atau dengan `pihole setpassword`.",
	},
}

func init() {
	add(piholeMessages)
}
