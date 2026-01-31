<script lang="ts">
  interface ContentTemplate {
    hook: string;
    caption: string;
    cta: string;
    hashtags: string[];
  }

  interface Recommendation {
    rank: number;
    platform: string;
    reasoning: string;
    content_template: ContentTemplate | null;
  }

  interface RecommendationsResult {
    recommendations: Recommendation[];
    strategic_advice: string;
    risks: string[];
  }

  let errors = '';
  let loading = false;
  let businessType = 'retail';
  let description = '';
  let location = '';
  let budget = '';
  let channels: string[] = [];
  let goal = 'awareness';
  let recommendations: RecommendationsResult | null = null;

  async function handleSubmit() {
    if (!description || !budget || channels.length === 0) {
      errors = 'Please fill in description, budget, and select at least one channel.';
      recommendations = null;
      return;
    }

    errors = '';
    loading = true;
    recommendations = null;

    try {
      const response = await fetch('http://192.168.89.196:8080/run-agent', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          business_type: businessType,
          description: `${description} based in ${location}`,
          monthly_budget: parseFloat(budget),
          goal: goal,
          channels: channels
        })
      });

      if (!response.ok) throw new Error(`Server error: ${response.statusText}`);

      const data = await response.json();

      if (data.strategic_advice === "LLM Init Failed") {
        errors = "The AI Agent is having trouble starting. Check backend API keys.";
      } else {
        recommendations = data;
        if (window.innerWidth < 900) {
          setTimeout(() => {
            document.getElementById('results')?.scrollIntoView({ behavior: 'smooth' });
          }, 100);
        }
      }
    } catch (err) {
      errors = "Could not connect to the backend. Ensure Go server is running on 8080.";
      console.error(err);
    } finally {
      loading = false;
    }
  }

  function copyToClipboard(text: string) {
    navigator.clipboard.writeText(text);
    alert('Template copied to clipboard!');
  }
</script>

