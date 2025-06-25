class Map 
{

    uint championTime = 0;
    string uid;
    uint authorTime;
    uint goldTime;
    uint silverTime;
    uint bronzeTime;
    string name;
    string author; 

    Map(CGameCtnChallengeInfo@ mapInfo) 
    {
        uid = mapInfo.MapUid;
        name = Text::StripFormatCodes(mapInfo.Name);
        author = mapInfo.AuthorNickName;

        authorTime = mapInfo.TMObjective_AuthorTime;
        goldTime = mapInfo.TMObjective_GoldTime;
        silverTime = mapInfo.TMObjective_SilverTime;
        bronzeTime = mapInfo.TMObjective_BronzeTime;
#if DEPENDENCY_CHAMPIONMEDALS
        uint currentChampionTime = ChampionMedals::GetCMTime();
        if (championTime != 0) 
        {
            championTime = currentChampionTime;
            DebugPrint("Champions Medal: " + championTime);
        }
#else
        DebugPrint("Champions Medals is not installed.");
#endif
    }

    Json::Value@ ToJson() 
    {
        Json::Value data = Json::Object();
        data["championTime"] = championTime;
        data["uid"] = uid;
        data["authorTime"] = authorTime;
        data["goldTime"] = goldTime;
        data["silverTime"] = silverTime;
        data["bronzeTime"] = bronzeTime;
        data["name"] = name;
        data["author"] = author;
        return data;
    }
};

bool CheckNewMap(Map@ previousMap, Map@ currentMap) 
{
    if (previousMap is null || currentMap is null) return previousMap !is currentMap;
    return previousMap.uid != currentMap.uid;
}

void SendUpdateMap(Map@ map) 
{
    if (!Setting_EnableMaps) 
    {
        DebugPrint("Sending Maps is deactivated.");
        return;
    }
    Json::Value body = map is null ? Json::Null : map.ToJson();
    PostWithRetries("trackmania/update/map", body, Setting_RetriesMap);
}

Map@ GetCurrentMap() 
{
    if (g_app is null || g_app.RootMap is null || g_app.RootMap.MapInfo.MapUid == "") 
    {
        return null;
    }
    return Map(g_app.RootMap.MapInfo);
}
