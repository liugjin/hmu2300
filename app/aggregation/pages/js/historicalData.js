let typeData = [{ name: '温湿度', typeid: '206', typeIconUrl: '../pic/collectionPage/humidity.svg', typeImg: '../pic/collectionPage/humiture.png', devicelist: [], }, { name: 'ups', typeid: '401', typeIconUrl: '../pic/collectionPage/UPS.svg', typeImg: '../pic/collectionPage/UPS.png', devicelist: [], },
{ name: '漏电检测', typeid: '415', typeIconUrl: '../pic/collectionPage/leakageDetection.svg', typeImg: '../pic/collectionPage/leakageDetection.png', devicelist: [], }, { name: '电表', typeid: '413', typeIconUrl: '../pic/collectionPage/wattHourMeter.svg', typeImg: '../pic/collectionPage/wattHourMeter.png', devicelist: [], },
{ name: '开关电源', typeid: '416', typeIconUrl: '../pic/collectionPage/STS.svg', typeImg: '../pic/collectionPage/STS.png', devicelist: [], }, { name: '空调', typeid: '402', typeIconUrl: '../pic/collectionPage/airConditioner.svg', typeImg: '../pic/collectionPage/airConditioner.png', devicelist: [], },
{ name: '烟感', typeid: '203', typeIconUrl: '../pic/collectionPage/smokeSensation.svg', typeImg: '../pic/collectionPage/smokeSensation.png', devicelist: [], }, { name: '水浸', typeid: '205', typeIconUrl: '../pic/collectionPage/waterOut.svg', typeImg: '../pic/collectionPage/waterOut.png', devicelist: [], },
{ name: '电子锁', typeid: '605', typeIconUrl: '../pic/collectionPage/electronicLock.svg', typeImg: '../pic/collectionPage/electronicLock.png', devicelist: [], }, { name: '防雷器', typeid: '419', typeIconUrl: '../pic/collectionPage/lightningArrester.svg', typeImg: '../pic/collectionPage/lightningArrester.png', devicelist: [], },
{ name: '摄像头', typeid: '601', typeIconUrl: '../pic/collectionPage/camera.svg', typeImg: '../pic/collectionPage/camera.png', devicelist: [], }, { name: 'PDU', typeid: '404', typeIconUrl: '../pic/collectionPage/PDU.svg', typeImg: '../pic/collectionPage/PDU.png', devicelist: [], },
{ name: '逆变器', typeid: '417', typeIconUrl: '../pic/collectionPage/UPS.svg', typeImg: '../pic/collectionPage/inverter2.png', devicelist: [], }, { name: '锂电池', typeid: '409', typeIconUrl: '../pic/collectionPage/UPS.svg', typeImg: '../pic/collectionPage/lidanci.png', devicelist: [], },]
let devicelist = []
let devicelistData = {}

