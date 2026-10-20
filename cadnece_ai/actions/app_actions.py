class AppActions:
          @staticmethod
          def format_theme_mode (theme_value:str) -> dict:
                  mode = theme_value.strip().lower() if theme_value else "light"
                  if not mode in ["dark", "light"]:
                          mode = "light"
                  return{
                          "action_type": "ui_update",
                          "target":"theme",
                          "payload":{
                             "mode": mode     
                          }
                  }
          @staticmethod
          def format_navigation(screen_name:str) -> dict:
                  route_map = {
                        "settings": "/settings" 
                  }
                  target_route = route_map.get(screen_name.strip().lower(), "/home")
                  return {
                    "action_type": "navigation",
                    "target": "route",
                    "payload": {
                    "path": target_route
                    }
        }
