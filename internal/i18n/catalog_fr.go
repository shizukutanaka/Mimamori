package i18n

var catalogFr = map[string]string{
	"title":                  "Mimamori — Vérification de la sécurité du réseau domestique",
	"safe":                   "Sûre",
	"caution":                "À surveiller",
	"danger":                 "Danger",
	"score":                  "Sécurité du réseau domestique : %s  %d points",
	"device_counts":          "Appareils : %d (%d nouveaux)",
	"findings_heading":       "Problèmes à résoudre",
	"devices_heading":        "Appareils",
	"no_problems":            "Aucun problème détecté",
	"first_run":              "Les appareils de votre réseau ont été enregistrés. Les nouveaux appareils seront signalés à partir de la prochaine vérification.",
	"new_badge":              "nouveau",
	"router":                 "Routeur",
	"randomized_device":      "Téléphone ou ordinateur (adresse Wi-Fi privée)",
	"telnet_open":            "Telnet est activé sur %s. Ce service de commande à distance, ancien et souvent détourné, peut permettre de prendre le contrôle d’un appareil.",
	"telnet_open_fix":        "Désactivez Telnet dans les réglages de l’appareil et installez les mises à jour proposées. Changez aussi tout mot de passe par défaut.",
	"ftp_open":               "L’accès aux fichiers par FTP est activé sur %s. Des fichiers ou des données de compte pourraient être exposés.",
	"ftp_open_fix":           "Désactivez FTP si vous n’en avez pas besoin. Sinon, choisissez une méthode plus sûre et un mot de passe unique.",
	"vnc_open":               "L’accès à distance VNC est activé sur %s. Une personne qui peut s’y connecter pourrait contrôler l’écran.",
	"vnc_open_fix":           "Désactivez VNC si vous n’avez pas besoin de l’accès à distance. Sinon, limitez les connexions et choisissez un mot de passe robuste.",
	"router_upnp":            "UPnP est activé sur le routeur. Les appareils de la maison peuvent ouvrir automatiquement un accès depuis Internet.",
	"router_upnp_fix":        "Si les consoles de jeu et autres appareils n’en ont pas besoin, désactivez UPnP dans les réglages du routeur.",
	"new_device":             "Un nouvel appareil s’est connecté : %s",
	"new_device_fix":         "Si vous ne le reconnaissez pas, changez le mot de passe Wi-Fi et reconnectez uniquement vos propres appareils.",
	"port_forward":           "%s est accessible depuis Internet via le port %d (ouvert par UPnP).",
	"port_forward_fix":       "Si vous ne reconnaissez pas cette règle, supprimez-la dans les réglages du routeur ou désactivez UPnP. Les consoles de jeu et certaines applications ouvrent parfois ces ports volontairement.",
	"port_forward_risky":     "L’accès à distance ou le partage de fichiers sur %s est ouvert à tout Internet via le port %d. Les pirates recherchent constamment ce type d’accès.",
	"port_forward_risky_fix": "Supprimez cette redirection de port dans les réglages du routeur et désactivez UPnP. Si vous avez besoin d’un accès à distance, utilisez plutôt la fonction VPN du routeur.",
	"points_recovered":       "(+%d points)",
	"press_enter":            "Appuyez sur Entrée pour fermer...",
}

func init() {
	messages["fr"] = catalogFr
}
