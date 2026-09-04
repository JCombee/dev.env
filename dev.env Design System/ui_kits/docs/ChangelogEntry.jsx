// ChangelogEntry — one version with its commit groups.
function ChangelogEntry({ version, date, sections, latest = false }) {
  return (
    <article style={{ display: 'grid', gridTemplateColumns: '160px 1fr', gap: 32, padding: '32px 0', borderTop: '1px solid var(--border-1)' }}>
      <div style={{ position: 'sticky', top: 'calc(var(--header-h) + 24px)', alignSelf: 'flex-start' }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: 8, marginBottom: 4 }}>
          <h3 style={{ fontFamily: 'var(--font-mono)', fontWeight: 600, fontSize: 18, color: 'var(--fg-1)', margin: 0 }}>v{version}</h3>
          {latest && <Badge variant="new">latest</Badge>}
        </div>
        <div style={{ fontFamily: 'var(--font-mono)', fontSize: 12, color: 'var(--fg-3)' }}>{date}</div>
      </div>
      <div>
        {sections.map((s, i) => (
          <section key={i} style={{ marginBottom: 24 }}>
            <div className="eyebrow" style={{ marginBottom: 12, color: s.tone === 'feat' ? 'var(--color-success)' : s.tone === 'fix' ? 'var(--color-warning)' : 'var(--fg-3)' }}>
              {s.label}
            </div>
            <ul style={{ listStyle: 'none', padding: 0, margin: 0, display: 'grid', gap: 8 }}>
              {s.entries.map((e, j) => (
                <li key={j} style={{ display: 'flex', alignItems: 'flex-start', gap: 12, fontSize: 14, lineHeight: 1.55, color: 'var(--fg-1)' }}>
                  <span style={{ color: s.tone === 'feat' ? 'var(--color-lime-deep)' : 'var(--fg-3)', flexShrink: 0, marginTop: 2 }}>—</span>
                  <span>
                    {e.msg}{' '}
                    {e.sha && <a href={`https://github.com/JCombee/dev.env/commit/${e.sha}`} target="_blank" rel="noopener" style={{ fontFamily: 'var(--font-mono)', fontSize: 12, color: 'var(--fg-3)', borderBottom: '1px solid var(--border-1)' }}>{e.sha.slice(0, 7)}</a>}
                  </span>
                </li>
              ))}
            </ul>
          </section>
        ))}
      </div>
    </article>
  );
}

Object.assign(window, { ChangelogEntry });
