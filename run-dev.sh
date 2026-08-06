docker compose down
docker compose --profile prod down
docker compose -f docker-compose.yml -f docker-compose.dev.yml up -d --build