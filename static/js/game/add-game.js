var videoLink = document.getElementById("videoLink")

var gameElements = [
    {
        par: document.getElementById("par1"),
        syn: {
            punt: document.getElementById("SpointsPart1"),
            n25: document.getElementById("S25part1"),
            per: document.getElementById("ScharDropDownPart1"),
            valFe: document.getElementById("SFEDropDownPart1")
        },
        red: {
            punt: document.getElementById("RpointsPart1"),
            n25: document.getElementById("R25part1"),
            per: document.getElementById("RcharDropDownPart1"),
            valFe: document.getElementById("RFEDropDownPart1")
        }
    },
    {
        par: document.getElementById("par2"),
        syn: {
            punt: document.getElementById("SpointsPart2"),
            n25: document.getElementById("S25part2"),
            per: document.getElementById("ScharDropDownPart2"),
            valFe: document.getElementById("SFEDropDownPart2")
        },
        red: {
            punt: document.getElementById("RpointsPart2"),
            n25: document.getElementById("R25part2"),
            per: document.getElementById("RcharDropDownPart2"),
            valFe: document.getElementById("RFEDropDownPart2")
        }
    },
    {
        par: document.getElementById("par3"),
        syn: {
            punt: document.getElementById("SpointsPart3"),
            n25: document.getElementById("S25part3"),
            per: document.getElementById("ScharDropDownPart3"),
            valFe: document.getElementById("SFEDropDownPart3")
        },
        red: {
            punt: document.getElementById("RpointsPart3"),
            n25: document.getElementById("R25part3"),
            per: document.getElementById("RcharDropDownPart3"),
            valFe: document.getElementById("RFEDropDownPart3")
        }
    }
]

