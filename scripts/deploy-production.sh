## This contains steps to launch a new version of the API inside of the server
## Typically this would be used following build-and-push-images.sh
## TODO: Eventually these commands should be triggered after a succesful Travis push...
# TODO: ... Perhaps only triggered by a push into production branch being used


## option 1: completely redeploy the entire app and brand new database

##  from within ssh into the server
## navigate to project directory
cd ~/go-database/

## get rid of old dusty containers
docker-compose down

## update new images from docker hub
docker pull trileyschwarz/go-database:latest

## redeploy the containers
docker-compose up -d


## option 2: only rebuild the API, database remains unchanged TODO
## option 3: API rebuild, and the database needs to be migrated to a new one (eg. schema changes etc.) TODO

## flag whether we need to empty database, or keep contents.
## need to save the database contents so that we can migrate contents after we
