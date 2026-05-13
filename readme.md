添加新小程序的方法：
打开 e:\\ZZX\\Projects\\addata\\fetch\_wechat\_ad\_data.js ，在 MINI\_PROGRAMS 数组中添加新配置：


const MINI\_PROGRAMS = \[
{
name: '光小充',
appid: 'wx92518cf320f09758',
appsecret: 'cb3e003deafa47458e06f511def0020d'
},
// 添加更多小程序
// {
//     name: '您的小程序名称',
//     appid: 'xxx',
//     appsecret: 'xxx'
// }
];



文件清单：

* fetch\_wechat\_ad\_data.js - 主程序
* run\_fetch.bat - 快捷启动脚本
* create\_scheduled\_task.ps1 - 定时任务脚本
* SCHEDULE\_TASK\_README.txt - 定时任务配置说明
如需设置每天自动拉取数据，可以参考 SCHEDULE\_TASK\_README.txt 创建Windows定时任务。

