const socket = new WebSocket("ws://localhost:8081/ws")

socket.onopen = () => {
    socket.send('test');
};

socket.onmessage = (event) => {
    console.log("recieved msg: ", event.data);
};
