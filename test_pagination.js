console.log('========== 微信广告数据API分页逻辑测试 ==========\n');

let mockCallCount = 0;
let currentResponseIndex = 0;

const mockResponses = [
    { list: Array(90).fill(null).map((_, i) => ({id: i + 1, date: '2026-04-01'})), total_num: 250 },
    { list: Array(90).fill(null).map((_, i) => ({id: 90 + i + 1, date: '2026-04-01'})), total_num: 250 },
    { list: Array(70).fill(null).map((_, i) => ({id: 180 + i + 1, date: '2026-04-01'})), total_num: 250 }
];

function sleep(ms) {
    return new Promise(r => setTimeout(r, ms));
}

async function mockHttpGet(url) {
    mockCallCount++;
    console.log(`  [模拟API调用 #${mockCallCount}]`);

    const pageMatch = url.match(/page=(\d+)/);
    const page = pageMatch ? parseInt(pageMatch[1]) : 1;

    console.log(`    URL: ${url}`);
    console.log(`    返回第 ${page} 页数据`);

    const response = mockResponses[page - 1] || { list: [], total_num: 0 };
    console.log(`    数据量: ${response.list.length} 条`);
    console.log(`    总记录数: ${response.total_num}`);

    await sleep(100);
    return response;
}

async function fetchDataWithRetry(token, action, appid, appsecret, extraParams = {}) {
    const params = new URLSearchParams({
        access_token: token,
        action,
        page: extraParams.page || '1',
        page_size: extraParams.page_size || '90',
    });

    Object.entries(extraParams).forEach(([k, v]) => {
        if (k !== 'page' && k !== 'page_size') {
            params.set(k, String(v));
        }
    });

    const url = `https://api.weixin.qq.com/publisher/stat?${params}`;
    return await mockHttpGet(url);
}

async function fetchAllPagesByDateRange(token, action, startDate, endDate, appid, appsecret) {
    let allList = [];
    let page = 1;
    const pageSize = 90;
    let hasMore = true;
    let lastData = null;

    while (hasMore) {
        const data = await fetchDataWithRetry(token, action, appid, appsecret, {
            start_date: startDate,
            end_date: endDate,
            page: String(page),
            page_size: String(pageSize)
        });

        lastData = data;
        const list = data.list || [];
        allList = allList.concat(list);

        const totalNum = data.total_num || 0;
        if (page * pageSize >= totalNum || list.length < pageSize) {
            hasMore = false;
        } else {
            page++;
            await sleep(500);
        }
    }

    return { list: allList, totalNum: lastData?.total_num || allList.length };
}

async function runTests() {
    console.log('测试场景：模拟获取250条数据（每页90条，需要3页）\n');

    mockCallCount = 0;
    const result = await fetchAllPagesByDateRange(
        'test_token',
        'publisher_adpos_general',
        '2026-04-01',
        '2026-04-10',
        'test_appid',
        'test_appsecret'
    );

    console.log('\n========== 测试结果 ==========');
    console.log(`API调用次数: ${mockCallCount}`);
    console.log(`获取的总记录数: ${result.list.length}`);
    console.log(`API返回的总记录数: ${result.totalNum}`);

    const test1Pass = result.list.length === 250 && mockCallCount === 3;
    console.log(`\n${test1Pass ? '✅' : '❌'} 测试1 - 大数据量分页:`);
    console.log(`   预期: 250条数据, 3次API调用`);
    console.log(`   实际: ${result.list.length}条数据, ${mockCallCount}次API调用`);

    console.log('\n========== 边界测试 ==========\n');

    const testCases = [
        {
            name: '空数据',
            responses: [{ list: [], total_num: 0 }],
            expectedCount: 0,
            expectedCalls: 1
        },
        {
            name: '单页数据（<90条）',
            responses: [{ list: Array(50).fill({id: 1}), total_num: 50 }],
            expectedCount: 50,
            expectedCalls: 1
        },
        {
            name: '刚好90条（1页）',
            responses: [{ list: Array(90).fill({id: 1}), total_num: 90 }],
            expectedCount: 90,
            expectedCalls: 1
        },
        {
            name: '91条数据（2页）',
            responses: [
                { list: Array(90).fill({id: 1}), total_num: 91 },
                { list: [{id: 91}], total_num: 91 }
            ],
            expectedCount: 91,
            expectedCalls: 2
        },
        {
            name: '180条数据（2页）',
            responses: [
                { list: Array(90).fill({id: 1}), total_num: 180 },
                { list: Array(90).fill({id: 2}), total_num: 180 }
            ],
            expectedCount: 180,
            expectedCalls: 2
        },
        {
            name: '1000条数据（12页）',
            responses: Array(12).fill(null).map((_, i) => ({
                list: i < 11 ? Array(90).fill({id: i}) : Array(10).fill({id: i}),
                total_num: 1000
            })),
            expectedCount: 1000,
            expectedCalls: 12
        }
    ];

    let allPassed = test1Pass;

    for (const testCase of testCases) {
        mockResponses.length = 0;
        mockResponses.push(...testCase.responses);
        mockCallCount = 0;

        const result = await fetchAllPagesByDateRange(
            'token', 'action', '2026-04-01', '2026-04-10', 'appid', 'secret'
        );

        const countMatch = result.list.length === testCase.expectedCount;
        const callsMatch = mockCallCount === testCase.expectedCalls;
        const passed = countMatch && callsMatch;

        allPassed = allPassed && passed;

        console.log(`${passed ? '✅' : '❌'} ${testCase.name}:`);
        console.log(`   数据: ${result.list.length}/${testCase.expectedCount} 条`);
        console.log(`   API: ${mockCallCount}/${testCase.expectedCalls} 次`);

        if (!passed) {
            console.log(`   原因: ${!countMatch ? '数据量不匹配 ' : ''}${!callsMatch ? 'API调用次数不正确' : ''}`);
        }
    }

    console.log('\n========== 最终结果 ==========');
    if (allPassed) {
        console.log('✅ 所有测试通过！分页逻辑工作正常。\n');
    } else {
        console.log('❌ 部分测试失败，请检查分页逻辑。\n');
    }

    return allPassed;
}

runTests().then(passed => {
    process.exit(passed ? 0 : 1);
}).catch(err => {
    console.error('测试执行错误:', err);
    process.exit(1);
});
