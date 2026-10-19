package i18n

import (
	"strings"
	"testing"
)

func TestZigbee2MQTTCatalogCoverage(t *testing.T) {
	for _, lang := range Supported() {
		message, ok := messages[lang]["zigbee2mqtt_noauth"]
		if !ok || message == "" {
			t.Errorf("%s catalog is missing %q", lang, "zigbee2mqtt_noauth")
		}
		fix, ok := messages[lang]["zigbee2mqtt_noauth_fix"]
		if !ok || fix == "" {
			t.Errorf("%s catalog is missing %q", lang, "zigbee2mqtt_noauth_fix")
		}
		if count := strings.Count(message, "%s"); count != 1 {
			t.Errorf("%s zigbee2mqtt_noauth has %d %%s placeholders, want 1", lang, count)
		}
		if strings.Contains(fix, "%s") {
			t.Errorf("%s zigbee2mqtt_noauth_fix has an unexpected %%s placeholder", lang)
		}
	}
	if got := T("en", "zigbee2mqtt_noauth", "Office"); got != "The Zigbee2MQTT web interface on Office lets anyone on your network control your Zigbee devices and run programs without a password." {
		t.Errorf("English Zigbee2MQTT message = %q", got)
	}
	if got := T("en", "zigbee2mqtt_noauth_fix"); got != "In Zigbee2MQTT's configuration.yaml, set frontend.auth_token to a long random value, then restart Zigbee2MQTT." {
		t.Errorf("English Zigbee2MQTT fix = %q", got)
	}
	if got := T("ja", "zigbee2mqtt_noauth", "居間"); got != "「居間」の Zigbee2MQTT の管理画面が、パスワードなしで家の中の誰でも Zigbee 機器を操作したり、プログラムを実行させたりできる状態です。" {
		t.Errorf("Japanese Zigbee2MQTT message = %q", got)
	}
	if got := T("ja", "zigbee2mqtt_noauth_fix"); got != "Zigbee2MQTT の configuration.yaml で frontend.auth_token に長いランダムな値を設定し、Zigbee2MQTT を再起動してください。" {
		t.Errorf("Japanese Zigbee2MQTT fix = %q", got)
	}
}
