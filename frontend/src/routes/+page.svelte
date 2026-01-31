<script lang="ts">
  // ... your existing logic stays exactly the same ...
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
  let businessType = 'retail';
  let description = '';
  let location = '';
  let budget = '';
  let channels: string[] = [];
  let goal = 'awareness';
  let recommendations: RecommendationsResult | null = null;

  async function handleSubmit() {
    // Example validation
    if (!description || !location || !budget || channels.length === 0) {
      errors = 'Please fill in all fields and select at least one channel.';
      recommendations = null;
      return;
    }
    errors = '';
    // Example: Replace with your API call
    recommendations = {
      recommendations: [
        {
          rank: 1,
          platform: 'Instagram',
          reasoning: 'Best for visual products and local reach.',
          content_template: {
            hook: 'Show your unique mugs!',
            caption: 'Handmade with love in Austin.',
            cta: 'Order now!',
            hashtags: ['#ceramics', '#handmade', '#austin']
          }
        },
        {
          rank: 2,
          platform: 'Facebook',
          reasoning: 'Great for community engagement.',
          content_template: {
            hook: 'Join our mug lovers group!',
            caption: 'Exclusive deals for members.',
            cta: 'Join today!',
            hashtags: ['#muglife', '#shoplocal']
          }
        },
        {
          rank: 3,
          platform: 'GMB',
          reasoning: 'Boosts local search visibility.',
          content_template: null
        }
      ],
      strategic_advice: 'Focus on Instagram and Facebook for best ROI.',
      risks: ['Ad fatigue', 'Budget constraints']
    };
  }
</script>

<div class="container">
  <h1>Micro-Business Marketing Consultant</h1>

  {#if errors}
    <p class="error">{errors}</p>
  {/if}

  <form on:submit|preventDefault={handleSubmit}>
    <label>
      Business Type:
      <select bind:value={businessType}>
        <option value="retail">Retail</option>
        <option value="service">Service</option>
        <option value="digital">Digital</option>
      </select>
    </label>

    <label>
      Description:
      <input type="text" placeholder="e.g. Handmade ceramic mugs" bind:value={description} />
    </label>

    <label>
      Location:
      <input type="text" placeholder="e.g. Austin, TX" bind:value={location} />
    </label>

    <label>
      Monthly Budget ($):
      <input type="number" bind:value={budget} />
    </label>

    <label>
      Channels:
      <div class="checkbox-group">
        <label class="check-item"><input type="checkbox" value="Instagram" bind:group={channels} /> Instagram</label>
        <label class="check-item"><input type="checkbox" value="Facebook" bind:group={channels} /> Facebook</label>
        <label class="check-item"><input type="checkbox" value="TikTok" bind:group={channels} /> TikTok</label>
        <label class="check-item"><input type="checkbox" value="Google My Business" bind:group={channels} /> GMB</label>
      </div>
    </label>

    <label>
      Primary Goal:
      <select bind:value={goal}>
        <option value="awareness">Awareness</option>
        <option value="sales">Sales</option>
      </select>
    </label>

    <button type="submit">Get Recommendations</button>
  </form>

  {#if recommendations}
    <div class="results-section">
      <h2>Top 3 Recommendations</h2>
      {#each recommendations.recommendations as rec}
        <div class="rec-card">
          <h3>{rec.rank}. {rec.platform}</h3>
          <p class="reasoning">{rec.reasoning}</p>
          
          {#if rec.content_template}
            <div class="template-box">
              <h4>Content Template</h4>
              <p><strong>Hook:</strong> {rec.content_template.hook}</p>
              <p><strong>Caption:</strong> {rec.content_template.caption}</p>
              <p><strong>CTA:</strong> {rec.content_template.cta}</p>
              <p class="tags">{rec.content_template.hashtags.join(', ')}</p>
            </div>
          {/if}
        </div>
      {/each}

      <div class="advice-card">
        <h3>Strategic Advice</h3>
        <p>{recommendations.strategic_advice}</p>

        <h3>Risks to Consider</h3>
        <ul>
          {#each recommendations.risks as risk}
            <li>{risk}</li>
          {/each}
        </ul>
      </div>
    </div>
  {/if}
</div>

<style>
  /* Base Page Styling */
  :global(body) {
    background-color: #f8fafc;
    font-family: 'Inter', system-ui, sans-serif;
    color: #1e293b;
    margin: 0;
    padding: 2rem;
  }

  .container {
    max-width: 600px;
    margin: 0 auto;
  }

  h1 { color: #0f172a; margin-bottom: 2rem; }

  form {
    background: white;
    padding: 1rem;
    border-radius: 16px;
    box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1);
    border: 1px solid #e2e8f0;
  }

  label {
    display: block;
    margin-bottom: 1rem;
    font-weight: 600;
    font-size: 0.875rem;
    color: #475569;
  }

  input, select {
    width: 100%;
    box-sizing: border-box; /* Prevents input from overflowing */
    margin-top: 0.5rem;
    padding: 0.75rem;
    border: 1px solid #cbd5e1;
    border-radius: 8px;
    font-size: 1rem;
  }

  .checkbox-group {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 8px;
    margin-top: 0.5rem;
    background: #f1f5f9;
    padding: 1rem;
    border-radius: 8px;
  }

  .check-item {
    font-weight: normal;
    margin-bottom: 0;
    display: flex;
    align-items: center;
    gap: 8px;
  }

  button {
    width: 100%;
    background-color: #6366f1;
    color: white;
    padding: 1rem;
    border: none;
    border-radius: 8px;
    font-weight: 700;
    font-size: 1rem;
    cursor: pointer;
    margin-top: 1rem;
    transition: all 0.2s;
  }

  button:hover { background-color: #4f46e5; transform: translateY(-1px); }

  /* Recommendations */
  .rec-card {
    background: white;
    border-left: 6px solid #6366f1;
    padding: 1.5rem;
    margin: 1rem 0;
    border-radius: 12px;
    box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
  }

  .template-box {
    background: #f8fafc;
    padding: 1rem;
    border-radius: 8px;
    margin-top: 1rem;
    border: 1px dashed #cbd5e1;
  }

  .advice-card {
    background: #eef2ff;
    padding: 1.5rem;
    border-radius: 12px;
    border: 1px solid #c7d2fe;
    margin-top: 2rem;
  }

  .error {
    background: #fef2f2;
    color: #dc2626;
    padding: 1rem;
    border-radius: 8px;
    border: 1px solid #fecaca;
    text-align: center;
  }
  
  .tags { color: #6366f1; font-size: 0.85rem; font-family: monospace; }
</style>