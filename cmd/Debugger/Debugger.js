
const bodyParser = require('body-parser'); // Middleware
const rio = require("socket.io-client");
var express=require("express");
var core=express();
var http=require('http');
var httpServer=http.createServer(core);
var io=require('socket.io', { maxHttpBufferSize: 1e8 })(httpServer);

function fb64(encData){
	'use strict'
	let buff = new Buffer.from(encData, "base64");
	let text = buff.toString('ascii');
	return text;
}
function tb64(decData){
	'use strict'	
	let buff = new Buffer.from(decData);
	let base64data = buff.toString('base64');
	return base64data;
}


io.on('connection',function(socket){
	console.log("############################################################################")
	console.log("############################################################################")
	console.log("#########################                        ###########################")
	console.log("#########################    Start Socket Data   ###########################")
	console.log("#########################                        ###########################")
	console.log("############################################################################")
	console.log("############################################################################")
    console.log(socket)
	console.log("############################################################################")
	console.log("############################################################################")
	console.log("#########################                        ###########################")
	console.log("#########################    End Socket Data     ###########################")
	console.log("#########################                        ###########################")
	console.log("############################################################################")
	console.log("############################################################################")
    socket.on("debug", function(data){
        console.log(data)
    })
    socket.on("repl", function(data){
		console.log("REPLY MESSAGE")
        console.log(data)
		console.log("\n\n")
    })
    socket.on("register", function(data){
		console.log("REGISTER MESSAGE")
        console.log(data)
		console.log("\n\n")
    })
    socket.on("reg", function(data){
		console.log("REG MESSAGE")
        console.log(data)
		console.log(fb64(data.res))
		console.log(fb64(data.tag))
		console.log("\n\n")
    })
    socket.on("cinfo", function(data){
		console.log("CLIENT INFO")
        console.log(data)
		console.log("\n\n")
    })
    socket.on("error", function(data){
		console.log("ERROR MESSAGE")
        console.log(data)
		console.log("\n\n")
    })
})

console.log(httpServer.listen(55555, "0.0.0.0"))