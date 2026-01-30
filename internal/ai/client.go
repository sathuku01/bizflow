package ai

import (
	"context"
	"fmt"
	"os"
	"sort"

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

	model := client.GenerativeModel("gemini-3-flash-preview") // works with generateContent

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
func RunAgent(input BusinessInput) AgentOutput {
	session := AgentSession{Input: input}

	// Step 1: Filter platforms deterministically
	session.Platforms = FilterPlatforms(input)

	client, err := NewLLMClient()
	if err != nil {
		fmt.Println("LLM Error:", err)
		return ErrorOutput("AI backend unavailable")
	}

	// Step 2: Generate reasoning for each platform
	for i, pf := range session.Platforms {
		prompt := fmt.Sprintf(
			"Business: %s\nDescription: %s\nBudget: %.2f\nGoal: %s\nExplain why the platform %s is suitable.",
			input.BusinessType, input.Description, input.MonthlyBudget, input.PrimaryGoal, pf.Name,
		)
		text, err := client.GenerateText(prompt)
		if err != nil {
			text = "Fallback reasoning (LLM failed)."
		}
		session.Platforms[i].Reason = text
		// Optional: assign a mock score based on deterministic rules
		session.Platforms[i].Score = mockScore(input, pf.Name)
	}

	// Step 3: Rank platforms
	sort.Slice(session.Platforms, func(i, j int) bool {
		return session.Platforms[i].Score > session.Platforms[j].Score
	})

	// Step 4: Generate content template for top platform
	if len(session.Platforms) > 0 {
		top := session.Platforms[0]
		contentPrompt := fmt.Sprintf(
			"Generate an example social media content template for a %s business using %s platform. Include hook, caption, CTA, hashtags.",
			input.BusinessType, top.Name,
		)
		contentText, err := client.GenerateText(contentPrompt)
		if err != nil {
			contentText = "Fallback content template"
		}
		session.TopContent = ContentTemplate{
			Hook:     "Generated automatically",
			Caption:  contentText,
			CTA:      "CTA coming soon",
			Hashtags: []string{"#agent"},
		}
	}

	// Step 5: Generate strategic advice & risks
	advPrompt := fmt.Sprintf(
		"Given the ranked platforms for a %s business, provide concise strategic advice and potential risks.",
		input.BusinessType,
	)
	advText, err := client.GenerateText(advPrompt)
	if err != nil {
		advText = "Fallback strategic advice"
	}
	session.Advice = advText
	session.Risks = []string{"Next step: add tools"}

	// Step 6: Build output
	var recommendations []Recommendation
	for idx, pf := range session.Platforms {
		recommendations = append(recommendations, Recommendation{
			Rank:            idx + 1,
			Platform:        pf.Name,
			Reasoning:       pf.Reason,
			ContentTemplate: &session.TopContent,
		})
	}

	return AgentOutput{
		Recommendations: recommendations,
		StrategicAdvice: session.Advice,
		Risks:           session.Risks,
	}
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
