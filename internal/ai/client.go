package ai

import (
	"context"
	"fmt"
	"os"
	"strings"

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

func extractValue(text, marker string) string {
	parts := strings.Split(text, marker)
	if len(parts) < 2 {
		return ""
	}
	result := strings.Split(parts[1], "\n")[0]
	return strings.TrimSpace(result)
}

func RunAgent(input BusinessInput, notion *NotionClient) AgentOutput {
    llm, err := NewLLMClient()
    if err != nil {
        return AgentOutput{StrategicAdvice: "LLM Init Failed"}
    }

    platforms := []string{"Instagram", "Facebook", "TikTok", "Google My Business"}
    var recs []Recommendation

    for _, p := range platforms {
        // 1. Check if template exists
        hook, caption, cta, tags, _ := notion.FetchTemplate(p)
        wasEmpty := (hook == "")

        prompt := fmt.Sprintf(`You are a micro-business marketing consultant.
Analyze %s for a %s business: %s.
Goal: %s | Budget: $%.2f`,
            p, input.BusinessType, input.Description, input.PrimaryGoal, input.MonthlyBudget)

        if wasEmpty {
            prompt += `\nIMPORTANT: No template found. Invent a creative Hook, Caption, and CTA.`
        } else {
            prompt += fmt.Sprintf(`\nUse this template context:\nHook: %s\nCaption: %s\nCTA: %s`, hook, caption, cta)
        }

        prompt += `\nReturn your response in this format:
REASONING: [Your 2-3 paragraphs]
NEW_HOOK: [The hook to use]
NEW_CAPTION: [The caption to use]
NEW_CTA: [The CTA to use]
HASHTAGS: [space separated tags without #]` // Added this to get the tags from AI

        aiRawResponse, err := llm.GenerateText(prompt)
        if err != nil {
            fmt.Printf("❌ AI GENERATION ERROR for %s: %v\n", p, err)
            continue
        }

        // 2. Extract AI generated content
        fHook := hook
        if fHook == "" { fHook = extractValue(aiRawResponse, "NEW_HOOK:") }
        fCaption := caption
        if fCaption == "" { fCaption = extractValue(aiRawResponse, "NEW_CAPTION:") }
        fCTA := cta
        if fCTA == "" { fCTA = extractValue(aiRawResponse, "NEW_CTA:") }
        
        // Extract tags from AI to avoid sending an empty slice to the DB
        rawTags := extractValue(aiRawResponse, "HASHTAGS:")
        fTags := tags
        if len(fTags) == 0 {
            fTags = strings.Fields(rawTags)
        }

        // 3. NEW: If it didn't exist, save it to Content DB
        if wasEmpty && fHook != "" {
            fmt.Printf("📝 Template for %s is missing. Attempting to save new AI version...\n", p)
            // Updated to match your NotionClient.SaveTemplate(platform, hook, caption, cta, tags)
            errT := notion.SaveTemplate(p, fHook, fCaption, fCTA, fTags) 
            if errT != nil {
                fmt.Printf("🔴 SAVE TEMPLATE ERROR [%s]: %v\n", p, errT)
            } else {
                fmt.Printf("✨ SUCCESS: Saved new %s template to library!\n", p)
            }
        }

        recs = append(recs, Recommendation{
            Rank:      len(recs) + 1,
            Platform:  p,
            Reasoning: extractValue(aiRawResponse, "REASONING:"),
            ContentTemplate: &ContentTemplate{
                Hook: fHook, Caption: fCaption, CTA: fCTA, Hashtags: fTags,
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
        fmt.Printf("❌ HISTORY SAVE ERROR: %v\n", errSave)
    } else {
        fmt.Println("✅ SUCCESS: History written to Notion")
    }

    return output
}
