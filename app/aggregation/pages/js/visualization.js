localStorage.setItem("supportEquipment",JSON.stringify(json));
$(document).ready(function () {
    $('.tooltipped').tooltip({
        delay: 50
    });
    if ($(window).width() < 600) {
        $(".pc-view").remove();
    } else {
        $(".phone-view").remove();
    }
    $('.modal').modal({
        dismissible: false
    });
    $('select').material_select();
    // lan口接入
    $.fn.addSvgClass = function (className) {
        return this.each(function () {
            var attr = $(this).attr('class')
            if (attr) {
                if (attr.indexOf(className) < 0) {
                    $(this).attr('class', attr + ' ' + className)
                }
            } else {
                $(this).attr('class', className)
            }
        })
    };
    $.fn.removeSvgClass = function (className) {
        return this.each(function () {
            var attr = $(this).attr('class')
            attr = attr.replace(' ' + className, '')
            $(this).attr('class', attr)
        })
    };

    $.ajax({ //如果是摄像机
        url: host + "/video/",
        type: "GET",
        success: function (res) {
            for (n in res.data) {
                $("#camera").append(
                    '<div class="equip-logo">' +
                    '<div class="equip-type"></div>' +
                    '<div class="equip-img"><img src="../pic/vis/camera.svg"></div>' +
                    '<div class="equip-name">' + res.data[n].cameraName + '</div>' +
                    '<div class="equip-line"></div>' +
                    '</div>'
                )
            }
            Materialize.updateTextFields();
        },
        error:function(resp1){
            if (!resp1.responseText){
                alert(resp1.statusText);
            }else{
                alert(resp1.responseText);
            }
        },
    })
    // 获取网关串口数据
    $.ajax({
        url: host + '/mu/',
        type: 'GET',
        dataType: 'json',
        error:function(resp1){
            if (!resp1.responseText){
                alert(resp1.statusText);
            }else{
                alert(resp1.responseText);
            }
        },
        success: function (res) {
            sessionStorage.setItem("muid",res.data[0].id);
            var r1 = [];
            var r2 = [];
            var l1 = [];
            var l2 = [];
            var port = new Array(6);
            var port1 = res.data[0].ports;
            for (j in port1) {
                if (port1[j].id == "wifi") {
                    port[0] = port1[j];
                } else if (port1[j].id == "lan1") {
                    port[1] = port1[j];
                } else if (port1[j].id == "lan2") {
                    port[2] = port1[j];
                } else if (port1[j].id == "rs3") {
                    port[3] = port1[j];
                } else if (port1[j].id == "rs2") {
                    port[4] = port1[j];
                } else if (port1[j].id == "rs1") {
                    port[5] = port1[j];
                }
            }

            for (i in port) {
                var sampleUnits = port[i].sampleUnits
                //端口1-1和1-2的address不能重复,为了不重复for循环，在此做处理后存入sessionStorage
                if (port[i].id == 'rs1' || port[i].id == 'rs2') {
                    for (p in port[i].sampleUnits) {
                        r1.push(port[i].sampleUnits[p].setting.address);
                    }
                } else if (port[i].id == 'rs3') {
                    for (p in port[i].sampleUnits) {
                        r2.push(port[i].sampleUnits[p].setting.address);
                    }
                } else if (port[i].id == 'lan1') {
                    for (p in port[i].sampleUnits) {
                        l1.push(port[i].sampleUnits[p].id);
                    }
                } else if (port[i].id == 'lan2') {
                    for (p in port[i].sampleUnits) {
                        l2.push(port[i].sampleUnits[p].id);
                    }
                }
                if ($(window).width() > 600) {
                    $("#pcView").append(
                        '<div class="port-2" id="' + port[i].id + '">' +
                        '<div style="text-align:center;"><span class="port-name"></span></div>' +
                        '</div>'
                    )
                    if (port[i].id == 'lan1' || port[i].id == 'lan2' || port[i].id == 'wifi') {
                        $("." + port[i].id).attr(
                            'data-tooltip', '<div>IP地址：' + port[i].setting.port + '</div>'
                        );
                    } else {
                        $("." + port[i].id).attr(
                            'data-tooltip',
                            '<div>波特率：' + port[i].setting.baudRate + '</div>' +
                            '<div>端口：' + port[i].setting.port + '</div>'
                        );
                    }
                    $('.' + port[i].id + '-color').removeSvgClass('port-none');
                    $('.' + port[i].id + '-color').addSvgClass('port-had');
                    $('.tooltipped').tooltip({
                        delay: 50
                    });
                    for (o in sampleUnits) {
                        var logo = sampleUnits[o].element.split('.');
                        $("#" + port[i].id).append(
                            '<div class="equip-logo tooltipped" data-position="right" data-delay="80" data-html="true" data-tooltip="<div class=tooltipLeft>设备地址：' + sampleUnits[o].setting.address + '</div><div class=tooltipLeft>采集周期：' + sampleUnits[o].period + '</div><div class=tooltipLeft>最大通讯错误数：' + sampleUnits[o].maxCommunicationErrors + '</div>">' +
                            '<div class="equip-type" id="ws' + sampleUnits[o].id + '"></div>' +
                            '<div class="equip-delete" onclick="deleteEquip(this)" data-id="' + sampleUnits[o].id + '"><img src="../pic/vis/delete.svg"></div>' +
                            '<div class="equip-img"><img src="../pic/vis/' + logo[0] + '.svg" onerror="imgLocation(this)"></div>' +
                            '<div class="equip-name">' + sampleUnits[o].name + '</div>' +
                            '<div class="equip-line"></div>' +
                            '</div>'
                        )
                        $('.tooltipped').tooltip({
                            delay: 50
                        });
                    }
                } else {
                    $("#tabSon").append(
                        '<div id="phone' + port[i].id + '" class="phone-page"></div>'
                    )
                    $('ul.tabs').tabs();
                    if (port[i].id == "lan1" || port[i].id == "lan2" || port[i].id == "wifi") {
                        var listName = "lanList";
                    } else {
                        var listName = "rsList";
                    }
                    if (port[i].sampleUnits == null) {
                        $("#phone" + port[i].id).removeClass("phone-page");
                        $("#phone" + port[i].id).addClass("phone-ne-page");
                        $("#phone" + port[i].id).append(
                            '<div onclick="' + listName + '(this)" data-set="' + port[i].id + '">添加设备</div>' +
                            '<object style="width: 50%;" data="../pic/noEquip.svg"></object>' +
                            '<div>当前端口没有设备，请添加设备！</div>'
                        )
                    } else {
                        $("#phone" + port[i].id).append('<div onclick="' + listName + '(this)" data-set=' + port[i].id + '>添加设备</div>')
                        for (o in sampleUnits) {
                            var logo = sampleUnits[o].element.split('.');
                            $("#phone" + port[i].id).append(
                                '<div class="equip-logo" onclick="phoneInfoModal(' + "'" + sampleUnits[o].name + "'" + ',' + sampleUnits[o].setting.address + ',' + sampleUnits[o].period + ',' + sampleUnits[o].maxCommunicationErrors + ')">' +
                                '<div class="equip-type" id="ws' + sampleUnits[o].id + '"></div>' +
                                '<div class="equip-img"><img src="../pic/vis/' + logo[0] + '.svg" onerror="imgLocation(this)"></div>' +
                                '<div class="equip-name">' + sampleUnits[o].name + '</div>' +
                                '<div class="equip-line"></div>' +
                                '</div>'
                            )
                        }
                    }

                }
            }
            var wsHost = host.replace("http", "ws");
            var ws = new WebSocket(wsHost + "/ws");
            ws.onmessage = function (evt) {
                var id = JSON.parse(evt.data)
                for (a in id) {
                    if (id[a].state == 1) {
                        $("#ws" + id[a].suid).addClass("green11");
                    } else {
                        $("#ws" + id[a].suid).removeClass("green11");
                    }
                }
            }
            sessionStorage.setItem('lan1Length', JSON.stringify(l1.sort()))
            sessionStorage.setItem('lan2Length', JSON.stringify(l2.sort()))
            sessionStorage.setItem('r1Address', JSON.stringify(r1.sort()))
            sessionStorage.setItem('r2Address', JSON.stringify(r2.sort()))
        }
    });
});

