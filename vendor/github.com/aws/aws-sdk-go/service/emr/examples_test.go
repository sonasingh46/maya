// THIS FILE IS AUTOMATICALLY GENERATED. DO NOT EDIT.

package emr_test

import (
	"bytes"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/emr"
)

var _ time.Duration
var _ bytes.Buffer

func ExampleEMR_AddInstanceGroups() {
	sess, err := session.NewSession()
	if err != nil {
		fmt.Println("failed to create session,", err)
		return
	}

	svc := emr.New(sess)

	params := &emr.AddInstanceGroupsInput{
		InstanceGroups: []*emr.InstanceGroupConfig{ // Required
			{ // Required
				InstanceCount: aws.Int64(1),                   // Required
				InstanceRole:  aws.String("InstanceRoleType"), // Required
				InstanceType:  aws.String("InstanceType"),     // Required
				AutoScalingPolicy: &emr.AutoScalingPolicy{
					Constraints: &emr.ScalingConstraints{ // Required
						MaxCapacity: aws.Int64(1), // Required
						MinCapacity: aws.Int64(1), // Required
					},
					Rules: []*emr.ScalingRule{ // Required
						{ // Required
							Action: &emr.ScalingAction{ // Required
								SimpleScalingPolicyConfiguration: &emr.SimpleScalingPolicyConfiguration{ // Required
									ScalingAdjustment: aws.Int64(1), // Required
									AdjustmentType:    aws.String("AdjustmentType"),
									CoolDown:          aws.Int64(1),
								},
								Market: aws.String("MarketType"),
							},
							Name: aws.String("String"), // Required
							Trigger: &emr.ScalingTrigger{ // Required
								CloudWatchAlarmDefinition: &emr.CloudWatchAlarmDefinition{ // Required
									ComparisonOperator: aws.String("ComparisonOperator"), // Required
									MetricName:         aws.String("String"),             // Required
									Period:             aws.Int64(1),                     // Required
									Threshold:          aws.Float64(1.0),                 // Required
									Dimensions: []*emr.MetricDimension{
										{ // Required
											Key:   aws.String("String"),
											Value: aws.String("String"),
										},
										// More values...
									},
									EvaluationPeriods: aws.Int64(1),
									Namespace:         aws.String("String"),
									Statistic:         aws.String("Statistic"),
									Unit:              aws.String("Unit"),
								},
							},
							Description: aws.String("String"),
						},
						// More values...
					},
				},
				BidPrice: aws.String("XmlStringMaxLen256"),
				Configurations: []*emr.Configuration{
					{ // Required
						Classification: aws.String("String"),
						Configurations: []*emr.Configuration{
						// Recursive values...
						},
						Properties: map[string]*string{
							"Key": aws.String("String"), // Required
							// More values...
						},
					},
					// More values...
				},
				EbsConfiguration: &emr.EbsConfiguration{
					EbsBlockDeviceConfigs: []*emr.EbsBlockDeviceConfig{
						{ // Required
							VolumeSpecification: &emr.VolumeSpecification{ // Required
								SizeInGB:   aws.Int64(1),         // Required
								VolumeType: aws.String("String"), // Required
								Iops:       aws.Int64(1),
							},
							VolumesPerInstance: aws.Int64(1),
						},
						// More values...
					},
					EbsOptimized: aws.Bool(true),
				},
				Market: aws.String("MarketType"),
				Name:   aws.String("XmlStringMaxLen256"),
			},
			// More values...
		},
		JobFlowId: aws.String("XmlStringMaxLen256"), // Required
	}
	resp, err := svc.AddInstanceGroups(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleEMR_AddJobFlowSteps() {
	sess, err := session.NewSession()
	if err != nil {
		fmt.Println("failed to create session,", err)
		return
	}

	svc := emr.New(sess)

	params := &emr.AddJobFlowStepsInput{
		JobFlowId: aws.String("XmlStringMaxLen256"), // Required
		Steps: []*emr.StepConfig{ // Required
			{ // Required
				HadoopJarStep: &emr.HadoopJarStepConfig{ // Required
					Jar: aws.String("XmlString"), // Required
					Args: []*string{
						aws.String("XmlString"), // Required
						// More values...
					},
					MainClass: aws.String("XmlString"),
					Properties: []*emr.KeyValue{
						{ // Required
							Key:   aws.String("XmlString"),
							Value: aws.String("XmlString"),
						},
						// More values...
					},
				},
				Name:            aws.String("XmlStringMaxLen256"), // Required
				ActionOnFailure: aws.String("ActionOnFailure"),
			},
			// More values...
		},
	}
	resp, err := svc.AddJobFlowSteps(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleEMR_AddTags() {
	sess, err := session.NewSession()
	if err != nil {
		fmt.Println("failed to create session,", err)
		return
	}

	svc := emr.New(sess)

	params := &emr.AddTagsInput{
		ResourceId: aws.String("ResourceId"), // Required
		Tags: []*emr.Tag{ // Required
			{ // Required
				Key:   aws.String("String"),
				Value: aws.String("String"),
			},
			// More values...
		},
	}
	resp, err := svc.AddTags(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleEMR_CancelSteps() {
	sess, err := session.NewSession()
	if err != nil {
		fmt.Println("failed to create session,", err)
		return
	}

	svc := emr.New(sess)

	params := &emr.CancelStepsInput{
		ClusterId: aws.String("XmlStringMaxLen256"),
		StepIds: []*string{
			aws.String("XmlStringMaxLen256"), // Required
			// More values...
		},
	}
	resp, err := svc.CancelSteps(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleEMR_CreateSecurityConfiguration() {
	sess, err := session.NewSession()
	if err != nil {
		fmt.Println("failed to create session,", err)
		return
	}

	svc := emr.New(sess)

	params := &emr.CreateSecurityConfigurationInput{
		Name: aws.String("XmlString"), // Required
		SecurityConfiguration: aws.String("String"), // Required
	}
	resp, err := svc.CreateSecurityConfiguration(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleEMR_DeleteSecurityConfiguration() {
	sess, err := session.NewSession()
	if err != nil {
		fmt.Println("failed to create session,", err)
		return
	}

	svc := emr.New(sess)

	params := &emr.DeleteSecurityConfigurationInput{
		Name: aws.String("XmlString"), // Required
	}
	resp, err := svc.DeleteSecurityConfiguration(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleEMR_DescribeCluster() {
	sess, err := session.NewSession()
	if err != nil {
		fmt.Println("failed to create session,", err)
		return
	}

	svc := emr.New(sess)

	params := &emr.DescribeClusterInput{
		ClusterId: aws.String("ClusterId"), // Required
	}
	resp, err := svc.DescribeCluster(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleEMR_DescribeJobFlows() {
	sess, err := session.NewSession()
	if err != nil {
		fmt.Println("failed to create session,", err)
		return
	}

	svc := emr.New(sess)

	params := &emr.DescribeJobFlowsInput{
		CreatedAfter:  aws.Time(time.Now()),
		CreatedBefore: aws.Time(time.Now()),
		JobFlowIds: []*string{
			aws.String("XmlString"), // Required
			// More values...
		},
		JobFlowStates: []*string{
			aws.String("JobFlowExecutionState"), // Required
			// More values...
		},
	}
	resp, err := svc.DescribeJobFlows(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleEMR_DescribeSecurityConfiguration() {
	sess, err := session.NewSession()
	if err != nil {
		fmt.Println("failed to create session,", err)
		return
	}

	svc := emr.New(sess)

	params := &emr.DescribeSecurityConfigurationInput{
		Name: aws.String("XmlString"), // Required
	}
	resp, err := svc.DescribeSecurityConfiguration(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleEMR_DescribeStep() {
	sess, err := session.NewSession()
	if err != nil {
		fmt.Println("failed to create session,", err)
		return
	}

	svc := emr.New(sess)

	params := &emr.DescribeStepInput{
		ClusterId: aws.String("ClusterId"), // Required
		StepId:    aws.String("StepId"),    // Required
	}
	resp, err := svc.DescribeStep(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleEMR_ListBootstrapActions() {
	sess, err := session.NewSession()
	if err != nil {
		fmt.Println("failed to create session,", err)
		return
	}

	svc := emr.New(sess)

	params := &emr.ListBootstrapActionsInput{
		ClusterId: aws.String("ClusterId"), // Required
		Marker:    aws.String("Marker"),
	}
	resp, err := svc.ListBootstrapActions(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleEMR_ListClusters() {
	sess, err := session.NewSession()
	if err != nil {
		fmt.Println("failed to create session,", err)
		return
	}

	svc := emr.New(sess)

	params := &emr.ListClustersInput{
		ClusterStates: []*string{
			aws.String("ClusterState"), // Required
			// More values...
		},
		CreatedAfter:  aws.Time(time.Now()),
		CreatedBefore: aws.Time(time.Now()),
		Marker:        aws.String("Marker"),
	}
	resp, err := svc.ListClusters(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleEMR_ListInstanceGroups() {
	sess, err := session.NewSession()
	if err != nil {
		fmt.Println("failed to create session,", err)
		return
	}

	svc := emr.New(sess)

	params := &emr.ListInstanceGroupsInput{
		ClusterId: aws.String("ClusterId"), // Required
		Marker:    aws.String("Marker"),
	}
	resp, err := svc.ListInstanceGroups(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleEMR_ListInstances() {
	sess, err := session.NewSession()
	if err != nil {
		fmt.Println("failed to create session,", err)
		return
	}

	svc := emr.New(sess)

	params := &emr.ListInstancesInput{
		ClusterId:       aws.String("ClusterId"), // Required
		InstanceGroupId: aws.String("InstanceGroupId"),
		InstanceGroupTypes: []*string{
			aws.String("InstanceGroupType"), // Required
			// More values...
		},
		InstanceStates: []*string{
			aws.String("InstanceState"), // Required
			// More values...
		},
		Marker: aws.String("Marker"),
	}
	resp, err := svc.ListInstances(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleEMR_ListSecurityConfigurations() {
	sess, err := session.NewSession()
	if err != nil {
		fmt.Println("failed to create session,", err)
		return
	}

	svc := emr.New(sess)

	params := &emr.ListSecurityConfigurationsInput{
		Marker: aws.String("Marker"),
	}
	resp, err := svc.ListSecurityConfigurations(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleEMR_ListSteps() {
	sess, err := session.NewSession()
	if err != nil {
		fmt.Println("failed to create session,", err)
		return
	}

	svc := emr.New(sess)

	params := &emr.ListStepsInput{
		ClusterId: aws.String("ClusterId"), // Required
		Marker:    aws.String("Marker"),
		StepIds: []*string{
			aws.String("XmlString"), // Required
			// More values...
		},
		StepStates: []*string{
			aws.String("StepState"), // Required
			// More values...
		},
	}
	resp, err := svc.ListSteps(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleEMR_ModifyInstanceGroups() {
	sess, err := session.NewSession()
	if err != nil {
		fmt.Println("failed to create session,", err)
		return
	}

	svc := emr.New(sess)

	params := &emr.ModifyInstanceGroupsInput{
		ClusterId: aws.String("ClusterId"),
		InstanceGroups: []*emr.InstanceGroupModifyConfig{
			{ // Required
				InstanceGroupId: aws.String("XmlStringMaxLen256"), // Required
				EC2InstanceIdsToTerminate: []*string{
					aws.String("InstanceId"), // Required
					// More values...
				},
				InstanceCount: aws.Int64(1),
				ShrinkPolicy: &emr.ShrinkPolicy{
					DecommissionTimeout: aws.Int64(1),
					InstanceResizePolicy: &emr.InstanceResizePolicy{
						InstanceTerminationTimeout: aws.Int64(1),
						InstancesToProtect: []*string{
							aws.String("InstanceId"), // Required
							// More values...
						},
						InstancesToTerminate: []*string{
							aws.String("InstanceId"), // Required
							// More values...
						},
					},
				},
			},
			// More values...
		},
	}
	resp, err := svc.ModifyInstanceGroups(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleEMR_PutAutoScalingPolicy() {
	sess, err := session.NewSession()
	if err != nil {
		fmt.Println("failed to create session,", err)
		return
	}

	svc := emr.New(sess)

	params := &emr.PutAutoScalingPolicyInput{
		AutoScalingPolicy: &emr.AutoScalingPolicy{ // Required
			Constraints: &emr.ScalingConstraints{ // Required
				MaxCapacity: aws.Int64(1), // Required
				MinCapacity: aws.Int64(1), // Required
			},
			Rules: []*emr.ScalingRule{ // Required
				{ // Required
					Action: &emr.ScalingAction{ // Required
						SimpleScalingPolicyConfiguration: &emr.SimpleScalingPolicyConfiguration{ // Required
							ScalingAdjustment: aws.Int64(1), // Required
							AdjustmentType:    aws.String("AdjustmentType"),
							CoolDown:          aws.Int64(1),
						},
						Market: aws.String("MarketType"),
					},
					Name: aws.String("String"), // Required
					Trigger: &emr.ScalingTrigger{ // Required
						CloudWatchAlarmDefinition: &emr.CloudWatchAlarmDefinition{ // Required
							ComparisonOperator: aws.String("ComparisonOperator"), // Required
							MetricName:         aws.String("String"),             // Required
							Period:             aws.Int64(1),                     // Required
							Threshold:          aws.Float64(1.0),                 // Required
							Dimensions: []*emr.MetricDimension{
								{ // Required
									Key:   aws.String("String"),
									Value: aws.String("String"),
								},
								// More values...
							},
							EvaluationPeriods: aws.Int64(1),
							Namespace:         aws.String("String"),
							Statistic:         aws.String("Statistic"),
							Unit:              aws.String("Unit"),
						},
					},
					Description: aws.String("String"),
				},
				// More values...
			},
		},
		ClusterId:       aws.String("ClusterId"),       // Required
		InstanceGroupId: aws.String("InstanceGroupId"), // Required
	}
	resp, err := svc.PutAutoScalingPolicy(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleEMR_RemoveAutoScalingPolicy() {
	sess, err := session.NewSession()
	if err != nil {
		fmt.Println("failed to create session,", err)
		return
	}

	svc := emr.New(sess)

	params := &emr.RemoveAutoScalingPolicyInput{
		ClusterId:       aws.String("ClusterId"),       // Required
		InstanceGroupId: aws.String("InstanceGroupId"), // Required
	}
	resp, err := svc.RemoveAutoScalingPolicy(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleEMR_RemoveTags() {
	sess, err := session.NewSession()
	if err != nil {
		fmt.Println("failed to create session,", err)
		return
	}

	svc := emr.New(sess)

	params := &emr.RemoveTagsInput{
		ResourceId: aws.String("ResourceId"), // Required
		TagKeys: []*string{ // Required
			aws.String("String"), // Required
			// More values...
		},
	}
	resp, err := svc.RemoveTags(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleEMR_RunJobFlow() {
	sess, err := session.NewSession()
	if err != nil {
		fmt.Println("failed to create session,", err)
		return
	}

	svc := emr.New(sess)

	params := &emr.RunJobFlowInput{
		Instances: &emr.JobFlowInstancesConfig{ // Required
			AdditionalMasterSecurityGroups: []*string{
				aws.String("XmlStringMaxLen256"), // Required
				// More values...
			},
			AdditionalSlaveSecurityGroups: []*string{
				aws.String("XmlStringMaxLen256"), // Required
				// More values...
			},
			Ec2KeyName:                    aws.String("XmlStringMaxLen256"),
			Ec2SubnetId:                   aws.String("XmlStringMaxLen256"),
			EmrManagedMasterSecurityGroup: aws.String("XmlStringMaxLen256"),
			EmrManagedSlaveSecurityGroup:  aws.String("XmlStringMaxLen256"),
			HadoopVersion:                 aws.String("XmlStringMaxLen256"),
			InstanceCount:                 aws.Int64(1),
			InstanceGroups: []*emr.InstanceGroupConfig{
				{ // Required
					InstanceCount: aws.Int64(1),                   // Required
					InstanceRole:  aws.String("InstanceRoleType"), // Required
					InstanceType:  aws.String("InstanceType"),     // Required
					AutoScalingPolicy: &emr.AutoScalingPolicy{
						Constraints: &emr.ScalingConstraints{ // Required
							MaxCapacity: aws.Int64(1), // Required
							MinCapacity: aws.Int64(1), // Required
						},
						Rules: []*emr.ScalingRule{ // Required
							{ // Required
								Action: &emr.ScalingAction{ // Required
									SimpleScalingPolicyConfiguration: &emr.SimpleScalingPolicyConfiguration{ // Required
										ScalingAdjustment: aws.Int64(1), // Required
										AdjustmentType:    aws.String("AdjustmentType"),
										CoolDown:          aws.Int64(1),
									},
									Market: aws.String("MarketType"),
								},
								Name: aws.String("String"), // Required
								Trigger: &emr.ScalingTrigger{ // Required
									CloudWatchAlarmDefinition: &emr.CloudWatchAlarmDefinition{ // Required
										ComparisonOperator: aws.String("ComparisonOperator"), // Required
										MetricName:         aws.String("String"),             // Required
										Period:             aws.Int64(1),                     // Required
										Threshold:          aws.Float64(1.0),                 // Required
										Dimensions: []*emr.MetricDimension{
											{ // Required
												Key:   aws.String("String"),
												Value: aws.String("String"),
											},
											// More values...
										},
										EvaluationPeriods: aws.Int64(1),
										Namespace:         aws.String("String"),
										Statistic:         aws.String("Statistic"),
										Unit:              aws.String("Unit"),
									},
								},
								Description: aws.String("String"),
							},
							// More values...
						},
					},
					BidPrice: aws.String("XmlStringMaxLen256"),
					Configurations: []*emr.Configuration{
						{ // Required
							Classification: aws.String("String"),
							Configurations: []*emr.Configuration{
							// Recursive values...
							},
							Properties: map[string]*string{
								"Key": aws.String("String"), // Required
								// More values...
							},
						},
						// More values...
					},
					EbsConfiguration: &emr.EbsConfiguration{
						EbsBlockDeviceConfigs: []*emr.EbsBlockDeviceConfig{
							{ // Required
								VolumeSpecification: &emr.VolumeSpecification{ // Required
									SizeInGB:   aws.Int64(1),         // Required
									VolumeType: aws.String("String"), // Required
									Iops:       aws.Int64(1),
								},
								VolumesPerInstance: aws.Int64(1),
							},
							// More values...
						},
						EbsOptimized: aws.Bool(true),
					},
					Market: aws.String("MarketType"),
					Name:   aws.String("XmlStringMaxLen256"),
				},
				// More values...
			},
			KeepJobFlowAliveWhenNoSteps: aws.Bool(true),
			MasterInstanceType:          aws.String("InstanceType"),
			Placement: &emr.PlacementType{
				AvailabilityZone: aws.String("XmlString"), // Required
			},
			ServiceAccessSecurityGroup: aws.String("XmlStringMaxLen256"),
			SlaveInstanceType:          aws.String("InstanceType"),
			TerminationProtected:       aws.Bool(true),
		},
		Name:           aws.String("XmlStringMaxLen256"), // Required
		AdditionalInfo: aws.String("XmlString"),
		AmiVersion:     aws.String("XmlStringMaxLen256"),
		Applications: []*emr.Application{
			{ // Required
				AdditionalInfo: map[string]*string{
					"Key": aws.String("String"), // Required
					// More values...
				},
				Args: []*string{
					aws.String("String"), // Required
					// More values...
				},
				Name:    aws.String("String"),
				Version: aws.String("String"),
			},
			// More values...
		},
		AutoScalingRole: aws.String("XmlString"),
		BootstrapActions: []*emr.BootstrapActionConfig{
			{ // Required
				Name: aws.String("XmlStringMaxLen256"), // Required
				ScriptBootstrapAction: &emr.ScriptBootstrapActionConfig{ // Required
					Path: aws.String("XmlString"), // Required
					Args: []*string{
						aws.String("XmlString"), // Required
						// More values...
					},
				},
			},
			// More values...
		},
		Configurations: []*emr.Configuration{
			{ // Required
				Classification: aws.String("String"),
				Configurations: []*emr.Configuration{
				// Recursive values...
				},
				Properties: map[string]*string{
					"Key": aws.String("String"), // Required
					// More values...
				},
			},
			// More values...
		},
		JobFlowRole: aws.String("XmlString"),
		LogUri:      aws.String("XmlString"),
		NewSupportedProducts: []*emr.SupportedProductConfig{
			{ // Required
				Args: []*string{
					aws.String("XmlString"), // Required
					// More values...
				},
				Name: aws.String("XmlStringMaxLen256"),
			},
			// More values...
		},
		ReleaseLabel:          aws.String("XmlStringMaxLen256"),
		ScaleDownBehavior:     aws.String("ScaleDownBehavior"),
		SecurityConfiguration: aws.String("XmlString"),
		ServiceRole:           aws.String("XmlString"),
		Steps: []*emr.StepConfig{
			{ // Required
				HadoopJarStep: &emr.HadoopJarStepConfig{ // Required
					Jar: aws.String("XmlString"), // Required
					Args: []*string{
						aws.String("XmlString"), // Required
						// More values...
					},
					MainClass: aws.String("XmlString"),
					Properties: []*emr.KeyValue{
						{ // Required
							Key:   aws.String("XmlString"),
							Value: aws.String("XmlString"),
						},
						// More values...
					},
				},
				Name:            aws.String("XmlStringMaxLen256"), // Required
				ActionOnFailure: aws.String("ActionOnFailure"),
			},
			// More values...
		},
		SupportedProducts: []*string{
			aws.String("XmlStringMaxLen256"), // Required
			// More values...
		},
		Tags: []*emr.Tag{
			{ // Required
				Key:   aws.String("String"),
				Value: aws.String("String"),
			},
			// More values...
		},
		VisibleToAllUsers: aws.Bool(true),
	}
	resp, err := svc.RunJobFlow(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleEMR_SetTerminationProtection() {
	sess, err := session.NewSession()
	if err != nil {
		fmt.Println("failed to create session,", err)
		return
	}

	svc := emr.New(sess)

	params := &emr.SetTerminationProtectionInput{
		JobFlowIds: []*string{ // Required
			aws.String("XmlString"), // Required
			// More values...
		},
		TerminationProtected: aws.Bool(true), // Required
	}
	resp, err := svc.SetTerminationProtection(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleEMR_SetVisibleToAllUsers() {
	sess, err := session.NewSession()
	if err != nil {
		fmt.Println("failed to create session,", err)
		return
	}

	svc := emr.New(sess)

	params := &emr.SetVisibleToAllUsersInput{
		JobFlowIds: []*string{ // Required
			aws.String("XmlString"), // Required
			// More values...
		},
		VisibleToAllUsers: aws.Bool(true), // Required
	}
	resp, err := svc.SetVisibleToAllUsers(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleEMR_TerminateJobFlows() {
	sess, err := session.NewSession()
	if err != nil {
		fmt.Println("failed to create session,", err)
		return
	}

	svc := emr.New(sess)

	params := &emr.TerminateJobFlowsInput{
		JobFlowIds: []*string{ // Required
			aws.String("XmlString"), // Required
			// More values...
		},
	}
	resp, err := svc.TerminateJobFlows(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}
