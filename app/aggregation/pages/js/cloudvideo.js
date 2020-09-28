/* 
    加载已有摄像机列表
*/
$(document).ready(function () {
    $('.modal').modal();
    $(".button-collapse").sideNav();
    $.ajax({
        url: host + "/video/",
        type: "GET",
        success: function (res) {
            var number = [];
            var list = res.data
            if(list.length>0){
                for(num in list){
                    var e = list[num].cameraId.slice(-1);
                    number.push(e);
                }
                var maxNum = Math.max.apply(null,number);
            }else{
                var maxNum = 0;
            }
            if ($("#camera1").attr("data-type") == "new") {
                $("#cameraId").val("camera" + (parseInt(maxNum)+parseInt(1)))
            }
            if (list.length > 0) {
                $(".cam-info").remove();
                for (a in list) {
                    $("#menu").append(
                        '<li class="cam-info">' +
                        '<a class="waves-effect" onclick="cameraSetting(this)" data-type="old" data-id=' + list[a].cameraId + '>' + list[a].cameraName + '</a>' +
                        '</li>'
                    )
                }
            }
            Materialize.updateTextFields();
        }
    })
    $.ajax({
        url: host + "/mu/id",
        type: "GET",
        success: function (re) {
            sessionStorage.setItem("uuid", re.data.uuid);
        }
    })
})
/* 
    事件
*/


/* 新增摄像机，暂时废除
    $("#addCamera").on("click", function () {
    $(".cam-info").last().parent().append(
        '<li class="cam-info">' +
        '<a class="waves-effect" onclick="cameraSetting(this)" data-type="new">' + '新摄像机' + '</a>' +
        '</li>'
    )
}) */

//保存摄像机
$("#submit1").on("click", function () {
    var camera1 = $("#camera1").serializeArray();
    var data1 = {};
    for (i in camera1) {
        var t = camera1[i].name;
        if (y !== "") {
            var y = camera1[i].value;
        } else {
            alert("您有未填写信息！");
            return false;
        }
        data1[t] = y
    }
    var subType = $("#camera1").attr("data-type");
    if (subType == "old") {
        var requestType = "PUT"
    } else if (subType == "new") {
        var requestType = "POST"
    }
    var data = JSON.stringify(data1);
    $.ajax({
        url: host + "/video/",
        type: requestType,
        data: data,
        success: function (res) {
            var list = res.data
            if (list.length > 0) {
                $(".cam-info").remove();
                for (a in list) {
                    $("#menu").append(
                        '<li class="cam-info">' +
                        '<a class="waves-effect" onclick="cameraSetting(this)" data-type="old" data-id=' + list[a].cameraId + '>' + list[a].cameraName + '</a>' +
                        '</li>'
                    )
                }
            }
            $('#modalAddCamera').modal('open');
        },
        error: function (param) {
            if (param.responseJSON.status == "803") {
                alert("添加摄像机达到上限");
                $("#camera1")[0].reset();
                Materialize.updateTextFields();
            } else if (param.responseJSON.status == "802") {
                alert("此摄像机不存在");
            } else if (param.responseJSON.status == "801") {
                alert("此摄像机已存在");
            }
        }

    })
})

//删除摄像机
$(".button-grop").on("click", "#deleteCamera", function () {
    var id = $("#camera1").attr("data-id");
    $.ajax({
        url: host + "/video/" + id,
        type: "DELETE",
        success: function (res) {
            var list = res.data
            if (list.length > 0) {
                $(".cam-info").remove();
                for (a in list) {
                    $("#menu").append(
                        '<li class="cam-info">' +
                        '<a class="waves-effect" onclick="cameraSetting(this)" data-type="old" data-id=' + list[a].cameraId + '>' + list[a].cameraName + '</a>' +
                        '</li>'
                    )
                }
            }
        }
    })
    $("#camera1")[0].reset();
    $("#camera1").attr("data-type","new");
    $("#camera1").attr("data-id","");
    $("#pageTitle").html("新增摄像机配置");
    Materialize.updateTextFields();
    $('#modalDeleteCamera').modal('open');
})

// 手机版侧边栏
$(".mobile-back").on("click", function () {
    $("#camList").fadeOut()
})
$(".toc").on("click", function () {
    $("#camList").fadeIn()
})

