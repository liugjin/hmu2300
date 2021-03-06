var json = {
    "userChoose": [{
        "id": "deviceType",
        "name": "设备类型"
    }, {
        "id": "protocolType",
        "name": "协议类型"
    }, {
        "id": "brandType",
        "name": "品牌"
    }],
    "deviceType": [{
            "id": "humiture",
            "name": "温湿度"
        },
        {
            "id": "air-conditioner",
            "name": "空调"
        },
        {
            "id": "ups",
            "name": "UPS"
        },
        {
            "id": "converter",
            "name": "变频器"
        },
        {
            "id": "power-meter",
            "name": "电力仪表"
        },
        {
            "id": "sensorflow",
            "name": "sensorflow"
        },
        {
            "id": "camera",
            "name": "摄像机"
        },
        {
            "id": "other",
            "name": "其它"
        }
    ],
    "protocolType": [{
            "id": "modbus",
            "name": "MODBUS"
        },
        {
            "id": "pmbus",
            "name": "电总"
        },
        {
            "id": "snmp",
            "name": "SNMP"
        },
        {
            "id": "video",
            "name": "Video"
        },
        {
            "id": "other",
            "name": "其它"
        }
    ],
    "brandType": [{
            "id": "best",
            "name": "百斯特"
        },
        {
            "id": "emerson",
            "name": "艾默生"
        },
        {
            "id": "hyiot",
            "name": "华远云联"
        },
        {
            "id": "yada",
            "name": "雅达"
        },
        {
            "id": "iteaq",
            "name": "艾特网能"
        },
        {
            "id": "zhenyang",
            "name": "真扬"
        },
        {
            "id": "elebest",
            "name": "贝斯特"
        },
        {
            "id": "hy-electric",
            "name": "华远电气"
        },
        {
            "id": "acrel",
            "name": "安科瑞"
        },
        {
            "id": "hiki",
            "name": "海康"
        },
        {
            "id": "other",
            "name": "其它"
        }
    ],
    "serialEquipment": [{
            "parameters": {
                "deviceType": {
                    "id": "humiture",
                    "name": "温湿度"
                },
                "protocolType": {
                    "id": "modbus",
                    "name": "MODBUS"
                },
                "brandType": {
                    "id": "best",
                    "name": "百斯特"
                }
            },
            "content": {
                "id": "best-th",
                "name": "百斯特温湿度",
                "model": "THD-0102",
                "image": "../pic/vis/best-th.svg",
                "desc": "百斯特温湿度",
                "default": {
                    "samplePort": {
                        "protocol": "protocol-modbus-serial",
                        "setting": {
                            "baudRate": "9600"
                        }
                    },
                    "sampleUnit": {
                        "cuperiod": "3000",
                        "cutimeout": "1000",
                        "cudelay": "500",
                        "cuthrottle": "500",
                        "cumaxNum": "5",
                        "library": "best-th.json",
                        "setting": {
                            "address": "1"
                        }
                    }
                }
            }
        },
        {
            "parameters": {
                "deviceType": {
                    "id": "humiture",
                    "name": "温湿度"
                },
                "protocolType": {
                    "id": "modbus",
                    "name": "MODBUS"
                },
                "brandType": {
                    "id": "emerson",
                    "name": "艾默生"
                }
            },
            "content": {
                "id": "emerson-th",
                "name": "艾默生温湿度",
                "model": "IRM-S02TH",
                "image": "../pic/vis/emerson-th.svg",
                "desc": "艾默生温湿度",
                "default": {
                    "samplePort": {
                        "protocol": "protocol-modbus-serial",
                        "setting": {
                            "baudRate": "9600"
                        }
                    },
                    "sampleUnit": {
                        "cuperiod": "3000",
                        "cutimeout": "1000",
                        "cudelay": "500",
                        "cuthrottle": "500",
                        "cumaxNum": "5",
                        "library": "emerson-th.json",
                        "setting": {
                            "address": "1"
                        }
                    }
                }
            }
        },
        {
            "parameters": {
                "deviceType": {
                    "id": "humiture",
                    "name": "温湿度"
                },
                "protocolType": {
                    "id": "modbus",
                    "name": "MODBUS"
                },
                "brandType": {
                    "id": "yada",
                    "name": "雅达"
                }
            },
            "content": {
                "id": "yada-th",
                "name": "雅达温湿度",
                "model": "YD8771Y",
                "image": "../pic/vis/yada-th.svg",
                "desc": "雅达温湿度",
                "default": {
                    "samplePort": {
                        "protocol": "protocol-modbus-serial",
                        "setting": {
                            "baudRate": "9600"
                        }
                    },
                    "sampleUnit": {
                        "cuperiod": "3000",
                        "cutimeout": "1000",
                        "cudelay": "500",
                        "cuthrottle": "500",
                        "cumaxNum": "5",
                        "library": "yada-th.json",
                        "setting": {
                            "address": "1"
                        }
                    }
                }
            }
        },
        {
            "parameters": {
                "deviceType": {
                    "id": "air-conditioner",
                    "name": "空调"
                },
                "protocolType": {
                    "id": "modbus",
                    "name": "MODBUS"
                },
                "brandType": {
                    "id": "emerson",
                    "name": "艾默生"
                }
            },
            "content": {
                "id": "ac-pex",
                "name": "PEX空调",
                "model": "ac-pex",
                "image": "../pic/vis/ac-pex.svg",
                "desc": "艾默生PEX空调",
                "default": {
                    "samplePort": {
                        "protocol": "protocol-modbus-serial",
                        "setting": {
                            "baudRate": "9600"
                        }
                    },
                    "sampleUnit": {
                        "cuperiod": "3000",
                        "cutimeout": "1000",
                        "cudelay": "500",
                        "cuthrottle": "500",
                        "cumaxNum": "5",
                        "library": "ac-pex.json",
                        "setting": {
                            "address": "1"
                        }
                    }
                }
            }
        },
        {
            "parameters": {
                "deviceType": {
                    "id": "air-conditioner",
                    "name": "空调"
                },
                "protocolType": {
                    "id": "modbus",
                    "name": "MODBUS"
                },
                "brandType": {
                    "id": "iteaq",
                    "name": "艾特网能"
                }
            },
            "content": {
                "id": "cool-row5000",
                "name": "艾特网能Cool Row5000(金戈)变频空调",
                "model": "cool-row5000",
                "image": "../pic/vis/cool-row5000.svg",
                "desc": "艾特网能Cool Row5000(金戈)变频空调",
                "default": {
                    "samplePort": {
                        "protocol": "protocol-modbus-serial",
                        "setting": {
                            "baudRate": "9600"
                        }
                    },
                    "sampleUnit": {
                        "cuperiod": "3000",
                        "cutimeout": "1000",
                        "cudelay": "500",
                        "cuthrottle": "500",
                        "cumaxNum": "5",
                        "library": "cool-row5000-frequency-conversion.json",
                        "setting": {
                            "address": "1"
                        }
                    }
                }
            }
        },
        {
            "parameters": {
                "deviceType": {
                    "id": "air-conditioner",
                    "name": "空调"
                },
                "protocolType": {
                    "id": "pmbus",
                    "name": "电总"
                },
                "brandType": {
                    "id": "emerson",
                    "name": "艾默生"
                }
            },
            "content": {
                "id": "acm03u1",
                "name": "ACM03U1空调",
                "model": "ACM03U1",
                "image": "../pic/vis/acm03u1.svg",
                "desc": "艾默生ACM03U1空调",
                "default": {
                    "samplePort": {
                        "protocol": "protocol-pmbus",
                        "setting": {
                            "baudRate": "9600"
                        }
                    },
                    "sampleUnit": {
                        "cuperiod": "3000",
                        "cutimeout": "1000",
                        "cudelay": "500",
                        "cuthrottle": "500",
                        "cumaxNum": "5",
                        "library": "ACM03U1.json",
                        "setting": {
                            "address": "1"
                        }
                    }
                }
            }
        },
        {
            "parameters": {
                "deviceType": {
                    "id": "ups",
                    "name": "UPS"
                },
                "protocolType": {
                    "id": "pmbus",
                    "name": "电总"
                },
                "brandType": {
                    "id": "emerson",
                    "name": "艾默生"
                }
            },
            "content": {
                "id": "gxe-ups",
                "name": "GXE-UPS",
                "model": "GXE-UPS",
                "image": "../pic/vis/GXE-UPS.svg",
                "desc": "艾默生GXE-UPS",
                "default": {
                    "samplePort": {
                        "protocol": "protocol-pmbus",
                        "setting": {
                            "baudRate": "9600"
                        }
                    },
                    "sampleUnit": {
                        "cuperiod": "3000",
                        "cutimeout": "1000",
                        "cudelay": "500",
                        "cuthrottle": "500",
                        "cumaxNum": "5",
                        "library": "GXE-UPS.json",
                        "setting": {
                            "address": "1"
                        }
                    }
                }
            }
        },
        {
            "parameters": {
                "deviceType": {
                    "id": "ups",
                    "name": "UPS"
                },
                "protocolType": {
                    "id": "pmbus",
                    "name": "电总"
                },
                "brandType": {
                    "id": "emerson",
                    "name": "艾默生"
                }
            },
            "content": {
                "id": "ita2",
                "name": "ITA2",
                "model": "ita2",
                "image": "../pic/vis/ita2.svg",
                "desc": "艾默生ITA2-UPS",
                "default": {
                    "samplePort": {
                        "protocol": "protocol-pmbus",
                        "setting": {
                            "baudRate": "9600"
                        }
                    },
                    "sampleUnit": {
                        "cuperiod": "3000",
                        "cutimeout": "1000",
                        "cudelay": "500",
                        "cuthrottle": "500",
                        "cumaxNum": "5",
                        "library": "ita2.json",
                        "setting": {
                            "address": "1"
                        }
                    }
                }
            }
        },
        {
            "parameters": {
                "deviceType": {
                    "id": "converter",
                    "name": "变频器"
                },
                "protocolType": {
                    "id": "modbus",
                    "name": "MODBUS"
                },
                "brandType": {
                    "id": "emerson",
                    "name": "艾默生"
                }
            },
            "content": {
                "id": "mv-fc",
                "name": "MV系列变频器",
                "model": "mv-frequency-converter",
                "image": "../pic/vis/mv-fc.svg",
                "desc": "艾默生MV系列变频器",
                "default": {
                    "samplePort": {
                        "protocol": "protocol-modbus-serial",
                        "setting": {
                            "baudRate": "9600"
                        }
                    },
                    "sampleUnit": {
                        "cuperiod": "3000",
                        "cutimeout": "1000",
                        "cudelay": "500",
                        "cuthrottle": "500",
                        "cumaxNum": "5",
                        "library": "mv-frequency-converter.json",
                        "setting": {
                            "address": "1"
                        }
                    }
                }
            }
        },
        {
            "parameters": {
                "deviceType": {
                    "id": "converter",
                    "name": "变频器"
                },
                "protocolType": {
                    "id": "modbus",
                    "name": "MODBUS"
                },
                "brandType": {
                    "id": "hy-electric",
                    "name": "华远电气"
                }
            },
            "content": {
                "id": "s1-fc",
                "name": "S1系列变频器",
                "model": "s1-frequency-converter",
                "image": "../pic/vis/s1-fc.svg",
                "desc": "华远电气S1系列变频器",
                "default": {
                    "samplePort": {
                        "protocol": "protocol-modbus-serial",
                        "setting": {
                            "baudRate": "9600"
                        }
                    },
                    "sampleUnit": {
                        "cuperiod": "3000",
                        "cutimeout": "1000",
                        "cudelay": "500",
                        "cuthrottle": "500",
                        "cumaxNum": "5",
                        "library": "s1-frequency-converter.json",
                        "setting": {
                            "address": "1"
                        }
                    }
                }
            }
        },
        {
            "parameters": {
                "deviceType": {
                    "id": "power-meter",
                    "name": "电力仪表"
                },
                "protocolType": {
                    "id": "modbus",
                    "name": "MODBUS"
                },
                "brandType": {
                    "id": "acrel",
                    "name": "安科瑞"
                }
            },
            "content": {
                "id": "aew100",
                "name": "AEW100电力仪表",
                "model": "AEW100",
                "image": "../pic/vis/AEW100.svg",
                "desc": "安科瑞AEW100电力仪表",
                "default": {
                    "samplePort": {
                        "protocol": "protocol-modbus-serial",
                        "setting": {
                            "baudRate": "9600"
                        }
                    },
                    "sampleUnit": {
                        "cuperiod": "3000",
                        "cutimeout": "1000",
                        "cudelay": "500",
                        "cuthrottle": "500",
                        "cumaxNum": "5",
                        "library": "AEW100.json",
                        "setting": {
                            "address": "1"
                        }
                    }
                }
            }
        },
        {
            "parameters": {
                "deviceType": {
                    "id": "power-meter",
                    "name": "电力仪表"
                },
                "protocolType": {
                    "id": "modbus",
                    "name": "MODBUS"
                },
                "brandType": {
                    "id": "acrel",
                    "name": "安科瑞"
                }
            },
            "content": {
                "id": "amc72l-e4",
                "name": "AMC72L-E4电力仪表",
                "model": "AMC72L-E4",
                "image": "../pic/vis/AMC72L-E4.svg",
                "desc": "安科瑞AMC72L-E4电力仪表",
                "default": {
                    "samplePort": {
                        "protocol": "protocol-modbus-serial",
                        "setting": {
                            "baudRate": "9600"
                        }
                    },
                    "sampleUnit": {
                        "cuperiod": "3000",
                        "cutimeout": "1000",
                        "cudelay": "500",
                        "cuthrottle": "500",
                        "cumaxNum": "5",
                        "library": "AMC72L-E4.json",
                        "setting": {
                            "address": "1"
                        }
                    }
                }
            }
        },
        {
            "parameters": {
                "deviceType": {
                    "id": "power-meter",
                    "name": "电力仪表"
                },
                "protocolType": {
                    "id": "modbus",
                    "name": "MODBUS"
                },
                "brandType": {
                    "id": "elebest",
                    "name": "贝斯特"
                }
            },
            "content": {
                "id": "se-pm",
                "name": "SE系列电力仪表",
                "model": "se-series",
                "image": "../pic/vis/se-pm.svg",
                "desc": "贝斯特SE系列电力仪表",
                "default": {
                    "samplePort": {
                        "protocol": "protocol-modbus-serial",
                        "setting": {
                            "baudRate": "9600"
                        }
                    },
                    "sampleUnit": {
                        "cuperiod": "3000",
                        "cutimeout": "1000",
                        "cudelay": "500",
                        "cuthrottle": "500",
                        "cumaxNum": "5",
                        "library": "se-series-multi-function-power-meter.json",
                        "setting": {
                            "address": "1"
                        }
                    }
                }
            }
        },
        {
            "parameters": {
                "deviceType": {
                    "id": "sensorflow",
                    "name": "sensorflow"
                },
                "protocolType": {
                    "id": "modbus",
                    "name": "MODBUS"
                },
                "brandType": {
                    "id": "hyiot",
                    "name": "华远云联"
                }
            },
            "content": {
                "id": "rfid",
                "name": "sensorflow RFID标签",
                "model": "sensorflow-tag-rfid",
                "image": "../pic/vis/sensorflow-tag-rfid.svg",
                "desc": "华远云联sensorflow RFID标签",
                "default": {
                    "samplePort": {
                        "protocol": "protocol-modbus-serial",
                        "setting": {
                            "baudRate": "38400",
                            "keyNumber": "6",
                            "WANInterface": "eth0.2",
                            "WifiInterface": "br-lan"
                        }
                    },
                    "sampleUnit": {
                        "cuperiod": "1000",
                        "cutimeout": "1000",
                        "cudelay": "500",
                        "cuthrottle": "500",
                        "cumaxNum": "5",
                        "library": "sensorflow-tag-rfid.json",
                        "setting": {
                            "address": "1"
                        }
                    }
                }
            }
        },
        {
            "parameters": {
                "deviceType": {
                    "id": "sensorflow",
                    "name": "sensorflow"
                },
                "protocolType": {
                    "id": "modbus",
                    "name": "MODBUS"
                },
                "brandType": {
                    "id": "hyiot",
                    "name": "华远云联"
                }
            },
            "content": {
                "id": "acousto-optic",
                "name": "sensorflow声光标签",
                "model": "sensorflow-tag-acousto-optic",
                "image": "../pic/vis/sensorflow-tag-acousto-optic.svg",
                "desc": "华远云联sensorflow声光标签",
                "default": {
                    "samplePort": {
                        "protocol": "protocol-modbus-serial",
                        "setting": {
                            "baudRate": "38400",
                            "keyNumber": "6",
                            "WANInterface": "eth0.2",
                            "WifiInterface": "br-lan"
                        }
                    },
                    "sampleUnit": {
                        "cuperiod": "1000",
                        "cutimeout": "1000",
                        "cudelay": "500",
                        "cuthrottle": "500",
                        "cumaxNum": "5",
                        "library": "sensorflow-tag-acousto-optic.json",
                        "setting": {
                            "address": "1"
                        }
                    }
                }
            }
        },
        {
            "parameters": {
                "deviceType": {
                    "id": "sensorflow",
                    "name": "sensorflow"
                },
                "protocolType": {
                    "id": "modbus",
                    "name": "MODBUS"
                },
                "brandType": {
                    "id": "hyiot",
                    "name": "华远云联"
                }
            },
            "content": {
                "id": "door",
                "name": "sensorflow门磁标签",
                "model": "sensorflow-tag-door",
                "image": "../pic/vis/sensorflow-tag-door.svg",
                "desc": "华远云联sensorflow门磁标签",
                "default": {
                    "samplePort": {
                        "protocol": "protocol-modbus-serial",
                        "setting": {
                            "baudRate": "38400",
                            "keyNumber": "6",
                            "WANInterface": "eth0.2",
                            "WifiInterface": "br-lan"
                        }
                    },
                    "sampleUnit": {
                        "cuperiod": "1000",
                        "cutimeout": "1000",
                        "cudelay": "500",
                        "cuthrottle": "500",
                        "cumaxNum": "5",
                        "library": "sensorflow-tag-door.json",
                        "setting": {
                            "address": "1"
                        }
                    }
                }
            }
        },
        {
            "parameters": {
                "deviceType": {
                    "id": "sensorflow",
                    "name": "sensorflow"
                },
                "protocolType": {
                    "id": "modbus",
                    "name": "MODBUS"
                },
                "brandType": {
                    "id": "hyiot",
                    "name": "华远云联"
                }
            },
            "content": {
                "id": "th",
                "name": "ensorflow温湿度标签",
                "model": "sensorflow-tag-th",
                "image": "../pic/vis/sensorflow-tag-th.svg",
                "desc": "华远云联sensorflow温湿度标签",
                "default": {
                    "samplePort": {
                        "protocol": "protocol-modbus-serial",
                        "setting": {
                            "baudRate": "38400",
                            "keyNumber": "6",
                            "WANInterface": "eth0.2",
                            "WifiInterface": "br-lan"
                        }
                    },
                    "sampleUnit": {
                        "cuperiod": "1000",
                        "cutimeout": "1000",
                        "cudelay": "500",
                        "cuthrottle": "500",
                        "cumaxNum": "5",
                        "library": "sensorflow-tag-th.json",
                        "setting": {
                            "address": "1"
                        }
                    }
                }
            }
        },
        {
            "parameters": {
                "deviceType": {
                    "id": "sensorflow",
                    "name": "sensorflow"
                },
                "protocolType": {
                    "id": "modbus",
                    "name": "MODBUS"
                },
                "brandType": {
                    "id": "hyiot",
                    "name": "华远云联"
                }
            },
            "content": {
                "id": "pt100",
                "name": "sensorflow PT100标签",
                "model": "sensorflow-tag-pt100",
                "image": "../pic/vis/sensorflow-tag-pt100.svg",
                "desc": "华远云联sensorflow PT100标签",
                "default": {
                    "samplePort": {
                        "protocol": "protocol-modbus-serial",
                        "setting": {
                            "baudRate": "38400",
                            "keyNumber": "6",
                            "WANInterface": "eth0.2",
                            "WifiInterface": "br-lan"
                        }
                    },
                    "sampleUnit": {
                        "cuperiod": "1000",
                        "cutimeout": "1000",
                        "cudelay": "500",
                        "cuthrottle": "500",
                        "cumaxNum": "5",
                        "library": "sensorflow-tag-pt100.json",
                        "setting": {
                            "address": "1"
                        }
                    }
                }
            }
        },
        {
            "parameters": {
                "deviceType": {
                    "id": "sensorflow",
                    "name": "sensorflow"
                },
                "protocolType": {
                    "id": "modbus",
                    "name": "MODBUS"
                },
                "brandType": {
                    "id": "hyiot",
                    "name": "华远云联"
                }
            },
            "content": {
                "id": "current",
                "name": "sensorflow电流标签",
                "model": "sensorflow-tag-current",
                "image": "../pic/vis/sensorflow-tag-current.svg",
                "desc": "华远云联sensorflow电流标签",
                "default": {
                    "samplePort": {
                        "protocol": "protocol-modbus-serial",
                        "setting": {
                            "baudRate": "38400",
                            "keyNumber": "6",
                            "WANInterface": "eth0.2",
                            "WifiInterface": "br-lan"
                        }
                    },
                    "sampleUnit": {
                        "cuperiod": "1000",
                        "cutimeout": "1000",
                        "cudelay": "500",
                        "cuthrottle": "500",
                        "cumaxNum": "5",
                        "library": "sensorflow-tag-current.json",
                        "setting": {
                            "address": "1"
                        }
                    }
                }
            }
        },
        {
            "parameters": {
                "deviceType": {
                    "id": "sensorflow",
                    "name": "sensorflow"
                },
                "protocolType": {
                    "id": "modbus",
                    "name": "MODBUS"
                },
                "brandType": {
                    "id": "hyiot",
                    "name": "华远云联"
                }
            },
            "content": {
                "id": "voltage",
                "name": "sensorflow电压标签",
                "model": "sensorflow-tag-voltage",
                "image": "../pic/vis/sensorflow-tag-voltage.svg",
                "desc": "华远云联sensorflow电压标签",
                "default": {
                    "samplePort": {
                        "protocol": "protocol-modbus-serial",
                        "setting": {
                            "baudRate": "38400",
                            "keyNumber": "6",
                            "WANInterface": "eth0.2",
                            "WifiInterface": "br-lan"
                        }
                    },
                    "sampleUnit": {
                        "cuperiod": "1000",
                        "cutimeout": "1000",
                        "cudelay": "500",
                        "cuthrottle": "500",
                        "cumaxNum": "5",
                        "library": "sensorflow-tag-voltage.json",
                        "setting": {
                            "address": "1"
                        }
                    }
                }
            }
        },
        {
            "parameters": {
                "deviceType": {
                    "id": "sensorflow",
                    "name": "sensorflow"
                },
                "protocolType": {
                    "id": "modbus",
                    "name": "MODBUS"
                },
                "brandType": {
                    "id": "hyiot",
                    "name": "华远云联"
                }
            },
            "content": {
                "id": "lampwith",
                "name": "sensorflow灯带标签",
                "model": "sensorflow-tag-lampwith",
                "image": "../pic/vis/sensorflow-tag-lampwith.svg",
                "desc": "华远云联sensorflow灯带标签",
                "default": {
                    "samplePort": {
                        "protocol": "protocol-modbus-serial",
                        "setting": {
                            "baudRate": "38400",
                            "keyNumber": "6",
                            "WANInterface": "eth0.2",
                            "WifiInterface": "br-lan"
                        }
                    },
                    "sampleUnit": {
                        "cuperiod": "1000",
                        "cutimeout": "1000",
                        "cudelay": "500",
                        "cuthrottle": "500",
                        "cumaxNum": "5",
                        "library": "sensorflow-tag-lampwith.json",
                        "setting": {
                            "address": "1"
                        }
                    }
                }
            }
        },
        {
            "parameters": {
                "deviceType": {
                    "id": "sensorflow",
                    "name": "sensorflow"
                },
                "protocolType": {
                    "id": "modbus",
                    "name": "MODBUS"
                },
                "brandType": {
                    "id": "hyiot",
                    "name": "华远云联"
                }
            },
            "content": {
                "id": "u",
                "name": "sensorflow U位",
                "model": "sensorflow-u-tracker",
                "image": "../pic/vis/sensorflow-u-tracker.svg",
                "desc": "华远云联sensorflow U位",
                "default": {
                    "samplePort": {
                        "protocol": "protocol-modbus-serial",
                        "setting": {
                            "baudRate": "38400",
                            "keyNumber": "6",
                            "WANInterface": "eth0.2",
                            "WifiInterface": "br-lan"
                        }
                    },
                    "sampleUnit": {
                        "cuperiod": "1000",
                        "cutimeout": "1000",
                        "cudelay": "500",
                        "cuthrottle": "500",
                        "cumaxNum": "5",
                        "library": "sensorflow-u-tracker.json",
                        "setting": {
                            "address": "1"
                        }
                    }
                }
            }
        },
        {
            "parameters": {
                "deviceType": {
                    "id": "other",
                    "name": "其它"
                },
                "protocolType": {
                    "id": "modbus",
                    "name": "MODBUS"
                },
                "brandType": {
                    "id": "other",
                    "name": "其它"
                }
            },
            "content": {
                "id": "electriccloset",
                "name": "配电柜",
                "model": "electriccloset",
                "image": "../pic/vis/electriccloset.svg",
                "desc": "配电柜",
                "default": {
                    "samplePort": {
                        "protocol": "protocol-modbus-serial",
                        "setting": {
                            "baudRate": "9600"
                        }
                    },
                    "sampleUnit": {
                        "cuperiod": "3000",
                        "cutimeout": "1000",
                        "cudelay": "500",
                        "cuthrottle": "500",
                        "cumaxNum": "5",
                        "library": "electriccloset.json",
                        "setting": {
                            "address": "1"
                        }
                    }
                }
            }
        },
        {
            "parameters": {
                "deviceType": {
                    "id": "other",
                    "name": "其它"
                },
                "protocolType": {
                    "id": "modbus",
                    "name": "MODBUS"
                },
                "brandType": {
                    "id": "emerson",
                    "name": "艾默生"
                }
            },
            "content": {
                "id": "irm-4di",
                "name": "IRM-S04DIF",
                "model": "IRM-S04DIF",
                "image": "../pic/vis/IRM-S04DIF.svg",
                "desc": "艾默生IRM-S04DIF",
                "default": {
                    "samplePort": {
                        "protocol": "protocol-modbus-serial",
                        "setting": {
                            "baudRate": "9600"
                        }
                    },
                    "sampleUnit": {
                        "cuperiod": "3000",
                        "cutimeout": "1000",
                        "cudelay": "500",
                        "cuthrottle": "500",
                        "cumaxNum": "5",
                        "library": "IRM-S04DIF.json",
                        "setting": {
                            "address": "1"
                        }
                    }
                }
            }
        },
        {
            "parameters": {
                "deviceType": {
                    "id": "other",
                    "name": "其它"
                },
                "protocolType": {
                    "id": "modbus",
                    "name": "MODBUS"
                },
                "brandType": {
                    "id": "zhenyang",
                    "name": "真扬"
                }
            },
            "content": {
                "id": "mhtpa-i",
                "name": "MHTPA-I",
                "model": "MHTPA-I",
                "image": "../pic/vis/MHTPA-I.svg",
                "desc": "真扬红外测温产品MHTPA-I",
                "default": {
                    "samplePort": {
                        "protocol": "protocol-modbus-serial",
                        "setting": {
                            "baudRate": "38400"
                        }
                    },
                    "sampleUnit": {
                        "cuperiod": "3000",
                        "cutimeout": "1000",
                        "cudelay": "500",
                        "cuthrottle": "500",
                        "cumaxNum": "5",
                        "library": "MHTPA-I.json",
                        "setting": {
                            "address": "1"
                        }
                    }
                }
            }
        }
    ],
    "netEquipment": [{
            "parameters": {
                "deviceType": {
                    "id": "ups",
                    "name": "UPS"
                },
                "protocolType": {
                    "id": "snmp",
                    "name": "SNMP"
                },
                "brandType": {
                    "id": "emerson",
                    "name": "艾默生"
                }
            },
            "content": {
                "id": "ups",
                "name": "ups",
                "model": "emerson-ups",
                "image": "../pic/vis/emerson-ups.svg",
                "desc": "艾默生UPS",
                "default": {
                    "samplePort": {
                        "protocol": "protocol-snmp-manager",
                        "setting": {}
                    },
                    "sampleUnit": {
                        "cuperiod": "3000",
                        "cutimeout": "1000",
                        "cudelay": "500",
                        "cuthrottle": "500",
                        "cumaxNum": "5",
                        "library": "emerson-ups.json",
                        "setting": {
                            "target": "192.168.1.163",
                            "port": "161",
                            "version": "v2c",
                            "readCommunity": "public",
                            "writeCommunity": "private"
                        }
                    }
                }
            }
        },
        {
            "parameters": {
                "deviceType": {
                    "id": "other",
                    "name": "其它"
                },
                "protocolType": {
                    "id": "modbus",
                    "name": "MODBUS"
                },
                "brandType": {
                    "id": "other",
                    "name": "其它"
                }
            },
            "content": {
                "id": "mbn",
                "name": "未知",
                "model": "未知",
                "image": "../pic/vis/th.svg",
                "desc": "未知",
                "default": {
                    "samplePort": {
                        "protocol": "protocol-modbus-serial",
                        "setting": {
                            "netport": "502"
                        }
                    },
                    "sampleUnit": {
                        "cuperiod": "3000",
                        "cutimeout": "1000",
                        "cudelay": "500",
                        "cuthrottle": "500",
                        "cumaxNum": "5",
                        "library": "th.json",
                        "setting": {
                            "address": "1"
                        }
                    }
                }
            }
        },
        {
            "parameters": {
                "deviceType": {
                    "id": "camera",
                    "name": "摄像机"
                },
                "protocolType": {
                    "id": "video",
                    "name": "Video"
                },
                "brandType": {
                    "id": "hiki",
                    "name": "海康"
                }
            },
            "content": {
                "id": "video",
                "name": "海康摄像机",
                "model": "AA-B1-1001",
                "image": "../pic/vis/hiki-video.svg",
                "desc": "海康摄像机",
                "default": {
                    "samplePort": {
                        "protocol": "protocol-camera",
                        "setting": {}
                    },
                    "sampleUnit": {
                        "cuperiod": "3000",
                        "cutimeout": "1000",
                        "cudelay": "500",
                        "cuthrottle": "500",
                        "cumaxNum": "5",
                        "library": "video.json",
                        "setting": {}
                    }
                }
            }
        },
        {
            "parameters": {
                "deviceType": {
                    "id": "camera",
                    "name": "摄像机"
                },
                "protocolType": {
                    "id": "video",
                    "name": "Video"
                },
                "brandType": {
                    "id": "zhenyang",
                    "name": "真扬"
                }
            },
            "content": {
                "id": "face-ipc",
                "name": "face-ipc",
                "model": "face-ipc",
                "image": "../pic/vis/face-ipc.svg",
                "desc": "真扬摄像机",
                "default": {
                    "samplePort": {
                        "protocol": "protocol-face-ipc",
                        "setting": {
                            "host": "192.168.1.1",
                            "port": "20020",
                            "uploadServer": "lab.huayuan-iot.com",
                            "author": "hyiot",
                            "project": "video",
                            "token": "b2b8ec80-8a3a-11e8-9083-afae74b81b2b",
                            "user": "hyiot"
                        }
                    },
                    "sampleUnit": {
                        "cuperiod": "2000",
                        "cutimeout": "1000",
                        "cudelay": "500",
                        "cuthrottle": "500",
                        "cumaxNum": "5",
                        "library": "face-ipc.json",
                        "setting": {
                            "address": "1"
                        }
                    }
                }
            }
        }
    ]
}