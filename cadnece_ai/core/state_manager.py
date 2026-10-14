import threading
class CoreStateManager:
          def __init__(self):
                  self._is_processing = False
                  self._lock = threading.Lock()
          def try_processing_lock(self) -> bool:
                  with self._lock:
                          if self._is_processing: # when we accept input something should set it to be true 
                                    print("State Warning: System is busy. Ignoring duplicate voice inputs.")
                                    return False
                          self._is_processing = True
                          return True
          def release_processing_lock(self) -> bool:
                  with self._lock:
                          self._is_processing = False
                              
          