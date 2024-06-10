const connectContainer = document.querySelector("#connectContainer")
const chatWrapper = document.querySelector("#chatWrapper")
const chat = document.querySelector("#chat")
const inputChatter = document.querySelector("#chatter")
const inputMessage = document.querySelector("#chatMessage")

/**
 * @type {WebSocket | null}
 */
let socket;

async function ConnectToServer() {
        connectContainer.style.display = "none"
        chatWrapper.style.display = "flex"
        const chatter = inputChatter.value

        const data = {
                "username": chatter
        }

        await fetch("http://localhost:8081/register", {
                method: 'POST',
                headers: {
                        'Content-Type': 'application/json'
                },
                body: JSON.stringify(data)
        })
                .then(async (res) => {
                        const token = JSON.parse(await res.text()).token
                        console.log("token", token)

                        sessionStorage.setItem("token", token)

                        ConnectWS()
                })
                .catch((err) => {
                        console.error(err)
                        LeaveChat()
                })
}


function ConnectWS() {
        const token = sessionStorage.getItem("token")
        socket = new WebSocket(`ws://localhost:8081/ws?token=${token}`)

        socket.addEventListener("error", (event) => {
                console.error("WebSocket error: " + event);

                LeaveChat()
        });

        socket.onclose = (event) => {
                console.log("WebSocket close: " + event);

                LeaveChat()
        }

        socket.onopen = () => {
                console.log("Chatter:", chatter)
        };

        socket.onmessage = (event) => {
                console.log("recieved msg: ", event.data)

                let newMsg = document.createElement("div")
                newMsg.classList.add("msg")
                newMsg.classList.add("otherMessage")
                newMsg.innerHTML = event.data

                chat.appendChild(newMsg)
        };
}


function LeaveChat() {
        if (socket != null) {
                socket.close()
        }

        connectContainer.style.display = "flex"
        chatWrapper.style.display = "none"
}


function SendMessage() {
        const message = inputMessage.value
        const token = sessionStorage.getItem("token")
        inputMessage.value = ""

        const data = {
                message: message,
                token: token
        }

        socket.send(JSON.stringify(data))

        let newMsg = document.createElement("div")
        newMsg.classList.add("msg")
        newMsg.classList.add("myMessage")
        newMsg.innerHTML = message

        chat.appendChild(newMsg)
}

