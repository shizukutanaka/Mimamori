package i18n

var historyMessages = map[string]map[string]string{
	"en": {
		"history_unavailable": "The device history file could not be read, so new devices cannot be reported and this check was not saved. Make sure the history file can be read and written.",
	},
	"ja": {
		"history_unavailable": "機器の履歴ファイルを読み込めなかったため、新しい機器をお知らせできません。今回の確認も保存していません。履歴ファイルを読み書きできるか確かめてください。",
	},
	"zh-Hans": {
		"history_unavailable": "无法读取设备历史文件，因此无法提示新设备，本次检查也未保存。请确认历史文件可以读写。",
	},
	"zh-Hant": {
		"history_unavailable": "無法讀取裝置歷史檔案，因此無法提示新裝置，本次檢查也未儲存。請確認歷史檔案可以讀寫。",
	},
	"es": {
		"history_unavailable": "No se pudo leer el historial de dispositivos, así que no se pueden avisar los dispositivos nuevos y esta comprobación no se guardó. Compruebe que el archivo de historial se pueda leer y escribir.",
	},
	"pt": {
		"history_unavailable": "Não foi possível ler o histórico de dispositivos, então dispositivos novos não podem ser avisados e esta verificação não foi salva. Verifique se o arquivo de histórico pode ser lido e gravado.",
	},
	"fr": {
		"history_unavailable": "L’historique des appareils n’a pas pu être lu : les nouveaux appareils ne peuvent donc pas être signalés et cette vérification n’a pas été enregistrée. Vérifiez que le fichier d’historique est accessible en lecture et en écriture.",
	},
	"de": {
		"history_unavailable": "Der Geräteverlauf konnte nicht gelesen werden. Neue Geräte können daher nicht gemeldet werden, und diese Prüfung wurde nicht gespeichert. Prüfen Sie, ob die Verlaufsdatei gelesen und geschrieben werden kann.",
	},
	"ru": {
		"history_unavailable": "Не удалось прочитать историю устройств, поэтому о новых устройствах сообщить нельзя, а эта проверка не сохранена. Убедитесь, что файл истории доступен для чтения и записи.",
	},
	"ko": {
		"history_unavailable": "기기 기록 파일을 읽을 수 없어 새 기기를 알려드릴 수 없으며, 이번 확인도 저장하지 않았습니다. 기록 파일을 읽고 쓸 수 있는지 확인해 주세요.",
	},
	"id": {
		"history_unavailable": "File riwayat perangkat tidak dapat dibaca, sehingga perangkat baru tidak dapat dilaporkan dan pemeriksaan ini tidak disimpan. Pastikan file riwayat dapat dibaca dan ditulis.",
	},
}

func init() {
	add(historyMessages)
}
