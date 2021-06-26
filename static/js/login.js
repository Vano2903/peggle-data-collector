//TODO add login with google (maybe github tho)
//TODO jwt implementation

var user

function getLoginData() {
    var user = document.getElementById("user").value;
    var psw = document.getElementById("password").value;
    return { username: user, password: psw }
}

async function checkLogin(code) {
    document.getElementById("loader-wrapper").style.display = "block";

    if (code == undefined) {
        user = getLoginData()
    } else {
        var ele = code.split(";");
        user = { username: ele[0], password: ele[1] };
    }
    console.log(user);
    console.log(JSON.stringify(user))
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
function checkResponse(cont, resp, user) {
    errore = document.getElementById("errore");
    if (!cont.includes("text/html")) {
        const respJson = JSON.parse(resp);
        console.log(respJson.message);
        errore.innerHTML = respJson.message;
        errore.style.display = "block";
    } else {
        // localStorage.setItem("user", user)
        // tok = suppCode == undefined ? supp : suppCode;
        document.write(resp);
        document.close();
    }
}