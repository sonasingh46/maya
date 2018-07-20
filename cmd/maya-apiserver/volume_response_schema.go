package main

// Server Error
// swagger:response serverError
type ServerError struct {
// in: body
Body struct {
// The Error Message
// example: Something went wrong.
// Required: true
Message string `json:"message"`
}
}

// Summary volume data that matches the query
// swagger:response volumeSummary
type VolumeSummary struct {
	// Volume Details
	// in: body
	Body struct {
		// List of openebs volumes.
		// Required: true
		Items []Item `json:"items,omitempty"`
		// Required: true
		ListMetaData ListMetaData `json:"metadata"`
	}
}
type Item struct {
	Status Status `json:"status,omitempty"`
	ItemMetadata ItemMetadata `json:"metadata,omitempty"`
}
// Standard list metadata.
type ListMetaData struct {
}
type Status struct {
	// example: null
	StatusProperty1 string `json:"Message"`
	// example: NotRunning
	Property2 string `json:"Phase"`
	// example: null
	Property3 string `json:"Reason"`
}
type ItemMetadata struct {
	// Volume Name
	//example: OpenEBS Volume
	//Required: true
	VolumeName string `json:"name"`
	// The time the snapshot was successfully created.
	// example: null
	// Required: true
	CreationTimestamp string `json:"creationTimestamp"`
	Labels Labels `json:"labels",omitempty`
	Annotations Annotations `json:"annotations"`

}
type Labels struct{

}
type Annotations struct{
	// Status of replica containers in a pod.
	// example: Waiting
	// Required: true
	Property2 string `json:"openebs.io/replica-container-status"`
	// Name of storage pool.
	// example: default
	// Required: true
	Property4 string `json:"openebs.io/storage-pool"`
	// Volume capacity.
	// example: 5G
	// Required: true
	Property7 string `json:"openebs.io/capacity"`
	// Status of controller containers in a pod.
	// example: Running
	// required: true
	Property8 string `json:"openebs.io/controller-container-status"`
	// Volume iqn.
	// example: iqn.2016-09.com.openebs.jiva:test
	// required: true
	Property9 string `json:"openebs.io/jiva-iqn"`
	// Volume type.
	// example: jiva
	// Required: true
	Property10 string `json:"openebs.io/volume-type"`
	//
	// example: 1
	// Required: true
	Property11 string `json:"deployment.kubernetes.io/revision"`
	// Jiva replica IP addresses.
	// example: 172.17.0.9,nil,nil
	// Required: true
	Property13 string `json:"openebs.io/jiva-replica-ips"`
	// Jiva controller cluster IP.
	// example: 10.0.0.205
	// Required: true
	Property14 string `json:"openebs.io/jiva-controller-cluster-ip"`
	// Jiva target portal address.
	// example: 10.0.0.205:3260
	// required: true
	Property16 string `json:"openebs.io/jiva-target-portal"`
	// Monitoring details.
	// example: false
	// Required: true
	Property18 string `json:"openebs.io/volume-monitor"`
	// Jiva controller IP addresses.
	// example: 172.17.0.8
	// Required: true
	Property20 string `json:"openebs.io/jiva-controller-ips"`
	// Jiva replica status.
	// example: Running,Pending,Pending
	// Required: true
	Property22 string `json:"openebs.io/jiva-replica-status"`
	// Jiva replica count.
	// example: 3
	// Required: true
	Property23 string `json:"openebs.io/jiva-replica-count"`
	// Jiva controller status.
	// example: Running
	// Required: true
	Property24 string `json:"openebs.io/jiva-controller-status"`
}

// Summary volume data that matches the query
// swagger:response SpecificVolumeSummary
type SpecificVolumeSummary struct {
	// Volume Details
	// in: body
	Body struct {
		// Volume data
		// Required: true
		Status Status `json:"status,omitempty"`
		// Required: true
		ItemMetadata ItemMetadata `json:"metadata,omitempty"`
	}
}

// Volume Created
// swagger:response VolumeCreated
type MyJsonName struct {
	// Volume Created
	// in: body
	Body struct {
		// List of volumes' data
		// Required: true
		ItemMetadata ItemMetadata `json:"metadata","metadata"`
		// Required: true
		Status Status `json:"status"`
	}
}

// No Such Volume
// swagger:response noSuchVolume
type noSuchVolume struct {
	// No Such Volume
	// in: body
	Body struct {
		// The Error Message
		// example: Volume 'example-vol-name' not found.
		// Required: true
		Message string `json:"message"`
	}
}
//Volume Deleted
// swagger:response successDelete
type VolumeDelete struct {
	// in: body
	Body struct {
		// The Success Message
		// example: Volume 'example-volume' deleted successfully
		// Required: true
		Message string `json:"message"`
	}
}
//Bad Request
// swagger:response badRequest
type BadRequest struct {
	// in: body
	Body struct {
		// The Error Message
		// example: Volume name missing.
		// Required: true
		Message string `json:"message"`
	}
}