
let typeData = [{ name: '温湿度', typeid: '206', typeIconUrl: '../pic/collectionPage/humidity.svg', typeImg: '../pic/collectionPage/humiture.png', devicelist: [], }, { name: 'ups', typeid: '401', typeIconUrl: '../pic/collectionPage/UPS.svg', typeImg: '../pic/collectionPage/UPS.png', devicelist: [], },
{ name: '漏电检测', typeid: '415', typeIconUrl: '../pic/collectionPage/leakageDetection.svg', typeImg: '../pic/collectionPage/leakageDetection.png', devicelist: [], }, { name: '电表', typeid: '413', typeIconUrl: '../pic/collectionPage/wattHourMeter.svg', typeImg: '../pic/collectionPage/wattHourMeter.png', devicelist: [], },
{ name: '开关电源', typeid: '416', typeIconUrl: '../pic/collectionPage/STS.svg', typeImg: '../pic/collectionPage/STS.png', devicelist: [], }, { name: '空调', typeid: '402', typeIconUrl: '../pic/collectionPage/airConditioner.svg', typeImg: '../pic/collectionPage/airConditioner.png', devicelist: [], },
{ name: '烟感', typeid: '203', typeIconUrl: '../pic/collectionPage/smokeSensation.svg', typeImg: '../pic/collectionPage/smokeSensation.png', devicelist: [], }, { name: '水浸', typeid: '205', typeIconUrl: '../pic/collectionPage/waterOut.svg', typeImg: '../pic/collectionPage/waterOut.png', devicelist: [], },
{ name: '电子锁', typeid: '605', typeIconUrl: '../pic/collectionPage/electronicLock.svg', typeImg: '../pic/collectionPage/electronicLock.png', devicelist: [], }, { name: '防雷器', typeid: '419', typeIconUrl: '../pic/collectionPage/lightningArrester.svg', typeImg: '../pic/collectionPage/lightningArrester.png', devicelist: [], },
{ name: '摄像头', typeid: '601', typeIconUrl: '../pic/collectionPage/camera.svg', typeImg: '../pic/collectionPage/camera.png', devicelist: [], }, { name: 'PDU', typeid: '404', typeIconUrl: '../pic/collectionPage/PDU.svg', typeImg: '../pic/collectionPage/PDU.png', devicelist: [], },
{ name: '逆变器', typeid: '417', typeIconUrl: '../pic/collectionPage/UPS.svg', typeImg: '../pic/collectionPage/inverter2.png', devicelist: [], }, { name: '锂电池', typeid: '409', typeIconUrl: '../pic/collectionPage/UPS.svg', typeImg: '../pic/collectionPage/lidanci.png', devicelist: [], },]
let devicelist = []
let devicelistData = {}


