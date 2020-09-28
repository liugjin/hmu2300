$(document).ready(function () {
    $.ajax({
        url: host + "/ap/",
        type: "GET",
        dataType: "JSON",
        success: function (param) {
            console.log(param)
            $("#apssid").val(param.data.ssid);
            $("#encryption1").val(param.data.encryption);
            $("#channel").val(param.data.channel);
            $("#hide").val(param.data.hidden);
            Materialize.updateTextFields();
            $('select').material_select();
        }
    })
})
$("#encryption1").on("change", function () {
    var psw = $("#encryption1").val();
    if (psw == "none") {
        $("#safepassword").hide();
    } else {
        $("#safepassword").show();
    }
})
$("#apkey").on("change", function () {
    var tp = $("#encryption1").val();
    var y = $(this).val();
    var te = /^[a-zA-Z0-9]{8,64}$/;
    real = te.test(y);
    if (tp !== "none") {
        if (real == false) {
            alert("请输入大于8位小于64位的密码！");
            $(this).val("");
            return false;
        }
    }
})

$("#step2").on("click", function () {
    var psw = $("#encryption1").val();
    if (psw == "none") {
        if($("#apssid").val() == ""){
            alert("您有未填写的内容！");
            return false;
        }
        var data = JSON.stringify({
            "ssid": $("#apssid").val(),
            "encryption": $("#encryption1").val(),
            "key": "",
            "channel": $("#channel").val(),
            "hide": $("#hide").val()
        })
    } else {
        if ($("#apkey").val() == "" || $("#apkey").val().length < 8 || $("#apkey").val().length > 64) {
            alert("请输入大于8位小于64位的密码！");
            $(this).val("");
            return false;
        } else {
            if($("#apssid").val() == ""){
                alert("您有未填写的内容！");
                return false;
            }
            var data = JSON.stringify({
                "ssid": $("#apssid").val(),
                "encryption": $("#encryption1").val(),
                "key": $("#apkey").val(),
                "channel": $("#channel").val(),
                "hide": $("#hide").val()
            })
        }
    }
    $.ajax({
        url: host + "/ap/",
        data: data,
        type: "PUT",
        dataType: "JSON",
        success: function (res) {
            console.log(res.status);
            if (res.status == "0") {
                console.log("请求成功");
                location.href = "3lan.html";
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
