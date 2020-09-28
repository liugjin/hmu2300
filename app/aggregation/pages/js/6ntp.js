$(document).ready(function () {
    $.ajax({
        url: host + "/ntp/",
        type: "GET",
        dataType: "JSON",
        success: function (param) {
            $("#ntp1").val(param.data.ntp1);
            $("#ntp2").val(param.data.ntp2);
            $("#ntp3").val(param.data.ntp3);
            Materialize.updateTextFields();
        }
    })
})

$("#6ntp").on("click", function () {
    let data = JSON.stringify({
        "ntp1":$("#ntp1").val(),
        "ntp2":$("#ntp2").val(),
        "ntp3":$("#ntp3").val()
    })

    $.ajax({
        url: host + "/ntp/",
        data: data,
        type: "PUT",
        dataType: "JSON",
        success: function (res) {
            console.log(res)
            if (res.status == 0) {
                console.log("请求成功");
                location.href = "4basic.html";
            } else {
                alert("配置失败！")
            }
        },
        error: function (res1) {
            console.log("失败");
            alert(res1.message);
        }
    })
})



// $("#lanipaddr").on("change", function () {
//     var b = $(this).val();
//     var test = /^(\d|[1-9]\d|1\d{2}|2[0-5][0-5])\.(\d|[1-9]\d|1\d{2}|2[0-5][0-5])\.(\d|[1-9]\d|1\d{2}|2[0-5][0-5])\.(\d|[1-9]\d|1\d{2}|2[0-5][0-5])$/
//     real = test.test(b);
//     if (real == false) {
//         alert("请输入正确的IP地址！");
//         $(this).val("");
//         return false;
//     }
// })

// $("#lannetmask").on("change", function () {
//     var b1 = $(this).val();
//     var test1 = /^((128|192)|2(24|4[08]|5[245]))(\.(0|(128|192)|2((24)|(4[08])|(5[245])))){3}$/;
//     real1 = test1.test(b1);
//     if (real1 == false) {
//         alert("请输入正确的子网掩码！");
//         $(this).val("");
//         return false;
//     }
// })