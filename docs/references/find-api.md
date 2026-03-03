# Semantic Find API (`/find`)

Find UI elements by natural language description using semantic matching.

## Endpoint

```
POST /find
POST /find?tab={tabId}
```

## Request Body

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `query` | string | *required* | Natural language description of the target element |
| `threshold` | float | `0.3` | Minimum score threshold (0.0–1.0) |
| `topK` | int | `3` | Maximum matches to return |

### Example Request

```bash
curl -s -X POST http://localhost:9868/find \
  -H "Content-Type: application/json" \
  -d '{"query": "submit button", "threshold": 0.3, "topK": 3}'
```

## Response

```json
{
  "best_ref": "e5",
  "score": 0.85,
  "confidence": "high",
  "strategy": "combined:lexical+embedding:hashing",
  "element_count": 42,
  "threshold": 0.3,
  "latency_ms": 2,
  "matches": [
    { "ref": "e5", "score": 0.85, "role": "button", "name": "Submit" },
    { "ref": "e12", "score": 0.42, "role": "button", "name": "Reset Form" }
  ]
}
```

### Response Fields

| Field | Type | Description |
|-------|------|-------------|
| `best_ref` | string | Ref of the highest-scoring element (use with `/action`) |
| `score` | float | Best match score (0.0–1.0) |
| `confidence` | string | `"high"` (≥0.8), `"medium"` (≥0.6), or `"low"` |
| `strategy` | string | Matching strategy used |
| `element_count` | int | Total elements evaluated on the page |
| `threshold` | float | Echo of the threshold used |
| `latency_ms` | int | Processing time in milliseconds |
| `matches` | array | Top-K matches above threshold, sorted by score desc |

## Matching Strategy

PinchTab uses a **combined matcher** that fuses two scoring methods:

### 1. Lexical Matching (60% weight)

Jaccard-based token overlap with:
- **Stopword removal** — "the", "a", "is" are ignored
- **Token frequency weighting** — rare tokens score higher
- **Role keyword boost** — matches on "button", "textbox", "link" etc. get a 0.15 bonus

### 2. Embedding Matching (40% weight)

Feature-hashing embedder (128-dim vectors) with:
- **Word unigrams** — exact word overlap
- **Character n-grams** (2–4) — sub-word similarity (e.g., "btn" ↔ "button")
- **Role-aware features** — extra dimensions for UI role keywords
- **Cosine similarity** — normalized vector comparison

### Combined Score

```
final_score = 0.6 × lexical_score + 0.4 × embedding_score
```

## Usage with /action

Chain `/find` → `/action` for natural language element interaction:

```bash
# 1. Find the element
REF=$(curl -s -X POST http://localhost:9868/find \
  -d '{"query": "login button"}' | jq -r '.best_ref')

# 2. Act on it
curl -s -X POST http://localhost:9868/action \
  -d "{\"ref\": \"$REF\", \"action\": \"click\"}"
```

## Confidence Levels

| Level | Score Range | Guidance |
|-------|------------|----------|
| `high` | ≥ 0.80 | Safe to act automatically |
| `medium` | 0.60 – 0.79 | Likely correct, may want confirmation |
| `low` | < 0.60 | Multiple candidates possible, re-query advised |

## Error Responses

| Code | Condition |
|------|-----------|
| `400` | Missing `query` field |
| `404` | Tab not found |
| `500` | No accessibility snapshot available (call `/text` first) |

## Performance

Benchmarks on Intel i5-4300U @ 1.90GHz:

| Operation | Elements | Latency | Allocations |
|-----------|----------|---------|-------------|
| Lexical Find | 16 | ~71 µs | 134 allocs |
| HashingEmbedder (single) | 1 | ~11 µs | 3 allocs |
| HashingEmbedder (batch) | 16 | ~171 µs | 49 allocs |
| Embedding Find | 16 | ~180 µs | 98 allocs |
| **Combined Find** | **16** | **~233 µs** | **263 allocs** |
| Combined Find | 100 | ~1.5 ms | 1685 allocs |
