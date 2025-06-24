string TEST_ENVIRONMENT_URL = "http://localhost:3000/api/";

string GetBackend(string endpoint) 
{
    string url = Setting_TestLocal ? TEST_ENVIRONMENT_URL : Setting_Backend_Url;
    return url + endpoint;
}

int PostRequestAsync(const string &in path, const Json::Value &in data)
{
    auto request = Net::HttpRequest();
    request.Method = Net::HttpMethod::Post;
    request.Body = Json::Write(data);
    request.Headers['Content-Type'] = 'application/json';
    request.Headers['Authorization'] = Settings_Token;
    request.Url = GetBackend(path);
    request.Start();
    
    while(!request.Finished()) 
    {
        yield();
    }

    int responseCode = request.ResponseCode();
    if (responseCode == 401) 
    {
        UI::ShowNotification(Icons::ExclamationTriangle + " WheelGPT", "Your WheelGPT Token is not valid. Please request a new one on https://wheelgpt.dev", UI::HSV(.1, .8, .8));
        DebugPrint("Please set a valid token.");
    }
    return responseCode;
}

void PostWithRetries(const string &in path, const Json::Value &in data, int retries = Setting_RetriesMap)        
{

    // At least try to send the Map once
    int responseCode = PostRequestAsync(path, data);
    if (responseCode == 200) 
    {
        DebugPrint("Successfully sent data to Server.");
        return;
    }

    // Sending the Map failed, check if we have retries left
    if (retries <= 0) 
    {
        DebugPrint("No Retries left.");
        return;
    }

    DebugPrint("Problem while sending data to" + path + " (" + responseCode + "). Retries left: " + retries);
    sleep(5000);
    PostWithRetries(path, data, retries - 1);
}


