package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/sirupsen/logrus"
)

func (sds *SdServer) handleGetData(w http.ResponseWriter, r *http.Request) {
	request := &SdsRequest{}
	err := parseRequest(r, request)
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	data, err := GetData(*request)
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	marshaledData, err := json.Marshal(data)
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(marshaledData)
	return
}

func parseRequest(r *http.Request, request *SdsRequest) error {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logrus.WithError(err).Error("Failed to parse request body")
		return err
	}
	err = json.Unmarshal(reqBody, request)
	if err != nil {
		logrus.WithError(err).Error("Failed to unmarshal request as ")
		return err
	}
	return nil
}
