from .audio_stream import AudioStream
from .state_manager import CoreStateManager
from .wake_word import WakeWordDetector
audio_stream_manager = AudioStream()
core_state_manager = CoreStateManager()
wake_word_detector = WakeWordDetector()
__all__ = ["audio_stream_manager", "core_state_manager", "wake_word_detector"]