function closeModal(modal) {
    var name = $(modal).attr("data-parent");
    $("#" + name).modal("close");
}

//rs口
function fileType(a) {
    $('#modalRs1').children().remove();
    var typeState = [];
    for (a in collect) {
        var state = $("#type_" + collect[a].type).prop('checked');
        typeState.push(state);
        if (state == true) {
            for (t in collect[a].content)
                $('#modalRs1').append(
                    '<div class="col s6 m6 l2" id="group' + collect[a].content[t].type + '">' +
                    '<div class="card small" style="cursor:pointer;" onclick="addRs(this)" data-equipid="' + collect[a].content[t].id + '" data-name="' + collect[a].content[t].name + '">' +
                    '<div class="card-image">' +
                    '<img src="' + collect[a].content[t].parameter.image + '">' +
                    '<span class="card-title" style="color:#000000;">' + collect[a].content[t].parameter.name + '</span>' +
                    '</div>' +
                    '<div class="card-content">' +
                    '<p>型号：' + collect[a].content[t].parameter.model + '</p>' +
                    '<p>厂商：' + collect[a].content[t].parameter.brand + '</p>' +
                    '</div>' +
                    '</div>' +
                    '</div>'
                )
        }
    }
    if (typeState.indexOf(true) < 0) {
        $('#modalRs1').children().remove();
        for (a in collect) {
            for (t in collect[a].content)
                $('#modalRs1').append(
                    '<div class="col s6 m6 l2" id="group' + collect[a].content[t].type + '">' +
                    '<div class="card small" style="cursor:pointer;" onclick="addRs(this)" data-equipid="' + collect[a].content[t].id + '" data-name="' + collect[a].content[t].name + '">' +
                    '<div class="card-image">' +
                    '<img src="' + collect[a].content[t].parameter.image + '">' +
                    '<span class="card-title" style="color:#000000;">' + collect[a].content[t].parameter.name + '</span>' +
                    '</div>' +
                    '<div class="card-content">' +
                    '<p>型号：' + collect[a].content[t].parameter.model + '</p>' +
                    '<p>厂商：' + collect[a].content[t].parameter.brand + '</p>' +
                    '</div>' +
                    '</div>' +
                    '</div>'
                )
        }
    }
}

