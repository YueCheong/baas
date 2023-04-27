layui.use('table', function(){
    var table = layui.table;
    table.render({
        elem: '#blockchain-table-toolbar'
        ,url:"api/blockchains"
        ,toolbar: '#blockchain-table-toolbar-toolbarDemo'
        ,title: 'Blockchain Control'
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
            for(var i=0; i<result.length; i++) {
                result[i].OrdererOrg = result[i].OrdererOrg.Name;
                var temp_peer = "";
                for(var j=0; result[i].PeerOrg!=null && j<result[i].PeerOrg.length; j++) {
                    temp_peer += result[i].PeerOrg[j].Name + " ";
                }
                result[i].PeerOrg = temp_peer;
                var temp_channel = "";
                for(var k=0; result[i].Channels != null && k<result[i].Channels.length; k++) {
                    temp_channel += result[i].Channels[k].Name;
                }
                result[i].Channels = temp_channel;
            }
            return {
                "code":0
                ,"msg":""
                ,"count":result.length
                ,"data":result
            }
        }
        ,cols: [[
            {field:'ID', width:100, title: 'ID', sort: true}
            ,{field:'Name', width:200, title: 'Name', sort: true}
            ,{field:'OrdererOrg', width:200, title: 'OrdererOrg', sort: true}
            ,{field:'PeerOrg', width:200, title: 'PeerOrg', sort: true}
            ,{field:'Config', width:200, title: 'Config', sort: true}
            ,{field:'Channels', width:200, title: 'Channels', sort: true}
            ,{field:'Status', width:200, title: 'Status', sort: true}
            ,{field:'Netname', width:200, title: 'Netname', sort: true}
            ,{width:500, align:'center', fixed: 'right', toolbar: '#blockchain-table-toolbar-barDemo'}
        ]]
    });
    
    table.on('tool(blockchain-table-toolbar)', function(obj) {
        var data = obj.data;
        if(obj.event == 'del') {
            obj.del(obj);
            api("api/blockchains/" + data.ID, null, 'DELETE');
        }
        if(obj.event == 'init') {
            var temp_data = {};
            temp_data["BlockchainID"] = data.ID;
            temp_data["Operation"] = 3;
            api("api/blockchains/", JSON.stringify(temp_data), 'post').then(resove => {
                if (resove > 299) {
                    alert("发生错误");
                } else {
                    jump("/blockchain");
                }
            })
        }
        if(obj.event == 'start') {
            var temp_data = {};
            temp_data["Operation"] = 1;
            temp_data["BlockchainID"] = data.ID;
            api("api/blockchains/", JSON.stringify(temp_data), 'post').then(resove => {
                if (resove > 299) {
                    alert("发生错误");
                } else {
                    jump("/blockchain");
                }
            });
        }
        if(obj.event == 'stop') {
            var temp_data = {};
            temp_data["Operation"] = 2;
            temp_data["BlockchainID"] = data.ID;
            api("api/blockchains/", JSON.stringify(temp_data), 'post').then(resove => {
                if (resove > 299) {
                    alert("发生错误");
                } else {
                    jump("/blockchain");
                }
            });
        }
        if(obj.event == 'addOrderer') {
            window.ID = data.ID;
            layer.open({
                //调整弹框的大小
                area:['900px','600px'],
                shadeClose:false,//点击旁边地方自动关闭
                //动画
                anim:2,
                //弹出层的基本类型
                type: 2,
                //刚才定义的弹窗页面
                content: '/pop/addOrderer',
                end: function (){
                    jump("/blockchain");
                }
            });
        }
        if(obj.event == 'addPeer') {
            window.ID = data.ID;
            layer.open({
                //调整弹框的大小
                area:['900px','600px'],
                shadeClose:false,//点击旁边地方自动关闭
                //动画
                anim:2,
                //弹出层的基本类型
                type: 2,
                //刚才定义的弹窗页面
                content: '/pop/addPeer',
                end: function (){
                    jump("/blockchain");
                }
            });
        }
    })
    var $ = layui.$, active = {
        addBlockchain: function(){ // 添加网络
            layer.open({
                //调整弹框的大小
                area:['900px','600px'],
                shadeClose:false,//点击旁边地方自动关闭
                //动画
                anim:2,
                //弹出层的基本类型
                type: 2,
                //刚才定义的弹窗页面
                content: '/pop/addBlockchain',
                end: function (){
                    jump("/blockchain");
                }
            });
        }
    };

    $('.demoTable .layui-btn').on('click', function(){
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

var host=window.location.host;
function jump(to) {
    window.location.href='http://'+host+to;
}
