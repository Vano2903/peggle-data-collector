function initUserArea() {
    document.getElementById("user-welcome").innerHTML = "Benvenuto " + user.username;
}

function showError(msg, id) {
    $(id).innerHTML = msg;
}

async function updatePassword() {
    currentPass = document.getElementById("currentPassword").value;
    newPass = document.getElementById("newPassword").value;
    confirmNewPass = document.getElementById("confirmNewPassword").value;
    if (currentPass != user.password) {
        $("#errCurrentPass").innerHTML = "la password inserita é diversa da quella corrente";
        return
    }
    if (newPass == "") {
        $("#errNewPassword").innerHTML = "la nuova password non puó essere vuota";
        return
    }
    if (newPass == currentPass) {
        $("#errNewPassword").innerHTML = "la nuova password non puó essere uguale a quella corrente";
        return
    }
    if (newPass != confirmNewPass) {
        $("#errConfirmPass").innerHTML = "le due password non corrispondono";
        return
    }
    const res = await fetch('/users/customization/?password=' + newPass, {
        method: "POST",
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(user)
    });
    const resp = await res.json();
    if (res.status != 200) {
        alert("qualcosa é andato storto, riprova")
        return
    }
    alert("modifica andata a buon fine, verrá richiesto di fare login nuovamente")
    //logout()
}

async function updatePfp() {
    newPfp = document.getElementById("newPfpUrl").value;
    const res = await fetch('/users/customization/?pfp=' + newPfp, {
        method: "POST",
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(user)
    });
    const resp = await res.json();
    if (res.status != 200) {
        alert("qualcosa é andato storto, riprova")
        return
    }
    alert("modifica andata a buon fine, verrá richiesto di fare login nuovamente")
    //logout()
}

function hideForms() {
    document.getElementById("updatePasswordForm").style.display = "none";
    document.getElementById("updatePfpForm").style.display = "none";
}

function showPfpUpdateForm() {
    hideForms()
    document.getElementById("updatePfpForm").style.display = "block";
}

function showPasswordUpdateForm() {
    hideForms()
    document.getElementById("updatePasswordForm").style.display = "block";
}
