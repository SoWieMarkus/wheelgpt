string TEST_ENVIRONMENT_URL = "http://localhost:3000/api/";

string GetBackendUrl(string endpoint)
{
    string url = Setting_TestLocal ? TEST_ENVIRONMENT_URL : Setting_Backend_Url;
    return url + endpoint;
}

// Sends one request to the WheelGPT API and returns the status code.
int SendRequestAsync(Net::HttpMethod method, const string &in path, const Json::Value &in data, bool hasBody)
{
    auto request = Net::HttpRequest();
    request.Method = method;
    if (hasBody)
    {
        request.Body = Json::Write(data);
        request.Headers['Content-Type'] = 'application/json';
    }
    request.Headers['Accept'] = 'application/json';
    request.Headers['Authorization'] = 'Bearer ' + Settings_Token;
    request.Url = GetBackendUrl(path);
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

bool IsSuccess(int responseCode)
{
    return responseCode >= 200 && responseCode < 300;
}

// Network errors, rate limits and server errors are worth another try.
// Other client errors (bad token, invalid data) will fail again.
bool ShouldRetry(int responseCode)
{
    return responseCode == 0 || responseCode == 429 || responseCode >= 500;
}

void SendWithRetries(Net::HttpMethod method, const string &in path, const Json::Value &in data, bool hasBody, int retries)
{
    int responseCode = SendRequestAsync(method, path, data, hasBody);
    if (IsSuccess(responseCode))
    {
        DebugPrint("Successfully sent data to " + path + " (" + responseCode + ").");
        return;
    }

    if (retries <= 0 || !ShouldRetry(responseCode))
    {
        DebugPrint("Giving up on " + path + " (" + responseCode + ").");
        return;
    }

    DebugPrint("Problem while sending data to " + path + " (" + responseCode + "). Retries left: " + retries);
    sleep(5000);
    SendWithRetries(method, path, data, hasBody, retries - 1);
}

void PutWithRetries(const string &in path, const Json::Value &in data, int retries = Setting_RetriesMap)
{
    SendWithRetries(Net::HttpMethod::Put, path, data, true, retries);
}

void PostWithRetries(const string &in path, const Json::Value &in data, int retries = Setting_RetriesMap)
{
    SendWithRetries(Net::HttpMethod::Post, path, data, true, retries);
}

void DeleteWithRetries(const string &in path, int retries = Setting_RetriesMap)
{
    SendWithRetries(Net::HttpMethod::Delete, path, Json::Object(), false, retries);
}
