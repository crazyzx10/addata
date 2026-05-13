# 微信小程序广告数据定时拉取

这是一个Windows计划任务配置脚本，用于每天自动执行数据拉取。

## 使用方法

### 方法一：手动创建计划任务

1. 打开Windows任务计划程序
2. 创建基本任务
3. 任务名称：`微信广告数据每日拉取`
4. 触发器：每天早上 9:00 执行
5. 操作：启动程序
   - 程序：`cmd.exe`
   - 参数：`/c "e:\ZZX\Projects\addata\run_fetch.bat"`
   - 起始位置：`e:\ZZX\Projects\addata`

### 方法二：使用命令创建计划任务

以管理员权限运行PowerShell，执行以下命令：

```powershell
$action = New-ScheduledTaskAction -Execute "cmd.exe" -Argument '/c "e:\ZZX\Projects\addata\run_fetch.bat"' -WorkingDirectory "e:\ZZX\Projects\addata"
$trigger = New-ScheduledTaskTrigger -Daily -At "09:00"
$settings = New-ScheduledTaskSettingsSet -AllowStartIfOnBatteries -DontStopIfGoingOnBatteries
Register-ScheduledTask -TaskName "微信广告数据每日拉取" -Action $action -Trigger $trigger -Settings $settings -Description "每天自动拉取微信小程序广告数据"
```

### 方法三：使用已创建的脚本创建计划任务

以管理员权限运行PowerShell，执行：

```powershell
.\create_scheduled_task.ps1
```

## 注意事项

1. 确保Node.js已安装并添加到系统PATH
2. 确保MySQL数据库可访问
3. 确保微信小程序凭证有效
4. 首次运行会拉取全部历史数据（从2016-01-01开始），可能需要较长时间

## 手动执行

直接双击 `run_fetch.bat` 或在命令行运行：

```bash
node fetch_wechat_ad_data.js
```
