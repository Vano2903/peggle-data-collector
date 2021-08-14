"use strict"
let id = window.location.pathname.slice(1);
let gameData;
document.getElementsByTagName('title')[0].innerHTML += id;

async function fetchGameData(){
    url = `/games/search?id=${id}`;
    const res = await fetch(url, {
        method: "GET"
    });
    const resp = await res.json();
    console.log(resp[0]);
    return resp[0];
}

async function genPage(){
    gameData = await fetchGameData();
    console.log(id);
    console.log(gameData);
}