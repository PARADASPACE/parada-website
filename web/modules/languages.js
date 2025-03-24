var userLang = navigator.language || navigator.userLanguage; 
console.log(userLang);
var flag = document.getElementById("language-btn");

if (document.cookie == "") {
  if (userLang == "cs-CZ") {
    document.cookie = "language=czech;";
  } else {
    document.cookie = "language=english;";
  }
}
displayLanguage();

flag.addEventListener('click', () => changeLanguage());

async function displayLanguage() {
    const data = await fetchData();
    const lang = await getCookieLang('language');

    var obj = data[lang === "czech" ? 0 : 1];
    for(var key in obj){
      var value = obj[key];
      console.log("key: " + key + " value: " +value );
      if(document.querySelector(('.'+key)) != null){
        document.querySelector(('.'+key)).textContent = value;
      }
    }

  
    flag.textContent = lang === "czech" ? "\uD83C\uDDE8\uD83C\uDDFF" : "\uD83C\uDDEC\uD83C\uDDE7";
}

async function changeLanguage() {
    const currentLang = await getCookieLang('language');
    document.cookie = `language=${currentLang === "czech" ? "english" : "czech"}; path=/`;
    displayLanguage();
}


async function getCookieLang(cookieName) {
  let name = cookieName+"=";
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
    console.log("Header Paragraph (CZ):", data[1].header_paragraph); // Access specific value
    return data;
  } catch (error) {
    console.error("Error fetching JSON:", error);
  }
}

