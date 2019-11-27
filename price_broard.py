from websocket import create_connection
ws = create_connection("ws://127.0.0.1:8080/socket")
ws.send('{"channel": "price_board","event": {"type": "UNSUBSCRIBE","payload": null}}')
result =  ws.recv()
print("Received '%s'" % result)
ws.close()
