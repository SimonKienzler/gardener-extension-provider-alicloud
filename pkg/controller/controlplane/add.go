// Copyright (c) 2019 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package controlplane

import (
	"sync/atomic"

	extensionscontroller "github.com/gardener/gardener/extensions/pkg/controller"
	"github.com/gardener/gardener/extensions/pkg/controller/controlplane"
	"github.com/gardener/gardener/extensions/pkg/controller/controlplane/genericactuator"
	"github.com/gardener/gardener/extensions/pkg/util"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	"github.com/gardener/gardener-extension-provider-alicloud/pkg/alicloud"
	"github.com/gardener/gardener-extension-provider-alicloud/pkg/apis/config"
	"github.com/gardener/gardener-extension-provider-alicloud/pkg/imagevector"
)

var (
	// DefaultAddOptions are the default AddOptions for AddToManager.
	DefaultAddOptions = AddOptions{}
)

// AddOptions are options to apply when adding the Alicloud controlplane controller to the manager.
type AddOptions struct {
	// Controller are the controller.Options.
	Controller controller.Options
	// IgnoreOperationAnnotation specifies whether to ignore the operation annotation or not.
	IgnoreOperationAnnotation bool
	// ShootWebhookConfig specifies the desired Shoot MutatingWebhooksConfiguration.
	ShootWebhookConfig *atomic.Value
	// CSI specifies the configuration for CSI components
	CSI config.CSI
}

// AddToManagerWithOptions adds a controller with the given Options to the given manager.
// The opts.Reconciler is being set with a newly instantiated actuator.
func AddToManagerWithOptions(mgr manager.Manager, opts AddOptions) error {
	return controlplane.Add(mgr, controlplane.AddArgs{
		Actuator: genericactuator.NewActuator(alicloud.Name, secretConfigsFunc, shootAccessSecretsFunc, nil, nil,
			nil, controlPlaneChart, controlPlaneShootChart, controlPlaneShootCRDsChart, storageClassChart, nil,
			NewValuesProvider(opts.CSI), extensionscontroller.ChartRendererFactoryFunc(util.NewChartRendererForShoot),
			imagevector.ImageVector(), "", opts.ShootWebhookConfig, mgr.GetWebhookServer().Port),
		ControllerOptions: opts.Controller,
		Predicates:        controlplane.DefaultPredicates(opts.IgnoreOperationAnnotation),
		Type:              alicloud.Type,
	})
}

// AddToManager adds a controller with the default Options.
func AddToManager(mgr manager.Manager) error {
	return AddToManagerWithOptions(mgr, DefaultAddOptions)
}
