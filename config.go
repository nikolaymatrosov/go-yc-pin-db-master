package main

type ClusterConfig struct {
	DbType    string `yaml:"dbType"`
	ClusterId string `yaml:"clusterId"`
	TargetAZ  string `yaml:"targetAz"`
}
