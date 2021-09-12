"use strict"

let last_date = undefined;
let cardContainer = document.getElementById("card-container");

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
