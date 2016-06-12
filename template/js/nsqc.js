var BindHtml = {
  init: function(pageid) {
    BindHtml.conNsq("test", "9527", "18321421187");

    if (typeof(EventSource) !== "undefined") {
      var source = new EventSource("ReceiveMsg"); //从服务器内存取
      //var source = new EventSource("GetMsgDB/?topic='test'&sort='-senddate'&limit=10"); //从数据库取
      source.onmessage = function(event) {
        // document.getElementById("result").innerHTML+=event.data + "<br />";
        // return;
        var jsondata = JSON.parse(event.data);

        if (pageid == "home") {
          BindHtml.getMsg(jsondata)
        }
        if (pageid == "index") {
          BindHtml.getBarrageMsg(jsondata)
        }


      };
    } else {
      document.getElementById("result").innerHTML = "对不起，您的浏览器不支持服务器发送的事件…";
    }
  },

  conNsq: function(topic, channel, userid) {

    $.ajax({
      url: "ConMsq/",
      type: 'GET',
      data: {
        consumerid:localStorage["consumerid"],
        topic: topic,
        channel: channel,
        userid: userid
      },
      success: function(data) {
          if(data == "err"){
            console.log("conNsq no:"+data);
            return;
          }else{
            localStorage["consumerid"] = data;
            console.log("conNsq id:"+data);
          }
      },
      error: function(XMLHttpRequest, textStatus, errorThrown) {
        console.log("conNsq err" + XMLHttpRequest.status);
        console.log("conNsq err" + XMLHttpRequest.readyState);
        console.log("conNsq err" + textStatus);
        return;
      },
    });
  },
  disNsq: function() {
    $.ajax({
      url: "DisMsq",
      type: 'GET',
      success: function(data) {
        console.log("disNsq success" + data);
        return;
      },
      error: function(XMLHttpRequest, textStatus, errorThrown) {
        console.log("conNsq err" + XMLHttpRequest.status);
        console.log("conNsq err" + XMLHttpRequest.readyState);
        console.log("conNsq err" + textStatus);
        return;
      },
    });
  },
  sendMsgNsq: function() {

    var msg = $("#iptCon").val();
    $.ajax({
      url: "http://nsq-ttthzygi35.tenxcloud.net:20157/put?topic=test",
      type: 'POST',
      data: msg,
      success: function(data) {},
      error: function(data) {
        alert("消息推送完毕！");
      },
    });
  },
  sendMsgLocal: function(msg) {

    $.ajax({
      url: "/SendMsg/",
      type: 'GET',
      data: {
        sendmsg: msg
      },
      success: function(data) {
        var jsondata = JSON.parse(data);
        console.log(jsondata.message);
      },
      error: function(data) {
        console.log(data);
      },
    });
  },
  getMsgDB: function(topic) {
    $.ajax({
      url: "/GetMsgDB/",
      type: 'GET',
      data: {
        topic: topic,
        sort: "-senddate",
        limit: 10
      },
      success: function(jsondata) {
        console.log(jsondata);
      },
      error: function(data) {
        console.log(data);
      },
    });
  },
  getMsg: function(jsondata) {
      var li = $("#ul_msg li[id='msgid" + jsondata.MssageID + "']").val();
      if (typeof(li) == "undefined") {
        $("#ul_msg").append("<li id='msgid" + jsondata.MssageID + "'>" + jsondata.Mssage + "</li>");
      }
  },
  getBarrageMsg: function(jsondata) {

      var data = {
        id: jsondata.MssageID,
        text: jsondata.Mssage,
        color: "#6f9",
        fixed: false,
        shadow: true
      };

      damoo.emit(data);

  },
  pageOut: function() {
    $.ajax({
      url: "/StopConsumer/",
      type: 'GET',
      data: {
        consumerid: localStorage["consumerid"]
      },
      success: function(data) {
        console.log(data);
      },
      error: function(XMLHttpRequest, textStatus, errorThrown) {
        console.log("pageOut err" + XMLHttpRequest.status);
        console.log("pageOut err" + XMLHttpRequest.readyState);
        console.log("pageOut err" + textStatus);
        return;
      },
    });
  }
}
