# golang-websocket
Example websocket contain client and server for simulate realwold working.
# Server 
Server waiting client connect with the request channel which used to be a channel of redis.
when server open connection system will subscript channel in redis (PUB/SUB) and if the system can receive the message then system will forward the message to client side.
PS. to run this, needed to publish the message from redis mannually

# Client
Client just send the channel request to sever and receive all message
