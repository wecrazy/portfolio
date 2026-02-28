// Package translate provides a lightweight automatic translation client backed
// by the unofficial Google Translate endpoint (translate.googleapis.com).
//
// No API key is required. Results are cached in memory so identical (from, to,
// text) tuples are only round-tripped once per process lifetime.
package translate

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

// Translator translates arbitrary text between languages.
type Translator interface {
	// Translate converts text from the source language to the target language.
	// langFrom and langTo are BCP-47 codes (e.g. "en", "id", "fr").
	Translate(ctx context.Context, text, langFrom, langTo string) (string, error)

	// TranslateMany translates a slice of texts in one call.
	// The returned slice matches the input slice length and order.
	TranslateMany(ctx context.Context, texts []string, langFrom, langTo string) ([]string, error)
}

// ─────────────────────────────────────────────
// Google Translate (unofficial) implementation
// ─────────────────────────────────────────────

// googleTranslateURL is the unofficial Google Translate API endpoint.
// It requires no API key and imposes no strict daily quota.
const googleTranslateURL = "https://translate.googleapis.com/translate_a/single"

type cacheKey struct {
	from, to, text string
}

type googleTranslator struct {
	mu     sync.RWMutex
	cache  map[cacheKey]string
	client *http.Client
}

// New returns a Translator backed by the unofficial Google Translate API with
// built-in in-memory caching.
func New() Translator {
	return &googleTranslator{
		cache: make(map[cacheKey]string),
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (g *googleTranslator) Translate(ctx context.Context, text, langFrom, langTo string) (string, error) {
	// Identical source and target, or empty text → nothing to do.
	if langFrom == langTo || strings.TrimSpace(text) == "" {
		return text, nil
	}

	k := cacheKey{from: langFrom, to: langTo, text: text}

	// Fast path: cache hit.
	g.mu.RLock()
	if v, ok := g.cache[k]; ok {
		g.mu.RUnlock()
		return v, nil
	}
	g.mu.RUnlock()

	// Build request.
	// client=gtx  → returns the nested-array format we parse below.
	// dt=t        → include translated text segments.
	endpoint := fmt.Sprintf("%s?client=gtx&sl=%s&tl=%s&dt=t&q=%s",
		googleTranslateURL,
		url.QueryEscape(langFrom),
		url.QueryEscape(langTo),
		url.QueryEscape(text),
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return text, fmt.Errorf("translate: build request: %w", err)
	}
	// Mimic a browser User-Agent to avoid being blocked.
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; portfolio-bot/1.0)")

	resp, err := g.client.Do(req)
	if err != nil {
		return text, fmt.Errorf("translate: request: %w", err)
	}
	defer resp.Body.Close() //nolint:errcheck

	if resp.StatusCode != http.StatusOK {
		return text, fmt.Errorf("translate: HTTP %d", resp.StatusCode)
	}

	// Response is a deeply nested JSON array, e.g.:
	//   [ [ ["Halo dunia","Hello world",null,null,1], ... ], null, "en" ]
	// We unmarshal into []interface{} and concatenate all first elements of
	// the inner arrays (segments) to reconstruct the full translated string.
	var raw []interface{}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return text, fmt.Errorf("translate: decode response: %w", err)
	}

	translated, err := extractGoogleTranslation(raw)
	if err != nil || translated == "" {
		return text, err
	}

	// Store in cache.
	g.mu.Lock()
	g.cache[k] = translated
	g.mu.Unlock()

	return translated, nil
}

// extractGoogleTranslation pulls the translated text out of the raw response.
// raw[0] is a slice of segments; each segment is [ translated, original, ... ].
func extractGoogleTranslation(raw []interface{}) (string, error) {
	if len(raw) == 0 {
		return "", fmt.Errorf("translate: empty response")
	}
	segments, ok := raw[0].([]interface{})
	if !ok {
		return "", fmt.Errorf("translate: unexpected response shape")
	}

	var b strings.Builder
	for _, seg := range segments {
		parts, ok := seg.([]interface{})
		if !ok || len(parts) == 0 {
			continue
		}
		if s, ok := parts[0].(string); ok {
			b.WriteString(s)
		}
	}
	return b.String(), nil
}

func (g *googleTranslator) TranslateMany(ctx context.Context, texts []string, langFrom, langTo string) ([]string, error) {
	out := make([]string, len(texts))

	type result struct {
		idx int
		val string
		err error
	}
	ch := make(chan result, len(texts))

	for i, t := range texts {
		go func(idx int, txt string) {
			v, err := g.Translate(ctx, txt, langFrom, langTo)
			ch <- result{idx: idx, val: v, err: err}
		}(i, t)
	}

	var firstErr error
	for range texts {
		r := <-ch
		out[r.idx] = r.val
		if r.err != nil && firstErr == nil {
			firstErr = r.err
		}
	}
	return out, firstErr
}
