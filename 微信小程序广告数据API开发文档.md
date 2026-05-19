# 微信小程序广告数据 API 开发文档

> **版本**: v1.0  
> **更新日期**: 2026-05-11  
> **适用对象**: 研发部开发同学

---

## 一、概述

本文档描述如何通过微信官方 API 获取小程序广告数据，包括广告位清单、汇总数据、细分数据和结算数据。

### 数据获取流程

```
获取 access_token → 调用业务 API 获取数据 → 数据转换（分→元）→ 写入目标存储
```

### API 基础信息

| 项目 | 说明 |
|------|------|
| Base URL | `https://api.weixin.qq.com` |
| 请求方法 | **全部为 GET**（参数通过 Query String 传递） |
| 鉴权方式 | access_token（URL 参数） |
| 数据格式 | JSON |
| 金额单位 | **所有金额字段均为"分"**，需除以 100 转换为"元" |
| access_token 有效期 | 2 小时（7200 秒），需缓存并在过期前刷新 |

---

## 二、鉴权：获取 access_token

### 请求

```
GET https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid={APPID}&secret={APPSECRET}
```

### 参数

| 参数 | 必填 | 说明 |
|------|------|------|
| grant_type | 是 | 固定值 `client_credential` |
| appid | 是 | 小程序的 AppID |
| secret | 是 | 小程序的 AppSecret |

### 响应示例

```json
{
  "access_token": "103_J7UJyrN2B3Z_5xKdE8cFvA0bRqWmT4hLpYsIoXnCjG6fDuSaQV",
  "expires_in": 7200
}
```

### 错误码

| errcode | 说明 |
|---------|------|
| -1 | 系统繁忙，稍后重试 |
| 40001 | AppSecret 错误 |
| 40013 | AppID 无效 |

### ⚠️ 注意事项

- **必须缓存 access_token**，不要每次调用业务 API 都重新获取
- 建议在过期前 5 分钟刷新
- 多个小程序需分别获取各自的 access_token

---

## 三、业务 API

所有业务 API 统一入口：

```
GET https://api.weixin.qq.com/publisher/stat?access_token={TOKEN}&action={ACTION}&{其他参数}
```

### 通用参数

| 参数 | 必填 | 说明 |
|------|------|------|
| access_token | 是 | 通过第二步获取 |
| action | 是 | 业务动作名称，见下表 |
| page | 否 | 页码，从 1 开始 |
| page_size | 否 | 每页条数，最大 90 |

### 通用响应结构

```json
{
  "ret": 0,              // 0=成功，非0=失败
  "err_msg": "ok",       // 错误信息
  "base_resp": {
    "ret": 0,
    "err_msg": "ok"
  },
  "total_num": 100,      // 总记录数
  "list": [...]          // 数据列表
}
```

### 通用错误码

| ret/errcode | 说明 |
|-------------|------|
| 0 | 成功 |
| 45009 | API 调用频率超限 |
| 45010 | 无权限（action 不支持或未开通） |
| 48001 | API 未授权 |

---

## 四、API 1：广告位清单

获取小程序下所有广告位的配置信息。

### 请求

```
GET https://api.weixin.qq.com/publisher/stat?access_token={TOKEN}&action=get_adunit_list&page=1&page_size=999
```

### 参数

| 参数 | 值 | 说明 |
|------|------|------|
| action | `get_adunit_list` | 固定值 |

### 响应示例

