let lgsMap = [
    {
        "id": "agora-map",
        "name": "Agora Hobby",
        "address": "French Rd, #05-164 Blk 809, Singapore 200809",
        "iframe": "https://www.google.com/maps/embed?pb=!1m18!1m12!1m3!1d3988.778050505021!2d103.85967687451628!3d1.3084089617085968!2m3!1f0!2f0!3f0!3m2!1i1024!2i768!4f13.1!3m3!1m2!1s0x31da19c9f7d7f74d%3A0xeaa1a66df7d4bcd6!2sAgora%20Hobby!5e0!3m2!1sen!2ssg!4v1702820213937!5m2!1sen!2ssg",
    },
    {
        "id": "cards-citadel-map",
        "name": "Cards Citadel",
        "address": "464 Crawford Ln, #02-01, Singapore 190464",
        "iframe": "https://www.google.com/maps/embed?pb=!1m18!1m12!1m3!1d3988.783678524258!2d103.85966947451631!3d1.3048646617197366!2m3!1f0!2f0!3f0!3m2!1i1024!2i768!4f13.1!3m3!1m2!1s0x31da190c9e183751%3A0xa2119a95d1e683f2!2sCards%20Citadel!5e0!3m2!1sen!2ssg!4v1702820792196!5m2!1sen!2ssg",
    },
    {
        "id": "dueller-point-map",
        "name": "Dueller's Point",
        "address": "450 Hougang Ave 10, B1-541, Singapore 530450",
        "iframe": "https://www.google.com/maps/embed?pb=!1m18!1m12!1m3!1d3988.662159756766!2d103.89300967451602!3d1.3793695614811952!2m3!1f0!2f0!3f0!3m2!1i1024!2i768!4f13.1!3m3!1m2!1s0x31da163eecb250ff%3A0xc7c259e72671dc62!2sDueller&#39;s%20Point!5e0!3m2!1sen!2ssg!4v1702820876967!5m2!1sen!2ssg",
    },
    {
        "id": "flagship-games-map",
        "name": "Flagship Games",
        "address": "5 Jln Pemimpin, #03-01A, Singapore 577197",
        "iframe": "https://www.google.com/maps/embed?pb=!1m18!1m12!1m3!1d3988.7100448274964!2d103.8376813757648!3d1.3505012580854614!2m3!1f0!2f0!3f0!3m2!1i1024!2i768!4f13.1!3m3!1m2!1s0x31da173ef6ffcc0b%3A0x880386dee363a253!2sFlagship%20Games!5e0!3m2!1sen!2ssg!4v1702820928970!5m2!1sen!2ssg",
    },
    {
        "id": "games-haven-pl-map",
        "name": "Games Haven - Paya Lebar",
        "address": "736 Geylang Rd, Singapore 389647",
        "iframe": "https://www.google.com/maps/embed?pb=!1m18!1m12!1m3!1d63819.358332241325!2d103.79905633083244!3d1.350592080054757!2m3!1f0!2f0!3f0!3m2!1i1024!2i768!4f13.1!3m3!1m2!1s0x31da1817d10ac901%3A0x2cacb3a0679089a2!2sGames%20Haven!5e0!3m2!1sen!2ssg!4v1702821045126!5m2!1sen!2ssg",
    },
    {
        "id": "games-haven-ct-map",
        "name": "Games Haven - Chinatown",
        "address": "531 Upper Cross St, #02-02, Singapore 050531",
        "iframe": "https://www.google.com/maps/embed?pb=!1m18!1m12!1m3!1d63821.04849133675!2d103.76901934863278!3d1.2846211999999948!2m3!1f0!2f0!3f0!3m2!1i1024!2i768!4f13.1!3m3!1m2!1s0x31da196e26ea750d%3A0x688a1d11efebe621!2sGames%20Haven%20-%20Chinatown!5e0!3m2!1sen!2ssg!4v1702821183177!5m2!1sen!2ssg",
    },
    {
        "id": "games-haven-amk-map",
        "name": "Games Haven - Ang Mo Kio",
        "address": "51 Ang Mo Kio Ave 3, #03-01, Singapore 569922",
        "iframe": "https://www.google.com/maps/embed?pb=!1m18!1m12!1m3!1d63821.04849133675!2d103.76901934863278!3d1.2846211999999948!2m3!1f0!2f0!3f0!3m2!1i1024!2i768!4f13.1!3m3!1m2!1s0x31da17d78ce21027%3A0x2ce73b73e4a0a9b7!2sGames%20Haven%20-%20Ang%20Mo%20Kio!5e0!3m2!1sen!2ssg!4v1702821196724!5m2!1sen!2ssg",
    },
    {
        "id": "games-haven-je-map",
        "name": "Games Haven - Jurong East",
        "address": "131 Jurong Gateway Rd, #02-245, Singapore 600131",
        "iframe": "https://www.google.com/maps/embed?pb=!1m18!1m12!1m3!1d63818.86192364747!2d103.77101964863279!3d1.3693644999999974!2m3!1f0!2f0!3f0!3m2!1i1024!2i768!4f13.1!3m3!1m2!1s0x31da10057f577829%3A0x63e8b2fbab922947!2sGames%20Haven%20-%20Jurong%20East!5e0!3m2!1sen!2ssg!4v1702821214358!5m2!1sen!2ssg",
    },
    {
        "id": "grey-ogre-map",
        "name": "Grey Ogre Games",
        "address": "83 Club St, Singapore 069451",
        "iframe": "https://www.google.com/maps/embed?pb=!1m18!1m12!1m3!1d3988.8199964760065!2d103.84085797576442!3d1.2817574584814586!2m3!1f0!2f0!3f0!3m2!1i1024!2i768!4f13.1!3m3!1m2!1s0x31da190d242b70db%3A0x965b932c3bc19eda!2sGrey%20Ogre%20Games!5e0!3m2!1sen!2ssg!4v1702821297360!5m2!1sen!2ssg",
    },
    {
        "id": "hideout-map",
        "name": "Hideout",
        "address": "803 King George's Ave, #02-190, Singapore 200803",
        "iframe": "https://www.google.com/maps/embed?pb=!1m18!1m12!1m3!1d15955.112777358516!2d103.84179288715819!3d1.3083185!2m3!1f0!2f0!3f0!3m2!1i1024!2i768!4f13.1!3m3!1m2!1s0x31da19e075f4c4f5%3A0x60e4a2c61816be63!2sHideout!5e0!3m2!1sen!2ssg!4v1702821327690!5m2!1sen!2ssg",
    },
    {
        "id": "manapro-map",
        "name": "Mana Pro",
        "address": "BLK 203 Choa Chu Kang Ave 1, B1-41, Singapore 680203",
        "iframe": "https://www.google.com/maps/embed?pb=!1m18!1m12!1m3!1d3988.6584888121897!2d103.74693327451605!3d1.3815577614740542!2m3!1f0!2f0!3f0!3m2!1i1024!2i768!4f13.1!3m3!1m2!1s0x31da1176665e2737%3A0x3b8608ab4d67724f!2sMana%20Pro!5e0!3m2!1sen!2ssg!4v1702821359528!5m2!1sen!2ssg",
    },
    {
        "id": "mox-map",
        "name": "Mox & Lotus",
        "address": "789 Geylang Rd, #3rd Floor, Singapore 389675",
        "iframe": "https://www.google.com/maps/embed?pb=!1m18!1m12!1m3!1d3988.7673317965505!2d103.88685357576463!3d1.3151327582901533!2m3!1f0!2f0!3f0!3m2!1i1024!2i768!4f13.1!3m3!1m2!1s0x31da19d89d198d6b%3A0xb3e238feedd6c90d!2sMox%20%26%20Lotus!5e0!3m2!1sen!2ssg!4v1702821402075!5m2!1sen!2ssg",
    },
    {
        "id": "mtg-asia-map",
        "name": "MTG Asia",
        "address": "261 Waterloo St, #03-28, Singapore 180261",
        "iframe": "https://www.google.com/maps/embed?pb=!1m18!1m12!1m3!1d3988.7930896678654!2d103.8493947744998!3d1.2989162986887468!2m3!1f0!2f0!3f0!3m2!1i1024!2i768!4f13.1!3m3!1m2!1s0x31da19bb4a2bee83%3A0x28725aa3a3e2a51!2sMTG-Asia!5e0!3m2!1sen!2ssg!4v1703085334392!5m2!1sen!2ssg",
    },
    {
        "id": "onemtg-map",
        "name": "One MTG",
        "address": "100 Jln Sultan, #03-11 Sultan Plaza, Singapore 199001",
        "iframe": "https://www.google.com/maps/embed?pb=!1m18!1m12!1m3!1d3988.7866900551694!2d103.85910407451628!3d1.3029641617257042!2m3!1f0!2f0!3f0!3m2!1i1024!2i768!4f13.1!3m3!1m2!1s0x31da19180d91f3a1%3A0x75c807bf93d430a4!2sOne%20MTG!5e0!3m2!1sen!2ssg!4v1702821425238!5m2!1sen!2ssg",
    },
    {
        "id": "sanctuary-map",
        "name": "Sanctuary Gaming",
        "address": "277 Orchard Rd, #04-09 Orchard, Gateway 238858",
        "iframe": "https://www.google.com/maps/embed?pb=!1m18!1m12!1m3!1d3988.7903463145303!2d103.83467127576455!3d1.3006530583733655!2m3!1f0!2f0!3f0!3m2!1i1024!2i768!4f13.1!3m3!1m2!1s0x31da19972abb859d%3A0x2a935be2b43c260d!2sSanctuary%20Gaming!5e0!3m2!1sen!2ssg!4v1702821457166!5m2!1sen!2ssg",
    }
]

