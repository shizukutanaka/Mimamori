package i18n

var noticeMessages = map[string]map[string]string{
	"en": {
		"incomplete_notice": "Some checks could not run, so this result may be incomplete. Running it again on the home network may help.",
		"scanning_notice":   "Checking your home network. This takes a little while…",
	},
	"ja": {
		"incomplete_notice": "一部の確認ができなかったため、この結果は不完全な可能性があります。自宅のネットワークに接続した状態で、もう一度実行してください。",
		"scanning_notice":   "家のネットワークを確認しています。少し時間がかかります…",
	},
	"zh-Hans": {
		"incomplete_notice": "部分检查未能运行，因此结果可能不完整。请在连接到家庭网络后再运行一次。",
		"scanning_notice":   "正在检查您的家庭网络。这需要一些时间…",
	},
	"zh-Hant": {
		"incomplete_notice": "部分檢查無法執行，因此結果可能不完整。請在連上家中網路後再執行一次。",
		"scanning_notice":   "正在檢查您的家庭網路。這需要一些時間…",
	},
	"es": {
		"incomplete_notice": "No se pudieron realizar algunas comprobaciones, así que el resultado puede estar incompleto. Vuelve a ejecutarlo conectado a la red de casa.",
		"scanning_notice":   "Estamos comprobando tu red doméstica. Esto tardará un poco…",
	},
	"pt": {
		"incomplete_notice": "Não foi possível executar algumas verificações, então o resultado pode estar incompleto. Tente novamente conectado à rede de casa.",
		"scanning_notice":   "Verificando sua rede doméstica. Isso vai levar um pouco…",
	},
	"fr": {
		"incomplete_notice": "Certaines vérifications n’ont pas pu être effectuées ; le résultat peut donc être incomplet. Relancez l’outil lorsque vous êtes connecté au réseau de la maison.",
		"scanning_notice":   "Vérification de votre réseau domestique. Cela prend un peu de temps…",
	},
	"de": {
		"incomplete_notice": "Einige Prüfungen konnten nicht ausgeführt werden. Das Ergebnis ist daher möglicherweise unvollständig. Führen Sie den Check erneut im Heimnetzwerk aus.",
		"scanning_notice":   "Ihr Heimnetzwerk wird überprüft. Das dauert einen Moment…",
	},
	"ru": {
		"incomplete_notice": "Некоторые проверки не удалось выполнить, поэтому результат может быть неполным. Запустите проверку ещё раз, подключившись к домашней сети.",
		"scanning_notice":   "Проверяем домашнюю сеть. Это займёт немного времени…",
	},
	"ko": {
		"incomplete_notice": "일부 검사를 실행하지 못해 결과가 불완전할 수 있습니다. 집 네트워크에 연결한 상태에서 다시 실행해 보세요.",
		"scanning_notice":   "집 네트워크를 확인하고 있습니다. 잠시 시간이 걸립니다…",
	},
	"id": {
		"incomplete_notice": "Beberapa pemeriksaan tidak dapat dijalankan, jadi hasil ini mungkin belum lengkap. Coba jalankan kembali saat terhubung ke jaringan rumah.",
		"scanning_notice":   "Sedang memeriksa jaringan rumah Anda. Proses ini memerlukan sedikit waktu…",
	},
}

var noticeCatalogTargets = map[string]map[string]string{
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
	for lang, catalog := range noticeMessages {
		target := noticeCatalogTargets[lang]
		for key, value := range catalog {
			target[key] = value
		}
		messages[lang] = target
	}
}
