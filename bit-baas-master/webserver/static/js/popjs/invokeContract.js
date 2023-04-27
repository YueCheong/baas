var host = window.location.host;

function jump(to) {
    window.location.href = 'http://' + host + to;
}

layui.use(['form', 'upload', 'layer'], function () {
    var form = layui.form;
    var upload = layui.upload;
    var flag = true;


    // form.on('select(invoke-type)', function (data) {
    //     var data = form.val('invoke-contract');
    //     layui.form.render(); //重新渲染 固定写法
    //     // flag = false;
    //     var str =
    //         '<div class="layui-form-item">' +
    //         '<label class="layui-form-label">parameter</label>' +
    //         '<div class="layui-input-inline">' +
    //         '<input type="text" name="parameter" class="layui-input input-double-width">' +
    //         '</div>' +
    //         '</div>';
    //     layui.$("#invoke-form").append(str);
    // });

    var data = form.val('invoke-contract');
    layui.form.render(); //重新渲染 固定写法
    // flag = false;
    var str =
        '<div class="layui-form-item">' +
        '<label class="layui-form-label">parameter</label>' +
        '<div class="layui-input-inline">' +
        '<input type="text" name="parameter" class="layui-input input-double-width">' +
        '</div>' +
        '</div>';
    layui.$("#invoke-form").append(str);


    layui.$('#invoke-contract-sure-btn').on('click', function () {
        var data = form.val('invoke-contract');
        var post_data = new Map();
        post_data["ID"] = parent.ID;
        post_data["BlockchainID"] = Number(parent.BlockchainID);
        post_data["ContractName"] = parent.ContractName;
        post_data["ContractVersion"] = parent.ContractVersion;
        post_data["InvokeType"] = Number(data["invoke-type"]);
        if(post_data["InvokeType"] == 0) {
            post_data["Args"] = data["parameter"];
        }
        else if(post_data["InvokeType"] == 1) {
            post_data["Args"] = data["parameter"];
        }
        console.log(JSON.stringify(post_data));
        api('../api/contractcall/', JSON.stringify(post_data), 'post').then(resove => {
            if (resove > 299) {
                alert("发生错误");
            } else {
                alert("transaction id: \n" + resove.Package.Txid + "\n" + "amount: \n" + resove.Package.Payload);
                var index = parent.layer.index; //获取当前弹层的索引号
                parent.layer.close(index); //关闭当前弹层
                jump("/chaincode");
            }
        });
    });

    layui.$('#invoke-contract-cancel-btn').on('click', function () {
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
        var postData = params.join('&');
        if (methods.toUpperCase() === 'POST') {
            xmlHttp.open('POST', url, true);
            xmlHttp.setRequestHeader('Content-Type', 'application/x-www-form-urlencoded;charset=utf-8');
            xmlHttp.send(opt);
        } else if (methods.toUpperCase() === 'GET') {
            xmlHttp.open('GET', url + '?' + postData, true);
            xmlHttp.send(null);
        } else if (methods.toUpperCase() === 'DELETE') {
            xmlHttp.open('DELETE', url + '?' + postData, true);
            xmlHttp.send(null);
        }
        xmlHttp.onreadystatechange = function () {
            if (xmlHttp.readyState == 4 && xmlHttp.status == 200) {
                console.log("api function xmltp", xmlHttp.responseText);
                resove(JSON.parse(xmlHttp.responseText));
            }
        };
    });
}