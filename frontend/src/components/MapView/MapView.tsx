import { useMemo } from 'react';
import { MapContainer, ImageOverlay } from 'react-leaflet';
import L from 'leaflet';
import 'leaflet/dist/leaflet.css';

import { minimapUrl } from '../../api/client';
import type { MapInfo, MatchJourney } from '../../types';
import { useFilters } from '../../context/FilterContext';
import { useHeatmapData } from '../../hooks/useHeatmapData';
import { filterEventsByTime } from '../../hooks/usePlaybackClock';
import { PlayerPathLayer } from './PlayerPath';
import { EventMarker } from './EventMarker';
import { HeatmapLayer } from './HeatmapLayer';

const MAP_SIZE = 1024;
const bounds: L.LatLngBoundsExpression = [[0, 0], [MAP_SIZE, MAP_SIZE]];

interface MapViewProps {
  mapInfo: MapInfo | null;
  journey: MatchJourney | null;
  currentTs: number;
  playbackActive: boolean;
}

export function MapView({ mapInfo, journey, currentTs, playbackActive }: MapViewProps) {
  const { showHumans, showBots, heatmapType } = useFilters();
  const { points: heatPoints } = useHeatmapData();

  const visiblePaths = useMemo(() => {
    if (!journey) return [];
    return journey.paths.filter((p) => (p.isBot ? showBots : showHumans));
  }, [journey, showHumans, showBots]);

  const visibleEvents = useMemo(() => {
    if (!journey) return [];
    const events = journey.events.filter((e) => (e.isBot ? showBots : showHumans));
    return playbackActive ? filterEventsByTime(events, currentTs) : events;
  }, [journey, showHumans, showBots, playbackActive, currentTs]);

  if (!mapInfo) {
    return <div className="map-placeholder">Select a map to begin</div>;
  }

  const imageUrl = minimapUrl(mapInfo.imageUrl);

  return (
    <div className="map-shell">
      <MapContainer
        crs={L.CRS.Simple}
        bounds={bounds}
        maxBounds={bounds}
        style={{ height: '100%', width: '100%', background: '#0b1220' }}
        zoomControl
        attributionControl={false}
      >
        <ImageOverlay url={imageUrl} bounds={bounds} opacity={0.92} />
        <HeatmapLayer points={heatPoints} active={heatmapType != null} />
        {visiblePaths.map((path) => (
          <PlayerPathLayer
            key={path.userId}
            path={path}
            currentTs={currentTs}
            playbackActive={playbackActive}
          />
        ))}
        {visibleEvents.map((event, i) => (
          <EventMarker key={`${event.userId}-${event.ts}-${event.eventType}-${i}`} event={event} />
        ))}
      </MapContainer>
    </div>
  );
}
