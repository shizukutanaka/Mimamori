# Mimamori — Home Safety Check

Mimamori answers a simple question: **“Is my home network safe right now?”** It looks for a few high-impact risks and explains what to do, rather than presenting packet details. It is a single static Go binary for Windows, macOS, and Linux; it needs no administrator/root access. It also checks whether another private address uses the router's MAC address, a possible sign of ARP spoofing.

It checks devices visible in the local neighbor table, the default gateway, UPnP device descriptions, and the connected Wi-Fi security using read-only operating-system tools. On every address in the home subnet (up to 254), including devices this PC has not talked to, it checks a small set of risky services: Telnet (TCP 23/2323), FTP (21), VNC (5900), Android Debug Bridge (ADB, 5555), Windows file sharing (SMB, 445), RTSP cameras (554/8554), MQTT brokers (1883), the go2rtc camera relay's web API (1984), the Frigate NVR's internal API (5000), the Node-RED editor (1880), the ESPHome dashboard (6052), the Docker Engine API (2375, and 2376 without client-certificate checks), Home Assistant (8123), Zigbee2MQTT (8080), Moonraker (7125), Tasmota and Shelly smart devices (80), WLED light controllers (80), Pi-hole v6 (80), the Homebridge UI (8581), the Syncthing web GUI (8384), Transmission's remote control (9091), Jupyter (8888), Ollama (11434), and code-server (8443). The PC running Mimamori is not checked, because a connection to itself does not show what other devices can reach; to check it, run Mimamori on another computer on the same network. For FTP, VNC, ADB, RTSP, MQTT, go2rtc, Frigate, Node-RED, ESPHome, Docker, Home Assistant, Zigbee2MQTT, Moonraker, Tasmota, Shelly, WLED, Pi-hole v6, Homebridge, Syncthing, Jupyter, and Ollama, it then sends standard read-only requests without a password to see whether the service lets anyone in (FTP's anonymous login, a VNC or ADB handshake, an RTSP `DESCRIBE`, an MQTT `CONNECT`, a go2rtc `GET /api`, Frigate's `GET /api/stats`, Node-RED's `GET /settings`, ESPHome's `GET /devices`, Docker's `GET /version` followed by `OPTIONS /containers/create`, a safe method that creates nothing and has no request body (both requests are also sent over TLS on port 2376; a 401 or 403 refusal, as from a read-only docker-socket-proxy, is not flagged), Home Assistant's `GET /api/onboarding`, a WebSocket handshake to Zigbee2MQTT's `/api` that only reads the first message, a `GET /server/info` to Moonraker, Tasmota's `GET /` start page, Shelly's `GET /shelly` (sent over HTTPS to the same device if port 80 redirects there), WLED: `GET /json/info`, then `GET /update` (200 only when no Settings PIN and no OTA lock), Pi-hole v6's one read-only `GET /api/auth` (warns only when `session.message` is "no password set"; Pi-hole v5 is not checked), Homebridge's `GET /api/auth/settings` (warns only when `formAuth` is false), or Syncthing's one `GET /meta.js` (warns only when it answers without a password; GUIs using HTTPS are not checked), or Transmission's one `GET /transmission/rpc` (warns only when it answers with a session ID and no password; no RPC command is sent), or Jupyter: one `GET /api/status`; warns only when it answers without a token; Ollama: `GET /` and `GET /api/version` only; never sends prompts or deletes; code-server: `GET /healthz` and `GET /login` only; warns only when login redirects home (no password)). For SMB, it sends one connection request to see whether the device still accepts the old SMBv1 protocol. It never guesses passwords, opens files, watches video, or sends commands, and it changes nothing. Router UPnP port forwards are read with a read-only query, and devices whose UPnP description reports a known-vulnerable libupnp or MiniUPnPd version are flagged. Many routers list only forwards that apps created through UPnP, so forwards set by hand in the router's admin page may not appear; check those in the router's settings. The score ranges from 0 to 100; findings explain the recommended fix and points recovered.

