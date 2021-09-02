//TODO add login with google (maybe github tho)
//TODO jwt implementation

var user

$(document).ready(
    function () {
        var localsUser = localStorage.getItem("user")
        if (localsUser !== null) {
            user = JSON.parse(localsUser)
            checkLogin("def")
        }
    }
)


function getLoginData() {
    var user = document.getElementById("user").value;
    var psw = document.getElementById("password").value;
    return { username: user, password: psw }
}

async function checkLogin(code) {
    document.getElementById("loader-wrapper").style.display = "block";

    if (code != "def") {
        if (code == undefined) {
            user = getLoginData()
        } else {
            var ele = code.split(";");
            user = { username: ele[0], password: ele[1] };
        }
    }
    const res = await fetch('/users/login', {
        method: "POST",
        headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(user)
    });
    const contentType = res.headers.get("Content-Type");
    const resp = await res.text();
    document.getElementById("loader-wrapper").style.display = "none";
    checkResponse(contentType, resp);
}

//mini easteregg :D oggi é il 9/6/21 xD
function checkResponse(cont, resp) {
    errore = document.getElementById("errore");
    if (!cont.includes("text/html")) {
        const respJson = JSON.parse(resp);
        errore.innerHTML = respJson.message;
        errore.style.display = "block";
    } else {
        errore.style.display = "none";
        // tok = suppCode == undefined ? supp : suppCode;
        document.write(resp);
        localStorage.setItem("user", JSON.stringify(user))
        document.close();
    }
}

var a = {'error': 'nessuno si é ancora loggato'}