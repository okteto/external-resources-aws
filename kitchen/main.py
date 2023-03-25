import asyncio
import concurrent
import json
import logging
import os
import time

from typing import List

from fastapi import FastAPI, WebSocket, WebSocketDisconnect
from fastapi.concurrency import run_in_threadpool
from fastapi.staticfiles import StaticFiles
from pydantic import BaseModel

import boto3
import uvicorn
import pydantic


class FoodCooked(BaseModel):
    orderId: str
    item: str

class FoodOrder(BaseModel):
    orderId: str
    item: List[str]


class SQSConsumer:
    def __init__(self):
        sqsQueName = os.getenv("QUEUE")
        self._sqs = boto3.resource('sqs')
        self._sqsQueue = self._sqs.get_queue_by_name(QueueName=sqsQueName)
        self._queue = asyncio.Queue()
        self._task: asyncio.Task

    async def get_food_order(self):
        return await self._queue.get();

    async def start(self):
        self._task = asyncio.create_task(self._pump_messages)

    async def shutdown(self):
        if self._task.done():
            self._task.result()
        else:
            self._task.cancel()
    
    def _pump_messages(self):
        while True:
            messages = self._sqsQueue.receive_messages()
            for message in messages:
                try:
                    order = json.dumps(message.body, default=pydantic_encoder)
                    self.queue.put(order)
                    message.delete()
                except Exception as e:
                    logging.error('error:', e)
                    break

class WSConnectionManager:
    def __init__(self):
        self.active_connections: List[WebSocket] = []

    async def connect(self, websocket: WebSocket):
        await websocket.accept()
        self.active_connections.append(websocket)

    def disconnect(self, websocket: WebSocket):
        self.active_connections.remove(websocket)

    async def send_food_order(self, order: FoodOrder, websocket: WebSocket):
        await websocket.send_json(order.json())

connectionManager = WSConnectionManager()
sqsConsumer = SQSConsumer()
app = FastAPI()

@app.post("/cooked", status_code=201)
async def cooked(cooked: FoodCooked):
    logging.info("item cooked:", cooked.item)


@app.websocket("/ws")
async def websocket_endpoint(websocket: WebSocket):
    await connectionManager.connect(websocket)
    try:
        while True:
            order = await sqsConsumer.get_food_order()
            await connectionManager.send_food_order(order)
    except WebSocketDisconnect:
        connectionManager.disconnect(websocket)
    
@app.on_event("startup")
async def app_startup():
    await sqsConsumer.start()
    return

@app.on_event("shutdown")
async def app_shutdown():
    await sqsConsumer.shutdown()
    return

app.mount("/", StaticFiles(directory="public",html=True), name="public")

if __name__ == "__main__":
    
    uvicorn.run("main:app", host="0.0.0.0", port=8000, reload=True)