# authService

# Supported Curls
1) curl -X POST -d "email=test@example.com&password=testpass" http://localhost:8080/signup <br />
2) curl -X POST -d "email=test@example.com&password=testpass" http://localhost:8080/signin <br />
3) curl -X GET -H "Authorization: <JWT_TOKEN>" http://localhost:8080/protected <br />
4) curl -X POST -H "Authorization: <JWT_TOKEN>" http://localhost:8080/refresh

# Docker Command
Checkout */authService and run below command : <br />
docker-compose up --build
