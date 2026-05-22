# WMAM 界面设计文档

## 概述

本文档描述了WMAM多用户版本的界面设计规范，采用Vercel风格设计，支持明暗主题切换。

---

## 1. 设计理念

### 1.1 风格参考

**Vercel 风格特点：
- 极简扁平，无过度装饰
- 大量留白，视觉舒适
- 卡片式布局，层次清晰
- 左侧固定导航栏
- 顶部操作区域
- 平滑的过渡动画
- 明暗主题无缝切换

### 1.2 核心原则

1. **简洁优先** - 去除一切不必要的元素
2. **清晰可识别的设计 - 一眼能找到需要的功能
3. **一致性 - 统一的交互模式
4. **响应友好 - 适配多种屏幕尺寸

---

## 2. 配色方案

### 2.1 亮色主题 (Light Mode)

| 变量名 | 颜色值 | 用途 |
|--------|--------|------|
| --bg-primary | #ffffff | 主背景色 |
| --bg-secondary | #fafafa | 次要背景色（导航栏）|
| --bg-card | #ffffff | 卡片背景 |
| --bg-hover | #f5f5f5 | 悬浮背景 |
| --text-primary | #000000 | 主要文字 |
| --text-secondary | #666666 | 次要文字 |
| --text-muted | #999999 | 弱化文字 |
| --border | #e5e5e5 | 边框 |
| --accent | #000000 | 强调色（按钮） |
| --accent-hover | #333333 | 强调色悬停 |
| --success | #22c55e | 成功状态 |
| --warning | #f59e0b | 警告状态 |
| --error | #ef4444 | 错误状态 |

### 2.2 暗色主题 (Dark Mode)

| 变量名 | 颜色值 | 用途 |
|--------|--------|------|
| --bg-primary | #0a0a0a | 主背景色 |
| --bg-secondary | #111111 | 次要背景色（导航栏）|
| --bg-card | #111111 | 卡片背景 |
| --bg-hover | #1a1a1a | 悬浮背景 |
| --text-primary | #ffffff | 主要文字 |
| --text-secondary | #a1a1aa | 次要文字 |
| --text-muted | #71717a | 弱化文字 |
| --border | #27272a | 边框 |
| --accent | #ffffff | 强调色（按钮） |
| --accent-hover | #e5e5e5 | 强调色悬停 |
| --success | #22c55e | 成功状态 |
| --warning | #f59e0b | 警告状态 |
| --error | #ef4444 | 错误状态 |

### 2.3 主题切换

- 切换按钮位于右上角
- 采用平滑过渡动画（0.3s ease）
- 使用 localStorage 保存用户偏好
- 首次访问根据系统主题自动判断

---

## 3. 布局架构

### 3.1 整体布局

```
┌──────────────────────────────────────────────────────┐
│ Logo        Theme      用户名  │ ← 顶部栏
├──────────┬───────────────────────────────────────┤
│          │                                         │
│  侧边栏  │          主内容区                       │
│          │                                         │
│ 📊仪表盘 │                                         │
│ 📱小程序 │                                         │
│ ⏱️执行拉取 │                                         │
│ 📜日志    │                                         │
│ 👤用户管 │                                         │
│          │                                         │
└──────────┴───────────────────────────────────────┘
```

### 3.2 布局规范

- **侧边栏宽度**：240px（固定）
- **顶部栏高度**：64px（固定）
- **内容区边距**：24px
- **卡片间距**：16px
- **内容区域最大宽度**：1200px（居中）

---

## 4. 组件设计

### 4.1 顶部栏 (Header)

**样式规范**

- Logo 居左
- 主题切换按钮
- 用户信息（用户名、头像、登出）

**交互**

- 主题切换：点击切换明暗主题
- 用户头像点击：显示下拉菜单（登出选项）

### 4.2 侧边栏 (Sidebar)

**导航项：

- 📊 仪表盘 (Dashboard) - 默认首页
- 📱 小程序 (Mini Programs) - 小程序列表管理
- ⏱️ 执行拉取 (Execute Fetch) - 数据拉取执行
- 📜 操作日志 (Activity Log) - 操作历史记录
- 👤 用户管理 (User Management) - 用户管理（仅管理员可见）

**选中状态：

- 背景色：--bg-hover
- 左侧指示条：3px，颜色 --accent
- 文字加粗

### 4.3 卡片 (Card)

**样式**

- 圆角：8px
- 边框：1px solid --border
- 背景色：--bg-card
- 内边距：24px
- 阴影：无（极简风格）
- 悬停效果：轻微上浮 +2px

### 4.4 按钮 (Button)

**主按钮 (Primary)

- 背景色：--accent
- 文字色：亮色主题白/暗色主题黑
- 圆角：6px
- 内边距：10px 20px
- 悬停：--accent-hover
- 文字：600 字重

**次要按钮 (Secondary)

- 背景色：透明
- 边框：1px solid --border
- 文字色：--text-primary
- 悬停：--bg-hover

