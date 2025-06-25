void CheckNewPersonalBest(int previousBestTime, int currentBestTime) 
{
    // When entering a new map, the previous best time is reset to DEFAULT_BEST_TIME
    // This means we would trigger a new personal best update on map join,
    // since the current best time on this map is not equal to DEFAULT_BEST_TIME.

    // This ensures that we only send updates when the personal best actually changes
    // and not when the player is just joining a new map
    // e.g. previous = DEFAULT_BEST_TIME, current = 1000
    // don't trigger an update 
    // e.g. previous = 1000, current = 1000
    // don't trigger an update
    // e.g. previous = 1000, current = 900
    // trigger an update
    return currentBestTime != DEFAULT_BEST_TIME && previousBestTime != currentBestTime;
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
    PostWithRetries("trackmania/update/pb", body, Setting_RetriesPB);
}

int GetCurrentPersonalBest(CTrackManiaNetwork@ network) 
{
    if (network.ClientManiaAppPlayground is null) 
    {
        return DEFAULT_BEST_TIME;
    }

    auto userManager = network.ClientManiaAppPlayground.UserMgr;
    MwId userId = uint(-1);
    if (userManager.Users.Length > 0) 
    {
        userId = userManager.Users[0].Id;
    }

    auto scoreManager = network.ClientManiaAppPlayground.ScoreMgr;
    return scoreManager.Map_GetRecord_v2(userId, g_app.RootMap.MapInfo.MapUid, "PersonalBest", "", "TimeAttack", "");
}