# 创建Windows计划任务脚本
# 以管理员权限运行此脚本

$taskName = "微信广告数据每日拉取"

$scriptPath = Split-Path -Parent $MyInvocation.MyCommand.Path
$batPath = Join-Path $scriptPath "run_fetch.bat"

Write-Host "========================================"
Write-Host "微信广告数据定时任务安装脚本"
Write-Host "========================================"
Write-Host ""

$action = New-ScheduledTaskAction -Execute "cmd.exe" -Argument "/c `"$batPath`"" -WorkingDirectory $scriptPath
$trigger = New-ScheduledTaskTrigger -Daily -At "09:00"
$settings = New-ScheduledTaskSettingsSet -AllowStartIfOnBatteries -DontStopIfGoingOnBatteries

try {
    $existingTask = Get-ScheduledTask -TaskName $taskName -ErrorAction SilentlyContinue
    if ($existingTask) {
        Write-Host "任务已存在，正在更新..."
        Unregister-ScheduledTask -TaskName $taskName -Confirm:$false
    }

    Register-ScheduledTask -TaskName $taskName -Action $action -Trigger $trigger -Settings $settings -Description "每天自动拉取微信小程序广告数据" | Out-Null

    Write-Host "计划任务创建成功！" -ForegroundColor Green
    Write-Host ""
    Write-Host "任务名称: $taskName"
    Write-Host "执行时间: 每天 09:00"
    Write-Host "执行脚本: $batPath"
    Write-Host ""
    Write-Host "可使用以下命令查看和管理任务：" -ForegroundColor Yellow
    Write-Host "  查看任务: Get-ScheduledTask -TaskName `"$taskName`""
    Write-Host "  立即执行: Start-ScheduledTask -TaskName `"$taskName`""
    Write-Host "  删除任务: Unregister-ScheduledTask -TaskName `"$taskName`" -Confirm:`$false"
} catch {
    Write-Host "创建计划任务失败: $_" -ForegroundColor Red
    Write-Host "请确保以管理员权限运行此脚本" -ForegroundColor Yellow
}
