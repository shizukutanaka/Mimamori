package i18n

var wifiMessages = map[string]map[string]string{
	"en": {
		"wifi_open":     "This PC is connected to a Wi-Fi network with no password. People nearby can join it and may see what you send.",
		"wifi_open_fix": "Set a Wi-Fi password in the router settings (choose WPA2 or WPA3), or use a different network.",
		"wifi_wep":      "This Wi-Fi uses WEP, an old protection that can be broken in minutes.",
		"wifi_wep_fix":  "In the router settings, switch the Wi-Fi security to WPA2 or WPA3 and set a new password.",
		"wifi_weak":     "This Wi-Fi uses old WPA/TKIP protection, which is weaker than current standards.",
		"wifi_weak_fix": "In the router settings, choose WPA2 (AES) or WPA3. If an old device can't connect afterwards, consider replacing it.",
	},
	"ja": {
		"wifi_open":     "この PC はパスワードのない Wi-Fi に接続しています。近くの人が接続し、送信内容を見られる可能性があります。",
		"wifi_open_fix": "ルーターの設定で WPA2 または WPA3 を選び、Wi-Fi にパスワードを設定してください。別のネットワークを使う方法もあります。",
		"wifi_wep":      "この Wi-Fi は、数分で破られることがある古い WEP 方式を使っています。",
		"wifi_wep_fix":  "ルーターの設定で Wi-Fi の安全方式を WPA2 または WPA3 に切り替え、新しいパスワードを設定してください。",
		"wifi_weak":     "この Wi-Fi は古い WPA/TKIP 方式を使っており、現在の方式より安全性が低い状態です。",
		"wifi_weak_fix": "ルーターの設定で WPA2（AES）または WPA3 を選んでください。その後つながらない古い機器があれば、買い替えを検討してください。",
	},
	"zh-Hans": {
		"wifi_open":     "这台电脑连接的 Wi-Fi 没有密码。附近的人可以加入网络，也可能看到你发送的内容。",
		"wifi_open_fix": "在路由器设置中为 Wi-Fi 设置密码（选择 WPA2 或 WPA3），或改用其他网络。",
		"wifi_wep":      "这个 Wi-Fi 使用旧的 WEP 保护，几分钟内就可能被破解。",
		"wifi_wep_fix":  "在路由器设置中将 Wi-Fi 安全方式改为 WPA2 或 WPA3，并设置新密码。",
		"wifi_weak":     "这个 Wi-Fi 使用旧的 WPA/TKIP 保护，安全性低于当前标准。",
		"wifi_weak_fix": "在路由器设置中选择 WPA2（AES）或 WPA3。如果旧设备无法连接，可考虑更换设备。",
	},
	"zh-Hant": {
		"wifi_open":     "這台電腦連線的 Wi-Fi 沒有密碼。附近的人可以加入網路，也可能看到你傳送的內容。",
		"wifi_open_fix": "在路由器設定中為 Wi-Fi 設定密碼（選擇 WPA2 或 WPA3），或改用其他網路。",
		"wifi_wep":      "這個 Wi-Fi 使用舊的 WEP 保護，幾分鐘內就可能被破解。",
		"wifi_wep_fix":  "在路由器設定中將 Wi-Fi 安全方式改為 WPA2 或 WPA3，並設定新密碼。",
		"wifi_weak":     "這個 Wi-Fi 使用舊的 WPA/TKIP 保護，安全性低於目前標準。",
		"wifi_weak_fix": "在路由器設定中選擇 WPA2（AES）或 WPA3。如果舊裝置無法連線，可考慮更換裝置。",
	},
	"es": {
		"wifi_open":     "Este PC está conectado a una red Wi-Fi sin contraseña. Las personas cercanas pueden conectarse y quizá ver lo que envías.",
		"wifi_open_fix": "Pon una contraseña de Wi-Fi en la configuración del router (elige WPA2 o WPA3) o usa otra red.",
		"wifi_wep":      "Esta red Wi-Fi usa WEP, una protección antigua que se puede romper en pocos minutos.",
		"wifi_wep_fix":  "En la configuración del router, cambia la seguridad Wi-Fi a WPA2 o WPA3 y pon una contraseña nueva.",
		"wifi_weak":     "Esta red Wi-Fi usa la antigua protección WPA/TKIP, más débil que los estándares actuales.",
		"wifi_weak_fix": "En la configuración del router, elige WPA2 (AES) o WPA3. Si algún dispositivo antiguo deja de conectarse, considera cambiarlo.",
	},
	"pt": {
		"wifi_open":     "Este PC está ligado a uma rede Wi-Fi sem senha. Pessoas por perto podem entrar e talvez ver o que você envia.",
		"wifi_open_fix": "Defina uma senha para o Wi-Fi nas configurações do roteador (escolha WPA2 ou WPA3) ou use outra rede.",
		"wifi_wep":      "Esta rede Wi-Fi usa WEP, uma proteção antiga que pode ser quebrada em poucos minutos.",
		"wifi_wep_fix":  "Nas configurações do roteador, mude a segurança do Wi-Fi para WPA2 ou WPA3 e defina uma nova senha.",
		"wifi_weak":     "Esta rede Wi-Fi usa a antiga proteção WPA/TKIP, mais fraca que os padrões atuais.",
		"wifi_weak_fix": "Nas configurações do roteador, escolha WPA2 (AES) ou WPA3. Se algum aparelho antigo não conectar depois, considere trocá-lo.",
	},
	"fr": {
		"wifi_open":     "Ce PC est connecté à un réseau Wi-Fi sans mot de passe. Les personnes à proximité peuvent s'y connecter et voir ce que vous envoyez.",
		"wifi_open_fix": "Définissez un mot de passe Wi-Fi dans les réglages du routeur (choisissez WPA2 ou WPA3), ou utilisez un autre réseau.",
		"wifi_wep":      "Ce réseau Wi-Fi utilise WEP, une ancienne protection qui peut être cassée en quelques minutes.",
		"wifi_wep_fix":  "Dans les réglages du routeur, passez la sécurité Wi-Fi à WPA2 ou WPA3 et définissez un nouveau mot de passe.",
		"wifi_weak":     "Ce réseau Wi-Fi utilise l'ancienne protection WPA/TKIP, moins sûre que les normes actuelles.",
		"wifi_weak_fix": "Dans les réglages du routeur, choisissez WPA2 (AES) ou WPA3. Si un ancien appareil ne peut plus se connecter, envisagez de le remplacer.",
	},
	"de": {
		"wifi_open":     "Dieser PC ist mit einem WLAN ohne Passwort verbunden. Menschen in der Nähe können sich verbinden und möglicherweise Ihre Daten mitlesen.",
		"wifi_open_fix": "Richten Sie in den Router-Einstellungen ein WLAN-Passwort ein (WPA2 oder WPA3), oder nutzen Sie ein anderes Netzwerk.",
		"wifi_wep":      "Dieses WLAN nutzt WEP, einen alten Schutz, der in wenigen Minuten geknackt werden kann.",
		"wifi_wep_fix":  "Stellen Sie in den Router-Einstellungen die WLAN-Sicherheit auf WPA2 oder WPA3 um und vergeben Sie ein neues Passwort.",
		"wifi_weak":     "Dieses WLAN nutzt den alten WPA/TKIP-Schutz, der schwächer als heutige Standards ist.",
		"wifi_weak_fix": "Wählen Sie in den Router-Einstellungen WPA2 (AES) oder WPA3. Falls sich ein altes Gerät danach nicht mehr verbinden kann, sollten Sie es ersetzen.",
	},
	"ru": {
		"wifi_open":     "Этот компьютер подключён к сети Wi-Fi без пароля. Люди поблизости могут подключиться к ней и, возможно, видеть ваши данные.",
		"wifi_open_fix": "Задайте пароль Wi-Fi в настройках роутера (выберите WPA2 или WPA3) либо используйте другую сеть.",
		"wifi_wep":      "Эта сеть Wi-Fi использует старую защиту WEP, которую можно взломать за несколько минут.",
		"wifi_wep_fix":  "В настройках роутера переключите защиту Wi-Fi на WPA2 или WPA3 и задайте новый пароль.",
		"wifi_weak":     "Эта сеть Wi-Fi использует старую защиту WPA/TKIP, которая слабее современных стандартов.",
		"wifi_weak_fix": "В настройках роутера выберите WPA2 (AES) или WPA3. Если старое устройство перестанет подключаться, подумайте о его замене.",
	},
	"ko": {
		"wifi_open":     "이 PC는 비밀번호가 없는 Wi-Fi에 연결되어 있습니다. 주변 사람이 네트워크에 접속해 전송 내용을 볼 수도 있습니다.",
		"wifi_open_fix": "라우터 설정에서 WPA2 또는 WPA3를 선택하고 Wi-Fi 비밀번호를 설정하거나 다른 네트워크를 사용하세요.",
		"wifi_wep":      "이 Wi-Fi는 몇 분 안에 뚫릴 수 있는 오래된 WEP 보호 방식을 사용합니다.",
		"wifi_wep_fix":  "라우터 설정에서 Wi-Fi 보안을 WPA2 또는 WPA3로 바꾸고 새 비밀번호를 설정하세요.",
		"wifi_weak":     "이 Wi-Fi는 현재 표준보다 약한 오래된 WPA/TKIP 보호 방식을 사용합니다.",
		"wifi_weak_fix": "라우터 설정에서 WPA2(AES) 또는 WPA3를 선택하세요. 이후 오래된 기기가 연결되지 않으면 교체를 고려하세요.",
	},
	"id": {
		"wifi_open":     "PC ini terhubung ke jaringan Wi-Fi tanpa kata sandi. Orang di dekatnya dapat bergabung dan mungkin melihat data yang Anda kirim.",
		"wifi_open_fix": "Atur kata sandi Wi-Fi di pengaturan router (pilih WPA2 atau WPA3), atau gunakan jaringan lain.",
		"wifi_wep":      "Wi-Fi ini memakai perlindungan WEP lama yang dapat dibobol dalam beberapa menit.",
		"wifi_wep_fix":  "Di pengaturan router, ubah keamanan Wi-Fi ke WPA2 atau WPA3 dan buat kata sandi baru.",
		"wifi_weak":     "Wi-Fi ini memakai perlindungan WPA/TKIP lama yang lebih lemah daripada standar saat ini.",
		"wifi_weak_fix": "Di pengaturan router, pilih WPA2 (AES) atau WPA3. Jika perangkat lama tidak bisa tersambung setelahnya, pertimbangkan untuk menggantinya.",
	},
}

