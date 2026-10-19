package i18n

var nodeREDMessages = map[string]map[string]string{
	"en": {
		"nodered_noauth":     "Anyone can use the Node-RED editor (%s) without a password (port 1880)",
		"nodered_noauth_fix": "Set adminAuth in Node-RED's settings.js to require a password, then restart Node-RED.",
	},
	"ja": {
		"nodered_noauth":     "Node-RED（%s）の編集画面を、パスワードなしで誰でも操作できます（1880番）",
		"nodered_noauth_fix": "Node-RED の settings.js で adminAuth を設定してパスワードを付け、Node-RED を再起動してください。",
	},
	"zh-Hans": {
		"nodered_noauth":     "任何人都可以在没有密码的情况下使用 Node-RED 编辑器（%s）（端口 1880）。",
		"nodered_noauth_fix": "在 Node-RED 的 settings.js 中设置 adminAuth 以启用密码，然后重启 Node-RED。",
	},
	"zh-Hant": {
		"nodered_noauth":     "任何人都可以在沒有密碼的情況下使用 Node-RED 編輯器（%s）（連接埠 1880）。",
		"nodered_noauth_fix": "在 Node-RED 的 settings.js 中設定 adminAuth 以啟用密碼，然後重新啟動 Node-RED。",
	},
	"es": {
		"nodered_noauth":     "Cualquiera puede usar el editor de Node-RED (%s) sin contraseña (puerto 1880).",
		"nodered_noauth_fix": "Configure adminAuth en settings.js de Node-RED para exigir una contraseña y reinicie Node-RED.",
	},
	"pt": {
		"nodered_noauth":     "Qualquer pessoa pode usar o editor do Node-RED (%s) sem senha (porta 1880).",
		"nodered_noauth_fix": "Configure adminAuth no settings.js do Node-RED para exigir uma senha e reinicie o Node-RED.",
	},
	"fr": {
		"nodered_noauth":     "N’importe qui peut utiliser l’éditeur Node-RED (%s) sans mot de passe (port 1880).",
		"nodered_noauth_fix": "Configurez adminAuth dans settings.js de Node-RED pour exiger un mot de passe, puis redémarrez Node-RED.",
	},
	"de": {
		"nodered_noauth":     "Jede Person kann den Node-RED-Editor (%s) ohne Passwort verwenden (Port 1880).",
		"nodered_noauth_fix": "Legen Sie adminAuth in Node-REDs settings.js fest, um ein Passwort zu verlangen, und starten Sie Node-RED neu.",
	},
	"ru": {
		"nodered_noauth":     "Любой может использовать редактор Node-RED (%s) без пароля (порт 1880).",
		"nodered_noauth_fix": "Настройте adminAuth в settings.js Node-RED, чтобы требовать пароль, затем перезапустите Node-RED.",
	},
	"ko": {
		"nodered_noauth":     "누구나 비밀번호 없이 Node-RED 편집기(%s)를 사용할 수 있습니다(포트 1880).",
		"nodered_noauth_fix": "Node-RED의 settings.js에서 adminAuth를 설정해 비밀번호를 요구하고 Node-RED를 다시 시작하세요.",
	},
	"id": {
		"nodered_noauth":     "Siapa pun dapat menggunakan editor Node-RED (%s) tanpa kata sandi (port 1880).",
		"nodered_noauth_fix": "Atur adminAuth di settings.js Node-RED agar meminta kata sandi, lalu mulai ulang Node-RED.",
	},
}

var nodeREDCatalogTargets = map[string]map[string]string{
	"en":      messages["en"],
	"ja":      messages["ja"],
	"zh-Hans": catalogZhHans,
	"zh-Hant": catalogZhHant,
	"es":      catalogEs,
	"pt":      catalogPt,
	"fr":      catalogFr,
	"de":      catalogDe,
	"ru":      catalogRu,
	"ko":      catalogKo,
	"id":      catalogId,
}

func init() {
	for lang, catalog := range nodeREDMessages {
		target := nodeREDCatalogTargets[lang]
		for key, value := range catalog {
			target[key] = value
		}
		messages[lang] = target
	}
}
