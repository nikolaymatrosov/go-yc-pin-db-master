package main

import (
	"context"
	"fmt"
	"os"

	"github.com/yandex-cloud/go-genproto/yandex/cloud/mdb/mysql/v1"
	"github.com/yandex-cloud/go-genproto/yandex/cloud/mdb/postgresql/v1"
	ycsdk "github.com/yandex-cloud/go-sdk"
)

func checkPostgres(ctx context.Context, sdk *ycsdk.SDK, cluster ClusterConfig) error {
	var (
		master  *postgresql.Host
		targets []*postgresql.Host
	)

	hostIter := sdk.MDB().PostgreSQL().Cluster().ClusterHostsIterator(ctx, cluster.ClusterId)

	for hostIter.Next() {
		host := hostIter.Value()
		if host.Role == postgresql.Host_MASTER {
			master = host
			if master.ZoneId == cluster.TargetAZ {
				return nil
			}
		}
		if host.ZoneId == cluster.TargetAZ && host.Health == postgresql.Host_ALIVE {
			targets = append(targets, host)
		}
	}
	if len(targets) == 0 {
		return fmt.Errorf("no awailable target hosts")
	}
	_, err := sdk.MDB().PostgreSQL().Cluster().StartFailover(ctx, &postgresql.StartClusterFailoverRequest{
		ClusterId: cluster.ClusterId,
		HostName:  targets[0].Name,
	})
	if err != nil {
		return err
	}

	return nil
}
func checkMySql(ctx context.Context, sdk *ycsdk.SDK, cluster ClusterConfig) error {
	var (
		master  *mysql.Host
		targets []*mysql.Host
	)

	hostIter := sdk.MDB().MySQL().Cluster().ClusterHostsIterator(ctx, cluster.ClusterId)

	for hostIter.Next() {
		host := hostIter.Value()
		if host.Role == mysql.Host_MASTER {
			master = host
			if master.ZoneId == cluster.TargetAZ {
				return nil
			}
		}
		if host.ZoneId == cluster.TargetAZ && host.Health == mysql.Host_ALIVE {
			targets = append(targets, host)
		}
	}
	if len(targets) == 0 {
		return fmt.Errorf("no awailable target hosts")
	}
	_, err := sdk.MDB().MySQL().Cluster().StartFailover(ctx, &mysql.StartClusterFailoverRequest{
		ClusterId: cluster.ClusterId,
		HostName:  targets[0].Name,
	})
	if err != nil {
		return err
	}

	return nil
}

func PinHandler(ctx context.Context) error {
	cluster := ClusterConfig{
		DbType:    os.Getenv("DB_TYPE"),
		ClusterId: os.Getenv("CLUSTER_ID"),
		TargetAZ:  os.Getenv("TARGET_AZ"),
	}

	sdk, err := ycsdk.Build(ctx, ycsdk.Config{
		// Вызов InstanceServiceAccount автоматически запрашивает IAM-токен и формирует
		// при помощи него данные для авторизации в SDK
		Credentials: ycsdk.InstanceServiceAccount(),
	})
	if err != nil {
		return err
	}

	switch cluster.DbType {
	case "postgres":
		return checkPostgres(ctx, sdk, cluster)
	case "mysql":
		return checkMySql(ctx, sdk, cluster)
	default:
		return fmt.Errorf("unknown db type")
	}
}
