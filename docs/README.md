BizFlow 🚀

Micro-Business Marketing Consultant Agent






🧠 Project Overview

BizFlow is an AI-driven Marketing Consultant Agent tailored to the unique needs of micro-businesses. It offers fast, contextual marketing strategy recommendations and ready-to-use content templates for social platforms — all from a simple form input. The tool is ideal for entrepreneurs, freelancers, and marketers seeking actionable ideas without the complexity of traditional marketing platforms.

BizFlow does not automate posting or handle persistent data yet— instead, it focuses on high-quality strategic insights coupled with clear rationales, making it a powerful idea generation assistant rather than a full campaign manager.

⚙️ How It Works

When a user submits their business details through an input form, BizFlow:

Normalizes the Form Input

User inputs are captured via a web form and internally transformed into structured JSON for processing.

Platform Iteration & AI Reasoning

The system evaluates a curated list of platforms (e.g., Instagram, Facebook, TikTok).

It then prompts an LLM to generate both reasoning and content templates tailored to the business context.

Template Backend Support

BizFlow attempts to fetch existing templates from a Notion database to guide generation.

If none exist, the LLM creates new ones which are saved back to Notion for future reuse.

Structured Output

Users get a ranked list of platform recommendations, strategic advice, and content ideas in a structured response.

🏛 Architecture Overview
bizflow/
├── cmd/agent/main.go            # Web server entrypoint
├── internal/
│   ├── ai/                     # Core AI logic and API integration
│   │   ├── client.go           # Agent workflow
│   │   └── content_generator.go # Data models for generation
│   └── core/                   # (Unused) business abstractions
│   └── filters/                # (Unused) filter logic
│   └── reasoning/              # (Unused) reasoning modules
│   └── scoring/                # (Unused) scoring modules
└── go.mod

📊 Data Flow Summary

A POST /run-agent request receives business input.

The handler calls ai.RunAgent.

The agent orchestrates LLM generation with optional Notion support.

Generated insights are compiled and returned as JSON.

All interactions are archived in Notion for reference.

🧩 Key Design Principles

Explainability Over Scores
Rather than opaque metrics, every recommendation includes human-readable rationale.

Focused Scope for Hackathon Success
Stateless interaction and single-call processing keep complexity low.

Template-Aware AI
Notion serves as a lightweight knowledge base, improving content quality over time.

🚀 Quick Start
🎯 Requirements

Create a .env file with:

OPENROUTER_API_KEY="sk-..."
NOTION_API_KEY="secret_..."
NOTION_HISTORY_DB_ID="..."
NOTION_CONTENT_DB_ID="..."

📦 Run Locally
go mod tidy
go run cmd/agent/main.go


The server will start on port 8080.

🧪 Example Request / Response
Example Input (form → JSON under the hood)
{
  "business_type": "retail",
  "description": "Used bookstore in a college town.",
  "location": "Cambridge, MA",
  "monthly_budget": 500,
  "goal": "Increase foot traffic",
  "channels": ["instagram"]
}

Example Output
{
  "recommendations": [
    {
      "platform": "Instagram",
      "reasoning": "...",
      "content_template": { ... }
    },
    { "platform": "Google My Business", ... }
  ],
  "strategic_advice": "...",
  "risks": ["High competition..."]
}

⚠️ Current Limitations & Non-Goals

BizFlow is intentionally scoped for idea generation, not automation or full campaign execution. Current limitations include:

❌ Does not execute or schedule posts
❌ Cannot manage ad campaigns
❌ Has no persistent user database
❌ No performance analytics or tracking
❌ No built-in A/B testing support
❌ Platform set is fixed and selective

🌱 What’s Next — Future Roadmap

Here are key enhancements that would elevate BizFlow into a more robust marketing platform:

✨ Automated Publishing & Scheduling

Enable users to queue and schedule posts directly from the platform.

📊 Ad Campaign Management & Analytics

Integrate tools to create, track, and optimize paid social campaigns.

🧠 Performance Tracking

Capture engagement metrics and business KPIs over time.

🧪 A/B Testing Insights

Propose structured testing plans and analyze performance across content variants.

🌎 Universal Platform Support

Expand recommendations beyond current channels to include email, SEO, and niche platforms.

💾 Persistent User Profiles

Allow users to store business profiles and past outputs for trend comparison.

🤝 Contributing

We welcome contributions! To contribute:

Fork the repository

Create a feature branch

Open a pull request detailing your changes

Optionally add tests and follow the existing test conventions.

📄 License

This project is licensed under the MIT License — see the LICENSE file for details.