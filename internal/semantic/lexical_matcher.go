package semantic

import (
	"context"
	"sort"
	"strings"
	"unicode"
)

// LexicalMatcher implements ElementMatcher using Jaccard similarity
// with stopword removal, token frequency weighting, and role-aware boosting.
// Zero external dependencies.
type LexicalMatcher struct{}

// NewLexicalMatcher creates a new LexicalMatcher.
func NewLexicalMatcher() *LexicalMatcher {
	return &LexicalMatcher{}
}

// Strategy returns "lexical".
func (m *LexicalMatcher) Strategy() string { return "lexical" }

// Find scores all elements against the query using lexical similarity,
// filters by threshold, sorts descending, and returns the top-K matches.
func (m *LexicalMatcher) Find(_ context.Context, query string, elements []ElementDescriptor, opts FindOptions) (FindResult, error) {
	if opts.TopK <= 0 {
		opts.TopK = 3
	}

	type scored struct {
		desc  ElementDescriptor
		score float64
	}

	var candidates []scored
	for _, el := range elements {
		composite := el.Composite()
		score := LexicalScore(query, composite)
		if score >= opts.Threshold {
			candidates = append(candidates, scored{desc: el, score: score})
		}
	}

	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].score > candidates[j].score
	})

	if len(candidates) > opts.TopK {
		candidates = candidates[:opts.TopK]
	}

	result := FindResult{
		Strategy:     "lexical",
		ElementCount: len(elements),
	}

	for _, c := range candidates {
		result.Matches = append(result.Matches, ElementMatch{
			Ref:   c.desc.Ref,
			Score: c.score,
			Role:  c.desc.Role,
			Name:  c.desc.Name,
		})
	}

	if len(result.Matches) > 0 {
		result.BestRef = result.Matches[0].Ref
		result.BestScore = result.Matches[0].Score
	}

	return result, nil
}

// --- lexical scoring internals ---

// tokenize splits a string into lowercase tokens, removing punctuation.
func tokenize(s string) []string {
	s = strings.ToLower(s)
	return strings.FieldsFunc(s, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	})
}

// tokenFreq returns token → count map.
func tokenFreq(tokens []string) map[string]int {
	m := make(map[string]int, len(tokens))
	for _, t := range tokens {
		m[t]++
	}
	return m
}

// tokenSet converts a slice of tokens to a set (map).
func tokenSet(tokens []string) map[string]bool {
	m := make(map[string]bool, len(tokens))
	for _, t := range tokens {
		m[t] = true
	}
	return m
}

// roleKeywords are element roles that carry strong semantic signal.
var roleKeywords = map[string]bool{
	"button":   true,
	"input":    true,
	"link":     true,
	"submit":   true,
	"form":     true,
	"textbox":  true,
	"checkbox": true,
	"radio":    true,
	"select":   true,
	"option":   true,
	"tab":      true,
	"menu":     true,
	"search":   true,
}

// LexicalScore computes a similarity between a query and an element
// description using Jaccard overlap on tokens with:
//   - lowercase normalization
//   - stopword removal
//   - token frequency weighting (repeated tokens count proportionally)
//   - role keyword boost (+0.15 if a role keyword overlaps)
//
// Returns a value in [0, 1].
func LexicalScore(query, desc string) float64 {
	qTokens := removeStopwords(tokenize(query))
	dTokens := removeStopwords(tokenize(desc))

	if len(qTokens) == 0 || len(dTokens) == 0 {
		return 0
	}

	qFreq := tokenFreq(qTokens)
	dFreq := tokenFreq(dTokens)

	// Weighted intersection: min(freq_q, freq_d) for each shared token.
	var intersectW float64
	for t, qc := range qFreq {
		if dc, ok := dFreq[t]; ok {
			minC := qc
			if dc < minC {
				minC = dc
			}
			intersectW += float64(minC)
		}
	}

	// Weighted union: max(freq_q, freq_d) for each token in either set.
	allTokens := tokenSet(append(qTokens, dTokens...))
	var unionW float64
	for t := range allTokens {
		qc := qFreq[t]
		dc := dFreq[t]
		maxC := qc
		if dc > maxC {
			maxC = dc
		}
		unionW += float64(maxC)
	}

	if unionW == 0 {
		return 0
	}

	jaccard := intersectW / unionW

	// Role boost: if a role keyword appears in both query and description.
	roleBoost := 0.0
	qSet := tokenSet(qTokens)
	dSet := tokenSet(dTokens)
	for t := range qSet {
		if roleKeywords[t] && dSet[t] {
			roleBoost = 0.15
			break
		}
	}

	score := jaccard + roleBoost
	if score > 1.0 {
		score = 1.0
	}
	return score
}
