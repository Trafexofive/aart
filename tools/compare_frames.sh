#!/bin/bash
# Compare two frames from an .aart file

FILE="${1:-test_giphy.aart}"
FRAME1="${2:-0}"
FRAME2="${3:-1}"

echo "=== Frame $FRAME1 ==="
awk "/^frame: $FRAME1$/,/^---$/" "$FILE" | head -30

echo ""
echo "=== Frame $FRAME2 ==="
awk "/^frame: $FRAME2$/,/^---$/" "$FILE" | head -30
