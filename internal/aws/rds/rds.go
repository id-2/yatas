package rds

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/stangirard/yatas/internal/config"
	"github.com/stangirard/yatas/internal/logger"
	"github.com/stangirard/yatas/internal/types"
)

func GetListRDS(s *session.Session) []*rds.DBInstance {
	logger.Debug("Getting list of RDS instances")
	svc := rds.New(s)

	params := &rds.DescribeDBInstancesInput{}
	resp, err := svc.DescribeDBInstances(params)
	if err != nil {
		panic(err)
	}

	logger.Debug(fmt.Sprintf("%v", resp.DBInstances))
	return resp.DBInstances
}

func checkIfEncryptionEnabled(s *session.Session, instances []*rds.DBInstance, testName string, c *[]types.Check) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check types.Check
	check.Name = "RDS Encryption"
	check.Id = testName
	check.Description = "Check if RDS encryption is enabled"
	check.Status = "OK"
	svc := rds.New(s)
	for _, instance := range instances {
		params := &rds.DescribeDBInstancesInput{
			DBInstanceIdentifier: instance.DBInstanceIdentifier,
		}
		resp, err := svc.DescribeDBInstances(params)
		if err != nil {
			panic(err)
		}
		if *resp.DBInstances[0].StorageEncrypted == false {
			check.Status = "FAIL"
			status := "FAIL"
			Message := "RDS encryption is not enabled on " + *instance.DBInstanceIdentifier
			check.Results = append(check.Results, types.Result{Status: status, Message: Message})
		} else {
			status := "OK"
			Message := "RDS encryption is enabled on " + *instance.DBInstanceIdentifier
			check.Results = append(check.Results, types.Result{Status: status, Message: Message})
		}
	}
	*c = append(*c, check)
}

func checkIfBackupEnabled(s *session.Session, instances []*rds.DBInstance, testName string, c *[]types.Check) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check types.Check
	check.Name = "RDS Backup"
	check.Id = testName
	check.Description = "Check if RDS backup is enabled"
	check.Status = "OK"
	svc := rds.New(s)
	for _, instance := range instances {
		params := &rds.DescribeDBInstancesInput{
			DBInstanceIdentifier: instance.DBInstanceIdentifier,
		}
		resp, err := svc.DescribeDBInstances(params)
		if err != nil {
			panic(err)
		}
		if *resp.DBInstances[0].BackupRetentionPeriod == 0 {
			check.Status = "FAIL"
			status := "FAIL"
			Message := "RDS backup is not enabled on " + *instance.DBInstanceIdentifier
			check.Results = append(check.Results, types.Result{Status: status, Message: Message})
		} else {
			status := "OK"
			Message := "RDS backup is enabled on " + *instance.DBInstanceIdentifier
			check.Results = append(check.Results, types.Result{Status: status, Message: Message})
		}
	}
	*c = append(*c, check)
}

func checkIfAutoUpgradeEnabled(s *session.Session, instances []*rds.DBInstance, testName string, c *[]types.Check) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check types.Check
	check.Name = "RDS Minor Auto Upgrade"
	check.Id = testName
	check.Description = "Check if RDS minor auto upgrade is enabled"
	check.Status = "OK"
	svc := rds.New(s)
	for _, instance := range instances {
		params := &rds.DescribeDBInstancesInput{
			DBInstanceIdentifier: instance.DBInstanceIdentifier,
		}
		resp, err := svc.DescribeDBInstances(params)
		if err != nil {
			panic(err)
		}
		if *resp.DBInstances[0].AutoMinorVersionUpgrade == false {
			check.Status = "FAIL"
			status := "FAIL"
			Message := "RDS auto upgrade is not enabled on " + *instance.DBInstanceIdentifier
			check.Results = append(check.Results, types.Result{Status: status, Message: Message})
		} else {
			status := "OK"
			Message := "RDS auto upgrade is enabled on " + *instance.DBInstanceIdentifier
			check.Results = append(check.Results, types.Result{Status: status, Message: Message})
		}
	}
	*c = append(*c, check)
}

func checkIfRDSPrivateEnabled(s *session.Session, instances []*rds.DBInstance, testName string, c *[]types.Check) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check types.Check
	check.Name = "RDS Private"
	check.Id = testName
	check.Description = "Check if RDS private is enabled"
	check.Status = "OK"
	svc := rds.New(s)
	for _, instance := range instances {
		params := &rds.DescribeDBInstancesInput{
			DBInstanceIdentifier: instance.DBInstanceIdentifier,
		}
		resp, err := svc.DescribeDBInstances(params)
		if err != nil {
			panic(err)
		}
		if *resp.DBInstances[0].PubliclyAccessible == true {
			check.Status = "FAIL"
			status := "FAIL"
			Message := "RDS private is not enabled on " + *instance.DBInstanceIdentifier
			check.Results = append(check.Results, types.Result{Status: status, Message: Message})
		} else {
			status := "OK"
			Message := "RDS private is enabled on " + *instance.DBInstanceIdentifier
			check.Results = append(check.Results, types.Result{Status: status, Message: Message})
		}
	}
	*c = append(*c, check)
}

func RunRDSTests(s *session.Session, c *config.Config) []types.Check {
	var checks []types.Check
	instances := GetListRDS(s)
	config.CheckTest(c, "AWS_RDS_001", checkIfEncryptionEnabled)(s, instances, "AWS_RDS_001", &checks)
	config.CheckTest(c, "AWS_RDS_002", checkIfBackupEnabled)(s, instances, "AWS_RDS_002", &checks)
	config.CheckTest(c, "AWS_RDS_003", checkIfAutoUpgradeEnabled)(s, instances, "AWS_RDS_003", &checks)
	config.CheckTest(c, "AWS_RDS_004", checkIfRDSPrivateEnabled)(s, instances, "AWS_RDS_004", &checks)

	return checks
}
