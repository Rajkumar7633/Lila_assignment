import { useEffect, useState } from 'react';
import { getDates, getMaps, getMatchJourney, getMatches } from '../api/client';
import type { MapInfo, MatchJourney, MatchSummary } from '../types';
import { useFilters } from '../context/FilterContext';

export function useCatalog() {
  const [maps, setMaps] = useState<MapInfo[]>([]);
  const [dates, setDates] = useState<string[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    let cancelled = false;
    (async () => {
      try {
        const [m, d] = await Promise.all([getMaps(), getDates()]);
        if (!cancelled) {
          setMaps(m);
          setDates(d);
        }
      } catch (e) {
        if (!cancelled) setError(e instanceof Error ? e.message : 'Failed to load catalog');
      } finally {
        if (!cancelled) setLoading(false);
      }
    })();
    return () => {
      cancelled = true;
    };
  }, []);

  return { maps, dates, loading, error };
}

export function useMatchList() {
  const { mapId, dateLabel } = useFilters();
  const [matches, setMatches] = useState<MatchSummary[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!mapId) {
      setMatches([]);
      return;
    }
    let cancelled = false;
    setLoading(true);
    getMatches(mapId, dateLabel)
      .then((data) => {
        if (!cancelled) setMatches(data);
      })
      .catch((e) => {
        if (!cancelled) setError(e instanceof Error ? e.message : 'Failed to load matches');
      })
      .finally(() => {
        if (!cancelled) setLoading(false);
      });
    return () => {
      cancelled = true;
    };
  }, [mapId, dateLabel]);

  return { matches, loading, error };
}

export function useMatchData() {
  const { matchId } = useFilters();
  const [journey, setJourney] = useState<MatchJourney | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!matchId) {
      setJourney(null);
      return;
    }
    let cancelled = false;
    setLoading(true);
    getMatchJourney(matchId)
      .then((data) => {
        if (!cancelled) setJourney(data);
      })
      .catch((e) => {
        if (!cancelled) setError(e instanceof Error ? e.message : 'Failed to load journey');
      })
      .finally(() => {
        if (!cancelled) setLoading(false);
      });
    return () => {
      cancelled = true;
    };
  }, [matchId]);

  return { journey, loading, error };
}
