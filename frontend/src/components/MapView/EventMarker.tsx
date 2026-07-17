import { Marker } from 'react-leaflet';
import L from 'leaflet';
import type { GameEvent } from '../../types';

export const EVENT_STYLE: Record<
  string,
  { color: string; label: string; shape: 'triangle' | 'cross' | 'diamond' | 'circle' }
> = {
  kill: { color: '#22c55e', label: 'Kill', shape: 'triangle' },
  death: { color: '#ef4444', label: 'Death', shape: 'cross' },
  storm: { color: '#a855f7', label: 'Storm death', shape: 'diamond' },
  loot: { color: '#f59e0b', label: 'Loot', shape: 'circle' },
};

function markerHTML(color: string, shape: string): string {
  const base = `display:flex;align-items:center;justify-content:center;width:18px;height:18px;`;
  switch (shape) {
    case 'triangle':
      return `<div style="${base}"><div style="width:0;height:0;border-left:7px solid transparent;border-right:7px solid transparent;border-bottom:12px solid ${color};filter:drop-shadow(0 0 2px #000)"></div></div>`;
    case 'cross':
      return `<div style="${base}font-size:16px;font-weight:900;color:${color};text-shadow:0 0 2px #000">✕</div>`;
    case 'diamond':
      return `<div style="${base}"><div style="width:11px;height:11px;background:${color};transform:rotate(45deg);border:1px solid #0f172a"></div></div>`;
    default:
      return `<div style="${base}"><div style="width:10px;height:10px;border-radius:50%;background:${color};border:2px solid #0f172a"></div></div>`;
  }
}

const iconCache = new Map<string, L.DivIcon>();

function getIcon(category: string): L.DivIcon {
  if (iconCache.has(category)) return iconCache.get(category)!;
  const style = EVENT_STYLE[category];
  const icon = L.divIcon({
    className: 'event-marker-icon',
    html: markerHTML(style.color, style.shape),
    iconSize: [18, 18],
    iconAnchor: [9, 9],
  });
  iconCache.set(category, icon);
  return icon;
}

interface EventMarkerProps {
  event: GameEvent;
}

export function EventMarker({ event }: EventMarkerProps) {
  const style = EVENT_STYLE[event.category];
  if (!style) return null;

  return (
    <Marker
      position={[event.position.y, event.position.x]}
      icon={getIcon(event.category)}
    />
  );
}