function rsList(rs) {
    $('#modalRs1').children().remove();
    $('#selectRs').children().remove();
    $('#modalRs').modal('open');
    var port = $(rs).attr('data-set');
    $('#modalRs').attr('data-rs', port);
    for (a in collect) {
        $('#selectRs').append(
            '<span>' +
            '<input onclick="fileType(this)" type="checkbox" data-id="' + collect[a].type + '" id="type_' + collect[a].type + '" />' +
            '<label for="type_' + collect[a].type + '">' + collect[a].typeName + '</label>' +
            '</span>'
        )
        for (v in collect[a].content)
            $('#modalRs1').append(
                '<div class="col s6 m6 l2">' +
                '<div class="card small" style="cursor:pointer;" onclick="addRs(this)" data-equipid="' + collect[a].content[v].id + '" data-name="' + collect[a].content[v].name + '">' +
                '<div class="card-image">' +
                '<img src="' + collect[a].content[v].parameter.image + '">' +
                '<span class="card-title" style="color:#000000;">' + collect[a].content[v].parameter.name + '</span>' +
                '</div>' +
                '<div class="card-content">' +
                '<p>型号：' + collect[a].content[v].parameter.model + '</p>' +
                '<p>厂商：' + collect[a].content[v].parameter.brand + '</p>' +
                '</div>' +
                '</div>' +
                '</div>'
            )
    }
}

