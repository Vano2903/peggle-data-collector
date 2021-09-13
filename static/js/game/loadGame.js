"use strict"

/**
 * this function runs onload and show the loader, generate the page, hide the loader and run the countup animation
 */
document.getElementById("loader-wrapper").style.display = "block";
window.onload = async function () {
    await genPage()
    document.getElementById("loader-wrapper").style.display = "none";
    runAnimations(".countup0")
}

//get the id from the url
let id = window.location.pathname.slice(1);
//game object
let gameData;
//set the title of the page to the id of the game
document.getElementsByTagName('title')[0].innerHTML += id;

/**
 * this function will fetch the game by id and return the json,
 * if the id is not found the 404 page will be loaded
 * @returns object of the game found
 */
async function fetchGameData() {
    let url = `/games/search?id=${id}`;
    const res = await fetch(url, {
        method: "GET"
    });
    const status = res.status;
    if (status === 404) {
        const res = await fetch('/404', {
            method: "GET"
        });
        const resp = await res.text();
        document.write(resp);
        document.close();
    }
    const resp = await res.json();
    console.log(resp[0]);
    return resp[0];
}

/**
 * this function will get the url of the user who added this game
 * @param {string} name name of the user who added the game
 * @returns the url of the user's profile picture
 */
async function getUserProfilePicture(name) {
    let url = "/users/pfp/" + name;
    const res = await fetch(url, {
        method: "GET"
    });
    const status = res.status;
    if (status === 400) {
        return ""
    }
    const resp = await res.json();
    console.log(resp[0])
    return resp[0];
}

/**
 * this function will set all the information on the html
 */
async function genPage() {
    gameData = await fetchGameData();
    document.getElementById("title").innerHTML = gameData.videoData.title;
    if (gameData.wonBy == 1) {
        document.getElementById("scrown").style.display = "block";
    } else if (gameData.wonBy == 0) {
        document.getElementById("rcrown").style.display = "block";
    } else {
        document.getElementById("scrown").style.display = "block";
        document.getElementById("rcrown").style.display = "block";
    }
    for (let i = 0; i < 4; i++) {
        setSection(i, gameData)
    }
    document.getElementById("youtube-embed").src = `https://www.youtube.com/embed/${id}?rel=0`
    document.getElementById("userName").innerHTML += gameData.addedBy;
    document.getElementById("userPfp").src = await getUserProfilePicture(gameData.addedBy);
}

/**
 * hide all the section of the game page
 */
function hideAllSections() {
    for (let i = 0; i < 4; i++) {
        document.getElementById(`par${i}`).style.display = "none";
        document.getElementById(`labButPar${i}`).classList.remove("btn-outline-secondary");
        document.getElementById(`labButPar${i}`).classList.remove("btn-outline-success");
        document.getElementById(`labButPar${i}`).classList.add("btn-outline-secondary");
    }
}

/**
 * show just the section choosen
 * @param {string} section name of the section
 */
function showSection(section) {
    hideAllSections()
    document.getElementById(`par${section}`).style.display = "block";
    document.getElementById(`labButPar${section}`).classList.remove("btn-outline-secondary");
    document.getElementById(`labButPar${section}`).classList.add("btn-outline-success");
}

/**
 * set the "overall" section of the game and the game parts 
 * @param {number} index index of the game (if 0 means overall, otherwise it's just a normal section)
 * @param {object} game game object
 */
function setSection(index, game) {
    if (index === 0) {
        document.getElementById("stpoints").innerHTML += game.stats.synergo.overall.tPoints;
        document.getElementById("st25").innerHTML += game.stats.synergo.overall.t25;
        document.getElementById("rtpoints").innerHTML += game.stats.redez.overall.tPoints;
        document.getElementById("rt25").innerHTML += game.stats.redez.overall.t25;
    } else {
        document.getElementById(`s${index}points`).innerHTML += game.stats.synergo["g" + index].points;
        document.getElementById(`s${index}25`).innerHTML += game.stats.synergo["g" + index].n25;
        document.getElementById(`s${index}fe`).innerHTML += game.stats.synergo["g" + index].valFe;
        document.getElementById(`s${index}charName`).innerHTML += game.stats.synergo["g" + index].character;
        document.getElementById(`s${index}charImg`).src = "/static/images/" + game.stats.synergo["g" + index].character.slice(0, 4) + ".png";

        document.getElementById(`r${index}points`).innerHTML += game.stats.redez["g" + index].points;
        document.getElementById(`r${index}25`).innerHTML += game.stats.redez["g" + index].n25;
        document.getElementById(`r${index}fe`).innerHTML += game.stats.redez["g" + index].valFe;
        document.getElementById(`r${index}charName`).innerHTML += game.stats.redez["g" + index].character;
        document.getElementById(`r${index}charImg`).src = "/static/images/" + game.stats.redez["g" + index].character.slice(0, 4) + ".png";
    }
}