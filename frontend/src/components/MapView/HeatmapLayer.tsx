import { useEffect } from 'react';
import { useMap } from 'react-leaflet';
import L from 'leaflet';
import 'leaflet.heat';
import type { HeatmapPoint } from '../../types';

interface HeatmapLayerProps {
  points: HeatmapPoint[];
  active: boolean;
}

export function HeatmapLayer({ points, active }: HeatmapLayerProps) {
  const map = useMap();

  useEffect(() => {
    if (!active || points.length === 0) return;

    const data: [number, number, number][] = points.map((p) => [p.y, p.x, p.intensity]);
    const layer = L.heatLayer(data, {
      radius: 28,
      blur: 18,
      maxZoom: 2,
      minOpacity: 0.35,
      gradient: {
        0.2: '#1e3a5f',
        0.5: '#f59e0b',
        0.8: '#ef4444',
        1.0: '#fef08a',
      },
    });
    layer.addTo(map);
    return () => {
      map.removeLayer(layer);
    };
  }, [map, points, active]);

  return null;
}
