import asyncio
import websockets
import json

addr = "ws://localhost:63342/front-ws"
message = '{"message": "idk"}'

async def hello():
    async with websockets.connect(addr) as websocket:

        await websocket.send(message)
        # print(f">>> {message}")

        greeting = await websocket.recv()
        print(f"<<< {greeting}")

if __name__ == "__main__":
    asyncio.run(hello())