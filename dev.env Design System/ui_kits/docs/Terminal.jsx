// Terminal — stylized CLI frame with traffic-light bar, optional typing animation.
function Terminal({ lines = [], title = '~/code/my-project', minHeight = 220, animate = false }) {
  const [visibleCount, setVisibleCount] = useState(animate ? 0 : lines.length);
  useEffect(() => {
    if (!animate) return;
    setVisibleCount(0);
    let i = 0;
    const id = setInterval(() => {
      i += 1;
      setVisibleCount(i);
      if (i >= lines.length) clearInterval(id);
    }, 320);
    return () => clearInterval(id);
  }, [animate, lines.length]);

  const renderLine = (line, idx) => {
    if (typeof line === 'string') {
      return <div key={idx} style={{ color: 'var(--fg-1)' }}>{line}</div>;
    }
    const { glyph, glyphColor, text, dim } = line;
    return (
      <div key={idx}>
        {glyph ? <span style={{ color: glyphColor || 'var(--color-lime-deep)', display: 'inline-block', width: 18 }}>{glyph}</span> : <span style={{ display: 'inline-block', width: 18 }}></span>}
        <span style={{ color: dim ? 'var(--fg-3)' : 'var(--fg-1)' }}>{text}</span>
      </div>
    );
  };

  return (
    <div style={{ border: '1px solid var(--border-2)', borderRadius: 8, overflow: 'hidden', boxShadow: 'var(--elevation-2)', background: 'var(--bg-1)' }}>
      <div style={{ background: 'var(--bg-2)', padding: '10px 14px', borderBottom: '1px solid var(--border-1)', display: 'flex', gap: 6, alignItems: 'center' }}>
        <span style={{ width: 10, height: 10, borderRadius: 999, background: 'var(--stone-300)' }}></span>
        <span style={{ width: 10, height: 10, borderRadius: 999, background: 'var(--stone-300)' }}></span>
        <span style={{ width: 10, height: 10, borderRadius: 999, background: 'var(--stone-300)' }}></span>
        <span style={{ marginLeft: 12, fontFamily: 'var(--font-mono)', fontSize: 12, color: 'var(--fg-3)' }}>{title}</span>
      </div>
      <pre style={{ margin: 0, padding: 16, background: 'var(--bg-1)', fontFamily: 'var(--font-mono)', fontSize: 13, lineHeight: 1.7, color: 'var(--fg-1)', minHeight, whiteSpace: 'pre-wrap', borderRadius: 0 }}>
        {lines.slice(0, visibleCount).map(renderLine)}
        <span style={{ display: 'inline-block', width: 8, height: 14, background: 'var(--color-lime)', marginLeft: 6, verticalAlign: -2, animation: 'devenvBlink 1s steps(2) infinite' }}></span>
      </pre>
      <style>{`@keyframes devenvBlink { 0%,49%{opacity:1} 50%,100%{opacity:0} }`}</style>
    </div>
  );
}

Object.assign(window, { Terminal });
