package es

// Index template
type ElasticIndex struct {
	Settings ElasticSettings          `json:"settings,omitempty"`
	Aliases  map[string]ElasticAlias  `json:"aliases,omitempty"`
	Mappings ElasticMappingProperties `json:"mappings,omitempty"`
}

// ElasticSettings
type ElasticSettings struct {
	NumberOfShards int                     `json:"number_of_shards,omitempty"`
	Analysis       ElasticSettingsAnalysis `json:"analysis,omitempty"`
}

// settings.analysis
type ElasticSettingsAnalysis struct {
	Analyzer   map[string]ElasticAnalyzer   `json:"analyzer,omitempty"`
	Normalizer map[string]ElasticNormalizer `json:"normalizer,omitempty"`
}

// ElasticAnalyzer is an alias entry
type ElasticAnalyzer struct {
	CharFilter []string `json:"char_filter,omitempty"`
	Tokenizer  string   `json:"tokenizer,omitempty"`
	Filter     []string `json:"filter,omitempty"`
}

type ElasticNormalizer struct {
	Type   string   `json:"type,omitempty"`
	Filter []string `json:"filter,omitempty"`
}

type ElasticAlias struct {
	Filter  map[string]map[string]map[string]map[string]string `json:"filter,omitempty"`
	Routing int                                                `json:"routing,omitempty"`
}
type ElasticMappingProperties struct {
	Properties map[string]ElasticMappingProperty `json:"properties,omitempty"`
}

type ElasticMappingProperty struct {
	Type       string            `json:"type,omitempty"`
	Normalizer string            `json:"normalizer,omitempty"`
	Meta       map[string]string `json:"meta,omitempty"`
}
