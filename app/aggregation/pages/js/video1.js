/* 
    author:Billy
    desc:录像设置
*/
$(document).ready(function () {
    $('#modalVideo').modal();
    var cameraId = window.location.search;
    cameraId = cameraId.substr(8, cameraId.length - 1)
    if (cameraId == ""){
        alert("请选择您需要设置的摄像机！");
        location.href = "cloudvideo.html";
    }
    $.ajax({
        url: host + '/storage/',
        type: 'GET',
        dataType: 'JSON',
        success: function (res) {
            //总容量
            var totalDisc = res.data.total - 200;
            //其他文件
            var otherDisc = parseInt(res.data.otherUsed);
            var other = (100 * otherDisc / totalDisc).toFixed(2);
            //配置容量
            var limitDisc = parseInt(res.data.limit);
            var setted = (100 * limitDisc / totalDisc).toFixed(2);
            //录像容量
            var recordDisc = parseInt(res.data.recordUsed);
            var used = (100 * recordDisc / totalDisc).toFixed(2);
            //剩余容量
            var leaveDisc = totalDisc - otherDisc;
            discInfo = JSON.stringify({
                totalDisc,
                otherDisc,
                limitDisc,
                recordDisc
            })
            sessionStorage.setItem("discInfo", discInfo)

            $(".total-disc").css("width", "100%");
            $(".total-disc").attr("title", totalDisc + "M");

            $(".limit-disc").css("width", setted + "%");
            $(".limit-disc").attr("title", limitDisc + "M");

            $(".slide-bar").css("width", 100 - other + "%");

            $(".used-disc").css("width", used + "%");
            $(".used-disc").attr("title", recordDisc + "M");

            $(".other-disc").css("width", other + "%");
            $(".other-disc").attr("title", otherDisc + "M");

            $("#nowUsed").val(recordDisc)

            $("#rangeBar").attr("value", limitDisc);

            $("#limit").val(limitDisc);
            $("#wqe1").val(recordDisc);
            Materialize.updateTextFields();
        }
    })
    //获取录像列表
    // $.ajax({
    //     url: host + '/record/' + cameraId,
    //     type: 'GET',
    //     dataType: 'JSON',
    //     success: function (res) {
    //         sessionStorage.setItem("videoList", JSON.stringify(res.data))
    //         for (video in res.data) {
    //             var f = res.data[video]
    //             f = f.split("_");
    //             date = f[2].split("T")[0];
    //             $("#videotape").append(
    //                 '<li class="videotape" onclick="videoModal(this)" title="' + res.data[video] + '">' + f[1] + "..." + f[2] + '</li>'
    //             )
    //         }
    //     },
    //     error: function (res) {
    //         console.log(res)
    //     }
    // })
    //获取设置时间
    $.ajax({
        url: host + '/video/' + cameraId,
        type: 'GET',
        dataType: 'JSON',
        success: function (res) {
            console.log(res)
            if (res.data.record.mode == "enable") {
                //设置了
                $("#videoSwitch").prop("checked", true);
                $("#videoStart").val(res.data.record.startTime);
                $("#videoEnd").val(res.data.record.endTime);
                $("#videoStart").removeAttr("disabled");
                $("#videoEnd").removeAttr("disabled");
            } else if (res.data.record.mode == "disable") {
                //未设置
                $("#videoSwitch").prop("checked", false);
            } else if (res.data.record.mode == "wholeday") {
                //全天
                $("#videoSwitch").prop("checked", true);
                $("#videoStart").val("全天");
                $("#videoStart").removeAttr("disabled");
            }

            if (res.data.sync.mode == "enable") {
                //设置了
                $("#cloudSwitch").prop("checked", true);
                $("#cloudStart").val(res.data.sync.startTime);
                $("#cloudEnd").val(res.data.sync.endTime);
                $("#cloudStart").removeAttr("disabled");
                $("#cloudEnd").removeAttr("disabled");
            } else if (res.data.sync.mode == "disable") {
                //未设置
                $("#cloudSwitch").prop("checked", false);
            } else if (res.data.sync.mode == "wholeday") {
                //全天
                $("#cloudSwitch").prop("checked", true);
                $("#cloudStart").val("全天");
                $("#cloudStart").removeAttr("disabled");
            }

        }
    })
})


var slider = document.getElementById('slider');
var limit = document.getElementById('limit');
var limitDisc = JSON.parse(sessionStorage.getItem("discInfo")).limitDisc;
var otherDisc = JSON.parse(sessionStorage.getItem("discInfo")).otherDisc;
var totalDisc = JSON.parse(sessionStorage.getItem("discInfo")).totalDisc;
var leaveDisc = totalDisc - otherDisc;

