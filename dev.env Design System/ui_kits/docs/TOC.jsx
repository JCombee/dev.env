// TOC — right-side table of contents for the docs reader.
function TOC() {
  const items = [
    { id: 'getting-started', label: 'Getting started' },
    { id: 'installation',    label: 'Installation' },
    { id: 'init',            label: 'Initialise a project' },
    { id: 'config',          label: 'Project configuration' },
    { id: 'start',           label: 'Start the services' },
    { id: 'stop',            label: 'Stop when you\'re done' },
    { id: 'next',            label: 'Where next' },
  ];
  const [active, setActive] = useState('getting-started');
  return (
    <aside style={{ position: 'sticky', top: 'calc(var(--header-h) + 24px)', alignSelf: 'flex-start', width: 200, flexShrink: 0, paddingTop: 4 }}>
      <div className="eyebrow" style={{ marginBottom: 10, paddingLeft: 10 }}>on this page</div>
      <ul style={{ listStyle: 'none', padding: 0, margin: 0, display: 'grid', gap: 6, borderLeft: '1px solid var(--border-1)' }}>
        {items.map(i => (
          <li key={i.id}>
            <a
              href={`#${i.id}`}
              onClick={() => setActive(i.id)}
              style={{
                display: 'block',
                fontFamily: 'var(--font-mono)',
                fontSize: 12,
                color: active === i.id ? 'var(--fg-1)' : 'var(--fg-3)',
                padding: '3px 10px',
                marginLeft: -1,
                borderLeft: active === i.id ? '1px solid var(--color-lime-deep)' : '1px solid transparent',
              }}>{i.label}</a>
          </li>
        ))}
      </ul>
      <div style={{ marginTop: 24, padding: 12, border: '1px solid var(--border-1)', borderRadius: 6, background: 'var(--bg-2)' }}>
        <div className="eyebrow" style={{ marginBottom: 6 }}>edit</div>
        <a href="https://github.com/JCombee/dev.env" target="_blank" rel="noopener" style={{ fontFamily: 'var(--font-mono)', fontSize: 12, color: 'var(--fg-2)' }}>edit on github →</a>
      </div>
    </aside>
  );
}

Object.assign(window, { TOC });
