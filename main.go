package main

import (

	"context"
	flag "github.com/spf13/pflag"

	"github.com/charmbracelet/log"
	sailpoint "github.com/sailpoint-oss/golang-sdk/v2"
	"github.com/sailpoint-oss/golang-sdk/v2/api_v2024"

	"github.com/peepnoop/validators"
)

var (
	sailpointURL, clientID, clientSecret, action, payload string
)



func config(sailpointURL string,clientID string,clientSecret string) sailpoint.ClientConfiguration {

	var simpleConfig sailpoint.ClientConfiguration
	simpleConfig.BaseURL = sailpointURL
	simpleConfig.ClientId = clientID
	simpleConfig.ClientSecret = clientSecret
	simpleConfig.TokenURL = sailpointURL + "/oauth/token"

	return simpleConfig
}


func main() {

	// Validate the input
	parseAndValidateInput()

	// set up the context for the api
	ctx := context.TODO()
	// this is building a custom configuration object. This could be replaced with using the default ENV configuration
	configuration := sailpoint.NewConfiguration(config(sailpointURL,clientID,clientSecret))
	apiClient := sailpoint.NewAPIClient(configuration)
	configuration.HTTPClient.RetryMax = 10


	// Set up validators factory
	factory := validators.NewValidatorFactory()


	// Actions to preform
	if action == "create-transform" {

		// pass the payload to the validate function to ensure the transform is valid
		result,err := factory.Validate(action,payload)

		if err != nil {
			log.Error("Failed to validate transform", "err", err)
		}

		// assert the result to the transform type
		if transform, ok := result.(api_v2024.Transform); ok {

			// create the transform
			transformObj, _, err := apiClient.V2024.TransformsAPI.CreateTransform(ctx).Transform(transform).Execute()
			if err != nil {
				log.Error("Failed to create transform", "err", err)
			}
			log.Info("Transform created successfully", "transform", transformObj.Id)
		}
	} else if action == "update-transform"{
		// pass the payload to the validate function to ensure the transform is valid
		result,err := factory.Validate(action,payload)

		if err != nil {
			log.Error("Failed to validate transform", "err", err)
		}

		// assert the result to the transform type
		if transform, ok := result.(api_v2024.TransformRead); ok {

			updateTransform := api_v2024.Transform{Attributes: transform.Attributes, Type: transform.Type, Name: transform.Name}

			// actually create the transform
			transformObj, _, err := apiClient.V2024.TransformsAPI.UpdateTransform(ctx,transform.Id).Transform(updateTransform).Execute()
			if err != nil {
				log.Error("Failed to update transform", "err", err)
			}
			log.Info("Transform updated successfully", "transform", transformObj.Id)
		}
	} else if action == "delete-transform" {
		// attempt to delete the transform
		resp, err := apiClient.V2024.TransformsAPI.DeleteTransform(ctx, payload).Execute()

		if err != nil {
			log.Error("Failed to delete transform", "err", err)
		}

		log.Info("Transform deleted successfully", "transform", resp.StatusCode)
	} else if action == "create-role" {
		// pass the payload to the validate function to ensure the transform is valid
		result,err := factory.Validate(action,payload)

		if err != nil {
			log.Error("Failed to validate role", "err", err)
		}

		// assert the result to the role type
		if role, ok := result.(api_v2024.Role); ok {
			// create the role
			roleObj, _, err := apiClient.V2024.RolesAPI.CreateRole(ctx).Role(role).Execute()
			if err != nil {
				log.Error("Failed to create role", "err", err)
			}
			log.Info("Transform created successfully", "role", roleObj.Id)
		}
	} else if action == "update-role" {
		log.Error("Not implemented", "action", action)
	} else if action == "delete-role" {
		log.Error("Not implemented", "action", action)
	} else {
		log.Error("Unknown action", "action", action)
	}
}

func parseAndValidateInput() {

	flag.StringVar(&sailpointURL, "sailpoint-url", "", "Sailpoint URL")
	flag.StringVar(&clientID, "client-id", "", "Client ID")
	flag.StringVar(&clientSecret, "client-secret", "", "Client Secret")
	flag.StringVar(&action, "action", "", "Action to perform")
	flag.StringVar(&payload, "payload", "", "Payload for the action")
	flag.Parse()


	if sailpointURL == "" {
		log.Fatal("Sailpoint URL is required")
	}
	if clientID == "" {
		log.Fatal("Client ID is required")
	}
	if clientSecret == "" {
		log.Fatal("Client Secret is required")
	}
	if action == "" {
		log.Fatal("Action is required")
	}
}