function addRs(su) {
    var choose = $(su).attr('data-equipid');
    var spid = $("#modalRs").attr('data-rs');
    var r1Address = JSON.parse(sessionStorage.getItem('r1Address'));
    var r2Address = JSON.parse(sessionStorage.getItem('r2Address'));
    if (spid == 'rs1' || spid == 'rs2') {
        if (r1Address.length > 0) {
            var setAddress = r1Address[r1Address.length - 1] + 1;
        } else {
            var setAddress = 1;
        }
    } else if (spid == 'rs3') {
        if (r2Address.length > 0) {
            var setAddress = r2Address[r2Address.length - 1] + 1;
        } else {
            var setAddress = 1;
        }
    }
    for (h in collect) {
        for (g in collect[h].content) {
            if (collect[h].content[g].id == choose) {
                var addData = JSON.stringify({
                    "muid": sessionStorage.getItem("muid"),
                    "spid": spid,
                    "id": spid + '-' + choose + '-' + setAddress,
                    "name": $(su).find('.card-title').text(),
                    "period": parseInt(collect[h].content[g].default.cuperiod),
                    "timeout": parseInt(collect[h].content[g].default.cutimeout),
                    "maxCommunicationErrors": parseInt(collect[h].content[g].default.cumaxNum),
                    "element": collect[h].content[g].default.library,
                    "enable": true,
                    "setting": {
                        "address": parseInt(setAddress)
                    }
                })
                $.ajax({
                    url: host + "/su/",
                    data: addData,
                    type: "POST",
                    dataType: "json",
                    error:function(resp1){
                        if (!resp1.responseText){
                            alert(resp1.statusText);
                        }else{
                            alert(resp1.responseText);
                        }
                    },
                    success: function (res) {
                        location.reload();
                    }
                })
            }
        }

    }
}
//lan口****************************************************************************************************************
function lanList(lan) {
    $('#modalLan1').children().remove();
    $('#selectLan').children().remove();
    $('#modalLan').modal('open');
    var port = $(lan).attr('data-set');
    $('#modalLan').attr('data-lan', port);
    for (choose of json.userChoose) {
        $('#selectLan').append(
            `<div id=${choose.id}><span class='choose-name'>${choose.name}:</span></div>`
        )
        for (type of json[choose.id]) {
            $(`#${choose.id}`).append(
                `<span class='choose-type'>` +
                `<input onclick="lanFileType(this)" type="checkbox" data-id="${type.id}" id="${choose.id}_${type.id}" />` +
                `<label for="${choose.id}_${type.id}">${type.name}</label>` +
                `</span>`
            )
        }
    }
    for (equipment of json.netEquipment) {
        console.log(equipment.parameters)
    }
    // json.netEquipment

    // for (a in netEquipment) {
    //     for (v in netEquipment[a].content)
    //         $('#modalLan1').append(
    //             '<div class="col s6 m6 l2">' +
    //             '<div class="card small" style="cursor:pointer;" onclick="addLan(this)" data-equiptype="' + netEquipment[a].type + '" data-equipid="' + netEquipment[a].content[v].id + '" data-name="' + netEquipment[a].content[v].name + '">' +
    //             '<div class="card-image">' +
    //             '<img src="' + netEquipment[a].content[v].parameter.image + '">' +
    //             '<span class="card-title" style="color:#000000;">' + netEquipment[a].content[v].parameter.name + '</span>' +
    //             '</div>' +
    //             '<div class="card-content">' +
    //             '<p>型号：' + netEquipment[a].content[v].parameter.model + '</p>' +
    //             '<p>厂商：' + netEquipment[a].content[v].parameter.brand + '</p>' +
    //             '</div>' +
    //             '</div>' +
    //             '</div>'
    //         )
    // }
}



