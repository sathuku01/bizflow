package main

import (
	"fmt"
	"log"

	"github.com/yourusername/marketing-consultant-agent/internal/core"
	"github.com/yourusername/marketing-consultant-agent/internal/filters"
)

func main() {
	fmt.Println("=== Marketing Consultant Agent - Foundation Test ===\n")

	// Create a sample business input
	business := core.BusinessInput{
		Type:        core.Retail,
		Description: "Handmade jewelry business selling unique artisan pieces online and at local markets",
		Location:    "Austin, TX",
		Budget:      80.0,
		Channels:    []string{},
		Goal:        core.Awareness,
	}

	// Validate the input
	if err := business.Validate(); err != nil {
		log.Fatalf("Invalid business input: %v", err)
	}

	fmt.Printf("Business Input: %s\n\n", business.String())
	fmt.Printf("Budget Tier: %s\n", business.BudgetTier())
	fmt.Printf("Is Local: %v\n", business.IsLocal())
	fmt.Printf("Is Online Only: %v\n\n", business.IsOnlineOnly())

	// Create platform filter
	platformFilter := filters.NewPlatformFilter()

	// Get filtered platforms
	relevantPlatforms := platformFilter.ApplyAllFilters(business)

	fmt.Println("RELEVANT PLATFORMS:")
	fmt.Println("-------------------")
	for i, platform := range relevantPlatforms {
		metadata, _ := core.GetPlatformMetadata(platform)
		fmt.Printf("%d. %s (Effort: %s, Organic: %v)\n",
			i+1,
			platform,
			metadata.EffortLevel,
			metadata.IsOrganic,
		)
	}

	// Get filtering explanations
	fmt.Println("\nFILTERING REASONING:")
	fmt.Println("--------------------")
	explanations := platformFilter.ExplainFiltering(business)
	for key, explanation := range explanations {
		fmt.Printf("%s: %s\n", key, explanation)
	}

	// Validate constraints for each platform
	constraintValidator := filters.NewConstraintValidator()

	fmt.Println("\nCONSTRAINT ANALYSIS:")
	fmt.Println("--------------------")
	for _, platform := range relevantPlatforms {
		budgetConstraint := constraintValidator.ValidateBudgetConstraints(business.Budget, platform)
		effortConstraint := constraintValidator.ValidateEffortConstraints(business, platform)
		combinedPenalty := constraintValidator.GetCombinedPenalty(business, platform)

		fmt.Printf("\n%s:\n", platform)
		fmt.Printf("  Budget: %s (Penalty: %.2f)\n", budgetConstraint.Reason, budgetConstraint.Penalty)
		fmt.Printf("  Effort: %s (Penalty: %.2f)\n", effortConstraint.Reason, effortConstraint.Penalty)
		fmt.Printf("  Combined Penalty: %.2f\n", combinedPenalty)
	}

	fmt.Println("\n=== Foundation Test Complete ===")
}