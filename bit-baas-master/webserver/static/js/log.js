layui.use('laydate', function(){
    var startDate = layui.laydate;
    var endDate = layui.laydate;
    //执行一个laydate实例
    startDate.render({
        elem: '#startDate' //指定元素
    });
    endDate.render({
        elem:'#endDate'
    })
});

layui.use('table', function(){
    var table = layui.table;
    //第一个实例
    table.render({
        elem: '#log'
        ,height: 312
        ,url: '/demo/table/user/' //数据接口
        ,page: true //开启分页
        ,cols: [[ //表头
            {field: 'id', title: 'ID', width:80, sort: true, fixed: 'left'}
            ,{field: 'time', title: '操作时间', width:200}
            ,{field: 'type', title: '操作类型', width:200, sort: true}
            ,{field: 'object', title: '操作对象', width:200}
            ,{field: 'result', title: '操作结果', width: 200}
        ]]
    });

});