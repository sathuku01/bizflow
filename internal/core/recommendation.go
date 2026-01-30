package core

type ContentTemplate struct {
    Hook     string   `json:"hook"`
    Caption  string   `json:"caption"`
    CTA      string   `json:"cta"`
    Hashtags []string `json:"hashtags"`
}

type Recommendation struct {
    Rank            int              `json:"rank"`
    Platform        Platform         `json:"platform"`
    Reasoning       string           `json:"reasoning"`
    Score           float64          `json:"score"`
    ContentTemplate *ContentTemplate `json:"content_template,omitempty"`
}

type ConsultationResult struct {
    Recommendations []Recommendation `json:"recommendations"`
    StrategicAdvice string           `json:"strategic_advice"`
    Risks           []string         `json:"risks"`
    Persona         string           `json:"persona"`
}