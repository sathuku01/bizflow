package ai

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/jomei/notionapi"
)

type NotionClient struct {
	client      *notionapi.Client
	contentDBID notionapi.DatabaseID
	historyDBID notionapi.DatabaseID
}

func NewNotionClient() *NotionClient {
	token := os.Getenv("NOTION_API_KEY")
	contentDB := os.Getenv("NOTION_CONTENT_DB_ID")
	historyDB := os.Getenv("NOTION_HISTORY_DB_ID")

	if token == "" || contentDB == "" || historyDB == "" {
		panic("Notion environment variables missing")
	}

	return &NotionClient{
		client:      notionapi.NewClient(notionapi.Token(token)),
		contentDBID: notionapi.DatabaseID(contentDB),
		historyDBID: notionapi.DatabaseID(historyDB),
	}
}

// FetchTemplate retrieves content from the Library/Content database
func (n *NotionClient) FetchTemplate(platform string) (h, c, ct string, tags []string, err error) {
    query := &notionapi.DatabaseQueryRequest{
        Filter: &notionapi.PropertyFilter{
            Property: "Platform",
            RichText: &notionapi.TextFilterCondition{Equals: platform},
        },
        PageSize: 1,
    }
    resp, err := n.client.Database.Query(context.Background(), n.contentDBID, query)
    
    // 1. If no row exists at all
    if err != nil || len(resp.Results) == 0 {
        return "", "", "", nil, fmt.Errorf("not found")
    }

    page := resp.Results[0]
    extract := func(propName string) string {
        if p, ok := page.Properties[propName].(*notionapi.RichTextProperty); ok && len(p.RichText) > 0 {
            return strings.TrimSpace(p.RichText[0].Text.Content)
        }
        return ""
    }

    hook := extract("Hook")
    
    // 2. CRITICAL: If the row exists but the Hook is empty, treat as "not found"
    if hook == "" {
        return "", "", "", nil, fmt.Errorf("template exists but is empty")
    }

    return hook, extract("Caption"), extract("CTA"), []string{}, nil
}
// SaveQuery logs the interaction to the History database
func (n *NotionClient) SaveQuery(input BusinessInput, output AgentOutput) error {
	fmt.Printf("⏳ Saving record to History DB [%s]...\n", n.historyDBID)

	props := notionapi.Properties{}

	props["Business Type"] = &notionapi.TitleProperty{
		Title: []notionapi.RichText{{Text: &notionapi.Text{Content: input.BusinessType}}},
	}
	props["Description"] = &notionapi.RichTextProperty{
		RichText: []notionapi.RichText{{Text: &notionapi.Text{Content: input.Description}}},
	}
	props["Budget"] = &notionapi.NumberProperty{Number: input.MonthlyBudget}

	if input.PrimaryGoal != "" {
		// Ensure 'Goal' is a Select property in Notion
		props["Goal"] = &notionapi.SelectProperty{
			Select: notionapi.Option{Name: input.PrimaryGoal},
		}
	}

	if len(input.Channels) > 0 {
		var options []notionapi.Option
		for _, ch := range input.Channels {
			options = append(options, notionapi.Option{Name: ch})
		}
		props["Channels"] = &notionapi.MultiSelectProperty{MultiSelect: options}
	}

	if len(output.Recommendations) > 0 {
		props["Top Platform"] = &notionapi.SelectProperty{
			Select: notionapi.Option{Name: output.Recommendations[0].Platform},
		}
	}

	formattedText := formatAIOutput(output)
	props["AI Output"] = &notionapi.RichTextProperty{
		RichText: splitToRichText(formattedText, 1900),
	}

	resp, err := n.client.Page.Create(context.Background(), &notionapi.PageCreateRequest{
		Parent:     notionapi.Parent{DatabaseID: n.historyDBID},
		Properties: props,
	})

	if err != nil {
		fmt.Printf("❌ HISTORY ERROR: %+v\n", err)
		return err
	}

	fmt.Printf("✅ HISTORY SUCCESS: Page ID %s\n", resp.ID)
	return nil
}

// SaveTemplate creates a new entry in the Content/Library database if it was missing
func (n *NotionClient) SaveTemplate(platform, hook, caption, cta string, tags []string) error {
	fmt.Printf("🔍 DEBUG: Saving new AI template to Content DB for %s...\n", platform)

	props := notionapi.Properties{
		"Platform": &notionapi.TitleProperty{
			Title: []notionapi.RichText{{Text: &notionapi.Text{Content: platform}}},
		},
		"Hook": &notionapi.RichTextProperty{
			RichText: []notionapi.RichText{{Text: &notionapi.Text{Content: hook}}},
		},
		"Caption": &notionapi.RichTextProperty{
			RichText: []notionapi.RichText{{Text: &notionapi.Text{Content: caption}}},
		},
		"CTA": &notionapi.RichTextProperty{
			RichText: []notionapi.RichText{{Text: &notionapi.Text{Content: cta}}},
		},
		"Notes": &notionapi.RichTextProperty{
			RichText: []notionapi.RichText{{Text: &notionapi.Text{Content: "AI Generated"}}},
		},
	}

	// Format Multi-Select Hashtags (Notion strictly requires unique Names)
	var tagOptions []notionapi.Option
	seen := make(map[string]bool)
	for _, t := range tags {
		cleanTag := strings.TrimSpace(strings.TrimPrefix(t, "#"))
		if cleanTag != "" && !seen[cleanTag] {
			tagOptions = append(tagOptions, notionapi.Option{Name: cleanTag})
			seen[cleanTag] = true
		}
	}
	props["Hashtags"] = &notionapi.MultiSelectProperty{
		MultiSelect: tagOptions,
	}

	resp, err := n.client.Page.Create(context.Background(), &notionapi.PageCreateRequest{
		Parent:     notionapi.Parent{DatabaseID: n.contentDBID},
		Properties: props,
	})

	if err != nil {
		fmt.Printf("❌ CONTENT DB ERROR: %v\n", err)
		return err
	}
	fmt.Printf("✅ CONTENT DB SUCCESS: Created Template ID %s\n", resp.ID)
	return nil
}

func formatAIOutput(out AgentOutput) string {
	var sb strings.Builder
	sb.WriteString("STRATEGIC ADVICE:\n" + out.StrategicAdvice + "\n\n")
	for _, r := range out.Recommendations {
		sb.WriteString(fmt.Sprintf("--- %s (Rank %d) ---\n", r.Platform, r.Rank))
		sb.WriteString("REASONING: " + r.Reasoning + "\n\n")
	}
	return sb.String()
}

func splitToRichText(text string, limit int) []notionapi.RichText {
	var results []notionapi.RichText
	runes := []rune(text)
	for i := 0; i < len(runes); i += limit {
		end := i + limit
		if end > len(runes) {
			end = len(runes)
		}
		results = append(results, notionapi.RichText{
			Type: "text",
			Text: &notionapi.Text{Content: string(runes[i:end])},
		})
	}
	return results
}