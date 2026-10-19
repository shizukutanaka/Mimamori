package i18n

var smb1Messages = map[string]map[string]string{
	"en": {
		"smb1_enabled":     "%s accepts the old SMBv1 file-sharing protocol. Ransomware such as WannaCry spread between home PCs and NAS devices through SMBv1.",
		"smb1_enabled_fix": "Turn off SMBv1. On Windows, open 'Turn Windows features on or off' and clear 'SMB 1.0/CIFS File Sharing Support'. On a NAS or router, set the minimum SMB version to SMB2 or later in its settings, or update its firmware.",
	},
	"ja": {
		"smb1_enabled":     "「%s」が古いファイル共有の方式（SMBv1）を受け付けています。WannaCry などのランサムウェアは、SMBv1 を通じて家の中の PC や NAS に広がりました。",
		"smb1_enabled_fix": "SMBv1 をオフにしてください。Windows では「Windows の機能の有効化または無効化」で「SMB 1.0/CIFS ファイル共有のサポート」のチェックを外します。NAS やルーターでは、設定で SMB の最小バージョンを SMB2 以上にするか、ファームウェアを更新してください。",
	},
	"zh-Hans": {
		"smb1_enabled":     "“%s”仍接受旧版文件共享协议 SMBv1。WannaCry 等勒索软件曾通过 SMBv1 在家庭电脑和 NAS 设备之间传播。",
		"smb1_enabled_fix": "请您关闭 SMBv1。在 Windows 中，打开“启用或关闭 Windows 功能”，取消选中“SMB 1.0/CIFS 文件共享支持”。在 NAS 或路由器上，请您在设置中将最低 SMB 版本设为 SMB2 或更高版本，或更新固件。",
	},
	"zh-Hant": {
		"smb1_enabled":     "「%s」仍接受舊版檔案分享通訊協定 SMBv1。WannaCry 等勒索軟體曾透過 SMBv1 在家用電腦與 NAS 裝置之間傳播。",
		"smb1_enabled_fix": "請您關閉 SMBv1。在 Windows 中，開啟「開啟或關閉 Windows 功能」，取消勾選「SMB 1.0/CIFS 檔案共用支援」。在 NAS 或路由器上，請您於設定中將最低 SMB 版本設為 SMB2 或更新版本，或更新韌體。",
	},
	"es": {
		"smb1_enabled":     "%s acepta el antiguo protocolo de uso compartido de archivos SMBv1. El ransomware, como WannaCry, se propagó entre PC domésticos y dispositivos NAS mediante SMBv1.",
		"smb1_enabled_fix": "Desactive SMBv1. En Windows, abra «Activar o desactivar las características de Windows» y desmarque «Compatibilidad con el protocolo para compartir archivos SMB 1.0/CIFS». En un NAS o router, establezca la versión mínima de SMB en SMB2 o posterior en la configuración, o actualice su firmware.",
	},
	"pt": {
		"smb1_enabled":     "%s aceita o antigo protocolo de compartilhamento de arquivos SMBv1. Ransomwares como o WannaCry se espalharam entre PCs domésticos e dispositivos NAS por meio do SMBv1.",
		"smb1_enabled_fix": "Você deve desativar o SMBv1. No Windows, abra “Ativar ou desativar recursos do Windows” e desmarque “Suporte para Compartilhamento de Arquivos SMB 1.0/CIFS”. Em um NAS ou roteador, defina a versão mínima do SMB como SMB2 ou posterior nas configurações, ou atualize o firmware.",
	},
	"fr": {
		"smb1_enabled":     "« %s » accepte l’ancien protocole de partage de fichiers SMBv1. Des rançongiciels comme WannaCry se sont propagés entre des PC personnels et des appareils NAS via SMBv1.",
		"smb1_enabled_fix": "Désactivez SMBv1. Sous Windows, ouvrez « Activer ou désactiver des fonctionnalités Windows » et décochez « Prise en charge du partage de fichiers SMB 1.0/CIFS ». Sur un NAS ou un routeur, définissez la version minimale de SMB sur SMB2 ou une version ultérieure dans les paramètres, ou mettez son micrologiciel à jour.",
	},
	"de": {
		"smb1_enabled":     "„%s“ akzeptiert das alte Dateifreigabeprotokoll SMBv1. Ransomware wie WannaCry hat sich über SMBv1 zwischen PCs und NAS-Geräten in Privathaushalten verbreitet.",
		"smb1_enabled_fix": "Deaktivieren Sie SMBv1. Öffnen Sie unter Windows „Windows-Features aktivieren oder deaktivieren“ und entfernen Sie das Häkchen bei „Unterstützung für die SMB 1.0/CIFS-Dateifreigabe“. Legen Sie auf einem NAS oder Router in den Einstellungen die Mindestversion für SMB auf SMB2 oder höher fest oder aktualisieren Sie die Firmware.",
	},
	"ru": {
		"smb1_enabled":     "Устройство «%s» принимает старый протокол общего доступа к файлам SMBv1. Программы-вымогатели, такие как WannaCry, распространялись через SMBv1 между домашними ПК и NAS-устройствами.",
		"smb1_enabled_fix": "Отключите SMBv1. В Windows откройте «Включение или отключение компонентов Windows» и снимите флажок «Поддержка общего доступа к файлам SMB 1.0/CIFS». На NAS или маршрутизаторе установите в настройках минимальную версию SMB — SMB2 или выше — либо обновите прошивку.",
	},
	"ko": {
		"smb1_enabled":     "기기 %s에서 오래된 파일 공유 프로토콜 SMBv1을 허용합니다. WannaCry와 같은 랜섬웨어는 SMBv1을 통해 가정용 PC와 NAS 기기 사이로 확산되었습니다.",
		"smb1_enabled_fix": "SMBv1을 끄십시오. Windows에서는 “Windows 기능 켜기/끄기”를 열고 “SMB 1.0/CIFS 파일 공유 지원”의 선택을 해제하십시오. NAS나 라우터에서는 설정에서 최소 SMB 버전을 SMB2 이상으로 지정하거나 펌웨어를 업데이트하십시오.",
	},
	"id": {
		"smb1_enabled":     "%s menerima protokol berbagi file lama SMBv1. Ransomware seperti WannaCry menyebar di antara PC rumah dan perangkat NAS melalui SMBv1.",
		"smb1_enabled_fix": "Nonaktifkan SMBv1. Pada Windows, silakan buka “Turn Windows features on or off” dan hapus centang “SMB 1.0/CIFS File Sharing Support”. Pada NAS atau router, Anda dapat mengatur versi SMB minimum ke SMB2 atau yang lebih baru di pengaturan, atau memperbarui firmware.",
	},
}

func init() {
	for lang, catalog := range smb1Messages {
		target := wifiCatalogTargets[lang]
		for key, value := range catalog {
			target[key] = value
		}
		messages[lang] = target
	}
}
