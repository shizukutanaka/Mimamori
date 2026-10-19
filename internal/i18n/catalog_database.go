package i18n

var databaseMessages = map[string]map[string]string{
	"en": {
		"forward_database":     "The database on %s is open to the whole Internet through port %d. Attackers scan for exposed databases to steal or erase the data and demand a ransom.",
		"forward_database_fix": "Remove this port forwarding in the router settings and turn off UPnP. Databases should never be reachable from the Internet; if you need remote access, use the router's VPN function.",
	},
	"ja": {
		"forward_database":     "「%s」のデータベースがポート %d を通じてインターネット全体に公開されています。攻撃者は公開されたデータベースを探し、データを盗んだり消去したりして身代金を要求します。",
		"forward_database_fix": "ルーターの設定でこのポート転送を削除し、UPnP を無効にしてください。データベースをインターネットからアクセスできる状態にしてはいけません。遠隔アクセスが必要な場合は、ルーターの VPN 機能を使ってください。",
	},
	"zh-Hans": {
		"forward_database":     "%s 上的数据库通过端口 %d 向整个互联网开放。攻击者会扫描暴露的数据库，窃取或删除数据并索要赎金。",
		"forward_database_fix": "请在路由器设置中删除此端口转发并关闭 UPnP。数据库绝不应从互联网访问；如果您需要远程访问，请使用路由器的 VPN 功能。",
	},
	"zh-Hant": {
		"forward_database":     "%s 上的資料庫透過連接埠 %d 對整個網際網路開放。攻擊者會掃描暴露的資料庫，竊取或刪除資料並勒索贖金。",
		"forward_database_fix": "請在路由器設定中刪除此連接埠轉送並關閉 UPnP。資料庫絕不應從網際網路存取；如果您需要遠端存取，請使用路由器的 VPN 功能。",
	},
	"es": {
		"forward_database":     "La base de datos de %s está expuesta a todo Internet a través del puerto %d. Los atacantes buscan bases de datos expuestas para robar o borrar los datos y exigir un rescate.",
		"forward_database_fix": "Elimine esta redirección de puertos en la configuración del router y desactive UPnP. Las bases de datos nunca deben ser accesibles desde Internet; si necesita acceso remoto, utilice la función VPN del router.",
	},
	"pt": {
		"forward_database":     "O banco de dados de %s está aberto para toda a Internet pela porta %d. Atacantes procuram bancos de dados expostos para roubar ou apagar os dados e exigir resgate.",
		"forward_database_fix": "Remova este encaminhamento de porta nas configurações do roteador e desative o UPnP. Bancos de dados nunca devem ficar acessíveis pela Internet; se você precisar de acesso remoto, use a função VPN do roteador.",
	},
	"fr": {
		"forward_database":     "La base de données de %s est exposée à l’ensemble d’Internet via le port %d. Les attaquants recherchent les bases exposées pour voler ou effacer les données et réclamer une rançon.",
		"forward_database_fix": "Supprimez cette redirection de port dans les réglages du routeur et désactivez UPnP. Les bases de données ne doivent jamais être accessibles depuis Internet ; si vous avez besoin d’un accès à distance, utilisez la fonction VPN du routeur.",
	},
	"de": {
		"forward_database":     "Die Datenbank auf %s ist über Port %d für das gesamte Internet zugänglich. Angreifer suchen nach offengelegten Datenbanken, um Daten zu stehlen oder zu löschen und Lösegeld zu verlangen.",
		"forward_database_fix": "Entfernen Sie diese Portweiterleitung in den Router-Einstellungen und schalten Sie UPnP aus. Datenbanken dürfen niemals aus dem Internet erreichbar sein. Wenn Sie Fernzugriff benötigen, nutzen Sie die VPN-Funktion des Routers.",
	},
	"ru": {
		"forward_database":     "База данных на устройстве %s доступна всему Интернету через порт %d. Злоумышленники ищут открытые базы данных, чтобы украсть или удалить данные и потребовать выкуп.",
		"forward_database_fix": "Удалите это перенаправление порта в настройках роутера и отключите UPnP. Базы данных не должны быть доступны из Интернета; если вам нужен удалённый доступ, используйте VPN-функцию роутера.",
	},
	"ko": {
		"forward_database":     "%s의 데이터베이스가 포트 %d를 통해 인터넷 전체에 공개되어 있습니다. 공격자는 노출된 데이터베이스를 찾아 데이터를 훔치거나 삭제하고 몸값을 요구합니다.",
		"forward_database_fix": "라우터 설정에서 이 포트 전달을 삭제하고 UPnP를 끄세요. 데이터베이스는 인터넷에서 접근할 수 없어야 합니다. 원격 접속이 필요하다면 라우터의 VPN 기능을 사용하세요.",
	},
	"id": {
		"forward_database":     "Basis data %s terbuka untuk seluruh Internet melalui port %d. Penyerang memindai basis data yang terekspos untuk mencuri atau menghapus data dan meminta tebusan.",
		"forward_database_fix": "Silakan hapus penerusan port ini di pengaturan router dan nonaktifkan UPnP. Basis data tidak boleh diakses dari Internet; jika Anda memerlukan akses jarak jauh, gunakan fitur VPN router.",
	},
}

func init() {
	for lang, catalog := range databaseMessages {
		target := wifiCatalogTargets[lang]
		for key, value := range catalog {
			target[key] = value
		}
		messages[lang] = target
	}
}
