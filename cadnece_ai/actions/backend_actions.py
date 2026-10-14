import requests
from typing import Optional
from config import settings
from fastapi.concurrency import run_in_threadpool
class BackendActions:
          def __init__(self):
                  self.go_backend_url = getattr(settings, "GO_URL", "http://localhost:8080")
                  self.session = requests.Session()
          async def sync_to_golang(self, event_name:str, event_time:str, append_cookies:Optional[dict] = None) -> bool:
                  endpoint = f"{self.go_backend_url}/api/v1/event/create"
                  query_params_body = {
                    "event_time" : event_time      
                  }
                  json_body = {
                          "event_name" : event_name
                
                  }
                  try:
                          print(f"syncing py worker to golang 1.27.1 gin handlers via {endpoint}")
                          response = await run_in_threadpool(
                           self.session.post,
                           url=endpoint,
                           params=query_params_body,
                           json=json_body,
                           cookies=append_cookies,
                           timeout=5.0
                           )
                          if response.status_code in (200, 201):
                                  print(f"successfully synced to the go backend via {endpoint}")
                                  return True
                           
                          print(f"sync operation to our go backend returned bad status code {response.status_code}")
                          return False
                  except Exception as e:
                          print(f"error syncing to go backend:{e}")
                          return False
                          
                          