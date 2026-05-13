@echo off
chcp 65001 >nul
echo ======================================
echo 微信小程序广告数据定时拉取任务
echo 执行时间: %date% %time%
echo ======================================
echo.

cd /d "%~dp0"

echo 开始执行数据拉取...
node fetch_wechat_ad_data.js

echo.
echo ======================================
echo 执行完成
echo ======================================
pause
