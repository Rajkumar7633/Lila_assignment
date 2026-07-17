import { useEffect, useState } from 'react';
import { getHeatmap } from '../api/client';
import type { HeatmapPoint, HeatmapType } from '../types';
import { useFilters } from '../context/FilterContext';

export function useHeatmapData() {
  const { mapId, dateLabel, matchId, heatmapType } = useFilters();
  const [points, setPoints] = useState<HeatmapPoint[]>([]);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    if (!mapId || !heatmapType) {
      setPoints([]);
      return;
    }
    let cancelled = false;
    setLoading(true);
    getHeatmap(mapId, dateLabel, matchId, heatmapType)
      .then((data) => {
        if (!cancelled) setPoints(data);
      })
      .catch(() => {
        if (!cancelled) setPoints([]);
      })
      .finally(() => {
        if (!cancelled) setLoading(false);
      });
    return () => {
      cancelled = true;
    };
  }, [mapId, dateLabel, matchId, heatmapType]);

  return { points, loading, heatmapType: heatmapType as HeatmapType | null };
}
