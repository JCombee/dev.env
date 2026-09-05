// Sidebar — left nav on docs reader. Active link inferred from `current` prop.
function Sidebar({ current = 'getting-started' }) {
  const sections = [
    {
      label: 'start',
      items: [
        { id: 'introduction',     title: 'Introduction' },
        { id: 'getting-started',  title: 'Getting started' },
        { id: 'installation',     title: 'Installation' },
      ],
    },
    {
      label: 'commands',
      items: [
        { id: 'setup',   title: 'dev setup' },
        { id: 'init',    title: 'dev init' },
        { id: 'start',   title: 'dev start' },
        { id: 'stop',    title: 'dev stop' },
        { id: 'status',  title: 'dev status' },
        { id: 'exec',    title: 'dev exec' },
        { id: 'config',  title: 'dev config' },
        { id: 'global',  title: 'dev global' },
      ],
    },
    {
      label: 'config',
      items: [
        { id: 'project-config', title: 'Project config' },
        { id: 'env-gen',        title: '.env generation' },
        { id: 'env-map',        title: 'Variable mapping' },
        { id: 'local',          title: 'Local overrides' },
        { id: 'ports',          title: 'Port allocation' },
      ],
    },
    {
      label: 'services',
      items: [
        { id: 'databases',  title: 'Databases' },
        { id: 'search',     title: 'Search engines' },
        { id: 'caches',     title: 'Caches' },
        { id: 'queues',     title: 'Queues' },
        { id: 'storage',    title: 'Storage' },
        { id: 'broadcast',  title: 'Broadcasting' },
        { id: 'mail',       title: 'Mail' },
        { id: 'unsupported',title: 'Unsupported images' },
      ],
    },
  ];

  return (
    <aside style={{ position: 'sticky', top: 'calc(var(--header-h) + 24px)', alignSelf: 'flex-start', width: 220, flexShrink: 0, paddingTop: 4 }}>
      <div style={{ position: 'relative', marginBottom: 16 }}>
        <input
          placeholder="search docs"
          style={{ width: '100%', fontFamily: 'var(--font-mono)', fontSize: 13, padding: '8px 32px 8px 10px', border: '1px solid var(--border-2)', borderRadius: 6, background: 'var(--bg-1)', color: 'var(--fg-1)', outline: 'none' }}
        />
        <span style={{ position: 'absolute', right: 8, top: '50%', transform: 'translateY(-50%)' }}><Kbd>⌘K</Kbd></span>
      </div>
      <nav>
        {sections.map(s => (
          <div key={s.label} style={{ marginBottom: 22 }}>
            <div className="eyebrow" style={{ display: 'block', marginBottom: 8, paddingLeft: 10 }}>{s.label}</div>
            <ul style={{ listStyle: 'none', padding: 0, margin: 0, display: 'grid', gap: 1 }}>
              {s.items.map(it => {
                const active = it.id === current;
                return (
                  <li key={it.id}>
                    <a href={`#${it.id}`} style={{
                      display: 'block',
                      fontFamily: 'var(--font-mono)', fontSize: 13,
                      padding: '5px 10px', borderRadius: 4,
                      color: active ? 'var(--fg-1)' : 'var(--fg-2)',
                      background: active ? 'var(--bg-2)' : 'transparent',
                      borderLeft: active ? '2px solid var(--color-lime-deep)' : '2px solid transparent',
                      paddingLeft: 8,
                    }}>{it.title}</a>
                  </li>
                );
              })}
            </ul>
          </div>
        ))}
      </nav>
    </aside>
  );
}

Object.assign(window, { Sidebar });
