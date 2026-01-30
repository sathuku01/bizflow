# BizFlow Development Plan

This document outlines the development tasks for each team member based on the project's scope and existing codebase.

## Project Overview: BizFlow Marketing Agent

The goal is to build a Go-based AI agent that acts as a marketing consultant for micro-businesses. It takes business details as input, filters and scores potential marketing platforms, and outputs a ranked list of recommendations with actionable advice and content templates.

---

## Developer A: Core Engine & Foundation (Started)

Developer A is responsible for building the foundational, non-AI logic of the application. This involves defining the core data structures and implementing the deterministic filtering and constraint validation rules. Much of this work appears to be complete or in-progress, but it needs to be finalized and tested.

**Your Tasks (Dev A):**

1.  **Finalize `internal/core/business.go`:**
    *   **Implement `Validate()` method:** Add a `Validate()` method to the `BusinessInput` struct. This method should check:
        *   `Type` is one of the defined constants (`Retail`, `Service`, `Digital`).
        *   `Description` is not empty.
        *   `Budget` is a non-negative number.
        *   `Goal` is one of the defined constants (`Awareness`, `Sales`).
        *   Return an `error` if any validation fails.
    *   **Implement `String()` method:** Add a `String()` method for easy logging and debugging that returns a formatted string of the business input.

2.  **Finalize `internal/core/platform.go`:**
    *   **Implement `GetPlatformMetadata()` helper:** The file currently defines `AllPlatforms()`. Create a new helper function `GetPlatformMetadata(platform Platform) (PlatformMetadata, bool)` that retrieves a single platform's metadata from the map and returns it, along with a boolean indicating if it was found.

3.  **Review and Test `internal/filters/`:**
    *   **Review `platform_filter.go` and `constraint_validator.go`:** The logic in these files is mostly complete. Review them against the "Constraint Rules" and "Filtering Logic" in `readthis.txt` to ensure they match perfectly.
    *   **Write Unit Tests:** Create `platform_filter_test.go` and `constraint_validator_test.go`. Write comprehensive tests for all public functions, especially:
        *   `ApplyAllFilters`: Test with different business inputs (`Retail`, `Service`, `Digital`, different budgets) to ensure the platform list is filtered correctly.
        *   `GetCombinedPenalty`: Test various business/platform combinations to verify the penalty calculation is correct.
        *   `IsValidPlatform`: Ensure this correctly identifies platforms that violate hard constraints.

---

## Developer B: AI Integration & Content Generation

Developer B is responsible for integrating with a Large Language Model (LLM) to provide the AI-driven parts of the consultation: persona inference, strategic advice, risk assessment, and content generation.

**Your Tasks (Dev B):**

1.  **Implement `internal/ai/client.go`:**
    *   Create a generic `LLMClient` interface with one method: `Generate(prompt string) (string, error)`.
    *   Implement a concrete client for a specific LLM provider (e.g., `OpenAIClient` or `GoogleAIClient`).
    *   The client's constructor should read the API key from an environment variable (e.g., `OPENAI_API_KEY`).

2.  **Implement `internal/ai/prompts.go`:**
    *   This file will store all prompt templates. Create exported string constants for each required generation task. Use `fmt.Sprintf` style placeholders.
    *   **`PersonaInferencePrompt`:** "Based on this business description: '%s', describe the ideal customer persona in one sentence."
    *   **`ContentTemplatePrompt`:** "For a '%s' business on '%s', generate a content template with a hook, caption, call-to-action, and three relevant hashtags. Format the output as a JSON object: {"hook": "...", "caption": "...", "cta": "...", "hashtags": ["...", "..."]}"
    *   **`StrategicAdvicePrompt`:** "Given a recommendation for a '%s' business to use '%s', provide a short, actionable strategic next step."
    *   **`RiskAssessmentPrompt`:** "What are two potential risks for a small business using '%s' for marketing?"
    *   **`ReasoningPrompt`:** "Explain in business terms why '%s' is a good marketing platform for a '%s' business with a goal of '%s'. Do not mention scores."

3.  **Implement `internal/ai/persona_inferrer.go`:**
    *   Create a function: `InferPersona(client LLMClient, business core.BusinessInput) (string, error)`.
    *   This function will format the `PersonaInferencePrompt` with the `business.Description`, call `client.Generate()`, and return the LLM's response.

4.  **Implement `internal/ai/content_generator.go`:**
    *   Create a function: `GenerateContentTemplate(client LLMClient, business core.BusinessInput, platform core.Platform) (*core.ContentTemplate, error)`.
    *   It will format the `ContentTemplatePrompt`, call the LLM, and then **unmarshal** the returned JSON string into the `core.ContentTemplate` struct. Handle JSON parsing errors gracefully.

