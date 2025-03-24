var userLang = navigator.language || navigator.userLanguage; 
//console.log(userLang);
var flag = document.getElementById("language-btn");
var acceptBtn = document.getElementById("accept-btn");
var declineBtn = document.getElementById("decline-btn");
var tempLang = "english";
if (document.cookie == "" && cookieConsent == "yes") {
  if (userLang == "cs-CZ") {
    document.cookie = "language=czech;";
  } else {
    document.cookie = "language=english;";
  }
}
displayLanguage();

flag.addEventListener('click', () => changeLanguage());
acceptBtn.addEventListener('click', function(){
  document.cookie = "cookie-consent=yes;";
});
declineBtn.addEventListener('click', function(){
  document.cookie = "cookie-consent=no;";
});
//console.log(document.cookie);
//console.log("consent: " + cookieConsent());

async function displayLanguage() {
    const data = await fetchData();
    const lang = cookieConsent() == "yes" ? getCookie() : tempLang;
    console.log(cookieConsent());
    console.log(lang);
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
    const currentLang = cookieConsent() == "yes" ? getCookie() : tempLang;
    console.log("curr lang: " + currentLang);
    if(cookieConsent()){
      document.cookie = `language=${currentLang === "czech" ? "english" : "czech"}; path=/`;
    }
    tempLang = currentLang == "czech" ? "english" : "czech";
  
    displayLanguage();
}

function getCookie(cookieName = "language") {
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
  if(getCookie('cookie-consent') == "yes"){
    return 1;
  }
  return -1;
}
