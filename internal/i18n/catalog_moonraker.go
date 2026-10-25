package i18n

var moonrakerMessages = map[string]map[string]string{
	"en": {
		"moonraker_noauth":     "Moonraker (the 3D printer control service) on %s accepts commands from this PC without a password. With its common default settings, anyone on your network can heat the printer and start prints.",
		"moonraker_noauth_fix": "In moonraker.conf, under [authorization], remove whole network ranges such as 192.168.0.0/16 from trusted_clients and list only your own devices, or create a user and set force_logins: True. Then restart Moonraker.",
	},
	"ja": {
		"moonraker_noauth":     "「%s」の Moonraker（3D プリンターの操作窓口）が、この PC からパスワードなしで操作できる状態です。よく使われる初めの設定のままだと、家の中の誰でもプリンターを加熱したり、印刷を始めたりできます。",
		"moonraker_noauth_fix": "moonraker.conf の [authorization] にある trusted_clients から 192.168.0.0/16 のような範囲を消して自分の機器だけを書くか、ユーザーを作って force_logins: True を設定し、Moonraker を再起動してください。",
	},
	"zh-Hans": {
		"moonraker_noauth":     "%s 上的 Moonraker（3D 打印机控制服务）允许这台 PC 在无需密码的情况下发送命令。使用常见默认设置时，网络中的任何人都可以加热打印机并开始打印。",
		"moonraker_noauth_fix": "在 moonraker.conf 的 [authorization] 中，请您从 trusted_clients 中删除 192.168.0.0/16 等整个网段，只列出您自己的设备；或者创建用户并设置 force_logins: True。然后重启 Moonraker。",
	},
	"zh-Hant": {
		"moonraker_noauth":     "%s 上的 Moonraker（3D 印表機控制服務）允許這台 PC 在無需密碼的情況下傳送命令。使用常見預設設定時，網路中的任何人都可以加熱印表機並開始列印。",
		"moonraker_noauth_fix": "在 moonraker.conf 的 [authorization] 中，請您從 trusted_clients 中移除 192.168.0.0/16 等整個網段，只列出您自己的裝置；或者建立使用者並設定 force_logins: True。然後重新啟動 Moonraker。",
	},
	"es": {
		"moonraker_noauth":     "Moonraker (el servicio de control de impresoras 3D) en %s acepta comandos desde este PC sin contraseña. Con la configuración inicial habitual, cualquier persona de su red puede calentar la impresora e iniciar impresiones.",
		"moonraker_noauth_fix": "En moonraker.conf, dentro de [authorization], elimine rangos de red completos como 192.168.0.0/16 de trusted_clients y enumere solo sus propios dispositivos, o cree un usuario y establezca force_logins: True. Después, reinicie Moonraker.",
	},
	"pt": {
		"moonraker_noauth":     "O Moonraker (serviço de controle da impressora 3D) em %s aceita comandos deste PC sem senha. Com as configurações iniciais mais comuns, qualquer pessoa na sua rede pode aquecer a impressora e iniciar trabalhos de impressão.",
		"moonraker_noauth_fix": "Em moonraker.conf, na seção [authorization], remova de trusted_clients redes inteiras como 192.168.0.0/16 e liste apenas os seus dispositivos, ou crie um usuário e defina force_logins: True. Em seguida, reinicie o Moonraker.",
	},
	"fr": {
		"moonraker_noauth":     "Moonraker (le service de contrôle d’imprimante 3D) sur %s accepte des commandes de ce PC sans mot de passe. Avec les paramètres initiaux courants, toute personne de votre réseau peut chauffer l’imprimante et lancer des impressions.",
		"moonraker_noauth_fix": "Dans moonraker.conf, sous [authorization], supprimez de trusted_clients les plages réseau complètes comme 192.168.0.0/16 et indiquez uniquement vos appareils, ou créez un utilisateur et définissez force_logins: True. Ensuite, redémarrez Moonraker.",
	},
	"de": {
		"moonraker_noauth":     "Moonraker (der Steuerungsdienst für 3D-Drucker) auf %s nimmt Befehle von diesem PC ohne Passwort an. Mit den häufigen Standardeinstellungen kann jede Person in Ihrem Netzwerk den Drucker aufheizen und Drucke starten.",
		"moonraker_noauth_fix": "Entfernen Sie in moonraker.conf unter [authorization] ganze Netzwerkbereiche wie 192.168.0.0/16 aus trusted_clients und listen Sie nur Ihre eigenen Geräte auf, oder erstellen Sie einen Benutzer und setzen Sie force_logins: True. Starten Sie Moonraker danach neu.",
	},
	"ru": {
		"moonraker_noauth":     "Moonraker (служба управления 3D-принтером) на %s принимает команды с этого ПК без пароля. При распространённых настройках по умолчанию любой пользователь вашей сети может нагреть принтер и запустить печать.",
		"moonraker_noauth_fix": "В moonraker.conf в разделе [authorization] удалите из trusted_clients целые диапазоны сети, например 192.168.0.0/16, и укажите только свои устройства либо создайте пользователя и установите force_logins: True. Затем перезапустите Moonraker.",
	},
	"ko": {
		"moonraker_noauth":     "%s의 Moonraker(3D 프린터 제어 서비스)는 이 PC에서 비밀번호 없이 명령을 받습니다. 흔히 사용하는 기본 설정에서는 네트워크의 누구나 프린터를 가열하고 인쇄를 시작할 수 있습니다.",
		"moonraker_noauth_fix": "moonraker.conf의 [authorization]에서 trusted_clients에 있는 192.168.0.0/16 같은 전체 네트워크 범위를 제거하고 본인 기기만 나열하거나, 사용자를 만들고 force_logins: True로 설정한 다음 Moonraker를 다시 시작하세요.",
	},
	"id": {
		"moonraker_noauth":     "Moonraker (layanan kendali printer 3D) di %s menerima perintah dari PC ini tanpa kata sandi. Dengan pengaturan awal yang umum, siapa pun di jaringan Anda dapat memanaskan printer dan memulai pencetakan.",
		"moonraker_noauth_fix": "Di moonraker.conf, pada bagian [authorization], hapus seluruh rentang jaringan seperti 192.168.0.0/16 dari trusted_clients dan cantumkan hanya perangkat Anda, atau buat pengguna dan atur force_logins: True. Lalu mulai ulang Moonraker.",
	},
}

func init() {
	add(moonrakerMessages)
}
