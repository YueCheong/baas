layui.use(['table', 'form'], function(){
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
    renderTable(table, "api/channels");
    form.on('select(blockchain-name)', function (data) {
        var data = form.val('channel-display');
        console.log("channel control select" + JSON.stringify(data));
        var url = "api/channels";
        if(data["blockchain-name"] != "") {
            url += "?blockchainid=" + data["blockchain-name"].toString();
        }
        renderTable(table, url);
    });


    table.on('tool(channel-table-toolbar)', function(obj) {
        var data = obj.data;
        console.log("update channel button" + JSON.stringify(data));
        window.blockchainID = data.BlockchainID;
        window.channelName = data.Name;
        if(obj.event == 'update') {
            layer.open({
                //调整弹框的大小
                area:['900px','600px'],
                shadeClose:false,//点击旁边地方自动关闭
                //动画
                anim:2,
                //弹出层的基本类型
                type: 2,
                //刚才定义的弹窗页面
                content: '/pop/updateChannel',
                end: function (){
                    jump("/channel");
                }
            });
        }
    })

    var $ = layui.$, active = {
        addChannel: function(){ // 添加通道
            layer.open({
                //调整弹框的大小
                area:['900px','600px'],
                shadeClose:false,//点击旁边地方自动关闭
                //动画
                anim:2,
                //弹出层的基本类型
                type: 2,
                //刚才定义的弹窗页面
                content: '/pop/addChannel',
                end: function (){
                    jump("/channel");
                }
            });
        }
    };

    $('.demoTable .layui-btn').on('click', function(){
        var type = $(this).data('type');
        active[type] ? active[type].call(this) : '';
    });
});

function renderTable(table, url) {
    table.render({
        elem: '#channel-table-toolbar'
        ,url: url
        ,toolbar: '#channel-table-toolbar-toolbarDemo'
        ,title: 'Channel Control'
        ,page: true
        ,parseData:function (res) {
            console.log(res);
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
            return {
                "code":0
                ,"msg":""
                ,"count":result.length
                ,"data":result
            }
        }
        ,cols: [[
            {field:'Name', width:300, title: 'Name'}
            ,{field:'Peers', width:300, title: 'Peers'}
            ,{field:'AnchorPeers', width:300, title: 'AnchorPeers'}
            ,{field:'BlockchainID', width:300, title: 'BlockchainID'}
            ,{field:'Blockchainname', width:300, title: 'Blockchainname'}
            ,{width:215, align:'center', fixed: 'right', toolbar: '#channel-table-toolbar-barDemo'}
        ]]
    });
}

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
        }
        xmlHttp.onreadystatechange = function () {
            if (xmlHttp.readyState == 4 && xmlHttp.status == 200) {
                resove(JSON.parse(xmlHttp.responseText));
            }
        };
    });
}