package main

import (
	"bytes"
	"fmt"
	"time"

	"github.com/mcilloni/go-openbaton/catalogue"
	log "github.com/sirupsen/logrus"
)

type handl struct {
	*log.Logger
}

// ActionForResume uses the given VNFR and VNFCInstance to return a valid
// action for resume. NoActionSpecified is returned in case no such Action exists.
func (h handl) ActionForResume(vnfr *catalogue.VirtualNetworkFunctionRecord,
	vnfcInstance *catalogue.VNFCInstance) catalogue.Action {
	return catalogue.NoActionSpecified
}

// CheckInstantiationFeasibility allows the VNFM to verify if the VNF instantiation is possible.
func (h handl) CheckInstantiationFeasibility() error {
	return nil
}

func (h handl) Configure(vnfr *catalogue.VirtualNetworkFunctionRecord) (*catalogue.VirtualNetworkFunctionRecord, error) {
	time.Sleep(3 * time.Second)

	return vnfr, nil
}

func (h handl) HandleError(vnfr *catalogue.VirtualNetworkFunctionRecord) error {
	h.WithFields(log.Fields{
		"tag":       "dummy-handl-error",
		"vnfm-name": vnfr.Name,
	}).Error("error for VNFR")

	return nil
}

func (h handl) Heal(vnfr *catalogue.VirtualNetworkFunctionRecord,
	component *catalogue.VNFCInstance, cause string) (*catalogue.VirtualNetworkFunctionRecord, error) {
	h.WithFields(log.Fields{
		"tag":       "dummy-handl-heal",
		"vnfr-name": vnfr.Name,
	}).Info("handling heal")

	return vnfr, nil
}

// Instantiate allows to create a VNF instance.
func (h handl) Instantiate(vnfr *catalogue.VirtualNetworkFunctionRecord, scripts interface{},
	vimInstances map[string][]*catalogue.VIMInstance) (*catalogue.VirtualNetworkFunctionRecord, error) {
	h.WithFields(log.Fields{
		"tag":  "dummy-handl-instantiate",
		"vnfr-name": vnfr.Name,
	}).Info("instantiating VNFR")

	vnfr.Configurations.Append(&catalogue.ConfigurationParameter{
		ConfKey: "new_key",
		Value:   "new_value",
	})

	time.Sleep(3 * time.Second)

	/*if vnfr.Name == "dummy-client" {
		h.WithFields(log.Fields{
			"tag":  "dummy-handl-instantiate",
			"vnfr-name": vnfr.Name,
		}).Info("purposefully delaying VNFR")

		time.Sleep(3 * time.Second)
	}*/
	
	return vnfr, nil
}

// Modify allows making structural changes (e.g.configuration, topology, behavior, redundancy model) to a VNF instance.
func (h handl) Modify(vnfr *catalogue.VirtualNetworkFunctionRecord,
	dependency *catalogue.VNFRecordDependency) (*catalogue.VirtualNetworkFunctionRecord, error) {

	buf := bytes.NewBufferString("")

	for key, value := range dependency.Parameters {
		buf.WriteString(fmt.Sprintf("\t%s: %v", key, value.Parameters))
	}

	h.WithFields(log.Fields{
		"tag":                        "dummy-handl-modify",
		"vnfr-hb_version":            vnfr.HbVersion,
		"vnfr-name":                  vnfr.Name,
		"vnfr-dependency":            dependency,
		"vnfr-dependency-parameters": buf.String(),
	}).Info("modifying VNFR")

	time.Sleep(3 * time.Second)
	return vnfr, nil
}

// Query allows retrieving a VNF instance state and attributes. (not implemented)
func (h handl) Query() error {
	h.WithFields(log.Fields{
		"tag": "dummy-handl-query",
	}).Warn("query invoked, not implemented")

	return nil
}

func (h handl) Resume(vnfr *catalogue.VirtualNetworkFunctionRecord,
	vnfcInstance *catalogue.VNFCInstance,
	dependency *catalogue.VNFRecordDependency) (*catalogue.VirtualNetworkFunctionRecord, error) {

	h.WithFields(log.Fields{
		"tag":       "dummy-handl-resume",
		"vnfr-name": vnfr.Name,
		"vnfr-id":   vnfr.ID,
	}).Info("resuming VNFR")

	return vnfr, nil
}

