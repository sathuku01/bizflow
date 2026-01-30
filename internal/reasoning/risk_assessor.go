package reasoning

import (
	"biz-flow/internal/ai"
	"biz-flow/internal/core"
	"fmt"
	"strings"
)

// AssessRisks uses an AI client to identify potential risks associated with using a
// specific marketing platform.
func AssessRisks(client ai.LLMClient, platform core.Platform) ([]string, error) {
	prompt := fmt.Sprintf(ai.RiskAssessmentPrompt, platform)

	rawResponse, err := client.Generate(prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to generate risks for %s: %w", platform, err)
	}

	// The prompt asks for each risk on a new line. We parse the response here.
	risks := strings.Split(rawResponse, "\n")
	var cleanedRisks []string
	for _, r := range risks {
		// Trim whitespace and common list prefixes
		rimmed := strings.TrimSpace(strings.TrimPrefix(strings.TrimPrefix(r, "*"), "-"))
		if trimmed != "" {
			cleanedRisks = append(cleanedRisks, trimmed)
		}
	}

	return cleanedRisks, nil
}