noUiSlider.create(slider, {
    start: [limitDisc],
    tooltips: [true],
    connect: [true, false],
    range: {
        'min': 0,
        'max': leaveDisc
    },
    format: wNumb({
        decimals: 0
    })
});

slider.noUiSlider.on("update", function (value, handle) {
    if (parseInt(value[0]) + parseInt(otherDisc) > parseInt(totalDisc)) {
        alert("设置容量超出总容量，请重新调整录像空间大小！");
        return false;
    }
    $(".limit-disc").css("width", 100 * value[0] / totalDisc + "%");
    if (handle == 0) {
        $("#limit").val(value[0]);
    }
})

$(document).on("change", "#limit", function () {
    slider.noUiSlider.set([this.value]);
})

// function videoModal(video) {
//     var host1 = "192.168.1.194";
//     var cameraId = window.location.search;
//     cameraId = cameraId.substr(8, cameraId.length - 1)
//     $("#modalVideo").modal('open');
//     $("#videoAddress").children().remove();
//     $("#videoAddress").append(
//         '<video autoplay="autoplay" class="responsive-video" controls>' +
//         '<source src="' + "http://" + host1 + "/record/" + cameraId + "/" + $(video).attr("title") + '" type="video/mp4">' +
//         '</video>'
//     )
// }

$("#lookVideo").on("click", function () {
    var cameraId = window.location.search;
    cameraId = cameraId.substr(8, cameraId.length - 1)
    location.href = "video2.html?camera=" + cameraId;
})

$("#searchVideo").on("click", function () {
    var cameraId = window.location.search;
    cameraId = cameraId.substr(8, cameraId.length - 1)
    var start = Date.parse($("#datepickS").val());
    var end = Date.parse($("#datepickE").val());
    if (start >= end) {
        alert("您选择的时间段不正确！请重新选择");
        $("#datepickS").val("");
        $("#datepickE").val("");
        return false;
    } else {
        $.ajax({
            url: host + '/record/' + cameraId,
            type: 'GET',
            data: {
                startTime: $("#datepickS").val(),
                endTime: $("#datepickE").val()
            },
            dataType: "JSON",
            success: function (res) {
                var videoList = res.data;
                console.log(videoList)
                $("#videotape").children().remove();
                for (video in videoList) {
                    var f = videoList[video]
                    f = f.split("_");
                    date = f[2].split("T")[0];
                    $("#videotape").append(
                        '<li class="videotape" onclick="videoModal(this)" title="' + videoList[video] + '">' + f[1] + "..." + f[2] + '</li>'
                    )
                }
            }
        })
    }
})

$('#videoStart').timepicker({
    'timeFormat': 'H:i',
    'noneOption': [{
        'label': '全天',
        'className': 'allday',
        'value': '全天'
    }]
});

$('#videoEnd').timepicker({
    'timeFormat': 'H:i'
});

$('#cloudStart').timepicker({
    'timeFormat': 'H:i',
    'noneOption': [{
        'label': '全天',
        'className': 'allday',
        'value': '全天'
    }]
});

$('#cloudEnd').timepicker({
    'timeFormat': 'H:i'
});

$("#videoStart").on("change", function () {
    if ($(this).val() == "全天") {
        $("#videoEnd").attr("disabled", "true");
    } else {
        $("#videoEnd").removeAttr("disabled");
    }
})

$("#videoEnd").on("change", function () {
    var videoStart = $("#videoStart").val().split(":");
    var videoEnd = $("#videoEnd").val().split(":");
    if (videoStart[0] > videoEnd[0]) {
        $("#videoEnd").val("");
        alert("请选择正确的结束时间");
        return false;
    } else if (videoStart[0] == videoEnd[0] && videoStart[1] > videoEnd[1]) {
        $("#videoEnd").val("");
        alert("请选择正确的结束时间");
        return false;
    } else if (videoStart[0] == videoEnd[0] && videoStart[1] == videoEnd[1]) {
        $("#videoEnd").val("");
        alert("请选择正确的结束时间");
        return false;
    }
})

$("#cloudStart").on("change", function () {
    if ($(this).val() == "全天") {
        $("#cloudEnd").attr("disabled", "true");
    } else {
        $("#cloudEnd").removeAttr("disabled");
    }
})

$("#cloudEnd").on("change", function () {
    var cloudStart = $("#cloudStart").val().split(":");
    var cloudEnd = $("#cloudEnd").val().split(":");
    if (cloudStart[0] > cloudEnd[0]) {
        $("#cloudEnd").val("");
        alert("请选择正确的结束时间");
        return false;
    } else if (cloudStart[0] == cloudEnd[0] && cloudStart[1] > cloudEnd[1]) {
        $("#cloudEnd").val("");
        alert("请选择正确的结束时间");
        return false;
    } else if (cloudStart[0] == cloudEnd[0] && cloudStart[1] == cloudEnd[1]) {
        alert("请选择正确的结束时间");
        $("#cloudEnd").val("");
        return false;
    }
})

