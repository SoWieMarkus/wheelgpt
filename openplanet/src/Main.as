// Time to wait between each check in milliseconds
const int CHECK_DELAY = 1000;
// Used whenever we are not on a map
const int DEFAULT_BEST_TIME_NO_MAP = -2;
// Used when on a map with no best time
const int DEFAULT_NO_BEST_TIME = -1;
const uint TIMEOUT_ROOM_UPDATE = 60000;

CTrackMania@ g_app;
CTrackManiaNetwork@ g_network;

void Main() 
{

    @g_app = cast<CTrackMania>(GetApp());
    @g_network = cast<CTrackManiaNetwork>(g_app.Network);

    Map@ previousMap = null;
    int previousBestTime = DEFAULT_BEST_TIME_NO_MAP;
    Room@ previousRoom = null;

    // If the previous room update is more then 2 * TIMEOUT_ROOM_UPDATE ago, we set the timestamp to now - TIMEOUT_ROOM_UPDATE
    // This is to ensure we will never get an integer overflow when no room updates are sent for a long time
    uint previousRoomUpdateTime = Time::Now;

    while(true) 
    {
        sleep(CHECK_DELAY);

        Map@ currentMap = GetCurrentMap();
        if (CheckNewMap(previousMap, currentMap))
        {
            @previousMap = currentMap;
            DebugPrint("New current map: " + (previousMap is null ? "null" : previousMap.name));
            SendUpdateMap(previousMap);
            previousBestTime = DEFAULT_BEST_TIME_NO_MAP; // Reset the previous best time when changing maps
        }

        int currentPersonalBestTime = GetCurrentPersonalBest();
        // Ensure we are currently on a map to avoid unexpected behavior (e.g. during map loading)
        if (CheckNewPersonalBest(previousBestTime, currentPersonalBestTime) && currentMap !is null) 
        {
            DebugPrint("New personal best detected: " + currentPersonalBestTime);
            SendUpdatePersonalBest(currentPersonalBestTime);
        }
        // Update the previous best time to the current one for the next iteration
        // This is always done to initialize the previous best time when entering a new map
        DebugPrint("Current personal best: " + currentPersonalBestTime + ", Previous best: " + previousBestTime);
        previousBestTime = currentPersonalBestTime;

        Room@ currentRoom = GetCurrentRoom();
        uint delta = Time::Now - previousRoomUpdateTime;
        DebugPrint("Time delta "  + delta);
        if (CheckNewRoom(previousRoom, currentRoom) && delta > TIMEOUT_ROOM_UPDATE)
        {
            @previousRoom = currentRoom;
            DebugPrint("New current room: " + (currentRoom is null ? "null" : currentRoom.name));
            SendUpdateRoom(previousRoom);
            previousRoomUpdateTime = Time::Now;
        }

        if (Time::Now - previousRoomUpdateTime > 2 * TIMEOUT_ROOM_UPDATE) 
        {
            // Reset the previous room update time to now - TIMEOUT_ROOM_UPDATE
            // This ensures that we will not miss any room updates that might come in after a long period of inactivity
            // and avoids integer overflow issues.
            previousRoomUpdateTime = Time::Now - TIMEOUT_ROOM_UPDATE;
        }        
    }

}

