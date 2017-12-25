function connect(wsUri, recvFun) {
	ws = new MyWebSocket(wsUri)
	ws.SetOpen(function (evt) {
		ws.SendPileUpData()
	})
	ws.SetRecv(recvFun)
	return ws
}

function EncodeUint32(number, len=0) {
	var b = new Uint8Array(4+len)
	b[0] = number & 0xFF
	b[1] = (number >> 8) & 0xFF
	b[2] = (number >> 16) & 0xFF
	b[3] = (number >> 24) & 0xFF
	return b
}

function DecodeUint32(buff) {
	try{
		let num = Number(buff[3]) << 24
		num += Number(buff[2]) << 16
		num += Number(buff[1]) << 8
		num += Number(buff[0])
		return num
	}catch(err) {
		console.log(err)
		return 0
	}
}

function checkWebSocket() {
	if(window.WebSocket){
	    console.log('This browser supports WebSocket');
	}else{
	    alert('This browser does not supports WebSocket');
	}
}

function MyWebSocket(url) {
	this.websocket = new window.WebSocket(url)
	this.sendDatas = new Array()

	this.SetOpen = function (setFun) {
		this.websocket.onopen = setFun
	}
	this.SetRecv = function (setFun) {
		this.websocket.onmessage = function(evt) {
			// console.log(evt)
			var reader = new FileReader();
			reader.onload = function(e) {
				// console.log(e, setFun)
			    setFun(new Uint8Array(e.target.result))
			}; 
			reader.readAsArrayBuffer(evt.data);
		}	
	}
	this.SetClose = function (setFun) {
		this.websocket.onclose = setFun
	}
	this.websocket.onerror = function (evt) {
		console.log(
			"Error------------->\n",
			showAll(evt),
			"<------------------------\n"
			)
	}
	this.Close = function () {
		this.websocket.close()
	}
	this.SendMsg = function (buff) {
		if ( this.websocket.readyState == 1 ) {
			this.websocket.send(buff)
			// console.log(buff)
		} else {
			this.sendDatas.push(buff)
		}
	}
	// 发送堆积数据
	this.SendPileUpData = function () {
		if ( this.websocket.readyState != 1 ) {
			return
		}
		for (var i = 0; i < this.sendDatas.length; i++) {
			this.SendMsg(this.sendDatas[i])
			// console.log(i,this.sendDatas[i])
		}
		console.log("沉郁数据发送完毕！")
	}
	return 
}

function showAll(arg) {
	var str = ""
		for(var i in arg){
			str += i+" : "+arg[i] + "\n"
		}
	return str
}