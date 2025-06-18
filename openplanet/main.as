[Setting category="Info" name="Token" password description="Get your token on https://wheelgpt.dev/"]
string token = "";

[Setting category="Info" name="Send PBs to Server"]
bool enablePBs = true;

[Setting category="Info" name="Send Map Details to Server"]
bool enableMaps = true;

[Setting category="Dev" name="Debug"]
bool debug = false;

[Setting category="Dev" name="TestLocal"]
bool testLocal = false;

[Setting category="Dev" name="Backend"]
string backend = "https://wheelgpt.dev/api/";

string localBackend = "http://localhost:3000/api/";

string getBackend() {
    if (testLocal) return localBackend;
    return backend;
}

Net::HttpRequest@ PostRequestAsync(const string &in url, const Json::Value &in data){
    auto request = Net::HttpRequest();
    request.Method = Net::HttpMethod::Post;
    request.Body = Json::Write(data);
    request.Headers['Content-Type'] = 'application/json';
    request.Headers['Authorization'] = token;
    request.Url = url;
    request.Start();
    
    while(!request.Finished()) {
        yield();
    }
    return request;
}

const int delay = 1000;

Map@ currentMap = null;

int previousBestTime = -1;

int retriesMap = 5;

int retriesPB = 5;


void Main() {

    auto app = cast<CTrackMania>(GetApp());
    auto network = cast<CTrackManiaNetwork>(app.Network);

    while(true) {
        sleep(delay);

        auto rootMap = app.RootMap;
        if (rootMap is null || rootMap.MapInfo.MapUid == "") {
            retriesMap = 5;  
            if (currentMap !is null) {
                @currentMap = null;
                SendUpdateMap();
            }
                   
            continue;
        }

        bool newMap = false;
        if (currentMap is null || currentMap.uid != rootMap.MapInfo.MapUid) {
            @currentMap = Map(rootMap.MapInfo);
            debugPrint("New map " + currentMap.name);
            SendUpdateMap();
            newMap = true;
        }

#if DEPENDENCY_CHAMPIONMEDALS

        uint championTime = ChampionMedals::GetCMTime();
        if (!currentMap.HasChampionTime() && championTime != 0) {
            currentMap.championTime = championTime;
            debugPrint("Champions Medal: " + championTime);
            SendUpdateMap();
        }
#else
        debugPrint("Champions Medals is not installed.");
#endif

        int personalBestTime = -1;
        if (network.ClientManiaAppPlayground is null) {
            continue;
        }

        auto userManager = network.ClientManiaAppPlayground.UserMgr;
        MwId userId = uint(-1);
        if (userManager.Users.Length > 0) {
            userId = userManager.Users[0].Id;
        }

        auto scoreManager = network.ClientManiaAppPlayground.ScoreMgr;
        personalBestTime = scoreManager.Map_GetRecord_v2(userId, rootMap.MapInfo.MapUid, "PersonalBest", "", "TimeAttack", "");

        if (newMap) {
            debugPrint("Set new pb time from storage");
            previousBestTime = personalBestTime;
        }

        if (personalBestTime != -1 && previousBestTime != personalBestTime) {
            debugPrint("new best time");
            previousBestTime = personalBestTime;
            SendUpdatePersonalBest(personalBestTime);
        }

    }

}

void debugPrint(string message) {
    if (debug) print(message);
}

void SendUpdateMap() {
    if (!enableMaps) {
        debugPrint("Sending maps is deactivated.");
        return;
    }

    if (retriesMap == 0) {
        retriesMap = 5;
        debugPrint("No Retries left.");
        return;
    }

    string endpoint = getBackend() + "trackmania/update/map";
    Json::Value body = Json::Object();
    if (currentMap is null) {
        body["map"] = null;
    } else {
        body["map"] = currentMap.ToJson();
    }
    

    debugPrint("Sending map to " + endpoint);
    string mapName = "null";
    if (currentMap !is null) {
        mapName = currentMap.name;
    }
    debugPrint("Map " + mapName);

    auto request = PostRequestAsync(endpoint, body);
    int resultCode = request.ResponseCode();

    if (resultCode == 200) {
        debugPrint("Successfully sent Map to Server.");
        retriesMap = 5;
        return;
    }

    debugPrint("Problem while sending data. Retries left: " + retriesMap);
    if (resultCode == 401) {
        debugPrint("Please set a valid token.");
    }
    debugPrint("Response: "+resultCode);

    retriesMap = retriesMap - 1;
    sleep(5000);
    SendUpdateMap();
}

void SendUpdatePersonalBest(int time) {

    if (!enablePBs) {
        debugPrint("Sending PBs is deactivated.");
        return;
    }
    
    if (retriesPB == 0) {
        retriesPB = 5;
        debugPrint("No Retries left.");
        return;
    }

    string endpoint = getBackend() + "trackmania/update/pb";
    Json::Value body = Json::Object();
    body["time"] = time;

    debugPrint("Sending pb to " + endpoint);

    auto request = PostRequestAsync(endpoint, body);
    int resultCode = request.ResponseCode();

    if (resultCode == 200) {
        debugPrint("Successfully sent PB to Server.");
        retriesPB = 5;
        return;
    }

    retriesPB = retriesPB - 1;

    debugPrint("Problem while sending data. Retries left: " + retriesPB);
    debugPrint("Response: "+resultCode);
    sleep(5000);
    SendUpdatePersonalBest(time);
}