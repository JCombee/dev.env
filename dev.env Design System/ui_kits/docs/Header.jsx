// Header — sticky top nav. Active link inferred from `page` prop.
function Header({ page = 'home' }) {
  const [theme, setTheme] = useState(() => {
    try { return localStorage.getItem('dev-env-theme') || 'light'; } catch { return 'light'; }
  });
  useEffect(() => {
    document.documentElement.setAttribute('data-theme', theme);
    try { localStorage.setItem('dev-env-theme', theme); } catch {}
  }, [theme]);

  const links = [
    { id: 'home', label: 'home',      href: 'index.html' },
    { id: 'docs', label: 'docs',      href: 'docs.html' },
    { id: 'ref',  label: 'cli',       href: 'reference.html' },
    { id: 'log',  label: 'changelog', href: 'changelog.html' },
  ];

  return (
    <header className="app-header">
      <div className="container app-header__inner">
        <a href="index.html" className="app-header__wordmark" style={{ textDecoration: 'none' }}>
          <Logomark size={22} dark={theme === 'dark'}/>
          <Wordmark height={18}/>
          <Badge style={{ marginLeft: 4 }}>v0.6.5</Badge>
        </a>
        <nav className="app-header__nav">
          {links.map(l => (
            <a key={l.id} href={l.href} className={l.id === page ? 'active' : ''}>{l.label}</a>
          ))}
        </nav>
        <div className="app-header__right">
          <a className="app-header__nav" href="https://github.com/JCombee/dev.env" target="_blank" rel="noopener" style={{ display: 'flex', alignItems: 'center', gap: 6, padding: '6px 10px', fontFamily: 'var(--font-mono)', fontSize: 13, color: 'var(--fg-2)' }}>
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round"><path d="M9 19c-5 1.5-5-2.5-7-3m14 6v-3.87a3.37 3.37 0 0 0-.94-2.61c3.14-.35 6.44-1.54 6.44-7A5.44 5.44 0 0 0 20 4.77 5.07 5.07 0 0 0 19.91 1S18.73.65 16 2.48a13.38 13.38 0 0 0-7 0C6.27.65 5.09 1 5.09 1A5.07 5.07 0 0 0 5 4.77a5.44 5.44 0 0 0-1.5 3.78c0 5.42 3.3 6.61 6.44 7A3.37 3.37 0 0 0 9 18.13V22"/></svg>
            <span>github</span>
          </a>
          <button className="theme-toggle" onClick={() => setTheme(theme === 'dark' ? 'light' : 'dark')} aria-label="toggle theme">
            {theme === 'dark' ? '☼' : '☾'}
          </button>
        </div>
      </div>
    </header>
  );
}

Object.assign(window, { Header });
