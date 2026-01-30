package ai

import (
	"fmt"
	"strings"
)

// MockLLMClient is a fake LLM client for testing and development.
// It returns hardcoded responses based on the content of the prompt.
type MockLLMClient struct{}

// NewMockLLMClient creates a new mock client.
func NewMockLLMClient() *MockLLMClient {
	return &MockLLMClient{}
}

// Generate returns a mock response based on keywords in the prompt.
func (c *MockLLMClient) Generate(prompt string) (string, error) {
	// For Content Template
	if strings.Contains(prompt, "Format the output as a JSON object") {
		return `{ 
			"hook": "✨ This is a mock hook!",
			"caption": "This is a mock caption for your post. Isn't it great?",
			"cta": "Click the link in our bio to learn more!",
			"hashtags": ["#mock", "#testing", "#bizflow"]
		}`, nil
	}

	// For Persona
	if strings.Contains(prompt, "customer persona") {
		return "A mock customer who is very interested in your products.", nil
	}

	// For Reasoning
	if strings.Contains(prompt, "Explain in business terms") {
		return "This is a mock explanation. We recommend this platform because it aligns with your mock business goals.", nil
	}

	// For Risks
	if strings.Contains(prompt, "potential risks") {
		return "1. Mock Risk One\n2. Mock Risk Two", nil
	}

	// For Strategy
	if strings.Contains(prompt, "strategic next step") {
		return "A mock strategic step would be to engage with your mock audience.", nil
	}

	// Default fallback
	return fmt.Sprintf("Mock response for prompt: %s", prompt), nil
}
