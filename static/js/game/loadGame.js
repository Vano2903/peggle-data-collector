"use strict"
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
    }
    const resp = await res.json();
    console.log(resp[0]);
    console.log(status);
    return resp[0];
}

async function genPage(){
    gameData = await fetchGameData();
    console.log(id);
    console.log(gameData);
    document.getElementById("title").innerHTML = gameData.videoData.title;
    if (gameData.wonBy == 1) {
        document.getElementById("scrown").style.display = "block";
    }else if (gameData.wonBy == 0) {
        document.getElementById("rcrown").style.display = "block";
    }else{
        document.getElementById("scrown").style.display = "block";
        document.getElementById("rcrown").style.display = "block";
    }

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

function setOverall(game){
    
}

