from websocket import create_connection
ws = create_connection("ws://127.0.0.1:8080/socket")
ws.send('{"channel": "orders","event": {"type": "NEW_ORDER","hash": "0x70034a326ab444412ae57b84dd30b0566c3a86e3cba15717d319847429444d12"}}')
result =  ws.recv()
print("Received '%s'" % result)
ws.close()
