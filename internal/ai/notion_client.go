package ai

import (
	"context"
	"fmt"
	"os"

	"github.com/jomei/notionapi"
)

// NotionClient handles both content templates and user query history
type NotionClient struct {
	client      *notionapi.Client
	contentDBID notionapi.DatabaseID
	historyDBID notionapi.DatabaseID
}

// NewNotionClientWithDBs allows manual injection of DB IDs
func NewNotionClientWithDBs(contentDBID, historyDBID string) *NotionClient {
	token := os.Getenv("NOTION_API_KEY")
	if token == "" || contentDBID == "" || historyDBID == "" {
		panic("NOTION_API_KEY or DB IDs not set")
	}

	client := notionapi.NewClient(notionapi.Token(token))
	return &NotionClient{
		client:      client,
		contentDBID: notionapi.DatabaseID(contentDBID),
		historyDBID: notionapi.DatabaseID(historyDBID),
	}
}

// NewNotionClient initializes the client from environment variables
func NewNotionClient() *NotionClient {
	token := os.Getenv("NOTION_API_KEY")
	contentDB := os.Getenv("NOTION_CONTENT_DB_ID")
	historyDB := os.Getenv("NOTION_HISTORY_DB_ID")

	if token == "" || contentDB == "" || historyDB == "" {
		panic("NOTION_API_KEY, NOTION_CONTENT_DB_ID, or NOTION_HISTORY_DB_ID not set")
	}

	client := notionapi.NewClient(notionapi.Token(token))

	return &NotionClient{
		client:      client,
		contentDBID: notionapi.DatabaseID(contentDB),
		historyDBID: notionapi.DatabaseID(historyDB),
	}
}

// FetchTemplate fetches a template from the content database
func (n *NotionClient) FetchTemplate(platform string) (hook, caption, cta string, hashtags []string, err error) {
	query := &notionapi.DatabaseQueryRequest{
		Filter: &notionapi.PropertyFilter{
			Property: "Platform",
			RichText: &notionapi.TextFilterCondition{
				Equals: platform,
			},
		},
		PageSize: 1,
	}

	resp, err := n.client.Database.Query(context.Background(), n.contentDBID, query)
	if err != nil {
		return "", "", "", nil, fmt.Errorf("query failed: %v", err)
	}

	if len(resp.Results) == 0 {
		return "", "", "", nil, fmt.Errorf("no template found for platform %s", platform)
	}

	page := resp.Results[0]

	extract := func(prop notionapi.Property) string {
		if prop == nil {
			return ""
		}
		rtProp, ok := prop.(*notionapi.RichTextProperty)
		if !ok || len(rtProp.RichText) == 0 {
			return ""
		}
		return rtProp.RichText[0].Text.Content
	}

	hook = extract(page.Properties["Hook"])
	caption = extract(page.Properties["Caption"])
	cta = extract(page.Properties["CTA"])

	hashtags = []string{}
	if prop, ok := page.Properties["Hashtags"].(*notionapi.RichTextProperty); ok {
		for _, rt := range prop.RichText {
			hashtags = append(hashtags, rt.Text.Content)
		}
	}

	return hook, caption, cta, hashtags, nil
}

// SaveQuery saves a user's query & AI output to the history database
func (n *NotionClient) SaveQuery(input BusinessInput, output AgentOutput) error {
	props := notionapi.Properties{}

	// Notion requires a Title. If BusinessType is empty, we provide a placeholder.
	titleContent := input.BusinessType
	if titleContent == "" {
		titleContent = "New Business Query"
	}
	props["Business Type"] = &notionapi.TitleProperty{
		Title: []notionapi.RichText{
			{Text: &notionapi.Text{Content: titleContent}},
		},
	}

	props["Description"] = &notionapi.RichTextProperty{
		RichText: []notionapi.RichText{
			{Text: &notionapi.Text{Content: input.Description}},
		},
	}

	props["Budget"] = &notionapi.NumberProperty{Number: input.MonthlyBudget}

	// FIX: Only add the Goal property if it actually has a value.
	// Sending an empty string "" to a Select property causes a 400 error.
	if input.PrimaryGoal != "" {
		props["Goal"] = &notionapi.SelectProperty{
			Select: notionapi.Option{Name: input.PrimaryGoal},
		}
	}

	// Ensure Channels is at least an empty array [] so Notion doesn't see 'null'
	ms := make([]notionapi.Option, 0)
	for _, c := range input.Channels {
		if c != "" {
			ms = append(ms, notionapi.Option{Name: c})
		}
	}
	props["Channels"] = &notionapi.MultiSelectProperty{
		MultiSelect: ms,
	}

	outputText := fmt.Sprintf("%+v", output)
	props["AI Output"] = &notionapi.RichTextProperty{
		RichText: []notionapi.RichText{
			{
				Type: "text",
				Text: &notionapi.Text{Content: outputText},
			},
		},
	}

	_, err := n.client.Page.Create(context.Background(), &notionapi.PageCreateRequest{
		Parent:     notionapi.Parent{DatabaseID: n.historyDBID},
		Properties: props,
	})
	if err != nil {
		return fmt.Errorf("failed to save query to Notion: %v", err)
	}
	// Inside SaveQuery, right before "return nil"
	//fmt.Println("Save successful! Fetching updated history...")
	//n.ListLatestHistory()

	return nil
}

// func (n *NotionClient) ListLatestHistory() {
// 	query := &notionapi.DatabaseQueryRequest{
// 		Sorts: []notionapi.SortObject{
// 			{
// 				Timestamp: "created_time",
// 				Direction: "descending",
// 			},
// 		},
// 		PageSize: 5,
// 	}

// 	resp, err := n.client.Database.Query(context.Background(), n.historyDBID, query)
// 	if err != nil {
// 		fmt.Printf("Failed to verify history: %v\n", err)
// 		return
// 	}

// 	fmt.Println("--- Recent History Entries ---")
// 	for _, page := range resp.Results {
// 		// Attempt to extract the Business Type (Title)
// 		if prop, ok := page.Properties["Business Type"].(*notionapi.TitleProperty); ok && len(prop.Title) > 0 {
// 			fmt.Printf("- [%s] %s\n", page.CreatedTime.Format("15:04:05"), prop.Title[0].Text.Content)
// 		}
// 	}
// }