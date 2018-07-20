package main
// swagger:parameters GetSnapshot
type VolumeName struct {
	// Name of volume.
	// unique: true
	// in: path
	VolumeName string `json:"volumeName"`
}
// swagger:parameters CreateSnapshot RevertSnapshot
type CreateSnapshot struct {
	//in: body
	Body struct {
		//example: OpenEBS Volume
		//Required: true
		Metadata SnapshotCreateMetadata `json:"metadata,omitempty"`
		//Required: true
		Spec Spec `json:"spec"`
	}
}
type SnapshotCreateMetadata struct{
	//Snapshot name
	//example: snapshotName
	//Required: true
	Name string `json:"name"`
}
type Spec struct{
	///Volume name
	//example: pvc-10eb916d-4159-11e8-b222-42010a800202
	//Required: true
	VolumeName string `json:"volumeName"`
}