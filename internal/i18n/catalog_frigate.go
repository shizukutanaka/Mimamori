package i18n

var frigateMessages = map[string]map[string]string{
	"en": {
		"frigate_noauth":     "Anyone can control the Frigate NVR (%s) as an administrator without a password (port 5000)",
		"frigate_noauth_fix": "Update Frigate to 0.14 or later, stop publishing port 5000 (remove 5000:5000 from the Docker ports), and use port 8971, which requires a login.",
	},
	"ja": {
		"frigate_noauth":     "録画ソフト Frigate（%s）を、パスワードなしで誰でも管理者として操作できます（5000番）",
		"frigate_noauth_fix": "Frigate を 0.14 以降に更新し、5000番の公開（Docker の ports の 5000:5000）を外して、ログインが必要な 8971番を使ってください。",
	},
	"zh-Hans": {
		"frigate_noauth":     "任何人都可以在没有密码的情况下，以管理员身份控制 Frigate NVR（%s）（端口 5000）。",
		"frigate_noauth_fix": "将 Frigate 更新到 0.14 或更高版本，停止公开端口 5000（从 Docker 的 ports 中移除 5000:5000），并使用需要登录的 8971 端口。",
	},
	"zh-Hant": {
		"frigate_noauth":     "任何人都可以在沒有密碼的情況下，以管理員身分控制 Frigate NVR（%s）（連接埠 5000）。",
		"frigate_noauth_fix": "將 Frigate 更新至 0.14 或更新版本，停止公開連接埠 5000（從 Docker 的 ports 移除 5000:5000），並使用需要登入的 8971 連接埠。",
	},
	"es": {
		"frigate_noauth":     "Cualquier persona puede controlar el NVR Frigate (%s) como administrador sin contraseña (puerto 5000).",
		"frigate_noauth_fix": "Actualice Frigate a la versión 0.14 o posterior, deje de publicar el puerto 5000 (quite 5000:5000 de los puertos de Docker) y use el puerto 8971, que requiere iniciar sesión.",
	},
	"pt": {
		"frigate_noauth":     "Qualquer pessoa pode controlar o NVR Frigate (%s) como administrador sem senha (porta 5000).",
		"frigate_noauth_fix": "Atualize o Frigate para a versão 0.14 ou posterior, pare de publicar a porta 5000 (remova 5000:5000 das portas do Docker) e use a porta 8971, que exige login.",
	},
	"fr": {
		"frigate_noauth":     "Toute personne peut contrôler le NVR Frigate (%s) en tant qu’administrateur sans mot de passe (port 5000).",
		"frigate_noauth_fix": "Mettez Frigate à jour vers la version 0.14 ou ultérieure, cessez de publier le port 5000 (retirez 5000:5000 des ports Docker) et utilisez le port 8971, qui exige une connexion.",
	},
	"de": {
		"frigate_noauth":     "Jede Person kann den Frigate-NVR (%s) ohne Passwort als Administrator steuern (Port 5000).",
		"frigate_noauth_fix": "Aktualisieren Sie Frigate auf Version 0.14 oder neuer, veröffentlichen Sie Port 5000 nicht mehr (entfernen Sie 5000:5000 aus den Docker-Ports) und verwenden Sie Port 8971, für den eine Anmeldung erforderlich ist.",
	},
	"ru": {
		"frigate_noauth":     "Любой может управлять Frigate NVR (%s) с правами администратора без пароля (порт 5000).",
		"frigate_noauth_fix": "Обновите Frigate до версии 0.14 или новее, перестаньте публиковать порт 5000 (удалите 5000:5000 из списка портов Docker) и используйте порт 8971, для которого требуется вход в систему.",
	},
	"ko": {
		"frigate_noauth":     "누구나 비밀번호 없이 Frigate NVR(%s)을 관리자 권한으로 제어할 수 있습니다(포트 5000).",
		"frigate_noauth_fix": "Frigate를 0.14 이상으로 업데이트하고, 포트 5000 공개를 중단한 뒤(Docker ports에서 5000:5000 제거), 로그인이 필요한 포트 8971을 사용하세요.",
	},
	"id": {
		"frigate_noauth":     "Siapa pun dapat mengendalikan NVR Frigate (%s) sebagai administrator tanpa kata sandi (port 5000).",
		"frigate_noauth_fix": "Perbarui Frigate ke versi 0.14 atau lebih baru, hentikan publikasi port 5000 (hapus 5000:5000 dari ports Docker), dan gunakan port 8971 yang memerlukan login.",
	},
}

var frigateCatalogTargets = map[string]map[string]string{
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
	for lang, catalog := range frigateMessages {
		target := frigateCatalogTargets[lang]
		for key, value := range catalog {
			target[key] = value
		}
		messages[lang] = target
	}
}
