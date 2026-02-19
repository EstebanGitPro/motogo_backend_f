package moderation

import (
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// IsOffensive checks if the given text contains offensive language.
// It tokenizes and normalizes the text, then checks against a dictionary of banned words.
// Returns true if the text contains at least one offensive word or phrase.
func IsOffensive(text string) bool {
	if strings.TrimSpace(text) == "" {
		return false
	}

	normalized := normalizeText(text)

	// Check single tokens
	tokens := strings.Fields(normalized)
	for _, token := range tokens {
		if offensiveWords[token] {
			return true
		}
	}

	// Check bigrams and trigrams (multi-word phrases)
	for _, phrase := range offensivePhrases {
		if strings.Contains(normalized, phrase) {
			return true
		}
	}

	return false
}

// normalizeText converts text to lowercase, removes accents, and strips non-alphanumeric characters.
func normalizeText(text string) string {
	// Lowercase
	lower := strings.ToLower(text)

	// Remove accents using unicode normalization
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	result, _, err := transform.String(t, lower)
	if err != nil {
		// Fallback to lowercase-only if normalization fails
		result = lower
	}

	return result
}

// offensiveWords is a dictionary of single offensive words in Spanish.
// Covers common Colombian Spanish profanity and insults.
var offensiveWords = map[string]bool{
	// Common Spanish profanity
	"mierda":       true,
	"puta":         true,
	"puto":         true,
	"hijueputa":    true,
	"hijuepucha":   true,
	"jueputa":      true,
	"gonorrea":     true,
	"malparido":    true,
	"malparida":    true,
	"cabron":       true,
	"cabrona":      true,
	"pendejo":      true,
	"pendeja":      true,
	"idiota":       true,
	"imbecil":      true,
	"estupido":     true,
	"estupida":     true,
	"basura":       true,
	"desgraciado":  true,
	"desgraciada":  true,
	"maldito":      true,
	"maldita":      true,
	"carajo":       true,
	"verga":        true,
	"culo":         true,
	"marica":       true,
	"maricon":      true,
	"huevon":       true,
	"guevon":       true,
	"cojudo":       true,
	"cerdo":        true,
	"cerda":        true,
	"asqueroso":    true,
	"asquerosa":    true,
	"inutil":       true,
	"mediocre":     true,
	"incompetente": true,
	"ladron":       true,
	"ladrona":      true,
	"estafador":    true,
	"estafadora":   true,
	"tramposo":     true,
	"tramposa":     true,
	"mentiroso":    true,
	"mentirosa":    true,
	"miserable":    true,
	"porqueria":    true,
	"zorra":        true,
	"perro":        true,
	"perra":        true,
	"rata":         true,
	"sapo":         true,
	"sapa":         true,
	"cretino":      true,
	"cretina":      true,
	"baboso":       true,
	"babosa":       true,
	"tarado":       true,
	"tarada":       true,
	"inepto":       true,
	"inepta":       true,
	"inservible":   true,
	"patético":     true,
	"patetico":     true,
	"vergüenza":    true,
	"verguenza":    true,

	// Colombian-specific profanity
	"hijuemadre":  true,
	"malparecido": true,
	"lampara":     true,
	"casposo":     true,
	"casposa":     true,
	"chirrete":    true,
	"ñero":        true,
	"nero":        true,
	"gamín":       true,
	"gamin":       true,

	// Threatening or aggressive terms
	"amenaza":  true,
	"demanda":  true,
	"denuncio": true,
	"denuncia": true,
}

// offensivePhrases contains multi-word offensive expressions.
var offensivePhrases = []string{
	"hijo de puta",
	"hijo de perra",
	"vete a la mierda",
	"que asco",
	"de mierda",
	"son unos ladrones",
	"son ladrones",
	"son estafadores",
	"me robaron",
	"me estafaron",
	"los voy a demandar",
	"los denuncio",
	"una porqueria",
	"una basura",
	"no sirven para nada",
}