Mimamori never changes settings, contacts the Internet, or sends telemetry. Network discovery and checks stay on the private LAN. Its only persistent data is a local device-history file, used to report newly seen devices. It stores each device's MAC address, last IP address, name, and first/last seen dates at `%APPDATA%\Mimamori\devices.json` (Windows), `~/Library/Application Support/Mimamori/devices.json` (macOS), or `~/.config/Mimamori/devices.json` (Linux), readable only by your user account. Delete that file to forget all devices, for example before giving the PC away.

On Linux, the Wi-Fi security check uses `nmcli`, falling back to `wpa_cli status` when NetworkManager is unavailable.

Run Mimamori only on networks you own or manage, such as your home network. On café, hotel, school, or work networks it would probe other people's devices, which may break that network's rules or set off its intrusion alarms.

## Download and run

Download the latest release from [GitHub Releases](https://github.com/shizukutanaka/Mimamori/releases/latest). Choose the file for your PC; `<version>` is the release tag shown on that page:

| Your PC | File |
| --- | --- |
| Most Windows PCs | `mimamori-<version>-windows-amd64.zip` |
| Windows on ARM (for example, Snapdragon) | `mimamori-<version>-windows-arm64.zip` |
| Mac with an Apple chip | `mimamori-<version>-darwin-arm64.tar.gz` |
| Intel Mac | `mimamori-<version>-darwin-amd64.tar.gz` |
| Most Linux PCs | `mimamori-<version>-linux-amd64.tar.gz` |
| Older 32-bit Linux PCs | `mimamori-<version>-linux-386.tar.gz` |
| 64-bit Raspberry Pi | `mimamori-<version>-linux-arm64.tar.gz` |
| 32-bit Raspberry Pi | `mimamori-<version>-linux-arm.tar.gz` |

Requirements: Windows 10 or later, macOS 13 Ventura or later, or Linux with kernel 3.2 or later (32-bit Raspberry Pi: ARMv6 or later). Windows 7/8 and macOS 12 or earlier are not supported, because the Go toolchain used for the release builds no longer supports them.

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

On Windows, run `mimamori.exe` (for example, `.\mimamori.exe --lang ja`) or double-click it in Explorer. When launched by double-click, the window stays open until you press Enter so you can read the results. The first run learns the devices on your network; new-device notices begin on the next run. By default, history is stored under your OS user configuration directory in `Mimamori/devices.json`. The history is not separated by network. If you also check another network you manage (for example, your parents' home), use a separate file for it, such as `--state ./parents.json`. Otherwise, every device there is reported as new on the first check, and those devices are no longer reported as new if they later appear at home.

On macOS, Terminal and SSH sessions can reach your home network without extra permission. Other apps that run Mimamori, such as iTerm2 or the VS Code terminal, need Local Network access: allow it when macOS asks, or turn it on for that app in **System Settings → Privacy & Security → Local Network**. Without it, Mimamori cannot check the devices on your network.

Supported languages: English (`en`), Japanese (`ja`), Simplified Chinese (`zh-Hans`), Traditional Chinese (`zh-Hant`), Spanish (`es`), Portuguese (`pt`), French (`fr`), German (`de`), Russian (`ru`), Korean (`ko`), and Indonesian (`id`). By default, Mimamori checks `LANGUAGE`, `LC_ALL`, `LC_MESSAGES`, then `LANG`; when these are empty, it uses the display language on Windows and the preferred language (System Settings → General → Language & Region) on macOS. Use `--lang` with a language code or regional tag, such as `--lang pt-BR` or `--lang zh_TW`, to choose explicitly. Use `--json` for indented JSON output. Its `checked` object (`ports`, `port_forwards`, `wifi`) shows which checks actually ran; `false` means not checked, not safe.

Many fixes say "in the router settings": open a browser and go to the address shown on the Router line of the device list (for example `http://192.168.1.1`), or use your router maker's app. The login details are usually printed on a label on the router.

Example:

```text
Mimamori — Home Safety Check  v0.1.0
Home safety: Danger  75 points
Devices: 4 (0 new)

Things to address
  [Danger] Telnet is open on Camera. This old remote-control service is often abused to take over devices. (192.168.1.30)
         → Turn off Telnet in the device settings, update its firmware, and change any default password.  (+25 points)

Devices
  192.168.1.1      aa:bb:cc:dd:ee:ff  Router
```

If a device is listed only by its IP address, Mimamori could not learn its name. Your router's admin page or app usually lists connected devices with the same IP or MAC address; alternatively, unplug the device you suspect, wait a minute or two (your PC remembers recently seen devices for a short while), and run Mimamori again to see which line disappears.

Exit codes: `0` means no danger finding (caution findings may still be present), `2` means at least one danger finding, and `1` means a flag or operational error; it is also returned when no device could be checked, and a corrupt history file is kept as `devices.json.corrupt` (or `.corrupt.1`, `.corrupt.2`, … if older backups exist). A danger finding still returns `2` even if the device history cannot be saved. `0` is also returned when part of the check could not run; the text output then shows an incomplete notice and `--json` sets `"incomplete": true`, so scheduled scripts should check that field too.

On Windows, a run started directly by Task Scheduler may also wait for Enter, as with a double-click, so the task keeps running; add `--json` for scheduled runs (it never waits).

## Build from source

Go 1.22 or later is required. The project uses only the Go standard library.

```sh
go build -o mimamori .
VERSION=v0.1.0 scripts/dist.sh
```

To add a language, add a catalog file in `internal/i18n` and list its code in `Supported()`.

---

# Mimamori — 家の安全診断

Mimamori は「今、家のネットワークは安全だろうか？」という疑問に答える診断ツールです。通信の細かな情報ではなく、特に注意が必要な状態と対処方法をわかりやすく伝えます。Windows・macOS・Linux で動く、管理者権限がいらない単一の Go バイナリです。ルーターと同じ MAC アドレスを使う別のプライベート IP アドレスがないかも確認し、ARP スプーフィングの兆候を調べます。

近くで確認できる機器、標準ゲートウェイ、UPnP の機器情報に加え、OS の機能だけを使って接続中 Wi-Fi の安全方式も読み取り専用で確認します。さらに、自宅サブネット内の全アドレス（最大 254 件）について、この PC と通信したことがない機器も含め、危険につながりやすいサービスの入口を確認します。対象は Telnet（TCP 23・2323）、FTP（21）、VNC（5900）、Android Debug Bridge（ADB、5555）、Windows のファイル共有の SMB（445）、防犯カメラの RTSP（554・8554）、スマートホームの MQTT（1883）、カメラの映像を中継する go2rtc の管理画面（1984）、録画ソフト Frigate の内部用の窓口（5000）、スマートホームの自動化に使う Node-RED の編集画面（1880）、ESPHome の管理画面（6052）、Docker の操作窓口（2375番、クライアント証明書を確かめない2376番）、Home Assistant（8123）、Zigbee2MQTT（8080）、Moonraker（7125）、スマート家電の Tasmota・Shelly（80）、WLED の照明制御機器（80）、Pi-hole v6（80）、Homebridge の管理画面（8581）、Syncthing の管理画面（8384）、Transmission の操作窓口（9091）、Jupyter（8888）、Ollama（11434）、code-server（8443）です。診断している PC 自身は調べません。自分への接続では、ほかの機器からどう見えるかが分からないためです。この PC を調べるには、同じネットワークの別のパソコンで Mimamori を実行してください。FTP・VNC・ADB・RTSP・MQTT・go2rtc・Frigate・Node-RED・ESPHome・Docker・Home Assistant・Zigbee2MQTT・Moonraker・Tasmota・Shelly・WLED・Pi-hole v6・Homebridge・Syncthing・Jupyter・Ollama では、パスワードなしで誰でも使える状態かを確かめるため、読み取り専用の標準手順で問い合わせます（FTP の匿名ログイン、VNC・ADB の接続開始、RTSP の `DESCRIBE`、MQTT の `CONNECT`、go2rtc の `GET /api`、Frigate の `GET /api/stats`、Node-RED の `GET /settings`、ESPHome の `GET /devices`、Docker の `GET /version` の後に何も作らずリクエスト本文もない安全な `OPTIONS /containers/create`（2376番では両方とも TLS で送信し、読み取り専用 docker-socket-proxy などから 401 または 403 で拒否された場合は警告しません）、Home Assistant の `GET /api/onboarding`、Zigbee2MQTT の `/api` への WebSocket 接続（最初の1通を読むだけ）、Moonraker への `GET /server/info`、Tasmota の最初のページの `GET /`、Shelly の `GET /shelly`（80番が同じ機器の HTTPS へ移すときはそちらへ）、WLED: `GET /json/info` の後に `GET /update`（Settings PIN と OTA ロックがない場合のみ 200）、Pi-hole v6 の読み取り専用の `GET /api/auth` を 1 回だけ送信（`session.message` が "no password set" の場合のみ警告し、Pi-hole v5 は対象外）、Homebridge の `GET /api/auth/settings`（`formAuth` が false の場合のみ警告）、Syncthing の `GET /meta.js` を 1 回だけ送信（パスワードなしで答えた場合のみ警告し、HTTPS の管理画面は対象外）、Transmission の `GET /transmission/rpc` を 1 回だけ送信（パスワードなしでセッション ID を返した場合のみ警告し、命令は送らない）、Jupyter の `GET /api/status` を 1 回だけ送信（token なしで答えた場合のみ警告）、Ollama の `GET /` と `GET /api/version` だけを送信（AI への問い合わせや削除はしない）、code-server の GET /healthz と GET /login だけを送信（ログインなしで / に移す場合のみ警告））。SMB では、古い SMBv1 をまだ受け付けるかを確かめるため、接続の申し込みを 1 回だけ送ります。パスワードを推測したり、ファイルや映像を開いたり、操作の命令を送ったりはせず、設定も変更しません。ルーターの UPnP ポート転送も読み取り専用の問い合わせで調べ、既知の脆弱性がある古い libupnp・MiniUPnPd を使う機器も知らせます。多くのルーターは、アプリが UPnP で作った転送だけを返すため、ルーターの設定画面で手で作った転送は表示されないことがあります。それらはルーターの設定画面で確かめてください。点数は 0〜100 点で、見つかった問題ごとに対処方法と回復する点数を表示します。

設定を変更せず、インターネットへ接続せず、利用状況も送信しません。機器の確認は家のネットワーク内だけで行います。保存するのは、新しい機器を次回から知らせるための機器履歴ファイルだけです。このファイルには、機器ごとの MAC アドレス・最後の IP アドレス・名前・初めて見た日と最後に見た日が入り、`%APPDATA%\Mimamori\devices.json`（Windows）、`~/Library/Application Support/Mimamori/devices.json`（macOS）、`~/.config/Mimamori/devices.json`（Linux）に、自分のアカウントだけが読める権限で置かれます。PC を手放すときなど、記録を消したい場合はこのファイルを削除してください。

Linux では Wi-Fi の安全方式を `nmcli` で確認し、NetworkManager が利用できない場合は `wpa_cli status` を使います。

Mimamori は、自宅など自分で管理しているネットワークでだけ実行してください。カフェ・ホテル・学校・会社などのネットワークで実行すると、他人の機器のポートまで確認することになり、そのネットワークの規約に反したり、侵入検知の警報が鳴ったりするおそれがあります。

## ダウンロードと実行

お使いの PC に合うファイルを [GitHub Releases](https://github.com/shizukutanaka/Mimamori/releases/latest) からダウンロードしてください。`<version>` はリリースページに表示されているタグです。

| PC | ファイル |
| --- | --- |
| ほとんどの Windows PC | `mimamori-<version>-windows-amd64.zip` |
| ARM 版 Windows（Snapdragon など） | `mimamori-<version>-windows-arm64.zip` |
| Apple チップ搭載 Mac | `mimamori-<version>-darwin-arm64.tar.gz` |
| Intel Mac | `mimamori-<version>-darwin-amd64.tar.gz` |
| ほとんどの Linux PC | `mimamori-<version>-linux-amd64.tar.gz` |
| 32 ビットの古い Linux PC | `mimamori-<version>-linux-386.tar.gz` |
| 64 ビット Raspberry Pi | `mimamori-<version>-linux-arm64.tar.gz` |
| 32 ビット Raspberry Pi | `mimamori-<version>-linux-arm.tar.gz` |

動作環境: Windows 10 以降、macOS 13 Ventura 以降、または Linux カーネル 3.2 以降（32 ビット Raspberry Pi は ARMv6 以降）。Windows 7・8 と macOS 12 以前には対応していません。配布物のビルドに使う Go がこれらに対応していないためです。

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

Windows では `mimamori.exe` を実行してください（例: `.\mimamori.exe --lang ja`）。エクスプローラーからダブルクリックしても起動でき、その場合は結果を読めるよう Enter キーを押すまでウィンドウが開いたままになります。初回は家の機器を記録し、新しい機器のお知らせは次回から表示します。履歴は通常、OS のユーザー設定フォルダーにある `Mimamori/devices.json` に保存されます。履歴はネットワークごとに分かれていません。実家など、自分で管理する別のネットワークも確認する場合は、`--state ./parents.json` のように別のファイルを指定してください。指定しないと、そのネットワークの機器が初回からすべて「新しい機器」として表示され、その機器が後で家に来ても「新しい機器」として表示されなくなります。

macOS では、「ターミナル」と SSH からは、追加の許可なしで家のネットワークを確認できます。iTerm2 や VS Code のターミナルなど、ほかのアプリから実行する場合は、「ローカルネットワーク」の許可が必要です。macOS に尋ねられたら許可するか、**システム設定 → プライバシーとセキュリティ → ローカルネットワーク** で、そのアプリをオンにしてください。許可しないと、家の機器を確認できません。

対応言語: 英語（`en`）、日本語（`ja`）、簡体字中国語（`zh-Hans`）、繁体字中国語（`zh-Hant`）、スペイン語（`es`）、ポルトガル語（`pt`）、フランス語（`fr`）、ドイツ語（`de`）、ロシア語（`ru`）、韓国語（`ko`）、インドネシア語（`id`）。言語は `LANGUAGE`、`LC_ALL`、`LC_MESSAGES`、`LANG` の順で自動判定し、これらが空の場合は Windows では表示言語、macOS では「システム設定 → 一般 → 言語と地域」の優先する言語を使います。`--lang pt-BR` や `--lang zh_TW` のように、言語コードや地域タグで指定できます。`--json` で整形済み JSON を出力できます。JSON の `checked`（`ports`・`port_forwards`・`wifi`）は実際に確認できた項目を示し、`false` は「安全」ではなく「確認できなかった」という意味です。

機器が IP アドレスだけで表示されるのは、名前を取得できなかった機器です。ルーターの管理画面やアプリの接続機器一覧に、同じ IP アドレスか MAC アドレスで載っていることが多いです。見当がつかなければ、疑わしい機器の電源を切り、1〜2 分待ってから（PC は最近見た機器をしばらく覚えています）もう一度実行し、消えた行を確かめてください。

終了コード: `0` は危険な問題なし（注意事項が表示される場合があります）、`2` は危険な問題あり、`1` はオプションまたは動作上のエラーです。確認できる機器がなかった場合にも `1` が返され、壊れた履歴ファイルは `devices.json.corrupt`（古いバックアップがある場合は `.corrupt.1`、`.corrupt.2` など）として残されます。履歴を保存できなかった場合でも、危険な問題があれば `2` を返します。確認の一部ができなかった場合も `0` を返します。その場合は、画面に「不完全な可能性」の案内が表示され、`--json` では `"incomplete": true` になります。定期実行のスクリプトでは、この項目も確かめてください。

Windows では、タスク スケジューラーから直接起動した場合も、ダブルクリックと同じく Enter キーを待ち、タスクが終わらないことがあります。定期実行では `--json` を付けてください（`--json` では待ちません）。

直し方に「ルーターの設定で」とあるときは、機器一覧の「ルーター」の行にあるアドレス（例: `http://192.168.1.1`）をブラウザで開くか、ルーターのメーカーのアプリを使います。ログイン情報はルーター本体のラベルに書かれていることが多いです。

```text
Mimamori — 家の安全診断  v0.1.0
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
