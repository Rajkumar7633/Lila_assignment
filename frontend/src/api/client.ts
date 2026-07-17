import type { HeatmapPoint, HeatmapType, MapInfo, MatchJourney, MatchSummary } from '../types';

const API_BASE = import.meta.env.VITE_API_URL ?? 'http://localhost:8080';

async function fetchJSON<T>(path: string): Promise<T> {
  const res = await fetch(`${API_BASE}${path}`);
  if (!res.ok) {
    const body = await res.text();
    throw new Error(body || `Request failed: ${res.status}`);
  }
  return res.json() as Promise<T>;
}

export function getMaps() {
  return fetchJSON<MapInfo[]>('/api/maps');
}

export function getDates() {
  return fetchJSON<string[]>('/api/dates');
}

export function getMatches(mapId: string, dateLabel: string) {
  const params = new URLSearchParams();
  if (mapId) params.set('map', mapId);
  if (dateLabel) params.set('date', dateLabel);
  return fetchJSON<MatchSummary[]>(`/api/matches?${params}`);
}

export function getMatchJourney(matchId: string) {
  return fetchJSON<MatchJourney>(`/api/matches/${encodeURIComponent(matchId)}/journey`);
}

export function getHeatmap(mapId: string, dateLabel: string, matchId: string, type: HeatmapType) {
  const params = new URLSearchParams({ map: mapId, type });
  if (dateLabel) params.set('date', dateLabel);
  if (matchId) params.set('match', matchId);
  return fetchJSON<HeatmapPoint[]>(`/api/heatmap?${params}`);
}

export function minimapUrl(path: string) {
  if (path.startsWith('http')) return path;
  return `${API_BASE}${path}`;
}
