@echo off
echo Stopping PostgreSQL...

docker stop chatroom-postgres
docker rm chatroom-postgres

echo Finish
pause