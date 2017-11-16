package epub

import "testing"

func TestTransliterate(t *testing.T) {
	for _, tc := range [][]string{
		{"Привет!", "Privet!"},
		{"Съешь еще этих мягких булочек и выпей чаю", "Sesh esche etih myagkih bulochek i vypej chayu"},
		{"Здарова, John! Как life?", "Zdarova, John! Kak life?"},
		{"Hello, 世界", "Hello, 世界"},
	} {
		res := transliterate(tc[0])
		if res != tc[1] {
			t.Errorf("transliteration of %q must be %q, got %q", tc[0], tc[1], res)
		}
	}
}
