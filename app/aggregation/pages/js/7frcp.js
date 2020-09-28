$(document).ready(function () {
    $.ajax({
        url: host + "/map/",
        type: "GET",
        dataType: "JSON",
        success: function (param) {
            let list = param.data.agencys;
            $("#serverIP").val(param.data.serverIP);
            $("#serverPort").val(param.data.serverPort);
            for (a in list) {
                $("#list").append(
                    `<li data-num="${parseInt(a) + 1}">
                        <div class="collapsible-header">代理${parseInt(a) + 1}<i class="material-icons right" onclick="deleteOption(this)">delete</i></div>
                        <div class="collapsible-body row" style="margin: 0;">
                            <form class="test">
                                <div class="input-field col l6"><input type="text" name="name" placeholder="名称" value="${list[a].name}"><label class="active">名称</label></div>
                                <div class="input-field col l6"><input type="text" name="localIP" placeholder="本地IP" value="${list[a].localIP}"><label class="active">本地IP</label></div>
                                <div class="input-field col l6"><input type="text" name="localPort" placeholder="本地端口" value="${list[a].localPort}"><label class="active">本地端口</label></div>
                                <div class="input-field col l6"><input type="text" name="remotePort" placeholder="远程端口" value="${list[a].remotePort}"><label class="active">远程端口</label></div>
                            </form>
                        </div>
                    </li>
                    `
                )
            }
            $('.modal').modal();
            Materialize.updateTextFields();
            $('.collapsible').collapsible();
        }
    })
})

let sliceArray = (array, size) => {
    var result = [];
    for (var x = 0; x < Math.ceil(array.length / size); x++) {
        var start = x * size;
        var end = start + size;
        result.push(array.slice(start, end));
    }
    return result;
}

let deleteOption = (cc) => {
    selfModal("是否删除此选项", null, function (ok) {
        if (ok) {
            $(cc).parent().parent().remove();
            $("#selfModal").modal("close");
            Materialize.toast('删除成功!', 3000)
        } else {
            $("#selfModal").modal("close");
        }
    })
}

$("#add").on("click", function () {
    let length = $('#list').find('li').length
    $("#list").append(
        `<li data-num="${parseInt(length) + 1}">
            <div class="collapsible-header">代理${parseInt(length) + 1}<i class="material-icons right" onclick="deleteOption(this)">delete</i></div>
            <div class="collapsible-body row" style="margin: 0;">
                <form class="test">
                    <div class="input-field col l6"><input type="text" name="name" placeholder="名称" value=""><label class="active">名称</label></div>
                    <div class="input-field col l6"><input type="text" name="localIP" placeholder="本地IP" value=""><label class="active">本地IP</label></div>
                    <div class="input-field col l6"><input type="text" name="localPort" placeholder="本地端口" value=""><label class="active">本地端口</label></div>
                    <div class="input-field col l6"><input type="text" name="remotePort" placeholder="远程端口" value=""><label class="active">远程端口</label></div>
                </form>
            </div>
        </li>`
    )
})

$("#restart").on("click", function () {
    $.ajax({
        url: host + "/map/restart",
        type: "PUT",
        success: function (res) {
            if (res.status == 0) {
                alert('重启成功')
            } else {
                alert("配置失败！")
            }
        },
        error: function (res1) {
            console.log("失败");
            alert(res1.message);
        }
    })
})

$("#7frcp").on("click", function () {
    let data = {};
    let dataList = [];
    
    data['serverIP'] = $("#serverIP").val();
    data['serverPort'] = $("#serverPort").val();

    let list = sliceArray($(".test").serializeArray(),4);
    for (item of list) {
        let ff = {}
        for (child of item) {
            ff[child.name] = child.value;
        }
        dataList.push (ff)
    }
    data['agencys'] = dataList;

    data = JSON.stringify(data);

    $.ajax({
        url: host + "/map/",
        data: data,
        type: "POST",
        dataType: "JSON",
        success: function (res) {
            if (res.status == 0) {
                location.reload();
            } else {
                alert("配置失败！")
            }
        },
        error: function (res1) {
            console.log("失败");
            alert(res1.message);
        }
    })
})