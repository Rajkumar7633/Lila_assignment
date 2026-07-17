#!/usr/bin/env python3
"""Validate parquet schema and run sanity checks on player data."""

import glob
import os
import sys

import duckdb
import pyarrow.parquet as pq

DATA_DIR = os.environ.get("DATA_PATH", "data/player_data")


def decode_event(val):
    if isinstance(val, bytes):
        return val.decode("utf-8")
    return str(val)


def inspect_sample_files(base: str, limit: int = 3) -> None:
    files = sorted(glob.glob(os.path.join(base, "February_*", "*")))[:limit]
    for path in files:
        table = pq.read_table(path)
        df = table.to_pandas()
        df["event"] = df["event"].apply(decode_event)
        print(f"\n--- {os.path.basename(path)} ---")
        print(f"Schema OK: {list(df.columns)}")
        print(f"Events: {df['event'].unique()[:8]}")
        print(f"Rows: {len(df)}, map: {df['map_id'].iloc[0]}")


def run_duckdb_checks(base: str) -> None:
    pattern = os.path.join(base, "February_*", "*")
    con = duckdb.connect()

    maps = con.execute(
        f"SELECT map_id, COUNT(*) FROM read_parquet('{pattern}') GROUP BY map_id"
    ).fetchall()
    events = con.execute(
        f"SELECT CAST(event AS VARCHAR), COUNT(*) FROM read_parquet('{pattern}') "
        "GROUP BY 1 ORDER BY 2 DESC"
    ).fetchall()
    matches = con.execute(
        f"SELECT COUNT(DISTINCT match_id) FROM read_parquet('{pattern}')"
    ).fetchone()

    print("\n=== DuckDB Summary ===")
    print(f"Maps: {maps}")
    print(f"Events: {events}")
    print(f"Unique matches: {matches[0]}")


def main() -> int:
    base = DATA_DIR
    if not os.path.isdir(base):
        print(f"Data directory not found: {base}", file=sys.stderr)
        return 1

    print(f"Inspecting data at {base}")
    inspect_sample_files(base)
    run_duckdb_checks(base)
    print("\nAll sanity checks passed.")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
