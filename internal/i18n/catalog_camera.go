package i18n

var cameraMessages = map[string]map[string]string{
	"en": {
		"forward_camera":     "The camera or recorder %s exposes its video/control port %d to the whole Internet.",
		"forward_camera_fix": "Disable this port forward or UPnP on the router, update the camera firmware, and change its password. For remote viewing, use the vendor's cloud service/app or a VPN.",
	},
	"ja": {
		"forward_camera":     "カメラ・録画機「%s」の映像・操作用ポート %d がインターネット全体に公開されています。",
		"forward_camera_fix": "ルーターでこのポート転送または UPnP を無効にし、カメラのファームウェアを更新してパスワードを変更してください。遠隔視聴にはメーカーのクラウドサービス・アプリまたは VPN を使ってください。",
	},
	"zh-Hans": {
		"forward_camera":     "摄像头或录像机 %s 的视频/控制端口 %d 可从整个互联网访问。",
		"forward_camera_fix": "请在路由器上关闭此端口转发或 UPnP，更新摄像头固件并更改密码。需要远程观看时，请使用厂商的云服务/应用或 VPN。",
	},
	"zh-Hant": {
		"forward_camera":     "攝影機或錄影機 %s 的影像／控制連接埠 %d 可從整個網際網路存取。",
		"forward_camera_fix": "請在路由器上關閉此連接埠轉送或 UPnP，更新攝影機韌體並更改密碼。需要遠端觀看時，請使用廠商的雲端服務／應用程式或 VPN。",
	},
	"es": {
		"forward_camera":     "La cámara o el grabador %s expone su puerto de vídeo/control %d a todo Internet.",
		"forward_camera_fix": "Desactive esta redirección de puertos o UPnP en el router, actualice el firmware de la cámara y cambie su contraseña. Para ver las imágenes a distancia, use el servicio en la nube o la aplicación del fabricante, o una VPN.",
	},
	"pt": {
		"forward_camera":     "A câmera ou o gravador %s expõe sua porta de vídeo/controle %d a toda a Internet.",
		"forward_camera_fix": "Desative este encaminhamento de porta ou o UPnP no roteador, atualize o firmware da câmera e troque a senha. Para visualizar remotamente, use a nuvem ou o aplicativo do fabricante, ou uma VPN.",
	},
	"fr": {
		"forward_camera":     "La caméra ou l’enregistreur %s expose son port vidéo/de commande %d à l’ensemble d’Internet.",
		"forward_camera_fix": "Désactivez cette redirection de port ou UPnP sur le routeur, mettez à jour le micrologiciel de la caméra et changez son mot de passe. Pour regarder les images à distance, utilisez le cloud ou l’application du fabricant, ou un VPN.",
	},
	"de": {
		"forward_camera":     "Die Kamera bzw. der Rekorder %s ist über den Video-/Steuerungsport %d aus dem gesamten Internet erreichbar.",
		"forward_camera_fix": "Deaktivieren Sie diese Portweiterleitung oder UPnP im Router, aktualisieren Sie die Kamera-Firmware und ändern Sie ihr Passwort. Nutzen Sie für die Fernansicht die Cloud oder App des Herstellers oder ein VPN.",
	},
	"ru": {
		"forward_camera":     "Камера или видеорегистратор %s открывает видеопорт или порт управления %d для всего Интернета.",
		"forward_camera_fix": "Отключите это перенаправление порта или UPnP в настройках роутера, обновите прошивку камеры и смените пароль. Для удалённого просмотра используйте облачный сервис или приложение производителя либо VPN.",
	},
	"ko": {
		"forward_camera":     "카메라 또는 녹화기 %s의 영상/제어 포트 %d가 인터넷 전체에 노출되어 있습니다.",
		"forward_camera_fix": "라우터에서 이 포트 전달 또는 UPnP를 끄고 카메라 펌웨어를 업데이트한 뒤 비밀번호를 변경하세요. 원격으로 볼 때는 제조사의 클라우드/앱이나 VPN을 사용하세요.",
	},
	"id": {
		"forward_camera":     "Kamera atau perekam %s memiliki port video/kontrol %d yang dapat diakses dari seluruh Internet.",
		"forward_camera_fix": "Nonaktifkan penerusan port ini atau UPnP di router, perbarui firmware kamera, dan ganti kata sandinya. Untuk melihat dari jarak jauh, gunakan layanan cloud/aplikasi vendor atau VPN.",
	},
}

func init() {
	for lang, catalog := range cameraMessages {
		target := wifiCatalogTargets[lang]
		for key, value := range catalog {
			target[key] = value
		}
		messages[lang] = target
	}
}
