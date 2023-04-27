layui.use('table', function(){
    var table = layui.table;
    table.render({
        elem: '#contract-log-table-toolbar'
        ,url:"/api/contractlog/" + parent.window.ID
        ,title: 'Contract Log'
        ,page: true
        ,parseData:function (res) {
            console.log(res.Package);
            if(res.Package == null) {
                return {
                    "code":0
                    ,"msg":""
                    ,"count":0
                    ,"data":res.Package
                }
            }
            // 这是个引用，修改会影响到原res
            var result = res.Package;
            var result_data = [];
            for(var i=0; i<result.length; i++) {
                var temp_result = {};
                temp_result["RecordID"] = result[i]["RecordID"];
                temp_result["InvokeType"] = result[i]["InvokeType"];
                temp_result["Args"] = JSON.stringify(result[i]["Args"]);
                temp_result["StatusCode"] = result[i]["StatusCode"];
                temp_result["TransactionID"] = result[i]["TransactionID"];
                temp_result["Payload"] = result[i]["Payload"];
                temp_result["Time"] = result[i]["Time"];
                result_data.push(temp_result);
            }
            console.log(result_data);
            return {
                "code":0
                ,"msg":""
                ,"count":result_data.length
                ,"data":result_data
            }
        }
        ,cols: [[
            {field:'RecordID', width:200, title: 'RecordID', sort: true}
            ,{field:'InvokeType', width:300, title: 'InvokeType'}
            ,{field:'Args', width:300, title: 'Args'}
            ,{field:'StatusCode', width:300, title: 'StatusCode'}
            ,{field:'TransactionID', width:300, title: 'TransactionID'}
            ,{field:'Payload', width:300, title: 'Payload'}
            ,{field:'Time', width:300, title: 'Time'}
        ]]
    });

    layui.$('.demoTable .layui-btn').on('click', function(){
        var type = $(this).data('type');
        active[type] ? active[type].call(this) : '';
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
            xmlHttp.send(postData);
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

var host=window.location.host;
function jump(to) {
    window.location.href='http://'+host+to;
}
