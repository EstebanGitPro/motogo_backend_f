package moderation

import "testing"

func TestIsOffensive_CleanText(t *testing.T) {
	tests := []struct {
		name string
		text string
	}{
		{"simple positive", "Excelente servicio, muy recomendado"},
		{"neutral feedback", "El servicio estuvo bien, cumplieron con lo prometido"},
		{"constructive criticism", "Podrían mejorar los tiempos de entrega"},
		{"numbers only", "5 estrellas"},
		{"rating comment", "Buen trabajo, la moto quedó como nueva"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if IsOffensive(tt.text) {
				t.Errorf("IsOffensive(%q) = true, want false", tt.text)
			}
		})
	}
}

func TestIsOffensive_OffensiveText(t *testing.T) {
	tests := []struct {
		name string
		text string
	}{
		{"single word", "Son unos hijueputa"},
		{"colombian slang", "que gonorrea de servicio"},
		{"insult", "son unos malparidos"},
		{"profanity", "mierda de taller"},
		{"aggressive phrase", "Son unos ladrones que me estafaron"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !IsOffensive(tt.text) {
				t.Errorf("IsOffensive(%q) = false, want true", tt.text)
			}
		})
	}
}

func TestIsOffensive_EmptyText(t *testing.T) {
	tests := []struct {
		name string
		text string
	}{
		{"empty string", ""},
		{"only spaces", "   "},
		{"only tabs", "\t\t"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if IsOffensive(tt.text) {
				t.Errorf("IsOffensive(%q) = true, want false", tt.text)
			}
		})
	}
}

func TestIsOffensive_AccentVariations(t *testing.T) {
	tests := []struct {
		name string
		text string
	}{
		{"with accent", "Son unos estúpidos"},
		{"without accent", "Son unos estupidos"},
		{"mixed case", "Son unos IMBECIL"},
		{"mixed accents", "Qué IMBÉCIL de mecánico"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !IsOffensive(tt.text) {
				t.Errorf("IsOffensive(%q) = false, want true", tt.text)
			}
		})
	}
}

func TestIsOffensive_Phrases(t *testing.T) {
	tests := []struct {
		name string
		text string
	}{
		{"theft accusation", "me robaron en este taller, no vuelvo"},
		{"fraud accusation", "me estafaron con el precio"},
		{"phrase in context", "Este taller es una porqueria total"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !IsOffensive(tt.text) {
				t.Errorf("IsOffensive(%q) = false, want true", tt.text)
			}
		})
	}
}

func TestNormalizeText(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"lowercase", "HELLO", "hello"},
		{"remove accents", "estúpido", "estupido"},
		{"mixed", "Ñero Malparído", "nero malparido"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := normalizeText(tt.input)
			if got != tt.want {
				t.Errorf("normalizeText(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
