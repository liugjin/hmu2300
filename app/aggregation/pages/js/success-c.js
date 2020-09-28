function execute() {
  new Vue({
    el: "#body-box",
    data() {
      return {
        restart: false,
        upgrade: false,
        cloudupgrade: false,
        cancelTimer: null
      };
    },
    mounted() {},
    methods: {
      //   本地升级
      upgradeClick() {
        const vm = this;
        let u = confirm("是否升级设备？");
        if (u == true) {
          $.ajax({
            url: `${host}/uploadRestartFileForm`,
            type: "POST",
            cache: false,
            data: new FormData($("#uploadFileForm")[0]),
            contentType: false,
            processData: false,
            success: function(data) {
              if (data.code == 0) {
                vm.upgrade = true;
                vm.changeApp();
              }
              if (data.code == 1) {
                alert("此版本为测试版本");
              }
            },
            error: function() {
              alert("异常");
            }
          });
        }
      },
      //   设备重启
      restartClick() {
        console.log("函数执行");
        const vm = this;
        let u = confirm("是否重启设备？");
        if (u == true) {
          vm.restart = true;
          $.ajax({
            url: `${host}/mu/reboot`,
            type: "PUT",
            error: function(resp) {
              alert(resp.responseText);
            },
            success: function(param) {
              console.log("param", param);
              vm.cancelTimer = setInterval(() => vm.pingPong(), 5000);
            }
          });
        }
      },
      // 云升级
      urlSend() {
        const vm = this;
        $.ajax({
          url: host + "/updateUrlForm",
          type: "POST",
          data: $("#cloudupdateUrlForm").serialize(),
          success: function(data) {
            if (data.code == 0) {
              alert("升级中且升级包大小为：" + data.size);
              vm.cloudupgrade = true;
              vm.changeApp();
            }
            if (data.code == -1) {
              alert("云升级失败，错误原因：" + data.err);
              document.location.reload(); //重新加载当前页面
            }
          },
          error: function() {
            alert("异常");
          }
        });
      },
      // 升级
      changeApp() {
        const vm = this;
        $.ajax({
          url: host + "/changeApp",
          type: "POST",
          success: function(data) {
            console.log("changeApp", data);
            if (data.code == 0) {
              vm.cloudupgrade = false;
              vm.cancelTimer = setInterval(() => vm.pingPong(), 5000);
            }
          },
          error: function(data) {
            console.log("异常", data);
            vm.cancelTimer = setInterval(() => vm.pingPong(), 5000);
          }
        });
      },
      pingPong() {
        const vm = this;
        $.ajax({
          url: host + "/ping",
          type: "GET",
          error: function(resp) {
            console.log(resp.responseText);
            console.log("ping");
          },
          success: function(res) {
            console.log("ping", res);
            clearInterval(vm.cancelTimer);
            location.href = "../pages/login.html";
            //  $("#loading").hide();
          }
        });
      }
    }
  });
}
execute();
