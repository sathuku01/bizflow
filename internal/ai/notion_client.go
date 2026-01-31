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

func (n *NotionClient) FetchTemplate(platform string) (h, c, ct string, tags []string, err error) {
	query := &notionapi.DatabaseQueryRequest{
		Filter: &notionapi.PropertyFilter{
			Property: "Platform",
			RichText: &notionapi.TextFilterCondition{Equals: platform},
		},
		PageSize: 1,
	}
	resp, err := n.client.Database.Query(context.Background(), n.contentDBID, query)
	if err != nil || len(resp.Results) == 0 {
		return "", "", "", nil, fmt.Errorf("no template found")
	}

	page := resp.Results[0]
	extract := func(propName string) string {
		if p, ok := page.Properties[propName].(*notionapi.RichTextProperty); ok && len(p.RichText) > 0 {
			return p.RichText[0].Text.Content
		}
		return ""
	}
	return extract("Hook"), extract("Caption"), extract("CTA"), []string{}, nil
}

func (n *NotionClient) SaveQuery(input BusinessInput, output AgentOutput) error {
    fmt.Printf("Attempting to save to Database ID: %s\n", n.historyDBID)
    
    props := notionapi.Properties{}

    // 1. Business Type (MUST be a 'Title' type column in Notion)
    props["Business Type"] = &notionapi.TitleProperty{
        Title: []notionapi.RichText{{Text: &notionapi.Text{Content: input.BusinessType}}},
    }

    // 2. Description (Text/Rich Text)
    props["Description"] = &notionapi.RichTextProperty{
        RichText: []notionapi.RichText{{Text: &notionapi.Text{Content: input.Description}}},
    }

    // 3. Budget (Number)
    props["Budget"] = &notionapi.NumberProperty{Number: input.MonthlyBudget}

    // 4. Goal (Select)
    if input.PrimaryGoal != "" {
        props["Goal"] = &notionapi.SelectProperty{
            Select: notionapi.Option{Name: input.PrimaryGoal},
        }
    }

    // 5. Channels (Multi-select)
    if len(input.Channels) > 0 {
        var options []notionapi.Option
        for _, ch := range input.Channels {
            options = append(options, notionapi.Option{Name: ch})
        }
        props["Channels"] = &notionapi.MultiSelectProperty{MultiSelect: options}
    }

    // 6. Top Platform (Select)
    if len(output.Recommendations) > 0 {
        props["Top Platform"] = &notionapi.SelectProperty{
            Select: notionapi.Option{Name: output.Recommendations[0].Platform},
        }
    }

    // 7. AI Output (Text/Rich Text)
    formattedText := formatAIOutput(output)
    props["AI Output"] = &notionapi.RichTextProperty{
        RichText: splitToRichText(formattedText, 1900),
    }

    // CREATE THE PAGE
    resp, err := n.client.Page.Create(context.Background(), &notionapi.PageCreateRequest{
        Parent:     notionapi.Parent{DatabaseID: n.historyDBID},
        Properties: props,
    })

    if err != nil {
        return fmt.Errorf("Notion API Error: %w", err)
    }

    fmt.Printf("✅ Success! New Page Created: %s\n", resp.ID)
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
		if end > len(runes) { end = len(runes) }
		results = append(results, notionapi.RichText{
			Type: "text",
			Text: &notionapi.Text{Content: string(runes[i:end])},
		})
	}
	return results
}