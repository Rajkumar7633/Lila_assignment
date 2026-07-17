export interface MapInfo {
  id: string;
  label: string;
  scale: number;
  originX: number;
  originZ: number;
  imageUrl: string;
  imageWidth: number;
  imageHeight: number;
}

export interface MatchSummary {
  matchId: string;
  mapId: string;
  dateLabel: string;
  playerCount: number;
  humanCount: number;
  botCount: number;
  startTs: number;
  endTs: number;
  durationMs: number;
}

export interface Point {
  x: number;
  y: number;
}

export interface GameEvent {
  userId: string;
  isBot: boolean;
  eventType: string;
  category: 'movement' | 'kill' | 'death' | 'storm' | 'loot' | 'other';
  ts: number;
  position: Point;
}

export interface PlayerPath {
  userId: string;
  isBot: boolean;
  points: Point[];
  times: number[];
}

export interface MatchJourney {
  matchId: string;
  mapId: string;
  dateLabel: string;
  startTs: number;
  endTs: number;
  humanCount: number;
  botCount: number;
  eventCounts: EventCounts;
  paths: PlayerPath[];
  events: GameEvent[];
}

export interface EventCounts {
  kills: number;
  deaths: number;
  loot: number;
  storm: number;
}

export interface HeatmapPoint {
  x: number;
  y: number;
  intensity: number;
}

export type HeatmapType = 'kills' | 'deaths' | 'traffic';

export interface FilterState {
  mapId: string;
  dateLabel: string;
  matchId: string;
}
