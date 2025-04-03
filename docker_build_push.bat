@echo off
setlocal

:: === 配置项 ===
set DOCKER_USERNAME=jiangfz
set IMAGE_NAME=online-chatroom-apiserver
set FULL_IMAGE=%DOCKER_USERNAME%/%IMAGE_NAME%:latest

echo.
echo 🚀 构建镜像: %FULL_IMAGE%
docker build -t %FULL_IMAGE% .

if %errorlevel% neq 0 (
    echo ❌ Docker build 失败！
    exit /b 1
)

echo.
echo 🚚 推送到 Docker Hub...
docker push %FULL_IMAGE%

if %errorlevel% neq 0 (
    echo ❌ Docker push 失败！
    exit /b 1
)

echo.
echo ✅ 镜像推送成功！
echo 🌐 镜像地址: %FULL_IMAGE%

endlocal
pause
