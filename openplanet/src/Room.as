class Room 
{
    string login;
    string name; 
    int numberOfPlayers = 0;
    int maxPlayers = 0;

    Room(CGameCtnNetServerInfo@ serverInfo, int numPlayers) 
    {
        login = serverInfo.ServerLogin;
		name = StripFormatCodes(serverInfo.ServerName);

		numberOfPlayers = numPlayers;
        maxPlayers = serverInfo.MaxPlayerCount;
    }

    Json::Value@ ToJson() 
    {
        Json::Value data = Json::Object();
        data["login"] = login;
        data["name"] = name;
        data["numberOfPlayers"] = numberOfPlayers;
        data["maxPlayers"] = maxPlayers;
        return data;
    }
};


Room@ GetCurrentRoom(CTrackManiaNetwork@ network)
{
    auto serverInfo = cast<CGameCtnNetServerInfo>(network.ServerInfo);
    if (serverInfo is null || serverInfo.ServerLogin == "") 
    {
        return null;
    }
    int numPlayers = network.PlayerInfos.Length;
    return Room(serverInfo, numPlayers);
}

void SendUpdateRoom(Room@ room) 
{
    if (!Setting_EnableRooms) 
    {
        DebugPrint("Sending Rooms is deactivated.");
        return;
    }
    Json::Value body = room is null ? Json::Null() : room.ToJson();
    PostWithRetries("trackmania/update/room", body, Setting_RetriesRoom);
}