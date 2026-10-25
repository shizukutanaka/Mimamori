package i18n

var vncMessages = map[string]map[string]string{
	"en": {
		"vnc_noauth":     "The screen sharing (VNC) on %s lets anyone on your network see and control its screen without a password.",
		"vnc_noauth_fix": "In the device's VNC or screen sharing settings, turn on a password, or turn VNC off if you don't use it.",
	},
	"ja": {
		"vnc_noauth":     "「%s」の画面共有（VNC）は、家の中の誰でもパスワードなしで画面を見たり操作したりできる状態です。",
		"vnc_noauth_fix": "機器の VNC（画面共有）の設定でパスワードを有効にしてください。使っていない場合は VNC 自体をオフにしてください。",
	},
	"zh-Hans": {
		"vnc_noauth":     "“%s”上的屏幕共享（VNC）允许您家中网络上的任何人无需密码查看并控制其屏幕。",
		"vnc_noauth_fix": "请在设备的 VNC 或屏幕共享设置中启用密码；如果您不使用 VNC，请将其关闭。",
	},
	"zh-Hant": {
		"vnc_noauth":     "「%s」上的螢幕分享（VNC）允許您家中網路上的任何人無需密碼查看並控制其螢幕。",
		"vnc_noauth_fix": "請在裝置的 VNC 或螢幕分享設定中啟用密碼；如果您不使用 VNC，請將其關閉。",
	},
	"es": {
		"vnc_noauth":     "El uso compartido de pantalla (VNC) de %s permite que cualquier persona de su red vea y controle su pantalla sin contraseña.",
		"vnc_noauth_fix": "En la configuración de VNC o de uso compartido de pantalla del dispositivo, active una contraseña o desactive VNC si no lo utiliza.",
	},
	"pt": {
		"vnc_noauth":     "O compartilhamento de tela (VNC) em %s permite que qualquer pessoa na sua rede veja e controle a tela sem senha.",
		"vnc_noauth_fix": "Nas configurações de VNC ou de compartilhamento de tela do dispositivo, ative uma senha ou desative o VNC se você não o usa.",
	},
	"fr": {
		"vnc_noauth":     "Le partage d’écran (VNC) sur %s permet à toute personne de votre réseau de voir et de contrôler son écran sans mot de passe.",
		"vnc_noauth_fix": "Dans les réglages VNC ou de partage d’écran de l’appareil, activez un mot de passe ou désactivez VNC si vous ne l’utilisez pas.",
	},
	"de": {
		"vnc_noauth":     "Die Bildschirmfreigabe (VNC) auf %s ermöglicht allen in Ihrem Netzwerk, den Bildschirm ohne Passwort zu sehen und zu steuern.",
		"vnc_noauth_fix": "Aktivieren Sie in den VNC- oder Bildschirmfreigabeeinstellungen des Geräts ein Passwort oder schalten Sie VNC aus, wenn Sie es nicht nutzen.",
	},
	"ru": {
		"vnc_noauth":     "Демонстрация экрана (VNC) на %s позволяет любому пользователю в вашей сети видеть экран устройства и управлять им без пароля.",
		"vnc_noauth_fix": "Включите пароль в настройках VNC или демонстрации экрана устройства либо отключите VNC, если вы им не пользуетесь.",
	},
	"ko": {
		"vnc_noauth":     "%s의 화면 공유(VNC)는 네트워크의 누구나 비밀번호 없이 화면을 보고 조작할 수 있게 합니다.",
		"vnc_noauth_fix": "기기의 VNC 또는 화면 공유 설정에서 비밀번호를 사용하도록 설정하세요. 사용하지 않는다면 VNC를 끄세요.",
	},
	"id": {
		"vnc_noauth":     "Berbagi layar (VNC) di %s memungkinkan siapa pun di jaringan Anda melihat dan mengendalikan layar perangkat tanpa kata sandi.",
		"vnc_noauth_fix": "Aktifkan kata sandi di pengaturan VNC atau berbagi layar perangkat, atau matikan VNC jika Anda tidak menggunakannya.",
	},
}

func init() {
	add(vncMessages)
}
