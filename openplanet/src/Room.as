class Room 
{
    string login;
    string name; 
    int numberOfPlayers = 0;
    int maxPlayers = 0;

    Room(CGameCtnNetServerInfo@ serverInfo, int numPlayers) 
    {
        login = serverInfo.ServerLogin;
        name = Text::StripFormatCodes(serverInfo.ServerName);
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


Room@ GetCurrentRoom()
{
    auto serverInfo = cast<CGameCtnNetServerInfo>(g_network.ServerInfo);
    if (serverInfo is null || serverInfo.ServerLogin == "") 
    {
        return null;
    }
    int numPlayers = g_network.PlayerInfos.Length;
    // Subtract 1 for some reason there is always one more player than the number of players in the room
    return Room(serverInfo, numPlayers - 1);
}

bool CheckNewRoom(Room@ previousRoom, Room@ currentRoom) 
{
    if (previousRoom is null || currentRoom is null) return previousRoom !is currentRoom;
    return previousRoom.login != currentRoom.login;
}

void SendUpdateRoom(Room@ room) 
{
    if (!Setting_EnableRooms) 
    {
        DebugPrint("Sending Rooms is deactivated.");
        return;
    }

    if (room is null) 
    {
        DebugPrint("No room to send.");
        DeleteWithRetries("trackmania/room", Setting_RetriesRoom);
        return;
    }
    Json::Value body = room.ToJson();
    PostWithRetries("trackmania/room", body, Setting_RetriesRoom);
}