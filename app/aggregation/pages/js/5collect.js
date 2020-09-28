$(document).ready(function() {
  var equipId = null
  $.ajax({
    url: host + '/mu/',
    type: 'GET',
    dataType: 'json',
    success: function(res) {
      var muport = res.data[0].ports;
      $("#equipid").val(res.data[0].id);
      $("#equipname").val(res.data[0].name);
      $("#equipUnit").html('<img class="circle-img" src="../pic/circle.png">' + "监控单元名称：" + res.data[0].name);
      $("#equipUnit").attr("data-mu", res.data[0].id)
      for (var i = 0; i < muport.length; i++) {
        var r = []; //用来保存id
        if (muport[i].id !== "sp-hmu") {
          $("#listTree").append(
            '<div ' + 'class=' + '"' + 'tree' + [i + 1] + '"' + '>' +
            '<div onclick="treeOpen(this);loadAttribute(this)" id=' + '"' + muport[i].id + '"' + ' class=' + 'dad-tree' + '><img class="circle-img" src="../pic/circle.png">采集端口名称：' + muport[i].name + '</div>' +
            '</div>'
          )
        }
        if (muport[i].sampleUnits !== null) {
          for (var j = 0; j < muport[i].sampleUnits.length; j++) {
            $('#listTree ' + '.tree' + [i + 1]).append(
              '<div onclick="loadAttribute(this)" class="a-class" id=' + muport[i].sampleUnits[j].id + '>采集单元名称：' + muport[i].sampleUnits[j].name + '</div>'
            )
          }
        }
        $(".a-class").hide();
      }
      Materialize.updateTextFields();
      $('select').material_select();
    }
  });
  $.ajax({
    url: host + '/pl/',
    type: 'GET',
    dataType: 'json',
    success: function(res1) {

      //协议填充
      for (var i = 0; i < res1.data.length; i++) {
        $("#protocolOption").append(
          "<option value='" + res1.data[i] + "'>" + res1.data[i] + "</option>"
        )
      }
      $('select').material_select();
    }
  });
  $.ajax({
    url: host + '/el/',
    type: 'GET',
    dataType: 'json',
    success: function(res2) {
      for (var i = 0; i < res2.data.length; i++) {
        $("#equipLab").append(
          "<option value='" + res2.data[i] + "'>" + res2.data[i] + "</option>"
        )
      }
      $('select').material_select();
    }
  });
});
//初始化结束

//树打开与关闭

function treeOpen(item) {
  $(item).parent().find(".a-class").toggle();
}

