package main

import (
	"flag"
	"strings"
	"time"

	"github.com/shizukutanaka/Mimamori/internal/i18n"
)

type runFlagValues struct {
	jsonOutput  *bool
	language    *string
	statePath   *string
	wait        *time.Duration
	showVersion *bool
}

func registerRunFlags(flags *flag.FlagSet) runFlagValues {
	return runFlagValues{
		jsonOutput:  flags.Bool("json", false, ""),
		language:    flags.String("lang", "auto", ""),
		statePath:   flags.String("state", "", ""),
		wait:        flags.Duration("wait", 3*time.Second, ""),
		showVersion: flags.Bool("version", false, ""),
	}
}

func preLanguage(args []string, getenv func(string) string, flags *flag.FlagSet) string {
	value := ""
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if arg == "--" || arg == "-" || !strings.HasPrefix(arg, "-") {
			break
		}
		nameAndValue := arg[1:]
		if strings.HasPrefix(arg, "--") {
			nameAndValue = arg[2:]
		}
		name, flagValue, hasValue := strings.Cut(nameAndValue, "=")
		registered := flags.Lookup(name)
		if registered == nil {
			break
		}
		if hasValue {
			if name == "lang" {
				value = flagValue
			}
			continue
		}
		if boolFlag, ok := registered.Value.(interface{ IsBoolFlag() bool }); ok && boolFlag.IsBoolFlag() {
			continue
		}
		if i+1 >= len(args) {
			break
		}
		if name == "lang" {
			value = args[i+1]
		}
		i++
	}
	if value != "auto" {
		if lang, ok := i18n.Match(value); ok {
			return lang
		}
	}
	return i18n.Detect(getenv)
}
