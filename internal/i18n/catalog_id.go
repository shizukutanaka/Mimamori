package i18n

var catalogId = map[string]string{
	"title":                  "Mimamori — Pemeriksaan Keamanan Jaringan Rumah",
	"safe":                   "Aman",
	"caution":                "Perlu perhatian",
	"danger":                 "Bahaya",
	"score":                  "Keamanan jaringan rumah: %s  %d poin",
	"device_counts":          "Perangkat: %d (%d baru)",
	"findings_heading":       "Hal yang perlu ditangani",
	"devices_heading":        "Perangkat",
	"no_problems":            "Tidak ada masalah yang ditemukan",
	"first_run":              "Perangkat di jaringan Anda dicatat pada pemeriksaan ini. Perangkat baru akan dilaporkan mulai pemeriksaan berikutnya.",
	"new_badge":              "baru",
	"router":                 "Router",
	"randomized_device":      "Ponsel atau PC (alamat Wi-Fi pribadi)",
	"telnet_open":            "Telnet terbuka di %s. Layanan kendali jarak jauh lama ini sering disalahgunakan untuk mengambil alih perangkat.",
	"telnet_open_fix":        "Matikan Telnet di pengaturan perangkat dan pasang pembaruan yang tersedia. Ganti juga kata sandi bawaan jika belum diubah.",
	"ftp_open":               "Akses file melalui FTP terbuka di %s. File atau informasi akun dapat terlihat oleh orang lain.",
	"ftp_open_fix":           "Matikan FTP jika tidak diperlukan. Jika diperlukan, gunakan cara yang lebih aman dan kata sandi khusus.",
	"vnc_open":               "Akses jarak jauh VNC terbuka di %s. Orang yang dapat terhubung mungkin bisa mengendalikan layar.",
	"vnc_open_fix":           "Matikan VNC jika akses jarak jauh tidak diperlukan. Jika diperlukan, batasi perangkat yang dapat terhubung dan gunakan kata sandi kuat.",
	"router_upnp":            "UPnP di router aktif. Perangkat di rumah dapat membuka jalur ke Internet secara otomatis.",
	"router_upnp_fix":        "Jika konsol game dan perangkat lain tidak memerlukannya, matikan UPnP di pengaturan router.",
	"new_device":             "Perangkat baru terhubung: %s",
	"new_device_fix":         "Jika Anda tidak mengenalinya, ganti kata sandi Wi-Fi dan sambungkan kembali hanya perangkat milik Anda.",
	"port_forward":           "%s dapat diakses dari Internet melalui port %d yang dibuka oleh UPnP.",
	"port_forward_fix":       "Jika Anda tidak mengenali aturan ini, hapus di pengaturan router atau matikan UPnP. Konsol game dan beberapa aplikasi sengaja membuka port seperti ini.",
	"port_forward_risky":     "Akses jarak jauh atau berbagi file di %s terbuka untuk seluruh Internet melalui port %d. Penyerang terus-menerus memindai akses seperti ini.",
	"port_forward_risky_fix": "Hapus penerusan port ini di pengaturan router dan matikan UPnP. Jika perlu akses jarak jauh, gunakan fitur VPN router.",
	"points_recovered":       "(+%d poin)",
	"press_enter":            "Tekan Enter untuk menutup...",
}

func init() {
	messages["id"] = catalogId
}
