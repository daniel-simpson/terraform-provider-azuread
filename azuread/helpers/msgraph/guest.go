package msgraph

import "fmt"

// Guest contains guest information returned from Microsoft Graph API
type Guest struct {
	id                string
	displayName       string
	mail              string
	userType          string
	userPrincipalName string
}

// GuestInvitationRequest describes the structure of data to be sent to the /invitations graph API endpoint
type GuestInvitationRequest struct {
	InvitedUserEmailAddress string `json:"invitedUserEmailAddress"`
	InviteRedirectURL       string `json:"inviteRedirectUrl"`
	SendInvitationMessage   bool   `json:"sendInvitationMessage"`
}

type invitedUser struct {
	id string
}

// GuestInvitationResponse describes the structure of data to be received from the /invitations graph API endpoint
type GuestInvitationResponse struct {
	invitedUser invitedUser
}

// GetGuest retrieves a Guest's details from the MS graph API
func (c *Client) GetGuest(id string) (Guest, error) {
	url := fmt.Sprintf("/users/%s?$select=id,displayName,mail,userType,userPrincipalName", id)
	req, requestError := c.newRequest("GET", url, nil)

	if requestError != nil {
		return Guest{}, requestError
	}

	response, responseError := c.httpClient.Do(req)

	if responseError != nil {
		return Guest{}, responseError
	}

	defer response.Body.Close()

	guestDetails := Guest{}

	parseError := parseResponse(response, &guestDetails)

	if parseError != nil {
		return Guest{}, parseError
	}

	return guestDetails, nil
}

// InviteGuest invites a user into a tenant from the MS graph API
func (c *Client) InviteGuest(email string) (string, error) {
	url := "/invitations"

	payloadBuffer, payloadBufferError := getPayloadBuffer(GuestInvitationRequest{
		InvitedUserEmailAddress: email,
		InviteRedirectURL:       "https://portal.azure.com",
		SendInvitationMessage:   true,
	})

	if payloadBufferError != nil {
		return "", payloadBufferError
	}

	req, requestError := c.newRequest("POST", url, payloadBuffer)

	if requestError != nil {
		return "", requestError
	}

	response, responseError := c.httpClient.Do(req)

	if responseError != nil {
		return "", responseError
	}

	defer response.Body.Close()

	guestInvite := GuestInvitationResponse{}

	parseError := parseResponse(response, &guestInvite)

	if parseError != nil {
		return "", parseError
	}

	return guestInvite.invitedUser.id, nil
}

// DeleteGuest removes a Guest user from a tenant.  Fails if the targeted user ID is not a "guest" userType
func (c *Client) DeleteGuest() error {
	// TODO: get user details and confirm that they are a "Guest" type

	c.newRequest("DELETE", "", nil)

	return nil
}