function execute() {
    new Vue({
        el: '#body-box',
        data() {
            return {
                typeData: typeData,
                devicelist: typeData[0].devicelist,
                devicelistData: devicelistData,
                vedioList: [],
                historicalData: [],
                typeSubscript: 0,
                startDate: null,
                endDate: null,
                pageIndex: 0,
                newPageIndex: 1,
                searchid: "",
                queryNumber: 10,
                totalNumber: 0,
                pagesNumber: 0
            }
        },
        mounted() {
            this.addDate()
            this.getDeviceType()
        },
        methods: {
            // 获取当天日期
            addDate() {
                let nowDate = new Date();
                let startDate = {
                    year: nowDate.getFullYear(),
                    month: nowDate.getMonth() + 1,
                    date: nowDate.getDate(),
                }
                // let lastDate1 = new Date(nowDate)
                // let lastDate2 = + lastDate1 + 1000 * 60 * 60 * 24
                // let lastDate3 = new Date(lastDate2)
                // let endDate = {
                //     year: lastDate3.getFullYear(),
                //     month: lastDate3.getMonth() + 1,
                //     date: lastDate3.getDate(),
                // }
                let endDate = {
                    year: nowDate.getFullYear(),
                    month: nowDate.getMonth() + 1,
                    date: nowDate.getDate(),
                }
                this.startDate = startDate.year + "-" + this.judge(startDate.month) + "-" + this.judge(startDate.date);
                this.endDate = endDate.year + "-" + this.judge(endDate.month) + "-" + this.judge(endDate.date);
            },
            // 判断
            judge(data) {
                if (data < 10) {
                    sdata = data.toString()
                    return 0 + sdata
                }
                if (data >= 10) {
                    return data
                }
            },
            // 获取设备类型
            getDeviceType() {
                const vm = this;
                $.ajax({
                    url: `${host}/mu/`,
                    type: 'GET',
                    dataType: 'json',
                    success: function (res) {
                        muport = res.data[0].ports;
                        for (let i = 0; i < muport.length; i++) {
                            for (let j = 0; j < muport[i].sampleUnits.length; j++) {
                                for (let k = 0; k < typeData.length; k++) {
                                    newid = muport[i].sampleUnits[j].id.slice(0, 3)
                                    if (newid.indexOf(typeData[k].typeid) != -1) {
                                        muport[i].sampleUnits[j].typeImg = typeData[k].typeImg
                                        muport[i].sampleUnits[j].typeIconUrl = typeData[k].typeIconUrl
                                        muport[i].sampleUnits[j].typeid = newid
                                        typeData[k].devicelist.push(muport[i].sampleUnits[j])
                                    }
                                }
                            }
                        }
                        vm.setDeviceBox(null, 0)
                    },
                    error: function () {
                        alert("异常");
                    }
                })
            },
            // 获取对应类型的设备列表
            getDeviceList(data, index) {
                this.typeSubscript = index
                this.devicelist = data.devicelist
                this.setDeviceBox(this.devicelist[0], 0)
            },
            // 设置设备根据类型调用不同设备
            setDeviceBox(data, index) {
                this.listSubscript = index
                this.pageIndex = 0
                this.searchid = ""
                data == null ? this.devicelistData = this.devicelist[0] : this.devicelistData = data
                if (this.devicelistData.typeid != '601') {
                    this.getHistoricalData()
                } else {
                    this.getpicture()
                }
            },
            //获取历史数据
            getHistoricalData() {
                const vm = this
                if(vm.historicalData){
                    vm.historicalData.splice(0, vm.historicalData.length)
                }
                let startDate = vm.startDate.split("-")[0] + vm.startDate.split("-")[1] + vm.startDate.split("-")[2]
                let endDate = vm.endDate.split("-")[0] + vm.endDate.split("-")[1] + vm.endDate.split("-")[2]
                let date1 = new Date(parseInt(vm.startDate.split("-")[0]), parseInt(vm.startDate.split("-")[1]) - 1, parseInt(vm.startDate.split("-")[2]), 0, 0, 0)
                let date2 = new Date(parseInt(vm.endDate.split("-")[0]), parseInt(vm.endDate.split("-")[1]) - 1, parseInt(vm.endDate.split("-")[2]), 0, 0, 0)
                if (date1.getTime() > date2.getTime()) {
                    alert('结束日期不能小于开始日期')
                    return
                } else {
                    let parameter = {
                        pageIndex: vm.pageIndex,
                        recordShow: vm.queryNumber,
                        deviceId: vm.searchid == "" ? vm.devicelistData.id : vm.searchid,
                        startDate: startDate,
                        endDate: endDate
                    }
                    $.post(`${host}/getHistoryDataDeviceIDByDate`, parameter, function (data) {
                        if (data) {
                            vm.pagesNumber = Math.ceil(data.count / vm.queryNumber)
                            vm.totalNumber = data.count
                            vm.historicalData = data.deviceRecords
                        }
                    });
                }
            },
            // 获取摄像头抓拍图片
            getpicture() {
                const vm = this
                cameraHoat = vm.devicelistData.setting.host;
                $.post(`${host}/capture/showPicName`, { host: cameraHoat }, function (data) {
                    if (data) {
                        vm.vedioList = data.sList.map(item => {
                            const arr_result = item.split("-");
                            let initTime = arr_result[1].split(".")[0];
                            let time = `${initTime.slice(0, 4)}-${initTime.slice(
                                4,
                                6
                            )}-${initTime.slice(6, 8)} ${initTime.slice(8, 10)}:${initTime.slice(
                                10,
                                12
                            )}:${initTime.slice(12, 14)}`;
                            return {
                                img: host + "/capture/" + item,
                                time: time
                            };
                        });
                    }
                });
            },
            //翻页
            pageTurning(data) {
                const vm = this
                if (data == "upper") {
                    if (vm.newPageIndex == 1) {
                        alert("当前已是第一页")
                        return
                    }
                    vm.pageIndex -= 1
                    vm.newPageIndex == 1 ? vm.newPageIndex = 1 : vm.newPageIndex = vm.newPageIndex -= 1
                    this.getHistoricalData()
                }
                if (data == "lower") {
                    if (vm.newPageIndex == vm.pagesNumber || vm.totalNumber ==0) {
                        alert("当前已是最后一页")
                        return
                    }
                    vm.pageIndex += 1
                    vm.newPageIndex == vm.pagesNumber ? vm.newPageIndex = vm.pagesNumber : vm.newPageIndex = vm.newPageIndex += 1
                    this.getHistoricalData()
                }
            },
            //搜索
            search() {
                if (this.searchid == "") {
                    alert("请输入要搜索的设备ID")
                    return
                } else {
                    this.getHistoricalData()
                }
            },
            //导出
            tableToExcel() {
                const vm = this
                if(!vm.historicalData){
                    alert("当前页面没有数据无法导出")
                    return
                }
                //列标题
                let str = '<tr><td>序号</td><td>设备ID</td><td>数据名称</td><td>通道ID</td><td>数值</td><td>单位</td><td>采集时间</td><td>创建时间</td></tr>';
                //循环遍历，每行加入tr标签，每个单元格加td标签
                for (let i = 0; i < vm.historicalData.length; i++) {
                    str += '<tr>';
                    for (let item in vm.historicalData[i]) {
                        //增加\t为了不让表格显示科学计数法或者其他格式
                        str += `<td>${vm.historicalData[i][item] + '\t'}</td>`;
                    }
                    str += '</tr>';
                }
                //Worksheet名
                let worksheet = 'Sheet1'
                let uri = 'data:application/vnd.ms-excel;base64,';

                //下载的表格模板数据
                let template = `<html xmlns:o="urn:schemas-microsoft-com:office:office" 
                    xmlns:x="urn:schemas-microsoft-com:office:excel" 
                    xmlns="http://www.w3.org/TR/REC-html40">
                    <head><!--[if gte mso 9]><xml><x:ExcelWorkbook><x:ExcelWorksheets><x:ExcelWorksheet>
                        <x:Name>${worksheet}</x:Name>
                        <x:WorksheetOptions><x:DisplayGridlines/></x:WorksheetOptions></x:ExcelWorksheet>
                        </x:ExcelWorksheets></x:ExcelWorkbook></xml><![endif]-->
                        </head><body><table>${str}</table></body></html>`;
                //下载模板
                window.location.href = uri + this.base64(template)

            },
            base64(s) { return window.btoa(unescape(encodeURIComponent(s))) }
        }
    })
}
execute()