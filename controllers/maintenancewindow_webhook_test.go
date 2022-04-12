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

var _ = Describe("MaintenanceWindow webhook", func() {
	const (
		MaintenanceWindowName = "test-maintenancewindow-wh"
		timeout               = time.Second * 60
		interval              = time.Millisecond * 250
	)

	Context("When a MaintenanceWindow is updated", func() {

		It("Should fail and return an error message that MaintenanceWindow CRs cannot be updated", func() {
			By("By creating a new MaintenanceWindow")
			ctx := context.Background()
			aMinuteFromNow := time.Now().UTC().Add(2 * time.Minute)
			maintenanceWindow := &windowv1alpha1.MaintenanceWindow{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "window.open-cluster-management.io/v1alpha1",
					Kind:       "MaintenanceWindow",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: MaintenanceWindowName,
				},
				Spec: windowv1alpha1.MaintenanceWindowSpec{
					Date:     aMinuteFromNow.Format("2006-01-02"),
					Time:     aMinuteFromNow.Format(time.Kitchen),
					TimeZone: "UTC",
					Duration: func(i int32) *int32 { return &i }(60),
				},
			}
			Expect(k8sClient.Create(ctx, maintenanceWindow)).Should(Succeed())

			maintenanceWindowLookupKey := types.NamespacedName{Name: MaintenanceWindowName}
			createdMaintenanceWindow := &windowv1alpha1.MaintenanceWindow{}

			Eventually(func() bool {
				err := k8sClient.Get(ctx, maintenanceWindowLookupKey, createdMaintenanceWindow)
				return err == nil
			}, timeout, interval).Should(BeTrue())

			By("Updating the Date")
			createdMaintenanceWindow.Spec.Date = aMinuteFromNow.Add(24 * time.Hour).Format(time.Kitchen)
			Expect(k8sClient.Update(context.Background(), createdMaintenanceWindow)).ShouldNot(Succeed())

			By("Updating the Time")
			createdMaintenanceWindow.Spec.Time = aMinuteFromNow.Add(1 * time.Hour).Format("2006-01-02")
			Expect(k8sClient.Update(context.Background(), createdMaintenanceWindow)).ShouldNot(Succeed())

			By("Updating the TimeZone")
			createdMaintenanceWindow.Spec.TimeZone = "CET"
			Expect(k8sClient.Update(context.Background(), createdMaintenanceWindow)).ShouldNot(Succeed())

			By("Updating the Duration")
			createdMaintenanceWindow.Spec.Duration = func(i int32) *int32 { return &i }(120)
			Expect(k8sClient.Update(context.Background(), createdMaintenanceWindow)).ShouldNot(Succeed())
		})
	})

})