```json
{
  "ret": 0,
  "total_num": 5,
  "ad_unit": [
    {
      "ad_unit_id": "adunit-376fd270a6514ce8",
      "ad_unit_name": "视频贴片广告",
      "ad_unit_type": "AD_UNIT_TYPE_VIDEO_BEGIN",
      "ad_slot": "SLOT_ID_WEAPP_VIDEO_BEGIN",
      "ad_unit_status": "AD_UNIT_STATUS_OFF",
      "slot_id": "7090665964306299",
      "appid": "wx92518cf320f09758",
      "ad_unit_size": [{"width": 582, "height": 166}],
      "is_allow_playable": true,
      "video_duration_min": 6,
      "video_duration_max": 86400
    },
    {
      "ad_unit_id": "adunit-ff3cc53c4a4a6409",
      "ad_unit_name": "扫码插屏广告",
      "ad_unit_type": "AD_UNIT_TYPE_INTERSTITIAL",
      "ad_slot": "SLOT_ID_WEAPP_INTERSTITIAL",
      "ad_unit_status": "AD_UNIT_STATUS_ON",
      "slot_id": "3030046789020061",
      "appid": "wx92518cf320f09758"
    },
    {
      "ad_unit_id": "adunit-b48085c884c6e7e2",
      "ad_unit_name": "其他模版广告",
      "ad_unit_type": "AD_UNIT_TYPE_TEMPLATE_CUSTOM",
      "ad_slot": "SLOT_ID_WEAPP_TEMPLATE",
      "ad_unit_status": "AD_UNIT_STATUS_ON",
      "slot_id": "4071202390577885",
      "appid": "wx92518cf320f09758",
      "templ_type_list": ["TEMPL_TYPE_BANNER_ADS"]
    },
    {
      "ad_unit_id": "adunit-bc8fa74828eaae9b",
      "ad_unit_name": "视频激励广告",
      "ad_unit_type": "AD_UNIT_TYPE_REWARED_VIDEO",
      "ad_slot": "SLOT_ID_WEAPP_REWARD_VIDEO",
      "ad_unit_status": "AD_UNIT_STATUS_ON",
      "slot_id": "1030436212907001",
      "appid": "wx92518cf320f09758"
    },
    {
      "ad_unit_id": "adunit-364cfb568c579145",
      "ad_unit_name": "扫码原生模版广告",
      "ad_unit_type": "AD_UNIT_TYPE_TEMPLATE_CUSTOM",
      "ad_slot": "SLOT_ID_WEAPP_TEMPLATE",
      "ad_unit_status": "AD_UNIT_STATUS_ON",
      "slot_id": "4071202390577885",
      "appid": "wx92518cf320f09758"
    }
  ]
}
```

### 关键字段说明

| 字段 | 类型 | 说明 |
|------|------|------|
| ad_unit_id | string | 广告位唯一ID（adunit-xxx 格式） |
| ad_unit_name | string | 广告位名称 |
| ad_slot | string | 广告位类型枚举（见下表） |
| ad_unit_status | string | 状态：`AD_UNIT_STATUS_ON`=正常，`AD_UNIT_STATUS_OFF`=暂停 |
| slot_id | string | 广告位数字ID（汇总/细分数据中用此关联） |

### 广告位类型枚举（ad_slot）

| 枚举值 | 中文名 | 说明 |
|--------|--------|------|
| SLOT_ID_WEAPP_VIDEO_BEGIN | 视频贴片 | 视频播放前展示 |
| SLOT_ID_WEAPP_INTERSTITIAL | 插屏 | 全屏弹窗广告 |
| SLOT_ID_WEAPP_BANNER | Banner | 底部横幅广告 |
| SLOT_ID_WEAPP_REWARD_VIDEO | 激励视频 | 看视频得奖励 |
| SLOT_ID_WEAPP_TEMPLATE | 原生模版 | 自定义模版，含原生/Banner等子类型 |
| SLOT_ID_WEAPP_COVER | 封面 | 开屏封面广告 |

> 💡 `SLOT_ID_WEAPP_TEMPLATE` 是一个大类，可能包含原生模版和 Banner 模版等子类型。同一 slot_id 可能对应多个 ad_unit_id。

---

## 五、API 2：汇总数据

按广告位类型汇总的每日广告数据。

### 请求

```
GET https://api.weixin.qq.com/publisher/stat?access_token={TOKEN}&action=publisher_adpos_general&page=1&page_size=90&start_date=2026-05-01&end_date=2026-05-09
```

### 参数

| 参数 | 必填 | 说明 |
|------|------|------|
| action | 是 | `publisher_adpos_general` |
| start_date | 是 | 起始日期，格式 `YYYY-MM-DD` |
| end_date | 是 | 结束日期，格式 `YYYY-MM-DD` |
| ad_slot | 否 | 按广告位类型筛选（不传则返回全部类型） |
| page | 是 | 页码，从 1 开始 |
| page_size | 是 | 每页条数，**最大 90** |

### 响应示例

