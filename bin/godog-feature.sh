set -euo pipefail

if [ -f ".env" ]; then
  cp .env.sample .env
else
  cp .env.sample .env
fi

echo ${feature}
godog run features/${feature}/*.feature --format=cucumber > report.json