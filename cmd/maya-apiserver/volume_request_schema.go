package main
// swagger:parameters DeleteVolume GetSpecificVolume
type VolumeNameParam struct {
	// Name of volume.
	//
	// unique: true
	// in: path
	VolumeName string `json:"name"`
}
// swagger:parameters CreateVolume
type CreateVolume struct {
	//in: body
	Body struct {
		//example: OpenEBS Volume
		//Required: true
		Metadata VolumeCreateMetadata `json:"metadata,omitempty"`
		VolumeCreateLabels VolumeCreateLabels `json:"labels"`
	}
}
type VolumeCreateMetadata struct {
	// Volume Name
	//example: OpenEBS Volume
	//Required: true
	VolumeName string `json:"name"`
}
type VolumeCreateLabels struct {
Storage string `json:"volumeprovisioner.mapi.openebs.io/storage-size"`
}