```json
{
  "ret": 0,
  "total_num": 8,
  "summary": {
    "req_succ_count": 152610,
    "exposure_count": 81696,
    "exposure_rate": 0.5353253390,
    "click_count": 1351,
    "click_rate": 0.0165369170,
    "income": 169419,
    "ecpm": 2073.7735017630
  },
  "list": [
    {
      "date": "2026-05-09",
      "ad_slot": "SLOT_ID_WEAPP_COVER",
      "slot_id": 5060180989186180,
      "slot_str": "5060180989186180",
      "req_succ_count": 12732,
      "exposure_count": 8424,
      "exposure_rate": 0.6616399620,
      "click_count": 292,
      "click_rate": 0.0346628680,
      "income": 41154,
      "ecpm": 4885.3276353280
    },
    {
      "date": "2026-05-09",
      "ad_slot": "SLOT_ID_WEAPP_INTERSTITIAL",
      "slot_id": 3030046789020061,
      "slot_str": "3030046789020061",
      "req_succ_count": 13362,
      "exposure_count": 7095,
      "exposure_rate": 0.5309833859999999,
      "click_count": 74,
      "click_rate": 0.010429880,
      "income": 4574,
      "ecpm": 644.6793516560
    }
  ]
}
```

### 关键字段说明

| 字段 | 类型 | 单位 | 说明 |
|------|------|------|------|
| date | string | - | 日期 `YYYY-MM-DD` |
| ad_slot | string | - | 广告位类型枚举 |
| slot_id | number | - | 广告位数字ID |
| slot_str | string | - | 广告位数字ID（字符串形式） |
| req_succ_count | number | 次 | **请求量**（成功请求次数） |
| exposure_count | number | 次 | **曝光量** |
| exposure_rate | number | - | **曝光率** = 曝光量/请求量 |
| click_count | number | 次 | **点击量** |
| click_rate | number | - | **点击率** = 点击量/曝光量（已是小数，如 0.0346 = 3.46%） |
| income | number | **分** | **总收入**（⚠️ 需 ÷100 转为元） |
| ecpm | number | **分** | **千次曝光收入**（⚠️ 需 ÷100 转为元） |

### ⚠️ 重要提示

1. **金额单位是"分"**！`income: 41154` 表示 411.54 元，`ecpm: 4885.33` 表示 48.85 元/千次曝光
2. **page_size 最大 90**，数据量超过 90 条需要分页
3. `summary` 字段是当前查询范围的汇总统计
4. 同一日期可能有多个广告位类型的数据
5. 最早可查 2016-01-01 起的数据

---

## 六、API 3：细分数据

按具体广告位（ad_unit）细分的每日广告数据，比汇总数据多出广告位维度。

### 请求

```
GET https://api.weixin.qq.com/publisher/stat?access_token={TOKEN}&action=publisher_adunit_general&page=1&page_size=90&start_date=2026-05-09&end_date=2026-05-09
```

### 参数

| 参数 | 必填 | 说明 |
|------|------|------|
| action | 是 | `publisher_adunit_general` |
| start_date | 是 | 起始日期 |
| end_date | 是 | 结束日期 |
| page | 是 | 页码 |
| page_size | 是 | 每页条数，最大 90 |

### 响应示例

```json
{
  "ret": 0,
  "total_num": 4,
  "list": [
    {
      "ad_unit_id": "adunit-364cfb568c579145",
      "ad_unit_name": "扫码原生模版广告",
      "stat_item": {
        "date": "2026-05-09",
        "ad_slot": "SLOT_ID_WEAPP_REWARD_VIDEO",
        "slot_str": "1030436212907001",
        "req_succ_count": 48,
        "exposure_count": 30,
        "exposure_rate": 0.625,
        "click_count": 0,
        "click_rate": 0.0,
        "income": 18,
        "publisher_income": 18,
        "agency_income": 0,
        "ecpm": 600.0,
        "is_smart_ads": 0,
        "parent_templ_type": "nul"
      }
    },
    {
      "ad_unit_id": "adunit-ff3cc53c4a4a6409",
      "ad_unit_name": "扫码插屏广告",
      "stat_item": {
        "date": "2026-05-09",
        "ad_slot": "SLOT_ID_WEAPP_INTERSTITIAL",
        "slot_str": "3030046789020061",
        "req_succ_count": 13362,
        "exposure_count": 7095,
        "exposure_rate": 0.5309833859999999,
        "click_count": 74,
        "click_rate": 0.010429880,
        "income": 4573,
        "publisher_income": 4573,
        "agency_income": 0,
        "ecpm": 644.5384073290001,
        "is_smart_ads": 0,
        "parent_templ_type": "nul"
      }
    }
  ]
}
```

### 与汇总数据的区别

| 维度 | 汇总数据 (API 2) | 细分数据 (API 3) |
|------|-------------------|-------------------|
| 粒度 | 按 **广告位类型**（ad_slot） | 按 **具体广告位**（ad_unit_id） |
| 结构 | list 中每项是一条记录 | list 中每项包含 `stat_item` 对象 |
| 额外字段 | 无 | `publisher_income`、`agency_income`、`is_smart_ads`、`parent_templ_type` |
| 同一 slot_id | 可能只有一条 | 可能有 **多条**（如两个模版广告共享同一 slot_id） |

