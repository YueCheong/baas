var host = window.location.host;

function jump(to) {
    window.location.href = 'http://' + host + to;
}

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
    // channel下拉框副初始值
    form.on('select(network-name)', function (data) {
        var data = form.val('add-contract');
        layui.$.ajax({
            url: '/api/channels?blockchainid=' + data["network-name"].toString(),
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
    });
    // upload
    upload.render({
        elem:'#upload-button',
        url:'/api/newcontract',
        accept:'file',
        bindAction: '#add-contract-sure-btn',
        auto: false, //一定要写这个,不能会自动触发上传
        //传入token值,否则会报401错误
        // header:{token:'token值'},
        before: function (){
            var form = layui.form.val('add-contract');
            this.data = {
                BlockchainID: Number(form["network-name"]),
                ContractLang: Number(form["CCLanguage"]),
                ContractName: form["contract-name"],
                ChannelName: form["channel-name"],
                ContractDesc: form["contract-desc"],
                ContractVersion: form["contract-version"]
            }
        },
        done: function (result) {
            //获取到文件导入后的响应数据
            console.log(result);
            var index = parent.layer.index; //获取当前弹层的索引号
            parent.layer.close(index); //关闭当前弹层
        }
    });
    layui.$('#add-contract-cancel-btn').on('click', function () {
        var index = parent.layer.index; //获取当前弹层的索引号
        parent.layer.close(index); //关闭当前弹层
    });
});