var wifiLocationMessages = map[string]map[string]string{
	"en": {
		"wifi_location_notice": "Wi-Fi security could not be checked because Windows location access is off. To check it, turn on Settings > Privacy & security > Location and run again. Mimamori does not use or send your location.",
	},
	"ja": {
		"wifi_location_notice": "Windows の位置情報がオフのため、Wi-Fi の暗号化を確認できませんでした。確認するには「設定 > プライバシーとセキュリティ > 位置情報」をオンにして、もう一度実行してください。Mimamori が位置情報を使ったり送信したりすることはありません。",
	},
	"zh-Hans": {
		"wifi_location_notice": "由于 Windows 的位置访问已关闭，无法检查 Wi-Fi 安全性。要进行检查，请打开“设置 > 隐私和安全性 > 位置”，然后重新运行。Mimamori 不会使用或发送您的位置。",
	},
	"zh-Hant": {
		"wifi_location_notice": "由於 Windows 的位置存取已關閉，無法檢查 Wi-Fi 安全性。若要進行檢查，請開啟「設定 > 隱私權與安全性 > 位置」，然後重新執行。Mimamori 不會使用或傳送您的位置。",
	},
	"es": {
		"wifi_location_notice": "No se pudo comprobar la seguridad de la red Wi-Fi porque el acceso a la ubicación de Windows está desactivado. Para comprobarla, activa la ubicación en Configuración > Privacidad y seguridad > Ubicación y vuelve a ejecutar la herramienta. Mimamori no usa ni envía tu ubicación.",
	},
	"pt": {
		"wifi_location_notice": "Não foi possível verificar a segurança do Wi-Fi porque o acesso à localização do Windows está desativado. Para verificar, ative a localização em Configurações > Privacidade e segurança > Localização e execute novamente. O Mimamori não usa nem envia sua localização.",
	},
	"fr": {
		"wifi_location_notice": "La sécurité Wi-Fi n'a pas pu être vérifiée car l'accès à la localisation de Windows est désactivé. Pour la vérifier, activez la localisation dans Paramètres > Confidentialité et sécurité > Localisation, puis relancez l'outil. Mimamori n'utilise pas et n'envoie pas votre localisation.",
	},
	"de": {
		"wifi_location_notice": "Die WLAN-Sicherheit konnte nicht überprüft werden, weil der Windows-Standortzugriff deaktiviert ist. Aktivieren Sie zur Prüfung den Standort unter Einstellungen > Datenschutz & Sicherheit > Standort und starten Sie das Programm erneut. Mimamori verwendet oder sendet Ihren Standort nicht.",
	},
	"ru": {
		"wifi_location_notice": "Не удалось проверить безопасность Wi-Fi, потому что доступ к местоположению Windows отключён. Чтобы выполнить проверку, включите местоположение в разделе Параметры > Конфиденциальность и безопасность > Расположение и запустите программу ещё раз. Mimamori не использует и не отправляет данные о вашем местоположении.",
	},
	"ko": {
		"wifi_location_notice": "Windows 위치 액세스가 꺼져 있어 Wi-Fi 보안을 확인할 수 없습니다. 확인하려면 설정 > 개인 정보 및 보안 > 위치에서 위치를 켜고 다시 실행하세요. Mimamori는 위치 정보를 사용하거나 전송하지 않습니다.",
	},
	"id": {
		"wifi_location_notice": "Keamanan Wi-Fi tidak dapat diperiksa karena akses lokasi Windows dinonaktifkan. Untuk memeriksanya, aktifkan lokasi di Pengaturan > Privasi & keamanan > Lokasi, lalu jalankan kembali. Mimamori tidak menggunakan atau mengirim lokasi Anda.",
	},
}

var wifiCatalogTargets = map[string]map[string]string{
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
	for lang, catalog := range wifiMessages {
		target := wifiCatalogTargets[lang]
		for key, value := range catalog {
			target[key] = value
		}
		messages[lang] = target
	}
	for lang, catalog := range wifiLocationMessages {
		target := wifiCatalogTargets[lang]
		for key, value := range catalog {
			target[key] = value
		}
		messages[lang] = target
	}
}
