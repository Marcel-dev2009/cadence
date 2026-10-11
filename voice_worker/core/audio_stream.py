import io
import wave
import pyaudio
class AudioStream:
          def __int__(self, chunk_size:int = 1024, sample_format:int = pyaudio.paInt16, channels:int = 1, rate:int = 16000):
                  self.chunk_size = chunk_size
                  self.sample_format = sample_format
                  self.rate = rate
                  self.channels = channels
                  self.p = pyaudio.PyAudio()
          def record_voice_input(self, silence_threshold:int = 500, max_seconds:int = 7) -> bytes:
                  stream = self.p.open(
                          format=self.sample_format,
                          channels=self.channels,
                          rate=self.rate,
                          frames_per_buffer=self.chunk_size,
                          input=True
                  )
                  print("Mic status: active, listening")
                  frames = []
                  max_chunk = int(self.rate / self.chunk_size * max_seconds)
                  for _ in range(0, max_chunk):
                          data = stream.read(self.chunk_size)
                          frames.append(data)
                  print("Mic status: idle. Recording closed")
                  stream.stop_stream()
                  stream.close()
                  wave_buffer = io.BytesIO()
                  wf = wave.open(wave_buffer, 'wb')
                  wf.setnchannels(self.channels)
                  wf.getsampwidth(self.p.get_sample_size(self.sample_format))
                  wf.setframerate(self.rate)
                  wf.writeframes(b''.join(frames))
                  wf.close()
                  return wave_buffer.getvalue()
        