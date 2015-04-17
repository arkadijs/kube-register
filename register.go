package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Meta struct {
	Name string `json:"name,omitempty"`
}

type Minion struct {
	APIVersion string `json:"apiVersion,omitempty"`
	Kind       string `json:"kind,omitempty"`
	ID         string `json:"id,omitempty"`
	Metadata   Meta   `json:"metadata,omitempty"`
	HostIP     string `json:"hostIP,omitempty"`
}

type MinionResp struct {
	Reason string `json:"reason,omitempty"`
}

func register(endpoint, addr string) error {
	m := &Minion{
		APIVersion: "v1beta1",
		Kind:       "Minion",
		ID:         addr,
		Metadata:   Meta{Name: addr},
		HostIP:     addr,
	}
	mr := &MinionResp{}
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%s/api/v1beta1/minions", endpoint)
	res, err := http.Post(url, "application/json", bytes.NewReader(data))
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode == 201 || res.StatusCode == 200 {
		log.Printf("registered machine: %s\n", addr)
		return nil
	}
	body, err := ioutil.ReadAll(res.Body)
	json.Unmarshal([]byte(body), &mr)
	if res.StatusCode == 409 && mr.Reason == "AlreadyExists" {
		return nil
	}
	reason := ""
	if err == nil {
		reason = ": " + string(body)
	}
	return errors.New(fmt.Sprintf("error registering %s: %#v\n%s", addr, res, reason))
}
