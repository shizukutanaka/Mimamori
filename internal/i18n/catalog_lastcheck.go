package i18n

var lastCheckMessages = map[string]map[string]string{
	"en": {
		"last_check_today": "Last check: today",
		"last_check_day":   "Last check: 1 day ago",
		"last_check_days":  "Last check: %d days ago",
	},
	"ja": {
		"last_check_today": "前回の確認: 今日",
		"last_check_day":   "前回の確認: 1日前",
		"last_check_days":  "前回の確認: %d日前",
	},
	"zh-Hans": {
		"last_check_today": "上次检查：今天",
		"last_check_day":   "上次检查：1 天前",
		"last_check_days":  "上次检查：%d 天前",
	},
	"zh-Hant": {
		"last_check_today": "上次檢查：今天",
		"last_check_day":   "上次檢查：1 天前",
		"last_check_days":  "上次檢查：%d 天前",
	},
	"es": {
		"last_check_today": "Última comprobación: hoy",
		"last_check_day":   "Última comprobación: hace 1 día",
		"last_check_days":  "Última comprobación: hace %d días",
	},
	"pt": {
		"last_check_today": "Última verificação: hoje",
		"last_check_day":   "Última verificação: há 1 dia",
		"last_check_days":  "Última verificação: há %d dias",
	},
	"fr": {
		"last_check_today": "Dernière vérification : aujourd’hui",
		"last_check_day":   "Dernière vérification : il y a 1 jour",
		"last_check_days":  "Dernière vérification : il y a %d jours",
	},
	"de": {
		"last_check_today": "Letzte Prüfung: heute",
		"last_check_day":   "Letzte Prüfung: vor 1 Tag",
		"last_check_days":  "Letzte Prüfung: vor %d Tagen",
	},
	"ru": {
		"last_check_today": "Последняя проверка: сегодня",
		"last_check_day":   "Последняя проверка: 1 день назад",
		"last_check_days":  "Последняя проверка: %d суток назад",
	},
	"ko": {
		"last_check_today": "마지막 확인: 오늘입니다",
		"last_check_day":   "마지막 확인: 1일 전입니다",
		"last_check_days":  "마지막 확인: %d일 전입니다",
	},
	"id": {
		"last_check_today": "Pemeriksaan terakhir: hari ini",
		"last_check_day":   "Pemeriksaan terakhir: 1 hari yang lalu",
		"last_check_days":  "Pemeriksaan terakhir: %d hari yang lalu",
	},
}

var lastCheckCatalogTargets = map[string]map[string]string{
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
	for lang, catalog := range lastCheckMessages {
		target := lastCheckCatalogTargets[lang]
		for key, value := range catalog {
			target[key] = value
		}
		messages[lang] = target
	}
}
