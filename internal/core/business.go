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