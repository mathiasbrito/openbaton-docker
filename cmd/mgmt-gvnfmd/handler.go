package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/mcilloni/openbaton-docker/mgmt"
	"github.com/openbaton/go-openbaton/catalogue"
	log "github.com/sirupsen/logrus"
)

type handl struct {
	*log.Logger
	acc mgmt.AMQPChannelAccessor
}

// ActionForResume uses the given VNFR and VNFCInstance to return a valid
// action for resume. NoActionSpecified is returned in case no such Action exists.
func (h *handl) ActionForResume(vnfr *catalogue.VirtualNetworkFunctionRecord,
	vnfcInstance *catalogue.VNFCInstance) catalogue.Action {
	return catalogue.NoActionSpecified
}

// CheckInstantiationFeasibility allows the VNFM to verify if the VNF instantiation is possible.
func (h *handl) CheckInstantiationFeasibility() error {
	return nil
}

func (h *handl) Configure(vnfr *catalogue.VirtualNetworkFunctionRecord) (*catalogue.VirtualNetworkFunctionRecord, error) {
	time.Sleep(3 * time.Second)

	return vnfr, nil
}

func (h *handl) HandleError(vnfr *catalogue.VirtualNetworkFunctionRecord) error {
	h.WithFields(log.Fields{
		"tag":       "mgmt-gvnfm-handl-error",
		"vnfm-name": vnfr.Name,
	}).Error("error for VNFR")

	return nil
}

func (h *handl) Heal(vnfr *catalogue.VirtualNetworkFunctionRecord,
	component *catalogue.VNFCInstance, cause string) (*catalogue.VirtualNetworkFunctionRecord, error) {
	h.WithFields(log.Fields{
		"tag":       "mgmt-gvnfm-handl-heal",
		"vnfr-name": vnfr.Name,
	}).Info("handling heal")

	return vnfr, nil
}

// Instantiate allows to create a VNF instance.
func (h *handl) Instantiate(vnfr *catalogue.VirtualNetworkFunctionRecord, scripts interface{},
	vimInstances map[string][]*catalogue.VIMInstance) (*catalogue.VirtualNetworkFunctionRecord, error) {
	h.WithFields(log.Fields{
		"tag":       "mgmt-gvnfm-handl-instantiate",
		"vnfr-name": vnfr.Name,
	}).Info("instantiating VNFR")

	if h.Level >= log.DebugLevel {
		h.WithFields(log.Fields{
			"tag":  "mgmt-gvnfm-handl-instantiate",
			"vnfr": JSON(vnfr),
		}).Debug("VNFR dump")
	}

	if vnfr.VDUs == nil {
		return nil, errors.New("no VDU provided")
	}

	for _, vdu := range vnfr.VDUs {
		for _, vnfc := range vdu.VNFCInstances {
			if _, err := h.instantiateVNFC(vnfc); err != nil {
				return nil, err
			}
		}
	}

	return vnfr, nil
}

// Modify allows making structural changes (e.g.configuration, topology, behavior, redundancy model) to a VNF instance.
func (h *handl) Modify(vnfr *catalogue.VirtualNetworkFunctionRecord,
	dependency *catalogue.VNFRecordDependency) (*catalogue.VirtualNetworkFunctionRecord, error) {

	md := collectMetadata(vnfr, dependency)

	if h.Level >= log.DebugLevel {
		h.WithFields(log.Fields{
			"tag":             "mgmt-gvnfm-handl-modify",
			"vnfr":            JSON(vnfr),
			"vnfr-dependency": JSON(dependency),
			"md":              md,
		}).Debug("modifying VNFR")
	}

	for _, vdu := range vnfr.VDUs {
		for _, vnfcInstance := range vdu.VNFCInstances {
			if err := h.applyVNFCMetadata(vnfcInstance, md); err != nil {
				return nil, err
			}
		}
	}

	return vnfr, nil
}

// Query allows retrieving a VNF instance state and attributes. (not implemented)
func (h *handl) Query() error {
	h.WithFields(log.Fields{
		"tag": "mgmt-gvnfm-handl-query",
	}).Warn("query invoked, not implemented")

	return nil
}

func (h *handl) Resume(vnfr *catalogue.VirtualNetworkFunctionRecord,
	vnfcInstance *catalogue.VNFCInstance,
	dependency *catalogue.VNFRecordDependency) (*catalogue.VirtualNetworkFunctionRecord, error) {

	h.WithFields(log.Fields{
		"tag":       "mgmt-gvnfm-handl-resume",
		"vnfr-name": vnfr.Name,
		"vnfr-id":   vnfr.ID,
	}).Info("resuming VNFR")

	return vnfr, nil
}

