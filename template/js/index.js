var source = new EventSource("ReceiveMsg"); //sendMessage后台的访问路径
source.onmessage = function (event) {
    var jsondata = {
        text: event.data,
        color: "#6f9",
        fixed: false,
        shadow: true
    };

    damoo.emit(jsondata);
};