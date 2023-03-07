let buildHead = `
    <meta charset="UTF-8">
    <link rel="stylesheet" href="static/css/base.css"/>
    <link rel="stylesheet" href="static/lib/astiloader/astiloader.css">
    <link rel="stylesheet" href="static/lib/astimodaler/astimodaler.css">
    <link rel="stylesheet" href="static/lib/astinotifier/astinotifier.css">
    <link rel="stylesheet" href="static/lib/font-awesome-4.7.0/css/font-awesome.min.css">
`

let listenerList = [0, ""]
let listener = {
    /*
    setStatus: function() {
        let s = document.getElementById("serverstatus").innerHTML
        if (s == "Online") {

        }
    }
*/
    about: function(html) {
        let c = document.createElement("div");
        c.innerHTML = html;
        asticode.modaler.setContent(c);
        asticode.modaler.show();
    },
    init: function() {
        // Init
        asticode.loader.init();
        asticode.modaler.init();
        asticode.notifier.init();
        
        // Wait for astilectron to be ready
        document.addEventListener('astilectron-ready', function() {
            astilectron.sendMessage("netmap", function(message) {
                if (message != "") {
                    console.log(message)
                    addReserved(message)
                }
            });
            // This will listen to messages sent by GO
            astilectron.onMessage(function(message) {
                // Process message
                switch (message) {
                    case "started":
                        listener.about("Port " + message.payload + " has been started");
                        return {payload: "payload"};
                        break;
                    case "stopped":
                        listener.about("Port " + message.payload + " has been stopped");
                        document.getElementById("port" + message.payload + "status").innerHTML = "Reserved";
                        return {payload: "payload"};
                        break;
                }
            });
            listener.listen()
        })
    },
    addReserved: function(existing) {
        let pairs = existing.split("|")
        let paircount = pairs.length
        pairs.forEach(element => {
            if (element != "") {
                pair = element.split("@")
                var table = document.getElementById("porttable");
                var row = document.createElement("tr");
                var tdn = document.createElement("td");
                tdn.innerHTML = pair[0]
                var tdp = document.createElement("td");
                tdp.innerHTML = pair[1]
                row.append(tdn)
                row.append(tdp)        
                table.append(row);
            }            
        });
    },
    addPort: function(n, p) {
        //listener.sendCommand("startPort", p)
        var table = document.getElementById("porttable");
        var row = document.createElement("tr");
        var tdn = document.createElement("td");
        tdn.innerHTML = n
        var tdp = document.createElement("td");
        tdp.innerHTML = p
        row.append(tdn)
        row.append(tdp)        
        table.append(row);
       // var pstat = document.createElement("td")
       // pstat.innerHTML = "Reserved"
       // var opts = document.createElement("td")
      //  opts.append(startBtn)
       // opts.append(stopBtn)
//        row.append(opts)
//        row.innerHTML = "<td>" + p + "</td><td id=\"port" + p + "status\">Reserved</td><td><button onclick=\"listener.sendCommand(\"startPort\", " + p +")\" id=\"start" + p + "btn\"><img src=\"imgs/play.png\" alt=\"Start\" width=\"30\" height=\"30\" id=\"start" + p + "\"></button>&emsp;<button onclick=\"listener.sendCommand(\"stopPort\", " + p +")\" id=\"stop" + p + "btn\"><img src=\"imgs/stop.png\" alt=\"Stop\" width=\"30\" height=\"30\" id=\"stop" + p +"\"></button></td>"

        console.log("Name: " + n + " Port: " + p)
        let v = n + "@" + p
        console.log(v)
    },
    startListener: function() {
        let n = document.getElementById("name").value
        let p = document.getElementById("port").value
        //listener.sendCommand("startPort", p)
        document.getElementById("port").value = ""
        document.getElementById("name").value = ""
        var table = document.getElementById("porttable");
        var row = document.createElement("tr");
        var tdn = document.createElement("td");
        tdn.innerHTML = n
        var tdp = document.createElement("td");
        tdp.innerHTML = p
        row.append(tdn)
        row.append(tdp)        
        table.append(row);
       // var pstat = document.createElement("td")
       // pstat.innerHTML = "Reserved"
       // var opts = document.createElement("td")
      //  opts.append(startBtn)
       // opts.append(stopBtn)
//        row.append(opts)
//        row.innerHTML = "<td>" + p + "</td><td id=\"port" + p + "status\">Reserved</td><td><button onclick=\"listener.sendCommand(\"startPort\", " + p +")\" id=\"start" + p + "btn\"><img src=\"imgs/play.png\" alt=\"Start\" width=\"30\" height=\"30\" id=\"start" + p + "\"></button>&emsp;<button onclick=\"listener.sendCommand(\"stopPort\", " + p +")\" id=\"stop" + p + "btn\"><img src=\"imgs/stop.png\" alt=\"Stop\" width=\"30\" height=\"30\" id=\"stop" + p +"\"></button></td>"

        console.log("Name: " + n + " Port: " + p)
        let v = n + "@" + p
        console.log(v)
        listener.sendCommand("startport", v)
    },    
    sendCommand: function(cmd, args) {   
        let msg = {"name": cmd, "payload": args}
        asticode.notifier.info(cmd);    
        asticode.notifier.info(args);    
        //let message = {"name": "startport", "payload": args};
        asticode.loader.show();
        console.log(msg)
        
        asticode.loader.show();
        astilectron.sendMessage(msg, function(message) {
            // Init
            asticode.loader.hide();

        })
        
    },
    close: function() {
        asticode.loader.show();
        astilectron.sendMessage({name: "close", payload: ""}, function(message) {
            // Init
            asticode.loader.hide();

        })
    },
    listen: function() {
        astilectron.onMessage(function(message) {
            switch (message.name) {
                case "started":
                    listener.about("Port " + message.payload + " has been started");
                    return {payload: "payload"};
                    break;
                case "alert":
                    asticode.notifier.info(message.payload);   
                    return {payload: "payload"};
                    break;
                case "ports":
                    let ports = message.payload
                    let data = ports.split("@")
                    addPort(data[0], data[1]) 
                    return {payload: "payload"};
                    break;
                case "stopped":
                    listener.about("Port " + message.payload + " has been stopped");
                    document.getElementById("port" + message.payload + "status").innerHTML = "Reserved";
                    return {payload: "payload"};
                    break;
                case "portmap":
                    let portVars = message.payload
                    console.log(portVars)
                    addReserved(portVars.name, portVars.port)

                    return {payload: "payload"}
                case "netmap":
                    let netVars = message.payload
                    console.log(netVars)
                    addReserved(netVars.name, netVars.port)

                    return {payload: "payload"}
            }
        });
    },
};