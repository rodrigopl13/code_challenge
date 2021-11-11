document.getElementById("login-form").addEventListener("submit", function (event){
    event.preventDefault();
    const url = "/v1/login"
    let data = {
        user_name: document.getElementById("username").value,
        password: document.getElementById("password").value
    }

    fetch(url, {
        method: "POST",
        headers: {"Content-Type": "application/json;charset=utf-8"},
        body: JSON.stringify(data)
    }).then(result => {
        result.json().then(body => {
            const response = body;
            if (result.ok){
                sessionStorage.setItem("token", result.headers.get("Authorization"));
                sessionStorage.setItem("user", JSON.stringify(response));
                window.location.href = "main.html";
            }else{
                const alertMsg = document.getElementById("alert-login");
                alertMsg.innerText = "Error authenticating user: " + JSON.stringify(response);
                alertMsg.style.visibility = "visible";
            }
        });

    });
    return false
});