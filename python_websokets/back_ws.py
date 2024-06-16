import asyncio
import websockets

async def hello(websocket):
    json_data = await websocket.recv()
    print(f"<<< {json_data}")


    # greeting = f"Hello {name}!"

    await websocket.send('Hello')
    # print(f">>> {greeting}")

async def main():
    async with websockets.serve(hello, "0.0.0.0", 63345):
        await asyncio.Future()  # run forever

if __name__ == "__main__":
    asyncio.run(main())