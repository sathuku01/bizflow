package filters

import (
	"fmt"
	"github.com/yourusername/marketing-consultant-agent/internal/core"
)

// ConstraintValidator validates specific constraints for platforms
type ConstraintValidator struct{}

// NewConstraintValidator creates a new constraint validator
func NewConstraintValidator() *ConstraintValidator {
	return &ConstraintValidator{}
}

// BudgetConstraint represents the result of a budget validation
type BudgetConstraint struct {
	IsValid     bool
	Reason      string
	Penalty     float64 // 0.0 (no penalty) to 1.0 (heavy penalty)
}

// EffortConstraint represents the result of an effort validation
type EffortConstraint struct {
	IsValid     bool
	Reason      string
	Penalty     float64 // 0.0 (no penalty) to 1.0 (heavy penalty)
}

// ValidateBudgetConstraints checks if a platform is feasible given the budget
func (cv *ConstraintValidator) ValidateBudgetConstraints(
	budget float64,
	platform core.Platform,
) BudgetConstraint {
	metadata, exists := core.GetPlatformMetadata(platform)
	if !exists {
		return BudgetConstraint{
			IsValid: false,
			Reason:  "Platform metadata not found",
			Penalty: 1.0,
		}
	}

	// Check if budget meets minimum requirement
	if budget < metadata.MinBudget {
		return BudgetConstraint{
			IsValid: false,
			Reason: fmt.Sprintf(
				"Budget ($%.2f/month) is below minimum required ($%.2f/month) for %s",
				budget,
				metadata.MinBudget,
				platform,
			),
			Penalty: 1.0,
		}
	}

	// Apply penalties based on budget tier
	if budget < 50 {
		// Low budget (<$50): Organic platforms are preferred
		if !metadata.IsOrganic {
			return BudgetConstraint{
				IsValid: true,
				Reason:  "Low budget makes paid platforms less effective",
				Penalty: 0.7, // Heavy penalty for paid platforms
			}
		}
		return BudgetConstraint{
			IsValid: true,
			Reason:  "Perfect fit for low-budget organic marketing",
			Penalty: 0.0,
		}
	}

	if budget >= 50 && budget <= 200 {
		// Medium budget ($50-$200): Organic preferred, paid soft-penalized
		if metadata.IsPaid && !metadata.IsOrganic {
			return BudgetConstraint{
				IsValid: true,
				Reason:  "Medium budget can support limited paid advertising",
				Penalty: 0.3, // Soft penalty for paid-only platforms
			}
		}
		return BudgetConstraint{
			IsValid: true,
			Reason:  "Good budget for consistent organic presence",
			Penalty: 0.0,
		}
	}

	// High budget (>$200): All platforms viable
	return BudgetConstraint{
		IsValid: true,
		Reason:  "Budget supports both organic and paid strategies",
		Penalty: 0.0,
	}
}

// ValidateEffortConstraints checks if a platform is feasible given effort capacity
func (cv *ConstraintValidator) ValidateEffortConstraints(
	business core.BusinessInput,
	platform core.Platform,
) EffortConstraint {
	metadata, exists := core.GetPlatformMetadata(platform)
	if !exists {
		return EffortConstraint{
			IsValid: false,
			Reason:  "Platform metadata not found",
			Penalty: 1.0,
		}
	}

	// Video platforms are high-effort for micro-businesses
	if metadata.RequiresVideo {
		// Exception: Retail businesses with visual products can handle video
		if business.Type == core.Retail {
			return EffortConstraint{
				IsValid: true,
				Reason:  "Visual products are well-suited for video content",
				Penalty: 0.2, // Small penalty (video still takes effort)
			}
		}

		// Service and Digital businesses: video is harder
		return EffortConstraint{
			IsValid: true,
			Reason:  "Video content requires significant production effort for service/digital businesses",
			Penalty: 0.6, // Moderate-to-heavy penalty
		}
	}

	// Effort level penalties
	switch metadata.EffortLevel {
	case core.HighEffort:
		return EffortConstraint{
			IsValid: true,
			Reason:  "High-effort platform may strain micro-business resources",
			Penalty: 0.4,
		}
	case core.MediumEffort:
		return EffortConstraint{
			IsValid: true,
			Reason:  "Moderate effort required, manageable for consistent posting",
			Penalty: 0.1,
		}
	case core.LowEffort:
		return EffortConstraint{
			IsValid: true,
			Reason:  "Low-effort platform, ideal for resource-constrained businesses",
			Penalty: 0.0,
		}
	default:
		return EffortConstraint{
			IsValid: true,
			Reason:  "Unknown effort level",
			Penalty: 0.0,
		}
	}
}

