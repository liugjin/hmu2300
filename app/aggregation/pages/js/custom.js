$(document).ready(function () {

    $("select").material_select();
    headTable();
    $.ajax({
        url: host + "/el/",
        dataType: "JSON",
        success: function (param) {
            for (var i = 0; i < param.data.length; i++) {
                $("#equipLibrary").append(
                    '<a href="#!" onclick="equipInfo(this)" class="collection-item">' +
                    param.data[i] +
                    "</a>"
                );
            }
        }
    });
});

$("#pmEnter").on("click", function () {
    $("#modalType").attr("data-type", "1");
})
$("#pmCancel").on("click", function () {
    $("#modalType").attr("data-type", "0");
})

function publicModal(message) {
    $('.modal').modal({
        dismissible: false,
        complete: function () {
            var click = $("#modalType").attr("data-type");
            if (click == "1") {
                var f = confirm("删除操作不可逆，请再次确认！");
                if (f == true) {
                    var equipName = $("#muid")
                        .val();
                    $.ajax({
                        url: host + "/el/" + equipName,
                        type: "DELETE",
                        success: function (res) {
                            console.log(res);
                            location.reload();
                        }
                    })
                }
            }
        }
    });
    $("#publicModal").modal("open");
    $("#pmTitle").text(message);
}

function deleteRow(i) {
    var con = confirm("是否确认删除此数据点");
    if (con == true) {
        $(i)
            .parent()
            .remove();
    }
}

function equipInfo(it) {
    $("#equipTable")
        .siblings()
        .remove();
    $(it).addClass("active");
    $(it)
        .siblings()
        .removeClass("active");
    $("#equipTable").attr("data-id", $(it).text());
    var name = $(it).html();
    $.ajax({
        url: host + "/el/" + name,
        dataType: "JSON",
        success: function (res) {
            var a = [];
            var b = [];
            var c = [];
            if (res.data.channels !== null) {
                for (var i = 0; i < res.data.channels.length; i++) {
                    a.push(res.data.channels[i]);
                }
            }
            if (res.data.mappings[0].mapping !== null) {
                for (var x = 0; x < res.data.mappings[0].mapping.length; x++) {
                    b.push(res.data.mappings[0].mapping[x]);
                }
            }

            //分页
            //组合数组
            for (var v = 0; v < a.length; v++) {
                c.push({
                    id: a[v].id,
                    name: a[v].name,
                    datatype: a[v].datatype,
                    value: a[v].value,
                    address: b[v].address,
                    channel: b[v].channel,
                    code: b[v].code,
                    expression: b[v].expression,
                    format: b[v].format,
                    quantity: b[v].quantity
                });
                $("#suset").append(
                    '<tr class="point-element">' +
                    '<td class="center xuh">' +
                    (v + 1) +
                    "</td>" +
                    '<td><input type="text" value="' +
                    a[v].id +
                    '"></td>' +
                    '<td><input type="text" value="' +
                    a[v].name +
                    '"></td>' +
                    '<td><input type="text" value="' +
                    a[v].datatype +
                    '"></td>' +
                    '<td><input type="text" value="' +
                    a[v].value +
                    '"></td>' +
                    '<td><input type="text" value="' +
                    b[v].code +
                    '"></td>' +
                    '<td><input type="text" value="' +
                    b[v].address +
                    '"></td>' +
                    '<td><input type="text" value="' +
                    b[v].quantity +
                    '"></td>' +
                    '<td><input type="text" value="' +
                    b[v].format +
                    '"></td>' +
                    '<td><input type="text" value="' +
                    b[v].cid1 +
                    '"></td>' +
                    '<td><input type="text" value="' +
                    b[v].cid2 +
                    '"></td>' +
                    '<td><input type="text" value="' +
                    b[v].command +
                    '"></td>' +
                    '<td><input type="text" value="' +
                    b[v].offset +
                    '"></td>' +
                    '<td><input type="text" value="' +
                    b[v].length +
                    '"></td>' +
                    '<td><input type="text" value="' +
                    b[v].oid +
                    '"></td>' +
                    '<td><input type="text" value="' +
                    b[v].nodeId +
                    '"></td>' +
                    '<td><input type="text" value="' +
                    b[v].expression +
                    '"></td>' +
                    '<td onclick="deleteRow(this)"><a class="waves-effect waves-teal btn-flat"><i class="material-icons">delete</i></a></td>' +
                    +"</tr>"
                );
            }
            $("#muid").val(res.data.id);
            $("#muname").val(res.data.name);
            $("#mutype").val(res.data.type);
            $("#mutype").attr("disabled", "true");
            $("#yz").val(res.data.mappings[0].setting.cov);
            pmTable();
            headTable();
            $("select").material_select();
        }
    });
}

