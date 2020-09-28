/* 
    author:Billy
    desc:录像查看
*/


var cameraId = window.location.search;
cameraId = cameraId.substr(8, cameraId.length - 1);
var date = new Date();
var today = date.toLocaleDateString();
var minute = date.getMinutes();
var hour = date.getHours();
if (minute > 30) {
    hour = hour + 1;
}
hour = hour.toString();
today = today.split("/");
if (today[1] < 10) {
    today[1] = "0" + today[1]
}
if (today[2] < 10) {
    today[2] = "0" + today[2]
}
today = today.join("-")
var nowTime = today + " " + hour + ":00";
$.datetimepicker.setLocale('ch')
$('#datetimepicker3').datetimepicker({
    inline: true
});

var noportHost = host.split(":");
noportHost.pop();
noportHost = noportHost.join(":");


// 每次进入页面加载当前日期的全部视频
$(document).ready(function () {
    $('.modal').modal();
    $.ajax({
        url: host + '/record/' + cameraId,
        type: 'GET',
        data: {
            startTime: today + " 00:00",
            endTime: today + " 24:00"
        },
        dataType: 'JSON',
        success: function (res) {
            for (video in res.data) {
                var f1 = res.data[video]
                f = f1.split("_");
                date = f[2].split("T")[0];
                time = f[2].split("T")[1].split(".")[0].replace(/-/g, ":");
                imgUrl = f1.replace(/mp4/, "jpg");
                var noportHost = host.split(":");
                noportHost.pop();
                noportHost = noportHost.join(":");
                $(".cam-right").append(
                    '<div class="col l3 s3 video-card">' +
                    '<div class="card">' +
                    '<div class="card-image waves-effect waves-block waves-light" data-movie=' + f1 + ' onclick="showVideo(this)">' +
                    '<img src="' + noportHost + '/capture/' + cameraId + '/' + imgUrl + '">' +
                    '</div>' +
                    '<div class="card-content">' +
                    '<span class="card-title activator grey-text text-darken-4">' + time +
                    '</span>' +
                    '</div>' +
                    '</div>' +
                    '</div>'
                )
            }
        },
        error: function (res) {
            console.log(res)
        }
    })

})


$("#datetimepicker3").on("change", function () {
    var time = $(this).val().replace(/\//g, "-");
    var time1 = time.split(" ");
    $.ajax({
        url: host + '/record/' + cameraId,
        type: 'GET',
        data: {
            startTime: time,
            endTime: time1[0] + " 24:00"
        },
        timeout:"3000",
        dataType: 'JSON',
        success: function (res) {
            $(".cam-right").children().remove();
            for (video in res.data) {
                var f1 = res.data[video]
                f = f1.split("_");
                date = f[2].split("T")[0];
                time = f[2].split("T")[1].split(".")[0].replace(/-/g, ":");
                imgUrl = f1.replace(/mp4/, "jpg");

                $(".cam-right").append(
                    '<div class="col l3 s3 video-card">' +
                    '<div class="card">' +
                    '<div class="card-image waves-effect waves-block waves-light" data-movie=' + f1 + ' onclick="showVideo(this)">' +
                    '<img src="' + noportHost + '/capture/' + cameraId + '/' + imgUrl + '">' +
                    '</div>' +
                    '<div class="card-content">' +
                    '<span class="card-title activator grey-text text-darken-4">' + time +
                    '</span>' +
                    '</div>' +
                    '</div>' +
                    '</div>'
                )
            }
        },
        error: function (res) {
            console.log(res)
        }
    })
})

function showVideo(item) {
    var host1 = "http://192.168.1.194:8081"
    console.log($(item).attr("data-movie"));
    $("#videoSrc").children().remove();
    $("#videoSrc").append(
        '<video class="responsive-video" controls style="width: 100%;" autoplay="autoplay">' +
            '<source src='+ host1 + '/record/' + cameraId + '/' + $(item).attr("data-movie") +' type="video/mp4">' +
        '</video>'
    )
    $('#modal1').modal('open');
}