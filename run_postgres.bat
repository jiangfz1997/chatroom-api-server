@echo off
echo 🚀 启动 PostgreSQL Docker 容器...

docker run --name chatroom-postgres ^
 -e POSTGRES_USER=chat ^
 -e POSTGRES_PASSWORD=123456 ^
 -e POSTGRES_DB=chatroom ^
 -v pg_data:/var/lib/postgresql/data ^
 -p 5432:5432 ^
 -d postgres

echo ✅ 启动完成！
pause