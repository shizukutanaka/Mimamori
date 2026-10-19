package i18n

var wledMessages = map[string]map[string]string{
	"en": {
		"wled_noauth":     "The WLED light controller on %s has no settings PIN, so anyone on your network can change its settings and replace its firmware.",
		"wled_noauth_fix": "Open the device's WLED web page, go to Config → Security & Updates, enter a 4-digit Settings PIN, and press Save.",
	},
	"ja": {
		"wled_noauth":     "「%s」の WLED（LED の照明）に設定の PIN がなく、家の中の誰でも設定を変えたり、ファームウェアを書き換えたりできる状態です。",
		"wled_noauth_fix": "WLED の画面を開き、「Config」→「Security & Updates」の「Settings PIN」に4桁の番号を入れて「Save」を押してください。",
	},
	"zh-Hans": {
		"wled_noauth":     "%s 上的 WLED 灯光控制器没有设置 PIN，因此您网络中的任何人都可以更改设置并替换固件。",
		"wled_noauth_fix": "请您打开设备的 WLED 网页，进入 Config → Security & Updates，输入 4 位 Settings PIN，然后按 Save。",
	},
	"zh-Hant": {
		"wled_noauth":     "%s 上的 WLED 燈光控制器沒有設定 PIN，因此您網路中的任何人都可以更改設定並替換韌體。",
		"wled_noauth_fix": "請您開啟裝置的 WLED 網頁，進入 Config → Security & Updates，輸入 4 位 Settings PIN，然後按 Save。",
	},
	"es": {
		"wled_noauth":     "El controlador de iluminación WLED en %s no tiene un PIN de configuración, por lo que cualquier persona de su red puede cambiar sus ajustes y reemplazar su firmware.",
		"wled_noauth_fix": "Abra la página web WLED del dispositivo, vaya a Config → Security & Updates, introduzca un Settings PIN de 4 dígitos y pulse Save.",
	},
	"pt": {
		"wled_noauth":     "O controlador de iluminação WLED em %s não tem um PIN de definições, por isso qualquer pessoa da sua rede pode alterar as definições e substituir o firmware.",
		"wled_noauth_fix": "Abra a página Web WLED do dispositivo, aceda a Config → Security & Updates, introduza um Settings PIN de 4 dígitos e prima Save.",
	},
	"fr": {
		"wled_noauth":     "Le contrôleur lumineux WLED sur %s n’a pas de code PIN pour les paramètres ; toute personne de votre réseau peut donc modifier ses réglages et remplacer son micrologiciel.",
		"wled_noauth_fix": "Ouvrez la page Web WLED de l’appareil, accédez à Config → Security & Updates, saisissez un Settings PIN à 4 chiffres, puis appuyez sur Save.",
	},
	"de": {
		"wled_noauth":     "Der WLED-Lichtcontroller auf %s hat keine Einstellungs-PIN. Daher kann jede Person in Ihrem Netzwerk die Einstellungen ändern und die Firmware ersetzen.",
		"wled_noauth_fix": "Öffnen Sie die WLED-Webseite des Geräts, gehen Sie zu Config → Security & Updates, geben Sie eine vierstellige Settings PIN ein und drücken Sie Save.",
	},
	"ru": {
		"wled_noauth":     "У контроллера освещения WLED на %s не задан PIN-код настроек, поэтому любой пользователь вашей сети может менять настройки и заменять прошивку.",
		"wled_noauth_fix": "Откройте веб-страницу WLED устройства, перейдите в Config → Security & Updates, установите 4-значный Settings PIN и нажмите Save.",
	},
	"ko": {
		"wled_noauth":     "%s의 WLED 조명 컨트롤러에는 설정 PIN이 없어 네트워크의 누구나 설정을 변경하거나 펌웨어를 교체할 수 있습니다.",
		"wled_noauth_fix": "기기의 WLED 웹 페이지를 열고 Config → Security & Updates로 이동해 4자리 Settings PIN을 입력한 뒤 Save를 누르세요.",
	},
	"id": {
		"wled_noauth":     "Pengontrol lampu WLED di %s tidak memiliki PIN pengaturan, sehingga siapa pun di jaringan Anda dapat mengubah pengaturan dan mengganti firmware.",
		"wled_noauth_fix": "Silakan buka halaman web WLED perangkat, buka Config → Security & Updates, masukkan Settings PIN 4 digit, lalu tekan Save.",
	},
}

var wledCatalogTargets = map[string]map[string]string{
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
	for lang, catalog := range wledMessages {
		target := wledCatalogTargets[lang]
		for key, value := range catalog {
			target[key] = value
		}
		messages[lang] = target
	}
}
