# Security Policy

Mimamori is a home-network safety check. People run it to decide whether their home is safe, so a bug that makes it miss a risk or crash can do real harm. Thank you for reporting it privately.

## How to report

Please do **not** open a public issue with vulnerability details.

1. Use GitHub's private reporting: open the [Security tab](https://github.com/shizukutanaka/Mimamori/security) and choose **Report a vulnerability**.
2. If that button is not available, open an issue titled "Security contact request" with no details. The maintainer will reply with a private channel.

Include the Mimamori version (`mimamori -version`), your OS, and the steps or input needed to reproduce. Reports in any language are welcome.

## Supported versions

Only the latest release receives security fixes.

## In scope

- Crafted responses from devices on the LAN, or crafted Wi-Fi names (SSIDs), that crash Mimamori, hang it, or hide a finding it should report.
- A result that shows "safe" even though a check could not run, instead of marking the result as incomplete.
- Mimamori sending anything outside the local network, changing any setting, or needing administrator or root rights.
- The local device history file (`devices.json` in the user configuration directory) being readable by other users, or a corrupted file causing data loss.

## Out of scope

- Vulnerabilities in the devices that Mimamori finds. Report those to the device vendor.
- Missing checks or false positives that are not caused by crafted input. Please open a normal issue.

---

## 脆弱性の報告（日本語）

脆弱性の詳細は、公開の Issue には書かないでください。[Security タブ](https://github.com/shizukutanaka/Mimamori/security)の **Report a vulnerability** から非公開で報告してください。このボタンがない場合は、詳細を書かずに「Security contact request」という題名で Issue を作成してください。非公開の連絡方法をお知らせします。
