document.getElementById("click").addEventListener("click", function(evt){
    let xhr = new XMLHttpRequest()
    xhr.open("GET", "/context-cancel", true)
    xhr.send()
    setTimeout(() => {
        xhr.abort()
    }, 1000)
})