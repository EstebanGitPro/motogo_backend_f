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
// Supports plural/suffix stripping to catch conjugated forms (e.g. "malparidos" → "malparido").
// Returns true if the text contains at least one offensive word or phrase.
func IsOffensive(text string) bool {
	if strings.TrimSpace(text) == "" {
		return false
	}

	normalized := normalizeText(text)

	// Check single tokens (exact match + stem candidates for plurals)
	tokens := strings.Fields(normalized)
	for _, token := range tokens {
		if offensiveWords[token] {
			return true
		}
		// Try possible singular forms for pluralized words
		for _, candidate := range stemCandidates(token) {
			if offensiveWords[candidate] {
				return true
			}
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

// stemCandidates generates possible singular forms from a pluralized Spanish word.
// It tries the most common Spanish pluralization rules in reverse:
//   - Remove trailing "es" (e.g. "ladrones" → "ladron", "imbeciles" → "imbecil")
//   - Remove trailing "s"  (e.g. "malparidos" → "malparido", "idiotas" → "idiota")
func stemCandidates(word string) []string {
	var candidates []string
	if len(word) > 3 && strings.HasSuffix(word, "es") {
		candidates = append(candidates, strings.TrimSuffix(word, "es"))
	}
	if len(word) > 2 && strings.HasSuffix(word, "s") {
		candidates = append(candidates, strings.TrimSuffix(word, "s"))
	}
	return candidates
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
// NOTE: Plural forms are handled by stripSuffix(), so only singular/base forms are needed here.
var offensiveWords = map[string]bool{
	// ── Common Spanish profanity ──
	"mierda":    true,
	"puta":      true,
	"puto":      true,
	"hijueputa": true,
	"jueputa":   true,
	"gonorrea":  true,
	"malparido": true,
	"malparida": true,
	"cabron":    true,
	"cabrona":   true,
	"pendejo":   true,
	"pendeja":   true,
	"idiota":    true,
	"imbecil":   true,
	"estupido":  true,
	"estupida":  true,
	"carajo":    true,
	"verga":     true,
	"culo":      true,
	"marica":    true,
	"maricon":   true,
	"huevon":    true,
	"guevon":    true,

	// ── Insults and degrading terms ──
	"basura":       true,
	"desgraciado":  true,
	"desgraciada":  true,
	"maldito":      true,
	"maldita":      true,
	"cretino":      true,
	"cretina":      true,
	"baboso":       true,
	"babosa":       true,
	"tarado":       true,
	"tarada":       true,
	"inepto":       true,
	"inepta":       true,
	"inutil":       true,
	"mediocre":     true,
	"incompetente": true,
	"inservible":   true,
	"patetico":     true,
	"patetica":     true,
	"asqueroso":    true,
	"asquerosa":    true,
	"miserable":    true,
	"porqueria":    true,
	"mugre":        true,
	"mugroso":      true,
	"mugrosa":      true,
	"bruto":        true,
	"bruta":        true,
	"tonto":        true,
	"tonta":        true,
	"bobo":         true,
	"boba":         true,
	"animal":       true, // cuando se usa como insulto

	// ── Colombian-specific profanity ──
	"hijuepucha":    true,
	"hijuemadre":    true,
	"malparecido":   true,
	"lampara":       true,
	"casposo":       true,
	"casposa":       true,
	"chirrete":      true,
	"nero":          true,
	"gamin":         true,
	"guache":        true,
	"zarrapastroso": true,
	"zarrapastrosa": true,
	"cacorro":       true,
	"pisca":         true,
	"culicagado":    true,
	"culicagada":    true,
	"guiso":         true,
	"ñero":          true,
	"loba":          true,
	"lobo":          true,
	"lacra":         true,
	"hp":            true, // abbreviation for hijueputa
	"hdp":           true, // hijo de puta abbreviation
	"hpta":          true, // hijueputa abbreviation
	"hptas":         true,
	"mk":            true, // marica abbreviation
	"mrk":           true,
	"gono":          true, // gonorrea abbreviated
	"gonorreas":     true,
	"csm":           true, // abbreviation
	"ptm":           true, // abbreviation
	"ctm":           true, // abbreviation

	// ── Costeño (Costa Caribe) slang ──
	"monda":       true, // vulgar: pene (Barranquilla/Cartagena)
	"mondá":       true,
	"vergajo":     true, // insulto fuerte costeño
	"vergaja":     true,
	"guevona":     true,
	"cachon":      true, // cornudo
	"cachona":     true,
	"cachudo":     true,
	"cachuda":     true,
	"awebao":      true, // ahuevado costeño
	"awebá":       true,
	"ahuevao":     true,
	"aguevao":     true,
	"ñerda":       true, // mierda costeño
	"nerda":       true, // normalized form
	"joda":        true, // molestia vulgar
	"jodedor":     true,
	"jodedora":    true,
	"mamaguevo":   true, // insulto fuerte costeño
	"mamagueva":   true,
	"mamagüevo":   true,
	"mamagüeva":   true,
	"mmgv":        true, // mamaguevo abbreviation
	"mmgvo":       true,
	"coñazo":      true, // golpe / insulto
	"conaso":      true, // normalized
	"coñaso":      true,
	"coño":        true,
	"cono":        true, // normalized form
	"mondao":      true, // pelado/tonto costeño
	"culipronta":  true,
	"culipronto":  true,
	"malandro":    true,
	"malandra":    true,
	"tierruo":     true, // despectivo costeño
	"tierrua":     true,
	"pela":        true, // golpiza
	"pelaito":     true, // despectivo
	"adefesio":    true,
	"guaricha":    true, // despectivo hacia mujeres
	"pajuo":       true, // tonto/pendejo costeño
	"pajua":       true,
	"pajudo":      true,
	"pajuda":      true,
	"monduo":      true,
	"lambon":      true, // lambón
	"lambona":     true,
	"chupamedias": true,

	"pichurria": true, // cosa mala/persona mala
	"arrecho":   true, // vulgar en costa
	"arrecha":   true,

	// ── Animals used as insults ──
	"zorra":   true,
	"perra":   true,
	"perro":   true,
	"cerdo":   true,
	"cerda":   true,
	"rata":    true,
	"sapo":    true,
	"sapa":    true,
	"culebra": true,

	// ── Fraud/theft accusations ──
	"ladron":       true,
	"ladrona":      true,
	"estafador":    true,
	"estafadora":   true,
	"tramposo":     true,
	"tramposa":     true,
	"mentiroso":    true,
	"mentirosa":    true,
	"corrupto":     true,
	"corrupta":     true,
	"sinverguenza": true,
	"cínico":       true,
	"cinico":       true,
	"atracador":    true,
	"atracadora":   true,
	"aprovechado":  true,
	"aprovechada":  true,
	"abusivo":      true,
	"abusiva":      true,

	// ── Sexual/vulgar terms ──
	"pene":   true,
	"vagina": true,
	"teta":   true,
	"nalga":  true,
	"coger":  true, // vulgar in Colombian context
	"cagar":  true,
	"joder":  true,
	"jodido": true,
	"jodida": true,
	"chimba": true, // Colombian vulgar

	"picha":  true,
	"cojudo": true,

	// ── Threatening / aggressive terms ──
	"amenaza":  true,
	"demanda":  true,
	"denuncio": true,
	"denuncia": true,
	"matar":    true,
	"muerte":   true,
}

// offensivePhrases contains multi-word offensive expressions.
var offensivePhrases = []string{
	// ── Profanity phrases ──
	"hijo de puta",
	"hijo de perra",
	"hijos de puta",
	"hijue puta",
	"hijue pucha",
	"vete a la mierda",
	"vaya a la mierda",
	"vayanse a la mierda",
	"que gonorrea",
	"que mierda",
	"que asco",
	"de mierda",
	"de porqueria",
	"una mierda",
	"una porqueria",
	"una basura",
	"pedazo de",
	"manga de",
	"bola de",
	"par de",

	// ── Fraud/theft phrases ──
	"son unos ladrones",
	"son ladrones",
	"son estafadores",
	"son unos estafadores",
	"son unos rateros",
	"me robaron",
	"me estafaron",
	"me tumbaron",
	"me clavaron",
	"me metieron los dedos",
	"los voy a demandar",
	"los denuncio",
	"poner una denuncia",

	// ── Quality insults ──
	"no sirven para nada",
	"no sirve para nada",
	"no saben hacer nada",
	"que asco de servicio",
	"que asco de taller",
	"el peor taller",
	"el peor servicio",
	"la peor experiencia",
	"pesimo servicio",
	"pesima atencion",
}