// Scale allows scaling (out / in) a VNF instance.
func (h *handl) Scale(action catalogue.Action,
	vnfr *catalogue.VirtualNetworkFunctionRecord,
	component catalogue.Component,
	scripts interface{},
	dependency *catalogue.VNFRecordDependency) (*catalogue.VirtualNetworkFunctionRecord, error) {

	h.WithFields(log.Fields{
		"tag":       "mgmt-gvnfm-handl-scale",
		"vnfr-name": vnfr.Name,
		"vnfr-id":   vnfr.ID,
		"action":    action,
	}).Debug("scaling VNFR")

	vnfcInstance, ok := component.(*catalogue.VNFCInstance)
	if !ok {
		return nil, fmt.Errorf("component is %T, expected %T", component, (*catalogue.VIMInstance)(nil))
	}

	switch action {
	case catalogue.ActionScaleIn:
		if err := h.scaleIn(vnfcInstance); err != nil {
			return nil, err
		}

	case catalogue.ActionScaleOut:
		if err := h.scaleOut(vnfr, vnfcInstance, dependency); err != nil {
			return nil, err
		}

	default:
		return nil, fmt.Errorf("invalid action for Scale: %v", action)
	}

	return vnfr, nil
}

// Start starts a VNFR.
func (h *handl) Start(vnfr *catalogue.VirtualNetworkFunctionRecord) (*catalogue.VirtualNetworkFunctionRecord, error) {
	h.WithFields(log.Fields{
		"tag":       "mgmt-gvnfm-handl-start",
		"vnfr-name": vnfr.Name,
	}).Debug("starting VNFR")

	for _, vdu := range vnfr.VDUs {
		for _, vnfcInstance := range vdu.VNFCInstances {
			if _, err := h.StartVNFCInstance(vnfr, vnfcInstance); err != nil {
				return nil, err
			}
		}
	}

	return vnfr, nil
}

func (h *handl) StartVNFCInstance(vnfr *catalogue.VirtualNetworkFunctionRecord,
	vnfcInstance *catalogue.VNFCInstance) (*catalogue.VirtualNetworkFunctionRecord, error) {

	h.WithFields(log.Fields{
		"tag":                "mgmt-gvnfm-handl-start_vnfc_instance",
		"vnfc_instance-name": vnfcInstance.Hostname,
		"vnfc_instance-id":   vnfcInstance.ID,
		"vim_instance-id":    vnfcInstance.VIMID,
	}).Info("starting VNFCInstance")

	if err := h.mgmt(vnfcInstance.VIMID).Start(vnfcInstance.Hostname); err != nil {
		return nil, err
	}

	return vnfr, nil
}

// Stop stops a previously created VNF instance.
func (h *handl) Stop(vnfr *catalogue.VirtualNetworkFunctionRecord) (*catalogue.VirtualNetworkFunctionRecord, error) {
	h.WithFields(log.Fields{
		"tag":       "mgmt-gvnfm-handl-stop",
		"vnfr-name": vnfr.Name,
	}).Info("stopping VNFR")

	return vnfr, nil
}

func (h *handl) StopVNFCInstance(vnfr *catalogue.VirtualNetworkFunctionRecord,
	vnfcInstance *catalogue.VNFCInstance) (*catalogue.VirtualNetworkFunctionRecord, error) {

	h.WithFields(log.Fields{
		"tag":                "mgmt-gvnfm-handl-stop_vnfc_instance",
		"vnfc_instance-name": vnfcInstance.Hostname,
		"vnfc_instance-id":   vnfcInstance.ID,
	}).Info("stopping VNFCInstance")

	return vnfr, nil
}

