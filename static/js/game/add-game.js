var videoLink = document.getElementById("videoLink")

function youtube_parser(url) {
    var regExp = /^.*((youtu.be\/)|(v\/)|(\/u\/\w\/)|(embed\/)|(watch\?))\??v?=?([^#&?]*).*/;
    var match = url.match(regExp);
    return (match && match[7].length == 11) ? match[7] : false;
}

videoLink.addEventListener('change', async function () {
    var alreadySent = await fetch("/game/check/" + youtube_parser(videoLink.value), {
        method: "GET",
        headers: {
            'Accept': 'application/json'
        }
    })
});