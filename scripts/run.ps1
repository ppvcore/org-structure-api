docker-compose down
docker-compose build
docker-compose up -d db migrate
Start-Sleep -Seconds 5
docker-compose run --rm migrate
docker-compose up -d app
docker-compose logs -f