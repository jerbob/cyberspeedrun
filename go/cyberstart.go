package main

import (
    "os"
    "fmt"
    "regexp"
    "net/url"
    "net/http"
    "io/ioutil"
    "encoding/json"
)

// Solution : Does a thing
type Solution struct {
    Attempt struct {
        Code string `json:"code"`
    } `json:"attempt"`
    Code string `json:"attempt.code"`
};

var csrfPattern = regexp.MustCompile(`name="csrf" value="([a-z0-9]+)"`)
var goEndpoint = "https://go.joincyberdiscovery.com"

func findCSRF(response string) (csrf string) {
    // Locate the CSRF token, given the response string.
    return csrfPattern.FindStringSubmatch(response)[1]
}

func postFlag(flag string, challenge string, csrf string) {
    // Post the flag for a challenge with the given session object.
    fmt.Printf("[+] Flag found: %s", flag)
    http.PostForm(
        fmt.Sprintf("%s/api", goEndpoint),
        url.Values{
            "action": {"flag_attempt"},
            "flag": {flag},
            "challenge": {challenge},
            "csrf": {csrf},
        },
    )
}

func solveChallenge(challenge string, solution interface{}) {
    // Solve a challenge, given a name and solution.
    response, err := http.Get(fmt.Sprintf("%s/challenges/%s", goEndpoint, challenge))
    if err != nil {
        fmt.Printf("%v", err)
        os.Exit(1)
    }
    defer response.Body.Close()
    body, _ := ioutil.ReadAll(response.Body)
    csrf := findCSRF(string(body))
    switch solution.(type) {
        default:
            flag := http.PostForm(
                fmt.Sprintf("%s/challenges/xhr/%s", goEndpoint, challenge),
                url.Values{
                    ""
                }
            )
            postFlag(flag.(string), challenge, csrf)
        case string:
            flag := solution.(string)
            postFlag(flag, challenge, csrf)
    }
}

func loadJSON(filename string) (result map[string]interface{}) {
    // Load a JSON file from a filepath.
    jsonFile, _ := os.Open(filename)
    bytes, _ := ioutil.ReadAll(jsonFile)
    json.Unmarshal([]byte(bytes), &result)
    jsonFile.Close()
    return result
}


func main() {
    solutions := loadJSON("../solutions.json")
    fmt.Printf("%v", solutions)
}
