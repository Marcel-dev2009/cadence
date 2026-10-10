import io
from faster_whisper import WhisperModel
from services import engine
class STTService:
          def __init__(self, model_size:str = "base"):
                  self.model = WhisperModel(model_size, device="cpu", compute_type="int8")
          def transcribe_and_route(self, audio_bytes:bytes) -> dict:
                  try:
                          audio_file = io.BytesIO(audio_bytes)
                          segements, _= self.model.transcribe(audio_file, beam_size=5,
                                                                  initial_prompt="Hello. This is a clean, continuous transcription without any stuttering or filler words.")
                          user_sentence = "".join([segement.text for segement in segements]).strip()
                          if not user_sentence:
                                  print("Stt service: no user voice input detected")
                                  return {"intent":"nil", "value":""}
                          print(f"Stt service, user voice input transcribed: {user_sentence}")
                          structure_json = engine.analyze_user_voice_text(user_sentence)
                          return structure_json
                  except Exception as e:
                           print(f"Error during Speech-to-Text transcription: {e}")
                           return {"intent": "nil", "value": ""}