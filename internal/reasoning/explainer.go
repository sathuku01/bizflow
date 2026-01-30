package reasoning

import (
	"biz-flow/internal/ai"
	"biz-flow/internal/core"
	"fmt"
)

// ExplainRecommendations uses an AI client to generate business-friendly explanations for
// why each of the recommended platforms was chosen.
func ExplainRecommendations(client ai.LLMClient, business core.BusinessInput, platforms []core.Platform) (map[core.Platform]string, error) {
	explanations := make(map[core.Platform]string)

	for _, platform := range platforms {
		prompt := fmt.Sprintf(
			ai.ReasoningPrompt,
			platform,
			business.Type,
			business.Goal,
		)

		reasoning, err := client.Generate(prompt)
		if err != nil {
			// In a real app, you might want to return a default message and log the error
			return nil, fmt.Errorf("failed to generate reasoning for %s: %w", platform, err)
		}
		explanations[platform] = reasoning
	}

	return explanations, nil
}
