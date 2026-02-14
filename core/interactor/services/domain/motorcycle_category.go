package domain

// MotorcycleCategory represents a distinct motorcycle category with its line count (HU41)
type MotorcycleCategory struct {
	Name      string `json:"name"`
	LineCount int    `json:"line_count"`
}

// CategoryLine represents a motorcycle line (model) within a specific category (HU41)
type CategoryLine struct {
	Model              string `json:"model"`
	BrandName          string `json:"brand"`
	EngineDisplacement int    `json:"engine_displacement"`
}

// EngineDisplacementRange represents a displacement range category (HU49)
type EngineDisplacementRange struct {
	Range string `json:"range"`
}

// RatingRange represents a valid rating value with its label (HU48)
type RatingRange struct {
	Value int    `json:"value"`
	Label string `json:"label"`
}
