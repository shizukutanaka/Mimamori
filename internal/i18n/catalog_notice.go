package i18n

var noticeMessages = map[string]map[string]string{
	"en": {
		"incomplete_notice": "Some checks could not run, so this result may be incomplete. Running it again on the home network may help.",
	},
	"ja": {
		"incomplete_notice": "一部の確認ができなかったため、この結果は不完全な可能性があります。自宅のネットワークに接続した状態で、もう一度実行してください。",
	},
	"zh-Hans": {
		"incomplete_notice": "部分检查未能运行，因此结果可能不完整。请在连接到家庭网络后再运行一次。",
	},
	"zh-Hant": {
		"incomplete_notice": "部分檢查無法執行，因此結果可能不完整。請在連上家中網路後再執行一次。",
	},
	"es": {
		"incomplete_notice": "No se pudieron realizar algunas comprobaciones, así que el resultado puede estar incompleto. Vuelve a ejecutarlo conectado a la red de casa.",
	},
	"pt": {
		"incomplete_notice": "Não foi possível executar algumas verificações, então o resultado pode estar incompleto. Tente novamente conectado à rede de casa.",
	},
	"fr": {
		"incomplete_notice": "Certaines vérifications n’ont pas pu être effectuées ; le résultat peut donc être incomplet. Relancez l’outil lorsque vous êtes connecté au réseau de la maison.",
	},
	"de": {
		"incomplete_notice": "Einige Prüfungen konnten nicht ausgeführt werden. Das Ergebnis ist daher möglicherweise unvollständig. Führen Sie den Check erneut im Heimnetzwerk aus.",
	},
	"ru": {
		"incomplete_notice": "Некоторые проверки не удалось выполнить, поэтому результат может быть неполным. Запустите проверку ещё раз, подключившись к домашней сети.",
	},
	"ko": {
		"incomplete_notice": "일부 검사를 실행하지 못해 결과가 불완전할 수 있습니다. 집 네트워크에 연결한 상태에서 다시 실행해 보세요.",
	},
	"id": {
		"incomplete_notice": "Beberapa pemeriksaan tidak dapat dijalankan, jadi hasil ini mungkin belum lengkap. Coba jalankan kembali saat terhubung ke jaringan rumah.",
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
