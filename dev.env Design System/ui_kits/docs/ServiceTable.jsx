// ServiceTable — supported-services table with monogram chips, mode badges, port hint.
function Monogram({ letters }) {
  return (
    <span style={{ display: 'inline-flex', width: 32, height: 32, borderRadius: 6, background: 'var(--bg-2)', border: '1px solid var(--border-2)', alignItems: 'center', justifyContent: 'center', fontFamily: 'var(--font-mono)', fontWeight: 600, fontSize: 12, color: 'var(--fg-1)', letterSpacing: '-0.02em' }}>
      {letters}
    </span>
  );
}

function ServiceTable() {
  const rows = [
    { cat: 'databases', svc: 'mysql',         m: 'MY', port: '3306', mode: 'shared',   note: 'separate database + user per project' },
    { cat: 'databases', svc: 'mariadb',       m: 'MA', port: '3306', mode: 'shared',   note: 'separate database + user per project' },
    { cat: 'databases', svc: 'postgres',      m: 'PG', port: '5432', mode: 'shared',   note: 'separate database + user per project' },
    { cat: 'databases', svc: 'mongo',         m: 'MO', port: '27017', mode: 'shared',  note: 'separate database per project' },
    { cat: 'databases', svc: 'cassandra',     m: 'CA', port: '9042', mode: 'shared',   note: 'separate keyspace per project' },
    { cat: 'search',    svc: 'elasticsearch', m: 'ES', port: '9200', mode: 'shared',   note: 'separate index per project' },
    { cat: 'search',    svc: 'opensearch',    m: 'OS', port: '9200', mode: 'shared',   note: 'separate index per project' },
    { cat: 'search',    svc: 'meilisearch',   m: 'MS', port: '7700', mode: 'shared',   note: 'separate index per project' },
    { cat: 'search',    svc: 'typesense',     m: 'TS', port: '8108', mode: 'shared',   note: 'separate collection per project' },
    { cat: 'search',    svc: 'solr',          m: 'SO', port: '8983', mode: 'shared',   note: 'separate core per project' },
    { cat: 'caches',    svc: 'redis',         m: 'RD', port: '6379', mode: 'shared',   note: 'separate DB index per project' },
    { cat: 'caches',    svc: 'memcached',     m: 'MC', port: '11211', mode: 'dedicated', note: 'no native namespacing' },
    { cat: 'queues',    svc: 'rabbitmq',      m: 'RQ', port: '5672', mode: 'shared',   note: 'separate vhost per project' },
    { cat: 'queues',    svc: 'kafka',         m: 'KF', port: '9092', mode: 'shared',   note: 'separate topic prefix per project' },
    { cat: 'storage',   svc: 'minio',         m: 'S3', port: '9000', mode: 'shared',   note: 'separate bucket per project' },
    { cat: 'mail',      svc: 'mailpit',       m: 'MP', port: '1025', mode: 'shared',   note: 'single shared inbox by design' },
    { cat: 'mail',      svc: 'mailhog',       m: 'MH', port: '1025', mode: 'shared',   note: 'single shared inbox by design' },
  ];
  const [filter, setFilter] = useState('');
  const filtered = rows.filter(r =>
    !filter || r.svc.toLowerCase().includes(filter.toLowerCase()) || r.cat.includes(filter.toLowerCase())
  );

  return (
    <section style={{ paddingTop: 96, paddingBottom: 96, borderTop: '1px solid var(--border-1)' }}>
      <div className="container">
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-end', marginBottom: 32, gap: 24, flexWrap: 'wrap' }}>
          <div>
            <Eyebrow>/ 04 — services</Eyebrow>
            <h2 style={{ fontFamily: 'var(--font-sans)', fontWeight: 600, fontSize: 36, lineHeight: 1.1, letterSpacing: '-0.015em', color: 'var(--fg-1)', margin: '12px 0 8px' }}>
              17 services supported. Any Docker image, too.
            </h2>
            <p style={{ fontSize: 16, color: 'var(--fg-2)', margin: 0, maxWidth: 600 }}>
              Supported services run shared (with automatic project isolation) or dedicated. Unsupported images still work — they just won't get auto-provisioning.
            </p>
          </div>
          <input
            value={filter}
            onChange={e => setFilter(e.target.value)}
            placeholder="filter…"
            style={{ fontFamily: 'var(--font-mono)', fontSize: 13, padding: '8px 12px', border: '1px solid var(--border-2)', borderRadius: 6, background: 'var(--bg-1)', color: 'var(--fg-1)', minWidth: 180, outline: 'none' }}
          />
        </div>
        <div style={{ border: '1px solid var(--border-2)', borderRadius: 6, overflow: 'hidden' }}>
          <div style={{ display: 'grid', gridTemplateColumns: '1.4fr 1fr 1.4fr 2fr', padding: '12px 20px', background: 'var(--bg-2)', borderBottom: '1px solid var(--border-1)', fontFamily: 'var(--font-mono)', fontSize: 11, textTransform: 'uppercase', letterSpacing: '0.06em', color: 'var(--fg-3)' }}>
            <span>service</span><span>port</span><span>mode</span><span>isolation</span>
          </div>
          {filtered.map((r, i) => (
            <div key={i} style={{ display: 'grid', gridTemplateColumns: '1.4fr 1fr 1.4fr 2fr', padding: '12px 20px', borderBottom: i === filtered.length - 1 ? 'none' : '1px solid var(--border-1)', alignItems: 'center', background: 'var(--bg-1)' }}>
              <span style={{ display: 'flex', alignItems: 'center', gap: 10, fontFamily: 'var(--font-mono)', fontSize: 14, color: 'var(--fg-1)' }}>
                <Monogram letters={r.m}/>
                {r.svc}
              </span>
              <span style={{ fontFamily: 'var(--font-mono)', fontSize: 13, color: 'var(--fg-2)' }}>{r.port}</span>
              <span>
                {r.mode === 'shared'
                  ? <Badge variant="shared">shared</Badge>
                  : <Badge variant="dedicated">dedicated only</Badge>}
              </span>
              <span style={{ fontSize: 13, color: 'var(--fg-2)' }}>{r.note}</span>
            </div>
          ))}
          {filtered.length === 0 && (
            <div style={{ padding: 32, textAlign: 'center', color: 'var(--fg-3)', fontFamily: 'var(--font-mono)', fontSize: 13 }}>no matches</div>
          )}
        </div>
      </div>
    </section>
  );
}

Object.assign(window, { ServiceTable, Monogram });
