import { useCallback, useEffect, useRef, useState } from 'react';
import type { MatchJourney } from '../types';

/** Wall-clock duration for a full match playback (level designers scrub ~30s, not 0.7s). */
const PLAYBACK_WALL_MS = 30_000;

const formatTime = (sec: number): string => {
  const m = Math.floor(sec / 60);
  const s = Math.floor(sec % 60);
  return `${m.toString().padStart(2, '0')}:${s.toString().padStart(2, '0')}`;
};

export function usePlaybackClock(journey: MatchJourney | null) {
  const [progress, setProgress] = useState(0);
  const [playing, setPlaying] = useState(false);
  const rafRef = useRef<number | null>(null);
  const lastFrameRef = useRef<number | null>(null);

  const startTs = journey?.startTs ?? 0;
  const endTs = journey?.endTs ?? 0;
  const dataSpan = Math.max(endTs - startTs, 1);

  useEffect(() => {
    setProgress(0);
    setPlaying(false);
  }, [journey?.matchId]);

  useEffect(() => {
    if (!playing || !journey) return;

    const tick = (now: number) => {
      if (lastFrameRef.current != null) {
        const delta = now - lastFrameRef.current;
        setProgress((prev) => {
          const next = prev + delta / PLAYBACK_WALL_MS;
          if (next >= 1) {
            setPlaying(false);
            return 1;
          }
          return next;
        });
      }
      lastFrameRef.current = now;
      rafRef.current = requestAnimationFrame(tick);
    };

    lastFrameRef.current = null;
    rafRef.current = requestAnimationFrame(tick);
    return () => {
      if (rafRef.current) cancelAnimationFrame(rafRef.current);
    };
  }, [playing, journey]);

  const currentTs = startTs + progress * dataSpan;

  const play = useCallback(() => {
    if (progress >= 1) setProgress(0);
    setPlaying(true);
  }, [progress]);

  const pause = useCallback(() => setPlaying(false), []);

  const scrub = useCallback((pct: number) => {
    setPlaying(false);
    setProgress(Math.min(Math.max(pct, 0), 1));
  }, []);

  return {
    currentTs,
    startTs,
    endTs,
    dataSpan,
    progress,
    playing,
    play,
    pause,
    scrub,
    playbackLabel: `Progress: ${Math.round(progress * 100)}% · In-game: ${formatTime(progress * dataSpan)} / ${formatTime(dataSpan)} · Scrub: ${(progress * (PLAYBACK_WALL_MS / 1000)).toFixed(1)}s / ${(PLAYBACK_WALL_MS / 1000).toFixed(0)}s`,
  };
}

export function filterPathByTime(
  points: { x: number; y: number }[],
  times: number[],
  currentTs: number,
) {
  const idx = times.findIndex((t) => t > currentTs);
  const end = idx === -1 ? points.length : Math.max(idx, 1);
  return points.slice(0, end);
}

export function filterEventsByTime<T extends { ts: number }>(events: T[], currentTs: number) {
  return events.filter((e) => e.ts <= currentTs);
}
