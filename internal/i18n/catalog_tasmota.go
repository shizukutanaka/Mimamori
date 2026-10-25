package i18n

var tasmotaMessages = map[string]map[string]string{
	"en": {
		"tasmota_noauth":     "The Tasmota smart device on %s has no web password, so anyone on your network can switch it on or off, run commands, and replace its firmware.",
		"tasmota_noauth_fix": "Open the device's Tasmota web page, go to Configuration → Other, enter a Web Admin Password, and press Save. Or, in the Console, run WebPassword followed by a password.",
	},
	"ja": {
		"tasmota_noauth":     "「%s」の Tasmota の機器に管理画面のパスワードがなく、家の中の誰でも電源を入れたり切ったり、命令を実行させたり、ファームウェアを書き換えたりできる状態です。",
		"tasmota_noauth_fix": "機器の Tasmota の画面を開き、「Configuration」→「Other」の「Web Admin Password」にパスワードを入れて「Save」を押してください。または「Console」で WebPassword に続けてパスワードを入力してください。",
	},
	"zh-Hans": {
		"tasmota_noauth":     "%s 上的 Tasmota 智能设备没有网页密码，因此网络中的任何人都可以开关设备、运行命令并替换固件。",
		"tasmota_noauth_fix": "请您打开设备的 Tasmota 网页，前往 Configuration → Other，在 Web Admin Password 中输入密码并按 Save；或者在 Console 中运行 WebPassword，后面跟上密码。",
	},
	"zh-Hant": {
		"tasmota_noauth":     "%s 上的 Tasmota 智慧裝置沒有網頁密碼，因此網路中的任何人都可以開關裝置、執行命令並替換韌體。",
		"tasmota_noauth_fix": "請您開啟裝置的 Tasmota 網頁，前往 Configuration → Other，在 Web Admin Password 中輸入密碼並按 Save；或者在 Console 中執行 WebPassword，後面接上密碼。",
	},
	"es": {
		"tasmota_noauth":     "El dispositivo inteligente Tasmota en %s no tiene contraseña web, por lo que cualquier persona de su red puede encenderlo o apagarlo, ejecutar comandos y reemplazar su firmware.",
		"tasmota_noauth_fix": "Abra la página web de Tasmota del dispositivo, vaya a Configuration → Other, introduzca una contraseña en Web Admin Password y pulse Save. O, en Console, ejecute WebPassword seguido de una contraseña.",
	},
	"pt": {
		"tasmota_noauth":     "O dispositivo inteligente Tasmota em %s está sem senha na interface web, então qualquer pessoa da sua rede pode ligá-lo ou desligá-lo, executar comandos e substituir o firmware.",
		"tasmota_noauth_fix": "Abra a página web do Tasmota no dispositivo, vá para Configuration → Other, informe uma senha em Web Admin Password e pressione Save. Ou, no Console, execute WebPassword seguido de uma senha.",
	},
	"fr": {
		"tasmota_noauth":     "L’appareil connecté Tasmota sur %s n’a pas de mot de passe Web ; toute personne de votre réseau peut donc l’allumer ou l’éteindre, exécuter des commandes et remplacer son micrologiciel.",
		"tasmota_noauth_fix": "Ouvrez la page Web Tasmota de l’appareil, allez dans Configuration → Other, saisissez un mot de passe dans Web Admin Password et appuyez sur Save. Ou, dans Console, exécutez WebPassword suivi d’un mot de passe.",
	},
	"de": {
		"tasmota_noauth":     "Das Tasmota-Smartgerät auf %s hat kein Web-Passwort. Daher kann jede Person in Ihrem Netzwerk es ein- oder ausschalten, Befehle ausführen und die Firmware ersetzen.",
		"tasmota_noauth_fix": "Öffnen Sie die Tasmota-Webseite des Geräts, wechseln Sie zu Configuration → Other, geben Sie unter Web Admin Password ein Passwort ein und drücken Sie Save. Oder führen Sie in Console den Befehl WebPassword gefolgt von einem Passwort aus.",
	},
	"ru": {
		"tasmota_noauth":     "У устройства Tasmota на %s нет пароля веб-интерфейса, поэтому любой пользователь вашей сети может включать и выключать его, выполнять команды и заменять прошивку.",
		"tasmota_noauth_fix": "Откройте веб-страницу Tasmota устройства, перейдите в Configuration → Other, введите пароль в Web Admin Password и нажмите Save. Или в Console выполните команду WebPassword, указав после неё пароль.",
	},
	"ko": {
		"tasmota_noauth":     "%s의 Tasmota 스마트 기기에는 웹 비밀번호가 없어 네트워크의 누구나 전원을 켜고 끄거나 명령을 실행하고 펌웨어를 바꿀 수 있습니다.",
		"tasmota_noauth_fix": "기기의 Tasmota 웹 페이지를 열고 Configuration → Other로 이동해 Web Admin Password에 비밀번호를 입력한 다음 Save를 누르세요. 또는 Console에서 WebPassword 뒤에 비밀번호를 입력해 실행하세요.",
	},
	"id": {
		"tasmota_noauth":     "Perangkat pintar Tasmota di %s tidak memiliki kata sandi web, sehingga siapa pun di jaringan Anda dapat menyalakan atau mematikannya, menjalankan perintah, dan mengganti firmware-nya.",
		"tasmota_noauth_fix": "Silakan buka halaman web Tasmota perangkat, buka Configuration → Other, masukkan kata sandi di Web Admin Password, lalu tekan Save. Atau, di Console, jalankan WebPassword diikuti kata sandi.",
	},
}

func init() {
	add(tasmotaMessages)
}
