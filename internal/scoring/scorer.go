package scoring

import "biz-flow/internal/core"

// Scorer defines the interface for any component that can score a marketing platform
// based on a business input. The returned score should be normalized between 0.0 (worst fit)
// and 1.0 (perfect fit).
type Scorer interface {
	Score(business core.BusinessInput, platform core.Platform) float64
}