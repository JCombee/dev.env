// Hero — homepage hero: copy on left, live terminal on right.
function Hero() {
  const lines = [
    { text: '$ dev start' },
    { glyph: '✓', text: 'mysql:8.0              already running (shared)' },
    { glyph: '✓', text: 'redis:latest           started (shared)' },
    { glyph: '✓', text: 'elasticsearch:8.11     started (dedicated)' },
    { glyph: '✓', text: 'Database "my-project" created', dim: false },
    { glyph: '✓', text: '.env written' },
  ];
  return (
    <section style={{ position: 'relative', overflow: 'hidden', paddingTop: 96, paddingBottom: 96 }}>
      <div className="dot-grid" style={{ position: 'absolute', inset: 0, opacity: 0.18, pointerEvents: 'none' }}></div>
      <div className="container" style={{ position: 'relative' }}>
        <div style={{ display: 'grid', gridTemplateColumns: 'minmax(0,1fr) minmax(0,1fr)', gap: 64, alignItems: 'center' }}>
          <div>
            <Eyebrow style={{ marginBottom: 16 }}>/ 01 — go cli · v0.6.5</Eyebrow>
            <h1 style={{ fontFamily: 'var(--font-sans)', fontWeight: 600, fontSize: 64, lineHeight: 1.02, letterSpacing: '-0.025em', color: 'var(--fg-1)', margin: '0 0 20px' }}>
              Shared dev services.<br/>
              <span style={{ color: 'var(--fg-3)' }}>One binary.</span><br/>
              <span style={{ color: 'var(--fg-3)' }}>Zero drift.</span>
            </h1>
            <p style={{ fontSize: 18, lineHeight: 1.55, color: 'var(--fg-2)', maxWidth: 520, margin: '0 0 28px' }}>
              Every project needs MySQL, Redis, Elasticsearch. The typical solution is a <code style={{ background: 'var(--bg-3)', border: '1px solid var(--border-1)', padding: '1px 6px', borderRadius: 4, fontFamily: 'var(--font-mono)', fontSize: '0.88em' }}>docker-compose.yml</code> per project. dev.env runs one shared pool instead.
            </p>
            <div style={{ display: 'flex', gap: 12, alignItems: 'center' }}>
              <Button as="a" href="docs.html" arrow>Read the docs</Button>
              <Button as="a" href="reference.html" variant="secondary">CLI reference</Button>
            </div>
            <div style={{ marginTop: 20, display: 'flex', gap: 16, fontFamily: 'var(--font-mono)', fontSize: 13, color: 'var(--fg-3)' }}>
              <span><span style={{ color: 'var(--color-lime-deep)' }}>✓</span> macOS</span>
              <span><span style={{ color: 'var(--color-lime-deep)' }}>✓</span> Linux</span>
              <span><span style={{ color: 'var(--color-lime-deep)' }}>✓</span> Windows</span>
              <span style={{ color: 'var(--fg-3)' }}>· requires Docker</span>
            </div>
          </div>
          <div>
            <Terminal title="~/code/my-project" lines={lines} animate={true} minHeight={260}/>
          </div>
        </div>
      </div>
    </section>
  );
}

Object.assign(window, { Hero });
