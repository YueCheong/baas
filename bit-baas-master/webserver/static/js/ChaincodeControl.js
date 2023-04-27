layui.use(['table', 'form'], function () {
    var table = layui.table;
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
    renderTable(table, "api/contracts/");
    form.on('select(blockchain-name)', function (data) {
        var data = form.val('contract-display');
        var url = "api/contracts/";
        if (data["blockchain-name"] != "") {
            url += "?blockchainid=" + data["blockchain-name"].toString();
            layui.$.ajax({
                url: '/api/channels?blockchainid=' + data["blockchain-name"].toString(),
                dataType: 'json',
                type: 'get',
                success: function (data) {
                    console.log(data);
                    layui.$.each(data.Package, function (index, item) {
                        layui.$('#channel-name').append(new Option(item.Name, item.Name));// 下拉菜单里添加元素
                    });
                    layui.form.render("select"); //重新渲染 固定写法
                }
            });
        }
        renderTable(table, url);

    });
    form.on('select(channel-name)', function (data) {
        var data = form.val('contract-display');
        var url = "api/contracts/";
        if (data["blockchain-name"] != "") {
            url += "?blockchainid=" + data["blockchain-name"].toString();
        }
        if (data["channel-name"] != "") {
            url += "&channelname=" + data["channel-name"].toString();
        }
        renderTable(table, url);
    });

    table.on('tool(contract-table-toolbar)', function (obj) {
        var data = obj.data;
        console.log("update contract button" + JSON.stringify(data));
        window.ID = data.ID;
        window.BlockchainID = data.BlockchainID;
        window.ContractName = data.ContractName;
        window.ContractName = data.ContractName;
        window.ContractVersion = data.ContractVersion
        if (obj.event == 'invoke') {
            layer.open({
                //调整弹框的大小
                area: ['900px', '600px'],
                shadeClose: false,//点击旁边地方自动关闭
                //动画
                anim: 2,
                //弹出层的基本类型
                type: 2,
                //刚才定义的弹窗页面
                content: '/pop/invokeContract',
                end: function () {
                    jump("/chaincode");
                }
            });
        }
        if (obj.event == 'instantiate') {
            layer.open({
                //调整弹框的大小
                area: ['900px', '600px'],
                shadeClose: false,//点击旁边地方自动关闭
                //动画
                anim: 2,
                //弹出层的基本类型
                type: 2,
                //刚才定义的弹窗页面
                content: '/pop/instantiateContract',
                end: function () {
                    jump("/chaincode");
                }
            });
        }
        if (obj.event == "log") {
            layer.open({
                //调整弹框的大小
                area: ['900px', '600px'],
                shadeClose: false,//点击旁边地方自动关闭
                //动画
                anim: 2,
                //弹出层的基本类型
                type: 2,
                //刚才定义的弹窗页面
                content: '/pop/contractLog',
                end: function () {
                    jump("/chaincode");
                }
            });
        }
        // 点击表格channelName单元格显示channel信息
        // if(obj.event == 'getChannel') {
        //     var data = obj.data;
        //     alert("channel弹窗1");
        //     if(obj.event === 'getChannel') {
        //         alert("channel弹窗2" + data.channelName);
        //     }
        // }
    })

    var $ = layui.$, active = {
        addContract: function () { // 添加合约
            layer.open({
                //调整弹框的大小
                area: ['900px', '600px'],
                shadeClose: false,//点击旁边地方自动关闭
                //动画
                anim: 2,
                //弹出层的基本类型
                type: 2,
                //刚才定义的弹窗页面
                content: '/pop/addContract',
                end: function () {
                    jump("/chaincode");
                }
            });
        }
    };

    $('.demoTable .layui-btn').on('click', function () {
        var type = $(this).data('type');
        active[type] ? active[type].call(this) : '';
    });
});

function renderTable(table, url) {
    table.render({
        elem: '#contract-table-toolbar'
        , url: url
        , toolbar: '#contract-table-toolbar-toolbarDemo'
        , title: 'Contract Control'
        , page: true
        , parseData: function (res) {
            console.log(res);
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
            return {
                "code": 0
                , "msg": ""
                , "count": result.length
                , "data": result
            }
        }
        , cols: [[
            {field: 'ID', width: 200, title: 'ID'}
            , {field: 'ContractName', width: 200, title: 'ContractName'}
            , {field: 'ContractDesc', width: 200, title: 'ContractDesc'}
            , {field: 'ContractVersion', width: 200, title: 'ContractVersion'}
            , {field: 'ContractLang', width: 200, title: 'ContractLang'}
            , {field: 'ChannelName', width: 200, title: 'ChannelName', event: 'getChannel'}
            , {width: 300, align: 'center', fixed: 'right', toolbar: '#contract-table-toolbar-barDemo'}
        ]]
    });
}

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