package analytics

import (
	menv "github.com/openebs/maya/pkg/env/v1alpha1"
	analytics "github.com/openebs/maya/pkg/usage"
)

// SendEventOrIgnore sends anonymous local-pv provision/delete events
func SendEventOrIgnore(pvName, capacity, stgType, method string) {
	if method == analytics.VolumeProvision {
		stgType = "local-" + stgType
	}
	if menv.Truthy(menv.OpenEBSEnableAnalytics) {
		analytics.New().Build().ApplicationBuilder().
			SetVolumeType(stgType, method).
			SetDocumentTitle(pvName).
			SetLabel(analytics.EventLabelCapacity).
			SetReplicaCount(analytics.LocalPVReplicaCount, method).
			SetCategory(method).
			SetVolumeCapacity(capacity).Send()
	}
}
