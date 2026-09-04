// CommandCard — used on the CLI reference page. One per command.
function CommandCard({ cmd, summary, synopsis, example, exampleOutput, flags = [], related = [] }) {
  return (
    <article style={{ border: '1px solid var(--border-2)', borderRadius: 6, padding: 28, marginBottom: 16, background: 'var(--bg-1)' }} id={cmd.replace(/\s+/g, '-')}>
      <header style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', marginBottom: 8, gap: 16, flexWrap: 'wrap' }}>
        <div style={{ display: 'flex', alignItems: 'baseline', gap: 12 }}>
          <h3 style={{ fontFamily: 'var(--font-mono)', fontWeight: 600, fontSize: 22, color: 'var(--fg-1)', margin: 0, letterSpacing: '-0.01em' }}>
            <span style={{ color: 'var(--fg-3)' }}>dev </span>{cmd}
          </h3>
        </div>
        <div style={{ display: 'flex', gap: 8 }}>
          {related.map(r => <Badge key={r}>{r}</Badge>)}
        </div>
      </header>
      <p style={{ fontSize: 15, color: 'var(--fg-2)', margin: '0 0 16px', maxWidth: 720, lineHeight: 1.55 }}>{summary}</p>
      <div className="eyebrow" style={{ marginBottom: 6 }}>synopsis</div>
      <pre style={{ margin: '0 0 16px', padding: 12, background: 'var(--bg-3)', border: '1px solid var(--border-1)', borderRadius: 4, fontFamily: 'var(--font-mono)', fontSize: 13, color: 'var(--fg-1)', lineHeight: 1.6 }}>{synopsis}</pre>
      {flags.length > 0 && (
        <>
          <div className="eyebrow" style={{ marginBottom: 6 }}>flags</div>
          <div style={{ display: 'grid', gridTemplateColumns: 'min-content 1fr', gap: '6px 16px', fontSize: 13, marginBottom: 16 }}>
            {flags.map(f => (
              <React.Fragment key={f.name}>
                <code style={{ background: 'var(--bg-3)', border: '1px solid var(--border-1)', padding: '1px 6px', borderRadius: 3, fontFamily: 'var(--font-mono)', whiteSpace: 'nowrap', alignSelf: 'start' }}>{f.name}</code>
                <span style={{ color: 'var(--fg-2)' }}>{f.desc}</span>
              </React.Fragment>
            ))}
          </div>
        </>
      )}
      <div className="eyebrow" style={{ marginBottom: 6 }}>example</div>
      <pre style={{ margin: 0, padding: 14, background: 'var(--bg-1)', border: '1px solid var(--border-2)', borderRadius: 6, fontFamily: 'var(--font-mono)', fontSize: 13, color: 'var(--fg-1)', lineHeight: 1.6, overflowX: 'auto' }}>
        <div style={{ color: 'var(--fg-1)' }}>{example}</div>
        {exampleOutput && (
          <div style={{ marginTop: 6 }}>
            {exampleOutput.split('\n').map((ln, i) => {
              const first = ln.trimStart()[0];
              let glyphColor = 'inherit';
              if (first === '✓') glyphColor = 'var(--color-lime-deep)';
              else if (first === '~') glyphColor = 'var(--fg-3)';
              else if (first === '!') glyphColor = 'var(--color-warning)';
              else if (first === '✗') glyphColor = 'var(--color-error)';
              return (
                <div key={i} style={{ color: 'var(--fg-2)' }}>
                  {ln && (first === '✓' || first === '~' || first === '!' || first === '✗')
                    ? <><span style={{ color: glyphColor }}>{first}</span>{ln.slice(ln.indexOf(first) + 1)}</>
                    : ln}
                </div>
              );
            })}
          </div>
        )}
      </pre>
    </article>
  );
}

Object.assign(window, { CommandCard });
