const socket = new WebSocket("ws://localhost:8080/ws");

socket.onopen = () => {
  console.log("Connected to server!");
  socket.send("ping");
};

socket.onmessage = (event) => {
  console.log("Server says:", event.data);
  setTimeout(() => {
    socket.send("ping");
  }, 1000);
};

socket.onclose = () => {
  console.log("Disconnected from server.");
};
