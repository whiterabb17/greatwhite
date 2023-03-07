let login = {
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
            // Listen
            login.listen();
            // Explore default path
        })
    },
    TryLogin: function(ip, port, usr, pwd, local) {
        let msg = ip + "#" + port + "#" + usr + "#" + pwd + "#" + local
        let message = {"name": "login", "payload": msg};
        console.log(message)
        // Send message
        asticode.loader.show();
        astilectron.sendMessage(message, function(message) {
            // Init
            console.log(message)
            // Check error
            if (message.name === "error") {
                asticode.loader.hide();
                asticode.notifier.error(message.payload);
                return
            }
            asticode.loader.hide()
        })
    },
    listen: function() {
        astilectron.onMessage(function(message) {
            switch (message.name) {                
                case "okay":
                    listen.about("Connected successfully to TeamServer");
                    return {payload: "payload"}
                    break;
                case "invalid":
                    asticode.notifier.error("Invalid Login Details or Teamserver Address");
                    return {payload: "payload"}
                    break;
            }
        });
    },
};