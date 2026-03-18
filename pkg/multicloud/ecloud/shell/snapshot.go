package shell

import (
	"yunion.io/x/pkg/util/shellutils"

	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/cloudmux/pkg/multicloud/ecloud"
)

func init() {
	type DiskSnapshotListOptions struct {
		DiskId string `help:"Disk ID"`
	}
	shellutils.R(&DiskSnapshotListOptions{}, "disk-snapshot-list", "List snapshots of a data disk", func(cli *ecloud.SRegion, args *DiskSnapshotListOptions) error {
		snapshots, err := cli.GetSnapshots("", args.DiskId, false)
		if err != nil {
			return err
		}
		printList(snapshots)
		return nil
	})

	type ServerSnapshotListOptions struct {
		ServerId string `help:"Server ID"`
	}
	shellutils.R(&ServerSnapshotListOptions{}, "server-snapshot-list", "List snapshots of a system disk by server", func(cli *ecloud.SRegion, args *ServerSnapshotListOptions) error {
		snapshots, err := cli.GetSnapshots("", args.ServerId, true)
		if err != nil {
			return err
		}
		printList(snapshots)
		return nil
	})

	type SnapshotShowOptions struct {
		SnapshotId string `help:"Snapshot ID"`
		DiskId     string `help:"Disk ID (for data disk snapshot)" optional:"true"`
		ServerId   string `help:"Server ID (for system disk snapshot)" optional:"true"`
	}
	shellutils.R(&SnapshotShowOptions{}, "snapshot-show", "Show snapshot detail", func(cli *ecloud.SRegion, args *SnapshotShowOptions) error {
		isSystem := args.ServerId != ""
		parentId := args.DiskId
		if isSystem {
			parentId = args.ServerId
		}
		snapshots, err := cli.GetSnapshots(args.SnapshotId, parentId, isSystem)
		if err != nil {
			return err
		}
		if len(snapshots) == 0 {
			return cloudprovider.ErrNotFound
		}
		printObject(&snapshots[0])
		return nil
	})

	type SnapshotCreateOptions struct {
		DiskId      string `help:"Disk ID"`
		Name        string `help:"Snapshot name"`
		Description string `help:"Description" optional:"true"`
	}
	shellutils.R(&SnapshotCreateOptions{}, "snapshot-create", "Create snapshot for disk", func(cli *ecloud.SRegion, args *SnapshotCreateOptions) error {
		snapshotId, err := cli.CreateEbsSnapshot(args.DiskId, args.Name, args.Description)
		if err != nil {
			return err
		}
		snapshots, err := cli.GetSnapshots(snapshotId, args.DiskId, false)
		if err != nil || len(snapshots) == 0 {
			return err
		}
		printObject(&snapshots[0])
		return nil
	})

	type SnapshotDeleteOptions struct {
		ID string `help:"Snapshot ID"`
	}
	shellutils.R(&SnapshotDeleteOptions{}, "snapshot-delete", "Delete snapshot", func(cli *ecloud.SRegion, args *SnapshotDeleteOptions) error {
		return cli.DeleteEbsSnapshot(args.ID)
	})
}

