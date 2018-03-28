package main

import (
  "fmt"
  "log"
  "encoding/json"
  "net/http"
  "os"
  "os/exec"
)

const (
  CONN_HOST="0.0.0.0"
  CONN_PORT="7890"
  LISTENER=CONN_HOST+":"+CONN_PORT
  CONN_TYPE="tcp"
)

type(
  //code shamelessly copied from github.com/bakins/alertmanager-webhook-example
  // to read in prometheus alertmanager json requests and spit out a string

  // HookMessage is the message we receive from Alertmanager
  HookMessage struct {
    Version           string            `json:"version"`
    GroupKey          string            `json:"groupKey"`
    Status            string            `json:"status"`
    Receiver          string            `json:"receiver"`
    GroupLabels       map[string]string `json:"groupLabels"`
    CommonLabels      map[string]string `json:"commonLabels"`
    CommonAnnotations map[string]string `json:"commonAnnotations"`
    ExternalURL       string            `json:"externalURL"`
    Alerts            []Alert           `json:"alerts"`
  }

// Alert is a single alert.
  Alert struct {
    Labels      map[string]string `json:"labels"`
    Annotations map[string]string `json:"annotations"`
    StartsAt    string            `json:"startsAt,omitempty"`
    EndsAt      string            `json:"EndsAt,omitempty"`
  }
)

func main() {
  fmt.Println("About to listen on " + LISTENER)

  http.HandleFunc("/alerts", postHandler)  
  log.Fatal(http.ListenAndServe(LISTENER,nil))
}

// handle incoming requests
func postHandler(w http.ResponseWriter, r *http.Request) {
  dec := json.NewDecoder(r.Body)
  defer r.Body.Close()

  var m HookMessage
  if err := dec.Decode(&m); err != nil {
    log.Printf("error decoding message: %v", err)
    http.Error(w, "invalid request body", 400)
    return
  }

  for _, element := range(m.Alerts) {
  //we need to break the map into key/val pairs, and then
  // check to see if we have summary/description
    summary := element.Annotations["summary"]
    description := element.Annotations["description"]

    cmd := "./stupid_echo"  //i made a shell script called stupid_echo also in this repo
    args := []string{summary, ":", description}
    if err := exec.Command(cmd, args...).Run(); err != nil{
      fmt.Fprintln(os.Stderr,err)
      os.Exit(1)
    }
  }

}
