package i18n

var ollamaMessages = map[string]map[string]string{
	"en": {
		"ollama_open":     "Ollama (local AI, port 11434) on %s is open to your network, and Ollama has no password. Anyone at home can run models on this computer's GPU, delete models, and use cloud models on its ollama.com account if it is signed in.",
		"ollama_open_fix": "If no other device needs it, set OLLAMA_HOST to 127.0.0.1:11434 (for Docker, publish with -p 127.0.0.1:11434:11434). In the Windows or Mac app, turn off \"Expose Ollama to the network\" in Settings.",
	},
	"ja": {
		"ollama_open":     "「%s」の Ollama（ローカルの AI、11434番）は、合言葉のしくみがないまま家のネットワークに開いています。家の誰でも、この PC の GPU で AI を動かしたり、モデルを消したりできます。ollama.com にサインインしていれば、その枠のクラウドのモデルも使えます。",
		"ollama_open_fix": "ほかの機器から使わないなら、OLLAMA_HOST を 127.0.0.1:11434 にしてください（Docker なら -p 127.0.0.1:11434:11434）。Windows・Mac のアプリなら、設定の「Expose Ollama to the network」を切ってください。",
	},
	"zh-Hans": {
		"ollama_open":     "%s 上的 Ollama（本地 AI，端口 11434）已向您的网络开放，而且 Ollama 没有密码。家中任何人都可以使用这台电脑的 GPU 运行模型、删除模型；如果账户已登录 ollama.com，还能使用该账户的云端模型。",
		"ollama_open_fix": "如果不需要其他设备访问，请将 OLLAMA_HOST 设为 127.0.0.1:11434（Docker 请使用 -p 127.0.0.1:11434:11434 发布端口）。在 Windows 或 Mac 应用的 Settings 中，请关闭“Expose Ollama to the network”。",
	},
	"zh-Hant": {
		"ollama_open":     "%s 上的 Ollama（本機 AI，連接埠 11434）已對您的網路開放，而且 Ollama 沒有密碼。家中任何人都可以使用這台電腦的 GPU 執行模型、刪除模型；若帳戶已登入 ollama.com，也能使用該帳戶的雲端模型。",
		"ollama_open_fix": "如果不需要其他裝置存取，請將 OLLAMA_HOST 設為 127.0.0.1:11434（Docker 請使用 -p 127.0.0.1:11434:11434 發佈連接埠）。在 Windows 或 Mac 應用程式的 Settings 中，請關閉「Expose Ollama to the network」。",
	},
	"es": {
		"ollama_open":     "Ollama (IA local, puerto 11434) en %s está abierto a tu red y no tiene contraseña. Cualquiera en casa puede ejecutar modelos en la GPU de este equipo, eliminarlos y usar modelos en la nube con su cuenta de ollama.com si ha iniciado sesión.",
		"ollama_open_fix": "Si ningún otro dispositivo lo necesita, configura OLLAMA_HOST como 127.0.0.1:11434 (en Docker, publica con -p 127.0.0.1:11434:11434). En la aplicación de Windows o Mac, desactiva «Expose Ollama to the network» en Settings.",
	},
	"pt": {
		"ollama_open":     "O Ollama (IA local, porta 11434) em %s está aberto à sua rede e não tem senha. Qualquer pessoa em casa pode executar modelos na GPU deste computador, excluí-los e usar modelos na nuvem com a conta ollama.com se ela estiver conectada.",
		"ollama_open_fix": "Se nenhum outro dispositivo precisar dele, defina OLLAMA_HOST como 127.0.0.1:11434 (no Docker, publique com -p 127.0.0.1:11434:11434). No aplicativo para Windows ou Mac, desative “Expose Ollama to the network” em Settings.",
	},
	"fr": {
		"ollama_open":     "Ollama (IA locale, port 11434) sur %s est ouvert à votre réseau et n’a pas de mot de passe. Toute personne à la maison peut exécuter des modèles sur le GPU de cet ordinateur, les supprimer et utiliser des modèles cloud avec son compte ollama.com s’il est connecté.",
		"ollama_open_fix": "Si aucun autre appareil n’en a besoin, définissez OLLAMA_HOST sur 127.0.0.1:11434 (avec Docker, publiez le port avec -p 127.0.0.1:11434:11434). Dans l’application Windows ou Mac, désactivez « Expose Ollama to the network » dans Settings.",
	},
	"de": {
		"ollama_open":     "Ollama (lokale KI, Port 11434) auf %s ist für Ihr Netzwerk offen und hat kein Passwort. Jede Person zu Hause kann Modelle auf der GPU dieses Computers ausführen, löschen und bei Anmeldung mit dem ollama.com-Konto Cloud-Modelle nutzen.",
		"ollama_open_fix": "Wenn kein anderes Gerät Ollama benötigt, setzen Sie OLLAMA_HOST auf 127.0.0.1:11434 (bei Docker veröffentlichen Sie den Port mit -p 127.0.0.1:11434:11434). Schalten Sie in der Windows- oder Mac-App unter Settings „Expose Ollama to the network“ aus.",
	},
	"ru": {
		"ollama_open":     "Ollama (локальный ИИ, порт 11434) на устройстве %s открыт в вашей сети и не защищён паролем. Любой человек дома может запускать модели на GPU этого компьютера, удалять их и использовать облачные модели учётной записи ollama.com, если она авторизована.",
		"ollama_open_fix": "Если другим устройствам Ollama не нужен, задайте OLLAMA_HOST=127.0.0.1:11434 (в Docker опубликуйте порт через -p 127.0.0.1:11434:11434). В приложении для Windows или Mac выключите в разделе Settings переключатель «Expose Ollama to the network».",
	},
	"ko": {
		"ollama_open":     "%s의 Ollama(로컬 AI, 포트 11434)가 네트워크에 열려 있고 비밀번호가 없습니다. 집 안의 누구나 이 컴퓨터의 GPU로 모델을 실행하거나 삭제할 수 있으며, ollama.com 계정에 로그인되어 있으면 클라우드 모델도 사용할 수 있습니다.",
		"ollama_open_fix": "다른 기기에서 사용할 필요가 없다면 OLLAMA_HOST를 127.0.0.1:11434로 설정하세요(Docker에서는 -p 127.0.0.1:11434:11434로 포트를 공개하세요). Windows 또는 Mac 앱에서는 Settings에서 “Expose Ollama to the network”를 끄세요.",
	},
	"id": {
		"ollama_open":     "Ollama (AI lokal, port 11434) di %s terbuka ke jaringan Anda dan tidak memiliki kata sandi. Siapa pun di rumah dapat menjalankan model di GPU komputer ini, menghapus model, dan menggunakan model cloud dengan akun ollama.com jika sedang masuk.",
		"ollama_open_fix": "Jika tidak diperlukan perangkat lain, atur OLLAMA_HOST ke 127.0.0.1:11434 (untuk Docker, publikasikan dengan -p 127.0.0.1:11434:11434). Di aplikasi Windows atau Mac, nonaktifkan “Expose Ollama to the network” di Settings.",
	},
}

var ollamaCatalogTargets = map[string]map[string]string{
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
	for lang, catalog := range ollamaMessages {
		target := ollamaCatalogTargets[lang]
		for key, value := range catalog {
			target[key] = value
		}
		messages[lang] = target
	}
}
