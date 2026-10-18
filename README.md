# Mimamori — Home Safety Check

Mimamori answers a simple question: **“Is my home network safe right now?”** It looks for a few high-impact risks and explains what to do, rather than presenting packet details. It is a single static Go binary for Windows, macOS, and Linux; it needs no administrator/root access.

It checks devices visible in the local neighbor table, the default gateway, UPnP device descriptions, a small set of risky services (Telnet, FTP, and VNC), and the connected Wi-Fi security using read-only operating-system tools. It also tries those risky-service ports on every address in the home subnet (up to 254), covering devices this PC has not talked to; it only attempts connections and changes nothing. Router UPnP port forwards are read with a read-only query. The score ranges from 0 to 100; findings explain the recommended fix and points recovered.

Mimamori never changes settings, contacts the Internet, or sends telemetry. Network discovery and checks stay on the private LAN. Its only persistent data is a local device-history file, used to report newly seen devices.

## Download and run

Download the latest release from [GitHub Releases](https://github.com/shizukutanaka/Mimamori/releases/latest). Choose the file for your PC; `<version>` is the release tag shown on that page:

| Your PC | File |
| --- | --- |
| Most Windows PCs | `mimamori-<version>-windows-amd64.zip` |
| Windows on ARM (for example, Snapdragon) | `mimamori-<version>-windows-arm64.zip` |
| Mac with an Apple chip | `mimamori-<version>-darwin-arm64.tar.gz` |
| Intel Mac | `mimamori-<version>-darwin-amd64.tar.gz` |
| Most Linux PCs | `mimamori-<version>-linux-amd64.tar.gz` |
| 64-bit Raspberry Pi | `mimamori-<version>-linux-arm64.tar.gz` |
| 32-bit Raspberry Pi | `mimamori-<version>-linux-arm.tar.gz` |

Download `SHA256SUMS` from the same release to verify the archive. Replace `<archive>` below with the downloaded file's name. Compare the hash printed on Windows or macOS with the matching line in `SHA256SUMS`; on Linux, run the check from the folder containing both files:

```text
Windows: certutil -hashfile <archive> SHA256
macOS:   shasum -a 256 <archive>
Linux:   sha256sum -c SHA256SUMS --ignore-missing
```

Extract the archive before running Mimamori. The binaries are unsigned, so your operating system may show a first-run warning:

- Windows SmartScreen: choose **More info → Run anyway**.
- macOS: choose **System Settings → Privacy & Security → Open Anyway**, or run `xattr -d com.apple.quarantine ./mimamori`.
- Linux: run `chmod +x ./mimamori`.

For maintainers: run `VERSION=v0.1.0 scripts/dist.sh`, then upload `dist/*` to a GitHub Release.

After extracting, run Mimamori from a terminal:

```sh
./mimamori
./mimamori --lang en
./mimamori --lang pt-BR
./mimamori --lang zh_TW
./mimamori --json
./mimamori --state ./devices.json --wait 5s
```

On Windows, run `mimamori.exe` (for example, `.\mimamori.exe --lang ja`) or double-click it in Explorer. When launched by double-click, the window stays open until you press Enter so you can read the results. The first run learns the devices on your network; new-device notices begin on the next run. By default, history is stored under your OS user configuration directory in `Mimamori/devices.json`.

Supported languages: English (`en`), Japanese (`ja`), Simplified Chinese (`zh-Hans`), Traditional Chinese (`zh-Hant`), Spanish (`es`), Portuguese (`pt`), French (`fr`), German (`de`), Russian (`ru`), Korean (`ko`), and Indonesian (`id`). By default, Mimamori checks `LANGUAGE`, `LC_ALL`, `LC_MESSAGES`, then `LANG`; on Windows, it falls back to the display language when these are empty. Use `--lang` with a language code or regional tag, such as `--lang pt-BR` or `--lang zh_TW`, to choose explicitly. Use `--json` for indented JSON output.

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

Exit codes: `0` means no danger finding (caution findings may still be present), `2` means at least one danger finding, and `1` means a flag or operational error; it is also returned when no device could be checked, and a corrupt history file is kept as `devices.json.corrupt` (or `.corrupt.1`, `.corrupt.2`, … if older backups exist).

## Build from source

Go 1.22 or later is required. The project uses only the Go standard library.

```sh
go build -o mimamori .
VERSION=v0.1.0 scripts/dist.sh
```

To add a language, add a catalog file in `internal/i18n` and list its code in `Supported()`.

---

# Mimamori — 家の安全診断

Mimamori は「今、家のネットワークは安全だろうか？」という疑問に答える診断ツールです。通信の細かな情報ではなく、特に注意が必要な状態と対処方法をわかりやすく伝えます。Windows・macOS・Linux で動く、管理者権限がいらない単一の Go バイナリです。

近くで確認できる機器、標準ゲートウェイ、UPnP の機器情報、危険につながりやすい Telnet・FTP・VNC の入口に加え、OS の機能だけを使って接続中 Wi-Fi の安全方式も読み取り専用で確認します。さらに、自宅サブネット内の全アドレス（最大 254 件）にも危険なサービスのポートへの接続を試み、この PC と通信したことがない機器も確認します。接続を試みるだけで、設定は変更しません。ルーターの UPnP ポート転送も読み取り専用の問い合わせで調べます。点数は 0〜100 点で、見つかった問題ごとに対処方法と回復する点数を表示します。

設定を変更せず、インターネットへ接続せず、利用状況も送信しません。機器の確認は家のネットワーク内だけで行います。保存するのは、新しい機器を次回から知らせるための機器履歴ファイルだけです。

## ダウンロードと実行

お使いの PC に合うファイルを [GitHub Releases](https://github.com/shizukutanaka/Mimamori/releases/latest) からダウンロードしてください。`<version>` はリリースページに表示されているタグです。

| PC | ファイル |
| --- | --- |
| ほとんどの Windows PC | `mimamori-<version>-windows-amd64.zip` |
| ARM 版 Windows（Snapdragon など） | `mimamori-<version>-windows-arm64.zip` |
| Apple チップ搭載 Mac | `mimamori-<version>-darwin-arm64.tar.gz` |
| Intel Mac | `mimamori-<version>-darwin-amd64.tar.gz` |
| ほとんどの Linux PC | `mimamori-<version>-linux-amd64.tar.gz` |
| 64 ビット Raspberry Pi | `mimamori-<version>-linux-arm64.tar.gz` |
| 32 ビット Raspberry Pi | `mimamori-<version>-linux-arm.tar.gz` |

同じリリースにある `SHA256SUMS` もダウンロードして、ファイルが正しいか確認できます。コマンド中の `<ファイル名>` はダウンロードしたファイル名に置き換えてください。Windows と macOS では表示されたハッシュを `SHA256SUMS` の該当行と比べてください。Linux では両方のファイルがあるフォルダーで実行します。

```text
Windows: certutil -hashfile <ファイル名> SHA256
macOS:   shasum -a 256 <ファイル名>
Linux:   sha256sum -c SHA256SUMS --ignore-missing
```

アーカイブを展開してから実行してください。バイナリは署名されていないため、初回起動時に OS の警告が出ることがあります。

- Windows SmartScreen: **「詳細情報」→「実行」** を選びます。
- macOS: **システム設定 → プライバシーとセキュリティ →「このまま開く」** を選ぶか、`xattr -d com.apple.quarantine ./mimamori` を実行します。
- Linux: `chmod +x ./mimamori` を実行します。

メンテナー向け: `VERSION=v0.1.0 scripts/dist.sh` を実行し、`dist/*` を GitHub Release にアップロードしてください。

展開後、ターミナルから実行します。

```sh
./mimamori
./mimamori --lang ja
./mimamori --lang pt-BR
./mimamori --lang zh_TW
./mimamori --json
./mimamori --state ./devices.json --wait 5s
```

Windows では `mimamori.exe` を実行してください（例: `.\mimamori.exe --lang ja`）。エクスプローラーからダブルクリックしても起動でき、その場合は結果を読めるよう Enter キーを押すまでウィンドウが開いたままになります。初回は家の機器を記録し、新しい機器のお知らせは次回から表示します。履歴は通常、OS のユーザー設定フォルダーにある `Mimamori/devices.json` に保存されます。

対応言語: 英語（`en`）、日本語（`ja`）、簡体字中国語（`zh-Hans`）、繁体字中国語（`zh-Hant`）、スペイン語（`es`）、ポルトガル語（`pt`）、フランス語（`fr`）、ドイツ語（`de`）、ロシア語（`ru`）、韓国語（`ko`）、インドネシア語（`id`）。言語は `LANGUAGE`、`LC_ALL`、`LC_MESSAGES`、`LANG` の順で自動判定し、Windows ではこれらが空の場合に表示言語を使います。`--lang pt-BR` や `--lang zh_TW` のように、言語コードや地域タグで指定できます。`--json` で整形済み JSON を出力できます。

終了コード: `0` は危険な問題なし（注意事項が表示される場合があります）、`2` は危険な問題あり、`1` はオプションまたは動作上のエラーです。確認できる機器がなかった場合にも `1` が返され、壊れた履歴ファイルは `devices.json.corrupt`（古いバックアップがある場合は `.corrupt.1`、`.corrupt.2` など）として残されます。

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
VERSION=v0.1.0 scripts/dist.sh
```

言語を追加する場合は `internal/i18n` にカタログファイルを追加し、`Supported()` にコードを登録してください。
