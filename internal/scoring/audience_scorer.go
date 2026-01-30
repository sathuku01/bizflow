package scoring

import "biz-flow/internal/core"

// AudienceScorer scores a platform based on its alignment with the business's type.
type AudienceScorer struct{}

// NewAudienceScorer creates a new AudienceScorer.
func NewAudienceScorer() *AudienceScorer {
	return &AudienceScorer{}
}

// Score assesses if the platform is a good fit for the business's industry type.
// It returns 1.0 if the business type is listed in the platform's 'BestFor' list,
// and 0.1 otherwise.
func (s *AudienceScorer) Score(business core.BusinessInput, platform core.Platform) float64 {
	metadata, exists := core.GetPlatformMetadata(platform)
	if !exists {
		return 0.0 // Platform not found, no score.
	}

	for _, businessType := range metadata.BestFor {
		if businessType == business.Type {
			return 1.0 // Perfect fit for this business type.
		}
	}

	return 0.1 // Poor fit for this business type.
}
