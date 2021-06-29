var videoLink = document.getElementById("videoLink")

videoLink.addEventListener('change', async function () {
    var alreadySent = await fetch("/game/check/" + videoLink.value, {
        method: "GET",
        headers: {
            'Accept': 'application/json'
        }
    })
});