// ValidateVisualRequirements checks if a business can meet visual content needs
func (cv *ConstraintValidator) ValidateVisualRequirements(
	business core.BusinessInput,
	platform core.Platform,
) EffortConstraint {
	metadata, exists := core.GetPlatformMetadata(platform)
	if !exists {
		return EffortConstraint{
			IsValid: false,
			Reason:  "Platform metadata not found",
			Penalty: 1.0,
		}
	}

	if !metadata.RequiresVisuals {
		// Platform doesn't require visuals, always valid
		return EffortConstraint{
			IsValid: true,
			Reason:  "Platform works well with text-based content",
			Penalty: 0.0,
		}
	}

	// Visual platforms
	switch business.Type {
	case core.Retail:
		// Retail products are inherently visual
		return EffortConstraint{
			IsValid: true,
			Reason:  "Retail products provide natural visual content opportunities",
			Penalty: 0.0,
		}
	case core.Service:
		// Services can show before/after, team photos, etc.
		return EffortConstraint{
			IsValid: true,
			Reason:  "Services can create visual content (before/after, testimonials, team)",
			Penalty: 0.2,
		}
	case core.Digital:
		// Digital products/services may struggle with visuals
		return EffortConstraint{
			IsValid: true,
			Reason:  "Digital products may require creative approaches to visual content",
			Penalty: 0.3,
		}
	default:
		return EffortConstraint{
			IsValid: true,
			Reason:  "Unknown business type",
			Penalty: 0.0,
		}
	}
}

// ValidateGoalAlignment checks if a platform aligns with marketing goal
func (cv *ConstraintValidator) ValidateGoalAlignment(
	goal core.MarketingGoal,
	platform core.Platform,
) BudgetConstraint {
	metadata, exists := core.GetPlatformMetadata(platform)
	if !exists {
		return BudgetConstraint{
			IsValid: false,
			Reason:  "Platform metadata not found",
			Penalty: 1.0,
		}
	}

	switch goal {
	case core.Awareness:
		// Awareness benefits from high reach potential
		if metadata.ReachPotential >= 8 {
			return BudgetConstraint{
				IsValid: true,
				Reason:  "Excellent reach potential for awareness campaigns",
				Penalty: 0.0,
			}
		} else if metadata.ReachPotential >= 6 {
			return BudgetConstraint{
				IsValid: true,
				Reason:  "Moderate reach potential for awareness",
				Penalty: 0.2,
			}
		}
		return BudgetConstraint{
			IsValid: true,
			Reason:  "Limited reach potential for awareness goals",
			Penalty: 0.4,
		}

	case core.Sales:
		// Sales benefits from high conversion focus
		if metadata.ConversionFocus >= 8 {
			return BudgetConstraint{
				IsValid: true,
				Reason:  "Excellent conversion potential for sales goals",
				Penalty: 0.0,
			}
		} else if metadata.ConversionFocus >= 6 {
			return BudgetConstraint{
				IsValid: true,
				Reason:  "Moderate conversion potential for sales",
				Penalty: 0.2,
			}
		}
		return BudgetConstraint{
			IsValid: true,
			Reason:  "Limited conversion potential for direct sales",
			Penalty: 0.4,
		}

	default:
		return BudgetConstraint{
			IsValid: true,
			Reason:  "Unknown goal",
			Penalty: 0.0,
		}
	}
}

// GetCombinedPenalty calculates the combined penalty from all constraints
func (cv *ConstraintValidator) GetCombinedPenalty(
	business core.BusinessInput,
	platform core.Platform,
) float64 {
	budgetConstraint := cv.ValidateBudgetConstraints(business.Budget, platform)
	effortConstraint := cv.ValidateEffortConstraints(business, platform)
	visualConstraint := cv.ValidateVisualRequirements(business, platform)
	goalConstraint := cv.ValidateGoalAlignment(business.Goal, platform)

	// Average the penalties (weighted equally for now)
	totalPenalty := budgetConstraint.Penalty +
		effortConstraint.Penalty +
		visualConstraint.Penalty +
		goalConstraint.Penalty

	return totalPenalty / 4.0
}

// IsValidPlatform checks if a platform is valid (no hard constraints violated)
func (cv *ConstraintValidator) IsValidPlatform(
	business core.BusinessInput,
	platform core.Platform,
) bool {
	budgetConstraint := cv.ValidateBudgetConstraints(business.Budget, platform)
	effortConstraint := cv.ValidateEffortConstraints(business, platform)
	visualConstraint := cv.ValidateVisualRequirements(business, platform)

	return budgetConstraint.IsValid &&
		effortConstraint.IsValid &&
		visualConstraint.IsValid
}