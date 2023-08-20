function setTheme(event) {
	var element = event.target;
	var navbar = document.getElementById('navbar');
	navbar.setAttribute('data-bs-theme',element.getAttribute('data-bs-theme-value'));
}

themeIdList = ['light', 'dark', 'auto'];
for (id in themeIdList) {
	document.getElementById('bd-theme-'+themeIdList[id]).addEventListener('click', setTheme);
}

document.querySelectorAll('.playlist-button').forEach(function(button) {
	button.addEventListener('click', onPlaylistButtonClick);
});

function setActiveClass(el) {
	if (!el.classList.contains("active")){
		  el.classList.add('active');
	}
}

function onPlaylistButtonClick(event) {
	var buttonArea = event.target.parentElement;
	document.querySelectorAll('.playlist-button').forEach(function(button) {
		button.classList.remove('active');
	});
	setActiveClass(buttonArea);
}

window.addEventListener("load", (event) => {
	if (window.location.href.indexOf("playlist")===-1) {
		return;
	}
	var playlistIdFromUrl = window.location.href.split("/")[4];
	var buttonArea = document.getElementById('playlist_'+playlistIdFromUrl);
	setActiveClass(buttonArea);
});