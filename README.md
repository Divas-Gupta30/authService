# authService

# Supported Curls
1) curl -X POST -H "Content-Type: application/json" -d '{"email": "test@example.com", "password": "testpass"}' http://localhost:8080/signup <br />
2) curl -X POST -H "Content-Type: application/json" -d '{"email": "test@example.com", "password": "testpass"}' http://localhost:8080/signin <br />
3) curl -X GET -H "Authorization: Bearer <JWT_TOKEN>" http://localhost:8080/anyprotectedroute <br />
4) curl -X POST -H "Authorization: Bearer <JWT_TOKEN>" http://localhost:8080/refresh <br />
 
5) curl -X POST -H "Authorization: Bearer <JWT_TOKEN>" http://localhost:8080/revoke <br />


# Docker Command
Checkout */authService and run below command : <br />
docker-compose up --build
