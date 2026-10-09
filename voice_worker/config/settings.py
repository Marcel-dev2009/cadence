import os
from pathlib import Path
from dotenv import load_dotenv
BASE_DIR = Path(__file__).resolve().parent.parent
load_dotenv(BASE_DIR / ".env")
class Settings:
          PROJECT_NAME:str = "cadence-ai-worker"
          AI_PROVIDER:str = os.getenv("AI_PROVIDER", "llama")
          AI_API_KEY:str  = os.getenv("AI_API_KEY", "")
          # the worker host being 0.0.0.0 means that literally any network host on our LAN can connect it 
          WORKER_HOST:str = os.getenv("WORKER_HOST", "0.0.0.0") 
          WORKER_PORT:int = int(os.getenv("WORKER_PORT", 8081))
          # Golang backend integration we'll test the url if we use base url or the postfixed url
          GO_URL:str = os.getenv("GO_URL", "http://localhost:8080/api/v1")
          def validate(self):
                  if not self.AI_API_KEY:
                          print("API key missing please, ensure it's provided on your envars")

settings = Settings()
settings.validate()
