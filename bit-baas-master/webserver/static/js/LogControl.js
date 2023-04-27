layui.use('table', function () {
    var table = layui.table;
    table.render({
        elem: '#log-table-toolbar'
        , url: "api/contractlog"
        ,toolbar: '#log-table-toolbar-toolbarDemo'
        , title: 'Log Control'
        , page: true
        , parseData: function (res) {
            console.log(res.Package);
            if (res.Package == null) {
                return {
                    "code": 0
                    , "msg": ""
                    , "count": 0
                    , "data": res.Package
                }
            }
            // 这是个引用，修改会影响到原res
            var result = res.Package;
            for (var i = 0; i < result.length; i++) {
                result["Args"] = JSON.stringify(result["Args"]);
            }
            return {
                "code": 0
                , "msg": ""
                , "count": result.length
                , "data": result
            }
        }
        , cols: [[
            {field: 'RecordID', width: 120, title: 'RecordID', sort: true, hide: true}
            , {field: 'BlockchainName', width: 120, title: 'BlockchainName'}
            , {field: 'ChannelName', width: 120, title: 'ChannelName'}
            , {field: 'ContractName', width: 120, title: 'ContractName'}
            , {field: 'ContractDesc', width: 120, title: 'ContractDesc', hide: true}
            , {field: 'ContractVer', width: 120, title: 'ContractVer', hide: true}
            , {field: 'InvokeType', width: 120, title: 'InvokeType'}
            , {field: 'Args', width: 120, title: 'Args', hide: true}
            , {field: 'StatusCode', width: 120, title: 'StatusCode'}
            , {field: 'TransactionID', width: 120, title: 'TransactionID'}
            , {field: 'Payload', width: 120, title: 'Payload', hide: true}
            , {field: 'Time', width: 120, title: 'Time', sort: true}
        ]]
    });


    $('.demoTable .layui-btn').on('click', function () {
        var type = $(this).data('type');
        active[type] ? active[type].call(this) : '';
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
        var params = [];
        for (var key in opt) {
            if (!!opt[key] || opt[key] === 0) {
                params.push(key + '=' + opt[key]);
            }
        }
        var postData = params.join('&');
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
            xmlHttp.send(postData);
        }
        xmlHttp.onreadystatechange = function () {
            if (xmlHttp.readyState == 4 && xmlHttp.status == 200) {
                resove(JSON.parse(xmlHttp.responseText));
            }
        };
    });
}

var host = window.location.host;

function jump(to) {
    window.location.href = 'http://' + host + to;
}
