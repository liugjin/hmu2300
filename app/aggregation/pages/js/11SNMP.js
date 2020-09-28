

// 获取Mib列表
function getMibSList() {
    $.ajax({
        url: `${host}/snmp/miblist `,
        type: 'GET',
        dataType: 'json',
        success: function (res) {
            console.log(res)
            if (res.sList) {
                muport = res.sList
                for (let i = 0; i < muport.length; i++) {
                    $('#linktype').append(
                        `<option value="" id="wifion">${muport[i]}</option>`
                    )
                }
            }

        },
        error: function () {
            alert("异常");
        }
    })
}

//获取community默认值
function getCommunity() {
    $.ajax({
        url: `${host}/snmp/`,
        type: 'GET',
        dataType: 'json',
        success: function (res) {
            console.log(res)
            if (res.data) {
                muport = res.data
                $("#read").val(muport.read)
                $("#write").val(muport.write)
                $("#trapd").val(muport.trapAddr)
            }
        },
        error: function () {
            alert("异常");
        }
    })
}

// 修改Community内容
function setCommunity() {
    $.ajax({
        url: `${host}/snmp/`,
        data: JSON.stringify({ read: $("#read").val(), write: $("#write").val(), trapAddr: $("#trapd").val() }),
        type: 'PUT',
        dataType: 'json',
        success: function (res) {
            console.log(res)
            if (res.message == "ok") {
                alert("Community内容修改成功");
            } else {
                alert("Community内容修改失败，请稍后再试");
            }
        },
        error: function () {
            alert("异常");
        }
    })
}
// 下载文件，前端通过file.js 把后端的文本转成txt文件并进行下载
function download() {
    let text = $("#linktype option:selected").text()
    console.log(text)
    if (text != "请选择mib列表") {
        $.ajax({
            url: `${host}/snmp/mib/${text}`,
            dataType: "text",
            success: function (res) {
                let file = new File([res], `${text}`, { type: "text/plain;charset=utf-8" });
                saveAs.saveAs(file, file.name)
                if (res.data) {
                    muport = res.data
                    $("#read").val(muport.read)
                    $("#write").val(muport.write)
                    $("#trapd").val(muport.trapAddr)
                }
            },
            error: function () {
                alert("异常");
            }
        })
    }
}


$(document).ready(function () {
    getMibSList()
    getCommunity()
    console.log("345",`${host}/snmp/miblist `)
    console.log("函数调用");
});