function loadAttribute(it) {
  $(it).addClass("ts-active");
  //单元去选项
  $(it).siblings().removeClass("ts-active");
  //端口去选项
  $(it).parent().siblings().children().removeClass("ts-active");
  //设备去选项
  $(it).parent().parent().siblings().removeClass("ts-active");
  $(it).siblings().children().children().removeClass("ts-active");
  equipId = $(it).attr("id");
  $.ajax({
    url: host + '/mu/',
    type: 'GET',
    dataType: 'json',
    success: function(respone) {
      console.log("respone",respone)
      var back = respone.data[0];
      var ports = [];
      var units = [];
      var dad = [];
      var su = [];
      //判断设备和端口
      if ($(it).hasClass("equipUnit")) { //监控单元
        $("#equipid").val(back.id);
        $("#equipname").val(back.name);
        $("#collecterequip").show();
        $("#collecterport").hide();
        $("#collecterunit").hide();
      } else if ($(it).hasClass("dad-tree")) { //端口
        for (var i = 0; i < back.ports.length; i++) {
          if (equipId == back.ports[i].id) {
            ports.push(back.ports[i])
            var sp = [];
            for (var index in back.ports[i].setting) {
              sp.push({
                name: index,
                val: back.ports[i].setting[index]
              });
            }
          }
        }
        //切换窗口
        $(".addvariable").remove(); //清空设备信息列表
        for (var e = 0; e < sp.length; e++) { //设备信息填充
          $("#spset").append(
            '<tr class="addvariable">' +
            '<td>' +
            '<input type="text" value="' + sp[e].name + '" />' +
            '</td>' +
            '<td>' +
            '<input type="text" value="' + sp[e].val + '" />' +
            '</td>' +
            '<td class="delete-row cter">' +
            '<img src="../pic/delete1.png">' +
            '</td>' +
            '</tr>'
          )
        }
        $("#collecterport").show();
        $("#collecterunit").hide();
        $("#collecterequip").hide();
        $("#collecterport").attr("data-id", ports[0].id);
        $("#collecterport").attr("data-type", "old");
        $("#cpid").val(ports[0].id);
        $("#cpid").attr("readonly", "readonly");
        $("#cpname").val(ports[0].name);
        $("#portswitch").prop("checked", ports[0].enable);
        $("#symbolOptin").val(ports[0].symbol);
        $("#protocolOption").val(ports[0].protocol);
        Materialize.updateTextFields();
        $('select').material_select();
      } else { //单元
        for (var i = 0; i < back.ports.length; i++) {
          if (back.ports[i].sampleUnits !== null) {
            for (var k = 0; k < back.ports[i].sampleUnits.length; k++) {
              if (equipId == back.ports[i].sampleUnits[k].id) {
                for (var index1 in back.ports[i].sampleUnits[k].setting) {
                  su.push({
                    name: index1,
                    val: back.ports[i].sampleUnits[k].setting[index1]
                  });
                }
                units.push(back.ports[i].sampleUnits[k]);
                dad.push(back.ports[i].id)
              }
            }
          }
        }
        //切换窗口
        $(".addvariable").remove(); //清空设备信息列表
        for (var e = 0; e < su.length; e++) { //设备信息填充
          $("#suset").append(
            '<tr class="addvariable">' +
            '<td>' +
            '<input type="text" value="' + su[e].name + '" />' +
            '</td>' +
            '<td>' +
            '<input type="text" value="' + su[e].val + '" />' +
            '</td>' +
            '<td class="delete-row cter">' +
            '<img src="../pic/delete1.png">' +
            '</td>' +
            '</tr>'
          )
        }
        $("#collecterport").hide();
        $("#collecterunit").show();
        $("#collecterequip").hide();
        $("#collecterport").attr("data-id", dad[0])
        $("#collecterunit").attr("data-type", "old");
        $("#cuid").val(units[0].id);
        $("#cuid").attr("readonly", "readonly");
        $("#cuname").val(units[0].name);
        $("#cuperiod").val(units[0].period);
        $("#cutimeout").val(units[0].timeout);
        $("#cudelay").val(units[0].delay);
        $("#cuthrottle").val(units[0].throttle);
        $("#cumaxNum").val(units[0].maxCommunicationErrors);
        $("#unitswitch").prop("checked", units[0].enable);
        $("#equipLab").val(units[0].element);
        Materialize.updateTextFields();
        $('select').material_select();
      }
      // units是当前的选项内容
    }
  })


}
//新增采集端口
$("#newport").on("click", function() {
  $("#collecterport").show();
  $("#collecterunit").hide();
  $("#collecterequip").hide();
  $("#collecterport").attr("data-type", "new");
  $(".addvariable").remove();
  $(".ts-active").removeClass("ts-active");
  $("#cpid").val("");
  $("#cpid").removeAttr("readonly");
  $("#cpname").val("");
  $("#portswitch").prop("checked", "flase");
  $("#symbolOptin").val("");
  $("#protocolOption").val("");
  Materialize.updateTextFields();
  $('select').material_select();
})

//新增采集单元
$("#newunit").on("click", function() {
  $("#collecterport").hide();
  $("#collecterunit").show();
  $("#collecterequip").hide();
  $("#collecterunit").attr("data-type", "new");
  $(".ts-active").removeClass("ts-active");
  $(".addvariable").remove();
  $("#cuid").val("");
  $("#cuid").removeAttr("readonly");
  $("#cuname").val("");
  $("#cuperiod").val("");
  $("#cutimeout").val("");
  $("#cudelay").val("");
  $("#cuthrottle").val("");
  $("#cumaxNum").val("");
  $("#unitswitch").prop("checked", "flase");
  Materialize.updateTextFields();
  $('select').material_select();
})

