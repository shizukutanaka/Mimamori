package i18n

var dockerMessages = map[string]map[string]string{
	"en": {
		"docker_noauth":     "The Docker remote API (port 2375) on %s has no password. Anyone on your network can run programs on it and take over the whole device.",
		"docker_noauth_fix": "Stop exposing Docker's unencrypted port 2375. For remote use, use TLS (port 2376) or SSH instead.",
	},
	"ja": {
		"docker_noauth":     "「%s」の Docker の操作窓口（2375番）に、パスワードがありません。家のネットワークの誰でも、その機器でプログラムを動かし、機器全体を乗っ取れます。",
		"docker_noauth_fix": "Docker の 2375番（暗号化なし）の公開をやめてください。離れた場所から使うときは、TLS（2376番）か SSH を使ってください。",
	},
	"zh-Hans": {
		"docker_noauth":     "%s 上的 Docker 远程 API（端口 2375）没有密码。网络中的任何人都可以在该设备上运行程序并完全控制设备。",
		"docker_noauth_fix": "请停止公开未加密的 Docker 2375 端口。远程使用时，请改用 TLS（端口 2376）或 SSH。",
	},
	"zh-Hant": {
		"docker_noauth":     "%s 上的 Docker 遠端 API（連接埠 2375）沒有密碼。網路中的任何人都可以在該裝置上執行程式並完全控制裝置。",
		"docker_noauth_fix": "請停止公開未加密的 Docker 2375 連接埠。遠端使用時，請改用 TLS（連接埠 2376）或 SSH。",
	},
	"es": {
		"docker_noauth":     "La API remota de Docker (puerto 2375) en %s no tiene contraseña. Cualquier persona de tu red puede ejecutar programas y tomar el control total del dispositivo.",
		"docker_noauth_fix": "Deje de exponer el puerto 2375 de Docker sin cifrar. Para el acceso remoto, utilice TLS (puerto 2376) o SSH.",
	},
	"pt": {
		"docker_noauth":     "A API remota do Docker (porta 2375) em %s não tem senha. Qualquer pessoa na sua rede pode executar programas e assumir o controle total do dispositivo.",
		"docker_noauth_fix": "Pare de expor a porta 2375 do Docker sem criptografia. Para uso remoto, use TLS (porta 2376) ou SSH.",
	},
	"fr": {
		"docker_noauth":     "L’API distante Docker (port 2375) sur %s n’a pas de mot de passe. Toute personne sur votre réseau peut y exécuter des programmes et prendre le contrôle complet de l’appareil.",
		"docker_noauth_fix": "Cessez d’exposer le port 2375 non chiffré de Docker. Pour un accès à distance, utilisez TLS (port 2376) ou SSH.",
	},
	"de": {
		"docker_noauth":     "Die Docker-Remote-API (Port 2375) auf %s hat kein Passwort. Jede Person in Ihrem Netzwerk kann dort Programme ausführen und das Gerät vollständig übernehmen.",
		"docker_noauth_fix": "Geben Sie den unverschlüsselten Docker-Port 2375 nicht mehr frei. Verwenden Sie für den Fernzugriff stattdessen TLS (Port 2376) oder SSH.",
	},
	"ru": {
		"docker_noauth":     "У удалённого API Docker (порт 2375) на устройстве %s нет пароля. Любой пользователь вашей сети может запускать на нём программы и полностью захватить устройство.",
		"docker_noauth_fix": "Перестаньте открывать незашифрованный порт Docker 2375. Для удалённого доступа используйте TLS (порт 2376) или SSH.",
	},
	"ko": {
		"docker_noauth":     "%s의 Docker 원격 API(포트 2375)에 비밀번호가 없습니다. 네트워크의 누구나 프로그램을 실행해 기기 전체를 장악할 수 있습니다.",
		"docker_noauth_fix": "암호화되지 않은 Docker 포트 2375를 더 이상 공개하지 마세요. 원격으로 사용할 때는 TLS(포트 2376) 또는 SSH를 사용하세요.",
	},
	"id": {
		"docker_noauth":     "API jarak jauh Docker (port 2375) pada %s tidak memiliki kata sandi. Siapa pun di jaringan Anda dapat menjalankan program dan mengambil alih seluruh perangkat.",
		"docker_noauth_fix": "Hentikan paparan port Docker 2375 yang tidak terenkripsi. Untuk penggunaan jarak jauh, gunakan TLS (port 2376) atau SSH.",
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
