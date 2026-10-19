package i18n

var spoofMessages = map[string]map[string]string{
	"en": {
		"gateway_spoof":     "Another address on your network (%s) is using the same hardware (MAC) address as your router. A device may be impersonating the router to watch your traffic (ARP spoofing). If your router intentionally uses two addresses, or you connect through a Wi-Fi repeater that answers on the router's behalf, this can be normal.",
		"gateway_spoof_fix": "If you connect through a Wi-Fi repeater or bridge, connect directly to the router and run Mimamori again; if the warning disappears, the repeater is answering on the router's behalf, which is normal for that kind of device. Otherwise, restart the router. If this warning returns, disconnect unfamiliar devices one at a time and run Mimamori again to find the cause, then update or remove that device.",
	},
	"ja": {
		"gateway_spoof":     "ネットワーク上の別のアドレス（%s）が、ルーターと同じハードウェアアドレス（MAC）を使っています。ある機器がルーターになりすまして通信をのぞき見ている（ARP スプーフィング）おそれがあります。ルーターがもともと2つのアドレスを使っている場合や、ルーターの代わりに応答する Wi-Fi 中継器を経由して接続している場合は問題ありません。",
		"gateway_spoof_fix": "Wi-Fi 中継器やブリッジを経由して接続している場合は、ルーターに直接つないで Mimamori を再実行してください。警告が消えれば、中継器がルーターの代わりに応答しているだけで、この種類の機器では正常な動作です。そうでない場合はルーターを再起動してください。それでも警告が出る場合は、見覚えのない機器を1台ずつネットワークから外しながら Mimamori を再実行して原因の機器を探し、その機器を更新するか取り除いてください。",
	},
	"zh-Hans": {
		"gateway_spoof":     "网络上的另一个地址（%s）正在使用与路由器相同的硬件（MAC）地址。某个设备可能正在冒充路由器以窥探您的通信（ARP 欺骗）。如果路由器本来就使用两个地址，或您通过代替路由器应答的 Wi-Fi 中继器连接，这可能是正常现象。",
		"gateway_spoof_fix": "如果您通过 Wi-Fi 中继器或网桥连接，请直接连接到路由器并重新运行 Mimamori；如果提示消失，说明中继器正在代替路由器应答，这种行为对这类设备来说是正常的。否则，请重启路由器。如果此提示再次出现，请逐个断开不熟悉的设备并重新运行 Mimamori，以找出问题设备，然后更新或移除该设备。",
	},
	"zh-Hant": {
		"gateway_spoof":     "網路上的另一個位址（%s）正在使用與路由器相同的硬體（MAC）位址。某個裝置可能正在冒充路由器以窺探您的流量（ARP 欺騙）。如果路由器原本就使用兩個位址，或您透過代替路由器回應的 Wi-Fi 中繼器連線，這可能是正常現象。",
		"gateway_spoof_fix": "如果您透過 Wi-Fi 中繼器或橋接器連線，請直接連線到路由器並重新執行 Mimamori；如果提示消失，表示中繼器正在代替路由器回應，對這類裝置而言這是正常的。否則，請重新啟動路由器。如果此提示再次出現，請逐一中斷不熟悉的裝置並重新執行 Mimamori，以找出問題裝置，然後更新或移除該裝置。",
	},
	"es": {
		"gateway_spoof":     "Otra dirección de su red (%s) está usando la misma dirección de hardware (MAC) que el router. Un dispositivo podría estar suplantando al router para observar su tráfico (suplantación ARP). Si el router usa intencionalmente dos direcciones, o se conecta a través de un repetidor Wi-Fi que responde en nombre del router, esto puede ser normal.",
		"gateway_spoof_fix": "Si se conecta a través de un repetidor Wi-Fi o un puente de red, conéctese directamente al router y vuelva a ejecutar Mimamori; si el aviso desaparece, el repetidor está respondiendo en nombre del router, algo normal en ese tipo de dispositivo. De lo contrario, reinicie el router. Si vuelve a aparecer este aviso, desconecte los dispositivos que no reconozca de uno en uno y vuelva a ejecutar Mimamori para encontrar la causa; luego actualice o retire ese dispositivo.",
	},
	"pt": {
		"gateway_spoof":     "Outro endereço na sua rede (%s) está usando o mesmo endereço de hardware (MAC) do roteador. Um dispositivo pode estar se passando pelo roteador para observar seu tráfego (falsificação ARP). Se o roteador usa intencionalmente dois endereços ou se você se conecta por meio de um repetidor Wi-Fi que responde em nome do roteador, isso pode ser normal.",
		"gateway_spoof_fix": "Se você se conecta por meio de um repetidor Wi-Fi ou de uma ponte, conecte-se diretamente ao roteador e execute o Mimamori novamente; se o aviso desaparecer, o repetidor está respondendo em nome do roteador, o que é normal para esse tipo de dispositivo. Caso contrário, reinicie o roteador. Se o aviso voltar, desconecte os dispositivos desconhecidos um de cada vez e execute o Mimamori novamente para encontrar a causa; depois, atualize ou remova esse dispositivo.",
	},
	"fr": {
		"gateway_spoof":     "Une autre adresse de votre réseau (%s) utilise la même adresse matérielle (MAC) que votre routeur. Un appareil peut se faire passer pour le routeur afin d’observer votre trafic (usurpation ARP). Si votre routeur utilise volontairement deux adresses ou si vous êtes connecté via un répéteur Wi-Fi qui répond au nom du routeur, cela peut être normal.",
		"gateway_spoof_fix": "Si vous êtes connecté via un répéteur Wi-Fi ou un pont, connectez-vous directement au routeur et relancez Mimamori ; si l’avertissement disparaît, le répéteur répond au nom du routeur, ce qui est normal pour ce type d’appareil. Sinon, redémarrez le routeur. Si cet avertissement réapparaît, déconnectez les appareils que vous ne reconnaissez pas un par un et relancez Mimamori pour trouver la cause, puis mettez à jour ou retirez l’appareil concerné.",
	},
	"de": {
		"gateway_spoof":     "Eine andere Adresse in Ihrem Netzwerk (%s) verwendet dieselbe Hardwareadresse (MAC) wie Ihr Router. Ein Gerät könnte sich als Router ausgeben, um Ihren Datenverkehr mitzulesen (ARP-Spoofing). Wenn Ihr Router absichtlich zwei Adressen verwendet oder Sie über einen WLAN-Repeater verbunden sind, der im Namen des Routers antwortet, kann das normal sein.",
		"gateway_spoof_fix": "Wenn Sie über einen WLAN-Repeater oder eine Bridge verbunden sind, verbinden Sie sich direkt mit dem Router und führen Sie Mimamori erneut aus. Wenn der Hinweis dann verschwindet, antwortet der Repeater im Namen des Routers, was bei dieser Art Gerät normal ist. Andernfalls starten Sie den Router neu. Wenn der Hinweis erneut erscheint, trennen Sie unbekannte Geräte einzeln vom Netzwerk und führen Sie Mimamori erneut aus, um die Ursache zu finden. Aktualisieren oder entfernen Sie das betreffende Gerät.",
	},
	"ru": {
		"gateway_spoof":     "Другой адрес в вашей сети (%s) использует тот же аппаратный адрес (MAC), что и роутер. Устройство может выдавать себя за роутер, чтобы следить за вашим трафиком (ARP-спуфинг). Если роутер намеренно использует два адреса или вы подключены через Wi-Fi-репитер, отвечающий от имени роутера, это может быть нормальным явлением.",
		"gateway_spoof_fix": "Если вы подключены через Wi-Fi-репитер или мост, подключитесь напрямую к роутеру и снова запустите Mimamori; если предупреждение исчезнет, значит, репитер отвечает от имени роутера, что нормально для устройств такого типа. В противном случае перезапустите роутер. Если предупреждение появится снова, по очереди отключайте незнакомые устройства от сети и запускайте Mimamori, чтобы найти причину. Затем обновите или удалите это устройство.",
	},
	"ko": {
		"gateway_spoof":     "네트워크의 다른 주소(%s)가 라우터와 같은 하드웨어(MAC) 주소를 사용하고 있습니다. 기기가 라우터를 가장해 통신 내용을 엿볼 수 있습니다(ARP 스푸핑). 라우터가 원래 두 주소를 사용하거나 라우터를 대신해 응답하는 Wi-Fi 중계기를 통해 연결한 경우에는 정상일 수 있습니다.",
		"gateway_spoof_fix": "Wi-Fi 중계기나 브리지를 통해 연결한 경우 라우터에 직접 연결한 뒤 Mimamori를 다시 실행하세요. 경고가 사라지면 중계기가 라우터를 대신해 응답한 것이며, 이런 종류의 기기에서는 정상적인 동작입니다. 그렇지 않으면 라우터를 다시 시작하세요. 경고가 다시 나타나면 낯선 기기를 하나씩 네트워크에서 분리하고 Mimamori를 다시 실행해 원인을 찾은 다음 해당 기기를 업데이트하거나 제거하세요.",
	},
	"id": {
		"gateway_spoof":     "Alamat lain di jaringan Anda (%s) menggunakan alamat perangkat keras (MAC) yang sama dengan router. Perangkat mungkin menyamar sebagai router untuk mengintip lalu lintas Anda (ARP spoofing). Jika router memang menggunakan dua alamat atau Anda terhubung melalui repeater Wi-Fi yang merespons atas nama router, hal ini mungkin normal.",
		"gateway_spoof_fix": "Jika Anda terhubung melalui repeater Wi-Fi atau bridge, sambungkan langsung ke router dan jalankan Mimamori lagi; jika peringatan hilang, repeater sedang merespons atas nama router, hal yang normal untuk perangkat jenis ini. Jika tidak, mulai ulang router. Jika peringatan ini muncul lagi, putuskan perangkat yang tidak dikenal satu per satu dan jalankan Mimamori lagi untuk menemukan penyebabnya, lalu perbarui atau hapus perangkat tersebut.",
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
