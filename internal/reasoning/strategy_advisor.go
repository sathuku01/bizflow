package reasoning

import (
	"biz-flow/internal/ai"
	"biz-flow/internal/core"
	"fmt"
)

// GenerateStrategy uses an AI client to create a short, actionable strategic recommendation
// for the business's top-ranked platform.
func GenerateStrategy(client ai.LLMClient, business core.BusinessInput, topPlatform core.Platform) (string, error) {
	prompt := fmt.Sprintf(
		ai.StrategicAdvicePrompt,
		business.Type,
		topPlatform,
	)

	strategy, err := client.Generate(prompt)
	if err != nil {
		return "", fmt.Errorf("failed to generate strategy: %w", err)
	}

	return strategy, nil
}