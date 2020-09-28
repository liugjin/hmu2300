$(document).ready(function () {
    $.ajax({
        url: host + "/internet/",
        type: "GET",
        dataType: "JSON",
        success: function (param) {
            console.log(param)
            if(param.data.internetmode == "eth"){
                $("#linktype").val("eth");
                $(".wifioption").hide();
                $(".fouroption").hide();
                $(".wiredoption").show();
                if(param.data.wanmode == "dhcp"){
                    $("#proto").val("dhcp");
                    $(".staticoption").hide();
                }else if(param.data.wanmode == "static"){
                    $("#proto").val("static");
                    $(".staticoption").show();
                    $("#ipaddr").val(param.data.staticip);
                    $("#netmask").val(param.data.staticmask);
                    $("#gateway").val(param.data.staticgateway);
                    $("#pdns").val(param.data.staticpdns);
                    $("#sdns").val(param.data.staticsdns);
                }
            }else if(param.data.internetmode == "wifi"){
                $("#linktype").val("wifi");
                $(".wifioption").show();
                $(".fouroption").hide();
                $(".wiredoption").hide();
                $("#ssid").val(param.data.wifissid);
                $("#key").val(param.data.wifipass);
            }else if(param.data.internetmode == "lte"){
                $("#linktype").val("lte");
                $(".wifioption").hide();
                $(".fouroption").show();
                $(".wiredoption").hide();
            }
            Materialize.updateTextFields();
            $('select').material_select();
        }
    })
})

function loading() {
    var setValue = $("#linktype").val();
    var seta = $("#proto").val();
    if (setValue == 'wifi') {
        $(".wifioption").show();
        $(".fouroption").hide();
        $(".wiredoption").hide();
        // $("#encryption").on("change",function () {
        //     if($("#encryption").val() == "psk"){
        //         $("#addpswinfo").removeClass("hide")
        //         $("#addpswinfo").addClass("show")
        //     }else if($("#encryption").val() == "psk2"){
        //         $("#addpswinfo").removeClass("hide")
        //         $("#addpswinfo").addClass("show")
        //     }else if($("#encryption").val() == "nopsk"){
        //         $("#addpswinfo").removeClass("show")
        //         $("#addpswinfo").addClass("hide")
        //     }
        // })
    } else if (setValue == 'lte') {
        $(".wifioption").hide();
        $(".fouroption").show();
        $(".wiredoption").hide();
    } else if (setValue == 'eth') {
        $(".wifioption").hide();
        $(".fouroption").hide();
        $(".wiredoption").show();
    }
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

$("#1networking").on("click", function () {
    var setValue = $("#linktype").val();
    var proto = $("#proto").val();
    var type="";
    var data={};
    if (setValue == "wifi") {
        if($("#ssid").val() == "" || $("#key").val() == ""){
            alert("您有未填写的内容！");
            return false;
        }
        type = "/wifi/";
        data = JSON.stringify({
            "ssid": $("#ssid").val(),
            "key": $("#key").val()
        })
    } else if (setValue == "eth") {
        if (proto == "dhcp") {
            type = "/eth/dhcp";
            data = {};
        } else if (proto == "static") {
            if($("#ipaddr").val() == "" || $("#netmask").val() == "" || $("#gateway").val() == "" || $("#pdns").val() == "" || $("#sdns").val() == ""){
                alert("您有未填写的内容！");
                return false;
            }
            type = "/eth/static";
            data = JSON.stringify({
                "ip": $("#ipaddr").val(),
                "mask": $("#netmask").val(),
                "gateway": $("#gateway").val(),
                "pdns": $("#pdns").val(),
                "sdns": $("#sdns").val()
            })
        }
    } else if (setValue == "lte") {
        type = "/lte/";
        data = {};
    }
    if (type.length == 0){
        alert("请选择接入方式");
        return
    }
    $.ajax({
        url: host + type,
        data: data,
        type: "PUT",
        dataType: "json",
        success: function (res) {
            console.log(res);
            if (res.status == "0") {
                console.log("请求成功");
                location.href = "2ap.html";
            } else {
                alert("配置失败！")
            }
        },
        error: function (resp1) {
            console.log("失败");
            if (!resp1.responseText){
                alert(resp1.statusText);
            }else{
                alert(resp1.responseText);
            }
        }
    })

})
