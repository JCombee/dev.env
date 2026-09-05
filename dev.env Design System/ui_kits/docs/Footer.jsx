// Footer — three-column, mono labels.
function Footer() {
  const cols = [
    {
      label: 'docs',
      links: [
        ['Getting started', 'docs.html'],
        ['Project config',  'docs.html#config'],
        ['Port allocation', 'docs.html#ports'],
        ['Local overrides', 'docs.html#local'],
      ],
    },
    {
      label: 'reference',
      links: [
        ['CLI commands',     'reference.html'],
        ['Supported services', 'index.html#services'],
        ['Changelog',        'changelog.html'],
      ],
    },
    {
      label: 'project',
      links: [
        ['GitHub',  'https://github.com/JCombee/dev.env'],
        ['Issues',  'https://github.com/JCombee/dev.env/issues'],
        ['Releases','https://github.com/JCombee/dev.env/releases'],
        ['License', '#'],
      ],
    },
  ];

  return (
    <footer style={{ borderTop: '1px solid var(--border-1)', paddingTop: 64, paddingBottom: 48, background: 'var(--bg-1)' }}>
      <div className="container">
        <div style={{ display: 'grid', gridTemplateColumns: '1.4fr repeat(3, 1fr)', gap: 48, marginBottom: 48 }}>
          <div>
            <div style={{ display: 'flex', alignItems: 'center', gap: 10, marginBottom: 14 }}>
              <Logomark size={26}/>
              <Wordmark height={20}/>
            </div>
            <p style={{ fontSize: 14, color: 'var(--fg-2)', maxWidth: 280, lineHeight: 1.55, margin: 0 }}>
              Shared dev services. One binary. Zero config drift.
            </p>
            <div style={{ marginTop: 14, fontFamily: 'var(--font-mono)', fontSize: 12, color: 'var(--fg-3)' }}>
              v0.6.5 · MIT · made by <a href="https://github.com/JCombee" target="_blank" rel="noopener" style={{ color: 'var(--fg-2)', borderBottom: '1px solid var(--border-2)' }}>JCombee</a>
            </div>
          </div>
          {cols.map(c => (
            <div key={c.label}>
              <Eyebrow style={{ marginBottom: 14 }}>{c.label}</Eyebrow>
              <ul style={{ listStyle: 'none', padding: 0, margin: 0, display: 'grid', gap: 8 }}>
                {c.links.map(([t, h]) => (
                  <li key={t}><a href={h} style={{ fontFamily: 'var(--font-mono)', fontSize: 13, color: 'var(--fg-2)' }}>{t}</a></li>
                ))}
              </ul>
            </div>
          ))}
        </div>
        <hr className="divider"/>
        <div style={{ paddingTop: 20, display: 'flex', justifyContent: 'space-between', fontFamily: 'var(--font-mono)', fontSize: 12, color: 'var(--fg-3)' }}>
          <span>© dev.env</span>
          <span>built with one binary and no JavaScript framework, mostly</span>
        </div>
      </div>
    </footer>
  );
}

Object.assign(window, { Footer });
