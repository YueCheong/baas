layui.use('table', function(){
    var table = layui.table;
    table.render({
        elem: '#network-table-toolbar'
        ,url:"api/networks"
        ,toolbar: '#network-table-toolbar-toolbarDemo'
        ,title: 'Network Control'
        ,page: true
        ,parseData:function (res) {
            if(res.Package == null) {
                return {
                    "code":0
                    ,"msg":""
                    ,"count":0
                    ,"data":res.Package
                }
            }
            // 这是个引用，修改会影响到原res
            var result = [];
            for(var i=0; i<res.Package.length; i++) {
                var temp = {"ID": i+1, "Name": res.Package[i]};
                result.push(temp);
            }
            return {
                "code":0
                ,"msg":""
                ,"count":result.length
                ,"data":result
            }
        }
        ,cols: [[
            {field:'ID', width:200, title: 'ID', sort: true}
            ,{field:'Name', width:300, title: 'Name'}
            ,{width:215, align:'center', fixed: 'right', toolbar: '#network-table-toolbar-barDemo'}
        ]]
    });
    
    table.on('tool(network-table-toolbar)', function(obj) {
        var data = obj.data;
        if(obj.event == 'del') {
            obj.del(obj);
            api("api/networks/" + data.Name, null, 'DELETE');
        }
    })
    var $ = layui.$, active = {
        addNetwork: function(){ // 添加网络
            layer.open({
                //调整弹框的大小
                area:['900px','600px'],
                shadeClose:false,//点击旁边地方自动关闭
                //动画
                anim:2,
                //弹出层的基本类型
                type: 2,
                //刚才定义的弹窗页面
                content: '/pop/addNetwork',
                end: function (){
                    jump("/network");
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

var host=window.location.host;
function jump(to) {
    window.location.href='http://'+host+to;
}
