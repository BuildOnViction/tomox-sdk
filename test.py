from websocket import create_connection
#import socket
ws = create_connection("ws://127.0.0.1:8080/socket")
ws.send('{"channel": "orders","event": {"type": "NEW_ORDER","hash": null,"payload": {}}}')
result =  ws.recv()
print("Received '%s'" % result)
ws.close()
