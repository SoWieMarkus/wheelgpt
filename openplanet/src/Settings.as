[Setting category="Info" name="Token" password description="Get your token on https://wheelgpt.dev/"]
string Settings_Token = "";

[Setting category="Info" name="Send PBs to WheelGPT" description="Enables sending your personal bests to WheelGPT. This triggers an event that publishes your PB in your Twitch chat."]
bool Setting_EnablePBs = true;

[Setting category="Info" name="Send Map to WheelGPT" description="Enables sending your current map to WheelGPT. Your Twitch Viewers can use !map to get your current map." ]
bool Setting_EnableMaps = true;

[Setting category="Info" name="Send Room to WheelGPT" description="Enables sending your current room to WheelGPT. Your Twitch Viewers can use !room to get your current room." ]
bool Setting_EnableRooms = true;

[Setting category="Dev" name="Debug Mode" description="Enables additional debug output in the console."]
bool Setting_DebugMode = false;

[Setting category="Dev" name="Test Local" description="If enabled, the backend URL will be set to http://localhost:3000/api/."]
bool Setting_TestLocal = false;

[Setting category="Dev" name="Retries for Map" min=0 max=10 description="The number of retries for sending the map to the backend. If set to 0, no retries will be made."]
int Setting_RetriesMap = 5;

[Setting category="Dev" name="Retries for PB" min=0 max=10 description="The number of retries for sending the personal best to the backend. If set to 0, no retries will be made."]
int Setting_RetriesPB = 5;

[Setting category="Dev" name="Retries for Room" min=0 max=10 description="The number of retries for sending the room to the backend. If set to 0, no retries will be made."]
int Setting_RetriesRoom = 5;

[Setting category="Dev" name="Backend" description="The URL of the backend API. This is only relevant if you are hosting your own WheelGPT instance."]
string Setting_Backend_Url = "https://wheelgpt.dev/api/";