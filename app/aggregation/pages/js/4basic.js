$(document).ready(function () {
  $.ajax({
    url: host + "/mu/id",
    type: "GET",
    dataType: "JSON",
    success: function (res) {
      //res是请求返回的字段，可打印出来，获取自己所需要的值判断是否登陆成功
      $("#collectorid").val(res.data.uuid);
      // console.log(res)
    },
  });
  $.ajax({
    url: host + "/mqtt/",
    type: "GET",
    dataType: "JSON",
    success: function (res) {
      $("#mqttip").val(res.data.host);
      $("#mqttport").val(res.data.port);
    },
    error: function (resp1) {
      $("#mqttip").val(resp1.data.host);
      $("#mqttport").val(resp1.data.port);
      console.log(resp1.data);
    },
  });
});

$("#4basic").on("click", function () {
  if ($("#mqttip").val() == "" || $("#mqttport").val() == "") {
    alert("您有未填写的内容！");
    return false;
  }
  var data = JSON.stringify({
    host: $("#mqttip").val(),
    port: $("#mqttport").val(),
  });
  $.ajax({
    url: host + "/mqtt/",
    data: data,
    type: "PUT",
    dataType: "JSON",
    success: function (res) {
      console.log(res.data);
      location.href = "visualization.html";
    },
    error: function (resp1) {
      if (!resp1.responseText) {
        alert(resp1.statusText);
      } else {
        alert(resp1.responseText);
      }
    },
  });
});
