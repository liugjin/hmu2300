function getSrvAddr() {
  return "http://" + window.location.host.split(":")[0] + ":8090";
  // return 'http://192.168.0.213:8090';
}

var host = getSrvAddr()
$("#3lan").on("click", function () {
  var lanIp = $("#lanipaddr").val();
  var mask = $("#lannetmask").val();

  if (lanIp == "") {
    alert("请输入正确的IP地址！");
    return false;
  }
  if (mask == "") {
    alert("请输入正确的子网掩码！");
    return false;
  }
  var g = confirm("修改LAN将会修改设备访问地址，请谨慎设置！");
  if (g == true) {
    var newIp = lanIp;
    if (location.port != "80") {
      newIp = lanIp + ":" + location.port;
    }
    var data = JSON.stringify({
      "ip": lanIp,
      "mask": mask
    })
    $.ajax({
      url: host + "/lan/",
      data: data,
      type: "PUT",
      dataType: "JSON",
      success: function (res) {
        if (res.status == "0") {
          console.log("请求成功");
          // 修改成功，等待机器生效
          // TODO:优化定时器跳转的等待
          setTimeout(function () {
            location.href = "http://" + newIp + "/4basic.html"
          }, 1000);
        } else {
          alert("配置失败！")
        }
      },
      error: function (resp1) {
        console.log("失败");
        if (!resp1.responseText) {
          alert(resp1.statusText);
        } else {
          alert(resp1.responseText);
        }
      }
    })

  }
})

function selfModal(title, message, callback) {
  this.button = function (ok) {
    callback(ok)
  }
  $("#modalHead").text(title);
  $("#modalContent").append(message);
  $("#selfModal").modal("open");
}
