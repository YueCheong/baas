var host=window.location.host;
function jump(to) {
    window.location.href='http://'+host+to;
}
layui.use('element', function(){
    var element = layui.element;
    element.on('nav(lNav)', function(elem){
        switch (elem[0].id){
            case "Index":
                jump("/");
                break;
            case "NetworkControl":
                jump("/network");
                break;
            case "BlockchainControl":
                jump("/blockchain");
                break;
            case "ChannelControl":
                jump("/channel");
                break;
            case "ChaincodeControl":
                jump("/chaincode");
                break;
            case "UserControl":
                jump("/user");
                break;
            case "Log":
                jump("/log");
                break;
        }
    });
    element.on('nav(userNav)',function (elem) {
        console.log(elem);
        switch (elem[0].id){
            case "UserProfile":
                jump("/profile");
                break;
            case "Message":
                jump("/message");
                break;
            case "LogOut":
                break;
        }
    })
    element.on('nav(mainNav)',function (elem) {
        switch (elem[0].id){
            case "Console":
                jump("/");
                break;
        }
    })
});