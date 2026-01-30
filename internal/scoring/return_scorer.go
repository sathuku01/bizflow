package scoring

import "biz-flow/internal/core"

// ReturnScorer scores a platform based on its potential to meet a specific marketing goal.
type ReturnScorer struct{}

// NewReturnScorer creates a new ReturnScorer.
func NewReturnScorer() *ReturnScorer {
	return &ReturnScorer{}
}

// Score assesses a platform's potential return based on the business's primary goal.
// For 'Awareness', score is based on ReachPotential.
// For 'Sales', score is based on ConversionFocus.
func (s *ReturnScorer) Score(business core.BusinessInput, platform core.Platform) float64 {
	metadata, exists := core.GetPlatformMetadata(platform)
	if !exists {
		return 0.0 // Platform not found, no score.
	}

	switch business.Goal {
	case core.Awareness:
		// Score is based on the platform's reach potential (1-10 scale).
		return float64(metadata.ReachPotential) / 10.0
	case core.Sales:
		// Score is based on the platform's conversion focus (1-10 scale).
		return float64(metadata.ConversionFocus) / 10.0
	default:
		return 0.0 // Unknown goal, no score.
	}
}