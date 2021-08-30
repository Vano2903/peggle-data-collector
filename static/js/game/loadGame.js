"use strict"
document.getElementById("loader-wrapper").style.display = "block";
window.onload = async function(){
    await genPage()
    document.getElementById("loader-wrapper").style.display = "none";
}

let id = window.location.pathname.slice(1);
let gameData;
document.getElementsByTagName('title')[0].innerHTML += id;

async function fetchGameData(){
    let url = `/games/search?id=${id}`;
    const res = await fetch(url, {
        method: "GET"
    });
    const status = res.status;
    if (status === 404){
        // window.location = "/404"
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

async function getUserProfilePicture(name){
    let url = "/users/pfp/"+ name;
    const res = await fetch(url, {
        method: "GET"
    });
    const status = res.status;
    if (status === 400){
        return ""
    }
    const resp = await res.json();
    return resp.url;
}

async function genPage(){
    gameData = await fetchGameData();
    document.getElementById("title").innerHTML = gameData.videoData.title;
    if (gameData.wonBy == 1) {
        document.getElementById("scrown").style.display = "block";
    }else if (gameData.wonBy == 0) {
        document.getElementById("rcrown").style.display = "block";
    }else{
        document.getElementById("scrown").style.display = "block";
        document.getElementById("rcrown").style.display = "block";
    }
    for (let i = 0; i < 4; i++){
        setSection(i, gameData)
    }
    document.getElementById("youtube-embed").src = `https://www.youtube.com/embed/${id}?rel=0`
    document.getElementById("userName").innerHTML += gameData.addedBy;
    document.getElementById("userPfp").src = await getUserProfilePicture(gameData.addedBy);
}

function hideAllSections(){
    for(let i = 0; i < 4; i++){
        document.getElementById(`par${i}`).style.display = "none";
        document.getElementById(`labButPar${i}`).classList.remove("btn-outline-secondary");
        document.getElementById(`labButPar${i}`).classList.remove("btn-outline-success");
        document.getElementById(`labButPar${i}`).classList.add("btn-outline-secondary");
    }
}

function showSection(section){
    hideAllSections()
    document.getElementById(`par${section}`).style.display = "block";
    document.getElementById(`labButPar${section}`).classList.remove("btn-outline-secondary");
    document.getElementById(`labButPar${section}`).classList.add("btn-outline-success");
}

function setSection(index, game){
    if (index === 0){
        document.getElementById("stpoints").innerHTML += game.stats.synergo.overall.tPoints;
        document.getElementById("st25").innerHTML += game.stats.synergo.overall.t25;
        document.getElementById("rtpoints").innerHTML += game.stats.redez.overall.tPoints;
        document.getElementById("rt25").innerHTML += game.stats.redez.overall.t25;
    }else{
        document.getElementById(`s${index}points`).innerHTML += game.stats.synergo["g"+index].points;
        document.getElementById(`s${index}25`).innerHTML += game.stats.synergo["g"+index].n25;
        document.getElementById(`s${index}fe`).innerHTML += game.stats.synergo["g"+index].valFe;
        document.getElementById(`s${index}charName`).innerHTML += game.stats.synergo["g"+index].character;
        document.getElementById(`s${index}charImg`).src = "/static/images/" + game.stats.synergo["g"+index].character.slice(0, 3)+".png";

        document.getElementById(`r${index}points`).innerHTML += game.stats.redez["g"+index].points;
        document.getElementById(`r${index}25`).innerHTML += game.stats.redez["g"+index].n25;
        document.getElementById(`r${index}fe`).innerHTML += game.stats.redez["g"+index].valFe;
        document.getElementById(`r${index}charName`).innerHTML += game.stats.redez["g"+index].character;
        document.getElementById(`r${index}charImg`).src = "/static/images/" + game.stats.redez["g"+index].character.slice(0, 3)+".png";
    }
}

