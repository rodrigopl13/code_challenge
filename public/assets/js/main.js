const userRaw = sessionStorage.getItem("user");
user = JSON.parse(userRaw);
token = sessionStorage.getItem("token")
if (!token) {
    window.location.replace(`${location.protocol}//${location.hostname}/v1`);
}
const tokenWords = token.split(/(\s+)/).filter( function(e) { return e.trim().length > 0; } );
//const wsUrl = "wss://localhost/v1/ws?bearer="+tokenWords[1]; ENCRYPTED PROTOCOL
const wsUrl = "ws://localhost/v1/ws?bearer="+tokenWords[1];
let wsSocket = new WebSocket(wsUrl);
const messagesURL = `${location.protocol}//${location.hostname}/v1/message`

function onLoadMain(){
    const welcomeMsg = document.getElementById("welcome");
    welcomeMsg.innerText = 'Welcome '+ user.user_name;

    const url = new URL(messagesURL);
    const params = [['limit', '50']];
    url.search = new URLSearchParams(params).toString()
    fetch(url, {method: 'GET', headers: {'Authorization': token}})
        .then(response => response.json())
        .then(messages => {
            messages.forEach(message => {
                appendMessage(new Date(message.created_at), message.user.user_name, message.message);
            });
        }).catch(error =>{
            console.log(error)
    })

}

wsSocket.onmessage = function (event){
    const msg = JSON.parse(event.data);
    appendMessage(new Date(msg.created_at), msg.user.user_name, msg.message);
};

wsSocket.onerror = function (error){
    const alertMsg = document.getElementById("alert-login");
    alertMsg.innerText = "websocket error: " + JSON.stringify(error);
    alertMsg.style.visibility = "visible";
};

document.getElementById("chat-form").addEventListener("submit", function (event){
   event.preventDefault();
   const inText = document.getElementById("usermsg");
   const textMsg = inText.value;
   inText.value = "";
   msg = {
       user: {
           id: user.id,
           first_name: user.first_name,
           last_name: user.last_name,
           user_name: user.user_name
       },
       message: textMsg,
       created_at: new Date().toISOString()
   }
   wsSocket.send(JSON.stringify(msg));
});

document.getElementById("exit").onclick = function (){
    sessionStorage.removeItem("user");
    window.location.replace(`${location.protocol}//${location.hostname}/v1`);
    return false;
};



function appendMessage(time, username, msg){
    if(countChatMessages() > 49){
        removeChatMessage();
    }
    const msgHTML = `<div class="msgln"><span class="chat-time">${formatDate(time)}</span> <b class="user-name">${username}</b>${msg}<br></div>`;
    document.getElementById("chatbox").insertAdjacentHTML("beforeend", msgHTML);
}

function formatDate(date) {
    const h = "0" + date.getHours();
    const m = "0" + date.getMinutes();

    return `${h.slice(-2)}:${m.slice(-2)}`;
}

function countChatMessages(){
    const chat = document.getElementById("chatbox");
    return chat.getElementsByClassName("msgln").length;
}

function removeChatMessage(){
    const chat = document.getElementById("chatbox");
    const messages = chat.getElementsByClassName("msgln");
    chat.removeChild(messages[0]);
}