import { useEffect, useMemo } from 'react';
import { FilterProvider, useFilters } from './context/FilterContext';
import { useCatalog, useMatchData, useMatchList } from './hooks/useMatchData';
import { usePlaybackClock } from './hooks/usePlaybackClock';
import { MapView } from './components/MapView/MapView';
import {
  DatePicker,
  HeatmapToggle,
  MapSelector,
  MatchSelector,
  PlayerToggle,
} from './components/Filters/Filters';
import { PlaybackControls } from './components/Timeline/PlaybackControls';
import { EventLegend } from './components/Legend/EventLegend';
import './App.css';

function Dashboard() {
  const { mapId, setMapId } = useFilters();
  const { maps, dates, loading: catalogLoading, error: catalogError } = useCatalog();
  const { matches, loading: matchesLoading } = useMatchList();
  const { journey, loading: journeyLoading } = useMatchData();
  const playback = usePlaybackClock(journey);

  useEffect(() => {
    if (!mapId && maps.length > 0) {
      setMapId(maps[0].id);
    }
  }, [maps, mapId, setMapId]);

  const activeMap = useMemo(
    () => maps.find((m) => m.id === (journey?.mapId ?? mapId)) ?? maps.find((m) => m.id === mapId) ?? null,
    [maps, mapId, journey],
  );

  return (
    <div className="app">
      <header className="header">
        <div>
          <p className="eyebrow">LILA BLACK · Level Design</p>
          <h1>Player Journey Visualization</h1>
        </div>
        <p className="subtitle">Explore movement, combat, loot, and storm deaths on production telemetry.</p>
      </header>

      {catalogError && <div className="banner error">{catalogError}</div>}

      <section className="toolbar">
        <MapSelector maps={maps} />
        <DatePicker dates={dates} />
        <MatchSelector matches={matches} loading={matchesLoading || catalogLoading} />
        <PlayerToggle />
        <HeatmapToggle />
      </section>

      <main className="main-grid">
        <MapView
          mapInfo={activeMap}
          journey={journey}
          currentTs={playback.currentTs}
          playbackActive={Boolean(journey)}
        />
        <aside className="sidebar">
          <EventLegend />
          <div className="stats-card">
            <h3>Selection</h3>
            {journeyLoading && <p className="muted">Loading journey…</p>}
            {journey && (
              <ul className="stats-list">
                <li><strong>Map</strong> {journey.mapId}</li>
                <li><strong>Humans</strong> {journey.humanCount}</li>
                <li><strong>Bots</strong> {journey.botCount}</li>
                <li><strong>Kills</strong> {journey.eventCounts?.kills ?? 0}</li>
                <li><strong>Deaths</strong> {journey.eventCounts?.deaths ?? 0}</li>
                <li><strong>Loot</strong> {journey.eventCounts?.loot ?? 0}</li>
                <li><strong>Storm</strong> {journey.eventCounts?.storm ?? 0}</li>
                <li><strong>In-game Duration</strong> {Math.floor(playback.dataSpan / 60)}m {Math.floor(playback.dataSpan % 60)}s</li>
              </ul>
            )}
            {!journey && !journeyLoading && (
              <p className="muted">Pick a match to render paths and events. Heatmaps work at map/date scope without a match.</p>
            )}
          </div>
        </aside>
      </main>

      <PlaybackControls
        playing={playback.playing}
        progress={playback.progress}
        label={playback.playbackLabel}
        disabled={!journey}
        onPlay={playback.play}
        onPause={playback.pause}
        onScrub={playback.scrub}
      />
    </div>
  );
}

export default function App() {
  return (
    <FilterProvider>
      <Dashboard />
    </FilterProvider>
  );
}
