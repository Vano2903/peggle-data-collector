//object that has all the getElementById of the main page (of users)
var divs = { stats: document.getElementById("stats"), users: document.getElementById("users-managment"), game: document.getElementById("game"), user: document.getElementById("user") }
//store what section the website is showing
var showing = "stats"

/**
 * show a section of the webiste (display show) and refresh the "showing" variable
 * @param {string} what the name of the section to show
 */
function show(what) {
    hideAll()
    divs[what].style.display = "block";
    showing = what
}

/**
 * set to display none all the section
 */
function hideAll() {
    for (var doc in divs) {
        divs[doc].style.display = "none";
    }
}

/**
 * remove from the local storage the user's data
 * and redirect the user to the login page
 */
function logout() {
    localStorage.removeItem("user");
    window.location.replace('/users/login');
}

/**
 * this function will fetch the profile picture of the user (will use the user saved in the local storage to get the information)
 * and will render the image in the icon on the top right and the dedicate user section
 */
async function getPfp() {
    console.log('/users/pfp/' + user.username);
    var res = await fetch('/users/pfp/' + user.username, {
        method: "GET"
    });
    var respJson = await res.json();
    console.log(respJson);
    document.getElementById("pfp").src = respJson[0]
    document.getElementById("pfp2").src = respJson[0]
}

/**
 * this function is assinged to a onclick event and will save as the last page saw the stats page
 * and will run the initCalendar function
 */
async function stats() {
    locationSaver('stats');
    await initCalendar()
}

/**
 * this function will either return the last page saw or save in local storage the last page
 * @param {string} position if not undefined the function will save the string in the local storage, will be used when the user logs in
 * @returns if position is undefined the function will return the last element in the local storage (and show that page to the user using the show function)
 */
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
