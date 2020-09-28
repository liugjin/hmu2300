function execute() {
    new Vue({
        el: '#body-box',
        data() {
            return {
                account: "",
                oldpwd: "",
                newpwd: ""
            }
        },
        mounted() {

        },
        methods: {
            submit() {
                let repswsub = JSON.stringify({
                    "username": this.account,
                    "oldpassword": this.oldpwd,
                    "newpassword":this.newpwd
                });
                if(this.oldpwd  == this.newpwd){
                    alert("新的密码不能和旧的密码一样！")
                    return
                }
                $.ajax({
                    url: `${host}/user/passwd`,
                    data: repswsub,
                    type: "PUT",
                    dataType: "json",
                    success: function (res) {
                        if(res.message =="ok"){
                            alert("密码修改成功，请到登录页面重新登录")
                            location.href = "login.html";
                        }
                        // if (res.data.msg == xxx) {//旧密码不正确！
                        //     alert("旧密码不正确！")
                        // } else if (res.data.msg == cccc) {                  //新的密码不能和旧的密码一样！
                        //     alert("新的密码不能和旧的密码一样！")
                        // } else {                                  //密码修改成功
                        //     alert("密码修改成功，请到登录页面重新登录")
                        //     location.href = "login.html";
                        // }
                    }
                });
            }
        },
    })
}
execute()