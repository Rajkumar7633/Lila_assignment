# Data Insights

Three patterns discovered by exploring the telemetry with this tool.

---

## 1. Combat is overwhelmingly bot-driven, not human PvP

### What caught my eye
On the kill/death heatmap, hotspots lit up everywhere — but filtering to human-only `Kill`/`Killed` events showed almost nothing (6 events total vs 3,115 bot combat events).

### Evidence
| Event type | Count |
|------------|------:|
| BotKill + BotKilled | 3,115 |
| Kill + Killed (human PvP) | 6 |
| KilledByStorm | 39 |

Human PvP accounts for **0.19%** of all combat events in the dataset.

### Actionable recommendation
- **Reduce bot density** in zones where designers want human engagement, or add objectives that force human confrontation (shared extraction, contested loot rooms).
- **Metrics affected:** human PvP encounters per match, average time-to-first-human-fight, player retention in mid-game.
- Add a dashboard filter (already partially supported) to compare bot vs human kill heatmaps side-by-side when tuning spawn tables.

### Why a level designer should care
If a zone looks "hot" on the aggregate heatmap but is actually bot-only, you may over-invest in cover/layout for fights that players never experience against each other. Separating bot and human combat layers prevents false confidence in encounter design.

---

## 2. Loot on Ambrose Valley is hyper-concentrated in one region

### What caught my eye
The loot heatmap on Ambrose Valley showed a single bright cluster; most of the map was comparatively cold despite being the most-played map (566 of 796 matches).

### Evidence
- Total loot events on Ambrose Valley: **9,955**
- Top 90×90 world-unit grid cell: **1,317 loot pickups (13.2% of all loot on the map)**
- Next-best cells: ~530–610 each — steep drop-off

Using README coordinates, the top cell centers around world `(x≈80, z≈-113)` — a small fraction of the 1024×1024 minimap area.

### Actionable recommendation
- **Redistribute loot spawns** or add secondary high-value loot routes in under-visited quadrants (visible as cold zones on the traffic heatmap).
- **Metrics affected:** map coverage %, average player path length, extraction success rate, time spent in low-loot areas.
- Place a new objective (key card, extraction unlock) in a cold quadrant to pull traffic — validate with the traffic overlay after the next data pull.

### Why a level designer should care
When 13% of all loot sits in one grid cell, pathing becomes predictable — players funnel, fights stack, and large map art goes unseen. The tool makes this visible in seconds vs digging through raw parquet.

---

## 3. Storm deaths cluster on smaller maps — extraction paths need work

### What caught my eye
Storm death markers (`KilledByStorm`) appeared evenly distributed by raw count, but per-map rates tell a different story given map size and match volume.

### Evidence
| Map | Storm deaths | Matches | Deaths / match |
|-----|-------------:|--------:|---------------:|
| AmbroseValley | 17 | 566 | 0.03 |
| Lockdown | 17 | 171 | **0.10** |
| GrandRift | 5 | 59 | 0.08 |

Lockdown has **3× the storm death rate per match** vs Ambrose Valley despite being the "close-quarters" map where players should reach extraction faster.

*(Note on GrandRift: While GrandRift also exhibits a high storm death rate of 0.08 per match, its sample size is very small — only 59 matches and 5 storm deaths. This makes Lockdown's rate of 0.10 across 171 matches the most statistically significant and reliable indicator of a systematic level design issue, though the GrandRift trend suggests smaller maps generally warrant a design review).*

Average storm-death positions on Lockdown cluster toward map edges (avg x≈−120, z≈−180 in world space vs broader spread on Ambrose).

### Actionable recommendation
- **Review extraction point placement** on Lockdown — add a mid-map fallback extract or widen safe corridors ahead of the one-directional storm.
- **Metrics affected:** storm death rate, extraction success %, average match survival time.
- Use timeline playback on Lockdown matches with storm deaths to see if players die while rotating vs while extracting — pinpoints whether the fix is path width or extract timing.

### Why a level designer should care
Storm deaths are the clearest signal that layout + timing failed the player. A high per-match storm rate on your smallest map suggests the storm/extract geometry is punishing before players can learn routes — fixable with data-backed iteration instead of guesswork.

---

*All stats computed from the provided Feb 10–14, 2026 production dataset (89,004 event rows, 796 matches).*
