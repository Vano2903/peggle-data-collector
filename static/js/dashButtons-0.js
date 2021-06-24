var divs = { stats: document.getElementById("stats"), users: document.getElementById("users-managment"), game: document.getElementById("game"), user: document.getElementById("user") }

function show(what) {
    hideAll()
    divs[what].style.display = "block";
}

function hideAll() {
    for (var doc in divs) {
        divs[doc].style.display = "none";
    }
}