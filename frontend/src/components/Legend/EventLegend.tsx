import { EVENT_STYLE } from '../MapView/EventMarker';

function ShapeSwatch({ category }: { category: string }) {
  const style = EVENT_STYLE[category];
  if (!style) return null;

  if (style.shape === 'triangle') {
    return (
      <span className="swatch-shape">
        <span className="shape-triangle" style={{ borderBottomColor: style.color }} />
      </span>
    );
  }
  if (style.shape === 'cross') {
    return <span className="swatch-shape shape-cross" style={{ color: style.color }}>✕</span>;
  }
  if (style.shape === 'diamond') {
    return <span className="swatch-shape shape-diamond" style={{ background: style.color }} />;
  }
  return <span className="swatch" style={{ background: style.color }} />;
}

export function EventLegend() {
  return (
    <div className="legend">
      <h3>Legend</h3>
      <div className="legend-row">
        <span className="swatch path-human" /> Human path (solid blue)
      </div>
      <div className="legend-row">
        <span className="swatch path-bot" /> Bot path (dashed gray)
      </div>
      {Object.entries(EVENT_STYLE).map(([key, style]) => (
        <div className="legend-row" key={key}>
          <ShapeSwatch category={key} /> {style.label}
        </div>
      ))}
    </div>
  );
}
