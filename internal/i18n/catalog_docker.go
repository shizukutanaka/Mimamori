package i18n

var dockerMessages = map[string]map[string]string{
	"en": {
		"docker_noauth":           "The Docker remote API (port 2375) on %s has no password. Anyone on your network can run programs on it and take over the whole device.",
		"docker_noauth_fix":       "Stop exposing Docker's unencrypted port 2375. For remote use, use TLS (port 2376) or SSH instead.",
		"docker_tls_noverify":     "The Docker remote API (port 2376) on %s uses TLS but doesn't check client certificates (no --tlsverify). Anyone on your network can run programs on it and take over the whole device.",
		"docker_tls_noverify_fix": "Start Docker with --tlsverify so only clients with a certificate from your CA can connect.",
	},
	"ja": {
		"docker_noauth":           "「%s」の Docker の操作窓口（2375番）に、パスワードがありません。家のネットワークの誰でも、その機器でプログラムを動かし、機器全体を乗っ取れます。",
		"docker_noauth_fix":       "Docker の 2375番（暗号化なし）の公開をやめてください。離れた場所から使うときは、TLS（2376番）か SSH を使ってください。",
		"docker_tls_noverify":     "「%s」の Docker の操作窓口（2376番）は、暗号化されていますが、つないでくる相手を確かめていません（--tlsverify なし）。家のネットワークの誰でも、その機器でプログラムを動かし、機器全体を乗っ取れます。",
		"docker_tls_noverify_fix": "Docker を --tlsverify 付きで起動し、証明書を持つ相手だけがつなげるようにしてください。",
	},
	"zh-Hans": {
		"docker_noauth":           "%s 上的 Docker 远程 API（端口 2375）没有密码。网络中的任何人都可以在该设备上运行程序并完全控制设备。",
		"docker_noauth_fix":       "请停止公开未加密的 Docker 2375 端口。远程使用时，请改用 TLS（端口 2376）或 SSH。",
		"docker_tls_noverify":     "%s 上的 Docker 远程 API（端口 2376）使用了 TLS，但不验证客户端证书（未启用 --tlsverify）。网络中的任何人都可以在该设备上运行程序并完全控制设备。",
		"docker_tls_noverify_fix": "请使用 --tlsverify 启动 Docker，以便只有持有您 CA 签发证书的客户端才能连接。",
	},
	"zh-Hant": {
		"docker_noauth":           "%s 上的 Docker 遠端 API（連接埠 2375）沒有密碼。網路中的任何人都可以在該裝置上執行程式並完全控制裝置。",
		"docker_noauth_fix":       "請停止公開未加密的 Docker 2375 連接埠。遠端使用時，請改用 TLS（連接埠 2376）或 SSH。",
		"docker_tls_noverify":     "%s 上的 Docker 遠端 API（連接埠 2376）使用 TLS，但不驗證用戶端憑證（未啟用 --tlsverify）。網路中的任何人都可以在該裝置上執行程式並完全控制裝置。",
		"docker_tls_noverify_fix": "請使用 --tlsverify 啟動 Docker，讓只有持有您的 CA 簽發憑證的用戶端才能連線。",
	},
	"es": {
		"docker_noauth":           "La API remota de Docker (puerto 2375) en %s no tiene contraseña. Cualquier persona de su red puede ejecutar programas y tomar el control total del dispositivo.",
		"docker_noauth_fix":       "Deje de exponer el puerto 2375 de Docker sin cifrar. Para el acceso remoto, utilice TLS (puerto 2376) o SSH.",
		"docker_tls_noverify":     "La API remota de Docker (puerto 2376) en %s usa TLS, pero no verifica los certificados de los clientes (sin --tlsverify). Cualquier persona de su red puede ejecutar programas y tomar el control total del dispositivo.",
		"docker_tls_noverify_fix": "Inicie Docker con --tlsverify para que solo puedan conectarse los clientes con un certificado de su CA.",
	},
	"pt": {
		"docker_noauth":           "A API remota do Docker (porta 2375) em %s não tem senha. Qualquer pessoa na sua rede pode executar programas e assumir o controle total do dispositivo.",
		"docker_noauth_fix":       "Pare de expor a porta 2375 do Docker sem criptografia. Para uso remoto, use TLS (porta 2376) ou SSH.",
		"docker_tls_noverify":     "A API remota do Docker (porta 2376) em %s usa TLS, mas não verifica os certificados dos clientes (sem --tlsverify). Qualquer pessoa na sua rede pode executar programas e assumir o controle total do dispositivo.",
		"docker_tls_noverify_fix": "Inicie o Docker com --tlsverify para que apenas clientes com um certificado da sua CA possam se conectar.",
	},
	"fr": {
		"docker_noauth":           "L’API distante Docker (port 2375) sur %s n’a pas de mot de passe. Toute personne sur votre réseau peut y exécuter des programmes et prendre le contrôle complet de l’appareil.",
		"docker_noauth_fix":       "Cessez d’exposer le port 2375 non chiffré de Docker. Pour un accès à distance, utilisez TLS (port 2376) ou SSH.",
		"docker_tls_noverify":     "L’API distante Docker (port 2376) sur %s utilise TLS, mais ne vérifie pas les certificats clients (sans --tlsverify). Toute personne sur votre réseau peut y exécuter des programmes et prendre le contrôle complet de l’appareil.",
		"docker_tls_noverify_fix": "Démarrez Docker avec --tlsverify afin que seuls les clients disposant d’un certificat de votre autorité de certification puissent se connecter.",
	},
	"de": {
		"docker_noauth":           "Die Docker-Remote-API (Port 2375) auf %s hat kein Passwort. Jede Person in Ihrem Netzwerk kann dort Programme ausführen und das Gerät vollständig übernehmen.",
		"docker_noauth_fix":       "Geben Sie den unverschlüsselten Docker-Port 2375 nicht mehr frei. Verwenden Sie für den Fernzugriff stattdessen TLS (Port 2376) oder SSH.",
		"docker_tls_noverify":     "Die Docker-Remote-API (Port 2376) auf %s verwendet TLS, überprüft aber keine Clientzertifikate (--tlsverify ist nicht aktiviert). Jede Person in Ihrem Netzwerk kann dort Programme ausführen und das Gerät vollständig übernehmen.",
		"docker_tls_noverify_fix": "Starten Sie Docker mit --tlsverify, damit nur Clients mit einem Zertifikat Ihrer CA eine Verbindung herstellen können.",
	},
	"ru": {
		"docker_noauth":           "У удалённого API Docker (порт 2375) на устройстве %s нет пароля. Любой пользователь вашей сети может запускать на нём программы и полностью захватить устройство.",
		"docker_noauth_fix":       "Перестаньте открывать незашифрованный порт Docker 2375. Для удалённого доступа используйте TLS (порт 2376) или SSH.",
		"docker_tls_noverify":     "Удалённый API Docker (порт 2376) на устройстве %s использует TLS, но не проверяет сертификаты клиентов (--tlsverify не включён). Любой пользователь вашей сети может запускать на нём программы и полностью захватить устройство.",
		"docker_tls_noverify_fix": "Запустите Docker с параметром --tlsverify, чтобы подключаться могли только клиенты с сертификатом вашего центра сертификации.",
	},
	"ko": {
		"docker_noauth":           "%s의 Docker 원격 API(포트 2375)에 비밀번호가 없습니다. 네트워크의 누구나 프로그램을 실행해 기기 전체를 장악할 수 있습니다.",
		"docker_noauth_fix":       "암호화되지 않은 Docker 포트 2375를 더 이상 공개하지 마세요. 원격으로 사용할 때는 TLS(포트 2376) 또는 SSH를 사용하세요.",
		"docker_tls_noverify":     "%s의 Docker 원격 API(포트 2376)는 TLS를 사용하지만 클라이언트 인증서를 확인하지 않습니다(--tlsverify 없음). 네트워크의 누구나 프로그램을 실행해 기기 전체를 장악할 수 있습니다.",
		"docker_tls_noverify_fix": "CA에서 발급한 인증서가 있는 클라이언트만 연결할 수 있도록 Docker를 --tlsverify와 함께 시작하세요.",
	},
	"id": {
		"docker_noauth":           "API jarak jauh Docker (port 2375) pada %s tidak memiliki kata sandi. Siapa pun di jaringan Anda dapat menjalankan program dan mengambil alih seluruh perangkat.",
		"docker_noauth_fix":       "Hentikan paparan port Docker 2375 yang tidak terenkripsi. Untuk penggunaan jarak jauh, gunakan TLS (port 2376) atau SSH.",
		"docker_tls_noverify":     "API jarak jauh Docker (port 2376) pada %s menggunakan TLS, tetapi tidak memeriksa sertifikat klien (tanpa --tlsverify). Siapa pun di jaringan Anda dapat menjalankan program dan mengambil alih seluruh perangkat.",
		"docker_tls_noverify_fix": "Jalankan Docker dengan --tlsverify agar hanya klien yang memiliki sertifikat dari CA Anda yang dapat terhubung.",
	},
}

var dockerCatalogTargets = map[string]map[string]string{
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
	for lang, catalog := range dockerMessages {
		target := dockerCatalogTargets[lang]
		for key, value := range catalog {
			target[key] = value
		}
		messages[lang] = target
	}
}