const mapDivBody = document.getElementById("map").getElementsByClassName("modal-body")[0];

let mapListHtml = "";
let mapItemsHtml = "";
let topBtnHtml = `<a href="#map-list" class="btn btn-primary" role="button">Back to top</a>`;
let closeBtnHtml = `<button type="button" class="btn btn-secondary" style="margin-left: 7px;" data-bs-dismiss="modal">Close</button>`;

// Menu
mapListHtml += `<div class="mb-4"><ul style="padding-left: 1rem">`;
for(let i = 0; i < lgsMap.length; i++) {
    mapListHtml += `<li><a href="#`+lgsMap[i].id+`" class="link-offset-2" alt="`+lgsMap[i].name+`">`+lgsMap[i].name+`</a></li>`;
    mapItemsHtml += `
        <div id="`+lgsMap[i].id+`" class="mb-4 map-item">
            <h5>`+lgsMap[i].name+`</h5>
            <div class="mb-2">`+lgsMap[i].address+`</div>
            <iframe class="w-100 h-100 border border-dark mb-3" src="`+lgsMap[i].iframe+`" allowfullscreen="" loading="lazy" referrerpolicy="no-referrer-when-downgrade"></iframe>
            <div>`+topBtnHtml+closeBtnHtml+`</div>
        </div>    
    `;
}
mapListHtml += `<ul></div>`;
mapDivBody.innerHTML = mapListHtml + mapItemsHtml;