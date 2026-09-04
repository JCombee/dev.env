// InstallBlock — tabbed installer with copy button.
function InstallBlock() {
  const [tab, setTab] = useState('macos');
  const [copied, setCopied] = useState(false);
  const snippets = {
    macos: `VERSION=$(curl -s https://api.github.com/repos/JCombee/dev.env/releases/latest | grep '"tag_name"' | cut -d'"' -f4) && \\
OS=$(uname -s | tr '[:upper:]' '[:lower:]') && \\
ARCH=$(uname -m | sed 's/x86_64/amd64/;s/aarch64/arm64/') && \\
curl -sSL "https://github.com/JCombee/dev.env/releases/download/\${VERSION}/dev_\${VERSION#v}_\${OS}_\${ARCH}.tar.gz" | tar -xz && \\
sudo mv dev /usr/local/bin/`,
    linux: `VERSION=$(curl -s https://api.github.com/repos/JCombee/dev.env/releases/latest | grep '"tag_name"' | cut -d'"' -f4) && \\
OS=$(uname -s | tr '[:upper:]' '[:lower:]') && \\
ARCH=$(uname -m | sed 's/x86_64/amd64/;s/aarch64/arm64/') && \\
curl -sSL "https://github.com/JCombee/dev.env/releases/download/\${VERSION}/dev_\${VERSION#v}_\${OS}_\${ARCH}.tar.gz" | tar -xz && \\
sudo mv dev /usr/local/bin/`,
    windows: `$release = Invoke-RestMethod https://api.github.com/repos/JCombee/dev.env/releases/latest
$asset = $release.assets | Where-Object { $_.name -like "*Windows_x86_64.zip" }
Invoke-WebRequest $asset.browser_download_url -OutFile dev.zip
Expand-Archive dev.zip -DestinationPath $env:LOCALAPPDATA\\dev
$env:PATH += ";$env:LOCALAPPDATA\\dev"`,
  };
  const tabs = [
    { id: 'macos',   label: 'macOS' },
    { id: 'linux',   label: 'Linux' },
    { id: 'windows', label: 'Windows' },
  ];
  const copy = () => {
    navigator.clipboard?.writeText(snippets[tab]).then(() => {
      setCopied(true);
      setTimeout(() => setCopied(false), 1500);
    });
  };
  return (
    <section style={{ paddingTop: 96, paddingBottom: 96, borderTop: '1px solid var(--border-1)', background: 'var(--bg-2)' }}>
      <div className="container">
        <div style={{ display: 'grid', gridTemplateColumns: 'minmax(0,1fr) minmax(0,1.4fr)', gap: 64, alignItems: 'flex-start' }}>
          <div>
            <Eyebrow>/ 03 — install</Eyebrow>
            <h2 style={{ fontFamily: 'var(--font-sans)', fontWeight: 600, fontSize: 36, lineHeight: 1.1, letterSpacing: '-0.015em', color: 'var(--fg-1)', margin: '12px 0 16px' }}>
              One curl. One binary.
            </h2>
            <p style={{ fontSize: 16, lineHeight: 1.55, color: 'var(--fg-2)', maxWidth: 380, margin: '0 0 20px' }}>
              No runtime, no plugin manager, no docker-compose to maintain. dev.env ships as a single Go binary — install it once and forget it's there.
            </p>
            <div style={{ display: 'flex', gap: 8, flexWrap: 'wrap' }}>
              <Badge>amd64</Badge>
              <Badge>arm64</Badge>
              <Badge>requires Docker</Badge>
            </div>
          </div>
          <div>
            <div style={{ border: '1px solid var(--border-2)', borderRadius: 8, overflow: 'hidden', background: 'var(--bg-1)' }}>
              <div style={{ display: 'flex', borderBottom: '1px solid var(--border-1)', background: 'var(--bg-1)' }}>
                {tabs.map(t => (
                  <button key={t.id} onClick={() => setTab(t.id)} style={{
                    border: 0, background: 'none', cursor: 'pointer',
                    fontFamily: 'var(--font-mono)', fontWeight: 500, fontSize: 13,
                    padding: '12px 18px', color: tab === t.id ? 'var(--fg-1)' : 'var(--fg-3)',
                    borderBottom: tab === t.id ? '2px solid var(--color-lime)' : '2px solid transparent',
                    marginBottom: -1,
                  }}>{t.label}</button>
                ))}
                <div style={{ marginLeft: 'auto', display: 'flex', alignItems: 'center', paddingRight: 12 }}>
                  <button onClick={copy} style={{ fontFamily: 'var(--font-mono)', fontSize: 12, padding: '5px 10px', borderRadius: 4, border: '1px solid var(--border-2)', background: 'var(--bg-1)', color: 'var(--fg-1)', cursor: 'pointer' }}>
                    {copied ? '✓ copied' : 'copy'}
                  </button>
                </div>
              </div>
              <pre style={{ margin: 0, padding: 18, background: 'var(--bg-1)', fontFamily: 'var(--font-mono)', fontSize: 13, lineHeight: 1.6, color: 'var(--fg-1)', overflowX: 'auto' }}>{snippets[tab]}</pre>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}

Object.assign(window, { InstallBlock });
