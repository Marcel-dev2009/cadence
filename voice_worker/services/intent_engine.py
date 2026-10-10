import json
import dateparser
from huggingface_hub import InferenceClient
from config import settings

class IntentEngine:
    def __init__(self):
        self.client = InferenceClient(
            model="meta-llama/Meta-Llama-3-8B-Instruct",
            token=settings.AI_API_KEY
        )

    def analyze_user_voice_text(self, user_sentence: str) -> dict:
        system_prompt = (
            "You are a strict backend assistant for a scheduling app named Cadence.\n"
            "Your job is to read the user's input sentence and classify their intent.\n"
            "You must respond ONLY with a raw JSON block. Do not include markdown code blocks.\n\n"
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
            response = self.client.text_generation(
                prompt=f"<|system|>\n{system_prompt}\n<|user|>\n{user_sentence}\n<|assistant|>\n",
                max_new_tokens=150,  # Bumped up slightly to safely accommodate list entries
                temperature=0.1
            )
            
            clean_response = response.strip()
            if clean_response.startswith("```"):
                clean_response = clean_response.strip("`").replace("json", "", 1).strip()
            
            payload = json.loads(clean_response)
            intent = payload.get("intent")
            
            if intent == "create_event":
                if "event_time" in payload and payload["event_time"]:
                    parsed_date = dateparser.parse(payload["event_time"])
                    if parsed_date:
                        payload["event_time"] = parsed_date.strftime("%Y-%m-%d")
                    else:
                        payload["event_time"] = None  # Fallback if text is unparsable

            elif intent == "change_theme":
                # Ensure a safe fallback string value if missing
                if "value" not in payload:
                    payload["value"] = "light"

            elif intent == "open_settings":
                # Match your required structural routing layout
                payload["route"] = "settings"

            return payload                 
            
        except Exception as e:
            print(f"❌ Error communicating with Llama 3 engine: {e}")
            return {"intent": "unknown", "value": ""}
