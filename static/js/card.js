var last_date = undefined

async function fetchCardData(date) {
    url = '/games/search?limit=20'
    console.log(date)
    if (date !== undefined) {
        url = `/games/search?upload=${date}&limit=20`
    }
    const res = await fetch(url, {
        method: "GET"
    });
    const resp = await res.json();
    return resp
}

