const pageTitle = "Gishath Fetch: MTG Price Checker for Singapore's LGS";
const form = document.getElementById("searchForm");
const lgsCheckboxesDiv = document.getElementById("lgsCheckboxes");
const searchInput = document.getElementById("search");
const submitBtn = document.getElementById("submitBtn");
const resultDiv = document.getElementById("result");
const resultCountDiv = document.getElementById("resultCount");
const lgsCheckboxes = document.getElementsByName('lgs[]');
const lgsOptions = [
    "Agora Hobby",
    "Cards Citadel",
    "Dueller's Point",
    "Flagship Games",
    "Games Haven",
    "Grey Ogre Games",
    "Hideout",
    "Mana Pro",
    "Mox & Lotus",
    "MTG Asia",
    "OneMtg",
    "Sanctuary Gaming"
];
let timeouts = [];
let baseUrl = "https://gishathfetch.com/";
let apiBaseUrl = "https://api.gishathfetch.com/";

if (window.location.hostname === "staging.gishathfetch.com") {
    baseUrl = "https://staging.gishathfetch.com/";
    apiBaseUrl = "https://staging-api.gishathfetch.com/";
}

setupConfig();

// Pre-select checkboxes and pre-fill search from cookie
function setupConfig() {
    appendLgsCheckboxes();
    setupEventListeners();
    onloadSearch();
}

function onloadSearch() {
    const urlParams = new URLSearchParams(window.location.search);
    const searchParam = urlParams.get('s');

    if (searchParam !== "") {
        searchInput.value = searchParam;
        submitBtn.click();
    }
}

function setupEventListeners() {
    form.addEventListener("submit", onFormSubmit);

    document.addEventListener("keypress", function(event) {
        if (event.keyCode === 13) {
            event.preventDefault();
            submitBtn.click();
        }
    });
}

function appendLgsCheckboxes() {
    let lgsSelected = [];

    if(localStorage.getItem('lgsSelected') !== null && localStorage.getItem('lgsSelected') !== undefined && localStorage.getItem('lgsSelected') !== "") {
        lgsSelected = decodeURIComponent(localStorage.getItem('lgsSelected')).split(",");
    } else {
        lgsSelected = lgsOptions;
    }

    lgsCheckboxesDiv.innerHTML = '';
    for(let i=0; i<lgsOptions.length; i++) {
        let isChecked = lgsSelected.includes(lgsOptions[i]) ? "checked" : "";
        lgsCheckboxesDiv.innerHTML += `
                <div class="form-check form-check-inline">
                  <input class="form-check-input" type="checkbox" id="lgsCheckbox`+i+`" class="lgsCheckboxes" value="`+lgsOptions[i]+`" name="lgs[]" `+isChecked+`>
                  <label class="form-check-label" for="lgsCheckbox`+i+`">`+lgsOptions[i]+`</label>
                </div>
              `;
    }
}

function clearTimeouts() {
    for (let i=0; i<timeouts.length; i++) {
        clearTimeout(timeouts[i]);
    }
}

// Timeout 15s in backend
function updateSubmitBtnProgress() {
    submitBtn.innerHTML = "Searching LGS"

    for(let i=1; i<=15; i++){
        timeouts.push(window.setTimeout(function(){
            submitBtn.innerHTML += " ."
        }, i * 1000));
    }
}

function resetResult() {
    resultDiv.innerHTML = "";
    resultCountDiv.innerHTML = "";
}

function resetSubmitBtn() {
    clearTimeouts();
    submitBtn.innerHTML = "Search";
    submitBtn.disabled = false;
}

function updatePageUrlTitle(searchStr, url) {
    window.history.pushState(searchStr.toLowerCase(), searchStr.toLowerCase() + " | " + pageTitle, url);
    document.title = searchStr.toLowerCase() + " | " + pageTitle;
}

function onFormSubmit(event) {
    event.preventDefault();

    let searchStr = searchInput.value.trim()

    // End if empty search str
    if (searchStr === "" || searchStr.length < 3) {
        return
    }

    // Tag search str
    gtag('event', 'search', {
        'search_term': searchStr.toLowerCase()
    });

    let lgsSelected = [];

    for(let i=0; i<lgsCheckboxes.length; i++) {
        if (lgsCheckboxes[i].checked) {
            lgsSelected.push(lgsCheckboxes[i].value)
        }
    }

    if (lgsSelected.length === 0) {
        lgsSelected = lgsOptions;
        for(let i=0; i<lgsCheckboxes.length; i++) {
            lgsCheckboxes[i].checked = true;
        }
    }

    // Set state to disabled
    submitBtn.disabled = true;
    // Reset result div
    resetResult();

    let request = new XMLHttpRequest();
    let searchQueryString = "?s="+encodeURIComponent(searchStr.toLowerCase());
    let searchUrl = apiBaseUrl + searchQueryString
    searchUrl += "&lgs=" + encodeURIComponent(lgsSelected.join(','));

    localStorage.setItem("lgsSelected", encodeURIComponent(lgsSelected.join(",")));

    request.open("GET", searchUrl);
    request.send();

    updateSubmitBtnProgress();

    request.onreadystatechange = function() {
        if (request.readyState === XMLHttpRequest.DONE) {
            let resultCount = 0;

            // Check the status of the response
            if (request.status === 200) {
                // Access the data returned by the server
                let result = JSON.parse(request.responseText);
                // Do something with the data
                if (result.hasOwnProperty("data")) {
                    if (result["data"] !== null && result["data"].length > 0) {
                        updatePageUrlTitle(searchStr, baseUrl + searchQueryString);
                        let html = `<div class="row">`;
                        for(let i = 0; i < result["data"].length; i++) {
                            if (result["data"][i].hasOwnProperty("url")
                                && result["data"][i].hasOwnProperty("img")
                                && result["data"][i].hasOwnProperty("name")
                                && result["data"][i].hasOwnProperty("price")
                                && result["data"][i].hasOwnProperty("src")) {
                                let h = `
                                  <div class="col-lg-3 col-6 mb-4">
                                    <div class="text-center mb-2">
                                      <a href="`+result["data"][i]["url"]+`" target="_blank">
                                        <img src="`+(result["data"][i]["img"]===""?`https://placehold.co/304x424?text=`+result["data"][i]["name"]:result["data"][i]["img"])+`" loading="lazy" class="img-fluid w-100" alt="`+result["data"][i]["name"]+`"/>
                                      </a>
                                    </div>
                                    <div class="text-center">
                                      <div class="fs-6 lh-sm fw-bold mb-1">`+result["data"][i]["name"]+`</div>
                                      `+((result["data"][i].hasOwnProperty("quality") && result["data"][i]["quality"]!=="")?`<div class="fs-6 lh-sm fw-bold mb-1">≪ `+result["data"][i]["quality"]+` ≫</div>`:``)+`
                                      <div class="fs-6 lh-sm">S$ `+result["data"][i]["price"].toFixed(2)+`</div>
                                      <div><a href="`+result["data"][i]["url"]+`" target="_blank" class="link-offset-2">`+result["data"][i]["src"]+`</a></div>
                                    </div>
                                  </div>`;
                                html += h
                                resultCount++;
                            }
                        }
                        html += `</div>`
                        resultDiv.innerHTML = html;
                    }
                }

                // Tag search str
                gtag('event', 'view_search_results', {
                    'search_term': searchStr.toLowerCase()
                });

            } else {
                // Handle error
            }

            resultCountDiv.innerHTML = `<div class="py-2">`+resultCount+` result`+(resultCount>1?"s":"")+` found</div>`;

            // Reset state
            resetSubmitBtn();
        }
    };
}
