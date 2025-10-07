bool CheckNewPersonalBest(int previousBestTime, int currentBestTime) 
{
    // When entering a new map, the previous best time is reset to DEFAULT_BEST_TIME_NO_MAP
    // This means we would trigger a new personal best update on map join,
    // since the current best time on this map is not equal to DEFAULT_BEST_TIME_NO_MAP.

    // This ensures that we only send updates when the personal best actually changes
    // and not when the player is just joining a new map
    // e.g. previous = DEFAULT_BEST_TIME_NO_MAP, current = 1000
    // don't trigger an update because we just loaded the map
    // e.g. previous = 1000, current = 1000
    // don't trigger an update because the personal best didn't change
    // e.g. previous = 1000, current = 900
    // if no personal best was set before, we still want to
    // e.g. previous = DEFAULT_NO_TIME, current = 900
    // trigger an update since we improved
    return previousBestTime != DEFAULT_BEST_TIME_NO_MAP && previousBestTime != currentBestTime;
}

void SendUpdatePersonalBest(int time) 
{
    if (!Setting_EnablePBs) 
    {
        DebugPrint("Sending PBs is deactivated.");
        return;
    }
    Json::Value body = Json::Object();
    body["time"] = time;
    PostWithRetries("trackmania/pb", body, Setting_RetriesPB);
}

int GetCurrentPersonalBest() 
{
    if (g_network.ClientManiaAppPlayground is null || 
        g_app.RootMap is null || 
        g_app.RootMap.MapInfo.MapUid == "") 
    {
        return DEFAULT_BEST_TIME_NO_MAP;
    }

    auto userManager = g_network.ClientManiaAppPlayground.UserMgr;
    MwId userId = uint(-1);
    if (userManager.Users.Length == 0) 
    {
        return DEFAULT_BEST_TIME_NO_MAP;
    }
    userId = userManager.Users[0].Id;

    auto scoreManager = g_network.ClientManiaAppPlayground.ScoreMgr;
    return scoreManager.Map_GetRecord_v2(userId, g_app.RootMap.MapInfo.MapUid, "PersonalBest", "", "TimeAttack", "");
}