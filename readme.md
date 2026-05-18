
# 微信小程序广告数据定时拉取

这是一个用于定时拉取微信小程序广告数据的自动化工具。

## 快速开始

### 1. 安装依赖

```bash
npm install
```

### 2. 配置环境变量

复制 `.env.example` 并重命名为 `.env`：

```bash
copy .env.example .env
```

编辑 `.env` 文件，填入你的配置：

```env
# 数据库配置
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_DATABASE=wechat_ad

# API配置
API_BASE=https://api.weixin.qq.com/publisher/stat
START_DATE=2025-07-01

# 小程序配置（替换为你的真实信息）
MINI_PROGRAM_1_NAME=你的小程序名称
MINI_PROGRAM_1_APPID=wx1234567890abcdef
MINI_PROGRAM_1_APPSECRET=your_app_secret
```

### 3. 运行程序

```bash
node fetch_wechat_ad_data.js
```

或双击 `run_fetch.bat`

## 添加新小程序

在 `.env` 文件中添加新的小程序配置：

```env
MINI_PROGRAM_1_NAME=小程序1
MINI_PROGRAM_1_APPID=wx1111111111111111
MINI_PROGRAM_1_APPSECRET=secret1

MINI_PROGRAM_2_NAME=小程序2
MINI_PROGRAM_2_APPID=wx2222222222222222
MINI_PROGRAM_2_APPSECRET=secret2
```

配置编号必须连续（如 1, 2, 3...）

## 设置定时任务

如需设置每天自动拉取数据，请参考 `SCHEDULE_TASK_README.txt`

## 文件清单

* `fetch_wechat_ad_data.js` - 主程序
* `run_fetch.bat` - 快捷启动脚本
* `create_scheduled_task.ps1` - 定时任务脚本
* `SCHEDULE_TASK_README.txt` - 定时任务配置说明
* `.env.example` - 环境变量配置示例

## 注意事项

1. 确保 MySQL 数据库已启动并可访问
2. 确保微信小程序凭证有效
3. 首次运行会拉取从 START_DATE 开始的所有历史数据
4. `.env` 文件包含敏感信息，请勿提交到版本控制
