from fastapi import FastAPI
from dotenv import load_dotenv
from app.schemas import CVCreateRequest


app = FastAPI
load_dotenv()


@app.get("/ping")
async def ping():
    return {"ping": "Pong!"}

@app.get("/cv_creeate")
async def create_cv_by_desctiption(
    desctiption CVCreateRequest
):
    