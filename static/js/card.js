var last_date = undefined

async function fetchCardData(date) {
    url = '/games/search?limit=20'
    console.log(date)
    if (date !== undefined) {
        console.log("nya")
        // url = `/games/search?upload=${date}&limit=20`
    }
    const res = await fetch(url, {
        method: "GET"
    });
    const resp = await res.json();
    console.log(resp)
    return resp
}

