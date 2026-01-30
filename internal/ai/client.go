package ai

import (
    "context"
    "fmt"
    "os"
    "github.com/google/generative-ai-go/genai"
    "google.golang.org/api/option"
)

// LLMClient wraps the Gemini client
type LLMClient struct {
    model *genai.GenerativeModel
}

var ErrMissingAPIKey = fmt.Errorf("GOOGLE_API_KEY not set")

// NewLLMClient initializes the client
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

    // Use a model known to support generateContent
    model := client.GenerativeModel("gemini-3-flash-preview")

    return &LLMClient{model: model}, nil
}

// GenerateText calls Gemini and returns a usable string
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

    // Collect text parts into a single string
    result := ""
    for _, p := range parts {
        switch v := p.(type) {
        case *genai.Text:
            result += string(*v)
        default:
            result += fmt.Sprintf("%v", v)
        }
    }

    return result, nil
}

// RunAgent uses the LLM and provides fallback if it fails
func RunAgent(input BusinessInput) AgentOutput {
    client, err := NewLLMClient()
    if err != nil {
        fmt.Println("LLM Error:", err)
        return ErrorOutput("AI backend unavailable")
    }

    prompt := BuildRecommendationPrompt(input)

    aiText, err := client.GenerateText(prompt)
    if err != nil {
        fmt.Println("AI Error:", err)
        aiText = "Fallback reasoning (LLM failed)."
    }

    return AgentOutput{
        Recommendations: []Recommendation{
            {
                Rank:      1,
                Platform:  "Instagram",
                Reasoning: aiText,
                ContentTemplate: &ContentTemplate{
                    Hook:     "Generated automatically — now agent improved.",
                    Caption:  "Your caption will improve with next tools.",
                    CTA:      "CTA coming soon.",
                    Hashtags: []string{"#agent"},
                },
            },
        },
        StrategicAdvice: "More detailed planning will be added when tool-calling is enabled.",
        Risks:           []string{"Next step: add tools"},
    }
}

// ErrorOutput returns fallback agent output
func ErrorOutput(msg string) AgentOutput {
    return AgentOutput{
        Recommendations: []Recommendation{},
        StrategicAdvice: msg,
        Risks:           []string{"AI unavailable"},
    }
}
