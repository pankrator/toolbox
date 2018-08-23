docker stop $(docker ps -q --filter name=pr-)
docker rm $(docker ps -a -q --filter name=pr-)
docker network rm $(docker network ls -q --filter name=pr-)