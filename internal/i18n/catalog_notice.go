package i18n

var noticeMessages = map[string]map[string]string{
	"en": {
		"incomplete_notice": "Some checks could not run, so this result may be incomplete. Running it again on the home network may help.",
		"no_devices_notice": "The scan could not see any device on your network, so the result would not be reliable. Make sure this PC is connected to your home network and run it again.",
		"scanning_notice":   "Checking your home network. This takes a little while…",
	},
	"ja": {
		"incomplete_notice": "一部の確認ができなかったため、この結果は不完全な可能性があります。自宅のネットワークに接続した状態で、もう一度実行してください。",
		"no_devices_notice": "家のネットワークに機器が 1 台も見つからなかったため、結果を信頼できません。この PC が自宅のネットワークに接続していることを確認して、もう一度実行してください。",
		"scanning_notice":   "家のネットワークを確認しています。少し時間がかかります…",
	},
	"zh-Hans": {
		"incomplete_notice": "部分检查未能运行，因此结果可能不完整。请在连接到家庭网络后再运行一次。",
		"no_devices_notice": "在您的网络上没有发现任何设备，因此结果不可靠。请确认这台电脑已连接到家庭网络，然后再运行一次。",
		"scanning_notice":   "正在检查您的家庭网络。这需要一些时间…",
	},
	"zh-Hant": {
		"incomplete_notice": "部分檢查無法執行，因此結果可能不完整。請在連上家中網路後再執行一次。",
		"no_devices_notice": "在您的網路上沒有發現任何裝置，因此結果不可靠。請確認這台電腦已連線到家庭網路，然後再執行一次。",
		"scanning_notice":   "正在檢查您的家庭網路。這需要一些時間…",
	},
	"es": {
		"incomplete_notice": "No se pudieron realizar algunas comprobaciones, así que el resultado puede estar incompleto. Vuelva a ejecutarlo conectado a la red de casa.",
		"no_devices_notice": "El análisis no encontró ningún dispositivo en su red, por lo que el resultado no sería fiable. Compruebe que este PC está conectado a la red de su casa y ejecútelo de nuevo.",
		"scanning_notice":   "Estamos comprobando su red doméstica. Esto tardará un poco…",
	},
	"pt": {
		"incomplete_notice": "Não foi possível executar algumas verificações, então o resultado pode estar incompleto. Tente novamente conectado à rede de casa.",
		"no_devices_notice": "A verificação não encontrou nenhum dispositivo na sua rede, por isso o resultado não seria confiável. Confirme que este PC está conectado à rede de casa e execute novamente.",
		"scanning_notice":   "Verificando sua rede doméstica. Isso vai levar um pouco…",
	},
	"fr": {
		"incomplete_notice": "Certaines vérifications n’ont pas pu être effectuées ; le résultat peut donc être incomplet. Relancez l’outil lorsque vous êtes connecté au réseau de la maison.",
		"no_devices_notice": "L’analyse n’a trouvé aucun appareil sur votre réseau ; le résultat ne serait donc pas fiable. Vérifiez que ce PC est connecté au réseau de la maison, puis relancez-la.",
		"scanning_notice":   "Vérification de votre réseau domestique. Cela prend un peu de temps…",
	},
	"de": {
		"incomplete_notice": "Einige Prüfungen konnten nicht ausgeführt werden. Das Ergebnis ist daher möglicherweise unvollständig. Führen Sie den Check erneut im Heimnetzwerk aus.",
		"no_devices_notice": "Die Prüfung hat kein Gerät in Ihrem Netzwerk gefunden, daher wäre das Ergebnis nicht verlässlich. Stellen Sie sicher, dass dieser PC mit Ihrem Heimnetzwerk verbunden ist, und führen Sie sie erneut aus.",
		"scanning_notice":   "Ihr Heimnetzwerk wird überprüft. Das dauert einen Moment…",
	},
	"ru": {
		"incomplete_notice": "Некоторые проверки не удалось выполнить, поэтому результат может быть неполным. Запустите проверку ещё раз, подключившись к домашней сети.",
		"no_devices_notice": "Проверка не обнаружила ни одного устройства в вашей сети, поэтому результат был бы ненадёжным. Убедитесь, что этот ПК подключён к домашней сети, и запустите проверку ещё раз.",
		"scanning_notice":   "Проверяем домашнюю сеть. Это займёт немного времени…",
	},
	"ko": {
		"incomplete_notice": "일부 검사를 실행하지 못해 결과가 불완전할 수 있습니다. 집 네트워크에 연결한 상태에서 다시 실행해 보세요.",
		"no_devices_notice": "네트워크에서 기기를 하나도 찾지 못해 결과를 신뢰할 수 없습니다. 이 PC가 집 네트워크에 연결되어 있는지 확인하고 다시 실행하세요.",
		"scanning_notice":   "집 네트워크를 확인하고 있습니다. 잠시 시간이 걸립니다…",
	},
	"id": {
		"incomplete_notice": "Beberapa pemeriksaan tidak dapat dijalankan, jadi hasil ini mungkin belum lengkap. Coba jalankan kembali saat terhubung ke jaringan rumah.",
		"no_devices_notice": "Pemeriksaan tidak menemukan perangkat apa pun di jaringan Anda, sehingga hasilnya tidak dapat diandalkan. Pastikan PC ini terhubung ke jaringan rumah, lalu jalankan lagi.",
		"scanning_notice":   "Sedang memeriksa jaringan rumah Anda. Proses ini memerlukan sedikit waktu…",
	},
}

func init() {
	add(noticeMessages)
}