function headTable() {
    if (
        $("#mutype").val() == "ModbusElement" ||
        $("#mutype").val() == "HMUElement" ||
        $("#mutype").val() == "CameraElement"
    ) {
        $(".data-code").show();
        $(".data-address").show();
        $(".data-quantity").show();
        $(".data-format").show();
        $(".data-cid1").hide();
        $(".data-cid2").hide();
        $(".data-command").hide();
        $(".data-offset").hide();
        $(".data-length").hide();
        $(".data-oid").hide();
        $(".data-nodeId").hide();
    } else if (
        $("#mutype").val() == "SnmpManagerElement"
    ) {
        $(".data-code").hide();
        $(".data-address").hide();
        $(".data-quantity").hide();
        $(".data-format").hide();
        $(".data-cid1").hide();
        $(".data-cid2").hide();
        $(".data-command").hide();
        $(".data-offset").hide();
        $(".data-length").hide();
        $(".data-oid").show();
        $(".data-nodeId").hide();
    } else if (
        $("#mutype").val() == "opcElement"
    ) {
        $(".data-code").hide();
        $(".data-address").hide();
        $(".data-quantity").hide();
        $(".data-format").hide();
        $(".data-cid1").hide();
        $(".data-cid2").hide();
        $(".data-command").hide();
        $(".data-offset").hide();
        $(".data-length").hide();
        $(".data-nodeId").show();
        $(".data-oid").hide();
    }
    else if (
        $("#mutype").val() == "PMBusElement" ||
        $("#mutype").val() == "OilMachineElement"
    ) {
        $(".data-code").hide();
        $(".data-address").hide();
        $(".data-quantity").hide();
        $(".data-format").hide();
        $(".data-cid1").show();
        $(".data-cid2").show();
        $(".data-command").show();
        $(".data-offset").show();
        $(".data-length").show();
        $(".data-oid").hide();
        $(".data-nodeId").hide();
    }
}

