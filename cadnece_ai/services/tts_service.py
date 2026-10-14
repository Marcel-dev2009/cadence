import pyttsx3
import uuid
import os
class TTSService:
          def __init__(self, voice_speed:int = 200):
                  self.voice_speed = voice_speed
          def synthesize_speech(self, text_to_speak:str) -> bytes:
                  # we need to generate a temporary file name to avoid concurrent overwrites
                  tempt_filename = f"{uuid.uuid4().hex}.wav" 
                  try:
                          tts_engine = pyttsx3.init()
                          tts_engine.setProperty('rate', self.voice_speed)
                          tts_engine.save_to_file(text_to_speak, tempt_filename)
                          tts_engine.runAndWait()
                          with open(tempt_filename, "rb") as audio_file:
                                  audio_bytes = audio_file.read()
                          return audio_bytes
                  except Exception as e:
                          print(f"error sythesizing speech:{e}")
                          return b""
                  finally:
                          if os.path.exists(tempt_filename):
                                  os.remove(tempt_filename)
                          
                          
                          