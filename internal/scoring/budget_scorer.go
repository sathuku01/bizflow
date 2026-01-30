package scoring

import (
	"biz-flow/internal/core"
	"biz-flow/internal/filters"
)

// BudgetScorer scores a platform based on its budget feasibility.
type BudgetScorer struct {
	validator *filters.ConstraintValidator
}

// NewBudgetScorer creates a new BudgetScorer.
func NewBudgetScorer() *BudgetScorer {
	return &BudgetScorer{
		validator: filters.NewConstraintValidator(),
	}
}

// Score calculates a score based on budget constraints.
// The score is calculated as 1.0 minus the penalty assessed by the ConstraintValidator.
// A lower penalty results in a higher score.
func (s *BudgetScorer) Score(business core.BusinessInput, platform core.Platform) float64 {
	constraint := s.validator.ValidateBudgetConstraints(business.Budget, platform)

	score := 1.0 - constraint.Penalty
	if score < 0 {
		return 0.0
	}
	return score
}