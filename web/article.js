const contentDiv = document.getElementById("content");
const navForm = document.getElementById("navForm");
const navSearchInput = document.getElementById("navSearch");
let baseUrl = "https://gishathfetch.com/";

if (window.location.hostname === "staging.gishathfetch.com") {
    baseUrl = "https://staging.gishathfetch.com/";
}

navForm.addEventListener("submit", onNavFormSubmit);

feather.replace();

function onNavFormSubmit(event) {
    event.preventDefault();

    let searchStr = navSearchInput.value.trim()

    // End if empty search str
    if (searchStr === "" || searchStr.length < 3) {
        return
    }

    window.open(baseUrl + "?s="+encodeURIComponent(searchStr.toLowerCase()));
}

function updateMetaTags(name, description, url){
    document.title = name;
    document.querySelector('meta[name="description"]').setAttribute("content", description);
    document.querySelector('meta[property="og:title"]').setAttribute("content", name);
    document.querySelector('meta[property="og:url"]').setAttribute("content", baseUrl + url);
    document.querySelector('meta[property="og:description"]').setAttribute("content", description);
}