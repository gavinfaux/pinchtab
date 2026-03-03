package semantic

// stopwords is a set of common English words that carry little semantic
// meaning and should be excluded from lexical matching to improve
// signal-to-noise ratio.
var stopwords = map[string]bool{
	"the": true, "a": true, "an": true, "is": true, "are": true,
	"was": true, "were": true, "be": true, "been": true, "being": true,
	"have": true, "has": true, "had": true, "do": true, "does": true,
	"did": true, "will": true, "would": true, "could": true, "should": true,
	"may": true, "might": true, "shall": true, "can": true,
	"to": true, "of": true, "in": true, "for": true, "on": true,
	"with": true, "at": true, "by": true, "from": true, "as": true,
	"into": true, "through": true, "about": true, "above": true,
	"after": true, "before": true, "between": true, "under": true,
	"and": true, "but": true, "or": true, "nor": true, "not": true,
	"so": true, "yet": true, "both": true, "either": true, "neither": true,
	"this": true, "that": true, "these": true, "those": true,
	"it": true, "its": true, "i": true, "me": true, "my": true,
	"we": true, "our": true, "you": true, "your": true, "he": true,
	"she": true, "his": true, "her": true, "they": true, "their": true,
}

// isStopword returns true if the token is a common English stopword.
func isStopword(token string) bool {
	return stopwords[token]
}

// removeStopwords filters out stopwords from a token list.
// If removal would empty the list, the original tokens are returned
// to avoid zero-signal matching.
func removeStopwords(tokens []string) []string {
	filtered := make([]string, 0, len(tokens))
	for _, t := range tokens {
		if !isStopword(t) {
			filtered = append(filtered, t)
		}
	}
	if len(filtered) == 0 {
		return tokens
	}
	return filtered
}
