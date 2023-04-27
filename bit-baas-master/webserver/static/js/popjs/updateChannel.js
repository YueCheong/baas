layui.use(['form'], function () {
    var form = layui.form;

    layui.$('#update-channel-sure-btn').on('click', function () {
        var data = form.val('update-channel');
        // data.OrdererOrg = "{" + data.OrdererOrg + "}";
        console.log(JSON.stringify(data));

        // dubug data
        data.ChannelName = parent.channelName;
        data.BlockchainID = Number(parent.blockchainID);
        if (data.Operation != "") {
            data.Operation = Number(data.Operation);
        }
        if (data.Args != "") {
            data.Args = JSON.parse(data.Args);
        }
        console.log("update channel pop window" + JSON.stringify(data));
        api('../api/channels/', JSON.stringify(data), 'post').then(resove => {
            if (resove > 299) {
                alert("发生错误");
            } else {
                var index = parent.layer.index; //获取当前弹层的索引号
                parent.layer.close(index); //关闭当前弹层
                jump("/channel");
            }
        });
    });


    layui.$('#add-channel-cancel-btn').on('click', function () {
        var index = parent.layer.index; //获取当前弹层的索引号
        parent.layer.close(index); //关闭当前弹层
    });
});

function api(url,opt,methods) {
    return new Promise(function(resove,reject){
        methods = methods || 'POST';
        var xmlHttp = null;
        if (XMLHttpRequest) {
            xmlHttp = new XMLHttpRequest();
        } else {
            xmlHttp = new ActiveXObject('Microsoft.XMLHTTP');
        };
        var params = [];
        for (var key in opt){
            if(!!opt[key] || opt[key] === 0){
                params.push(key + '=' + opt[key]);
            }
        };
        var postData = params.join(',');
        if (methods.toUpperCase() == 'POST') {
            xmlHttp.open('POST', url, true);
            xmlHttp.setRequestHeader('Content-Type', 'application/x-www-form-urlencoded;charset=utf-8');
            xmlHttp.send(opt);
        }else if (methods.toUpperCase() === 'GET') {
            xmlHttp.open('GET', url + '?' + postData, true);
            xmlHttp.send(null);
        }else if(methods.toUpperCase() === 'DELETE'){
            xmlHttp.open('DELETE', url + '?' + postData, true);
            xmlHttp.send(null);
        }else if(methods.toUpperCase() == 'PUT') {
            xmlHttp.open('PUT', url, true);
            xmlHttp.setRequestHeader('Content-Type', 'application/x-www-form-urlencoded;charset=utf-8');
            xmlHttp.send(opt);
        }
        xmlHttp.onreadystatechange = function () {
            if (xmlHttp.readyState == 4 && xmlHttp.status == 200) {
                resove(JSON.parse(xmlHttp.responseText));
            }
            else {
                resove(xmlHttp.status);
            }
        };
    });
}