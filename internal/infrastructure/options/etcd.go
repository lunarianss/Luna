// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package options

import (
	"fmt"

	"github.com/spf13/pflag"
)

// EtcdOptions defines options for etcd cluster.
type EtcdOptions struct {
	Endpoints            []string `json:"endpoints"               mapstructure:"endpoints"`
	Timeout              int      `json:"timeout"                 mapstructure:"timeout"`
	RequestTimeout       int      `json:"request-timeout"         mapstructure:"request-timeout"`
	LeaseExpire          int      `json:"lease-expire"            mapstructure:"lease-expire"`
	Username             string   `json:"username"                mapstructure:"username"`
	Password             string   `json:"password"                mapstructure:"password"`
	UseTLS               bool     `json:"use-tls"                 mapstructure:"use-tls"`
	CaCert               string   `json:"ca-cert"                 mapstructure:"ca-cert"`
	Cert                 string   `json:"cert"                    mapstructure:"cert"`
	Key                  string   `json:"key"                     mapstructure:"key"`
	HealthBeatPathPrefix string   `json:"health_beat_path_prefix" mapstructure:"health_beat_path_prefix"`
	HealthBeatIFaceName  string   `json:"health_beat_iface_name"  mapstructure:"health_beat_iface_name"`
	Namespace            string   `json:"namespace"               mapstructure:"namespace"`
}

// NewEtcdOptions create a `zero` value instance.
func NewEtcdOptions() *EtcdOptions {
	return &EtcdOptions{}
}

// Validate verifies flags passed to RedisOptions.
func (o *EtcdOptions) Validate() []error {
	errs := []error{}

	if len(o.Endpoints) == 0 {
		errs = append(errs, fmt.Errorf("etcd endpoints can not be empty"))
	}

	if o.RequestTimeout <= 0 {
		errs = append(errs, fmt.Errorf("--etcd.request-timeout cannot be negative"))
	}

	return errs
}

// AddFlags adds flags related to redis storage for a specific api-server to the specified FlagSet.
func (o *EtcdOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringSliceVar(&o.Endpoints, "etcd.endpoints", o.Endpoints, "Endpoints of etcd cluster.")
	fs.StringVar(&o.Username, "etcd.username", o.Username, "Username of etcd cluster.")
	fs.StringVar(&o.Password, "etcd.password", o.Password, "Password of etcd cluster.")
	fs.IntVar(&o.Timeout, "etcd.timeout", o.Timeout, "Etcd dial timeout in seconds.")
	fs.IntVar(&o.RequestTimeout, "etcd.request-timeout", o.RequestTimeout, "Etcd request timeout in seconds.")
	fs.IntVar(&o.LeaseExpire, "etcd.lease-expire", o.LeaseExpire, "Etcd expire timeout in seconds.")
	fs.BoolVar(&o.UseTLS, "etcd.use-tls", o.UseTLS, "Use tls transport to connect etcd cluster.")
	fs.StringVar(&o.CaCert, "etcd.ca-cert", o.CaCert, "Path to cacert for connecting to etcd cluster.")
	fs.StringVar(&o.Cert, "etcd.cert", o.Cert, "Path to cert file for connecting to etcd cluster.")
	fs.StringVar(&o.Key, "etcd.key", o.Key, "Path to key file for connecting to etcd cluster.")
	fs.StringVar(
		&o.HealthBeatPathPrefix,
		"etcd.health-beat-path-pre",
		o.HealthBeatPathPrefix,
		"health beat path prefix.",
	)
	fs.StringVar(
		&o.HealthBeatIFaceName,
		"etcd.health-beat-iface-name",
		o.HealthBeatIFaceName,
		"health beat registry iface name, such as eth0.",
	)
	fs.StringVar(&o.Namespace, "etcd.namespace", o.Namespace, "Etcd storage namespace.")
}
