function updateSong(title) {
	var xmlhttp = new XMLHttpRequest();
	xmlhttp.onreadystatechange = function() {
		if (xmlhttp.readyState === 4) {
			console.log(xmlhttp.responseText);
		}
	};
	xmlhttp.open("POST", "http://192.168.1.253:8080/set_song/");
	xmlhttp.send(title)
}

function checkTabs() {
	chrome.tabs.getAllInWindow(null, function(tabs){
	    for (var i = 0; i < tabs.length; i++) {
	    	if (tabs[i].audible) {
	    		updateSong(tabs[i].title)
	    		return
	    	}                        
	    }
	    updateSong("No Song")
	});
}

setInterval(checkTabs, 1000);

chrome.extension.onRequest.addListener(function(data, sender) {
	updateSong(data);
});
