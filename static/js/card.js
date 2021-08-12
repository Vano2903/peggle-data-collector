let last_date = undefined;
let cardContainer = document.getElementById("card-container"); 

async function fetchCardData(date) {
    url = '/games/search?limit=5'
    console.log(date)
    if (date !== undefined) {
        url = `/games/search?upload=<-${date}&limit=5`
    }
    const res = await fetch(url, {
        method: "GET"
    });
    const resp = await res.json();
    console.log(resp)
    last_date = resp[resp.length - 1].videoData.uploadDate;
    last_date = last_date.substr(0, last_date.indexOf('T'))
    console.log(last_date)
    return resp
}

async function createCards(){
    card_json = await fetchCardData(last_date)
    for(let i = 0; i < card_json.length; i++){
        var card = document.createElement('a');
        card.setAttribute('href', '/'+card_json[i].videoData.id);
        var col = document.createElement('div')
        col.setAttribute("class", "col")
        var c = document.createElement('div')
        c.setAttribute("class", "card")
        var img = document.createElement('img')
        img.setAttribute("src", card_json[i].videoData.thumbMaxResUrl)
        var cbody = document.createElement('div')
        cbody.setAttribute("class", "card-body")
        var wonBy = document.createElement('img')
        wonBy.setAttribute("src", "/static/images/synergo-png.png")
        if (card_json[i].wonBy == 0){ 
            wonBy.setAttribute("src", "/static/images/redez-png.png")
        }
        var datas = document.createElement('div')
        datas.setAttribute("class", "datas");
        var title = document.createElement("h5")
        title.setAttribute("class", "card-title");
        title.innerHTML = card_json[i].videoData.title
        var p = document.createElement("p")
        p.setAttribute("class", "card-text")
        datas.appendChild(title)
        datas.appendChild(p)
        cbody.appendChild(wonBy)
        cbody.appendChild(datas)
        c.appendChild(img)
        c.appendChild(cbody)
        col.appendChild(c)
        card.appendChild(col)
        cardContainer.appendChild(card)
    }
}
