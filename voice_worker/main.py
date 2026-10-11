import base64
import uvicorn
from fastapi import FastAPI, UploadFile, HTTPException
from fastapi.middleware import cors
from core import core_state_manager
from services import stt_service_engine, tts_service_engine
from actions import call_backend, app_action_instruct
app = FastAPI(
           title="Cadence AI Worker",
           description="Production-ready core audio processing and intent extraction matrix"
)
app.add_middleware(
          cors.CORSMiddleware,
          allow_origins=["*"],
          allow_credentials=True,
          allow_methods=["*"],
          allow_headers=["*"]
)
@app.post("/app/voice/process") # route that connects flutter to python
async def process_voice_file(file:UploadFile = File(...)):
          if not file:
                  raise HTTPException(status_code=400, detail="No audio file provided")
          if not core_state_manager.try_processing_lock():
                  raise HTTPException(
                          status_code=429,
                          detail="System is currently busy"
                  )
          try:
                  audio_bytes = await file.read()
                  ai_result = stt_service_engine.transcribe_and_route(audio_bytes)
                  intent = ai_result.get("intent", "unknown")
                  client_directive = {}
                  voice_reply_text = "Voice Operation Completed Succesfully"
                  operation_success = True
                  if intent == "change_theme":
                          theme_value = ai_result.get("value", "light")
                          client_directive = app_action_instruct.format_theme_mode(theme_value)
                          voice_reply_text = f"Changing theme to {theme_value}"
                  elif intent == "open_settings":
                          client_directive = app_action_instruct.format_navigation("settings")
                          voice_reply_text = "routing to settings screen"
                  elif  intent == "create_event":
                          time_str = ai_result.get("event_time","")
                          name = ai_result.get("event_name", "new event")
                          operation_success = call_backend.sync_to_golang(event_name=name, event_time=time_str)
                          if operation_success:
                                  voice_reply_text = f"Successfully scheduled an event for {name} to take place at {time_str}"
                                  client_directive = {
                                          "action_type":"ui_update",
                                          "target":"calender",
                                          "payload":{"refresh":True}
                                  }
                          else:
                                  voice_reply_text = f"Sorry I ran into an error scheduling {name}"
                                  client_directive = {"action_type": "error", "payload": {"message": "Go Backend connection failure"}}
                  else:
                          voice_reply_text = "I'm not completley sure how to process that command yet"
                          client_directive = {
                                  "action_type":"unknown",
                                  "target":"nil",
                                  "payload":{}
                          }
                  reply_audio_base64 = ""
                  if operation_success:
                          audio_data = tts_service_engine.synthesize_speech(voice_reply_text)
                          reply_audio_base64 = base64.b64encode(audio_data).decode('utf-8')  
                  return {
                    "intent": intent,
                    "directive": client_directive,
                    "spoken_text": voice_reply_text,
                    "audio_response": reply_audio_base64
                 } 
          except Exception as e:
                  print(f"❌ Critical breakdown on core main processing stream: {e}")
                  raise HTTPException(status_code=500, detail="Internal AI-worker pipeline error")
          finally:
                  await file.close()
                  core_state_manager.release_processing_lock()
if __name__ == "__main__":
        uvicorn.run("main:app", port=8000, host="0.0.0.0", reload=True)
          
                          
                                        

                                
     
                  
                  
