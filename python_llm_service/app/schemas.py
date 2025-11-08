from pydantic import BaseModel

class CVCreateRequest(BaseModel):
    description: str