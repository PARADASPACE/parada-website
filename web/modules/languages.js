var cookieForm = document.getElementById("cookie-consent");
if(cookieConsent() === 1){
  cookieForm.style.display = "none";
}
var userLang = navigator.language || navigator.userLanguage; 
var flag = document.getElementById("language-btn");
var acceptBtn = document.getElementById("accept-btn");
var declineBtn = document.getElementById("decline-btn");
var tempLang = getCookie("language") || (userLang === "cs-CZ" ? "czech" : "english");

if (cookieConsent() === 1 && getCookie("language") === "") {
  document.cookie = `language=${tempLang};`;
}


flag.addEventListener('click', () => changeLanguage());
acceptBtn.addEventListener('click', function(){
  document.cookie = "cookie-consent=yes;";
  document.cookie = `language=${tempLang};`;
  cookieForm.style.display = "none";
  displayLanguage();
});
declineBtn.addEventListener('click', function(){
  document.cookie = "cookie-consent=no; path=/; expires=Mon, 19 Feb 1979 00:00:00 UTC;";
  document.cookie = "language=; path=/; expires=Fri, 30 Mar 1979 00:00:00 UTC;";
  cookieForm.style.display = "none";
  displayLanguage();
});
//console.log(document.cookie);
//console.log("consent: " + cookieConsent());
displayLanguage();
async function displayLanguage() {
    const data = await fetchData();
    const lang = cookieConsent() === "1" ? getCookie() : tempLang;
    console.log("Cookie consent: " + cookieConsent());
    console.log("Language: " + lang);
    var obj = data[lang === "czech" ? 0 : 1];
    for(var key in obj){
      var value = obj[key];
      //console.log("key: " + key + " value: " +value );
      if(document.querySelector(('.'+key)) != null){
        document.querySelector(('.'+key)).textContent = value;
      }
    }

  
    flag.textContent = lang === "czech" ? "\uD83C\uDDE8\uD83C\uDDFF" : "\uD83C\uDDEC\uD83C\uDDE7";
}

async function changeLanguage() {
    const currentLang = cookieConsent() === 1 ? getCookie("language") : tempLang;
    const newLang = currentLang === "czech" ? "english" : "czech";

    if(cookieConsent()===1){
      document.cookie = `language=${newLang};`;
    }
    tempLang = newLang; 
  
    displayLanguage();
}


function getCookie(cookieName) {
  let name = cookieName + "=";
  let decodedCookie = decodeURIComponent(document.cookie);
  let val = decodedCookie.split(';');

  for (let i = 0; i < val.length; i++) {
    let c = val[i];
    while (c.charAt(0) == ' ') {
      c = c.substring(1);
    }
    if (c.indexOf(name) == 0) {
      return c.substring(name.length, c.length);
    }
  }
  return "";
}



async function fetchData() {
  try {
    const response = await fetch("/modules/languages.json")
    if (!response.ok) {
      throw new Error(`HTTP error! Status: ${response.status}`);
    }
    const data = await response.json(); 
    return data;
  } catch (error) {
    console.error("Error fetching JSON:", error);
  }
}

function cookieConsent(){
  return getCookie('cookie-consent') == "yes" ? 1 : -1;
}
