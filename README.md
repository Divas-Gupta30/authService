# authService

# Supported Curls
curl -X POST -d "email=test@example.com&password=testpass" http://localhost:8080/signup <br />
curl -X POST -d "email=test@example.com&password=testpass" http://localhost:8080/signin <br />
curl -X GET -H "Authorization: <JWT_TOKEN>" http://localhost:8080/protected <br />
curl -X POST -H "Authorization: <JWT_TOKEN>" http://localhost:8080/refresh

# Docker Command
Checkout */authService and run below commans
docker-compose up --build
