package i18n

var catalogEs = map[string]string{
	"title":                  "Mimamori — Revisión de seguridad del hogar",
	"safe":                   "Segura",
	"caution":                "Precaución",
	"danger":                 "Peligro",
	"score":                  "Seguridad de la red doméstica: %s  %d puntos",
	"device_counts":          "Dispositivos: %d (%d nuevos)",
	"findings_heading":       "Problemas que debe resolver",
	"devices_heading":        "Dispositivos",
	"no_problems":            "No se encontraron problemas",
	"first_run":              "En esta ejecución se guardaron los dispositivos de su red. A partir de la próxima ejecución, se avisará si se conecta un dispositivo nuevo.",
	"new_badge":              "nuevo",
	"router":                 "Router",
	"randomized_device":      "Teléfono o PC (dirección privada de Wi-Fi)",
	"telnet_open":            "Telnet está abierto en %s. Este servicio antiguo de control remoto suele usarse para tomar el control de dispositivos.",
	"telnet_open_fix":        "Desactive Telnet en la configuración del dispositivo e instale las actualizaciones disponibles. Cambie también cualquier contraseña predeterminada.",
	"ftp_open":               "El acceso FTP a archivos está abierto en %s. Esto podría exponer archivos o datos de cuentas.",
	"ftp_open_fix":           "Desactive FTP si no lo necesita. Si lo necesita, use una opción más segura y una contraseña exclusiva.",
	"vnc_open":               "El acceso remoto VNC está abierto en %s. Quien pueda conectarse quizá controle la pantalla.",
	"vnc_open_fix":           "Desactive VNC si no necesita el acceso remoto. Si lo necesita, limite quién puede conectarse y use una contraseña segura.",
	"router_upnp":            "UPnP está activado en el router. Los dispositivos de casa podrían abrir automáticamente una entrada hacia Internet.",
	"router_upnp_fix":        "Si dispositivos como las consolas no lo necesitan, desactive UPnP en la configuración del router.",
	"new_device":             "Se conectó un dispositivo nuevo: %s",
	"new_device_fix":         "Si no lo reconoce, cambie la contraseña de Wi-Fi y vuelva a conectar solo sus propios dispositivos.",
	"port_forward":           "%s puede ser accesible desde Internet mediante el puerto %d (abierto por UPnP).",
	"port_forward_fix":       "Si no reconoce esta regla, elimínela en la configuración del router o desactive UPnP. Las consolas y algunas aplicaciones abren estos puertos a propósito.",
	"port_forward_risky":     "El acceso remoto o el uso compartido de archivos de %s está abierto a todo Internet mediante el puerto %d. Los atacantes buscan estas entradas constantemente.",
	"port_forward_risky_fix": "Elimine esta redirección de puertos en la configuración del router y desactive UPnP. Si necesita acceso remoto, use la función VPN del router.",
	"points_recovered":       "(+%d puntos)",
	"press_enter":            "Pulse Enter para cerrar...",
}

func init() {
	messages["es"] = catalogEs
}
