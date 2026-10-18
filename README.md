# Mimamori — Home Safety Check

Mimamori answers a simple question: **“Is my home network safe right now?”** It looks for a few high-impact risks and explains what to do, rather than presenting packet details. It is a single static Go binary for Windows, macOS, and Linux; it needs no administrator/root access.

It checks devices visible in the local neighbor table, the default gateway, UPnP device descriptions, and a small set of risky services (Telnet, FTP, and VNC). The score ranges from 0 to 100; findings explain the recommended fix and points recovered.

Mimamori never changes settings, contacts the Internet, or sends telemetry. Network discovery and checks stay on the private LAN. Its only persistent data is a local device-history file, used to report newly seen devices.

## Download and run

Download the binary for your operating system and run it from a terminal:

```sh
./mimamori
./mimamori --lang en
./mimamori --json
./mimamori --state ./devices.json --wait 5s
```

On Windows, run `mimamori.exe` (for example, `.\mimamori.exe --lang ja`) or double-click it in Explorer. When launched by double-click, the window stays open until you press Enter so you can read the results. The first run learns the devices on your network; new-device notices begin on the next run. By default, history is stored under your OS user configuration directory in `Mimamori/devices.json`.

Language is detected from the system locale by default. Use `--lang en` or `--lang ja` to choose explicitly, and `--json` for indented JSON output.

Example:

```text
Mimamori — Home Safety Check
Home safety: Danger  75 points
Devices: 4 (0 new)

Things to address
  [Danger] Telnet is open on Camera. This old remote-control service is often abused to take over devices. (192.168.1.30)
         → Turn off Telnet in the device settings, update its firmware, and change any default password.  (+25 points)

Devices
  192.168.1.1      aa:bb:cc:dd:ee:ff  Router
```

Exit codes: `0` means no danger finding (caution findings may still be present), `2` means at least one danger finding, and `1` means a flag or operational error.

## Build from source

Go 1.22 or later is required. The project uses only the Go standard library.

```sh
go build -o mimamori .
./scripts/dist.sh
```

To add a language, add one catalog map in `internal/i18n/i18n.go`.

---

# Mimamori — 家の安全診断

Mimamori は「今、家のネットワークは安全だろうか？」という疑問に答える診断ツールです。通信の細かな情報ではなく、特に注意が必要な状態と対処方法をわかりやすく伝えます。Windows・macOS・Linux で動く、管理者権限がいらない単一の Go バイナリです。

近くで確認できる機器、標準ゲートウェイ、UPnP の機器情報を調べ、危険につながりやすい Telnet・FTP・VNC の入口を確認します。点数は 0〜100 点で、見つかった問題ごとに対処方法と回復する点数を表示します。

設定を変更せず、インターネットへ接続せず、利用状況も送信しません。機器の確認は家のネットワーク内だけで行います。保存するのは、新しい機器を次回から知らせるための機器履歴ファイルだけです。

## ダウンロードと実行

お使いの OS 向けのバイナリをダウンロードし、ターミナルから実行します。

```sh
./mimamori
./mimamori --lang ja
./mimamori --json
./mimamori --state ./devices.json --wait 5s
```

Windows では `mimamori.exe` を実行してください（例: `.\mimamori.exe --lang ja`）。エクスプローラーからダブルクリックしても起動でき、その場合は結果を読めるよう Enter キーを押すまでウィンドウが開いたままになります。初回は家の機器を記録し、新しい機器のお知らせは次回から表示します。履歴は通常、OS のユーザー設定フォルダーにある `Mimamori/devices.json` に保存されます。

言語はシステムの設定から自動で選ばれます。`--lang en` または `--lang ja` で指定でき、`--json` で整形済み JSON を出力できます。

終了コード: `0` は危険な問題なし（注意事項が表示される場合があります）、`2` は危険な問題あり、`1` はオプションまたは動作上のエラーです。

```text
Mimamori — 家の安全診断
家の安全状態: 安全  100点
機器: 4台（新しい機器 0台）

危険はありません
```

## ソースからビルド

Go 1.22 以降が必要です。Go 標準ライブラリ以外の依存はありません。

```sh
go build -o mimamori .
./scripts/dist.sh
```

言語を追加する場合は `internal/i18n/i18n.go` のメッセージ一覧に 1 つの map を追加してください。
