package search

import (
	"context"
	"fmt"
)

// Query contains all parameters for full-text search, filtering and pagination.
type Query struct {
	Keyword   string         // Main search keyword
	Filters   map[string]any // Exact filters (e.g., "status": "published")
	SortBy    string         // Field used for sorting (e.g., "created_at", "price")
	Ascending bool           // Sorting ascending or descending
	Limit     int            // Limit the number of records returned
	Offset    int            // Offset starting index to retrieve records
}

// Hit represents a matching search result document.
type Hit struct {
	ID       string  // Document ID
	Score    float64 // Relevancy score (keyword match)
	Document []byte  // Raw document data as JSON
}

// Result contains matching results and statistical metadata.
type Result struct {
	Hits      []Hit // List of matching documents
	TotalHits int64 // Total matching documents in the system
	TookMs    int64 // Search execution time (milliseconds)
}

// ProviderType defines the search provider.
type ProviderType string

const (
	ProviderMeiliSearch   ProviderType = "meilisearch"
	ProviderElasticSearch ProviderType = "elasticsearch"
	ProviderAlgolia       ProviderType = "algolia"
	ProviderInMemory      ProviderType = "inmemory"
)

// Config holds the configuration parameters for the search engine.
type Config struct {
	Provider ProviderType
	URL      string // Endpoint connection URL
	APIKey   string // Access API key / credentials
}

// Engine defines a common interface for all Search Engines (Elasticsearch, Meilisearch, Algolia, etc.).
type Engine interface {
	// Index indexes a new document or replaces the entire old document.
	Index(ctx context.Context, indexName string, documentID string, doc any) error

	// Update updates a portion of data of an existing document.
	Update(ctx context.Context, indexName string, documentID string, doc any) error

	// Delete deletes a document from the search index.
	Delete(ctx context.Context, indexName string, documentID string) error

	// Search performs full-text search based on keyword, filter and pagination.
	Search(ctx context.Context, indexName string, q Query) (*Result, error)

	// Close flushes and releases any resources used by the search engine.
	Close() error
}

// NewEngine is a Factory function to create the corresponding search Engine.
func NewEngine(cfg Config) (Engine, error) {
	switch cfg.Provider {
	case ProviderMeiliSearch:
		return newMeiliSearchEngine(cfg)
	case ProviderElasticSearch:
		return newElasticSearchEngine(cfg)
	case ProviderAlgolia:
		return newAlgoliaEngine(cfg)
	case ProviderInMemory:
		return newInMemoryEngine(cfg)
	default:
		return nil, fmt.Errorf("search: unsupported provider: %s", cfg.Provider)
	}
}

// ── MeiliSearch Implementation (Stub) ────────────────────────────────────────

type meilisearchEngine struct {
	cfg Config
}

func newMeiliSearchEngine(cfg Config) (*meilisearchEngine, error) {
	// TODO: Initialize Meilisearch client
	return &meilisearchEngine{cfg: cfg}, nil
}

func (e *meilisearchEngine) Index(ctx context.Context, indexName string, documentID string, doc any) error {
	return nil
}

func (e *meilisearchEngine) Update(ctx context.Context, indexName string, documentID string, doc any) error {
	return nil
}

func (e *meilisearchEngine) Delete(ctx context.Context, indexName string, documentID string) error {
	return nil
}

func (e *meilisearchEngine) Search(ctx context.Context, indexName string, q Query) (*Result, error) {
	return &Result{}, nil
}

func (e *meilisearchEngine) Close() error {
	return nil
}

// ── Elasticsearch Implementation (Stub) ──────────────────────────────────────

type elasticsearchEngine struct {
	cfg Config
}

func newElasticSearchEngine(cfg Config) (*elasticsearchEngine, error) {
	// TODO: Initialize Elasticsearch client
	return &elasticsearchEngine{cfg: cfg}, nil
}

func (e *elasticsearchEngine) Index(ctx context.Context, indexName string, documentID string, doc any) error {
	return nil
}

func (e *elasticsearchEngine) Update(ctx context.Context, indexName string, documentID string, doc any) error {
	return nil
}

func (e *elasticsearchEngine) Delete(ctx context.Context, indexName string, documentID string) error {
	return nil
}

func (e *elasticsearchEngine) Search(ctx context.Context, indexName string, q Query) (*Result, error) {
	return &Result{}, nil
}

func (e *elasticsearchEngine) Close() error {
	return nil
}

// ── Algolia Implementation (Stub) ────────────────────────────────────────────

type algoliaEngine struct {
	cfg Config
}

func newAlgoliaEngine(cfg Config) (*algoliaEngine, error) {
	// TODO: Initialize Algolia client
	return &algoliaEngine{cfg: cfg}, nil
}

func (e *algoliaEngine) Index(ctx context.Context, indexName string, documentID string, doc any) error {
	return nil
}

func (e *algoliaEngine) Update(ctx context.Context, indexName string, documentID string, doc any) error {
	return nil
}

func (e *algoliaEngine) Delete(ctx context.Context, indexName string, documentID string) error {
	return nil
}

func (e *algoliaEngine) Search(ctx context.Context, indexName string, q Query) (*Result, error) {
	return &Result{}, nil
}

func (e *algoliaEngine) Close() error {
	return nil
}

// ── InMemory Implementation (Stub/Mock) ──────────────────────────────────────

type inmemoryEngine struct {
	cfg Config
}

func newInMemoryEngine(cfg Config) (*inmemoryEngine, error) {
	return &inmemoryEngine{cfg: cfg}, nil
}

func (e *inmemoryEngine) Index(ctx context.Context, indexName string, documentID string, doc any) error {
	return nil
}

func (e *inmemoryEngine) Update(ctx context.Context, indexName string, documentID string, doc any) error {
	return nil
}

func (e *inmemoryEngine) Delete(ctx context.Context, indexName string, documentID string) error {
	return nil
}

func (e *inmemoryEngine) Search(ctx context.Context, indexName string, q Query) (*Result, error) {
	return &Result{}, nil
}

func (e *inmemoryEngine) Close() error {
	return nil
}
