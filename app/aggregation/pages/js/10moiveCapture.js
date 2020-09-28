// 绘制视频列表
function drawVedioList(data) {
    let vedioList = [];
    vedioList = data.sList.map(item => {
        const arr_result = item.split("-");
        let initTime = arr_result[1].split(".")[0];
        let time = `${initTime.slice(0, 4)}-${initTime.slice(
            4,
            6
        )}-${initTime.slice(6, 8)} ${initTime.slice(8, 10)}:${initTime.slice(
            10,
            12
        )}:${initTime.slice(12, 14)}`;
        return {
            host: host + "/capture/" + item,
            time: time
        };
    });
    $(".picture-box .picture-right")
        .children()
        .remove();
    for (var i = 0; i < vedioList.length; i++) {
        var dom =
            '<li class="vedio"><div class="vedio-image" style="background-image:url(' +
            vedioList[i].host +
            ')"></div>' +
            '<p class="time">' +
            vedioList[i].time +
            "</p></li>";
        $(".picture-box .picture-right").append(dom);
    }
}
// 获取视频列表
function getVedioList(nowhost) {
    var posthost = "";
    nowhost ? (posthost = nowhost) : (posthost = $("#menu .collection-item:first").attr("host"));
    $.post(host + "/capture/showPicName", { host: posthost }, function(
        data,
        status
    ) {
        drawVedioList(data);
    });
}
// 获取摄像机列表
function getCameraList() {
    $.get(host + "/mu/", function(data, status) {
        const ports = data.data[0].ports;
        let allProtocolCamera = [];
        let camera_list = [];
        ports.forEach(item => {
            item.protocol === "protocol-camera"
                ? allProtocolCamera.push(item)
                : "";
        });
        allProtocolCamera.forEach(item => {
            camera_list = camera_list.concat(item.sampleUnits);
        });
        for (var i = 0; i < camera_list.length; i++) {
            var dom =
                '<li class="collection-item" host=' +
                camera_list[i].setting.host +
                ">" +
                camera_list[i].name +
                "</li>";
            $("#menu").append(dom);
        }
        getVedioList();
        let destroyInterval = setInterval(()=>{
            getVedioList();
        },15000)
        $("#menu").on("click", ".collection-item", function(self) {
            const nowhost = $(self.target).attr("host");
            clearInterval(destroyInterval);
            getVedioList(nowhost);
            destroyInterval = setInterval(()=>{
                getVedioList(nowhost);
            },15000)
        });
    });
}
function init() {
    getCameraList();
}
$(document).ready(function() {
    init();
});