### 细分数据额外字段

| 字段 | 类型 | 单位 | 说明 |
|------|------|------|------|
| publisher_income | number | 分 | **流量主收入**（与 income 相同） |
| agency_income | number | 分 | 代理商收入（通常为 0） |
| is_smart_ads | number | - | 是否智能广告（0=否） |
| parent_templ_type | string | - | 父模版类型（"nul"=无） |

---

## 七、API 4：结算数据

获取流量主的结算汇总信息。

### 请求

```
GET https://api.weixin.qq.com/publisher/stat?access_token={TOKEN}&action=publisher_settlement&page=1&page_size=90
```

### 参数

| 参数 | 必填 | 说明 |
|------|------|------|
| action | 是 | `publisher_settlement` |

### 响应示例

```json
{
  "ret": 0,
  "total_num": 0,
  "revenue_all": 14114754,
  "settled_revenue_all": 12182857,
  "penalty_all": 0,
  "wyw_settled_summary": {
    "wyw_revenue_all": 0,
    "wyw_settled_revenue_all": 0,
    "wyw_penalty_all": 0
  }
}
```

### 关键字段说明

| 字段 | 类型 | 单位 | 说明 |
|------|------|------|------|
| revenue_all | number | **分** | 总预估收入 |
| settled_revenue_all | number | **分** | 总已结算收入 |
| penalty_all | number | **分** | 总罚金 |
| total_num | number | - | 结算明细条数（0=暂无明细） |
| wyw_settled_summary | object | - | 微信云开发相关结算汇总 |

> 💡 当 `total_num = 0` 时表示暂无结算明细记录，但 `revenue_all` 和 `settled_revenue_all` 仍可能显示累计数据。

---

## 八、数据转换规则

### 8.1 金额转换

所有金额字段（income、ecpm、publisher_income、agency_income、revenue_all、settled_revenue_all）的单位都是**分**，需转换为元：

```
元 = 分 ÷ 100
```

**示例**：
- `income: 41154` → **411.54 元**
- `ecpm: 4885.33` → **48.85 元/千次曝光**

### 8.2 点击率转换

API 返回的 `click_rate` 已经是小数形式，直接使用即可：

- `click_rate: 0.0347` → **3.47%**
- 展示时：`(click_rate * 100).toFixed(2) + '%'`

### 8.3 广告位类型名称映射

建议建立映射表，将枚举值转换为中文名称：

```javascript
const AD_SLOT_MAP = {
  'SLOT_ID_WEAPP_VIDEO_BEGIN': '视频贴片',
  'SLOT_ID_WEAPP_INTERSTITIAL': '插屏',
  'SLOT_ID_WEAPP_BANNER': 'Banner',
  'SLOT_ID_WEAPP_REWARD_VIDEO': '激励视频',
  'SLOT_ID_WEAPP_TEMPLATE': '原生模版',
  'SLOT_ID_WEAPP_COVER': '封面'
};
```

### 8.4 广告位状态映射

```javascript
const AD_STATUS_MAP = {
  'AD_UNIT_STATUS_ON': '正常',
  'AD_UNIT_STATUS_OFF': '暂停'
};
```

---

## 九、分页处理

### 分页规则

- `page` 从 **1** 开始
- `page_size` 最大 **90**
- `total_num` 返回符合条件的总记录数

### 分页算法

```javascript
async function fetchAllPages(token, action, params) {
  const all = [];
  let page = 1;
  const pageSize = 90;
  
  while (true) {
    const url = buildUrl(token, action, { ...params, page, page_size: pageSize });
    const data = await httpGet(url);
    
    if (data.ret !== 0) throw new Error(`API错误: ${data.err_msg}`);
    
    const list = data.list || data.ad_unit || [];
    if (list.length === 0) break;
    
    all.push(...list);
    
    // 已获取全部数据
    if (all.length >= data.total_num) break;
    
    page++;
    await sleep(500);  // 避免频率限制
  }
  
  return all;
}
```

### ⚠️ 频率限制

- 建议每次 API 调用间隔 **≥ 500ms**
- 超频会返回 `ret: 45009`
- 可用指数退避重试

---

## 十、完整代码示例（Node.js）

