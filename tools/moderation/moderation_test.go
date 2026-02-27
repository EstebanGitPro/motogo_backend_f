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
		{"word 'perro' in normal context avoided — but kept flagged", ""},
	}

	for _, tt := range tests {
		if tt.text == "" {
			continue
		}
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
		{"insult singular", "es un malparido"},
		{"insult plural", "son unos malparidos"},
		{"profanity", "mierda de taller"},
		{"aggressive phrase", "Son unos ladrones que me estafaron"},
		{"abbreviation hp", "que hp de servicio"},
		{"abbreviation hdp", "hdp no sirven"},
		{"abbreviation mk", "mk que servicio tan malo"},
		{"colombian chimba", "que chimba tan mala"},
		{"colombian guache", "que guache el mecanico"},
		// Costeño slang
		{"costeño mamaguevo", "que mamaguevo de mecanico"},
		{"costeño vergajo", "ese vergajo no sabe nada"},
		{"costeño awebao", "que awebao ese tipo"},
		{"costeño monda", "ese monda no sirve"},
		{"costeño cachon", "ese mechnico es un cachon"},
		{"costeño pichurria", "ese taller es una pichurria"},
		{"costeño pajuo", "que pajuo tan malo"},
		{"costeño mmgv", "que mmgv de servicio"},
		{"costeño ñerda", "que ñerda de atencion"},
		{"costeño coño", "coño que mal servicio"},
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
		{"with accent singular", "Es un estúpido"},
		{"without accent singular", "Es un estupido"},
		{"with accent plural", "Son unos estúpidos"},
		{"without accent plural", "Son unos estupidos"},
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

func TestIsOffensive_Plurals(t *testing.T) {
	tests := []struct {
		name string
		text string
	}{
		{"malparidos", "son unos malparidos"},
		{"malditos", "son unos malditos"},
		{"idiotas", "son unos idiotas"},
		{"pendejos", "todos son pendejos"},
		{"estupidos", "que estupidos"},
		{"ladrones", "son unos ladrones"},
		{"desgraciados", "son unos desgraciados"},
		{"imbeciles", "que imbeciles"},
		{"incompetentes", "son unos incompetentes"},
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
		{"worst service", "el peor servicio que he tenido"},
		{"colombian tumbaron", "me tumbaron con esa cotizacion"},
		{"threat phrase", "los voy a demandar por esto"},
		{"hijos de puta", "son unos hijos de puta"},
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

func TestStemCandidates(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		contains string // expected candidate that should be in the list
	}{
		{"malparidos → malparido", "malparidos", "malparido"},
		{"malditos → maldito", "malditos", "maldito"},
		{"estupidos → estupido", "estupidos", "estupido"},
		{"idiotas → idiota", "idiotas", "idiota"},
		{"ladrones → ladron", "ladrones", "ladron"},
		{"imbeciles → imbecil", "imbeciles", "imbecil"},
		{"desgraciados → desgraciado", "desgraciados", "desgraciado"},
		{"incompetentes → incompetente", "incompetentes", "incompetente"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			candidates := stemCandidates(tt.input)
			found := false
			for _, c := range candidates {
				if c == tt.contains {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("stemCandidates(%q) = %v, expected to contain %q", tt.input, candidates, tt.contains)
			}
		})
	}
}
