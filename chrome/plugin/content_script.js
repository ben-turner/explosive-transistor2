setInterval(function() {
	chrome.extension.sendRequest(document.title.replace(/ - Google Play Music$/, ""));
}, 5000);
