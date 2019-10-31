## This script is intended to be used to build and push the go-database images to docker-hub
## Should be run from the root of this project
docker build --tag trileyschwarz/go-database:latest .
docker push trileyschwarz/go-database:latest