function pmTable() {
    $("#suset tr:gt(0)").each(function () {
        var tr = $(this);
        if (
            $("#mutype").val() == "ModbusElement" ||
            $("#mutype").val() == "HMUElement" ||
            $("#mutype").val() == "CameraElement"
        ) {
            tr.find("td")
                .eq(5)
                .show();
            tr.find("td")
                .eq(6)
                .show();
            tr.find("td")
                .eq(7)
                .show();
            tr.find("td")
                .eq(8)
                .show();

            tr.find("td")
                .eq(9)
                .hide();
            tr.find("td")
                .eq(10)
                .hide();
            tr.find("td")
                .eq(11)
                .hide();
            tr.find("td")
                .eq(12)
                .hide();
            tr.find("td")
                .eq(13)
                .hide();
            tr.find("td")
                .eq(14)
                .hide();
        } else if (
            $("#mutype").val() == "SnmpManagerElement" ||
            $("#mutype").val() == "opcElement"
        ) {
            tr.find("td")
                .eq(5)
                .hide();
            tr.find("td")
                .eq(6)
                .hide();
            tr.find("td")
                .eq(7)
                .hide();
            tr.find("td")
                .eq(8)
                .hide();

            tr.find("td")
                .eq(9)
                .hide();
            tr.find("td")
                .eq(10)
                .hide();
            tr.find("td")
                .eq(11)
                .hide();
            tr.find("td")
                .eq(12)
                .hide();
            tr.find("td")
                .eq(13)
                .hide();
            tr.find("td")
                .eq(14)
                .show();
        } else if (
            $("#mutype").val() == "PMBusElement" ||
            $("#mutype").val() == "OilMachineElement"
        ) {
            tr.find("td")
                .eq(5)
                .hide();
            tr.find("td")
                .eq(6)
                .hide();
            tr.find("td")
                .eq(7)
                .hide();
            tr.find("td")
                .eq(8)
                .hide();

            tr.find("td")
                .eq(9)
                .show();
            tr.find("td")
                .eq(10)
                .show();
            tr.find("td")
                .eq(11)
                .show();
            tr.find("td")
                .eq(12)
                .show();
            tr.find("td")
                .eq(13)
                .show();
            tr.find("td")
                .eq(14)
                .hide();
        }
    });
}
$("#mutype").on("change", function () {
    headTable();
    pmTable();
});
//新增设备
$("#add-equip").on("click", function () {
    $(".point-element").remove();
    $("#muid").val("");
    $("#muname").val("");
    $("#mutype").val("ModbusElement");
    $("#yz").val("");
    headTable();
    pmTable();
    $("#mutype").removeAttr("disabled");
    $("select").material_select();
    $("#equipLibrary")
        .children()
        .removeClass("active");
});
//增加数据点
$("#addata-point").on("click", function () {
    var num = $(".point-element").length + 1;
    $("#suset").append(
        '<tr class="point-element">' +
        '<td class="center xuh">' +
        num +
        "</td>" +
        '<td><input type="text" value=""></td>' +
        '<td><input type="text" value=""></td>' +
        '<td><input type="text" value=""></td>' +
        '<td><input type="text" value=""></td>' +
        '<td><input type="text" value=""></td>' +
        '<td><input type="text" value=""></td>' +
        '<td><input type="text" value=""></td>' +
        '<td><input type="text" value=""></td>' +
        '<td><input type="text" value=""></td>' +
        '<td><input type="text" value=""></td>' +
        '<td><input type="text" value=""></td>' +
        '<td><input type="text" value=""></td>' +
        '<td><input type="text" value=""></td>' +
        '<td><input type="text" value=""></td>' +
        '<td><input type="text" value=""></td>' +
        '<td onclick="deleteRow(this)"><a class="waves-effect waves-teal btn-flat"><i class="material-icons">delete</i></a></td>' +
        +"</tr>"
    );
    // headTable();
    pmTable();
});

