#!/bin/bash

COVERAGE_FILE=$1
THRESHOLD=$2

if [ -z "$COVERAGE_FILE" ] || [ -z "$THRESHOLD" ]; then
  echo "Usage: $0 <coverage_file> <threshold>"
  exit 1
fi

if [ ! -f "$COVERAGE_FILE" ]; then
  echo "Error: Coverage file '$COVERAGE_FILE' not found."
  exit 1
fi

# Extract the total percentage.
# Example line: "total: (statements) 22.2%"
TOTAL_LINE=$(grep "total:" "$COVERAGE_FILE")
if [ -z "$TOTAL_LINE" ]; then
  echo "Error: Could not find 'total:' line in coverage file."
  exit 1
fi

# Extract number (remove % sign)
PERCENTAGE=$(echo "$TOTAL_LINE" | awk '{print $3}' | sed 's/%//')

echo "Total coverage: ${PERCENTAGE}%"
echo "Threshold: ${THRESHOLD}%"

# Compare using awk for floating point comparison
IS_BELOW=$(awk -v p="$PERCENTAGE" -v t="$THRESHOLD" 'BEGIN {print (p < t)}')

if [ "$IS_BELOW" -eq 1 ]; then
  echo "FAIL: Coverage is below threshold."
  exit 1
else
  echo "PASS: Coverage is above or equal to threshold."
  exit 0
fi
