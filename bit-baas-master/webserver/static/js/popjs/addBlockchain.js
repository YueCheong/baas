var host=window.location.host;
function jump(to) {
    window.location.href='http://'+host+to;
}

layui.use(['form'], function(){
    var form = layui.form;

    // Network下拉框赋初始值
    layui.$.ajax({
        url: '/api/networks',
        dataType: 'json',
        type: 'get',
        success: function (data) {
            layui.$.each(data.Package, function (index, item) {
                layui.$('#network-name').append(new Option(item, item));// 下拉菜单里添加元素
            });
            layui.form.render("select"); //重新渲染 固定写法
        }
    });

    //表单取值
    layui.$('#add-blockchain-sure-btn').on('click', function(){
        var data = form.val('add-blockchain');
        var put_data = {}

        if(data.OrdererOrg != null && data.OrdererOrg != "" && data.OrdererOrg != undefined) {
            put_data.OrdererOrg = JSON.parse(data.OrdererOrg);
        }
        if(data.PeerOrg != null && data.PeerOrg != "" && data.PeerOrg != undefined) {
            put_data.PeerOrg = JSON.parse(data.PeerOrg);
        }
        if(data["network-name"] != null && data["network-name"] != "" && data["network-name"] != undefined) {
            put_data.Netname = data["network-name"];
        }
        if(data.Name != null && data.Name != "" && data.Name != undefined) {
            put_data.Name = data["Name"];
        }
        api('../api/blockchains/', JSON.stringify(put_data), 'put').then(resove => {
            if(resove > 299) {
                alert("发生错误");
            }
            else {
                var index = parent.layer.index; //获取当前弹层的索引号
                parent.layer.close(index); //关闭当前弹层
                jump("/blockchain");
            }
        });
    });

    layui.$('#add-blockchain-cancel-btn').on('click', function(){
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
        if (methods.toUpperCase() === 'POST') {
            xmlHttp.open('POST', url, true);
            xmlHttp.setRequestHeader('Content-Type', 'application/x-www-form-urlencoded;charset=utf-8');
            xmlHttp.send(postData);
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