function lanFileType(a) {
    console.log(a);
    
    $('#modalLan1').children().remove();
    var typeState = [];
    for (a in netEquipment) {
        var state = $("#type_" + netEquipment[a].type).prop('checked');
        typeState.push(state);
        if (state == true) {
            for (t in netEquipment[a].content)
                $('#modalLan1').append(
                    '<div class="col s6 m6 l2" id="group' + netEquipment[a].content[t].type + '">' +
                    '<div class="card small" style="cursor:pointer;" onclick="addLan(this)" data-equiptype="' + netEquipment[a].type + '" data-equipid="' + netEquipment[a].content[t].id + '" data-name="' + netEquipment[a].content[t].name + '">' +
                    '<div class="card-image">' +
                    '<img src="' + netEquipment[a].content[t].parameter.image + '">' +
                    '<span class="card-title" style="color:#000000;">' + netEquipment[a].content[t].parameter.name + '</span>' +
                    '</div>' +
                    '<div class="card-content">' +
                    '<p>型号：' + netEquipment[a].content[t].parameter.model + '</p>' +
                    '<p>厂商：' + netEquipment[a].content[t].parameter.brand + '</p>' +
                    '</div>' +
                    '</div>' +
                    '</div>'
                )
        }
    }
    if (typeState.indexOf(true) < 0) {
        $('#modalLan1').children().remove();
        for (a in netEquipment) {
            for (t in netEquipment[a].content)
                $('#modalLan1').append(
                    '<div class="col s6 m6 l2" id="group' + netEquipment[a].content[t].type + '">' +
                    '<div class="card small" style="cursor:pointer;" onclick="addLan(this)" data-equiptype="' + netEquipment[a].type + '" data-equipid="' + netEquipment[a].content[t].id + '" data-name="' + netEquipment[a].content[t].name + '">' +
                    '<div class="card-image">' +
                    '<img src="' + netEquipment[a].content[t].parameter.image + '">' +
                    '<span class="card-title" style="color:#000000;">' + netEquipment[a].content[t].parameter.name + '</span>' +
                    '</div>' +
                    '<div class="card-content">' +
                    '<p>型号：' + netEquipment[a].content[t].parameter.model + '</p>' +
                    '<p>厂商：' + netEquipment[a].content[t].parameter.brand + '</p>' +
                    '</div>' +
                    '</div>' +
                    '</div>'
                )
        }
    }
}

function addLan(su) {
    var type = $(su).attr("data-equiptype");
    var choose = $(su).attr("data-equipid");
    var spid = $("#modalLan").attr('data-lan');
    var sortNum = []
    if (spid == "lan1") {
        var lan1Length = JSON.parse(sessionStorage.getItem("lan1Length"));
        if (lan1Length == 0) {
            var setAddress = 1;
        } else {
            for (c in lan1Length) {
                var numb = lan1Length[c].split("-");
                sortNum.push(numb[numb.length - 1])
            }
            sortNum.sort();
            var setAddress = parseInt(sortNum[sortNum.length - 1]) + 1;
        }
    } else if (spid == "lan2") {
        var lan2Length = JSON.parse(sessionStorage.getItem("lan2Length"));
        if (lan2Length == 0) {
            var setAddress = 1;
        } else {
            for (c in lan2Length) {
                var numb = lan2Length[c].split("-");
                sortNum.push(numb[numb.length - 1])
            }
            sortNum.sort();
            var setAddress = parseInt(sortNum[sortNum.length - 1]) + 1;
        }
    }
    //如果是摄像机类型,需要同步云视频和hmu的配置，需要两个请求添加
    if (type == "camera") {
        var content = "<form id = 'cameraAdd'><input type='text' name='cameraIp' placeholder = '请填写ip地址'/>" + "<input type='text' name='cameraAccount' placeholder = '请填写账户'/>" + "<input type='password' name='cameraPassword' placeholder = '请填写密码'/></form>"
        selfModal("请填写一下信息", content, function (ok) {
            if (ok == true) {
                var number = [];
                var list = getCameraList();
                if (list.length > 0) {
                    for (num in list) {
                        var e = list[num].cameraId.slice(-1);
                        number.push(e);
                    }
                    var maxNum = Math.max.apply(null, number) + 1;
                } else {
                    var maxNum = 0 + 1;
                }
                var uuid = getUuid();
                var cameraId = uuid + "_camera" + maxNum;
                var userMsg = $("#cameraAdd").serializeArray();
                addCloudCamera(userMsg, cameraId, function (aa) {
                    for (h in netEquipment) {
                        for (g in netEquipment[h].content) {
                            if (netEquipment[h].content[g].id == choose) {
                                var addData = JSON.stringify({
                                    "muid": sessionStorage.getItem("muid"),
                                    "spid": spid,
                                    "id": cameraId,
                                    "name": $(su).find('.card-title').text(),
                                    "period": parseInt(netEquipment[h].content[g].default.cuperiod),
                                    "timeout": parseInt(netEquipment[h].content[g].default.cutimeout),
                                    "maxCommunicationErrors": parseInt(netEquipment[h].content[g].default.cumaxNum),
                                    "element": netEquipment[h].content[g].default.library,
                                    "enable": true,
                                    "setting": {
                                        "ip": parseInt($("#lanIpAddress").val())
                                    }
                                });
                                $.ajax({
                                    url: host + "/su/",
                                    data: addData,
                                    type: "POST",
                                    dataType: "json",
                                    error:function(resp1){
                                        if (!resp1.responseText){
                                            alert(resp1.statusText);
                                        }else{
                                            alert(resp1.responseText);
                                        }
                                    },
                                    success: function (res) {
                                        location.reload();
                                    }
                                })
                            }
                        }
                    }
                });
            }
        })
    } else {
        for (h in netEquipment) {
            for (g in netEquipment[h].content) {
                if (netEquipment[h].content[g].id == choose) {
                    var addData = JSON.stringify({
                        "muid": sessionStorage.getItem("muid"),
                        "spid": spid,
                        "id": spid + '-' + choose + '-' + setAddress,
                        "name": $(su).find('.card-title').text(),
                        "period": parseInt(netEquipment[h].content[g].default.cuperiod),
                        "timeout": parseInt(netEquipment[h].content[g].default.cutimeout),
                        "maxCommunicationErrors": parseInt(netEquipment[h].content[g].default.cumaxNum),
                        "element": netEquipment[h].content[g].default.library,
                        "enable": true,
                        "setting": {
                            "ip": parseInt($("#lanIpAddress").val())
                        }
                    });
                    $.ajax({
                        url: host + "/su/",
                        data: addData,
                        type: "POST",
                        dataType: "json",
                        error:function(resp1){
                            if (!resp1.responseText){
                                alert(resp1.statusText);
                            }else{
                                alert(resp1.responseText);
                            }
                        },
                        success: function (res) {
                            location.reload();
                        }
                    })
                }
            }
        }
    }


}

