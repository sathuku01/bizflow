# Integration Expectations

This document outlines the functions and interfaces this module (Scoring & Reasoning) expects from other modules.

## From Developer A (Core Engine)

### `internal/core/platform.go`

To properly score platforms, the scoring engine needs a reliable way to fetch the metadata for any given platform. I expect the following function to be available in the `core` package:

```go
// GetPlatformMetadata retrieves the metadata for a specific platform.
// It returns the metadata and a boolean which is false if the platform
// is not found in the platform list.
func GetPlatformMetadata(platform Platform) (PlatformMetadata, bool)
```

This function is critical for the `audience_scorer` and `return_scorer`.

### `internal/core/business.go`

The web handler needs to validate the incoming business input. I expect the `BusinessInput` struct to have a `Validate` method.

```go
package core

// Validate checks if the business input contains valid data.
func (b *BusinessInput) Validate() error
```

---

## From Developer B (AI Integration)

The Reasoning Engine relies entirely on the AI module. I expect the following components to be available in the `internal/ai` package.

### `internal/ai/client.go`

A generic interface to interact with an LLM.

```go
package ai

// LLMClient defines the interface for an AI client that can generate text.
type LLMClient interface {
	Generate(prompt string) (string, error)
}
```

### `internal/ai/prompts.go`

Exported string constants for all required prompts. The reasoning module will use `fmt.Sprintf` to insert values into these prompts.

```go
package ai

const (
	// PersonaInferencePrompt is used to infer the customer persona.
	// Parameter: (1: Business Description)
	PersonaInferencePrompt = "Based on this business description: '%s', describe the ideal customer persona in one sentence."

	// ContentTemplatePrompt is used to generate a content template for a specific platform.
	// Parameters: (1: Business Type), (2: Platform Name)
	// The output MUST be a JSON object that can be unmarshalled into the ContentTemplate struct.
	ContentTemplatePrompt = "For a '%s' business on '%s', generate a content template with a hook, caption, call-to-action, and three relevant hashtags. Format the output as a JSON object: {\"hook\": \"...\", \"caption\": \"...\", \"cta\": \"...\", \"hashtags\": [\"...\", \"...\"]}"

	// ReasoningPrompt is used to generate a business-friendly explanation for a recommendation.
	// Parameters: (1: Platform Name), (2: Business Type), (3: Business Goal)
	ReasoningPrompt = "Explain in business terms why %s is a good marketing platform for a %s business with a goal of %s. Do not mention scores."

	// RiskAssessmentPrompt is used to identify potential risks.
	// Parameter: (1: Platform Name)
	RiskAssessmentPrompt = "What are two potential risks for a small business using %s for marketing? List each on a new line."

	// StrategicAdvicePrompt is used to generate a single, actionable next step.
	// Parameters: (1: Business Type), (2: Platform Name)
	StrategicAdvicePrompt = "Given a recommendation for a %s business to use %s, provide a short, actionable strategic next step."
)
```

### Other AI Functions

The orchestrator will call these functions, which are expected to use the client and prompts above.

```go
package persona_inferrer

// InferPersona infers the customer persona from the business description.
func InferPersona(client ai.LLMClient, business core.BusinessInput) (string, error)
```

```go
package content_generator

// GenerateContentTemplate creates a content template for the top-ranked platform.
func GenerateContentTemplate(client ai.LLMClient, business core.BusinessInput, platform core.Platform) (*core.ContentTemplate, error)
```
