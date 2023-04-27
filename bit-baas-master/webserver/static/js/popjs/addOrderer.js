layui.use(['form', 'upload', 'layer'], function () {
    var form = layui.form;
    var upload = layui.upload;

    // blockchain下拉框赋初始值
    layui.$.ajax({
        url: '/api/blockchains',
        dataType: 'json',
        type: 'get',
        success: function (data) {
            console.log(data);
            layui.$.each(data.Package, function (index, item) {
                layui.$('#network-name').append(new Option(item.Name, item.ID));// 下拉菜单里添加元素
            });
            layui.form.render("select"); //重新渲染 固定写法
        }
    });

    //表单取值
    layui.$('#add-orderer-sure-btn').on('click', function(){
        var data = form.val('add-orderer');
        var post_data = {};
        post_data.Operation = Number(data.Operation);
        data = data["orderer-config"];
        data = JSON.parse(data);

        // for(var key in data) {
        //     if(notEmpty(key)) {
        //         post_data[key] = JSON.parse(data[key]);
        //     }
        // }
        if(notEmpty(data["Name"])) {
            post_data.Name = data["Name"];
        }
        if(notEmpty(data["Domain"])) {
            post_data.Domain = data["Domain"];
        }
        if(notEmpty(data["MSPID"])) {
            post_data.MSPID = data["MSPID"];
        }
        if(notEmpty(data["Nodes"])) {
            post_data.Nodes = data["Nodes"];
        }
        post_data.BlockchainID = window.parent.ID;
        api('../api/blockchains', JSON.stringify(post_data), 'post').then(resove => {
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

});

function notEmpty(s) {
    if(s != null && s != "" && s != undefined) {
        return true;
    }
    return false;
}

function api(url,opt,methods) {
    return new Promise(function(resove,reject){
        methods = methods || 'POST';
        var xmlHttp = null;
        if (XMLHttpRequest) {
            xmlHttp = new XMLHttpRequest();
        } else {
            xmlHttp = new ActiveXObject('Microsoft.XMLHTTP');
        }
        var params = [];
        for (var key in opt){
            if(!!opt[key] || opt[key] === 0){
                params.push(key + '=' + opt[key]);
            }
        }
        var postData = params.join('&');
        if (methods.toUpperCase() === 'POST') {
            xmlHttp.open('POST', url, true);
            xmlHttp.setRequestHeader('Content-Type', 'application/x-www-form-urlencoded;charset=utf-8');
            xmlHttp.send(opt);
        }else if (methods.toUpperCase() === 'GET') {
            xmlHttp.open('GET', url + '?' + postData, true);
            xmlHttp.send(null);
        }else if(methods.toUpperCase() === 'DELETE'){
            xmlHttp.open('DELETE', url + '?' + postData, true);
            xmlHttp.send(null);
        }else if (methods.toUpperCase() == 'PUT') {
            xmlHttp.open('PUT', url, true);
            xmlHttp.setRequestHeader('Content-Type', 'application/x-www-form-urlencoded;charset=utf-8');
            xmlHttp.send(postData);
        }
        xmlHttp.onreadystatechange = function () {
            if (xmlHttp.readyState == 4 && xmlHttp.status == 200) {
                resove(JSON.parse(xmlHttp.responseText));
            }
        };
    });
}
