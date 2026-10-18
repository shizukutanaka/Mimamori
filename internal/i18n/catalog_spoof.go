package i18n

var spoofMessages = map[string]map[string]string{
	"en": {
		"gateway_spoof":     "Another address on your network (%s) is using the same hardware (MAC) address as your router. A device may be impersonating the router to watch your traffic (ARP spoofing). If your router intentionally uses two addresses, you can ignore this.",
		"gateway_spoof_fix": "Restart the router. If this warning returns, disconnect unfamiliar devices one at a time and run Mimamori again to find the cause, then update or remove that device.",
	},
	"ja": {
		"gateway_spoof":     "ネットワーク上の別のアドレス（%s）が、ルーターと同じハードウェアアドレス（MAC）を使っています。ある機器がルーターになりすまして通信をのぞき見ている（ARP スプーフィング）おそれがあります。ルーターがもともと2つのアドレスを使っている場合は問題ありません。",
		"gateway_spoof_fix": "ルーターを再起動してください。それでも警告が出る場合は、見覚えのない機器を1台ずつネットワークから外しながら Mimamori を再実行して原因の機器を探し、その機器を更新するか取り除いてください。",
	},
	"zh-Hans": {
		"gateway_spoof":     "网络上的另一个地址（%s）正在使用与路由器相同的硬件（MAC）地址。某个设备可能正在冒充路由器以窥探您的通信（ARP 欺骗）。如果路由器本来就使用两个地址，可以忽略此提示。",
		"gateway_spoof_fix": "请重启路由器。如果此提示再次出现，请逐个断开不熟悉的设备并重新运行 Mimamori，以找出问题设备，然后更新或移除该设备。",
	},
	"zh-Hant": {
		"gateway_spoof":     "網路上的另一個位址（%s）正在使用與路由器相同的硬體（MAC）位址。某個裝置可能正在冒充路由器以窺探您的流量（ARP 欺騙）。如果路由器原本就使用兩個位址，可以忽略此提示。",
		"gateway_spoof_fix": "請重新啟動路由器。如果此提示再次出現，請逐一中斷不熟悉的裝置並重新執行 Mimamori，以找出問題裝置，然後更新或移除該裝置。",
	},
	"es": {
		"gateway_spoof":     "Otra dirección de tu red (%s) está usando la misma dirección de hardware (MAC) que el router. Un dispositivo podría estar suplantando al router para observar tu tráfico (suplantación ARP). Si el router usa intencionalmente dos direcciones, puedes ignorar este aviso.",
		"gateway_spoof_fix": "Reinicia el router. Si vuelve a aparecer este aviso, desconecta los dispositivos que no reconozcas de uno en uno y vuelve a ejecutar Mimamori para encontrar la causa; luego actualiza o retira ese dispositivo.",
	},
	"pt": {
		"gateway_spoof":     "Outro endereço na sua rede (%s) está usando o mesmo endereço de hardware (MAC) do roteador. Um dispositivo pode estar se passando pelo roteador para observar seu tráfego (falsificação ARP). Se o roteador usa intencionalmente dois endereços, você pode ignorar este aviso.",
		"gateway_spoof_fix": "Reinicie o roteador. Se o aviso voltar, desconecte os dispositivos desconhecidos um de cada vez e execute o Mimamori novamente para encontrar a causa; depois, atualize ou remova esse dispositivo.",
	},
	"fr": {
		"gateway_spoof":     "Une autre adresse de votre réseau (%s) utilise la même adresse matérielle (MAC) que votre routeur. Un appareil peut se faire passer pour le routeur afin d’observer votre trafic (usurpation ARP). Si votre routeur utilise volontairement deux adresses, vous pouvez ignorer cet avertissement.",
		"gateway_spoof_fix": "Redémarrez le routeur. Si cet avertissement réapparaît, déconnectez les appareils que vous ne reconnaissez pas un par un et relancez Mimamori pour trouver la cause, puis mettez à jour ou retirez l’appareil concerné.",
	},
	"de": {
		"gateway_spoof":     "Eine andere Adresse in Ihrem Netzwerk (%s) verwendet dieselbe Hardwareadresse (MAC) wie Ihr Router. Ein Gerät könnte sich als Router ausgeben, um Ihren Datenverkehr mitzulesen (ARP-Spoofing). Wenn Ihr Router absichtlich zwei Adressen verwendet, können Sie diesen Hinweis ignorieren.",
		"gateway_spoof_fix": "Starten Sie den Router neu. Wenn der Hinweis erneut erscheint, trennen Sie unbekannte Geräte einzeln vom Netzwerk und führen Sie Mimamori erneut aus, um die Ursache zu finden. Aktualisieren oder entfernen Sie das betreffende Gerät.",
	},
	"ru": {
		"gateway_spoof":     "Другой адрес в вашей сети (%s) использует тот же аппаратный адрес (MAC), что и роутер. Устройство может выдавать себя за роутер, чтобы следить за вашим трафиком (ARP-спуфинг). Если роутер намеренно использует два адреса, это предупреждение можно игнорировать.",
		"gateway_spoof_fix": "Перезапустите роутер. Если предупреждение появится снова, по очереди отключайте незнакомые устройства от сети и запускайте Mimamori, чтобы найти причину. Затем обновите или удалите это устройство.",
	},
	"ko": {
		"gateway_spoof":     "네트워크의 다른 주소(%s)가 라우터와 같은 하드웨어(MAC) 주소를 사용하고 있습니다. 기기가 라우터를 가장해 통신 내용을 엿볼 수 있습니다(ARP 스푸핑). 라우터가 원래 두 주소를 사용하는 경우 이 경고를 무시해도 됩니다.",
		"gateway_spoof_fix": "라우터를 다시 시작하세요. 경고가 다시 나타나면 낯선 기기를 하나씩 네트워크에서 분리하고 Mimamori를 다시 실행해 원인을 찾은 다음 해당 기기를 업데이트하거나 제거하세요.",
	},
	"id": {
		"gateway_spoof":     "Alamat lain di jaringan Anda (%s) menggunakan alamat perangkat keras (MAC) yang sama dengan router. Perangkat mungkin menyamar sebagai router untuk mengintip lalu lintas Anda (ARP spoofing). Jika router memang menggunakan dua alamat, Anda dapat mengabaikan peringatan ini.",
		"gateway_spoof_fix": "Mulai ulang router. Jika peringatan ini muncul lagi, putuskan perangkat yang tidak dikenal satu per satu dan jalankan Mimamori lagi untuk menemukan penyebabnya, lalu perbarui atau hapus perangkat tersebut.",
	},
}

func init() {
	for lang, catalog := range spoofMessages {
		target := wifiCatalogTargets[lang]
		for key, value := range catalog {
			target[key] = value
		}
		messages[lang] = target
	}
}
