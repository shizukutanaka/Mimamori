# Translating Mimamori

This guide describes the translation code in `internal/i18n`.

## How Mimamori chooses a language

`--lang <code>` selects a language explicitly. The default is `auto`, which calls `i18n.Detect`.

In automatic mode, `i18n.Detect` checks `LANGUAGE` first (trying its colon-separated values in order), then the first nonempty variable in this order: `LC_ALL`, `LC_MESSAGES`, `LANG`. If that variable is not supported, the result is English. If none is set, Mimamori checks the system locale; on Windows, `systemLocale` reads the Windows UI language. English is the final fallback.

`i18n.Match` ignores case and surrounding whitespace, and accepts locale suffixes such as `.UTF-8` and `@latin`. It matches regular language tags by their primary language, so `pt_BR` maps to `pt`. For Chinese, an explicit `Hans` or `Hant` script takes precedence; `TW`, `HK`, and `MO` select `zh-Hant`, and other Chinese tags default to `zh-Hans` (`zh_TW` → `zh-Hant`, `zh_CN` → `zh-Hans`).

## Where strings live

The English and Japanese base messages are in `internal/i18n/i18n.go`. Other language catalogs are in `internal/i18n/catalog_<code>.go`, for example `catalog_fr.go`.

Feature-specific strings live in `internal/i18n/catalog_<feature>.go`, one file per check or screen (for example `catalog_wifi.go`, `catalog_notice.go`, `catalog_nodered.go`; about 35 files). Each holds a map keyed by language code and registers it with `add(...)` in its `init`, which merges the entries into `messages`. There are no per-file language lists to maintain.

## Fixing a translation

Change the translated value, not its key or the surrounding Go structure. Keep every format verb—such as `%s` and `%d`—in the same order as the English value. `TestCatalogKeysAndFormatVerbs` checks that each language has the same keys and format-verb sequence as English.

## Keep the existing register

Use the register already established in that language’s catalog:

- Spanish: formal **usted**
- French: **vous**
- German: formal **Sie**
- Portuguese: **você**
- Russian: **вы**
- Korean: polite **합니다/세요**
- Indonesian: **Anda**
- Chinese: **您**

Matching the catalog’s existing register is required, including when translating a new message.

## Adding a language

1. Add its code to `Supported()` in `internal/i18n/i18n.go`.
2. `i18n.Match` already maps most tags by their primary language code. Add special handling in `Match` only if the language needs it; `i18n.Detect` uses `Match`, so add detection logic only if the new language needs special environment or system-locale handling.
3. Add `internal/i18n/catalog_<code>.go` holding a `var catalog<Code> = map[string]string{...}` with every base message, and list it in the `messages` literal in `i18n.go` next to the other languages (`"fr": catalogFr,`). English and Japanese are defined inline in that literal; no `init` is needed.
4. Add a block for the new language code to every feature catalog `catalog_<feature>.go` (`ls internal/i18n/catalog_*.go`, excluding the language files and tests). `TestCatalogKeysAndFormatVerbs` fails for each key the new language is missing, so it doubles as the checklist.
5. Update `wantLanguages` and any relevant `Match`/`Detect` cases in `internal/i18n/i18n_test.go`. The catalog tests must pass for all base and feature keys, with matching format verbs.

Run `go test ./internal/i18n/` and `go test ./...` after the changes.

## Check your work

Run these commands with the new language code:

```sh
go test ./internal/i18n/
go run . --lang <code>
go run . --lang <code> --json
```

The first command checks catalog coverage and language matching. The text run checks the rendered output; the JSON run checks the localized messages in findings.
