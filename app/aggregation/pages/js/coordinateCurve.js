
// 获取是否有配置文件
function getcfgok() {
    let xData = []
    let value = {}
    let yAxis = { type: 'value', min: null, max: null }
    if (window["WebSocket"]) {//定义websocket 采集
        let wsHost = host.replace("http", "ws");
        let conn = new WebSocket(wsHost + "/real");
        conn.onclose = function (evt) {
            console.log("1", evt)
        };
        conn.onmessage = function (evt) {
            let messages = evt.data.split('\n');
            console.log("messages", messages)
            for (var i = 0; i < messages.length; i++) {
                if (messages[i] != "") {
                    let jsonMessages = JSON.parse(messages[i])
                    console.log("jsonMessages",jsonMessages)
                    if (jsonMessages.sampleUnitId == "acce") {
                        if (jsonMessages.channelId == "freq") {
                            xData.splice(0, xData.length)
                            for (var i = 1; i <= 5 * jsonMessages.value; i++) {
                                xData.push(i)
                            }
                            console.log("xData", xData)
                        }
                        if (jsonMessages.channelId == "points") {
                            value = JSON.parse(jsonMessages.value)
                            let data = setValue(value.x, value.y, value.z)
                            if (data.max < 500 || data.min > -500) {
                                yAxis.max = 500
                                yAxis.min = -500
                            }
                            if (data.max >= 500 || data.min <= -500) {
                                yAxis.max = null
                                yAxis.min = null
                            }
                            console.log("data", data)
                        }
                        echart(value, xData, yAxis)
                    }
                }
            }
        };
    }
}


function setValue(xArr, yArr, zArr) {
    c = xArr.concat(yArr, zArr)
    max = Math.max.apply(null, c)
    min = Math.min.apply(null, c)
    return yMaxMin = { max: max, min: min }
}



function echart(value, xData, yAxis) {
    dom = document.getElementById("vibration");
    myChart = echarts.init(dom)
    console.log("value2", value)
    option = {
        title: {
            text: '振动曲线'
        },
        tooltip: {
            trigger: 'axis'
        },
        legend: {
            data: ['x轴', 'y轴', 'z轴',]
        },
        grid: {
            left: '3%',
            right: '4%',
            bottom: '3%',
            containLabel: true
        },
        toolbox: {
            feature: {
                saveAsImage: {}
            }
        },
        xAxis: {
            type: 'category',
            boundaryGap: false,
            data: xData
        },
        yAxis: yAxis,
        series: [
            {
                name: 'x轴',
                type: 'line',
                stack: '总量',
                data: value.x
            },
            {
                name: 'y轴',
                type: 'line',
                stack: '总量',
                data: value.y
            },
            {
                name: 'z轴',
                type: 'line',
                stack: '总量',
                data: value.z
            }
        ]
    };
    myChart.setOption(option);
}
$(document).ready(function () {
    getcfgok()
    console.log("函数调用");
});