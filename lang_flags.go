package main

import (
	"fmt"
	"strings"

	"github.com/shizukutanaka/Mimamori/internal/i18n"
)

func preLanguage(args []string, getenv func(string) string) string {
	value := ""
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if arg == "--" {
			break
		}
		if arg == "-lang" || arg == "--lang" {
			value = ""
			if i+1 < len(args) && args[i+1] != "--" {
				value = args[i+1]
				i++
			}
			continue
		}
		name, flagValue, hasValue := strings.Cut(arg, "=")
		if hasValue && (name == "-lang" || name == "--lang") {
			value = flagValue
		}
	}
	if value != "auto" {
		if lang, ok := i18n.Match(value); ok {
			return lang
		}
	}
	return i18n.Detect(getenv)
}

func resolveLanguage(flagValue, lang string, getenv func(string) string) (string, error) {
	if flagValue == "auto" {
		return i18n.Detect(getenv), nil
	}
	if lang, ok := i18n.Match(flagValue); ok {
		return lang, nil
	}
	return "", fmt.Errorf("%s", i18n.T(lang, "invalid_lang", flagValue, strings.Join(i18n.Supported(), ", ")))
}