```javascript
const https = require('https');

// ==================== 配置 ====================
const APPID = 'your_appid';
const APPSECRET = 'your_appsecret';
const API_BASE = 'https://api.weixin.qq.com/publisher/stat';

// ==================== 工具函数 ====================

function httpGet(url) {
  return new Promise((resolve, reject) => {
    https.get(url, res => {
      let data = '';
      res.on('data', chunk => data += chunk);
      res.on('end', () => {
        try { resolve(JSON.parse(data)); }
        catch (e) { reject(new Error('JSON解析失败')); }
      });
    }).on('error', reject);
  });
}

function sleep(ms) { return new Promise(r => setTimeout(r, ms)); }

function fenToYuan(fen) { return (Number(fen) || 0) / 100; }

function fmtDate(d) {
  return `${d.getFullYear()}-${String(d.getMonth()+1).padStart(2,'0')}-${String(d.getDate()).padStart(2,'0')}`;
}

// ==================== API 调用 ====================

// 1. 获取 access_token
async function getToken() {
  const url = `https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=${APPID}&secret=${APPSECRET}`;
  const data = await httpGet(url);
  if (!data.access_token) throw new Error(`获取Token失败: ${JSON.stringify(data)}`);
  return data.access_token;
}

// 2. 通用分页请求
async function fetchAllPages(token, action, extraParams = {}) {
  const all = [];
  let page = 1;
  
  while (true) {
    const params = new URLSearchParams({
      access_token: token,
      action,
      page: String(page),
      page_size: '90',
      ...Object.fromEntries(
        Object.entries(extraParams).map(([k, v]) => [k, String(v)])
      )
    });
    
    const url = `${API_BASE}?${params}`;
    const data = await httpGet(url);
    
    if (data.ret !== 0) {
      throw new Error(`API错误 [${data.ret}]: ${data.err_msg}`);
    }
    
    const list = data.list || data.ad_unit || [];
    if (list.length === 0) break;
    
    all.push(...list);
    if (data.total_num !== undefined && all.length >= data.total_num) break;
    
    page++;
    await sleep(500);
  }
  
  return all;
}

// 3. 获取广告位清单
async function getAdunitList(token) {
  return await fetchAllPages(token, 'get_adunit_list');
}

// 4. 获取汇总数据
async function getSummaryData(token, startDate, endDate) {
  return await fetchAllPages(token, 'publisher_adpos_general', {
    start_date: startDate,
    end_date: endDate
  });
}

// 5. 获取细分数据
async function getDetailData(token, startDate, endDate) {
  return await fetchAllPages(token, 'publisher_adunit_general', {
    start_date: startDate,
    end_date: endDate
  });
}

// 6. 获取结算数据
async function getSettlementData(token) {
  const params = new URLSearchParams({
    access_token: token,
    action: 'publisher_settlement',
    page: '1',
    page_size: '90'
  });
  const data = await httpGet(`${API_BASE}?${params}`);
  if (data.ret !== 0) throw new Error(`结算API错误: ${data.err_msg}`);
  return data;
}

// ==================== 数据转换示例 ====================

const AD_SLOT_NAMES = {
  'SLOT_ID_WEAPP_VIDEO_BEGIN': '视频贴片',
  'SLOT_ID_WEAPP_INTERSTITIAL': '插屏',
  'SLOT_ID_WEAPP_BANNER': 'Banner',
  'SLOT_ID_WEAPP_REWARD_VIDEO': '激励视频',
  'SLOT_ID_WEAPP_TEMPLATE': '原生模版',
  'SLOT_ID_WEAPP_COVER': '封面'
};

// 转换汇总数据记录
function transformSummaryItem(item) {
  return {
    date: item.date,
    adSlotType: AD_SLOT_NAMES[item.ad_slot] || item.ad_slot,
    slotId: item.slot_str || String(item.slot_id),
    requestCount: item.req_succ_count,
    exposureCount: item.exposure_count,
    exposureRate: (item.exposure_rate * 100).toFixed(2) + '%',
    clickCount: item.click_count,
    clickRate: (item.click_rate * 100).toFixed(2) + '%',
    income: fenToYuan(item.income),         // 分→元
    ecpm: fenToYuan(item.ecpm)              // 分→元
  };
}

// 转换细分数据记录
function transformDetailItem(item) {
  const stat = item.stat_item;
  return {
    adUnitId: item.ad_unit_id,
    adUnitName: item.ad_unit_name,
    date: stat.date,
    adSlotType: AD_SLOT_NAMES[stat.ad_slot] || stat.ad_slot,
    slotId: stat.slot_str,
    requestCount: stat.req_succ_count,
    exposureCount: stat.exposure_count,
    clickCount: stat.click_count,
    clickRate: (stat.click_rate * 100).toFixed(2) + '%',
    income: fenToYuan(stat.income),
    publisherIncome: fenToYuan(stat.publisher_income),
    ecpm: fenToYuan(stat.ecpm)
  };
}

// ==================== 主流程 ====================

async function main() {
  console.log('获取 access_token...');
  const token = await getToken();
  console.log('Token 获取成功');
  
  // 广告位清单
  console.log('\n获取广告位清单...');
  const adunits = await getAdunitList(token);
  console.log(`广告位数量: ${adunits.length}`);
  adunits.forEach(u => {
    console.log(`  ${u.ad_unit_name} | ${AD_SLOT_NAMES[u.ad_slot] || u.ad_slot} | ${u.ad_unit_status === 'AD_UNIT_STATUS_ON' ? '正常' : '暂停'}`);
  });
  
  // 汇总数据（最近30天）
  const endDate = fmtDate(new Date());
  const startDate = fmtDate(new Date(Date.now() - 30 * 24 * 3600 * 1000));
  console.log(`\n获取汇总数据 (${startDate} ~ ${endDate})...`);
  const summary = await getSummaryData(token, startDate, endDate);
  console.log(`汇总数据条数: ${summary.length}`);
  
  // 计算总收入（元）
  const totalIncome = summary.reduce((sum, item) => sum + fenToYuan(item.income), 0);
  console.log(`总收入: ${totalIncome.toFixed(2)} 元`);
  
  // 细分数据
  console.log(`\n获取细分数据 (${startDate} ~ ${endDate})...`);
  const detail = await getDetailData(token, startDate, endDate);
  console.log(`细分数据条数: ${detail.length}`);
  
  // 结算数据
  console.log('\n获取结算数据...');
  const settlement = await getSettlementData(token);
  console.log(`总预估收入: ${fenToYuan(settlement.revenue_all).toFixed(2)} 元`);
  console.log(`总已结算收入: ${fenToYuan(settlement.settled_revenue_all).toFixed(2)} 元`);
}

main().catch(err => console.error('错误:', err.message));
```

