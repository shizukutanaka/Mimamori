package i18n

import (
	"strings"
	"testing"
)

func TestParseAppleLanguages(t *testing.T) {
	tests := []struct {
		name      string
		output    string
		want      string
		wantMatch string
	}{
		{
			name:   "preferred language first",
			output: "(\n    \"ja-JP\",\n    \"en-US\"\n)",
			want:   "ja-JP",
		},
		{
			name:   "unquoted entry",
			output: "(\n    en,\n    ja\n)",
			want:   "en",
		},
		{
			name:      "Chinese script tag",
			output:    "(\n    \"zh-Hant-TW\",\n    \"en-US\"\n)",
			want:      "zh-Hant-TW",
			wantMatch: "zh-Hant",
		},
		{
			name:   "Portuguese tag",
			output: "(\n    \"pt-BR\"\n)",
			want:   "pt-BR",
		},
		{name: "empty array", output: "()", want: ""},
		{name: "empty input", output: "", want: ""},
		{name: "garbage", output: "not a plist array", want: ""},
		{
			name:   "CRLF line endings",
			output: "(\r\n    \"ja-JP\",\r\n    \"en-US\"\r\n)",
			want:   "ja-JP",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := parseAppleLanguages(test.output)
			if got != test.want {
				t.Fatalf("parseAppleLanguages(%q) = %q, want %q", test.output, got, test.want)
			}
			if test.wantMatch != "" {
				lang, ok := Match(got)
				if !ok || lang != test.wantMatch {
					t.Fatalf("Match(%q) = (%q, %t), want (%q, true)", got, lang, ok, test.wantMatch)
				}
			}
		})
	}
}

func FuzzParseAppleLanguages(f *testing.F) {
	for _, output := range []string{
		"(\n    \"ja-JP\",\n    \"en-US\"\n)",
		"(\n    en,\n    ja\n)",
		"(\n    \"zh-Hant-TW\"\n)",
		"(\n    \"pt-BR\"\n)",
		"()",
		"",
		"not a plist array",
		"(\r\n    \"ja-JP\",\r\n)",
		"(\n    \"\xff\",\n)",
		strings.Repeat("x", 10*1024),
	} {
		f.Add(output)
	}

	f.Fuzz(func(t *testing.T, output string) {
		got := parseAppleLanguages(output)
		if strings.ContainsAny(got, "\",()\r\n") {
			t.Fatalf("parseAppleLanguages(%q) returned invalid entry %q", output, got)
		}
	})
}
