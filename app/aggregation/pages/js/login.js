function execute() {
    console.log("函数调用");
    new Vue({
        el: '#body-box',
        data() {
            return {
                account: "",
                password: "",
            }
        },
        mounted() {
            
        },
        methods: {
            login() {
                if (this.account == "") {
                    alert("请填写账号");
                    return false;
                }
                if (this.password == "") {
                    alert("请填写密码");
                    return false;
                }
                if (this.account != "" && this.password != "") {
                    let login = JSON.stringify({
                        "username": this.account,
                        "password": this.password,
                    });
                    console.log("login", login)
                    $.ajax({
                        url: `${host}/user/login`,
                        data: login,
                        type: 'POST',
                        dataType: 'json',
                        success: function (res) { //res是请求返回的字段，可打印出来，获取自己所需要的值判断是否登陆成功
                            console.log("res", res)
                            if (res.status == "0") {
                                location.href = "./homePage.html"; //通过验证后
                                console.log("登陆成功");
                            } else {
                                alert("账号密码错了,请重新输入。");
                            }
                        },
                        error: function (re) {
                            console.log("re", re);
                            alert("账号密码错了,请重新输入。");
                        }
                    });
                }
            },
            keyLogin() {
                console.log("keyCode", event.keyCode)
                if (event.keyCode == 13) {
                    this.login()
                }
            }
        }
    })
}
execute()


