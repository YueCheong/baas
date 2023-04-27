var host=window.location.host;
function jump(to) {
    window.location.href='http://'+host+to;
}
layui.use('element', function(){
    var element = layui.element;
});