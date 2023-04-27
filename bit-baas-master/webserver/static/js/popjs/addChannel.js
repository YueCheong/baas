var host = window.location.host;

function jump(to) {
    window.location.href = 'http://' + host + to;
}

layui.use(['form', 'upload', 'layer'], function () {
    var form = layui.form;

    // 下拉框赋初始值
    layui.$.ajax({
        url: '/api/blockchains',
        dataType: 'json',
        type: 'get',
        success: function (data) {
            console.log(data);
            layui.$.each(data.Package, function (index, item) {
                layui.$('#blockchain-name').append(new Option(item.Name, item.ID));// 下拉菜单里添加元素
            });
            layui.form.render("select"); //重新渲染 固定写法
        }
    });

    //表单取值
    layui.$('#add-channel-sure-btn').on('click', function () {
        var data = form.val('add-channel');
        // data.OrdererOrg = "{" + data.OrdererOrg + "}";
        console.log(JSON.stringify(data));

        if (data.Peers != "") {
            data.Peers = JSON.parse(data.Peers);
        }
        if (data.AnchorPeers != "") {
            data.AnchorPeers = JSON.parse(data.AnchorPeers);
        }
        if(data.AnchorPeers == "") {
            data.AnchorPeers = null;
        }
        data.BlockchainID = Number(data["blockchain-name"]);
        delete data["blockchain-name"];
        console.log("json prase" + JSON.stringify(data));
        api('../api/channels/', JSON.stringify(data), 'put').then(resove => {
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


function api(url, opt, methods) {
    return new Promise(function (resove, reject) {
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
                resove(JSON.parse(xmlHttp.responseText));
            } else {
                resove(xmlHttp.status);
            }
        };
    });
}