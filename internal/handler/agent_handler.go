package handler

import (
	"biz-flow/internal/ai"
	"biz-flow/internal/core"
	"biz-flow/internal/filters"
	"biz-flow/internal/reasoning"
	"biz-flow/internal/scoring"
	"sort"
)

// The following imports are for expected packages from Dev B
import (
	"biz-flow/internal/ai/content_generator"
	"biz-flow/internal/ai/persona_inferrer"
)

type scoredPlatform struct {
	Platform core.Platform
	Score    float64
}

// GenerateConsultation is the central orchestrator of the BizFlow agent.
// It takes a business input and an AI client, and runs the full pipeline to
// generate a complete marketing consultation.
func GenerateConsultation(businessInput core.BusinessInput, client ai.LLMClient) (*core.ConsultationResult, error) {
	// 1. Instantiate services
	platformFilter := filters.NewPlatformFilter()
	scorers := []scoring.Scorer{
		scoring.NewAudienceScorer(),
		scoring.NewBudgetScorer(),
		scoring.NewEffortScorer(),
		scoring.NewReturnScorer(),
	}

	// 2. Filter relevant platforms
	relevantPlatforms := platformFilter.ApplyAllFilters(businessInput)
	if len(relevantPlatforms) == 0 {
		return &core.ConsultationResult{
			StrategicAdvice: "No suitable marketing platforms were found based on your business profile. This might be due to very specific constraints.",
		}, nil
	}

	// 3. Score and Rank
	rankedPlatforms := make([]scoredPlatform, 0, len(relevantPlatforms))
	for _, platform := range relevantPlatforms {
		var compositeScore float64
		for _, s := range scorers {
			compositeScore += s.Score(businessInput, platform)
		}
		avgScore := compositeScore / float64(len(scorers))
		rankedPlatforms = append(rankedPlatforms, scoredPlatform{Platform: platform, Score: avgScore})
	}

	sort.Slice(rankedPlatforms, func(i, j int) bool {
		return rankedPlatforms[i].Score > rankedPlatforms[j].Score
	})

	// Limit to top 3
	topPlatforms := rankedPlatforms
	if len(topPlatforms) > 3 {
		topPlatforms = topPlatforms[:3]
	}

	// 4. Generate AI Content
	topPlatformNames := make([]core.Platform, len(topPlatforms))
	for i, p := range topPlatforms {
		topPlatformNames[i] = p.Platform
	}

	persona, err := persona_inferrer.InferPersona(client, businessInput)
	if err != nil {
		return nil, err
	}

	explanations, err := reasoning.ExplainRecommendations(client, businessInput, topPlatformNames)
	if err != nil {
		return nil, err
	}

	contentTemplate, err := content_generator.GenerateContentTemplate(client, businessInput, topPlatforms[0].Platform)
	if err != nil {
		return nil, err
	}

	risks, err := reasoning.AssessRisks(client, topPlatforms[0].Platform)
	if err != nil {
		return nil, err
	}

	strategy, err := reasoning.GenerateStrategy(client, businessInput, topPlatforms[0].Platform)
	if err != nil {
		return nil, err
	}

	// 5. Assemble Final Result
	result := &core.ConsultationResult{
		Persona:         persona,
		Risks:           risks,
		StrategicAdvice: strategy,
		Recommendations: make([]core.Recommendation, len(topPlatforms)),
	}

	for i, p := range topPlatforms {
		rec := core.Recommendation{
			Rank:      i + 1,
			Platform:  p.Platform,
			Score:     p.Score,
			Reasoning: explanations[p.Platform],
		}
		if i == 0 {
			rec.ContentTemplate = contentTemplate
		}
		result.Recommendations[i] = rec
	}

	return result, nil
}
