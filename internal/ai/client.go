package ai

import (
	"context"
	"fmt"
	"os"
	"log"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// --------------------
// Session & Structures
// --------------------
type AgentSession struct {
	Input       BusinessInput
	Platforms   []PlatformScore
	TopContent  ContentTemplate
	Advice      string
	Risks       []string
}

type PlatformScore struct {
	Name     string
	Score    float64
	Reason   string
}

// --------------------
// LLM Client
// --------------------
type LLMClient struct {
	model *genai.GenerativeModel
}

var ErrMissingAPIKey = fmt.Errorf("GOOGLE_API_KEY not set")

func NewLLMClient() (*LLMClient, error) {
	apiKey := os.Getenv("GOOGLE_API_KEY")
	if apiKey == "" {
		return nil, ErrMissingAPIKey
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to init Gemini client: %v", err)
	}

	model := client.GenerativeModel("gemini-2.0-flash") // works with generateContent

	return &LLMClient{model: model}, nil
}

func (l *LLMClient) GenerateText(prompt string) (string, error) {
	ctx := context.Background()
	resp, err := l.model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", fmt.Errorf("gemini API error: %v", err)
	}

	if len(resp.Candidates) == 0 {
		return "", fmt.Errorf("no candidates returned by LLM")
	}

	parts := resp.Candidates[0].Content.Parts
	if len(parts) == 0 {
		return "", fmt.Errorf("candidate has no content parts")
	}

	result := ""
	for _, p := range parts {
		result += fmt.Sprintf("%v", p)
	}

	return result, nil
}

// --------------------
// Core Agent Logic
// --------------------
func RunAgent(input BusinessInput, notion *NotionClient) AgentOutput {
	llm, err := NewLLMClient()
	if err != nil {
		log.Println("LLM unavailable:", err)
		return ErrorOutput("LLM unavailable, using fallback reasoning")
	}

	platforms := []string{"Instagram", "Facebook", "TikTok", "Google My Business"}
	var recommendations []Recommendation

	for _, p := range platforms {
		// Fetch template from Notion
		hook, caption, cta, hashtags, err := notion.FetchTemplate(p)
		if err != nil {
			log.Println("Failed to fetch template for", p, ":", err)
			hook, caption, cta, hashtags = "Generated automatically", "Fallback content template", "CTA coming soon", []string{"#agent"}
		}

		// Build prompt for LLM
		prompt := fmt.Sprintf(
			"You are a marketing strategist. Your task: generate a short reasoning why the platform '%s' is suitable for a business with description: '%s'. Use this content template to guide the messaging:\nHook: %s\nCaption: %s\nCTA: %s\nHashtags: %v\nOutput should be concise, persuasive, and business-oriented.",
			p, input.Description, hook, caption, cta, hashtags,
		)

		// Generate reasoning from LLM
		reasoning, err := llm.GenerateText(prompt)
		if err != nil {
			log.Println("LLM generation failed for", p, ":", err)
			reasoning = fmt.Sprintf("Fallback reasoning for %s based on business: %s", p, input.Description)
		}

		rec := Recommendation{
			Rank:      len(recommendations) + 1,
			Platform:  p,
			Reasoning: reasoning,
			ContentTemplate: &ContentTemplate{
				Hook:     hook,
				Caption:  caption,
				CTA:      cta,
				Hashtags: hashtags,
			},
		}

		recommendations = append(recommendations, rec)
	}

	output := AgentOutput{
		Recommendations: recommendations,
		StrategicAdvice: "Strategic advice placeholder",
		Risks:           []string{"Next step: add tools"},
	}

	// Save query to history DB
	if err := notion.SaveQuery(input, output); err != nil {
		log.Println("Warning: failed to save query to Notion:", err)
	}

	return output
}

// --------------------
// Deterministic Helpers
// --------------------
func FilterPlatforms(input BusinessInput) []PlatformScore {
	var platforms []string

	switch input.BusinessType {
	case "retail":
		platforms = []string{"Instagram", "Facebook", "TikTok", "Google My Business"}
	case "service":
		platforms = []string{"Google My Business", "Facebook", "WhatsApp Business", "Instagram"}
	case "digital":
		platforms = []string{"LinkedIn", "Email", "YouTube", "Instagram"}
	}

	var scores []PlatformScore
	for _, p := range platforms {
		scores = append(scores, PlatformScore{Name: p})
	}
	return scores
}

// Example scoring function (mock for now)
func mockScore(input BusinessInput, platform string) float64 {
	score := 5.0
	if platform == "Instagram" && input.BusinessType == "retail" {
		score += 3
	}
	if input.MonthlyBudget > 200 && (platform == "Facebook" || platform == "Google My Business") {
		score += 2
	}
	return score
}

// --------------------
// Fallback Output
// --------------------
func ErrorOutput(msg string) AgentOutput {
	return AgentOutput{
		Recommendations: []Recommendation{},
		StrategicAdvice: msg,
		Risks:           []string{"AI unavailable"},
	}
}
