const contentDiv = document.getElementById("content");
const navForm = document.getElementById("navForm");
const navSearchInput = document.getElementById("navSearch");
let baseUrl = "https://gishathfetch.com/";

if (window.location.hostname === "staging.gishathfetch.com") {
    baseUrl = "https://staging.gishathfetch.com/";
}

navForm.addEventListener("submit", onNavFormSubmit);

function onNavFormSubmit(event) {
    event.preventDefault();

    let searchStr = navSearchInput.value.trim()

    // End if empty search str
    if (searchStr === "" || searchStr.length < 3) {
        return
    }

    window.open(baseUrl + "?s="+encodeURIComponent(searchStr.toLowerCase()));
}