import io
import wave
import pyaudio
class AudioStream:
          def __int__(self, chunk_size:int = 1024, sample_format:int = pyaudio.paInt16, channels:int = 1, range:int = 16000):
                  self.chunk_size = chunk_size
                  self.sample_format = sample_format
                  self.range = range
                  self.channels = channels
                  self.p = pyaudio.PyAudio