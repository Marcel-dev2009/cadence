import base64
import uvicorn
import requests
from fastapi import FastAPI, UploadFile, HTTPException, File
from fastapi.middleware import cors
import shutil
import os
import tempfile  #Clean, unique temporary file generation
from core import core_state_manager
from fastapi.concurrency import run_in_threadpool
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

@app.post("/app/voice/process")
async def process_voice_file(file: UploadFile = File(...)):
    incoming_cookies = {"session_token":"c15e00fc7e48f32a57a374be7e139c2447afd056bcf4d06d60f15b7eaca90a9e"}  #hardcoded token, this is to be sent from flutter client
    # Use a unique per-request safe filename path to prevent cross-user file overwrites
    suffix = os.path.splitext(file.filename)[1] if file.filename else ".wav"
    with tempfile.NamedTemporaryFile(suffix=suffix, delete=False) as temp_buffer:
        temp_filename = temp_buffer.name
        shutil.copyfileobj(file.file, temp_buffer)
        
    try:
        with open(temp_filename, "rb") as f:
            audio_bytes = f.read()
                   
        # Corrected run_in_threadpool syntax (pass function reference first, then arguments)
        ai_result = await run_in_threadpool(stt_service_engine.transcribe_and_route, audio_bytes)
        
        intent = ai_result.get("intent", "unknown")
        client_directive = {}
        voice_reply_text = "Voice Operation Completed Successfully"
        operation_success = True
        
        if intent == "change_theme":
            theme_value = ai_result.get("value", "light")
            client_directive = app_action_instruct.format_theme_mode(theme_value)
            voice_reply_text = f"Changing theme to {theme_value}"
            
        elif intent == "open_settings":
            client_directive = app_action_instruct.format_navigation("settings")
            voice_reply_text = "routing to settings screen"
            
        elif intent == "create_event":
            time_str = ai_result.get("event_time", "")
            name = ai_result.get("event_name", "new event")
            
            # If sync_to_golang is an async coroutine, await it natively without run_in_threadpool
            operation_success = await call_backend.sync_to_golang(event_name=name, event_time=time_str, append_cookies=incoming_cookies)
            
            if operation_success:
                voice_reply_text = f"Successfully scheduled an event for {name} to take place at {time_str}"
                client_directive = {
                    "action_type": "ui_update",
                    "target": "calendar",
                    "payload": {"refresh": True}
                }
            else:
                voice_reply_text = f"Sorry I ran into an error scheduling {name}"
                client_directive = {"action_type": "error", "payload": {"message": "Go Backend connection failure"}}
        else:
            voice_reply_text = "I'm not completely sure how to process that command yet"
            client_directive = {
                "action_type": "unknown",
                "target": "nil",
                "payload": {},
            }
            
        reply_audio_base64 = ""
        if operation_success:
            # FIX: Corrected run_in_threadpool syntax for TTS engine processing
            audio_data = await run_in_threadpool(tts_service_engine.synthesize_speech, voice_reply_text)
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
        if os.path.exists(temp_filename):
            try:
                os.remove(temp_filename)
            except Exception:
                pass        
        core_state_manager.release_processing_lock()

if __name__ == "__main__":
    uvicorn.run("main:app", port=8000, host="0.0.0.0", reload=True)
