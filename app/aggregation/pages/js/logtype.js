 var tokencheck = {token : sessionStorage.getItem("token")};
//验证登陆状态
$(document).ready(function(){
    $.ajax({
        url: '/cgi-bin/checktoken.cgi',
        data: tokencheck,
        type: 'GET',
        dataType: 'json',
        success: function(res) {
            if(res.login == false){
                //验证token与后台是否相同，若相同则为登陆了否则返回首页重新登陆
                location.href = "login.html";
                return false;
            }
        }
    });    
});