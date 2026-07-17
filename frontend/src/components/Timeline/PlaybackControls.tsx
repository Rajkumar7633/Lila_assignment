interface PlaybackControlsProps {
  playing: boolean;
  progress: number;
  label: string;
  disabled: boolean;
  onPlay: () => void;
  onPause: () => void;
  onScrub: (pct: number) => void;
}

export function PlaybackControls({
  playing,
  progress,
  label,
  disabled,
  onPlay,
  onPause,
  onScrub,
}: PlaybackControlsProps) {
  return (
    <div className="timeline">
      <div className="timeline-controls">
        <button type="button" disabled={disabled || playing} onClick={onPlay}>
          ▶ Play
        </button>
        <button type="button" disabled={disabled || !playing} onClick={onPause}>
          ⏸ Pause
        </button>
        <span className="timeline-hint">
          {disabled ? 'Select a match to scrub timeline' : label}
        </span>
      </div>
      <input
        type="range"
        min={0}
        max={1000}
        value={Math.round(progress * 1000)}
        disabled={disabled}
        onChange={(e) => onScrub(Number(e.target.value) / 1000)}
        className="timeline-slider"
      />
    </div>
  );
}
