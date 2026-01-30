package filters

import (
	"biz-flow/internal/core"
)

// PlatformFilter handles filtering platforms based on business characteristics
type PlatformFilter struct{}

// NewPlatformFilter creates a new platform filter
func NewPlatformFilter() *PlatformFilter {
	return &PlatformFilter{}
}

// FilterByBusinessType returns platforms relevant to the business type
func (pf *PlatformFilter) FilterByBusinessType(businessType core.BusinessType) []core.Platform {
	switch businessType {
	case core.Retail:
		return []core.Platform{
			core.Instagram,
			core.Facebook,
			core.TikTok,
			core.GoogleBusiness,
		}
	case core.Service:
		return []core.Platform{
			core.GoogleBusiness,
			core.Facebook,
			core.WhatsApp,
			core.Instagram,
		}
	case core.Digital:
		return []core.Platform{
			core.LinkedIn,
			core.Email,
			core.YouTube,
			core.Instagram,
		}
	default:
		// Return all platforms if type is unknown
		return core.GetAllPlatformNames()
	}
}

// FilterByLocation filters platforms based on business location characteristics
func (pf *PlatformFilter) FilterByLocation(business core.BusinessInput, platforms []core.Platform) []core.Platform {
	// If business is local, prioritize platforms good for local reach
	// If online-only, filter out platforms that are primarily local
	
	if business.IsLocal() {
		// Local businesses should definitely have Google My Business
		return pf.ensureIncluded(platforms, core.GoogleBusiness)
	}
	
	// Online-only businesses don't need location-specific filtering
	return platforms
}

// FilterByBudget filters platforms based on budget constraints
func (pf *PlatformFilter) FilterByBudget(business core.BusinessInput, platforms []core.Platform) []core.Platform {
	allPlatformMetadata := core.AllPlatforms()
	filtered := make([]core.Platform, 0)
	
	for _, platform := range platforms {
		metadata, exists := allPlatformMetadata[platform]
		if !exists {
			continue
		}
		
		// Low budget: only include platforms with MinBudget = 0 (organic)
		if business.HasLowBudget() {
			if metadata.MinBudget == 0 {
				filtered = append(filtered, platform)
			}
			continue
		}
		
		// Medium/High budget: include all platforms within budget
		if business.Budget >= metadata.MinBudget {
			filtered = append(filtered, platform)
		}
	}
	
	return filtered
}

// FilterByEffort filters platforms based on effort feasibility
func (pf *PlatformFilter) FilterByEffort(business core.BusinessInput, platforms []core.Platform) []core.Platform {
	allPlatformMetadata := core.AllPlatforms()
	filtered := make([]core.Platform, 0)
	
	for _, platform := range platforms {
		metadata, exists := allPlatformMetadata[platform]
		if !exists {
			continue
		}
		
		// For micro-businesses (implied solo operation), filter out high-effort video platforms
		// unless the business is highly visual (retail)
		if metadata.RequiresVideo && metadata.EffortLevel == core.HighEffort {
			// Only include if business is retail (visual products)
			if business.Type == core.Retail {
				filtered = append(filtered, platform)
			}
			// Skip for service and digital businesses
			continue
		}
		
		// Include all other platforms
		filtered = append(filtered, platform)
	}
	
	return filtered
}

// ApplyAllFilters applies all filtering rules to get relevant platforms
func (pf *PlatformFilter) ApplyAllFilters(business core.BusinessInput) []core.Platform {
	// Step 1: Filter by business type (primary filter)
	platforms := pf.FilterByBusinessType(business.Type)
	
	// Step 2: Filter by location
	platforms = pf.FilterByLocation(business, platforms)
	
	// Step 3: Filter by budget
	platforms = pf.FilterByBudget(business, platforms)
	
	// Step 4: Filter by effort feasibility
	platforms = pf.FilterByEffort(business, platforms)
	
	return platforms
}

// ensureIncluded ensures a platform is in the list (helper method)
func (pf *PlatformFilter) ensureIncluded(platforms []core.Platform, platformToInclude core.Platform) []core.Platform {
	for _, p := range platforms {
		if p == platformToInclude {
			return platforms // Already included
		}
	}
	// Add it
	return append(platforms, platformToInclude)
}

// GetFilteredCount returns the number of platforms after filtering
func (pf *PlatformFilter) GetFilteredCount(business core.BusinessInput) int {
	return len(pf.ApplyAllFilters(business))
}

// ExplainFiltering returns a human-readable explanation of why platforms were filtered
func (pf *PlatformFilter) ExplainFiltering(business core.BusinessInput) map[string]string {
	explanations := make(map[string]string)
	
	// Business type filtering
	explanations["business_type"] = string(business.Type) + " businesses are best suited for " + 
		formatPlatforms(pf.FilterByBusinessType(business.Type))
	
	// Budget filtering
	if business.HasLowBudget() {
		explanations["budget"] = "Low budget (<$50/month) limits platforms to organic-only channels"
	} else if business.HasMediumBudget() {
		explanations["budget"] = "Medium budget ($50-$200/month) allows organic and some paid channels"
	} else {
		explanations["budget"] = "High budget (>$200/month) enables all channel types including paid advertising"
	}
	
	// Location filtering
	if business.IsLocal() {
		explanations["location"] = "Local business benefits from location-based platforms like Google My Business"
	} else if business.IsOnlineOnly() {
		explanations["location"] = "Online-only business can leverage any platform regardless of location"
	}
	
	return explanations
}

// formatPlatforms converts a slice of platforms to a readable string
func formatPlatforms(platforms []core.Platform) string {
	if len(platforms) == 0 {
		return "none"
	}
	if len(platforms) == 1 {
		return string(platforms[0])
	}
	
	result := ""
	for i, p := range platforms {
		if i == len(platforms)-1 {
			result += "and " + string(p)
		} else if i == len(platforms)-2 {
			result += string(p) + " "
		} else {
			result += string(p) + ", "
		}
	}
	return result
}