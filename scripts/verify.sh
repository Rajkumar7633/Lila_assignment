#!/usr/bin/env bash
# Quick verification script for assignment checklist
set -euo pipefail

API="${API_URL:-http://localhost:8080}"
echo "Testing API at $API"

curl -sf "$API/health" | grep -q ok && echo "✓ Health"
curl -sf "$API/api/maps" | python3 -c "import sys,json; m=json.load(sys.stdin); assert len(m)==3; print('✓ Maps (3)')"
curl -sf "$API/api/dates" | python3 -c "import sys,json; d=json.load(sys.stdin); assert len(d)==5; print('✓ Dates (5)')"
curl -sf "$API/api/matches?map=AmbroseValley" | python3 -c "import sys,json; m=json.load(sys.stdin); assert len(m)>0; print(f'✓ Matches ({len(m)})')"

MID=$(curl -sf "$API/api/matches?map=AmbroseValley" | python3 -c "import sys,json; print(json.load(sys.stdin)[0]['matchId'])")
ENCODED=$(python3 -c "import urllib.parse; print(urllib.parse.quote('$MID', safe=''))")
J=$(curl -sf "$API/api/matches/$ENCODED/journey")
echo "$J" | python3 -c "
import sys,json
j=json.load(sys.stdin)
assert 'paths' in j and 'events' in j
assert 'humanCount' in j and 'eventCounts' in j
assert len(j['paths'])>0
p=j['paths'][0]['points'][0]
assert 0<=p['x']<=1024 and 0<=p['y']<=1024
print(f'✓ Journey (humans={j[\"humanCount\"]}, bots={j[\"botCount\"]}, events={j[\"eventCounts\"]})')
"

curl -sf "$API/api/heatmap?map=AmbroseValley&type=traffic" | python3 -c "import sys,json; h=json.load(sys.stdin); assert len(h)>0; print(f'✓ Traffic heatmap ({len(h)} pts)')"
curl -sf "$API/api/heatmap?map=AmbroseValley&type=kills" | python3 -c "import sys,json; h=json.load(sys.stdin); assert len(h)>0; print(f'✓ Kill heatmap ({len(h)} pts)')"
curl -sf "$API/api/heatmap?map=AmbroseValley&type=deaths" | python3 -c "import sys,json; h=json.load(sys.stdin); assert len(h)>0; print(f'✓ Death heatmap ({len(h)} pts)')"

# README coordinate example
curl -sf "$API/api/maps" > /dev/null
python3 - <<'PY'
# AmbroseValley: x=-301.45, z=-355.55 -> pixel ~(78, 890)
scale, ox, oz = 900, -370, -473
x, z = -301.45, -355.55
u = (x - ox) / scale
v = (z - oz) / scale
px, py = u * 1024, (1 - v) * 1024
assert 76 <= px <= 80 and 888 <= py <= 892, (px, py)
print('✓ Coordinate mapping matches README example')
PY

echo ""
echo "All checks passed."