//复制数据点
$("#copylast-point").on("click", function () {
    var last = $(".point-element").last();
    $("#suset").append(last.clone());
    for (var i = 0; i < $(".xuh").length; i++) {
        $(".xuh")
            .eq(i)
            .html(i + 1);
    }
});
//删除
$("#delete").on("click", function () {
    publicModal('请确认是否要删除自定义设备库');
})
//保存
$("#save").on("click", function () {
    var upInfo = new Array();
    if (
        $("#mutype").val() == "ModbusElement" ||
        $("#mutype").val() == "HMUElement" ||
        $("#mutype").val() == "CameraElement"
    ) {
        $("#suset tr:gt(0)").each(function () {
            var tr = $(this);
            var info = {
                chid: tr
                    .find("td")
                    .eq(1)
                    .find("input")
                    .val(),
                name: tr
                    .find("td")
                    .eq(2)
                    .find("input")
                    .val(),
                type: tr
                    .find("td")
                    .eq(3)
                    .find("input")
                    .val(),
                value: parseInt(
                    tr
                        .find("td")
                        .eq(4)
                        .find("input")
                        .val()
                ),
                code: parseInt(
                    tr
                        .find("td")
                        .eq(5)
                        .find("input")
                        .val()
                ),
                address: parseInt(
                    tr
                        .find("td")
                        .eq(6)
                        .find("input")
                        .val()
                ),
                quantity: parseInt(
                    tr
                        .find("td")
                        .eq(7)
                        .find("input")
                        .val()
                ),
                format: tr
                    .find("td")
                    .eq(8)
                    .find("input")
                    .val(),
                expression: tr
                    .find("td")
                    .eq(15)
                    .find("input")
                    .val()
            };
            upInfo.push(info);
        });
    } else if (
        $("#mutype").val() == "SnmpManagerElement" ||
        $("#mutype").val() == "opcElement"
    ) {
        $("#suset tr:gt(0)").each(function () {
            var tr = $(this);
            var info = {
                chid: tr
                    .find("td")
                    .eq(1)
                    .find("input")
                    .val(),
                name: tr
                    .find("td")
                    .eq(2)
                    .find("input")
                    .val(),
                type: tr
                    .find("td")
                    .eq(3)
                    .find("input")
                    .val(),
                value: parseInt(
                    tr
                        .find("td")
                        .eq(4)
                        .find("input")
                        .val()
                ),
                oid: tr
                    .find("td")
                    .eq(14)
                    .find("input")
                    .val(),
                expression: tr
                    .find("td")
                    .eq(15)
                    .find("input")
                    .val()
            };
            upInfo.push(info);
        });
    } else if (
        $("#mutype").val() == "PMBusElement" ||
        $("#mutype").val() == "OilMachineElement"
    ) {
        $("#suset tr:gt(0)").each(function () {
            var tr = $(this);
            var info = {
                chid: tr
                    .find("td")
                    .eq(1)
                    .find("input")
                    .val(),
                name: tr
                    .find("td")
                    .eq(2)
                    .find("input")
                    .val(),
                type: tr
                    .find("td")
                    .eq(3)
                    .find("input")
                    .val(),
                value: parseInt(
                    tr
                        .find("td")
                        .eq(4)
                        .find("input")
                        .val()
                ),
                cid1: parseInt(
                    tr
                        .find("td")
                        .eq(9)
                        .find("input")
                        .val()
                ),
                cid2: parseInt(
                    tr
                        .find("td")
                        .eq(10)
                        .find("input")
                        .val()
                ),
                command: parseInt(
                    tr
                        .find("td")
                        .eq(11)
                        .find("input")
                        .val()
                ),
                offset: parseInt(
                    tr
                        .find("td")
                        .eq(12)
                        .find("input")
                        .val()
                ),
                length: parseInt(
                    tr
                        .find("td")
                        .eq(13)
                        .find("input")
                        .val()
                ),
                expression: tr
                    .find("td")
                    .eq(15)
                    .find("input")
                    .val()
            };
            upInfo.push(info);
        });
    }

    data1 = JSON.stringify({
        id: $("#muid").val(),
        name: $("#muname").val(),
        type: $("#mutype").val(),
        version: "1.0.0",
        description: "test",
        cov: parseFloat($("#yz").val()),
        channels: upInfo
    });
    $.ajax({
        url: host + "/el/",
        data: data1,
        dataType: "JSON",
        type: "POST",
        success: function (res) {
            console.log(res);
        }
    });
    setTimeout(function () {
        location.reload();
    }, 2000);
});

//上载json

$("#importJson").on("click", function () {
    var upFile = document.getElementById("upFile");
    upFile.click();
});

$("#upFile").on("change", function () {
    var fullSize = [];
    var files = $(this)[0].files;
    var all = 0;
    $.each(files, function (c) {
        fullSize.push(files[c].size);
    });
    for (var i = 0; i < fullSize.length; i++) {
        all += fullSize[i];
    }
    var upSize = (all / (1024 * 1024)).toFixed(2);
    // 大小是字节为单位
    if (upSize > 1) {
        alert("所选文件大小为" + upSize + "M超过了1M，请重新选取");
        return false;
    } else {
        var fileObj = new FormData();
        $.each(files, function (f) {
            fileObj.append("file", files[f]);
        });
        $.ajax({
            url: host + "/ellib/",
            data: fileObj,
            dataType: "JSON",
            type: "POST",
            processData: false,
            contentType: false,
            success: function (res) {
                console.log(res);
                location.reload();
            }
        })
    }
});
//下载json
$("#exportJson").on("click", function () {
    var jsonName = $("#equipLibrary")
        .find(".active")
        .text();
    $.ajax({
        url: host + "/ellib/" + jsonName,
        dataType: "JSON",
        type: "GET",
        success: function (res) {
            var blob = new Blob([JSON.stringify(res)], {
                type: "application/json"
            });
            var download = URL.createObjectURL(blob);
            var link = document.createElement("a");
            link.download = jsonName;
            link.href = download;
            link.textContent = "asdasdsadsadsad";
            link.id = "downloadJson";
            link.click();
            link.remove();
        }
    });
});