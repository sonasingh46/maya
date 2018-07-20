package main
// Summary snapshot data that matches the query
// swagger:response snapshotSummary
type SnapshotSummary struct {
	// in: body
	Body struct {
		// Required: true
		//example: volume-head-000.img
		Snapshot Snapshotschema `json:"SnapshotName,omitempty"`

	}
}
type Snapshotschema struct {
	// Name of snapshot
	// Required: true
	// example: volume-snap-ae1398a8-1151-4062-8e97-b4d074ee62a3.img
	Name        string   `json:"name"`
	// Name of parent snapshot
	// Required: true
	// example : volume-snap-ae1398a8-1151-4062-8e97-b4d074ee62a3.img
	Parent      string   `json:"parent"`
	// Name of children snapshot
	// Required: true
	Children    []string `json:"children"`
	// Required: true
	// example :false
	Removed     bool     `json:"removed"`
	// Required: true
	// example :false
	UserCreated bool     `json:"usercreated"`
	// Required: true
	// example: 2018-05-02T21:06:22Z
	Created     string   `json:"created"`
	// Required: true
	// example: 5
	Size        string   `json:"size"`
}

// Revert snapshot data that matches the query
// swagger:response snapshotReverted
type SnapshotReverted struct {
	// in: body
	Body struct {
		// Required: true
		// example: Reverting to snapshot [pvcsnap1] of volume [pvc-d07f8a7a-2d99-11e8-bbe0-42010a800243]
		Message string `json:"message,omitempty"`
	}
}

// Create snapshot data that matches the query
// swagger:response snapshotCreated
type SnapshotCreated struct {
	// in: body
	Body struct {
		// Required: true
		//example:snapshotOutput
		Typ string `json:"type,omitempty"`
		// Required: true
		//example: pvc-snap1
		Id string `json:"id,omitempty"`
		// Required: true
		Links Links `json:"links"`
		// Required: true
		Actions Actions `json:"actions"`

	}
}
type Links struct{
	// Required: true
	// example: http://10.32.1.39:9501/v1/snapshotoutputs/pvsnap1
	Self string `json:"self"`
}
type Actions struct{
}
