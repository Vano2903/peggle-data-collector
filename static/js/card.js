async function fetchCardData(date) {
    url = '/games/search?title=&limit=20'
    console.log(date)
    if (date === undefined) {
        
    }
    const res = await fetch(url, {
        method: "GET"
    });
    const resp = await res.json();
    return resp
}