---

## Developer C: Scoring, Reasoning & API Handler (Your Part)

You are responsible for the central orchestration layer. You will consume the foundational logic from Dev A and the AI services from Dev B to score platforms, build the final recommendation, and expose it all via an API endpoint.

**Your Tasks (Dev C):**

1.  **Implement the Scoring Engine in `internal/scoring/`:**
    *   **`scorer.go`**: Define a `Scorer` interface: `type Scorer interface { Score(business core.BusinessInput, platform core.Platform) float64 }`. The score should be normalized between 0.0 (worst fit) and 1.0 (perfect fit).
    *   **`budget_scorer.go`**: Implement the `Scorer` interface. The score can be `1.0 - penalty`, using the penalty from `constraint_validator.ValidateBudgetConstraints`.
    *   **`effort_scorer.go`**: Implement the `Scorer` interface. Calculate score as `1.0 - penalty` from `constraint_validator.ValidateEffortConstraints`.
    *   **`return_scorer.go`**: Implement the `Scorer` interface. The logic should be:
        *   Get platform metadata.
        *   If `business.Goal` is `Awareness`, the score is `float64(metadata.ReachPotential) / 10.0`.
        *   If `business.Goal` is `Sales`, the score is `float64(metadata.ConversionFocus) / 10.0`.
    *   **`audience_scorer.go`**: Implement the `Scorer` interface. Check if the `business.Type` is in the platform's `BestFor` list from its metadata. Return `1.0` if it is, `0.1` if it isn't.

2.  **Implement the Reasoning Engine in `internal/reasoning/`:**
    *   **`explainer.go`**: Create `ExplainRecommendations(client ai.LLMClient, business core.BusinessInput, platforms []core.Platform) (map[core.Platform]string, error)`. This function will loop through the given platforms, call the LLM with the `ReasoningPrompt` for each, and return a map of platform-to-reasoning strings.
    *   **`risk_assessor.go`**: Create `AssessRisks(client ai.LLMClient, platform core.Platform) ([]string, error)`. This will use the `RiskAssessmentPrompt` and parse the LLM response into a slice of strings.
    *   **`strategy_advisor.go`**: Create `GenerateStrategy(client ai.LLMClient, business core.BusinessInput, topPlatform core.Platform) (string, error)`. This will use the `StrategicAdvicePrompt`.

3.  **Implement the Central Orchestrator in `internal/handler/agent_handler.go`:**
    *   This is your main task. Create a function `GenerateConsultation(businessInput core.BusinessInput) (*core.ConsultationResult, error)`.
    *   **Step 1: Instantiate Services:** Create instances of the `PlatformFilter`, all your `Scorers`, and the `LLMClient`.
    *   **Step 2: Filter:** Call `platformFilter.ApplyAllFilters(businessInput)` to get the relevant platforms.
    *   **Step 3: Score & Rank:**
        *   For each filtered platform, calculate a `compositeScore`. A simple average is a good start: `(budgetScore + effortScore + returnScore + audienceScore) / 4`.
        *   Store these scores and sort the platforms to find the top 3.
    *   **Step 4: Generate AI Content (using Dev B's modules):**
        *   Call `persona_inferrer.InferPersona`.
        *   Call `reasoning.ExplainRecommendations` for the top 3 platforms.
        *   Call `content_generator.GenerateContentTemplate` for the #1 platform only.
        *   Call `reasoning.AssessRisks` for the #1 platform.
        *   Call `reasoning.GenerateStrategy` for the #1 platform.
    *   **Step 5: Assemble Final Result:** Create and populate the `core.ConsultationResult` struct with all the data you've gathered.

4.  **Implement the Web Server in `web/handler.go`:**
    *   Use the standard `net/http` package.
    *   Create a handler function, e.g., `consultationHandler`, for the route `/api/consult`.
    *   Inside the handler:
        *   Check for `http.MethodPost`.
        *   Decode the incoming request body (JSON) into a `core.BusinessInput` struct.
        *   Call `businessInput.Validate()` (from Dev A's work). Return a `400 Bad Request` if it fails.
        *   Call `handler.GenerateConsultation()` with the input.
        *   If there's an error, return a `500 Internal Server Error`.
        *   If successful, marshal the `ConsultationResult` to JSON and write it to the response with a `200 OK` status.
    *   In `cmd/agent/main.go`, replace the existing test code with the logic to start your web server.
