package bingocloud

import (
	"time"

	"yunion.io/x/jsonutils"

	"yunion.io/x/cloudmux/pkg/cloudprovider"
)

type SInstanceBackup struct {
	region             *SRegion
	BackupId           string
	InstanceId         string
	BackupName         string
	DisplayName        string
	OwnerId            string
	BackupSize         int
	BackupStatus       string
	Progress           string
	StorageId          string
	BlockDeviceMapping string
	Description        string
}

func (self SInstanceBackup) GetId() string {
	return self.BackupId
}

func (self SInstanceBackup) GetName() string {
	return self.BackupName
}

func (self SInstanceBackup) GetGlobalId() string {
	return self.BackupId
}

func (self SInstanceBackup) GetCreatedAt() time.Time {
	return time.Now()
}

func (self SInstanceBackup) GetStatus() string {
	status := self.BackupStatus
	switch self.BackupStatus {
	case "complete":
		status = "completed"
	}
	return status
}

func (self SInstanceBackup) IsEmulated() bool {
	return false
}

func (self SInstanceBackup) GetSysTags() map[string]string {
	return nil
}

func (self SInstanceBackup) GetTags() (map[string]string, error) {
	return nil, nil
}

func (self SInstanceBackup) SetTags(tags map[string]string, replace bool) error {
	return nil
}

func (self SInstanceBackup) GetProjectId() string {
	return ""
}

func (self SInstanceBackup) GetDescription() string {
	return self.Description
}

func (self SInstanceBackup) Delete() error {
	return self.region.deleteInstanceBackup(self.BackupId)
}

func (self SInstanceBackup) Refresh() error {
	newBackups, err := self.region.getInstanceBackups(self.InstanceId, self.BackupId)
	if err != nil {
		return err
	}
	if len(newBackups) == 1 {
		return jsonutils.Update(self, &newBackups[0])
	}
	return cloudprovider.ErrNotFound
}

func (self *SRegion) createInstanceBackup(instanceId, name string, desc string) (string, error) {
	params := map[string]string{}
	params["InstanceId"] = instanceId
	params["BackupName"] = name
	params["Description"] = desc

	resp, err := self.invoke("BackupInstance", params)
	if err != nil {
		return "", err
	}
	newId := ""
	err = resp.Unmarshal(&newId, "backupInstanceResult", "backup", "backupId")

	return newId, err
}

func (self *SRegion) getInstanceBackups(instanceId, backupId string) ([]SInstanceBackup, error) {
	params := map[string]string{}
	if instanceId != "" {
		params["InstanceId.1"] = instanceId
	}

	resp, err := self.invoke("DescribeInstanceBackups", params)
	if err != nil {
		return nil, err
	}

	var ret []SInstanceBackup
	_ = resp.Unmarshal(&ret, "describeInstanceBackupsResult", "instanceBackups")

	if backupId != "" && ret != nil {
		for _, backup := range ret {
			if backupId == backup.BackupId {
				return []SInstanceBackup{backup}, nil
			}
		}
	}

	return ret, err
}

func (self *SRegion) restoreInstanceBackup(backupId string) error {
	params := map[string]string{}
	params["BackupId"] = backupId
	_, err := self.invoke("RestoreFromInstanceBackup", params)
	return err
}

func (self *SRegion) deleteInstanceBackup(id string) error {
	params := map[string]string{}
	params["BackupId"] = id
	_, err := self.invoke("DeleteInstanceBackup", params)
	return err
}
