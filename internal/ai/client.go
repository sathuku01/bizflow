package ai

import (
	"context"
	"fmt"
	"os"

	"github.com/sashabaranov/go-openai"
)

type LLMClient struct {
	client *openai.Client
	model  string
}

func NewLLMClient() (*LLMClient, error) {
	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("OPENROUTER_API_KEY missing")
	}
	config := openai.DefaultConfig(apiKey)
	config.BaseURL = "https://openrouter.ai/api/v1"

	return &LLMClient{
		client: openai.NewClientWithConfig(config),
		model:  "google/gemini-2.0-flash-001",
	}, nil
}

func (l *LLMClient) GenerateText(prompt string) (string, error) {
	resp, err := l.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: l.model,
			Messages: []openai.ChatCompletionMessage{
				{Role: openai.ChatMessageRoleUser, Content: prompt},
			},
		},
	)
	if err != nil {
		return "", err
	}
	return resp.Choices[0].Message.Content, nil
}

func RunAgent(input BusinessInput, notion *NotionClient) AgentOutput {
	llm, err := NewLLMClient()
	if err != nil {
		return AgentOutput{StrategicAdvice: "LLM Init Failed"}
	}

	platforms := []string{"Instagram", "Facebook", "TikTok", "Google My Business"}
	var recs []Recommendation

	for _, p := range platforms {
		// 1. Fetch template from Notion
		hook, caption, cta, tags, _ := notion.FetchTemplate(p)

		// 2. Build the prompt
		prompt := fmt.Sprintf(`You are a micro-business marketing consultant.
Analyze %s for a %s business: %s.
Goal: %s | Budget: $%.2f

Instructions:
1. Provide 2-3 concise paragraphs of strategic reasoning.
2. IMPORTANT: If the following template fields are empty, invent creative ones.
   Hook: %s
   Caption: %s
   CTA: %s

Return your response in this format:
REASONING: [Your paragraphs]
NEW_HOOK: [Invented if missing]
NEW_CAPTION: [Invented if missing]
NEW_CTA: [Invented if missing]`,
			p, input.BusinessType, input.Description, input.PrimaryGoal, input.MonthlyBudget, hook, caption, cta)

		// 3. Get AI Response
		aiRawResponse, _ := llm.GenerateText(prompt)

		// 4. (Optional) Basic parsing logic could go here to extract the fields. 
		// For now, we'll keep the reasoning as the full AI text.
		recs = append(recs, Recommendation{
			Rank:      len(recs) + 1,
			Platform:  p,
			Reasoning: aiRawResponse, 
			ContentTemplate: &ContentTemplate{
				Hook: hook, Caption: caption, CTA: cta, Hashtags: tags,
			},
		})
	}

	output := AgentOutput{
		Recommendations: recs,
		StrategicAdvice: "Focus on visual consistency and leveraging local search intent.",
		Risks:           []string{"High competition in retail category"},
	}
	errSave := notion.SaveQuery(input, output)
	if errSave != nil {
    	fmt.Printf("❌ NOTION ERROR: %v\n", errSave)
	}
	_ = notion.SaveQuery(input, output)
	return output
}