function deleteEquip(equipment) {
    var equipId = $(equipment).attr("data-id");
    var delDate1 = JSON.stringify({
        "muid": sessionStorage.getItem("muid"),
        "spid": $(equipment).parent().parent().attr("id"),
        "id": equipId
    })
    selfModal("是否删除采集单元", null, function (ok) {
        if (ok == true) {
            if ($(equipment).attr("data-id").indexOf("camera") > 0) {
                deleteCloudCamera(equipId, function () {
                    $.ajax({
                        url: host + "/su/",
                        type: "DELETE",
                        data: delDate1,
                        error:function(resp1){
                            if (!resp1.responseText){
                                alert(resp1.statusText);
                            }else{
                                alert(resp1.responseText);
                            }
                        },
                        success: function (param) {
                            location.reload();
                        }
                    })
                })
            } else {
                $.ajax({
                    url: host + "/su/",
                    type: "DELETE",
                    data: delDate1,
                    error:function(resp1){
                        if (!resp1.responseText){
                            alert(resp1.statusText);
                        }else{
                            alert(resp1.responseText);
                        }
                    },
                    success: function (param) {
                        location.reload();
                    }
                })
            }
        }
    })
}

//监听默认图片
function imgLocation(defaultImg) {
    $(defaultImg).attr("src", "../pic/vis/currency.svg")
}

//手机模拟点击盒子口
function simulationClick(port) {
    $(".phoneActive").removeClass("port-had");
    var clickPort = $(port).attr("data-id");
    $(port).find(".phoneActive").addClass("port-had");
    $("#simPhone" + clickPort).click();
}
//手机信息弹窗
function phoneInfoModal(name, address, period, communication, state) {
    $("#phoneModalAppend").children().remove();
    $("#phoneModalAppend").append(
        '<div>设备名称：' + name + '</div>' +
        '<div>设备地址：' + address + '</div>' +
        '<div>采集周期：' + period + '</div>' +
        '<div>最大通讯错误数：' + communication + '</div>'
    )
    $("#phoneModal1").modal("open");
}

