$(document).ready(function () {
    $.ajax({
        url: '/cgi-bin/data.cgi',
        type: 'GET',
        dataType: 'json',
        success: function (res) {
            console.log(res)
            // 第一步
            if (res.step1.type == "wifi") { //联网配置，wifi接入填充
                $("#wifion").attr("selected", "true"); //WiFi接入选项
                $("#ssid").val(res.step1.type.wifi.ssid); //ssid填充
                $("#key").val(res.step1.type.wifi.key); //wifi密码填充
                if (res.step1.type.wifi.encryption == "psk") {
                    $("#optionpsk").attr("selected", "true") //PSK加密
                } else if (res.step1.type.wifi.encryption == "psk2") {
                    $("#optionpsk2").attr("selected", "true") //PSK2加密
                };
            } else if (res.step1.type == "online") { //有线接入
                $("#lineon").attr("selected", "true"); //有线接入选项
                if (res.step1.type.online.proto == "static") {
                    $("#optionstatic").attr("selected", "true"); //静态选项
                    $("#ipaddr").val(res.step1.type.online.proto.ipaddr); //IP地址填充
                    $("#netmask").val(res.step1.type.online.proto.netmask); //子网掩码填充
                    $("#gateway").val(res.step1.type.online.proto.gateway); //网关地址填充
                } else if (res.step1.type.online.proto == "dhcp") {
                    $("#optiondhcp").attr("selected", "true"); //dhcp选项
                }
            } else if (res.step1.code == "404") {
                return false;
                //如果没有获取到值，就不填充
            }
            //第二步
            if (res.step2.code == "404") {
                return fasle;
            } else {
                $("#server").val(res.step2.server);
                $("#port").val(res.step2.port);
            }
            //第三步
            if (res.step3.code == "404") {
                return false;
            } else {
                $("#apssid").val(res.step3.apssid);
                if (res.step3.encryption == "psk") {
                    $("#optionappsk").attr("selected", "true");
                } else if (res.step3.encryption == "psk2") {
                    $("#optionappsk2").attr("selected", "true");
                }
                $("#apkey").val(res.step3.apkey);
            }

            //第四步
            if (res.step4.code == "404") {
                return false;
            } else {
                $("#lanipaddr").val(res.step4.lanipaddr);
                $("#lannetmask").val(res.step4.lannetmask)
            }
            //初始化选项
            loading();
        }
    });
});
//初始化结束

function loading() {
    var setValue = $("#linktype").val();
    console.log(setValue);
    if (setValue == 'wifi') {
        $(".wifioption").show();
        $(".fouroption").hide();
        $(".wiredoption").hide();
    } else if (setValue == 'fourg') {
        $(".wifioption").hide();
        $(".fouroption").show();
        $(".wiredoption").hide();
    } else if (setValue == 'wired') {
        $(".wifioption").hide();
        $(".fouroption").hide();
        $(".wiredoption").show();
    }
    var seta = $("#proto").val();
    console.log(seta);
    if (seta == 'static') {
        $(".staticoption").show();
    } else {
        $(".staticoption").hide();
    }
};
$("#proto").on("change", function () {
    loading();
});
//联网配置首选项开始
$("#linktype").on("change", function () {
    loading();
});
//联网配置首选项结束

//联网配置提交
$("#step1").on("click", function () {
    var option = $("#linktype").val();
    if (option == 'wifi') { //wifi接入
        if ($("#ssid").val().length == 0) {
            alert("请填写SSID");
            return false;
        } else if ($("#key").val().length == 0) {
            alert("请填写WIFI密码");
            return false;
        } else if ($("select[name='encryption']").val() == null) {
            alert("请选择加密方式");
            return false;
        } else {
            var data1 = {
                linktype: $("#linktype").val(),
                ssid: $("#ssid").val(),
                key: $("#key").val(),
                encryption: $("select[name='encryprion']").val() //注意这个是选项，不是中文，为了防止传值乱码问题可以根据选项记录一下分别对应内容
            }
        }
    } else if (option == '4G') { //4G
        if ($("select[name='isp']").val() == null) {
            alert("请选择运营商");
            return false;
        } else {
            var data1 = {
                linktype: $("#linktype").val(),
                isp: $("select[name='isp']").val() //注意这个是选项，不是中文，为了防止传值乱码问题可以根据选项记录一下分别对应内容
            }
        }
        // var data1 = $("#four").serialize();
    } else if (option == 'wired') { //有线接入
        if ($("select[name='proto']").val() == null) {
            alert("请选择地址获取方式");
            return false;
        } else if ($("#ipaddr").val().length == 0) {
            alert("请填写IP地址");
            return false;
        } else if ($("#netmask").val().length == 0) {
            alert("请填写子网掩码");
            return false;
        } else if ($("#gateway").val().length == 0) {
            alert("请填写子网掩码");
            return false;
        } else {
            var data1 = {
                linktype: $("#linktype").val(),
                proto: $("select[name='proto']").val(), //注意这个是选项，不是中文，为了防止传值乱码问题可以根据选项记录一下分别对应内容
                ipaddr: $("#ipaddr").val(),
                netmask: $("#netmask").val(),
                gateway: $("#gateway").val()
            }
        }
    }
    // 请求接口去掉注释即可，data1数据已经序列化
    $.ajax({
        url: '/cgi-bin/step1.cgi',
        // data: data1,
        type: 'POST',
        dataType: 'json',
        success: function (res) {
            if (res.result == true) {
                location.href = "step3.html"
                //location.href = "success.html";
            }else{
                alert("提交失败，原因是"+res.data)
            }
        }
    });
    //location.href = "success.html";
});

