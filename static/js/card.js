"use strict"

//not in use, should be used for some buttons but was never really implemented cause idk how to
let last_date = undefined;
//node of the card container
let cardContainer = document.getElementById("card-container");

/**
 * this function will return the games after using a title as query
 * @param {string} date not in use
 * @param {string} title title of the game to use as a query, if undefined the title will be empty string
 * @returns array of the games found 
 */
async function fetchCardData(date, title) {
    if (title === undefined) {
        title = "";
    }
    let url = '/games/search?title=' + title;
    // console.log(date)
    // if (date !== undefined) {
    //     url = `/games/search?upload=<-${date}`
    // }
    const res = await fetch(url, {
        method: "GET"
    });
    const resp = await res.json();
    console.log(resp);
    last_date = resp[resp.length - 1].videoData.uploadDate;
    last_date = last_date.substr(0, last_date.indexOf('T'));
    console.log(last_date);
    return resp;
}

/**
 * given a title the function will get the array of games and create the card with the datas from the game object
 * @param {string} title title to use as a query
 */
async function createCards(title) {
    let card_json = await fetchCardData(last_date, title);
    cardContainer.innerHTML = "";
    for (let i = 0; i < card_json.length; i++) {
        let card = document.createElement('a');
        card.setAttribute('href', '/' + card_json[i].videoData.id);
        let col = document.createElement('div');
        col.setAttribute("class", "col");
        let c = document.createElement('div');
        c.setAttribute("class", "card");
        let img = document.createElement('img');
        img.setAttribute("src", card_json[i].videoData.thumbMaxResUrl);
        let cbody = document.createElement('div');
        cbody.setAttribute("class", "card-body");
        let wonBy = document.createElement('img');
        wonBy.setAttribute("src", "/static/images/synergo-png.png");
        if (card_json[i].wonBy == 0) {
            wonBy.setAttribute("src", "/static/images/redez-png.png");
        }
        let datas = document.createElement('div');
        datas.setAttribute("class", "datas");
        let title = document.createElement("h5");
        title.setAttribute("class", "card-title");
        title.innerHTML = card_json[i].videoData.title;
        let p = document.createElement("p");
        p.setAttribute("class", "card-text");
        datas.appendChild(title);
        datas.appendChild(p);
        cbody.appendChild(wonBy);
        cbody.appendChild(datas);
        c.appendChild(img);
        c.appendChild(cbody);
        col.appendChild(c);
        card.appendChild(col);
        cardContainer.appendChild(card);
    }
}
