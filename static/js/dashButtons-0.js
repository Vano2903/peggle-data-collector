var divs = { stats: document.getElementById("stats"), users: document.getElementById("users-managment"), game: document.getElementById("game"), user: document.getElementById("user") }
var showing = "stats"

function show(what) {
    hideAll()
    divs[what].style.display = "block";
    console.log(what)
    showing = what
}

function hideAll() {
    for (var doc in divs) {
        divs[doc].style.display = "none";
    }
}

function logout() {
    localStorage.removeItem("user");
    window.location.replace('/users/login');
}

async function getPfp() {
    var res = await fetch('/users/pfp', {
        method: "POST",
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(user)
    });
    var resp = await res.text();
    console.log(resp);
    var respJson = JSON.parse(resp);
    console.log(respJson);
    if (respJson.url != "") {
        document.getElementById("pfp").src = respJson.url
    }
}