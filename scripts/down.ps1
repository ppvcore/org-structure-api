docker compose down
docker rmi org-structure-api-app:latest
docker builder prune -a -f
docker system prune -a --volumes -f