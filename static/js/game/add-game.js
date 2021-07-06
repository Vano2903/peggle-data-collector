var videoLink = document.getElementById("videoLink")

function youtube_parser(url) {
    var regExp = /^.*((youtu.be\/)|(v\/)|(\/u\/\w\/)|(embed\/)|(watch\?))\??v?=?([^#&?]*).*/;
    var match = url.match(regExp);
    return (match && match[7].length == 11) ? match[7] : false;
}

$(document).ready(
    function () {
        videoLink.addEventListener('input', async function () {
            var res = await fetch("/games/check/" + youtube_parser(videoLink.value), {
                method: "GET",
                headers: {
                    'Accept': 'application/json'
                }
            })
            var resp = await res.text();
            console.log(JSON.parse(resp))
            // console.log(JSON.parse(alreadySent))
        });
    }
)