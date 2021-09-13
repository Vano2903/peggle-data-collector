//TODO add login with google (maybe github tho)
//TODO jwt implementation

//user object
var user;

/**
 * onload check if the user is stored in the local storage, if so check the login using the local storage datas
 */
$(document).ready(
    function () {
        var localsUser = localStorage.getItem("user")
        if (localsUser !== null) {
            user = JSON.parse(localsUser)
            checkLogin("def")
        }
    }
)

/**
 * get the values from user and password textbox
 * @returns object with username and password
 */
function getLoginData() {
    var user = document.getElementById("user").value;
    var psw = document.getElementById("password").value;
    return { username: user, password: psw }
}

/**
 * display the loading animation and send the user datas to the login api and check if the response is correct
 * @param {string} code if "def" is set as code the function wont get the datas from the textboxes to check the login, 
 * if undefined will get the datas from the textboxes, otherwise will split the code as user informations
 */
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
/**
 * this function will check the response of the login api
 * @param {string} cont content of the fetch (if is a "text/html" the login is succesful and will store the user in localstorage)
 * @param {string} resp response of the api
 */
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

var a = { 'error': 'nessuno si é ancora loggato' }