// Terminate allows terminating gracefully or forcefully a previously created VNF instance.
func (h *handl) Terminate(vnfr *catalogue.VirtualNetworkFunctionRecord) (*catalogue.VirtualNetworkFunctionRecord, error) {
	h.WithFields(log.Fields{
		"tag":             "mgmt-gvnfm-handl-terminate",
		"vnfr-name":       vnfr.Name,
		"vnfr-hb_version": vnfr.HbVersion,
	}).Info("terminating VNFR")

	for _, event := range vnfr.LifecycleEvents {
		if event.Event == catalogue.EventRelease {
			for _, vdu := range vnfr.VDUs {
				h.WithFields(log.Fields{
					"tag":       "mgmt-gvnfm-handl-terminate",
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
func (h *handl) UpdateSoftware(script *catalogue.Script,
	vnfr *catalogue.VirtualNetworkFunctionRecord) (*catalogue.VirtualNetworkFunctionRecord, error) {
	h.WithFields(log.Fields{
		"tag":       "mgmt-gvnfm-handl-update_software",
		"script":    script,
		"vnfr-name": vnfr.Name,
		"vnfr-id":   vnfr.ID,
	}).Info("updating software for VNFR")

	time.Sleep(3 * time.Second)

	return vnfr, nil
}

// UpgradeSoftware allows deploying a new software release to a VNF instance.
func (h *handl) UpgradeSoftware() error {
	h.WithFields(log.Fields{
		"tag": "mgmt-gvnfm-handl-update_software",
	}).Warn("UpgradeSoftware called - but it's no-op")

	return nil
}

// UserData returns a string containing UserData.
func (h *handl) UserData() string {
	h.WithFields(log.Fields{
		"tag": "mgmt-gvnfm-handl-user_data",
	}).Info("returning UserData")

	return "#!/usr/bin/env sh\n"
}

func (h *handl) applyVNFCMetadata(vnfcInstance *catalogue.VNFCInstance, md map[string]string) error {
	return h.mgmt(vnfcInstance.VIMID).AddMetadata(vnfcInstance.Hostname, md)
}

func (h *handl) instantiateVNFC(vnfc *catalogue.VNFCInstance) (*catalogue.Server, error) {
	srv, err := h.mgmt(vnfc.VIMID).Check(vnfc.Hostname)
	if err != nil {
		return nil, err
	}

	h.WithFields(log.Fields{
		"tag":      "mgmt-gvnfm-handl-instantiate",
		"srv-name": srv.Name,
	}).Debug("server is alive")

	return srv, err
}

func (h *handl) mgmt(vimID string) mgmt.VIMConnector {
	return mgmt.NewConnector(vimID, h.acc)
}

func (h *handl) scaleIn(vnfcInstance *catalogue.VNFCInstance) error {
	h.WithFields(log.Fields{
		"tag":           "mgmt-gvnfm-handl-scaleIn",
		"vnfc_instance": vnfcInstance.ID,
	}).Debug("scaling in")

	// do absolutely nothing

	return nil
}

func (h *handl) scaleOut(vnfr *catalogue.VirtualNetworkFunctionRecord,
	vnfcInstance *catalogue.VNFCInstance, dependency *catalogue.VNFRecordDependency) error {

	h.WithFields(log.Fields{
		"tag":           "mgmt-gvnfm-handl-scaleOut",
		"vnfc_instance": vnfcInstance.ID,
	}).Debug("scaling out")

	// scale out equals to instantiate + configure + start

	if _, err := h.instantiateVNFC(vnfcInstance); err != nil {
		return err
	}

	// configure

	if err := h.applyVNFCMetadata(vnfcInstance, collectMetadata(vnfr, dependency)); err != nil {
		return err
	}

	// start

	_, err := h.StartVNFCInstance(vnfr, vnfcInstance)

	return err
}

func JSON(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		return fmt.Sprintf("<<%v>>", err)
	}

	return string(b)
}

func collectMetadata(vnfr *catalogue.VirtualNetworkFunctionRecord,
	dependency *catalogue.VNFRecordDependency) map[string]string {

	md := make(map[string]string)

	for _, confParam := range vnfr.Provides.ConfigurationParameters {
		md[confParam.ConfKey] = confParam.Value
	}

	for _, confParam := range vnfr.Configurations.ConfigurationParameters {
		md[confParam.ConfKey] = confParam.Value
	}

	if dependency != nil {
		for ptype, depParam := range dependency.Parameters {
			if ptype == vnfr.Type {
				continue // skip yourself
			}

			for pkey, pval := range depParam.Parameters {
				key := strings.ToUpper(fmt.Sprintf("%s_%s", ptype, pkey))

				md[key] = pval
			}
		}

		for ptype, vnfcDepParam := range dependency.VNFCParameters {
			if ptype == vnfr.Type {
				continue // skip yourself
			}

			for _, depParam := range vnfcDepParam.Parameters {
				for pkey, pval := range depParam.Parameters {
					key := strings.ToUpper(fmt.Sprintf("%s_%s", ptype, pkey))

					md[key] = pval
				}
			}
		}
	}

	return md
}
