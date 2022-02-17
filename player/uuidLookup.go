package player

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

const ApiUrl = "https://api.thom.club"

type UuidLookup struct {
	Uuids map[string]string `json:"uuids"`
}

type UuidPostReq struct {
	Uuids []string `json:"uuids"`
}

type UsernameLookup struct {
	Usernames map[string]string `json:"usernames"`
}

type UsernamePostReq struct {
	Usernames []string `json:"usernames"`
}

func LookupUUIDs(uuids []string) *UuidLookup {
	postData, err := json.Marshal(UuidPostReq{Uuids: uuids})
	if err != nil {
		log.Fatal(err)
	}
	resp, err := http.Post(ApiUrl+"/bulk_uuids", "application/json", bytes.NewReader(postData))
	object := &UuidLookup{}
	if err != nil {
		log.Fatal(err)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(data, object)
	return object
}

func LookupSingularUUID(uuid string) string {
	data := LookupUUIDs([]string{uuid})
	return data.Uuids[uuid]
}

func LookupUsernames(usernames []string) *UsernameLookup {
	postData, err := json.Marshal(UsernamePostReq{Usernames: usernames})
	if err != nil {
		log.Fatal(err)
	}
	resp, err := http.Post(ApiUrl+"/bulk_usernames", "application/json", bytes.NewReader(postData))
	object := &UsernameLookup{}
	if err != nil {
		log.Fatal(err)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(data, object)
	return object
}
