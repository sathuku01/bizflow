# Strategy: Integrating Keyword & Trend Analysis

This document outlines a strategy for integrating an external keyword/trend analysis API into the BizFlow agent. The goal is to evolve the agent from a purely logic-based recommender into a more sophisticated, data-driven consultant.

## How This Enhances Agentic Capabilities

For an AI Agent Hackathon, demonstrating "agentic behavior" is key. An agent isn't just a model; it's a system that can reason, plan, and use tools to achieve a goal. Integrating keyword analysis directly enhances these capabilities in several ways:

1.  **Tool Use (Core Agentic Behavior):**
    *   The agent will gain the ability to use an external tool: the `KeywordAPI`. This is a fundamental agentic pattern. It demonstrates the capacity to interact with and leverage external systems to gather information that it does not possess internally.

2.  **Sophisticated Reasoning & Planning:**
    *   The agent's reasoning process becomes multi-step and more intelligent. Instead of a simple one-step recommendation, the agent's internal "thought process" becomes:
        1.  "First, I need to understand the business niche."
        2.  "Next, I will use my `KeywordAPI` tool to find high-volume, low-competition keywords for this niche."
        3.  "Then, I will evaluate which marketing platforms are best suited for these specific keywords."
        4.  "Finally, I will synthesize this data to form my final recommendation and explain *why* it's the best choice based on the data I found."
    *   This transforms the agent from a simple calculator into a system that formulates and executes a research plan.

3.  **Data-Driven Decision Making:**
    *   The recommendations will no longer be based solely on generic, hardcoded metadata. They will be dynamically influenced by real-world (or mocked real-world) data. This demonstrates a higher level of intelligence and adaptability. The agent's advice for a "retail" business could change from one day to the next if trend data shifts.

4.  **Enhanced Explanation Quality:**
    *   The agent can provide much more powerful and convincing explanations. Instead of saying, "Instagram is good for retail," it can say, "I recommend Instagram because your core keyword, 'handmade jewelry,' has shown a 30% increase in engagement there over the last month, presenting a clear market opportunity."

## Proposed Implementation Plan

Here is a phased approach to integrate this functionality.

### Phase 1: Define the "Tool" - The Keyword API Client

We will first define the interface for our new tool.

*   **Location:** Create a new package `internal/trends/`.
*   **Interface:** Define a `KeywordAPIClient` and its data structures.

```go
// internal/trends/client.go

package trends

import "biz-flow/internal/core"

// KeywordStats represents data for a single keyword.
type KeywordStats struct {
    Keyword          string
    SearchVolume     int     // e.g., Monthly search volume
    Competition      float64 // 0.0 (low) to 1.0 (high)
    SuggestedPlatforms []core.Platform // Platforms where this keyword is trending
}

// KeywordAPIClient defines the interface for a tool that fetches keyword data.
type KeywordAPIClient interface {
    GetKeywordData(niche string) ([]KeywordStats, error)
}
```

### Phase 2: Create a Mock Implementation

To allow for development and testing without a live API, we'll create a mock client that returns hardcoded data.

*   **Location:** `internal/trends/mock_client.go`
*   **Implementation:** The `MockKeywordClient` will implement the `KeywordAPIClient` interface and return a predefined list of `KeywordStats` when called. This allows the rest of the application to be built and tested.

### Phase 3: Create a New `KeywordOpportunityScorer`

This new scorer will use the data from the API client to score platforms.

*   **Location:** `internal/scoring/keyword_scorer.go`
*   **Logic:**
    1.  The `KeywordOpportunityScorer` will be initialized with a `KeywordAPIClient`.
    2.  Its `Score()` method will call `client.GetKeywordData()` using the `business.Description` or `business.Type` as the niche.
    3.  It will then analyze the results. For the platform being scored, if it appears in the `SuggestedPlatforms` list for a high-volume, low-competition keyword, it gets a high score (e.g., `1.0`). If not, it gets a low score (e.g., `0.1`).

### Phase 4: Integrate the Scorer and Enhance Reasoning

Finally, we'll plug the new scorer into the main pipeline and improve the AI's reasoning.

1.  **Update the Orchestrator (`internal/handler/agent_handler.go`):**
    *   Add the `NewKeywordOpportunityScorer()` to the list of scorers. The existing logic will automatically incorporate its score into the final ranking.

2.  **Enhance the AI Reasoning Prompt:**
    *   The prompt sent to the LLM for explanation will be augmented with the discovered keyword data.

    *   **Before:**
        `"Explain why Instagram is a good platform for a retail business..."`

    *   **After:**
        `"Explain why Instagram is a good platform for a retail business. Your explanation MUST incorporate the following data I found: The keyword 'handmade jewelry' has a high search volume and is strongly trending on Instagram. Emphasize this market opportunity."`

This structured approach allows us to add significant "agentic" intelligence to the BizFlow application, making it a much more compelling project for a hackathon.
