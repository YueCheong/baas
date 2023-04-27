layui.use(['form'], function () {
    var form = layui.form;
    var display = form.val('index-display');


    // ******
    // var httpRequest = new XMLHttpRequest();//第一步：建立所需的对象
    // httpRequest.open('GET', "/api/summary/", true);//第二步：打开连接  将请求参数写在url中  ps:"./Ptest.php?name=test&nameone=testone"
    // httpRequest.send();//第三步：发送请求  将请求参数写在URL中
    // /**
    //  * 获取数据后的处理程序
    //  */
    // httpRequest.onreadystatechange = function () {
    //     if (httpRequest.readyState == 4 && httpRequest.status == 200) {
    //         var data = httpRequest.responseText;//获取到json字符串，还需解析
    //         data = JSON.parse(data);
    //         data = data.Package
    //         layui.$("#Name").val("asdfasdf");
    //     }
    // };
    // ***

    layui.$.ajax({
        type: 'get',
        url: '/api/summary/',
        dataType: "json",
        success: function(result) {
            layui.$("#network-num").val("网络数暂无");
            layui.$("#channel-num").val(result.Package.ChannelTotal);
            layui.$("#chaincode-num").val(result.Package.ContractTotal);
            layui.$("#blockchain-num").val(result.Package.BlochainTotal);
        }
    });
});


function api(url, opt, methods) {
    methods = methods || 'POST';
    var xmlHttp = null;
    if (XMLHttpRequest) {
        xmlHttp = new XMLHttpRequest();
    } else {
        xmlHttp = new ActiveXObject('Microsoft.XMLHTTP');
    }
    ;
    var params = [];
    for (var key in opt) {
        if (!!opt[key] || opt[key] === 0) {
            params.push(key + '=' + opt[key]);
        }
    }
    ;
    var postData = params.join(',');
    if (methods.toUpperCase() === 'POST') {
        xmlHttp.open('POST', url, true);
        xmlHttp.setRequestHeader('Content-Type', 'application/x-www-form-urlencoded;charset=utf-8');
        xmlHttp.send(postData);
    } else if (methods.toUpperCase() === 'GET') {
        xmlHttp.open('GET', url + '?' + postData, true);
        xmlHttp.send(null);
    } else if (methods.toUpperCase() === 'DELETE') {
        xmlHttp.open('DELETE', url + '?' + postData, true);
        xmlHttp.send(null);
    } else if (methods.toUpperCase() == 'PUT') {
        xmlHttp.open('PUT', url, true);
        xmlHttp.setRequestHeader('Content-Type', 'application/x-www-form-urlencoded;charset=utf-8');
        xmlHttp.send(opt);
    }

    xmlHttp.onreadystatechange = function () {
        if (xmlHttp.readyState == 4 && xmlHttp.status == 200) {
            console.log("api", xmlHttp.responseText);
            return xmlHttp.responseText;
        } else {
            return xmlHttp.status;
        }
    };
}
