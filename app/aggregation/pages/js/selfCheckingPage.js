let typeData = [{ name: '温湿度', typeid: '206', typeIconUrl: '../pic/collectionPage/humidity.svg', typeImg: '../pic/collectionPage/humiture.png', devicelist: [], }, { name: 'ups', typeid: '401', typeIconUrl: '../pic/collectionPage/UPS.svg', typeImg: '../pic/collectionPage/UPS.png', devicelist: [], },
{ name: '漏电检测', typeid: '415', typeIconUrl: '../pic/collectionPage/leakageDetection.svg', typeImg: '../pic/collectionPage/leakageDetection.png', devicelist: [], }, { name: '电表', typeid: '413', typeIconUrl: '../pic/collectionPage/wattHourMeter.svg', typeImg: '../pic/collectionPage/wattHourMeter.png', devicelist: [], },
{ name: '开关电源', typeid: '416', typeIconUrl: '../pic/collectionPage/STS.svg', typeImg: '../pic/collectionPage/STS.png', devicelist: [], }, { name: '空调', typeid: '402', typeIconUrl: '../pic/collectionPage/airConditioner.svg', typeImg: '../pic/collectionPage/airConditioner.png', devicelist: [], },
{ name: '烟感', typeid: '203', typeIconUrl: '../pic/collectionPage/smokeSensation.svg', typeImg: '../pic/collectionPage/smokeSensation.png', devicelist: [], }, { name: '水浸', typeid: '205', typeIconUrl: '../pic/collectionPage/waterOut.svg', typeImg: '../pic/collectionPage/waterOut.png', devicelist: [], },
{ name: '电子锁', typeid: '605', typeIconUrl: '../pic/collectionPage/electronicLock.svg', typeImg: '../pic/collectionPage/electronicLock.png', devicelist: [], }, { name: '防雷器', typeid: '419', typeIconUrl: '../pic/collectionPage/lightningArrester.svg', typeImg: '../pic/collectionPage/lightningArrester.png', devicelist: [], },
{ name: '摄像头', typeid: '601', typeIconUrl: '../pic/collectionPage/camera.svg', typeImg: '../pic/collectionPage/camera.png', devicelist: [], }, { name: 'PDU', typeid: '404', typeIconUrl: '../pic/collectionPage/PDU.svg', typeImg: '../pic/collectionPage/PDU.png', devicelist: [], },
{ name: '逆变器', typeid: '417', typeIconUrl: '../pic/collectionPage/UPS.svg', typeImg: '../pic/collectionPage/inverter2.png', devicelist: [], }, { name: '锂电池', typeid: '409', typeIconUrl: '../pic/collectionPage/UPS.svg', typeImg: '../pic/collectionPage/lidanci.png', devicelist: [], },]
let deviceList = []
let deviceStater = []
function execute() {
    new Vue({
        el: '#body-box',
        data() {
            return {
                selfscan: [],
                start: false,
                barLength: 0,
                interval: null
            }
        },
        mounted() {

        },
        methods: {
            // 获取设备列表
            getDeviceList() {
                const vm = this
                $.ajax({
                    url: `${host}/mu/`,
                    type: 'GET',
                    dataType: 'json',
                    success: function (res) {
                        muport = res.data[0].ports;
                        deviceList.splice(0, deviceList.length)
                        for (let i = 0; i < muport.length; i++) {
                            // if (muport[i].id.indexOf("di") != -1 || muport[i].id.indexOf("do") != -1) {
                            //     continue
                            // } else {
                                
                            // }
                            for (let j = 0; j < muport[i].sampleUnits.length; j++) {
                                for (let k = 0; k < typeData.length; k++) {
                                    newid = muport[i].sampleUnits[j].id.slice(0, 3)
                                    if (newid.indexOf(typeData[k].typeid) != -1) {
                                        muport[i].sampleUnits[j].icon = typeData[k].typeIconUrl
                                        deviceList.push(muport[i].sampleUnits[j])
                                    }
                                }
                            }
                        }
                        vm.getStatus()
                    },
                    error: function () {
                        alert("异常");
                    }
                })
            },
            // 获取各设备状态
            getStatus() {
                const vm = this
                $.ajax({
                    url: `${host}/getStatus`,
                    type: 'GET',
                    dataType: 'json',
                    success: function (res) {
                        muport = JSON.parse(res.bs);
                        deviceStater.splice(0, deviceStater.length)
                        for (let i = 0; i < muport.length; i++) {
                            for (let j = 0; j < deviceList.length; j++) {
                                if (muport[i].suid == deviceList[j].id) {
                                    let stateData = { state: 0, name: null, icon: null };
                                    stateData.state = muport[i].state
                                    stateData.name = deviceList[j].name
                                    stateData.icon = deviceList[j].icon
                                    deviceStater.push(stateData)
                                }
                            }
                        }
                        vm.startSelfscan()
                    },
                    error: function () {
                        alert("异常");
                    }
                })
            },
            // 开始自检
            startSelfscan() {
                const vm = this
                let i = 0
                vm.start = true
                vm.barLength = 0
                vm.selfscan.splice(0, vm.selfscan.length)
                clearInterval(vm.interval)
                vm.interval = setInterval(function () {
                    vm.selfscan.push(deviceStater[i])
                    vm.barLength = vm.setbar(deviceStater.length, vm.selfscan.length) //自检进度条
                    i = i + 1
                    if (i >= deviceStater.length) {
                        clearInterval(vm.interval)
                    }
                }, 2000)
            },
            setbar(originalLength, setafterLength) {
                return Number((setafterLength / originalLength * 100).toFixed(0))
            },
        }

    })
}
execute()