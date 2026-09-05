// Shared primitives: Eyebrow, Badge, Kbd, Button
const { useState, useEffect, useRef } = React;

function Eyebrow({ children, style }) {
  return <span className="eyebrow" style={style}>{children}</span>;
}

function Badge({ children, variant = 'default', style }) {
  return <span className={`badge ${variant !== 'default' ? `badge-${variant}` : ''}`} style={style}>{children}</span>;
}

function Kbd({ children }) {
  return <span className="kbd">{children}</span>;
}

function Button({ children, variant = 'primary', as = 'button', href, onClick, style, arrow }) {
  const Tag = as;
  const cls = `btn btn-${variant}`;
  const inner = <>{children}{arrow ? <span className="arrow">→</span> : null}</>;
  if (as === 'a') return <a className={cls} href={href} onClick={onClick} style={style}>{inner}</a>;
  return <Tag className={cls} onClick={onClick} style={style}>{inner}</Tag>;
}

function Logomark({ size = 22, dark = false }) {
  const stroke = dark ? '#F7F6F2' : '#0E1116';
  return (
    <svg width={size} height={size} viewBox="0 0 64 64" style={{ flexShrink: 0 }}>
      <rect x="8" y="10" width="48" height="10" rx="2" fill="none" stroke={stroke} strokeWidth="3"/>
      <rect x="8" y="22" width="48" height="10" rx="2" fill="#C6F432" stroke="#C6F432" strokeWidth="3"/>
      <rect x="8" y="34" width="48" height="10" rx="2" fill="none" stroke={stroke} strokeWidth="3"/>
      <rect x="8" y="46" width="48" height="10" rx="2" fill="none" stroke={stroke} strokeWidth="3"/>
    </svg>
  );
}

function Wordmark({ height = 22 }) {
  return (
    <span style={{ display: 'inline-flex', alignItems: 'baseline', fontFamily: 'var(--font-mono)', fontWeight: 600, fontSize: height, letterSpacing: '-0.02em', color: 'var(--fg-1)', lineHeight: 1 }}>
      <span>dev</span>
      <span style={{ display: 'inline-block', width: height * 0.22, height: height * 0.22, background: 'var(--color-lime)', margin: `0 ${height * 0.04}px ${height * 0.05}px`, transform: 'translateY(0)' }}></span>
      <span>env</span>
    </span>
  );
}

Object.assign(window, { Eyebrow, Badge, Kbd, Button, Logomark, Wordmark });
