var divs = { stats: document.getElementById("stats"), users: document.getElementById("users-managment"), game: document.getElementById("game"), user: document.getElementById("user") }
var showing = "stats"

function show(what) {
    hideAll()
    divs[what].style.display = "block";
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
    var respJson = await res.json();
    console.log(respJson);
    if (respJson.url != "") {
        document.getElementById("pfp").src = respJson.url
        document.getElementById("pfp2").src = respJson.url
    }
}

async function stats() {
    locationSaver('stats');
    await initCalendar()
}

function locationSaver(position) {
    if (position == undefined) {
        if (localStorage.getItem("location") == undefined) {
            localStorage.setItem("location", "stats");
            show("stats");
            return "stats";
        } else {
            showing = localStorage.getItem("location");
            show(showing);
            return showing;

        }
    } else {
        localStorage["location"] = position;
        showing = position;
        show(position);
        return position;
    }
}