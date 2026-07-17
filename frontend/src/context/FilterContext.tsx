import { createContext, useContext, useMemo, useState, type ReactNode } from 'react';
import type { FilterState, HeatmapType } from '../types';

interface FilterContextValue extends FilterState {
  heatmapType: HeatmapType | null;
  showHumans: boolean;
  showBots: boolean;
  setMapId: (mapId: string) => void;
  setDateLabel: (date: string) => void;
  setMatchId: (matchId: string) => void;
  setHeatmapType: (type: HeatmapType | null) => void;
  setShowHumans: (v: boolean) => void;
  setShowBots: (v: boolean) => void;
}

const FilterContext = createContext<FilterContextValue | null>(null);

export function FilterProvider({ children }: { children: ReactNode }) {
  const [mapId, setMapId] = useState('');
  const [dateLabel, setDateLabel] = useState('');
  const [matchId, setMatchId] = useState('');
  const [heatmapType, setHeatmapType] = useState<HeatmapType | null>(null);
  const [showHumans, setShowHumans] = useState(true);
  const [showBots, setShowBots] = useState(true);

  const value = useMemo(
    () => ({
      mapId,
      dateLabel,
      matchId,
      heatmapType,
      showHumans,
      showBots,
      setMapId,
      setDateLabel,
      setMatchId,
      setHeatmapType,
      setShowHumans,
      setShowBots,
    }),
    [mapId, dateLabel, matchId, heatmapType, showHumans, showBots],
  );

  return <FilterContext.Provider value={value}>{children}</FilterContext.Provider>;
}

export function useFilters() {
  const ctx = useContext(FilterContext);
  if (!ctx) throw new Error('useFilters must be used within FilterProvider');
  return ctx;
}
