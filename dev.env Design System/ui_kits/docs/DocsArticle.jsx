// DocsArticle — the rendered article body for the docs reader. Uses .prose styles.
function Callout({ tone = 'warn', children }) {
  const borderColor = tone === 'warn' ? 'var(--color-warning)' : 'var(--color-info)';
  return (
    <blockquote style={{ borderLeftColor: borderColor }}>{children}</blockquote>
  );
}

function DocsArticle() {
  return (
    <article className="prose" style={{ flex: 1, minWidth: 0, paddingTop: 8 }}>
      <Eyebrow>/ start / getting started</Eyebrow>
      <h1 id="getting-started" style={{ marginTop: 12 }}>Getting started</h1>
      <p>
        dev.env runs <strong style={{ color: 'var(--fg-1)' }}>one shared pool of containers</strong> for all your projects. Services that support multi-tenancy isolate projects using their own built-in mechanisms: a dedicated database per project, a separate Redis DB index, a scoped Elasticsearch index. dev.env provisions all of that automatically when you run <code>dev start</code>.
      </p>
      <p>
        This page walks you through installing the binary, initialising your first project, and starting a shared MySQL + Redis on demand.
      </p>

      <h2 id="installation">Installation</h2>
      <p>Install on macOS or Linux:</p>
      <pre><code>{`VERSION=$(curl -s https://api.github.com/repos/JCombee/dev.env/releases/latest \\
  | grep '"tag_name"' | cut -d'"' -f4) && \\
curl -sSL "https://github.com/JCombee/dev.env/releases/download/\${VERSION}/dev_\${VERSION#v}_$(uname -s | tr A-Z a-z)_$(uname -m | sed 's/x86_64/amd64/').tar.gz" \\
  | tar -xz && sudo mv dev /usr/local/bin/`}</code></pre>
      <p>That's it. dev.env is a single Go binary with no runtime dependencies — only Docker needs to be running.</p>

      <h2 id="init">Initialise a project</h2>
      <p>From inside your project's directory, run <code>dev init</code>. The interactive wizard detects your project type, asks which services you need, and writes a <code>.dev.env.yaml</code>.</p>
      <pre><code>{`$ dev init
Detected project type: laravel

Select services (space to toggle, enter to confirm):
  [x] mysql
  [x] redis
  [ ] elasticsearch
  [ ] rabbitmq

MySQL tag [latest]: 8.0
MySQL mode: shared or dedicated? [shared]:
Enable .env generation? [y/N]: y

Writing .dev.env.yaml... done`}</code></pre>

      <h2 id="config">Project configuration</h2>
      <p>The wizard writes a file that looks like this:</p>
      <pre><code>{`project: my-project
type: laravel
env: true
services:
  - image: mysql
    tag: 8.0
  - redis
  - image: elasticsearch
    tag: 8.11
    dedicated: true`}</code></pre>
      <p>Every field is documented in <a href="#config-reference">Field reference</a>. The short version:</p>
      <table>
        <thead>
          <tr><th>field</th><th>required</th><th>what it does</th></tr>
        </thead>
        <tbody>
          <tr><td><code>project</code></td><td>yes</td><td>Unique identifier. Used as the database name, Redis prefix, etc.</td></tr>
          <tr><td><code>services</code></td><td>yes</td><td>List of services. Each can be a string or an object with <code>image</code>, <code>tag</code>, and <code>dedicated</code>.</td></tr>
          <tr><td><code>env</code></td><td>no</td><td>If <code>true</code>, write a <code>.env</code> file on <code>dev start</code>. Default: <code>false</code>.</td></tr>
          <tr><td><code>type</code></td><td>no</td><td>Project type. Used by <code>dev init</code> for auto-detection only.</td></tr>
        </tbody>
      </table>

      <h2 id="start">Start the services</h2>
      <p>Run <code>dev start</code> and dev.env brings everything up:</p>
      <pre><code>{`$ dev start
✓ mysql:8.0              already running (shared)
✓ redis:latest           started (shared)
✓ elasticsearch:8.11     started (dedicated)
✓ Database "my-project" created
✓ .env written`}</code></pre>
      <p>
        Shared containers already running for another project are reused. Dedicated containers are always project-private and never shared. On first run, dev.env provisions project-specific resources (creates DB, user, index, etc.). On subsequent runs it checks that everything is still in order.
      </p>

      <Callout tone="warn">
        <p style={{ margin: 0 }}><strong style={{ color: 'var(--fg-1)' }}>Heads up.</strong> Per-project port overrides force a service to run as a dedicated container, bypassing the shared pool dev.env is optimised for. Prefer a global override in <code>~/.dev.env/settings.yaml</code>.</p>
      </Callout>

      <h2 id="stop">Stop when you're done</h2>
      <p>Run <code>dev stop</code> and dev.env shuts down containers <em style={{ color: 'var(--fg-1)' }}>only if no other active project depends on them</em>:</p>
      <pre><code>{`$ dev stop
✓ redis:latest           stopped
✓ elasticsearch:8.11     stopped (dedicated)
~ mysql:8.0              kept running (in use by: api-project)`}</code></pre>
      <p>The <code>~</code> glyph means the container is still running because another project needs it. dev.env tracks every project's state in <code>~/.dev.env/projects/&lt;name&gt;/state.yaml</code> and consults all of them before stopping anything shared.</p>

      <h2 id="next">Where next</h2>
      <ul>
        <li><a href="#ports">Port allocation</a> — how dev.env keeps multiple versions of the same service from colliding</li>
        <li><a href="#env-gen">.env generation</a> — what variables each service writes and how to rename them</li>
        <li><a href="#local">Local overrides</a> — using <code>.dev.env.local.yaml</code> for per-developer customisation</li>
        <li><a href="reference.html">CLI reference</a> — every command and flag</li>
      </ul>
    </article>
  );
}

Object.assign(window, { DocsArticle, Callout });
