// Copyright © 2020 Intel Corporation. All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause

package functions

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/edgexfoundry/go-mod-core-contracts/v2/dtos"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/dtos/requests"
)

// SubscribeToNotificationService configures an email notification and submits
// it to the EdgeX notification service
func (boardStatus CheckBoardStatus) SubscribeToNotificationService() error {
	subscriptionClient := boardStatus.Service.SubscriptionClient()
	if subscriptionClient == nil {
		return errors.New("error notification service missing from client's configuration")
	}

	dto := dtos.Subscription{
		Id:   uuid.NewString(),
		Name: boardStatus.Configuration.NotificationName,
		Channels: []dtos.Address{
			{
				Type:         "EMAIL",
				EmailAddress: dtos.EmailAddress{Recipients: boardStatus.Configuration.NotificationEmailAddresses},
			},
		},
		Receiver: boardStatus.Configuration.NotificationReceiver,
		Labels: []string{
			boardStatus.Configuration.NotificationCategory,
		},
		Categories: []string{
			boardStatus.Configuration.NotificationCategory,
		},
		AdminState: boardStatus.Configuration.SubscriptionAdminState,
	}
	reqs := []requests.AddSubscriptionRequest{requests.NewAddSubscriptionRequest(dto)}
	_, err := subscriptionClient.Add(context.Background(), reqs)
	if err != nil {
		return fmt.Errorf("failed to subscribe to the EdgeX notification service: %s", err.Error())
	}

	return nil
}

func (boardStatus CheckBoardStatus) SendNotification(message string) error {
	notificationClient := boardStatus.Service.NotificationClient()
	if notificationClient == nil {
		return errors.New("error notification service missing from client's configuration")
	}

	dto := dtos.NewNotification(boardStatus.Configuration.NotificationLabels,
		boardStatus.Configuration.NotificationCategory,
		message,
		boardStatus.Configuration.NotificationSender,
		boardStatus.Configuration.NotificationSeverity,
	)

	req := requests.NewAddNotificationRequest(dto)
	_, err := notificationClient.SendNotification(context.Background(), []requests.AddNotificationRequest{req})
	if err != nil {
		return fmt.Errorf("failed to send the notification: %s", err.Error())
	}

	return nil
}
