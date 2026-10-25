package i18n

var privateAddressMessages = map[string]map[string]string{
	"en": {
		"new_device_private":     "A device using a private (changing) address connected: %s. Phones and PCs change this address over time, so it may be one of your own devices.",
		"new_device_private_fix": "Check whether everyone's phones/PCs are accounted for. If you see devices you don't recognize repeatedly, change the Wi-Fi password and reconnect only your own devices.",
	},
	"ja": {
		"new_device_private":     "プライベートアドレス（変化するアドレス）を使う機器が接続しました: %s。スマホや PC はこのアドレスを定期的に変えるため、ご家族の機器の可能性があります。",
		"new_device_private_fix": "家族のスマホや PC か確認してください。心当たりのない機器が繰り返し出る場合は、Wi-Fi のパスワードを変更し、自分の機器だけ再接続してください。",
	},
	"zh-Hans": {
		"new_device_private":     "使用私有（会变化）地址的设备已连接：%s。手机和电脑会随时间更换此地址，因此可能是您自己的设备。",
		"new_device_private_fix": "请确认家人的手机和电脑都已记录。如果反复出现您不认识的设备，请更改 Wi-Fi 密码，并只重新连接自己的设备。",
	},
	"zh-Hant": {
		"new_device_private":     "使用私人（會變動）位址的裝置已連線：%s。手機和電腦會隨時間更換此位址，因此可能是您自己的裝置。",
		"new_device_private_fix": "請確認家人的手機和電腦都已列入。如果反覆出現您不認識的裝置，請更改 Wi-Fi 密碼，並只重新連線自己的裝置。",
	},
	"es": {
		"new_device_private":     "Se conectó un dispositivo que usa una dirección privada (cambiante): %s. Los teléfonos y los PC cambian esta dirección con el tiempo, por lo que podría ser uno de sus dispositivos.",
		"new_device_private_fix": "Compruebe que los teléfonos y los PC de todos estén contemplados. Si siguen apareciendo dispositivos que no reconoce, cambie la contraseña del Wi-Fi y vuelva a conectar únicamente sus propios dispositivos.",
	},
	"pt": {
		"new_device_private":     "Um dispositivo com endereço privado (variável) se conectou: %s. Celulares e computadores mudam esse endereço com o tempo, então pode ser um dos seus aparelhos.",
		"new_device_private_fix": "Confira se os celulares e computadores da família estão na lista. Se aparelhos desconhecidos continuarem aparecendo, troque a senha do Wi-Fi e reconecte apenas os seus aparelhos.",
	},
	"fr": {
		"new_device_private":     "Un appareil utilisant une adresse privée (changeante) s’est connecté : %s. Les téléphones et les PC changent cette adresse au fil du temps ; il peut donc s’agir d’un de vos appareils.",
		"new_device_private_fix": "Vérifiez que les téléphones et les PC de chacun sont bien répertoriés. Si des appareils inconnus réapparaissent, changez le mot de passe Wi-Fi et reconnectez uniquement les vôtres.",
	},
	"de": {
		"new_device_private":     "Ein Gerät mit privater (wechselnder) Adresse hat sich verbunden: %s. Smartphones und PCs ändern diese Adresse im Laufe der Zeit, daher könnte es eines Ihrer Geräte sein.",
		"new_device_private_fix": "Prüfen Sie, ob die Smartphones und PCs aller Familienmitglieder erfasst sind. Wenn wiederholt unbekannte Geräte auftauchen, ändern Sie das WLAN-Passwort und verbinden Sie nur Ihre eigenen Geräte erneut.",
	},
	"ru": {
		"new_device_private":     "Подключилось устройство с частным (меняющимся) адресом: %s. Телефоны и компьютеры со временем меняют этот адрес, поэтому это может быть одно из ваших устройств.",
		"new_device_private_fix": "Проверьте, учтены ли телефоны и компьютеры всех членов семьи. Если незнакомые устройства появляются снова, смените пароль Wi-Fi и подключите заново только свои устройства.",
	},
	"ko": {
		"new_device_private":     "비공개(변경되는) 주소를 사용하는 기기가 연결되었습니다: %s. 휴대폰과 PC는 시간이 지나면 이 주소를 바꾸므로 가족의 기기일 수 있습니다.",
		"new_device_private_fix": "가족의 휴대폰과 PC가 모두 확인되었는지 살펴보세요. 모르는 기기가 반복해서 나타나면 Wi-Fi 비밀번호를 바꾸고 본인 기기만 다시 연결하세요.",
	},
	"id": {
		"new_device_private":     "Perangkat dengan alamat privat (berubah-ubah) terhubung: %s. Ponsel dan PC mengubah alamat ini dari waktu ke waktu, jadi perangkat itu mungkin milik Anda.",
		"new_device_private_fix": "Pastikan ponsel dan PC semua anggota keluarga sudah tercatat. Jika perangkat yang tidak dikenal terus muncul, ganti kata sandi Wi-Fi dan sambungkan kembali hanya perangkat milik Anda.",
	},
}

func init() {
	add(privateAddressMessages)
}