// Scale allows scaling (out / in, up / down) a VNF instance.
func (h handl) Scale(scaleInOrOut catalogue.Action,
	vnfr *catalogue.VirtualNetworkFunctionRecord,
	component catalogue.Component,
	scripts interface{},
	dependency *catalogue.VNFRecordDependency) (*catalogue.VirtualNetworkFunctionRecord, error) {

	h.WithFields(log.Fields{
		"tag":       "dummy-handl-scale",
		"vnfr-name": vnfr.Name,
		"vnfr-id":   vnfr.ID,
		"action":    scaleInOrOut,
	}).Info("scaling VNFR")

	time.Sleep(3 * time.Second)

	return vnfr, nil
}

// Start starts a VNFR.
func (h handl) Start(vnfr *catalogue.VirtualNetworkFunctionRecord) (*catalogue.VirtualNetworkFunctionRecord, error) {
	h.WithFields(log.Fields{
		"tag":       "dummy-handl-start",
		"vnfr-name": vnfr.Name,
	}).Info("starting VNFR")

	time.Sleep(3 * time.Second)
	return vnfr, nil
}

func (h handl) StartVNFCInstance(vnfr *catalogue.VirtualNetworkFunctionRecord,
	vnfcInstance *catalogue.VNFCInstance) (*catalogue.VirtualNetworkFunctionRecord, error) {

	h.WithFields(log.Fields{
		"tag":                "dummy-handl-start_vnfc_instance",
		"vnfc_instance-name": vnfcInstance.Hostname,
		"vnfc_instance-id":   vnfcInstance.ID,
	}).Info("starting VNFCInstance")

	return vnfr, nil
}

// Stop stops a previously created VNF instance.
func (h handl) Stop(vnfr *catalogue.VirtualNetworkFunctionRecord) (*catalogue.VirtualNetworkFunctionRecord, error) {
	h.WithFields(log.Fields{
		"tag":       "dummy-handl-stop",
		"vnfr-name": vnfr.Name,
	}).Info("stopping VNFR")

	//time.Sleep(3 * time.Second)
	return vnfr, nil
}

func (h handl) StopVNFCInstance(vnfr *catalogue.VirtualNetworkFunctionRecord,
	vnfcInstance *catalogue.VNFCInstance) (*catalogue.VirtualNetworkFunctionRecord, error) {

	h.WithFields(log.Fields{
		"tag":                "dummy-handl-stop_vnfc_instance",
		"vnfc_instance-name": vnfcInstance.Hostname,
		"vnfc_instance-id":   vnfcInstance.ID,
	}).Info("stopping VNFCInstance")

	return vnfr, nil
}

// Terminate allows terminating gracefully or forcefully a previously created VNF instance.
func (h handl) Terminate(vnfr *catalogue.VirtualNetworkFunctionRecord) (*catalogue.VirtualNetworkFunctionRecord, error) {
	h.WithFields(log.Fields{
		"tag":             "dummy-handl-terminate",
		"vnfr-name":       vnfr.Name,
		"vnfr-hb_version": vnfr.HbVersion,
	}).Info("terminating VNFR")

	for _, event := range vnfr.LifecycleEvents {
		if event.Event == catalogue.EventRelease {
			for _, vdu := range vnfr.VDUs {
				h.WithFields(log.Fields{
					"tag":       "dummy-handl-terminate",
					"vnfr-name": vnfr.Name,
					"vdu":       vdu,
				}).Debug("removing VDU")

				time.Sleep(3 * time.Second)
			}
		}
	}

	return vnfr, nil
}

// UpdateSoftware allows applying a minor / limited software update(e.g.patch) to a VNF instance.
func (h handl) UpdateSoftware(script *catalogue.Script,
	vnfr *catalogue.VirtualNetworkFunctionRecord) (*catalogue.VirtualNetworkFunctionRecord, error) {
	h.WithFields(log.Fields{
		"tag":       "dummy-handl-update_software",
		"script":    script,
		"vnfr-name": vnfr.Name,
		"vnfr-id":   vnfr.ID,
	}).Info("updating software for VNFR")

	time.Sleep(3 * time.Second)

	return vnfr, nil
}

// UpgradeSoftware allows deploying a new software release to a VNF instance.
func (h handl) UpgradeSoftware() error {
	h.WithFields(log.Fields{
		"tag": "dummy-handl-update_software",
	}).Warn("UpgradeSoftware called - but it's no-op")

	return nil
}

// UserData returns a string containing UserData.
func (h handl) UserData() string {
	h.WithFields(log.Fields{
		"tag": "dummy-handl-user_data",
	}).Info("returning UserData")

	return "#!/usr/bin/env sh\n"
}