---

## 十一、注意事项与常见问题

### Q1: 请求返回 `ret: 45010 not supported action`

**原因**: 请求方法错误。本 API 系列只支持 **GET** 请求，参数通过 Query String 传递。  
**解决**: 确保使用 GET 方法，不要用 POST。

### Q2: 请求返回 `errcode: 48001 api unauthorized`

**原因**: 当前小程序未开通该 API 权限。  
**解决**: 登录微信公众平台 → 小程序后台 → 流量主 → 开通数据 API 权限。

### Q3: 收入数据看起来异常大或异常小

**原因**: 金额单位是**分**，不是元。`income: 41154` 表示 411.54 元。  
**解决**: 所有金额字段除以 100。

### Q4: 同一 slot_id 出现多条数据

**原因**: `SLOT_ID_WEAPP_TEMPLATE`（原生模版）类型下可能包含多个 ad_unit（如"其他模版广告"和"扫码原生模版广告"），它们共享同一个 slot_id。  
**解决**: 用 `ad_unit_id` 区分，不要仅靠 `slot_id` 去重。

### Q5: 数据最多能查多久之前的？

微信广告数据 API 存储最早 **2016年1月1日** 起的数据。

### Q6: access_token 过期了怎么办？

access_token 有效期 2 小时。建议：
- 缓存 token 和获取时间
- 每次调用前检查是否过期
- 过期前 5 分钟主动刷新
- 若请求返回 `errcode: 42001`，说明 token 已过期，需重新获取

### Q7: 请求频率有限制吗？

有。建议每次请求间隔 ≥ 500ms。超频返回 `ret: 45009`。  
推荐使用指数退避重试策略：500ms → 1s → 2s → 4s → 放弃。

---

## 十二、API 调用速查表

| 数据类型 | action | 必需参数 | 请求方法 |
|---------|--------|---------|---------|
| 广告位清单 | `get_adunit_list` | - | GET |
| 汇总数据 | `publisher_adpos_general` | start_date, end_date | GET |
| 细分数据 | `publisher_adunit_general` | start_date, end_date | GET |
| 结算数据 | `publisher_settlement` | - | GET |

| 通用参数 | 说明 |
|---------|------|
| access_token | 鉴权令牌 |
| page | 页码（从1开始） |
| page_size | 每页条数（最大90） |

---

*文档结束。如有疑问请联系产品/运营同学确认需求细节。*
