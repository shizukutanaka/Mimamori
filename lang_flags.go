package main

import (
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