**危险按钮 (Danger)

- 背景色：--error
- 文字色：白色

### 4.5 输入框 (Input)

- 圆角：6px
- 边框：1px solid --border
- 背景色：--bg-primary
- 内边距：10px 14px
- 聚焦状态：边框 --accent，无阴影

---

## 5. 页面设计

### 5.1 登录页 (Login Page)

```
┌─────────────────────────────────┐
│                                 │
│        WMAM Logo            │
│                                 │
│      ┌───────────────────┐    │
│      │  用户名输入框    │    │
│      └───────────────────┘    │
│      ┌───────────────────┐    │
│      │  密码输入框      │    │
│      └───────────────────┘    │
│                                 │
│      ┌───────────────────┐    │
│      │   登录按钮      │    │
│      └───────────────────┘    │
│                                 │
│      Copyright ...           │
└─────────────────────────────────┘
```

**内容：

- Logo（顶部
- 输入框：用户名、密码
- 登录按钮
- 底部版权信息

### 5.2 仪表盘页 (Dashboard Page)

**统计卡片（4个）

```
┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐
│ 用户数 │ │ 小程序数  │ │ 今日拉取  │ │ 最后拉取  │
│  5     │ │  12     │ │  3     │ │ 2h前   │
└──────────┘ └──────────┘ └──────────┘ └──────────┘
```

**快速操作**

- 开始拉取数据按钮（主按钮）

### 5.3 小程序管理页 (Mini Programs Page)

**布局**

```
标题 + 添加按钮

小程序列表卡片
┌─────────────────────────────────────────┐
│ 小程序1  AppID: xxx  [编辑] [删除]│
│ 小程序2  AppID: xxx  [编辑] [删除]│
│ ...                               │
└─────────────────────────────────────────┘
```

**功能（仅管理员可见编辑/删除按钮）：

- 普通用户仅查看，无操作按钮

### 5.4 执行拉取页 (Execute Fetch Page)

**布局**

```
小程序状态卡片区

小程序1  [就绪] [就绪]
小程序2  [就绪] [就绪]

进度条 [100%]

操作按钮 [开始] [中断] [继续] [重新]

日志区域
> ...
```

**功能**

- 左侧：小程序状态卡片
- 中间：进度条 + 操作按钮
- 底部：实时日志
- 状态显示正在拉取中，显示提示信息

### 5.5 操作日志页 (Activity Log Page)

**日志列表，时间线样式

```
2025-05-22 14:30 用户A 开始拉取数据
2025-05-22 12:00 用户B 拉取完成
2025-05-22 10:00 管理员 添加小程序
```

### 5.6 用户管理页 (User Management Page)

**仅管理员可见**

```
用户列表

用户名  角色    [编辑] [禁用] [删除]

添加用户按钮
```

---

## 6. 交互设计

### 6.1 状态反馈

- **加载状态**：Spinner
- **成功提示**：Toast 消息（右上角，3秒自动消失）
- **错误提示**：Toast 消息
- **确认操作**：弹窗

### 6.2 动画

- **主题切换：0.3s ease
- **按钮交互：0.15s ease
- **页面过渡**：0.2s ease-in-out

---

## 7. 技术实现方案

### 7.1 CSS 变量定义

```css
:root {
  /* 亮色主题 */
  --bg-primary: #ffffff;
  --bg-secondary: #fafafa;
  --bg-card: #ffffff;
  --bg-hover: #f5f5f5;
  --text-primary: #000000;
  --text-secondary: #666666;
  --text-muted: #999999;
  --border: #e5e5e5;
  --accent: #000000;
  --accent-hover: #333333;
}

[data-theme="dark"] {
  --bg-primary: #0a0a0a;
  --bg-secondary: #111111;
  --bg-card: #111111;
  --bg-hover: #1a1a1a;
  --text-primary: #ffffff;
  --text-secondary: #a1a1aa;
  --text-muted: #71717a;
  --border: #27272a;
  --accent: #ffffff;
  --accent-hover: #e5e5e5;
}

* {
  transition: background-color 0.3s ease, border-color 0.3s ease, color 0.3s ease;
}
```

### 7.2 主题切换逻辑

```javascript
function initTheme() {
  const saved = localStorage.getItem('theme');
  const system = window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light';
  const theme = saved || system;
  document.documentElement.setAttribute('data-theme', theme);
  return theme;
}

function toggleTheme() {
  const current = document.documentElement.getAttribute('data-theme');
  const next = current === 'dark' ? 'light' : 'dark';
  document.documentElement.setAttribute('data-theme', next);
  localStorage.setItem('theme', next);
  return next;
}
```

---

## 8. 响应式适配

- **桌面端**：> 768px - 完整布局
- **平板/手机端**：<= 768px - 侧边栏折叠成汉堡菜单
- **移动端**：<= 480px - 单列布局

---

## 9. 图标方案

采用 SVG 图标，确保在明暗主题自适应
