let buildHead = `
    <meta charset="UTF-8">
    <link rel="stylesheet" href="static/css/base.css"/>
    <link rel="stylesheet" href="static/lib/astiloader/astiloader.css">
    <link rel="stylesheet" href="static/lib/astimodaler/astimodaler.css">
    <link rel="stylesheet" href="static/lib/astinotifier/astinotifier.css">
    <link rel="stylesheet" href="static/lib/font-awesome-4.7.0/css/font-awesome.min.css">
`

let logs = {
    notify: function(log) {
        let logList = document.getElementById("logs")
        let c = document.createElement("div");
        let span = document.createElement("span")
        let b = document.createElement("b")
        b.innerHTML = "Log"
        span.append(b)
        let br = document.createElement("br")
        let br2 = document.createElement("br")
        span.append(br)
        span.append(br2)
        let p = document.createElement("a")
        p.innerHTML = log
        span.append(p)
        c.append(span)
        logList.append(c)
    },
    init: function() {
        // Init
        asticode.loader.init();
        asticode.modaler.init();
        asticode.notifier.init();
        
        // Wait for astilectron to be ready
        document.addEventListener('astilectron-ready', function() {
            logs.listen();
        })
    },
    clearLogs: function() {     
        asticode.loader.show();   
        let c = document.createElement("div");
        c.innerHTML = html;
        asticode.modaler.setContent(c);
        asticode.loader.hide()
        asticode.modaler.show();
        document.getElementById("logs").html()
    },
    listen: function() {
        astilectron.onMessage(function(message) {
            switch (message.name) {
                case "addLog":
                    logs.notify(message.payload);
                    return {payload: "payload"};
                    break;
            }
        });
    },
};