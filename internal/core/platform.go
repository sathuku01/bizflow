package core

type Platform string

const (
    Instagram      Platform = "Instagram"
    Facebook       Platform = "Facebook"
    TikTok         Platform = "TikTok"
    GoogleBusiness Platform = "Google My Business"
    WhatsApp       Platform = "WhatsApp Business"
    Email          Platform = "Email/Newsletter"
    LinkedIn       Platform = "LinkedIn"
    YouTube        Platform = "YouTube"
)
type EffortLevel string

const (
    LowEffort    EffortLevel = "low"
    MediumEffort EffortLevel = "medium"
    HighEffort   EffortLevel = "high"
)

type PlatformMetadata struct {
    Name              Platform
    RequiresVisuals   bool
    RequiresVideo     bool
    MinBudget         float64
    EffortLevel       EffortLevel
    BestFor           []BusinessType
    SupportsHashtags  bool
    IsOrganic         bool
    IsPaid            bool
    ReachPotential    int         // 1-10 scale
    ConversionFocus   int         // 1-10 scale
}