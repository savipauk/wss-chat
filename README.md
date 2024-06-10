## wss-chat

university project for go class



clients register with a username and get a token

clients connect to ws server with the token they got

each websocket message needs to validate the token

tokens are stored in the `sessions` file

TODO tokens are automatically deleted after set amount of inactivity
