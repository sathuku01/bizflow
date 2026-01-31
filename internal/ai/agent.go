package ai

import "fmt"

// BuildRecommendationPrompt creates the prompt for the recommendation agent

func BuildRecommendationPrompt(input BusinessInput) string {
    return fmt.Sprintf(`
You are a micro-business marketing consultant.
Return ONLY valid JSON in the structure:

{
 "recommendations": [
   {
     "rank": 1,
     "platform": "",
     "reasoning": "",
     "content_template": {
       "hook": "",
       "caption": "",
       "cta": "",
       "hashtags": []
     }
   },
   {
     "rank": 2,
     "platform": "",
     "reasoning": ""
   },
   {
     "rank": 3,
     "platform": "",
     "reasoning": ""
   }
 ],
 "strategic_advice": "",
 "risks": []
}

Business info:
Type: %s
Description: %s
Location: %s
Monthly Budget: %.2f
Current Channels: %v
Primary Goal: %s
`, input.BusinessType, input.Description, input.Location, input.MonthlyBudget, input.CurrentChannels, input.PrimaryGoal)
}
