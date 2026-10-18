package i18n

var rtspMessages = map[string]map[string]string{
	"en": {
		"rtsp_noauth":     "The camera feed on %s can be viewed without a password (RTSP).",
		"rtsp_noauth_fix": "Turn on RTSP authentication (a password) in the camera's settings, or turn RTSP off if you don't use it, and update the firmware.",
	},
	"ja": {
		"rtsp_noauth":     "%s のカメラ映像が、パスワードなしで見られる状態です（RTSP）。",
		"rtsp_noauth_fix": "カメラの設定で RTSP の認証（パスワード）を有効にするか、使っていなければ RTSP を無効にし、ファームウェアを更新してください。",
	},
	"zh-Hans": {
		"rtsp_noauth":     "摄像头 %s 的 RTSP 画面无需密码即可观看。",
		"rtsp_noauth_fix": "请在摄像头设置中启用 RTSP 身份验证（密码），或在不使用 RTSP 时将其关闭，并更新固件。",
	},
	"zh-Hant": {
		"rtsp_noauth":     "攝影機 %s 的 RTSP 影像無需密碼即可觀看。",
		"rtsp_noauth_fix": "請在攝影機設定中啟用 RTSP 驗證（密碼），或在不使用 RTSP 時將其關閉，並更新韌體。",
	},
	"es": {
		"rtsp_noauth":     "La transmisión de la cámara %s puede verse sin contraseña (RTSP).",
		"rtsp_noauth_fix": "Active la autenticación de RTSP (con contraseña) en la configuración de la cámara o desactive RTSP si no lo utiliza, y actualice el firmware.",
	},
	"pt": {
		"rtsp_noauth":     "A transmissão da câmera %s pode ser vista sem senha (RTSP).",
		"rtsp_noauth_fix": "Ative a autenticação RTSP (com senha) nas configurações da câmera ou desative o RTSP se você não o usa, e atualize o firmware.",
	},
	"fr": {
		"rtsp_noauth":     "Le flux vidéo de la caméra %s est accessible sans mot de passe (RTSP).",
		"rtsp_noauth_fix": "Activez l’authentification RTSP (avec un mot de passe) dans les paramètres de la caméra ou désactivez RTSP si vous ne l’utilisez pas, puis mettez le micrologiciel à jour.",
	},
	"de": {
		"rtsp_noauth":     "Der Kamerastream von %s ist ohne Passwort abrufbar (RTSP).",
		"rtsp_noauth_fix": "Aktivieren Sie die RTSP-Authentifizierung (mit Passwort) in den Kameraeinstellungen oder deaktivieren Sie RTSP, wenn Sie es nicht verwenden, und aktualisieren Sie die Firmware.",
	},
	"ru": {
		"rtsp_noauth":     "Видеопоток камеры %s доступен без пароля (RTSP).",
		"rtsp_noauth_fix": "Включите аутентификацию RTSP (пароль) в настройках камеры или отключите RTSP, если не используете его, и обновите прошивку.",
	},
	"ko": {
		"rtsp_noauth":     "카메라 %s의 RTSP 영상이 비밀번호 없이 시청 가능한 상태입니다.",
		"rtsp_noauth_fix": "카메라 설정에서 RTSP 인증(비밀번호)을 켜거나 사용하지 않는다면 RTSP를 끄고 펌웨어를 업데이트하세요.",
	},
	"id": {
		"rtsp_noauth":     "Siaran kamera %s dapat dilihat tanpa kata sandi (RTSP).",
		"rtsp_noauth_fix": "Aktifkan autentikasi RTSP (kata sandi) di pengaturan kamera, atau matikan RTSP jika Anda tidak menggunakannya, lalu perbarui firmware.",
	},
}

var rtspCatalogTargets = map[string]map[string]string{
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
	for lang, catalog := range rtspMessages {
		target := rtspCatalogTargets[lang]
		for key, value := range catalog {
			target[key] = value
		}
		messages[lang] = target
	}
}
