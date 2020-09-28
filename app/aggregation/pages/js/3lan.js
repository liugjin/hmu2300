$(document).ready(function () {
    $.ajax({
        url: host + "/lan/",
        type: "GET",
        dataType: "JSON",
        success: function (param) {
            $("#lanipaddr").val(param.data.lanip);
            $("#lannetmask").val(param.data.lanmask);
            $("#macaddress").val(param.data.lanmac);
            Materialize.updateTextFields();
        }
    })
})

$("#lanipaddr").on("change", function () {
    var b = $(this).val();
    var test = /^(\d|[1-9]\d|1\d{2}|2[0-5][0-5])\.(\d|[1-9]\d|1\d{2}|2[0-5][0-5])\.(\d|[1-9]\d|1\d{2}|2[0-5][0-5])\.(\d|[1-9]\d|1\d{2}|2[0-5][0-5])$/
    real = test.test(b);
    if (real == false) {
        alert("请输入正确的IP地址！");
        $(this).val("");
        return false;
    }
})

$("#lannetmask").on("change", function () {
    var b1 = $(this).val();
    var test1 = /^((128|192)|2(24|4[08]|5[245]))(\.(0|(128|192)|2((24)|(4[08])|(5[245])))){3}$/;
    real1 = test1.test(b1);
    if (real1 == false) {
        alert("请输入正确的子网掩码！");
        $(this).val("");
        return false;
    }
})
function resetForm(id) {
    document.getElementById(id).reset()
}