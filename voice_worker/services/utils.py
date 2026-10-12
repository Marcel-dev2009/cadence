import os
import json
import re
from typing import Protocol, Optional
import tempfile
import dateparser
from openai import OpenAI
from faster_whisper import WhisperModel
from huggingface_hub import InferenceClient
from config import settings

class Engine(Protocol):
    # FIX 1: Updated signature to accept the text string
    def analyze_user_voice_text(self, user_sentence: str) -> dict:
        ...

class IntentEngine:
    def __init__(self):
        self.client = OpenAI(
            base_url="https://api.groq.com/openai/v1",
            api_key=settings.AI_API_KEY
        )

    def analyze_user_voice_text(self, user_sentence: str) -> dict:
        system_prompt = (
            "You are a strict backend assistant for a scheduling app named Cadence.\n"
            "Your job is to read the user's input sentence and classify their intent.\n"
            "You must respond ONLY with a raw JSON object. Do not include markdown code blocks.\n\n"
            "Available Actions:\n"
            "1. 'create_event' - Requires extracting 'event_name' and 'event_time'.\n"
            "2. 'change_theme' - Requires extracting 'value' (e.g., 'dark' or 'light').\n"
            "3. 'open_settings' - Used if they mention settings, profile, or configuration pages.\n"
            "4. 'unknown' - Used for anything else.\n\n"
            "JSON Output format constraints:\n"
            '{"intent": "change_theme", "value": "dark"}\n'
            '{"intent": "create_event", "event_name": "code review", "event_time": "3pm thursday"}\n'
            '{"intent": "open_settings", "route": "settings"}\n'
        )
        
        try:
            response = self.client.chat.completions.create(
                model="openai/gpt-oss-20b",
                messages=[
                    {"role": "system", "content": system_prompt},
                    {"role": "user", "content": user_sentence}
                ],
                max_tokens=150,
                temperature=0.1
            )
            
            clean_response = response.choices[0].message.content.strip()
            
            # FIX 2: Bulletproof JSON Extraction (Handles markdown blocks or conversational text safely)
            json_match = re.search(r'\{.*\}', clean_response, re.DOTALL)
            if json_match:
                clean_response = json_match.group(0)
            
            payload = json.loads(clean_response)
            intent = payload.get("intent")
            
            if intent == "create_event":
                if "event_time" in payload and payload["event_time"]:
                    parsed_date = dateparser.parse(payload["event_time"])
                    if parsed_date:
                        payload["event_time"] = parsed_date.strftime("%Y-%m-%d")
                    else:
                        payload["event_time"] = None

            elif intent == "change_theme":
                if "value" not in payload:
                    payload["value"] = "light"

            elif intent == "open_settings":
                payload["route"] = "settings"

            return payload                 
            
        except Exception as e:
            print(f"❌ Error communicating with Groq Llama engine: {e}")
            return {"intent": "unknown", "value": ""}

    
class STTService:
    def __init__(self, model_size: str = "base", engine: Optional[Engine] = None):
        self.engine = engine if engine is not None else IntentEngine()
        self.model = WhisperModel(model_size, device="cpu", compute_type="int8")
        
    def transcribe_and_route(self, audio_bytes: bytes) -> dict:
        # FIX 3: Safe, unique temporary file generation to prevent multi-user race conditions
        with tempfile.NamedTemporaryFile(suffix=".wav", delete=False) as temp_file:
            temp_filename = temp_file.name
            temp_file.write(audio_bytes)
            
        try:
            segments, _ = self.model.transcribe(
                temp_filename, 
                beam_size=5, 
                initial_prompt="Hello, This is a clean, continuous transcription without any stuttering or filler words."
            )
            user_sentence = "".join([s.text for s in segments]).strip()
            
            if not user_sentence:
                print("Stt service: no user voice input detected")
                return {"intent": "nil", "value": ""}
                
            print(f"Stt service, user voice input transcribed: {user_sentence}")
            structure_json = self.engine.analyze_user_voice_text(user_sentence)
            return structure_json
            
        except Exception as e:
            print(f"Error during Speech-to-Text transcription: {e}")
            return {"intent": "nil", "value": ""}
        finally:
            if os.path.exists(temp_filename):
                try:
                    os.remove(temp_filename)
                except Exception:
                    pass
