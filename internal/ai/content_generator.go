package ai

// BusinessInput holds the micro-business information
type BusinessInput struct {
	BusinessType  string   `json:"business_type"`  // retail/service/digital
	Description   string   `json:"description"`    // product or service description
	Location      string   `json:"location"`       
	MonthlyBudget float64  `json:"monthly_budget"` 
	// Change this tag to "goal" to match your curl command
	PrimaryGoal   string   `json:"goal"`           
	Channels      []string `json:"channels"`
}

// ContentTemplate holds an example content snippet for a platform
type ContentTemplate struct {
	Hook     string   `json:"hook"`
	Caption  string   `json:"caption"`
	CTA      string   `json:"cta"`
	Hashtags []string `json:"hashtags"`
}

// Recommendation represents a ranked platform recommendation
type Recommendation struct {
	Rank            int              `json:"rank"`
	Platform        string           `json:"platform"`
	Reasoning       string           `json:"reasoning"`
	ContentTemplate *ContentTemplate `json:"content_template,omitempty"`
}

// AgentOutput is the full structured response from the agent
type AgentOutput struct {
	Recommendations []Recommendation `json:"recommendations"`
	StrategicAdvice string           `json:"strategic_advice"`
	Risks           []string         `json:"risks"`
}