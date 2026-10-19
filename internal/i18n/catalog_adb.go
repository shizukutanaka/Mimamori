package i18n

var adbMessages = map[string]map[string]string{
	"en": {
		"adb_open":       "Android debugging (ADB) may be open on %s. Malware has spread between Android TV boxes and phones through this port.",
		"adb_open_fix":   "On the device, turn off USB debugging / wireless debugging in Developer options. If it is a TV box without that setting, update its firmware or consider replacing it.",
		"adb_noauth":     "Android debugging (ADB) on %s lets anyone on your network install apps and control the device without approval.",
		"adb_noauth_fix": "Turn off ADB (network debugging) in the device's Developer options. If you can't find it, contact the maker or stop using the device.",
	},
	"ja": {
		"adb_open":       "「%s」で Android の開発者向け接続（ADB）が開いている可能性があります。このポートを通じて、Android の TV ボックスやスマートフォンに感染を広げるマルウェアが確認されています。",
		"adb_open_fix":   "機器の「開発者向けオプション」で USB デバッグとワイヤレス デバッグをオフにしてください。その設定がない TV ボックスの場合は、ファームウェアを更新するか、買い替えを検討してください。",
		"adb_noauth":     "「%s」の Android のデバッグ機能（ADB）は、家の中の誰でも承認なしにアプリを入れたり機器を操作したりできる状態です。",
		"adb_noauth_fix": "機器の「開発者向けオプション」で ADB（ネットワーク デバッグ）をオフにしてください。見つからない場合は、メーカーに問い合わせるか、その機器の使用をやめることを検討してください。",
	},
	"zh-Hans": {
		"adb_open":       "“%s”上的 Android 调试 (ADB) 可能已开启。恶意软件曾通过此端口在 Android TV 盒子和手机之间传播。",
		"adb_open_fix":   "请在设备的“开发者选项”中关闭 USB 调试/无线调试。如果电视盒子没有此设置，请更新固件或考虑更换设备。",
		"adb_noauth":     "设备“%s”上的 Android 调试功能（ADB）允许您网络中的任何人在未经您批准的情况下安装应用并控制设备。",
		"adb_noauth_fix": "请在设备的“开发者选项”中关闭 ADB（网络调试）。如果找不到该设置，请联系制造商或停止使用该设备。",
	},
	"zh-Hant": {
		"adb_open":       "「%s」上的 Android 偵錯 (ADB) 可能已開啟。惡意軟體曾透過此連接埠在 Android TV 盒子與手機之間傳播。",
		"adb_open_fix":   "請在裝置的「開發者選項」中關閉 USB 偵錯／無線偵錯。如果電視盒子沒有此設定，請更新韌體或考慮更換設備。",
		"adb_noauth":     "「%s」上的 Android 偵錯功能（ADB）讓您網路中的任何人都能在未經您同意的情況下安裝應用程式並控制裝置。",
		"adb_noauth_fix": "請在裝置的「開發者選項」中關閉 ADB（網路偵錯）。如果找不到該設定，請聯絡製造商或停止使用該裝置。",
	},
	"es": {
		"adb_open":       "La depuración de Android (ADB) podría estar abierta en %s. Se ha propagado malware entre TV box Android y teléfonos a través de este puerto.",
		"adb_open_fix":   "Desactive la depuración USB o inalámbrica en Opciones de desarrollador del dispositivo. Si el TV box no tiene esa opción, actualice su firmware o considere reemplazarlo.",
		"adb_noauth":     "La depuración de Android (ADB) en %s permite que cualquiera en su red instale aplicaciones y controle el dispositivo sin su aprobación.",
		"adb_noauth_fix": "Desactive ADB (depuración por red) en las Opciones de desarrollador del dispositivo. Si no encuentra la opción, póngase en contacto con el fabricante o deje de usar el dispositivo.",
	},
	"pt": {
		"adb_open":       "A depuração do Android (ADB) pode estar aberta em %s. Malwares já se espalharam entre TV boxes e celulares Android por essa porta.",
		"adb_open_fix":   "No dispositivo, desative a depuração USB ou sem fio nas Opções do desenvolvedor. Se o TV box não tiver essa opção, atualize o firmware ou considere substituí-lo.",
		"adb_noauth":     "A depuração do Android (ADB) em %s permite que qualquer pessoa na sua rede instale aplicativos e controle o dispositivo sem a sua aprovação.",
		"adb_noauth_fix": "Desative o ADB (depuração pela rede) nas Opções do desenvolvedor do dispositivo. Se não encontrar essa opção, entre em contato com o fabricante ou pare de usar o dispositivo.",
	},
	"fr": {
		"adb_open":       "Le débogage Android (ADB) est peut-être ouvert sur %s. Des logiciels malveillants se sont propagés entre des boîtiers TV Android et des téléphones par ce port.",
		"adb_open_fix":   "Désactivez le débogage USB ou sans fil dans les options pour les développeurs de l’appareil. Si un boîtier TV ne propose pas ce réglage, mettez son micrologiciel à jour ou envisagez de le remplacer.",
		"adb_noauth":     "Le débogage Android (ADB) sur %s permet à toute personne de votre réseau d’installer des applications et de contrôler l’appareil sans votre accord.",
		"adb_noauth_fix": "Désactivez ADB (débogage réseau) dans les options pour les développeurs de l’appareil. Si vous ne trouvez pas ce réglage, contactez le fabricant ou cessez d’utiliser l’appareil.",
	},
	"de": {
		"adb_open":       "Android-Debugging (ADB) ist auf %s möglicherweise geöffnet. Schadsoftware hat sich über diesen Port zwischen Android-TV-Boxen und Smartphones verbreitet.",
		"adb_open_fix":   "Deaktivieren Sie auf dem Gerät USB-Debugging und WLAN-Debugging in den Entwickleroptionen. Falls eine TV-Box diese Einstellung nicht bietet, aktualisieren Sie die Firmware oder erwägen Sie einen Austausch.",
		"adb_noauth":     "Android-Debugging (ADB) auf %s ermöglicht allen in Ihrem Netzwerk, ohne Ihre Zustimmung Apps zu installieren und das Gerät zu steuern.",
		"adb_noauth_fix": "Deaktivieren Sie ADB (Netzwerk-Debugging) in den Entwickleroptionen des Geräts. Wenn Sie die Einstellung nicht finden, wenden Sie sich an den Hersteller oder verwenden Sie das Gerät nicht mehr.",
	},
	"ru": {
		"adb_open":       "На устройстве «%s» может быть открыт Android Debug Bridge (ADB). Через этот порт вредоносное ПО распространялось между Android TV-приставками и телефонами.",
		"adb_open_fix":   "Отключите отладку по USB и по беспроводной сети в разделе «Для разработчиков» на устройстве. Если в TV-приставке нет такой настройки, обновите её прошивку или подумайте о замене устройства.",
		"adb_noauth":     "Отладка Android (ADB) на устройстве «%s» позволяет любому в вашей сети устанавливать приложения и управлять устройством без вашего разрешения.",
		"adb_noauth_fix": "Отключите ADB (сетевую отладку) в разделе «Для разработчиков» на устройстве. Если вы не можете найти эту настройку, свяжитесь с производителем или прекратите пользоваться устройством.",
	},
	"ko": {
		"adb_open":       "기기 %s에서 Android 디버깅(ADB)이 열려 있을 수 있습니다. 이 포트를 통해 Android TV 박스와 휴대전화 사이로 악성코드가 퍼진 사례가 있습니다.",
		"adb_open_fix":   "기기의 개발자 옵션에서 USB 디버깅/무선 디버깅을 끄세요. TV 박스에 해당 설정이 없다면 펌웨어를 업데이트하거나 기기 교체를 고려하세요.",
		"adb_noauth":     "%s의 Android 디버깅(ADB)은 네트워크의 누구나 승인 없이 앱을 설치하고 기기를 제어할 수 있게 합니다.",
		"adb_noauth_fix": "기기의 개발자 옵션에서 ADB(네트워크 디버깅)를 끄세요. 설정을 찾을 수 없다면 제조사에 문의하거나 해당 기기의 사용을 중단하세요.",
	},
	"id": {
		"adb_open":       "Android debugging (ADB) mungkin terbuka di %s. Malware pernah menyebar antara Android TV box dan ponsel melalui port ini.",
		"adb_open_fix":   "Di perangkat, matikan USB debugging atau wireless debugging di Opsi developer. Jika TV box tidak memiliki pengaturan itu, perbarui firmware atau pertimbangkan untuk menggantinya.",
		"adb_noauth":     "Android debugging (ADB) pada %s memungkinkan siapa pun di jaringan Anda memasang aplikasi dan mengendalikan perangkat tanpa persetujuan Anda.",
		"adb_noauth_fix": "Matikan ADB (debugging jaringan) di Opsi developer perangkat. Jika Anda tidak menemukan pengaturan itu, hubungi produsen atau berhenti menggunakan perangkat tersebut.",
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
