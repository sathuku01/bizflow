package core

type BusinessType string

const (
	Retail  BusinessType = "retail"
	Service BusinessType = "service"
	Digital BusinessType = "digital"
)

type MarketingGoal string

const (
	Awareness MarketingGoal = "awareness"
	Sales     MarketingGoal = "sales"
)

type BusinessInput struct {
	Type        BusinessType   `json:"type"`
	Description string         `json:"description"`
	Location    string         `json:"location"`
	Budget      float64        `json:"budget"`
	Channels    []string       `json:"channels"`
	Goal        MarketingGoal  `json:"goal"`
}

// IsLocal checks if the business is local (not online-only)
func (b BusinessInput) IsLocal() bool {
	return b.Location != "" && b.Location != "online"
}

// IsOnlineOnly checks if the business is online-only
func (b BusinessInput) IsOnlineOnly() bool {
	return b.Location == "" || b.Location == "online"
}

// HasLowBudget checks if budget is low (<$50/month)
func (b BusinessInput) HasLowBudget() bool {
	return b.Budget < 50
}

// HasMediumBudget checks if budget is medium ($50-$200/month)
func (b BusinessInput) HasMediumBudget() bool {
	return b.Budget >= 50 && b.Budget <= 200
}

// HasHighBudget checks if budget is high (>$200/month)
func (b BusinessInput) HasHighBudget() bool {
	return b.Budget > 200
}