$("#videoSwitch").on("change", function () {
    if ($(this).prop("checked") == true) {
        $("#videoStart").removeAttr("disabled");
        $("#videoEnd").removeAttr("disabled");
    } else {
        $("#videoStart").attr("disabled", "true");
        $("#videoEnd").attr("disabled", "true");
    }
})

$("#cloudSwitch").on("change", function () {
    if ($(this).prop("checked") == true) {
        $("#cloudStart").removeAttr("disabled");
        $("#cloudEnd").removeAttr("disabled");
    } else {
        $("#cloudStart").attr("disabled", "true");
        $("#cloudEnd").attr("disabled", "true");
    }
})

$("#slideShow").on("mouseover", function () {
    $('.limit-disc').hide();
    $(".slide-bar").css("display", "inline-block");
})

$("#slideShow").on("mouseout", function () {
    $('.limit-disc').show();
    $(".slide-bar").css("display", "none");
})

$("#limit").on("change", function () {
    var info = sessionStorage.getItem("discInfo");
    var total = JSON.parse(info).totalDisc;
    var other = JSON.parse(info).otherDisc;
    // var oldLimit = JSON.parse(info).limitDisc
    var limit = parseInt($(this).val());
    // var used = $("#nowUsed").val();
    if (limit + other > total) {
        $(this).val(JSON.parse(info).limitDisc)
        alert("设置容量超出总容量，请重新调整录像空间大小！");
        return false;
    }

})

// $("#datepick").on("change", function () {
//     var videoList = JSON.parse(sessionStorage.getItem("videoList"));
//     $("#videotape").children().remove();
//     for (video in videoList) {
//         var f = videoList[video]
//         f = f.split("_");
//         date = f[2].split("T")[0];
//         if ($(this).val() == date) {
//             $("#videotape").append(
//                 '<li class="videotape" onclick="videoModal(this)" title="' + videoList[video] + '">' + f[1] + "..." + f[2] + '</li>'
//             )
//         } else if ($(this).val() == "") {
//             $("#videotape").append(
//                 '<li class="videotape" onclick="videoModal(this)" title="' + videoList[video] + '">' + f[1] + "..." + f[2] + '</li>'
//             )
//         }

//     }
// })

$("#submit1").on("click", function () {
    var limit = JSON.stringify({
        limit: $("#limit").val()
    });
    var recordStart = $("#videoStart").val();
    var recordEnd = $("#videoEnd").val();
    var syncStart = $("#cloudStart").val();
    var syncEnd = $("#cloudEnd").val();
    var videoSwitch = $("#videoSwitch").prop("checked");
    var cloudSwitch = $("#cloudSwitch").prop("checked");
    var cameraId = window.location.search;
    cameraId = cameraId.substr(8, cameraId.length - 1);

    var info = sessionStorage.getItem("discInfo");
    // var totalDisc = JSON.parse(info).totalDisc;
    var oldLimit = JSON.parse(info).limitDisc;
    var used = parseInt($("#nowUsed").val());
    var setLimit = parseInt($("#limit").val());


    if (setLimit < used) {
        var a = confirm("您设置录像存储容量小于已用容量，确认将会删除较早文件！")
        if (a !== true) {
            location.reload()
            return false;
        }
    }

    if (cloudSwitch == true) {
        if (syncStart !== "全天" && syncEnd == "") {
            alert("请选择结束时间！");
            return false;
        }
    }

    if (recordStart !== "全天" && videoSwitch == true) {
        if (recordEnd == "") {
            alert("请选择结束时间！");
            return false;
        }
    }

    if (recordStart == "全天") {
        recordStart = "00:00"
        recordEnd = "24:00"
    }
    if (syncStart == "全天") {
        syncStart = "00:00"
        syncEnd = "24:00"
    }
    var setTime = JSON.stringify({
        cameraId: cameraId,
        record: {
            enable: videoSwitch,
            startTime: recordStart,
            endTime: recordEnd
        },
        sync: {
            enable: cloudSwitch,
            startTime: syncStart,
            endTime: syncEnd
        }
    })

    //设置录像存储上限
    $.ajax({
        url: host + "/storage/",
        data: limit,
        type: "PUT",
        success: function (re) {
            console.log(re);
        }
    })

    //设置录像时间
    $.ajax({
        url: host + "/record/",
        data: setTime,
        type: "PUT",
        success: function (re) {
            console.log(re);
            alert("设置成功！")
        }
    })
})