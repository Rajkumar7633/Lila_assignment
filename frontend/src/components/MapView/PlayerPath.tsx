import { Polyline } from 'react-leaflet';
import type { PlayerPath } from '../../types';
import { filterPathByTime } from '../../hooks/usePlaybackClock';

interface PlayerPathLayerProps {
  path: PlayerPath;
  currentTs: number;
  playbackActive: boolean;
}

export function PlayerPathLayer({ path, currentTs, playbackActive }: PlayerPathLayerProps) {
  const visiblePoints = playbackActive ? filterPathByTime(path.points, path.times, currentTs) : path.points;
  if (visiblePoints.length < 2) return null;

  const latLngs = visiblePoints.map((p) => [p.y, p.x] as [number, number]);

  return (
    <Polyline
      positions={latLngs}
      pathOptions={{
        color: path.isBot ? '#94a3b8' : '#38bdf8',
        weight: path.isBot ? 2 : 3,
        opacity: path.isBot ? 0.55 : 0.85,
        dashArray: path.isBot ? '6 8' : undefined,
      }}
    />
  );
}
