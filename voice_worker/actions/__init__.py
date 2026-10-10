from .app_actions import AppActions
from .backend_actions import BackendActions
app_action_instruct = AppActions()
call_backend = BackendActions()
__all__ = ["app_action_instruct", "call_backend"]