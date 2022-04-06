/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	windowv1alpha1 "github.com/stolostron/maintenance-window-operator/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("MaintenanceWindow controller", func() {
	const (
		MaintenanceWindowName = "test-maintenancewindow"

		timeout  = time.Second * 10
		duration = time.Second * 10
		interval = time.Millisecond * 250
	)

	Context("When updating MaintenanceWindow status", func() {
		It("Should transition Status.State from SCHEDULED to OPENED to finally CLOSED", func() {
			By("By creating a new MaintenanceWindow")
			ctx := context.Background()
			now := time.Now()
			maintenanceWindow := &windowv1alpha1.MaintenanceWindow{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "window.open-cluster-management.io/v1alpha1",
					Kind:       "MaintenanceWindow",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: MaintenanceWindowName,
				},
				Spec: windowv1alpha1.MaintenanceWindowSpec{
					Date:     now.Format("2006-01-02"),
					Time:     now.Format(time.Kitchen),
					TimeZone: "UTC",
					Duration: func(i int32) *int32 { return &i }(1),
				},
			}
			Expect(k8sClient.Create(ctx, maintenanceWindow)).Should(Succeed())

			maintenanceWindowLookupKey := types.NamespacedName{Name: MaintenanceWindowName}
			createdMaintenanceWindow := &windowv1alpha1.MaintenanceWindow{}

			Eventually(func() bool {
				err := k8sClient.Get(ctx, maintenanceWindowLookupKey, createdMaintenanceWindow)
				return err == nil
			}, timeout, interval).Should(BeTrue())

			Expect(createdMaintenanceWindow.Status.State).Should(Equal("SCHEDULED"))
		})
	})

})
