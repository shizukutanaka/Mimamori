package i18n

var adbMessages = map[string]map[string]string{
	"en": {
		"adb_open":     "Android debugging (ADB) may be open on %s. Malware has spread between Android TV boxes and phones through this port.",
		"adb_open_fix": "On the device, turn off USB debugging / wireless debugging in Developer options. If it is a TV box without that setting, update its firmware or consider replacing it.",
	},
	"ja": {
		"adb_open":     "「%s」で Android の開発者向け接続（ADB）が開いている可能性があります。このポートを通じて、Android の TV ボックスやスマートフォンに感染を広げるマルウェアが確認されています。",
		"adb_open_fix": "機器の「開発者向けオプション」で USB デバッグとワイヤレス デバッグをオフにしてください。その設定がない TV ボックスの場合は、ファームウェアを更新するか、買い替えを検討してください。",
	},
	"zh-Hans": {
		"adb_open":     "“%s”上的 Android 调试 (ADB) 可能已开启。恶意软件曾通过此端口在 Android TV 盒子和手机之间传播。",
		"adb_open_fix": "请在设备的“开发者选项”中关闭 USB 调试/无线调试。如果电视盒子没有此设置，请更新固件或考虑更换设备。",
	},
	"zh-Hant": {
		"adb_open":     "「%s」上的 Android 偵錯 (ADB) 可能已開啟。惡意軟體曾透過此連接埠在 Android TV 盒子與手機之間傳播。",
		"adb_open_fix": "請在裝置的「開發者選項」中關閉 USB 偵錯／無線偵錯。如果電視盒子沒有此設定，請更新韌體或考慮更換設備。",
	},
	"es": {
		"adb_open":     "La depuración de Android (ADB) podría estar abierta en %s. Se ha propagado malware entre TV box Android y teléfonos a través de este puerto.",
		"adb_open_fix": "Desactive la depuración USB o inalámbrica en Opciones de desarrollador del dispositivo. Si el TV box no tiene esa opción, actualice su firmware o considere reemplazarlo.",
	},
	"pt": {
		"adb_open":     "A depuração do Android (ADB) pode estar aberta em %s. Malwares já se espalharam entre TV boxes e celulares Android por essa porta.",
		"adb_open_fix": "No dispositivo, desative a depuração USB ou sem fio nas Opções do desenvolvedor. Se o TV box não tiver essa opção, atualize o firmware ou considere substituí-lo.",
	},
	"fr": {
		"adb_open":     "Le débogage Android (ADB) est peut-être ouvert sur %s. Des logiciels malveillants se sont propagés entre des boîtiers TV Android et des téléphones par ce port.",
		"adb_open_fix": "Désactivez le débogage USB ou sans fil dans les options pour les développeurs de l’appareil. Si un boîtier TV ne propose pas ce réglage, mettez son micrologiciel à jour ou envisagez de le remplacer.",
	},
	"de": {
		"adb_open":     "Android-Debugging (ADB) ist auf %s möglicherweise geöffnet. Schadsoftware hat sich über diesen Port zwischen Android-TV-Boxen und Smartphones verbreitet.",
		"adb_open_fix": "Deaktivieren Sie auf dem Gerät USB-Debugging und WLAN-Debugging in den Entwickleroptionen. Falls eine TV-Box diese Einstellung nicht bietet, aktualisieren Sie die Firmware oder erwägen Sie einen Austausch.",
	},
	"ru": {
		"adb_open":     "На устройстве «%s» может быть открыт Android Debug Bridge (ADB). Через этот порт вредоносное ПО распространялось между Android TV-приставками и телефонами.",
		"adb_open_fix": "Отключите отладку по USB и по беспроводной сети в разделе «Для разработчиков» на устройстве. Если в TV-приставке нет такой настройки, обновите её прошивку или подумайте о замене устройства.",
	},
	"ko": {
		"adb_open":     "기기 %s에서 Android 디버깅(ADB)이 열려 있을 수 있습니다. 이 포트를 통해 Android TV 박스와 휴대전화 사이로 악성코드가 퍼진 사례가 있습니다.",
		"adb_open_fix": "기기의 개발자 옵션에서 USB 디버깅/무선 디버깅을 끄세요. TV 박스에 해당 설정이 없다면 펌웨어를 업데이트하거나 기기 교체를 고려하세요.",
	},
	"id": {
		"adb_open":     "Android debugging (ADB) mungkin terbuka di %s. Malware pernah menyebar antara Android TV box dan ponsel melalui port ini.",
		"adb_open_fix": "Di perangkat, matikan USB debugging atau wireless debugging di Opsi developer. Jika TV box tidak memiliki pengaturan itu, perbarui firmware atau pertimbangkan untuk menggantinya.",
	},
}

func init() {
	for lang, catalog := range adbMessages {
		target := wifiCatalogTargets[lang]
		for key, value := range catalog {
			target[key] = value
		}
		messages[lang] = target
	}
}