//保存采集端口
$("#saveport").on("click", function() {
  var portSet = {};
  $("#spset tr:gt(0)").each(function() {
    var tr = $(this);
    var v = tr.find("td").eq(0).find("input").val();
    var n = tr.find("td").eq(1).find("input").val();
    if (v == "baudRate" || v == "keyNumber") {
      portSet[v] = parseInt(n)
    } else {
      portSet[v] = n
    }
  })
  var addDate = JSON.stringify({
    "muid": $("#equipUnit").attr("data-mu"),
    "id": $("#cpid").val(),
    "symbol": $("#symbolOptin").val(),
    "protocol": $("#protocolOption").val(),
    "name": $("#cpname").val(),
    "enable": $("#portswitch").prop("checked"),
    "setting": portSet
  })
  if ($("#collecterport").attr("data-type") == "new") {
    var tp = "POST";
  } else if ($("#collecterport").attr("data-type") == "old") {
    var tp = "PUT";
  }
  $.ajax({
    url: host + "/sp/",
    type: tp,
    data: addDate,
    error: function(resp) {
      console.log(resp.responseText);
      alert(resp.responseText);
    },
    success: function(param) {
      console.log("新增/修改成功");
      location.reload();
    }
  })
  // setTimeout(function () {
  //     location.reload()
  // }, 2000)
})
//保存采集单元
$("#saveunit").on("click", function() {
  var portSet = {};
  $("#suset tr:gt(0)").each(function() {
    var tr = $(this);
    var v = tr.find("td").eq(0).find("input").val();
    var n = tr.find("td").eq(1).find("input").val();
    if (v == "address") {
      portSet[v] = parseInt(n)
    } else {
      portSet[v] = n
    }
  })
  var addDate = JSON.stringify({
    "muid": $("#equipUnit").attr("data-mu"),
    "spid": $("#collecterport").attr("data-id"),
    "id": $("#cuid").val(),
    "name": $("#cuname").val(),
    "period": parseInt($("#cuperiod").val()),
    "timeout": parseInt($("#cutimeout").val()),
    "maxCommunicationErrors": parseInt($("#cumaxNum").val()),
    "element": $("#equipLab").val(),
    "enable": $("#unitswitch").prop("checked"),
    "setting": portSet
  })
  if ($("#collecterunit").attr("data-type") == "new") {
    var tp = "POST";
  } else if ($("#collecterunit").attr("data-type") == "old") {
    var tp = "PUT";
  }
  $.ajax({
    url: host + "/su/",
    type: tp,
    data: addDate,
    error: function(resp) {
      console.log(resp.responseText);
      alert(resp.responseText);
    },
    error: function(resp) {
      console.log(resp.responseText);
      alert(resp.responseText);
    },
    success: function(param) {
      console.log("新增/修改成功");
      equipId = Number(equipId)+1
      // location.reload()
    }
  })
  // setTimeout(function () {
  //     location.reload()
  // }, 2000)
})

//复制采集单元
$("#copy").on("click", function() {
  console.log("equipId",equipId)
  $.ajax({
    url: host + '/mu/',
    type: 'GET',
    dataType: 'json',
    success: function(respone) {
      console.log("respone",respone)
      var back = respone.data[0];
      var units = [];
      var dad = [];
      var su = [];
      for (var i = 0; i < back.ports.length; i++) {
        if (back.ports[i].sampleUnits !== null) {
          for (var k = 0; k < back.ports[i].sampleUnits.length; k++) {
            if (equipId == back.ports[i].sampleUnits[k].id) {
              for (var index1 in back.ports[i].sampleUnits[k].setting) {
                su.push({
                  name: index1,
                  val: back.ports[i].sampleUnits[k].setting[index1]
                });
              }
              units.push(back.ports[i].sampleUnits[k]);
              dad.push(back.ports[i].id)
            }
          }
        }
      }
      arr=units[0].name.split("-")
      console.log("arr",arr)
      let unitsName = ""
      let num = null
      if(arr.length ==1){
        unitsName = units[0].name+"-"+1
        console.log("unitsName",unitsName)
      }else if(arr.length == 2){
        num = Number(arr[1])+1
        unitsName = arr[0]+"-"+ num
      }
      $("#collecterport").hide();
      $("#collecterunit").show();
      $("#collecterequip").hide();
      $("#collecterport").attr("data-id", dad[0])
      $("#collecterunit").attr("data-type", "new");
      $(".ts-active").removeClass("ts-active");
      $("#cuid").val(Number(units[0].id)+1);
      $("#cuid").removeAttr("readonly");
      $("#cuname").val(unitsName);
      $("#cuperiod").val(units[0].period);
      $("#cutimeout").val(units[0].timeout);
      $("#cudelay").val(units[0].delay);
      $("#cuthrottle").val(units[0].throttle);
      $("#cumaxNum").val(units[0].maxCommunicationErrors);
      $("#unitswitch").prop("checked", units[0].enable);
      $("#equipLab").val(units[0].element);
      Materialize.updateTextFields();
      $('select').material_select();
    }
  })
})




