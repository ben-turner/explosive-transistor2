function connect() {
    port = chrome.extension.connectNative('com.ben_turner.explosive_transistor');
}

chrome.extension.onRequest.addListener(function(data, sender) {
	var xmlhttp = new XMLHttpRequest();
	xmlhttp.onreadystatechange = function() {
		if (xmlhttp.readyState === 4) {
			console.log(xmlhttp.responseText);
		}
	};
	xmlhttp.open("POST", "http://192.168.1.253:8080/set_song/");
	xmlhttp.send(data)
});
