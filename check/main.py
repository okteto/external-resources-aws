import os
import uvicorn

from fastapi import FastAPI, HTTPException
from fastapi.staticfiles import StaticFiles

from pydantic import BaseModel
from typing import List


app = FastAPI()

checks = {}


class Item(BaseModel):
    name: str
    price: float | None = 0
    ready: bool | None = False

class Check(BaseModel):
    orderId: str
    items: List[Item]
    total: float | None = 0


@app.get("/healthz")
async def healthz():
    return {"message": "Check please 🫰!"}


@app.get("/checks", response_model=list[Check])
async def getChecks():
    response = []
    for checkID in checks:
        response.append(checks[checkID])

    print(response)
    return response

@app.post("/checks", status_code=200)
async def prepare_check(check: Check):
    # calculate price
    total = 0
    for i in range(len(check.items)):
        item = check.items[i]
        price = len(item.name)
        check.items[i].price = price
        total += price

    check.total = total
    checks[check.orderId] = check
    print(("The total for check {check_id} is: ${total} 🧮").format(check_id=check.orderId, total=check.total))


@app.get("/checks/{check_id}")
async def getCheck(check_id):
    if check_id in checks.keys():
        return checks[check_id]
    raise HTTPException(status_code=404, detail="Check not found 👎🏼")

@app.delete("/checks/{check_id}")
async def payCheck(check_id):
    if check_id in checks.keys():
        del checks[check_id]
        return
    
    raise HTTPException(status_code=404, detail="Check not found 👎🏼")

app.mount("/", StaticFiles(directory="public", html=True), name="public")

if __name__ == "__main__":
   reload=bool(os.getenv("RELOAD"))
   uvicorn.run("main:app", host="0.0.0.0", port=8000, reload=reload)