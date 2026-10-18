package i18n

var catalogPt = map[string]string{
	"title":                  "Mimamori — Verificação de segurança da rede doméstica",
	"safe":                   "Segura",
	"caution":                "Atenção",
	"danger":                 "Perigo",
	"score":                  "Segurança da rede doméstica: %s  %d pontos",
	"device_counts":          "Dispositivos: %d (%d novos)",
	"findings_heading":       "Problemas a resolver",
	"devices_heading":        "Dispositivos",
	"no_problems":            "Nenhum problema encontrado",
	"first_run":              "Nesta execução, os dispositivos da sua rede foram registrados. A partir da próxima execução, você será avisado quando um novo dispositivo se conectar.",
	"new_badge":              "novo",
	"router":                 "Roteador",
	"randomized_device":      "Celular ou computador (endereço privado de Wi-Fi)",
	"telnet_open":            "O Telnet está ativo em %s. Esse serviço antigo de controle remoto costuma ser usado para assumir o controle de dispositivos.",
	"telnet_open_fix":        "Desative o Telnet nas configurações do dispositivo e instale as atualizações disponíveis. Troque também qualquer senha padrão.",
	"ftp_open":               "O acesso FTP a arquivos está ativo em %s. Isso pode expor arquivos ou dados de contas.",
	"ftp_open_fix":           "Desative o FTP se não precisar dele. Se precisar, use uma opção mais segura e uma senha exclusiva.",
	"vnc_open":               "O acesso remoto VNC está ativo em %s. Quem conseguir se conectar talvez controle a tela.",
	"vnc_open_fix":           "Desative o VNC se não precisar de acesso remoto. Se precisar, limite quem pode se conectar e use uma senha forte.",
	"router_upnp":            "O UPnP está ativo no roteador. Dispositivos da casa podem abrir automaticamente uma entrada para a Internet.",
	"router_upnp_fix":        "Se consoles e outros dispositivos não precisarem, desative o UPnP nas configurações do roteador.",
	"new_device":             "Um novo dispositivo se conectou: %s",
	"new_device_fix":         "Se não reconhecer o dispositivo, troque a senha do Wi-Fi e conecte novamente apenas os seus aparelhos.",
	"port_forward":           "%s pode ser acessado pela Internet pela porta %d (aberta pelo UPnP).",
	"port_forward_fix":       "Se não reconhecer esta regra, remova-a nas configurações do roteador ou desative o UPnP. Consoles e alguns aplicativos abrem essas portas de propósito.",
	"port_forward_risky":     "O acesso remoto ou compartilhamento de arquivos em %s está aberto para toda a Internet pela porta %d. Atacantes procuram essas entradas o tempo todo.",
	"port_forward_risky_fix": "Remova esta regra de encaminhamento de portas nas configurações do roteador e desative o UPnP. Se precisar de acesso remoto, use a função VPN do roteador.",
	"points_recovered":       "(+%d pontos)",
	"press_enter":            "Pressione Enter para fechar...",
}

func init() {
	messages["pt"] = catalogPt
}
