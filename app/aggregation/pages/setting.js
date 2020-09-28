var collect = [
  {
    type: "collecter",
    typeName:"采集器",
    content: [{
      id: "sensorflow",
      parameter: {
        name: "sensorflow",
        brand: "华远云联",
        model: "AA-B1-1001",
        image: "../pic/vis/sensorflow.svg"
      },
      default: {
        //默认采集周期
        cuperiod: "2000",
        //默认采集超时
        cutimeout: "1000",
        //默认采集延时
        cudelay: "500",
        //默认采集节流时间
        cuthrottle: "500",
        //默认最大通讯错误数
        cumaxNum: "100",
        //默认设备库
        library: "sensorflow.json"
      },
      //默认变量
      variable: {
        address: "1"
      }
    }]
  },
  {
    type: "physic",
    typeName:"物理设备",
    content: [{
      id: "alarm",
      parameter: {
        name: "闹钟",
        brand: "华远云联",
        model: "AA-B1-1002",
        image: "../pic/vis/alarm.svg"
      },
      default: {
        cuperiod: "500",
        cutimeout: "100",
        cudelay: "800",
        cuthrottle: "500",
        cumaxNum: "100",
        library: "alarm.json"
      },
      variable: {
        address: "2"
      }
    },{
      id: "AEW100",
      parameter: {
        name: "AEW100",
        brand: "华远云联",
        model: "AA-B1-1002",
        image: "../pic/vis/alarm.svg"
      },
      default: {
        cuperiod: "500",
        cutimeout: "100",
        cudelay: "800",
        cuthrottle: "500",
        cumaxNum: "100",
        library: "alarm.json"
      },
      variable: {
        address: "2"
      }
    }]
  },
  {
    type: "huayuan",
    typeName:"华远",
    content: [{
      id: "hmu2000",
      parameter: {
        name: "hmu2000",
        brand: "华远云联",
        model: "AA-B1-1002",
        image: "../pic/vis/hmu2000.svg"
      },
      default: {
        cuperiod: "500",
        cutimeout: "100",
        cudelay: "800",
        cuthrottle: "500",
        cumaxNum: "100",
        library: "hmu2000.json"
      },
      variable: {
        address: "2"
      }
    }]
  }
];

var netEquipment = [
  {
    type: "camera",
    typeName:"海康摄像机",
    content: [{
      id: "hikiVision",
      parameter: {
        name: "hiki",
        brand: "海康摄像机",
        model: "Camera-1001",
        image: "../pic/vis/hiki.svg",
        rtspRule: "rtsp://[username]:[password]@[ip]:[port]/[codec]/[channel]/[subtype]/av_stream"
        /* 
        username: 用户名。例如admin。
        password: 密码。例如12345。
        ip: 为设备IP。例如 192.0.0.64。
        port: 端口号默认为554，若为默认可不填写。
        codec：有h264、MPEG-4、mpeg4这几种。
        channel: 通道号，起始为1。例如通道1，则为ch1。
        subtype: 码流类型，主码流为main，辅码流为sub。 
        */
      },
      default: {
        //默认采集周期
        cuperiod: "2000",
        //默认采集超时
        cutimeout: "1000",
        //默认采集延时
        cudelay: "500",
        //默认采集节流时间
        cuthrottle: "500",
        //默认最大通讯错误数
        cumaxNum: "100",
        //默认设备库
        library: "hiki.json"
      },
      //默认变量
      variable: {
        address: "1"
      }
    },{
      id: "huayuanVision",
      parameter: {
        name: "Huayuan",
        brand: "华远云联摄像机",
        model: "Camera-Huayuan",
        image: "../pic/vis/huayuan.svg",
        rtspRule: "rtsp://username:password@ip:port/cam/realmonitor?channel=1&subtype=0"
        /* 
        username: 用户名。例如admin。
        password: 密码。例如admin。
        ip: 为设备IP。例如 10.7.8.122。
        port: 端口号默认为554，若为默认可不填写。
        channel: 通道号，起始为1。例如通道2，则为channel=2。
        subtype: 码流类型，主码流为0（即subtype=0），辅码流为1（即subtype=1） 
        */
      },
      default: {
        //默认采集周期
        cuperiod: "2000",
        //默认采集超时
        cutimeout: "1000",
        //默认采集延时
        cudelay: "500",
        //默认采集节流时间
        cuthrottle: "500",
        //默认最大通讯错误数
        cumaxNum: "100",
        //默认设备库
        library: "Camera-Huayuan.json"
      },
      //默认变量
      variable: {
        address: "1"
      }
    }]
  },{
    type: "flask",
    typeName:"随便添加的",
    content: [{
      id: "random",
      parameter: {
        name: "random",
        brand: "随便测试",
        model: "random-1001",
        image: "../pic/vis/random.svg"
      },
      default: {
        //默认采集周期
        cuperiod: "2000",
        //默认采集超时
        cutimeout: "1000",
        //默认采集延时
        cudelay: "500",
        //默认采集节流时间
        cuthrottle: "500",
        //默认最大通讯错误数
        cumaxNum: "100",
        //默认设备库
        library: "sensorflow.json"
      },
      //默认变量
      variable: {
        address: "1"
      }
    }]
  },
];
