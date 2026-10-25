package i18n

var frigateMessages = map[string]map[string]string{
	"en": {
		"frigate_noauth":     "The Frigate NVR (port 5000) on %s has no password. Anyone on your network can watch your cameras and change its settings as an administrator.",
		"frigate_noauth_fix": "Update Frigate to 0.14 or later, stop publishing port 5000 (remove 5000:5000 from the Docker ports), and use port 8971, which requires a login.",
	},
	"ja": {
		"frigate_noauth":     "「%s」の録画ソフト Frigate（5000番）に、パスワードがありません。家のネットワークの誰でも、管理者としてカメラの映像を見たり、設定を変えたりできます。",
		"frigate_noauth_fix": "Frigate を 0.14 以降に更新し、5000番の公開（Docker の ports の 5000:5000）を外して、ログインが必要な 8971番を使ってください。",
	},
	"zh-Hans": {
		"frigate_noauth":     "%s 上的 Frigate NVR（端口 5000）没有密码。您网络中的任何人都可以以管理员身份查看您的摄像头并更改设置。",
		"frigate_noauth_fix": "将 Frigate 更新到 0.14 或更高版本，停止公开端口 5000（从 Docker 的 ports 中移除 5000:5000），并使用需要登录的 8971 端口。",
	},
	"zh-Hant": {
		"frigate_noauth":     "%s 上的 Frigate NVR（連接埠 5000）沒有密碼。您網路中的任何人都可以以管理員身分查看您的攝影機並更改設定。",
		"frigate_noauth_fix": "將 Frigate 更新至 0.14 或更新版本，停止公開連接埠 5000（從 Docker 的 ports 移除 5000:5000），並使用需要登入的 8971 連接埠。",
	},
	"es": {
		"frigate_noauth":     "El NVR Frigate en %s (puerto 5000) no tiene contraseña. Cualquier persona de su red puede ver sus cámaras y cambiar la configuración como administrador.",
		"frigate_noauth_fix": "Actualice Frigate a la versión 0.14 o posterior, deje de publicar el puerto 5000 (quite 5000:5000 de los puertos de Docker) y use el puerto 8971, que requiere iniciar sesión.",
	},
	"pt": {
		"frigate_noauth":     "O NVR Frigate em %s (porta 5000) não tem senha. Qualquer pessoa na sua rede pode ver as suas câmeras e alterar as configurações como administrador.",
		"frigate_noauth_fix": "Atualize o Frigate para a versão 0.14 ou posterior, pare de publicar a porta 5000 (remova 5000:5000 das portas do Docker) e use a porta 8971, que exige login.",
	},
	"fr": {
		"frigate_noauth":     "Le NVR Frigate sur %s (port 5000) n’a pas de mot de passe. Toute personne sur votre réseau peut regarder vos caméras et modifier ses paramètres en tant qu’administrateur.",
		"frigate_noauth_fix": "Mettez Frigate à jour vers la version 0.14 ou ultérieure, cessez de publier le port 5000 (retirez 5000:5000 des ports Docker) et utilisez le port 8971, qui exige une connexion.",
	},
	"de": {
		"frigate_noauth":     "Der Frigate-NVR auf %s (Port 5000) hat kein Passwort. Jede Person in Ihrem Netzwerk kann Ihre Kameras ansehen und als Administrator die Einstellungen ändern.",
		"frigate_noauth_fix": "Aktualisieren Sie Frigate auf Version 0.14 oder neuer, veröffentlichen Sie Port 5000 nicht mehr (entfernen Sie 5000:5000 aus den Docker-Ports) und verwenden Sie Port 8971, für den eine Anmeldung erforderlich ist.",
	},
	"ru": {
		"frigate_noauth":     "Frigate NVR на устройстве «%s» (порт 5000) не защищён паролем. Любой пользователь вашей сети может смотреть ваши камеры и менять настройки с правами администратора.",
		"frigate_noauth_fix": "Обновите Frigate до версии 0.14 или новее, перестаньте публиковать порт 5000 (удалите 5000:5000 из списка портов Docker) и используйте порт 8971, для которого требуется вход в систему.",
	},
	"ko": {
		"frigate_noauth":     "%s의 Frigate NVR(포트 5000)에 비밀번호가 없습니다. 네트워크의 누구나 관리자 권한으로 카메라를 보고 설정을 바꿀 수 있습니다.",
		"frigate_noauth_fix": "Frigate를 0.14 이상으로 업데이트하고, 포트 5000 공개를 중단한 뒤(Docker ports에서 5000:5000 제거), 로그인이 필요한 포트 8971을 사용하세요.",
	},
	"id": {
		"frigate_noauth":     "NVR Frigate di %s (port 5000) tidak memiliki kata sandi. Siapa pun di jaringan Anda dapat melihat kamera Anda dan mengubah pengaturannya sebagai administrator.",
		"frigate_noauth_fix": "Perbarui Frigate ke versi 0.14 atau lebih baru, hentikan publikasi port 5000 (hapus 5000:5000 dari ports Docker), dan gunakan port 8971 yang memerlukan login.",
	},
}

func init() {
	add(frigateMessages)
}