//删除采集端口
$("#delport").on("click", function() {
  var delDate = JSON.stringify({
    "muid": $("#equipUnit").attr("data-mu"),
    "id": $("#cpid").val()
  })
  var t = confirm("是否确认删除采集端口");
  if (t == true) {
    $.ajax({
      url: host + "/sp/",
      type: "DELETE",
      data: delDate,
      error: function(resp) {
        console.log(resp.responseText);
        alert(resp.responseText);
      },
      success: function(param) {
        console.log("删除成功");
        location.reload()
      }
    })
  }
  //   $.ajax({
  //       url: host + "/sp/",
  //       type: "DELETE",
  //       data: delDate,
  //       error:function(resp){
  //           console.log(resp.responseText);
  //           alert(resp.responseText);
  //       },
  //       success: function (param) {
  //           console.log("删除成功");
  //           location.reload()
  //       }
  //   })
  // setTimeout(function () {
  //     location.reload()
  // }, 2000)
})
//删除采集单元
$("#delunit").on("click", function() {
  var delDate1 = JSON.stringify({
    "muid": $("#equipUnit").attr("data-mu"),
    "spid": $("#collecterport").attr("data-id"),
    "id": $("#cuid").val()
  })
  // var delDate = JSON.stringify({
  //   muid: $("#equipUnit").attr("data-mu"),
  //   spid: $("#collecterport").attr("data-id"),
  //   id: $("#cuid").val(),
  //   name: $("#cuname").val(),
  //   period: parseInt($("#cuperiod").val()),
  //   timeout: parseInt($("#cutimeout").val()),
  //   maxCommunicationErrors: parseInt($("#cumaxNum").val()),
  //   element: $("#equipLab").val(),
  //   enable: $("#unitswitch").prop("checked"),
  //   setting: portSet
  // });
  var y = confirm("是否删除采集单元");
  if (y == true) {
    $.ajax({
      url: host + "/su/",
      type: "DELETE",
      data: delDate1,
      error: function(resp) {
        console.log(resp.responseText);
        alert(resp.responseText);
      },
      success: function(param) {
        location.reload()
        console.log("删除成功");
        // setTimeout(function () {
        //     location.reload()
        // }, 2000)
      }
    })
  }
})

//表格删除当前行
$(".table").on("click", ".delete-row", function() {
  //当前变量名
  var variableName = $(this).prev().prev().children().val();
  var sure = confirm("是否删除" + variableName + "?删除后将从数据库移除相关记录，不可恢复。");
  if (sure == true) {
    $(this).parent().remove();
  }
})

//表格新增一行
$(".addchangable button").on("click", function() {
  var newrow = '<tr class="addvariable">' +
    '<td><input type="text" value=""/></td>' +
    '<td><input type="text" value=""/></td>' +
    '<td class="delete-row cter">' +
    '<img src="../pic/delete1.png">' +
    '</td>' +
    '</tr>'
  $(".table").append(newrow)
})

//保存
$("#reboot").on("click", function() {
  var u = confirm("是否重启设备？");
  if (u == true) {
    $.ajax({
      url: host + "/mu/reboot",
      type: "PUT",
      error: function(resp) {
        console.log(resp.responseText);
        alert(resp.responseText);
      },
      success: function(param) {
        $("#loading").show();
        setInterval(pingPong, 5000);
      }
    })
  }
})

$("#restart").on("click", function() {
  var u = confirm("是否重启采集程序？");
  if (u == true) {
    $.ajax({
      url: host + "/mu/restart",
      type: "PUT",
      error: function(resp) {
        console.log(resp.responseText);
        alert(resp.responseText);
      },
      success: function(param) {
        $("#loading").show();
        setInterval(pingPong, 5000);
      }
    })
  }
})

$("#upgrade").on("click", function() {
  // 定义升级的操作
  var doUpgrade = function() {
    $.ajax({
      url: host + "/mu/upgrade",
      type: "POST",
      error: function(resp) {
        console.log(resp.responseText);
        alert(resp.responseText);
      },
      success: function(param) {
        $("#loading").show();
        setInterval(pingPong, 5000);
      }
    })
  };

  // 获取版本信息
  $.ajax({
    url: host + "/mu/upgrade",
    type: "GET",
    error: function(resp) {
      console.log(resp.responseText);
      alert(resp.responseText);
    },
    success: function(param) {
      if (param.status != "200") {
        alert(param.message);
        return;
      }
      var u = confirm(param.message);
      if (u == true) {
        doUpgrade();
      }
    }
  })
})

function pingPong() {
  $.ajax({
    url: host + "/ping",
    type: "GET",
    error: function(resp) {
      console.log(resp.responseText);
      alert(resp.responseText);
    },
    success: function(res) {
      location.href = "../pages/login.html";
      $("#loading").hide();
    }
  })
}
