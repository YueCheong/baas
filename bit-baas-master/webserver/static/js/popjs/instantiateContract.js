var host=window.location.host;
function jump(to) {
    window.location.href='http://'+host+to;
}

layui.use(['form'], function(){
    var form = layui.form;

    //表单取值
    layui.$('#instantiate-contract-sure-btn').on('click', function(){
        var data = form.val('instantiate-contract');
        // data.OrdererOrg = "{" + data.OrdererOrg + "}";
        var post_data = new Map();
        post_data['ID'] = parent.ID;
        post_data['BlockchainID'] = parent.BlockchainID;
        post_data['ContractName'] = parent.ContractName;
        post_data['ContractVersion'] = parent.ContractVersion;
        post_data['Operation'] = 1;
        post_data['Args'] = data['parameter'];
        api('../api/contracts/', JSON.stringify(post_data), 'post').then(resove => {
            if(resove > 299) {
                alert("发生错误");
            }
            else {
                alert("实例化成功");
                var index = parent.layer.index; //获取当前弹层的索引号
                parent.layer.close(index); //关闭当前弹层
                jump("/chaincode");
            }
        });
    });

    layui.$('#instantiate-contract-cancel-btn').on('click', function(){
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