set -euo pipefail

head -n -2 report.json > report.json.temp

cp report.json.temp report.json