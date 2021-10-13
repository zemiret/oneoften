socketUrl = () => "ws://127.0.0.1:8080/ws"

let livesSec, buzzerSec, seqNumberSecb
let ws;
let isDead = false;

const MessageDecreaseLive  = "DECREASE_LIVE";
const MessageBuzzer = "MESSAGE_BUZZER";
const MessagePlayerState = "PLAYER_STATE"

window.onload = () => {
    livesSec = document.getElementById("lives-sec");
    buzzerSec = document.getElementById("buzzer-sec");
    seqNumberSec = document.getElementById("seqNumber-sec");

    ws = connect(onMessage, () => console.log("CLOSED WS"))

    console.log('DOC LOADED');
}

onMessage = (message) => {
    console.log("GOT SOCKET MESSAGE ", message)

    if (!isDead) {
        switch (message.messageType) {
            case MessagePlayerState:
                console.log("PLAYER STATE", message.payload)

                if (message.payload.lives <= 0) {
                    isDead = true
                    renderDead();
                    break
                }

                renderLives(message.payload)
                renderSeqNumber(message.payload);
                break
            case MessageBuzzer:
                console.log("MESSAGE BUZZER", message.payload)
                renderBuzzer(message.payload);
                break
            default:
                console.log("UNKNOWN MESSAGE TYPE: ", message)
                break
        }
    }
}

connect = (
    onMessage,
    onClosed,
) => {
    let socket = new WebSocket(socketUrl());
    socket.onopen = () => {
        console.log("SOCKET OPEN")
    };
    socket.onerror = (err) => {
        console.log("SOCKET ERROR: ", err)
    };
    socket.onmessage = (event) => {
        onMessage(JSON.parse(event.data));
    };
    socket.onclose = onClosed;

    return socket;
};

renderDead = () => {
    document.body.innerHTML = '';
    document.body.style.backgroundColor = '#000000';
}

renderLives = (playerState) => {
    const text = document.createTextNode("ZYCIA: " + playerState.lives.toString());
    if (livesSec != null) {
        livesSec.innerHTML = '';
        livesSec.appendChild(text);
    }
}

renderSeqNumber = (playerState) => {
    const text = document.createTextNode("Pan(i) z numerem: " + playerState.seqNumber.toString());
    if (seqNumberSec != null) {
        seqNumberSec.innerHTML = '';
        seqNumberSec.appendChild(text);
    }
}

renderBuzzer = (buzzerMsg) => {
    const text = document.createTextNode("Odpowiada Pan(i) z numerem: " + buzzerMsg.seqNumber.toString());
    if (buzzerSec != null) {
        buzzerSec.innerHTML = '';
        buzzerSec.appendChild(text);
    }
}

onLiveDecrease = () => {
    console.log("LIVE DECREASE");

    ws.send(JSON.stringify({
        messageType: MessageDecreaseLive
    }))
}