function youtube_parser(url) {
    var regExp = /^.*((youtu.be\/)|(v\/)|(\/u\/\w\/)|(embed\/)|(watch\?))\??v?=?([^#&?]*).*/;
    var match = url.match(regExp);
    return (match && match[7].length == 11) ? match[7] : url;
}

function hideAllGameSections() {
    gameElements.forEach((item) => {
        item.par.style.display = "none";
    })
}

function showGameSections(index) {
    hideAllGameSections();
    gameElements[index].par.style.display = "block";
}

function checkIfAllComplete() {
    if (videoLink.value == "") {
        return false
    }
    for (var i = 0; i < gameElements.length; i++) {
        if (gameElements[i].syn.punt.value == "") {
            return false
        }
        if (gameElements[i].syn.n25.value == "") {
            return false
        }
        if (gameElements[i].syn.per.selectedIndex == -1) {
            return false
        }
        if (gameElements[i].syn.valFe.selectedIndex == -1) {
            return false
        }

        if (gameElements[i].red.punt.value == "") {
            return false
        }
        if (gameElements[i].red.n25.value == "") {
            return false
        }
        if (gameElements[i].red.per.selectedIndex == -1) {
            return false
        }
        if (gameElements[i].red.valFe.selectedIndex == -1) {
            return false
        }
    }
    return true
}

function getWhoWon() {
    let [syn, red] = [0, 0];
    for (var i = 0; i < gameElements.length; i++) {
        syn += parseInt(gameElements[i].syn.punt.value)
        red += parseInt(gameElements[i].red.punt.value)
    }
    if (syn > red) {
        return 1;
    } else if (syn < red) {
        return 0;
    } else {
        return -1
    }
}

function getOverall(player) {
    let overall = { tPoints: 0, t25: 0 };
    if (player == "syn") {
        for (var i = 0; i < gameElements.length; i++) {
            overall.tPoints += parseInt(gameElements[i].syn.punt.value)
            overall.t25 += parseInt(gameElements[i].syn.n25.value)
        }
    } else if (player == "red") {
        for (var i = 0; i < gameElements.length; i++) {
            overall.tPoints += parseInt(gameElements[i].red.punt.value)
            overall.t25 += parseInt(gameElements[i].red.n25.value)
        }
    }
    return overall;
}

async function addCommit() {
    const res = await fetch('/commit/add', {
        method: "POST",
        headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(user)
    });
    const resp = await res.text();
    const status = res.status;
    if (status != 200) {
        console.error("error sending the game: ", resp)
    } else {
        console.log("commit added correctly")
    }
}

function genGameData() {
    return {
        videoData: {
            id: youtube_parser(videoLink.value)
        },
        wonBy: getWhoWon(),
        comment: "ciao", //document.getElementById("comment").value
        addedBy: user.username,
        stats: {
            synergo: {
                overall: getOverall("syn"),
                g1: {
                    points: parseInt(gameElements[0].syn.punt.value),
                    n25: parseInt(gameElements[0].syn.n25.value),
                    valFe: parseInt(gameElements[0].syn.valFe.value),
                    character: gameElements[0].syn.per.value
                },
                g2: {
                    points: parseInt(gameElements[1].syn.punt.value),
                    n25: parseInt(gameElements[1].syn.n25.value),
                    valFe: parseInt(gameElements[1].syn.valFe.value),
                    character: gameElements[1].syn.per.value
                },
                g3: {
                    points: parseInt(gameElements[2].syn.punt.value),
                    n25: parseInt(gameElements[2].syn.n25.value),
                    valFe: parseInt(gameElements[2].syn.valFe.value),
                    character: gameElements[2].syn.per.value
                },
            },
            redez: {
                overall: getOverall("red"),
                g1: {
                    points: parseInt(gameElements[0].red.punt.value),
                    n25: parseInt(gameElements[0].red.n25.value),
                    valFe: parseInt(gameElements[0].red.valFe.value),
                    character: gameElements[0].red.per.value
                },
                g2: {
                    points: parseInt(gameElements[1].red.punt.value),
                    n25: parseInt(gameElements[1].red.n25.value),
                    valFe: parseInt(gameElements[1].red.valFe.value),
                    character: gameElements[1].red.per.value
                },
                g3: {
                    points: parseInt(gameElements[2].red.punt.value),
                    n25: parseInt(gameElements[2].red.n25.value),
                    valFe: parseInt(gameElements[2].red.valFe.value),
                    character: gameElements[2].red.per.value
                },
            }
        }
    }
}

async function uploadUpdateGame() {
    document.getElementById("loader-wrapper").style.display = "block";
    await Promise.all([
        uploadGame(),
        addCommit()
    ])
    document.getElementById("loader-wrapper").style.display = "none";
}

async function uploadGame() {
    if (checkIfAllComplete) {
        let game;
        game = genGameData()
        const res = await fetch('/games/add', {
            method: "POST",
            headers: {
                'Accept': 'application/json',
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(game)
        });
        const resp = await res.json();
        const status = res.status;
        if (status != 200) {
            console.error("error sending the game: ", resp)
        } else {
            console.log("%cgame sended correctly :D thanks", "color:green")
        }
    }
    return false;
}

function updataGame() {

}

$("#wholeForm").on("input", function () {
    document.getElementById("send_data").disabled = !checkIfAllComplete()
});

function fillGameSections(gameObject) {
    gameElements[0].syn.punt.value = gameObject.stats.synergo.g1.points
    gameElements[0].syn.n25.value = gameObject.stats.synergo.g1.n25
    gameElements[0].syn.per.value = gameObject.stats.synergo.g1.character
    gameElements[0].syn.valFe.value = gameObject.stats.synergo.g1.valFe

    gameElements[0].red.punt.value = gameObject.stats.redez.g1.points
    gameElements[0].red.n25.value = gameObject.stats.redez.g1.n25
    gameElements[0].red.per.value = gameObject.stats.redez.g1.character
    gameElements[0].red.valFe.value = gameObject.stats.redez.g1.valFe

    gameElements[1].syn.punt.value = gameObject.stats.synergo.g2.points
    gameElements[1].syn.n25.value = gameObject.stats.synergo.g2.n25
    gameElements[1].syn.per.value = gameObject.stats.synergo.g2.character
    gameElements[1].syn.valFe.value = gameObject.stats.synergo.g2.valFe

    gameElements[1].red.punt.value = gameObject.stats.redez.g2.points
    gameElements[1].red.n25.value = gameObject.stats.redez.g2.n25
    gameElements[1].red.per.value = gameObject.stats.redez.g2.character
    gameElements[1].red.valFe.value = gameObject.stats.redez.g2.valFe

    gameElements[2].syn.punt.value = gameObject.stats.synergo.g3.points
    gameElements[2].syn.n25.value = gameObject.stats.synergo.g3.n25
    gameElements[2].syn.per.value = gameObject.stats.synergo.g3.character
    gameElements[2].syn.valFe.value = gameObject.stats.synergo.g3.valFe

    gameElements[2].red.punt.value = gameObject.stats.redez.g3.points
    gameElements[2].red.n25.value = gameObject.stats.redez.g3.n25
    gameElements[2].red.per.value = gameObject.stats.redez.g3.character
    gameElements[2].red.valFe.value = gameObject.stats.redez.g3.valFe
}

$(document).ready(
    function () {
        videoLink.addEventListener('input', async function () {
            if (videoLink.value != 0) {
                var res = await fetch("/games/check/" + youtube_parser(videoLink.value), {
                    method: "GET",
                    headers: {
                        'Accept': 'application/json'
                    }
                })
                var resp = await res.text();
                let respJson = JSON.parse(resp)
                console.log(respJson)
                if (!('msg' in respJson)) {
                    fillGameSections(respJson[0])
                }
            }
        });
    }
)