//mqtt配置
// $("#step2").on("click", function () {
//     if ($("#server").val().length == 0) {
//         alert("请填写服务器地址");
//         return false;
//     } else if ($("#port").val().length == 0) {
//         alert("请填写端口号");
//         return false;
//     } else {
//         var mqtt = {
//             server: $("#server").val(),
//             port: $("#port").val()
//         }
//     }
//     $.ajax({
//         url: '/cgi-bin/mqtt.cgi',
//         data: mqtt,
//         type: 'POST',
//         dataType: 'json',
//         success: function (res) {
//             if (res.result == true) {
//                 location.href = "step3.html"
//             }else{
//                 alert("提交失败，原因是"+res.data)
//             }
//         }
//     });
//     //location.href = "step3.html";
// })

//AP配置
$("#step3").on("click", function () {
    if ($("#apssid").val().length == 0) {
        alert("请填写SSID");
        return false;
    } else if ($("select[name='encryption']").val() == null) {
        alert("请选择安全性");
        return false;
    } else if ($("#apkey").val().length == 0) {
        alert("请填写密码");
        return false;
    } else {
        var ap = {
            ssid: $("#apssid").val(),
            encryption: $("select[name='encryption']").val(),
            key: $("#apkey").val()
        }
    }
    $.ajax({
        url: '/cgi-bin/ap.cgi',
        // data: ap,
        type: 'POST',
        dataType: 'json',
        success: function (res) {
            if (res.result == true) {
                location.href = "step4.html"
            }else{
                alert("提交失败，原因是"+res.data)
            }
        }
    });
    //location.href = "success.html";
})

$("#step4").on("click", function () { //LAN配置
    if ($("#lanipaddr").val().length == 0) {
        alert("请填写IP地址");
        return false;
    } else if ($("#lannetmask").val().length == 0) {
        alert("请填写子网掩码");
        return false;
    } else {
        var lan = {
            ipaddr: $("#lanipaddr").val(),
            netmask: $("#lannetmask").val()
        }
    }
    $.ajax({
        url: '/cgi-bin/lan.cgi',
        // data: lan,
        type: 'POST',
        dataType: 'json',
        success: function (res) {
            if (res.result == true) {
                location.href = "step4.html"
            }else{
                alert("提交失败，原因是"+res.data)
            }
        }
    });
    //location.href = "step4.html";
})

$("#step4").on("click", function () { //LAN配置
    if ($("#collectorid").val().length == 0) {
        alert("请填写采集器ID");
        return false;
    } else if ($("#collectorip").val().length == 0) {
        alert("请填写采集器IP");
        return false;
    }else if ($("#mqttip").val().length == 0) {
        alert("请填写Mqtt IP");
        return false;
    }else if ($("#mqttport").val().length == 0) {
        alert("请填写Mqtt 端口");
        return false;
    } else {
        var basic = {
            collectorid: $("#lanipaddr").val(),
            collectorip: $("#lannetmask").val(),
            mqttip: $("#lanipaddr").val(),
            mqttport: $("#lannetmask").val()
        }
    }
    $.ajax({
        url: '/cgi-bin/basic.cgi',
        // data: basic,
        type: 'POST',
        dataType: 'json',
        success: function (res) {
            if (res.result == true) {
                location.href = "visualization.html"
            }else{
                alert("提交失败，原因是"+res.data)
            }
        }
    });
    //location.href = "step4.html";
})

$("#restart").on("click",function(){
    var v = confirm("是否确认重启")
    if (v == true){
        //重启代码
        location.href="../index.html"
    }
})