//取一些摄像机的数据
function getUuid() {
    var get1;
    $.ajax({
        url: host + "/mu/id",
        type: "GET",
        async: false,
        error:function(resp1){
            if (!resp1.responseText){
                alert(resp1.statusText);
            }else{
                alert(resp1.responseText);
            }
        },
        success: function (re) {
            get1 = re;
        }
    });
    return get1.data.uuid;
}

function getCameraList() {
    var camera;
    $.ajax({
        url: host + "/video/",
        type: "GET",
        async: false,
        error:function(resp1){
            if (!resp1.responseText){
                alert(resp1.statusText);
            }else{
                alert(resp1.responseText);
            }
        },
        success: function (res) {
            camera = res;
        }
    })
    return camera.data;
}

function addCloudCamera(info, id, callback) {
    var userInfo = {}
    for (f in info) {
        userInfo[info[f].name] = info[f].value
    }
    var cameraName = id.split("_")[1];
    var data = JSON.stringify({
        cameraId: id,
        cameraName: cameraName,
        rtspUrl: "rtsp://" + userInfo.cameraAccount + ":" + userInfo.cameraPassword + "@" + userInfo.cameraIp + ":554/Streaming/Channels/101/",
        serverUrl: "lab.huayuan-iot.com",
        streamId: userInfo.cameraIp,
        streamName: id,
        rtmpUrl: "rtmp://lab.huayuan-iot.com:9641/live/" + id,
        hlsUrl: "http://lab.huayuan-iot.com:9642/hls/" + id + "/index.m3u8"
    })
    $.ajax({
        url: host + "/video/",
        type: "POST",
        data: data,
        success: function (res) {
            if (typeof callback === "function") {
                callback(res);
            }
        },
        error: function (param) {
            if (param.responseJSON.status == "803") {
                alert("添加摄像机达到上限");
                Materialize.updateTextFields();
            } else if (param.responseJSON.status == "802") {
                alert("此摄像机不存在");
            } else if (param.responseJSON.status == "801") {
                alert("此摄像机已存在");
            }
        }

    })
}

function deleteCloudCamera(id, callback) {
    $.ajax({
        url: host + "/video/" + id,
        type: "DELETE",
        error:function(resp1){
            if (!resp1.responseText){
                alert(resp1.statusText);
            }else{
                alert(resp1.responseText);
            }
        },
        success: function (res) {
            if (typeof (callback) === "function") {
                callback();
            }
        }
    })
}

const testSwitch = (equipment) => {
    const checkType = $(equipment).prop("checked");
    const selectName = $(equipment).parent().parent().attr('data-id');
    const typeName =  $(equipment).attr('data-id');

    let supportEquipment = localStorage.getItem("supportEquipment");
    try {
        supportEquipment = JSON.parse(supportEquipment)
    } catch (error) {
        console.log(error);
        alert("请检查数据格式是否正确！");
        return false;
    }

    if (checkType) {
        supportEquipment[typeName].map((value) => {
            console.log(value)
        })
    }
}

const testModal = () => {
    $('#testModal').modal('open');
    let supportEquipment = localStorage.getItem("supportEquipment");
    let showList = [];
    try {
        supportEquipment = JSON.parse(supportEquipment)
    } catch (error) {
        console.log(error);
        alert("请检查数据格式是否正确！");
        return false;
    }
    for (type in supportEquipment.userChoose) {
        showList.push(supportEquipment.userChoose[type].id);
        $("#testModal .modal-content").append(
            `<div class=${supportEquipment.userChoose[type].id} data-id=${supportEquipment.userChoose[type].id}>` +
            `<span>${supportEquipment.userChoose[type].name}</span>`+
            `</div>`
        )
    }
    for (name in supportEquipment) {
        if (showList.indexOf(name) > -1) {
            supportEquipment[name].map((value,b,c) =>{
                $(`.${name}`).append(
                    `<span>
                    <input id=${name}-${value.id} data-id=${value.id} class="filled-in red" type="checkbox" onclick="testSwitch(this)"/>
                    <label for=${name}-${value.id}>${value.name}</label>
                    </span>`
                )
            })

        }
    }
    
}
