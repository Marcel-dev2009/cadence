from config import settings
import json
from huggingface_hub import InferenceClient
class IntentEngine:
          def __int__(self):
                  self.client = InferenceClient(
                          model="meta-llama/Meta-Llama-8B-Instruct",
                          token=settings.AI_API_KEY
                  )
          def analyze_user_voice_text(self, user_sentence:str) -> dict:
                  system_prompt = (
                          "You are a strict backend assistant for a scheduling app named Cadence.\n"
                          "Your job is to read the user's input sentence and classify their intent.\n"
                          "You must respond ONLY with a raw JSON block. Do not include a markdown code blocks.\n\n"
                          "Available Actions:\n"
                          "1. 'create_event' - Used if they want to schedule, add, or set an event.\n"
                          "2. 'change_theme' - Used if they mention app theme, dark mode or light mode.\n"
                          "3. 'open_settings'- Used if they mention settings, profile, or configuration pages.\n"
                          "4.  'todo_list' - Used if they mention checklist, todo-list bucket-list"
                          "5. 'unknown' - Used for anything else.\n\n"
                          "JSON Output format example:\n"
                          '{"intent": "change_theme", "value": "dark"}\n'
                          '{"intent": "create_event", "value": ""}'
                  )
                  try:
                          response = self.client.text_generation(
                                  prompt=f"<|system|>\n{system_prompt}\n<|user|>\n{user_sentence}\n<|assistant|>\n>",
                                  max_new_tokens=100,
                                  temperature=0.1
                          )
                          clean_response = response.strip()
                          return json.loads(clean_response)
                  except Exception as e:
                          print(f"❌ Error communicating with Llama 3 engine: {e}")
                          return {"intent": "unknown", "value": ""}
