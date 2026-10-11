from services import stt_service_engine
import io
class WakeWordDetector:
          def __init__(self, target_wake_word:str = "Hey, Cadence"):
                  self.target_wake_word = target_wake_word.lower().strip()
          def inspect_for_activation(self, rapid_audio_chunk:bytes) -> bytes:
                  try:
                          audio_stream = io.BytesIO(rapid_audio_chunk)
                          segements,_ = stt_service_engine.model.transcribe(audio_stream, beam_size=5)
                          spoken_text = "".join([s.text for s in segements]).lower().strip
                          if self.target_wake_word in spoken_text:
                                  print(f"Trigger word detected in spoken text: {spoken_text}")
                                  return True
                          return False
                  except Exception as e:
                          print(f"Trigger word not found in spoken text:{e}")
                          return False