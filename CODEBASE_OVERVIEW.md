# BizFlow Codebase Overview

This document provides a high-level explanation of the BizFlow application's structure, using the analogy of a consulting firm to clarify the roles of different code modules and how they interact.

## The "Consulting Firm" Analogy

Think of the entire application as a **small consulting firm**. A client submits a request, and different specialists work on it before a project manager assembles and delivers the final report.

### The Departments (Code Groups)

Here are the main "departments" in our firm, which correspond to the folders in your `internal/` directory:

1.  **`core/` - The Company Library**
    *   **Purpose:** This isn't a person, but the reference library that everyone uses. It contains the fundamental definitions: what a `BusinessInput` looks like, what a `Platform` is, and what the final `ConsultationResult` report should contain.
    *   **Who uses it:** Everyone.

2.  **`filters/` - The Junior Analyst (Dev A's Work)**
    *   **Purpose:** This analyst does the first, quick pass. Their job is to look at the client's business and throw out all the obviously bad marketing ideas based on simple, hardcoded rules (e.g., "This is a digital business, so let's not focus on local flyers").
    *   **Output:** A shorter, more manageable list of potential marketing platforms.

3.  **`scoring/` - The Team of Accountants (Dev C's Work)**
    *   **Purpose:** This is your team of data-crunchers. They are purely logical and don't write prose. They take the filtered list from the Analyst and assign a numerical score to each platform based on different criteria (Budget, Effort, Audience, etc.).
    *   **Output:** A ranked list of platforms with scores attached.

4.  **`ai/` - The Expert Writers (Dev B's Work)**
    *   **Purpose:** This department contains the raw writing talent. It has the connection to the Large Language Model (LLM) and the prompt templates needed to generate human-readable text. They don't know *what* to write about until they are told.
    *   **Output:** Raw text (like explanations, risks, and content ideas).

5.  **`reasoning/` - The Lead Strategist (Dev C's Work)**
    *   **Purpose:** You, acting as the lead strategist. You take the ranked list from the "Accountants" and instruct the "Expert Writers" on what specific pieces of text to generate (e.g., "Write me an explanation for this platform," "Assess the risks for that one").

### The Orchestrator (The Project Manager)

Your most important role is in **`handler/agent_handler.go`**. This is the office of the Project Manager. "Orchestration" simply means managing the workflow and telling everyone what to do and when.

Here is the step-by-step process you manage—this is the core of the "agent":

1.  **Request Arrives:** The `web/handler.go` (the Receptionist) receives a request and hands the client's file (`BusinessInput`) to you, the Project Manager in `agent_handler.go`.

2.  **You Call the Analyst:** You first send the file to the `filters` department. They do their quick analysis and hand you back a shorter list of relevant platforms.

3.  **You Call the Accountants:** You take this short list and give it to your `scoring` team. They crunch the numbers and give you a score for each platform. You then rank them to get your Top 3.

4.  **You Call the Writers:** The logical work is done. Now you need the report. You turn to the `reasoning` and `ai` departments with your Top 3 platforms and give them specific instructions:
    *   "Write a persuasive `Reasoning` paragraph for each of these three."
    *   "For the #1 platform, also write the `Risks`, `StrategicAdvice`, and a `ContentTemplate`."

5.  **You Assemble the Final Report:** The writers give you back all the text snippets. You take this text, combine it with the ranked list and scores from the accountants, and assemble everything into the final, polished `ConsultationResult` report.

6.  **You Deliver the Report:** You hand the finished report back to the Receptionist (`web/handler.go`), who sends it back to the client as a JSON response.

### Where the Devs Meet (The Handoffs)

*   **Dev A -> Dev C:** Dev A provides the foundational rules. You (**Dev C**) *use* those rules in your `scoring` and `handler` modules. The handoff is you calling their functions.

*   **Dev B -> Dev C:** Dev B provides the raw ability to write text (the `ai.LLMClient`). You (**Dev C**) *use* that ability in your `reasoning` and `handler` modules to generate specific, meaningful content. The handoff is you calling the AI client with your carefully crafted prompts.

In summary, when we say this is an **AI Agent**, we just mean it's a program that follows a plan (the orchestration steps above), uses tools (the scorers and the AI client), and produces a result that is more than the sum of its parts. Your code as **Dev C** is the brain that directs all these specialized tools.