//自动生成streamId和通道名称
$("#rtspUrl").on("change", function () {
    var uid = sessionStorage.getItem("uuid")
    var val = $(this).val();
    if (val.indexOf("@") > 0) {
        var removeHead = val.substr(parseInt(val.indexOf("@") + 1));
    } else {
        var removeHead = val.substr(7);
    }
    var local = removeHead.indexOf(":")
    var url = removeHead.slice(0, local)
    if (url.length > 0) {
        var test = /^((25[0-5]|2[0-4]\d|((1\d{2})|([1-9]?\d)))\.){3}(25[0-5]|2[0-4]\d|((1\d{2})|([1-9]?\d)))$/
        real = test.test(url);
        if (real == false) {
            alert("请输入正确的RTSP地址格式！");
            $(this).val("");
            $("#streamId").val("");
            $("#streamName").val("");
            Materialize.updateTextFields();
            return false;
        } else {
            if ($("#camera1").attr("data-type") == "new") {
                $("#streamId").val(url);
                $("#streamName").val(uid + "_" + $("#cameraId").val());
                Materialize.updateTextFields();
            }
        }
    } else {
        alert("请输入正确的RTSP地址格式！");
        $(this).val("");
        $("#streamId").val("");
        $("#streamName").val("");
        Materialize.updateTextFields();
        return false;
    }
})
//自动生成播放地址
$("#serverUrl").on("change", function () {
    var address = $(this).val();
    if ($("#streamId").val() !== "" && $("#streamName").val() !== "" && $("#serverUrl").val() !== "") {
        if(address.slice(0,7) == "http://"){
            serviceAddr = address.substr(7);
        }else if(address.slice(0,8) == "https://"){
            serviceAddr = address.substr(8);
        }else{
            serviceAddr = $(this).val();
        }
        //分隔
        if($("#camera1").attr("data-type") == "new"){
            var uid = sessionStorage.getItem("uuid") + "_" + $("#cameraId").val();
        }else if($("#camera1").attr("data-type") == "old"){
            var uid = $("#cameraId").val();
        }
        $("#rtmpUrl").val("rtmp://" + serviceAddr + ":9641/live/" + uid);
        $("#hlsUrl").val("http://" + serviceAddr + ":9642/hls/" + uid + "/index.m3u8");
        Materialize.updateTextFields();
    } else {
        $("#rtmpUrl").val("");
        $("#hlsUrl").val("");
        Materialize.updateTextFields();
    }
})

$("#video1").on("click",function () {
    var camera = $("#cameraId").val();
    location.href = "video1.html?camera=" + camera
  })

function cameraSetting(item) {
    $(".cam-right").show();
    $(".cam-right-default").remove();
    if (document.body.clientWidth < 600) {
        $("#camList").fadeOut();
    }
    $("#deleteCamera").remove();
    $(item).addClass("selected");
    $(item).parent().siblings().children().removeClass("selected");
    var id = $(item).attr("data-id");
    var type = $(item).attr("data-type");
    $("#camera1").attr("data-type", type);
    $("#camera1").attr("data-id", id);
    $("#camera1")[0].reset();
    Materialize.updateTextFields();
    $("#test1").show();
    $("#notice").hide();
    $("#submit1").show();
    $("#occupation1").show();
    $("#occupation2").hide();
    if (type == "old") {
        // $("#deleteCamera").show();
        // $(".button-grop").append(
        //     '<a id="deleteCamera" class="waves-effect btn z-depth-1">删除</a>'
        // )
        $.ajax({
            url: host + "/video/" + id,
            type: "GET",
            success: function (res) {
                $("#pageTitle").html(res.data.cameraName + "配置")
                $("#cameraId").val(res.data.cameraId);
                $("#cameraName").val(res.data.cameraName);
                $("#rtspUrl").val(res.data.rtspUrl);
                $("#streamId").val(res.data.streamId);
                $("#serverUrl").val(res.data.serverUrl);
                $("#streamName").val(res.data.streamName);
                $("#rtmpUrl").val(res.data.rtmpUrl);
                $("#hlsUrl").val(res.data.hlsUrl);
                /* $("#userName").val(res.data.userName);
                $("#password").val(res.data.password); */
                $("#camera1").attr("data-type", "old");
                Materialize.updateTextFields();
            }
        })
    }
}