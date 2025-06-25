// Time to wait between each check in milliseconds
const int CHECK_DELAY = 1000;
const int DEFAULT_BEST_TIME = -1;
const uint TIMEOUT_ROOM_UPDATE = 60000;

CTrackMania@ g_app;
CTrackManiaNetwork@ g_network;

void Main() 
{

    @g_app = cast<CTrackMania>(GetApp());
    @g_network = cast<CTrackManiaNetwork>(g_app.Network);

    Map@ previousMap = null;
    int previousBestTime = DEFAULT_BEST_TIME;
    Room@ previousRoom = null;
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
            previousBestTime = DEFAULT_BEST_TIME; // Reset the previous best time when changing maps
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
    }

}