<div class="container">
  <header>
    <h1>BizFlow AI</h1>
    <p class="subtitle">Micro-Business Marketing Consultant</p>
  </header>

  {#if errors}
    <div class="error-banner">
      <span class="icon">⚠️</span> {errors}
    </div>
  {/if}

  <div class="main-grid">
    <section class="form-section">
      <form on:submit|preventDefault={handleSubmit}>
        <div class="field">
          <label for="type">Business Type</label>
          <select id="type" bind:value={businessType}>
            <option value="retail">Retail</option>
            <option value="service">Service</option>
            <option value="digital">Digital</option>
          </select>
        </div>

        <div class="field">
          <label for="desc">Description</label>
          <input id="desc" type="text" placeholder="e.g. Handmade ceramic mugs" bind:value={description} />
        </div>

        <div class="field">
          <label for="loc">Location</label>
          <input id="loc" type="text" placeholder="e.g. Austin, TX" bind:value={location} />
        </div>

        <div class="field">
          <label for="bud">Monthly Budget ($)</label>
          <input id="bud" type="number" inputmode="decimal" bind:value={budget} />
        </div>

        <div class="field">
          <label>Target Channels</label>
          <div class="checkbox-group">
            <label class="check-item"><input type="checkbox" value="Instagram" bind:group={channels} /> Instagram</label>
            <label class="check-item"><input type="checkbox" value="Facebook" bind:group={channels} /> Facebook</label>
            <label class="check-item"><input type="checkbox" value="TikTok" bind:group={channels} /> TikTok</label>
            <label class="check-item"><input type="checkbox" value="Google My Business" bind:group={channels} /> GMB</label>
          </div>
        </div>

        <div class="field">
          <label for="goal">Primary Goal</label>
          <select id="goal" bind:value={goal}>
            <option value="awareness">Awareness</option>
            <option value="sales">Sales</option>
          </select>
        </div>

        <button type="submit" disabled={loading}>
          {#if loading}
            <span class="spinner"></span> Analyzing Market...
          {:else}
            Generate Strategy
          {/if}
        </button>
      </form>
    </section>

    <section class="results-section" id="results">
      {#if recommendations}
        <div class="advice-card">
          <h3>Strategic Advice</h3>
          <p>{recommendations.strategic_advice}</p>
        </div>

        <h2>Requested Recommendations</h2>

        {@const filteredRecs = recommendations.recommendations.filter(rec => {
          const platform = rec.platform.toLowerCase();
          return channels.some(c => platform.includes(c.toLowerCase()) || c.toLowerCase().includes(platform));
        })}

        {#if filteredRecs.length > 0}
          {#each filteredRecs as rec}
            <div class="rec-card">
              <div class="rec-header">
                <span class="rank">#{rec.rank}</span>
                <h3>{rec.platform}</h3>
              </div>
              <p class="reasoning">{rec.reasoning}</p>
              
              {#if rec.content_template}
                <div class="template-box">
                  <div class="template-header">
                    <span>Content Template</span>
                    <button class="copy-btn" on:click={() => copyToClipboard(`${rec.content_template?.hook}\n${rec.content_template?.caption}`)}>
                      Copy
                    </button>
                  </div>
                  <div class="template-content">
                    <p><strong>Hook:</strong> {rec.content_template.hook}</p>
                    <p><strong>Caption:</strong> {rec.content_template.caption}</p>
                    <p><strong>CTA:</strong> {rec.content_template.cta}</p>
                    {#if rec.content_template.hashtags}
                      <p class="tags">{rec.content_template.hashtags.map(t => (t.startsWith('#') ? t : '#' + t)).join(' ')}</p>
                    {/if}
                  </div>
                </div>
              {/if}
            </div>
          {/each}
        {:else}
          <div class="empty-state">
            <p>The AI suggests other channels might be better, or naming didn't match. Try broader selections!</p>
          </div>
        {/if}

        {#if recommendations.risks && recommendations.risks.length > 0}
          <div class="risk-card">
            <h3>Risks to Consider</h3>
            <ul>
              {#each recommendations.risks as risk}
                <li>{risk}</li>
              {/each}
            </ul>
          </div>
        {/if}
      {:else if !loading}
        <div class="empty-state">
          <p>Fill in details to generate your strategy.</p>
        </div>
      {/if}
    </section>
  </div>
</div>

<style>
  :global(body) {
    background-color: #f1f5f9;
    font-family: 'Inter', -apple-system, sans-serif;
    color: #0f172a;
    margin: 0;
  }

  .container { max-width: 1200px; margin: 0 auto; padding: 2rem; }
  header { margin-bottom: 2rem; text-align: center; }
  h1 { font-size: 2.5rem; font-weight: 800; margin: 0; color: #4f46e5; }
  .subtitle { color: #64748b; font-size: 1.1rem; }

  .main-grid { display: grid; grid-template-columns: 400px 1fr; gap: 2rem; align-items: start; }

  @media (min-width: 901px) {
    .form-section { position: sticky; top: 2rem; background: white; padding: 2rem; border-radius: 1rem; box-shadow: 0 4px 6px -1px rgb(0 0 0 / 0.1); }
  }

  @media (max-width: 900px) {
    .container { padding: 1rem; }
    .main-grid { grid-template-columns: 1fr; gap: 1.5rem; }
    .form-section { position: static; background: white; padding: 1.5rem; border-radius: 1rem; }
    h1 { font-size: 1.8rem; }
  }

  .field { margin-bottom: 1.25rem; }
  label { display: block; font-weight: 600; margin-bottom: 0.4rem; font-size: 0.875rem; }
  input, select { width: 100%; padding: 0.75rem; border: 1px solid #cbd5e1; border-radius: 0.5rem; font-size: 1rem; box-sizing: border-box; }

  .checkbox-group { display: grid; grid-template-columns: 1fr 1fr; gap: 0.5rem; background: #f8fafc; padding: 0.75rem; border-radius: 0.5rem; border: 1px solid #e2e8f0; }
  .check-item { font-weight: normal; font-size: 0.85rem; display: flex; align-items: center; gap: 0.5rem; cursor: pointer; }

  button { width: 100%; background: #4f46e5; color: white; padding: 1rem; border: none; border-radius: 0.5rem; font-weight: 700; cursor: pointer; transition: all 0.2s; }
  button:hover:not(:disabled) { background: #4338ca; transform: translateY(-1px); }
  button:disabled { background: #94a3b8; cursor: not-allowed; }

  .rec-card { background: white; padding: 1.5rem; border-radius: 1rem; margin-bottom: 1.5rem; box-shadow: 0 4px 6px -1px rgb(0 0 0 / 0.1); border-left: 5px solid #4f46e5; }
  .rec-header { display: flex; align-items: center; gap: 1rem; margin-bottom: 1rem; }
  .rank { background: #4f46e5; color: white; padding: 0.25rem 0.75rem; border-radius: 999px; font-weight: bold; font-size: 0.8rem; }
  
  .template-box { margin-top: 1rem; border: 1px solid #e2e8f0; border-radius: 0.5rem; overflow: hidden; }
  .template-header { background: #f8fafc; padding: 0.6rem 1rem; font-size: 0.75rem; font-weight: bold; color: #64748b; border-bottom: 1px solid #e2e8f0; display: flex; justify-content: space-between; align-items: center; }
  .template-content { padding: 1rem; font-size: 0.95rem; line-height: 1.5; }
  
  .copy-btn { width: auto; padding: 0.2rem 0.6rem; font-size: 0.7rem; background: #e2e8f0; color: #475569; }
  .copy-btn:hover { background: #cbd5e1; transform: none; }

  .tags { color: #4f46e5; font-family: monospace; font-size: 0.85rem; margin-top: 0.8rem; }
  .advice-card { background: #4f46e5; color: white; padding: 1.5rem; border-radius: 1rem; margin-bottom: 2rem; }
  .risk-card { background: #fff1f2; border: 1px solid #fecdd3; padding: 1.5rem; border-radius: 1rem; margin-bottom: 2rem; }
  .risk-card h3 { color: #be123c; margin-top: 0; }
  .error-banner { background: #fef2f2; border: 1px solid #fecaca; color: #dc2626; padding: 1rem; border-radius: 0.5rem; margin-bottom: 2rem; display: flex; align-items: center; gap: 0.5rem; }
  .empty-state { text-align: center; padding: 5rem 1rem; color: #94a3b8; border: 2px dashed #cbd5e1; border-radius: 1rem; }
  .spinner { display: inline-block; width: 1.2rem; height: 1.2rem; border: 3px solid rgba(255,255,255,0.3); border-radius: 50%; border-top-color: #fff; animation: spin 0.8s linear infinite; margin-right: 0.5rem; vertical-align: middle; }
  @keyframes spin { to { transform: rotate(360deg); } }
</style>