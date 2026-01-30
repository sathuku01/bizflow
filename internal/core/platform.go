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

// AllPlatforms returns a map of all platforms and their metadata
func AllPlatforms() map[Platform]PlatformMetadata {
	return map[Platform]PlatformMetadata{
		Instagram: {
			Name:            Instagram,
			RequiresVisuals: true,
			RequiresVideo:   false,
			MinBudget:       0,
			EffortLevel:     MediumEffort,
			BestFor:         []BusinessType{Retail, Service, Digital},
			SupportsHashtags: true,
			IsOrganic:       true,
			IsPaid:          true,
			ReachPotential:  9,
			ConversionFocus: 7,
		},
		Facebook: {
			Name:            Facebook,
			RequiresVisuals: true,
			RequiresVideo:   false,
			MinBudget:       0,
			EffortLevel:     MediumEffort,
			BestFor:         []BusinessType{Retail, Service},
			SupportsHashtags: false,
			IsOrganic:       true,
			IsPaid:          true,
			ReachPotential:  8,
			ConversionFocus: 8,
		},
		TikTok: {
			Name:            TikTok,
			RequiresVisuals: true,
			RequiresVideo:   true,
			MinBudget:       0,
			EffortLevel:     HighEffort,
			BestFor:         []BusinessType{Retail, Digital},
			SupportsHashtags: true,
			IsOrganic:       true,
			IsPaid:          true,
			ReachPotential:  10,
			ConversionFocus: 6,
		},
		GoogleBusiness: {
			Name:            GoogleBusiness,
			RequiresVisuals: true,
			RequiresVideo:   false,
			MinBudget:       0,
			EffortLevel:     LowEffort,
			BestFor:         []BusinessType{Retail, Service},
			SupportsHashtags: false,
			IsOrganic:       true,
			IsPaid:          true,
			ReachPotential:  7,
			ConversionFocus: 9,
		},
		WhatsApp: {
			Name:            WhatsApp,
			RequiresVisuals: false,
			RequiresVideo:   false,
			MinBudget:       0,
			EffortLevel:     LowEffort,
			BestFor:         []BusinessType{Service},
			SupportsHashtags: false,
			IsOrganic:       true,
			IsPaid:          false,
			ReachPotential:  5,
			ConversionFocus: 8,
		},
		Email: {
			Name:            Email,
			RequiresVisuals: false,
			RequiresVideo:   false,
			MinBudget:       0,
			EffortLevel:     MediumEffort,
			BestFor:         []BusinessType{Retail, Service, Digital},
			SupportsHashtags: false,
			IsOrganic:       true,
			IsPaid:          true,
			ReachPotential:  6,
			ConversionFocus: 9,
		},
		LinkedIn: {
			Name:            LinkedIn,
			RequiresVisuals: true,
			RequiresVideo:   false,
			MinBudget:       0,
			EffortLevel:     HighEffort,
			BestFor:         []BusinessType{Digital, Service},
			SupportsHashtags: false,
			IsOrganic:       true,
			IsPaid:          true,
			ReachPotential:  7,
			ConversionFocus: 8,
		},
		YouTube: {
			Name:            YouTube,
			RequiresVisuals: true,
			RequiresVideo:   true,
			MinBudget:       0,
			EffortLevel:     HighEffort,
			BestFor:         []BusinessType{Retail, Digital},
			SupportsHashtags: true,
			IsOrganic:       true,
			IsPaid:          true,
			ReachPotential:  9,
			ConversionFocus: 7,
		},
	}
}

// GetAllPlatformNames returns a list of all platform names
func GetAllPlatformNames() []Platform {
	return []Platform{
		Instagram,
		Facebook,
		TikTok,
		GoogleBusiness,
		WhatsApp,
		Email,
		LinkedIn,
		YouTube,
	}
}