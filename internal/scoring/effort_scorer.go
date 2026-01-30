package scoring

import (
	"biz-flow/internal/core"
	"biz-flow/internal/filters"
)

// EffortScorer scores a platform based on the estimated effort required from the user.
type EffortScorer struct {
	validator *filters.ConstraintValidator
}

// NewEffortScorer creates a new EffortScorer.
func NewEffortScorer() *EffortScorer {
	return &EffortScorer{
		validator: filters.NewConstraintValidator(),
	}
}

// Score calculates a score based on effort constraints.
// The score is calculated as 1.0 minus the penalty assessed by the ConstraintValidator.
// A lower penalty (less effort) results in a higher score.
func (s *EffortScorer) Score(business core.BusinessInput, platform core.Platform) float64 {
	constraint := s.validator.ValidateEffortConstraints(business, platform)

	score := 1.0 - constraint.Penalty
	if score < 0 {
		return 0.0
	}
	return score
}