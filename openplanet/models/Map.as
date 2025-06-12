class Map {

    uint championTime = 0;
    string uid;
    uint authorTime;
    uint goldTime;
    uint silverTime;
    uint bronzeTime;
    string name;
    string author; 

    Map(CGameCtnChallengeInfo@ MapInfo) {
        uid = MapInfo.MapUid;
        name = MapInfo.Name;
        author = MapInfo.AuthorNickName;

        authorTime = MapInfo.TMObjective_AuthorTime;
        goldTime = MapInfo.TMObjective_GoldTime;
        silverTime = MapInfo.TMObjective_SilverTime;
        bronzeTime = MapInfo.TMObjective_BronzeTime;
    }

    Json::Value@ ToJson() {
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

    bool HasChampionTime() {
        return championTime > 0;
    }

}