// Time to wait between each check in milliseconds
const int CHECK_DELAY = 1000;
const int DEFAULT_BEST_TIME = -1;

void Main() 
{

    CTrackMania@ app = cast<CTrackMania>(GetApp());
    CTrackManiaNetwork@ network = cast<CTrackManiaNetwork>(app.Network);

    Map@ previousMap = null;
    int previousBestTime = DEFAULT_BEST_TIME;

    while(true) 
    {
        sleep(CHECK_DELAY);

        auto currentMap = GetCurrentMap(app);

        if (CheckNewMap(previousMap, currentMap))
        {
            @previousMap = Map(currentMap.MapInfo);
            DebugPrint("New current map: " + previousMap is null ? "null" : previousMap.name);
            SendUpdateMap(previousMap);
            previousBestTime = DEFAULT_BEST_TIME; // Reset the previous best time when changing maps
        }


        int currentPersonalBestTime = GetCurrentPersonalBest(network);
        // Ensure we are currently on a map to avoid unexpected behavior (e.g. during map loading)
        if (CheckNewPersonalBest(previousBestTime, currentPersonalBestTime) && currentMap !is null) {
            DebugPrint("New personal best detected: " + currentPersonalBestTime);
            SendUpdatePersonalBest(personalBestTime);
        }
        // Update the previous best time to the current one for the next iteration
        // This is always done to initialize the previous best time when entering a new map
        previousBestTime = currentPersonalBestTime;

    }

}

