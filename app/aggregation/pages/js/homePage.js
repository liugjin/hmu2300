

// 获取是否有配置文件
function getcfgok() {
    $.ajax({
        url: `${host}/mu/cfgok`,
        type: 'GET',
        dataType: 'json',
        success: function (res) {
            if (res.data) {
                muport = res.data
                if (muport.buscfgok == false) {
                    $("#message").removeClass('message0')
                    $("#message").addClass('message1')
                }
                if (muport.mucfgok == false) {
                    $("#message").removeClass('message0')
                    $("#message").addClass('message1')
                }
                if (muport.mucfgok == false && muport.mucfgok == false) {
                    $("#message").removeClass('message0')
                    $("#message").addClass('message1')
                }
            }
        },
        error: function () {
            alert("异常");
        }
    })
}




$(document).ready(function () {
    getcfgok()
    console.log("函数调用");
});