// 获取设备类型
function execute() {
    new Vue({
        el: '#body-box',
        data() {
            return {
                typeData: typeData,
                devicelist: typeData[0].devicelist,
                devicelistData: devicelistData,
                channels: [],
                newchannels: [],
                vedioList: [],
                typeSubscript: 0,
                listSubscript: 0,
                isinfo: false,
                status: null,
                gatewayID: null,
            }
        },
        mounted() {
            this.getStatus()
        },
        methods: {
            //获取设备状态
            getStatus() {
                const vm = this;
                $.ajax({
                    url: `${host}/getStatus`,
                    type: 'GET',
                    dataType: 'json',
                    success: function (res) {
                        console.log("res",res)
                        vm.status = JSON.parse(res.bs);
                        vm.getDeviceType()
                    },
                    error: function () {
                        alert("异常");
                    }
                })
            },
            // 获取设备类型
            getDeviceType() {
                const vm = this;
                $.ajax({
                    url: `${host}/mu/`,
                    type: 'GET',
                    dataType: 'json',
                    success: function (res) {
                        vm.gatewayID = res.data[0].id
                        muport = res.data[0].ports;
                        for (let i = 0; i < muport.length; i++) {
                            for (let j = 0; j < muport[i].sampleUnits.length; j++) {
                                for (let k = 0; k < typeData.length; k++) {
                                    newid = muport[i].sampleUnits[j].id.slice(0, 3)
                                    if (newid.indexOf(typeData[k].typeid) != -1) {
                                        muport[i].sampleUnits[j].typeImg = typeData[k].typeImg
                                        muport[i].sampleUnits[j].typeIconUrl = typeData[k].typeIconUrl
                                        muport[i].sampleUnits[j].typeid = newid
                                        muport[i].sampleUnits[j].interfaceType = muport[i].id
                                        vm.setStatus(muport[i].sampleUnits[j])
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
                data == null ? this.devicelistData = this.devicelist[0] : this.devicelistData = data
                this.setStatus(this.devicelistData)
                if (this.devicelistData.typeid != '601') {
                    this.getDeviceData()
                } else {
                    this.getpicture()
                }
            },
            //设置状态
            setStatus(data) {
                const vm = this;
                if (vm.status.length > 0) {
                    for (let i = 0; i < vm.status.length; i++) {
                        if (vm.status[i].suid == data.id) {
                            data.state = vm.status[i].state
                            // if (data.interfaceType.indexOf("di") != -1 || data.interfaceType.indexOf("do") != -1) {
                            //     data.state = -1
                            // } else {
                            //     data.state = vm.status[i].state
                            // }
                        }
                    }
                }
            },
            // 获取设备实时数据
            getDeviceData() {
                const vm = this;
                this.isinfo = false;
                $.ajax({
                    url: `${host}/el/${vm.devicelistData.element}`,
                    type: 'GET',
                    dataType: 'json',
                    success: function (res) {
                        vm.channels = res.data.channels
                        if (window["WebSocket"]) {//定义websocket 采集
                            let wsHost = host.replace("http", "ws"); 
                            let conn = new WebSocket(wsHost + "/real");
                            conn.onclose = function (evt) {
                                console.log("1", evt)
                            };
                            conn.onmessage = function (evt) {
                                let messages = evt.data.split('\n');
                                for (var i = 0; i < messages.length; i++) {
                                    if (messages[i] != "") {
                                        let jsonMessages = JSON.parse(messages[i])
                                        if (jsonMessages.sampleUnitId == vm.devicelistData.id) {
                                            for (let i = 0; i < vm.channels.length; i++) {
                                                if (jsonMessages.channelId == vm.channels[i].id) {
                                                    vm.channels[i].value = jsonMessages.value;
                                                }
                                            }
                                        }
                                    }
                                }
                            };
                        }
                        vm.clickinfo()
                    },
                    error: function () {
                        alert("异常");
                    }
                })
            },
            // 获取摄像头抓拍图片
            getpicture() {
                const vm = this
                cameraHoat = vm.devicelistData.setting.host;
                $.post(`${host}/capture/showPicName`, { host: cameraHoat }, function (data) {
                    console.log("datadata",data)
                    if (data && data.sList.length>0) {
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
            //控制门禁
            controlLock(id) {
                const vm = this
                let commandValue = null
                let data = null
                for (var i = 0; i < vm.channels.length; i++) {
                    if (vm.channels[i].id == "OnOrOffLock") {
                        vm.channels[i].value == "0" ? commandValue = "1" : commandValue = "0"
                        data = {
                            commandMuid: vm.gatewayID,
                            commandUnitId: id,
                            commandChannelId: vm.channels[i].id,
                            commandKey: "value",
                            commandValue: commandValue,
                            commandType: "int"
                        }
                    }
                }
                $.post(`${host}/commandForm`, data, function (data) {
                    if (data) {
                        console.log("data", data)
                        // alert(data+"控制成功")
                    }
                });
            },
            //点击显示更多
            clickinfo() {
                this.isinfo = !this.isinfo
                this.newchannels.splice(0, this.newchannels.length)
                if (this.isinfo) {
                    this.newchannels = this.channels.slice(0, 10)
                } else {
                    this.newchannels = this.channels.slice(0, this.channels.length)

                }
            },
        }
    })
}
execute()


