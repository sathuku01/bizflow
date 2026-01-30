package ai

type BusinessInput struct {
    BusinessType    string   `json:"business_type"`  // retail/service/digital
    Description     string   `json:"description"`
    Location        string   `json:"location"`       // city/country or "online"
    MonthlyBudget   float64  `json:"monthly_budget"` // USD
    CurrentChannels []string `json:"current_channels"`
    PrimaryGoal     string   `json:"primary_goal"`   // awareness/sales
}


type ContentTemplate struct {
    Hook     string   `json:"hook"`
    Caption  string   `json:"caption"`
    CTA      string   `json:"cta"`
    Hashtags []string `json:"hashtags"`
}

type Recommendation struct {
    Rank            int              `json:"rank"`
    Platform        string           `json:"platform"`
    Reasoning       string           `json:"reasoning"`
    ContentTemplate *ContentTemplate `json:"content_template,omitempty"`
}

type AgentOutput struct {
    Recommendations  []Recommendation `json:"recommendations"`
    StrategicAdvice  string           `json:"strategic_advice"`
    Risks            []string         `json:"risks"`
}
