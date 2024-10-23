package main

import (
	"context"
	"fmt"
	"bufio"
	"encoding/json"
	"os"

	sailpoint "github.com/sailpoint-oss/golang-sdk/v2"
	"github.com/sailpoint-oss/golang-sdk/v2/api_v2024"
	"github.com/charmbracelet/log"

)

// ValidatorResult interface
type ValidatorResult interface {
	Unwrap() interface{}
}
// Define a wrapper type for api_v2024.Transform
type TransformWrapper struct {
    Transform api_v2024.Transform
}

type RoleWrapper struct {
	Role api_v2024.Role
}

// Implement ValidatorResult for TransformWrapper
func (tw TransformWrapper) Unwrap() interface{} {
    return tw.Transform
}

func (rw RoleWrapper) Unwrap() interface{} {
    return rw.Role
}

func config(sailpointURL string,clientID string,clientSecret string) sailpoint.ClientConfiguration {

	var simpleConfig sailpoint.ClientConfiguration
	simpleConfig.BaseURL = sailpointURL
	simpleConfig.ClientId = clientID
	simpleConfig.ClientSecret = clientSecret
	simpleConfig.TokenURL = sailpointURL + "/oauth/token"

	return simpleConfig
}

func validator(action string, payload string) (ValidatorResult, error){

	// check the type of action to see what needs to be validated

	// TRANSFORMS //
	if action == "create-transform" {
		// validate the payload

		var transform api_v2024.Transform
		var decoder *json.Decoder

		if payload != "" {
			file, err := os.Open(payload)
			if err != nil {
				return nil,err
			}
			defer file.Close()
			decoder = json.NewDecoder(bufio.NewReader(file))
		} else {
			decoder = json.NewDecoder(bufio.NewReader(os.Stdin))
		}

		if err := decoder.Decode(&transform); err != nil {
			return nil,err
		}

		log.Debug("Filepath", "path", payload)

		log.Debug("Transform", "transform", transform)

		if transform.GetName() == "" {
			return nil,fmt.Errorf("The transform must have a name")
		}
		return TransformWrapper{Transform: transform},nil
	}
	if action == "update-transform" {
		var transform api_v2024.TransformRead

			filepath := payload
			if filepath != "" {
				file, err := os.Open(filepath)
				if err != nil {
					return nil,err
				}
				defer file.Close()

				err = json.NewDecoder(file).Decode(&transform)
				if err != nil {
					return nil,err
				}
			} else {
				err := json.NewDecoder(os.Stdin).Decode(&transform)
				if err != nil {
					return nil,err
				}
			}

			if transform.Id == "" {
				return nil,fmt.Errorf("The transform must have an Id to update")
			}
	}

	// ROLES //
	if action == "create-role" {
		// validate the payload

		var role api_v2024.Role
		var decoder *json.Decoder

		if payload != "" {
			file, err := os.Open(payload)
			if err != nil {
				return nil,err
			}
			defer file.Close()
			decoder = json.NewDecoder(bufio.NewReader(file))
		} else {
			decoder = json.NewDecoder(bufio.NewReader(os.Stdin))
		}

		if err := decoder.Decode(&role); err != nil {
			return nil,err
		}

		log.Debug("Filepath", "path", payload)
		log.Debug("Role", "role", role)

		if role.GetName() == "" {
			return nil,fmt.Errorf("The role must have a name")
		}
		return RoleWrapper{Role: role},nil
	}
	return nil,fmt.Errorf("Unknown action: %s", action)
}

func main() {
	sailpointURL := os.Getenv("INPUT_SAILPOINT_URL")
	clientID := os.Getenv("INPUT_CLIENT_ID")
	clientSecret := os.Getenv("INPUT_CLIENT_SECRET")
	action := os.Getenv("INPUT_ACTION")
	payload := os.Getenv("INPUT_PAYLOAD")


	// TESTING VALUES REMOVE BEFORE PUSHING

	// set up the context for the api
	ctx := context.TODO()
	// this is building a custom configuration object. This could be replaced with using the default ENV configuration
	configuration := sailpoint.NewConfiguration(config(sailpointURL,clientID,clientSecret))
	apiClient := sailpoint.NewAPIClient(configuration)
	configuration.HTTPClient.RetryMax = 10


	// TESTING BLOCK NEEDS //
	// TO BE MOVED TO VALIDATION //


	// // get the list of sources for validation if necessary
	// sources, _, err := apiClient.V2024.SourcesAPI.ListSources(ctx).Execute()
	// source_info := make(map[string]string)

	// if err != nil {
	// 	log.Error("Failed to get sources", "err", err)
	// }

	// // Loop sources and get the id for each source
	// for _, source := range sources {
	// 	// append the id and name to a map of source ids and names
	// 	source_info[source.GetId()] = source.Name
		
	// }

	// // DEBUG
	// log.Debug("Sources", "sources", source_info)
	// // END DEBUG

	// END TESTING BLOCK //


	// Transforms
	if action == "create-transform" {

		// pass the payload to the validate function to ensure the transform is valid
		result,err := validator(action,payload)

		if err != nil {
			log.Error("Failed to validate transform", "err", err)
		}

		// Perform type assertion
		transformWrapper, ok := result.(TransformWrapper)
		if !ok {
			log.Fatal("Result is not of type api_v2024.Transform")
		}

		transform := transformWrapper.Transform

		// get the source id from the transform

		// actually create the transform
		transformObj, _, err := apiClient.V2024.TransformsAPI.CreateTransform(ctx).Transform(transform).Execute()
		if err != nil {
			log.Error("Failed to create transform", "err", err)
		}
		log.Info("Transform created successfully", "transform", transformObj.Id)

	}

	if action == "update-transform"{
		// pass the payload to the validate function to ensure the transform is valid
		result,err := validator(action,payload)

		if err != nil {
			log.Error("Failed to validate transform", "err", err)
		}

		// Perform type assertion
		transformWrapper, ok := result.(TransformWrapper)
		if !ok {
			log.Fatal("Result is not of type api_v2024.Transform")
		}

		transform := transformWrapper.Transform

		updateTransform := api_v2024.Transform{Attributes: transform.Attributes, Type: transform.Type, Name: transform.Name}

		// actually create the transform
		transformObj, _, err := apiClient.V2024.TransformsAPI.UpdateTransform(ctx,transform.Id).Transform(updateTransform).Execute()
		if err != nil {
			log.Error("Failed to update transform", "err", err)
		}
		log.Info("Transform updated successfully", "transform", transformObj.Id)

	}

	if action == "delete-transform" {
		// attempt to delete the transform
		resp, err := apiClient.V2024.TransformsAPI.DeleteTransform(ctx, payload).Execute()

		if err != nil {
			log.Error("Failed to delete transform", "err", err)
		}

		log.Info("Transform deleted successfully", "transform", resp.StatusCode)
	}

	// Sources

	// Identity Profiles

	// Roles
	if action == "create-role" {
		// pass the payload to the validate function to ensure the json is valid
		result,err := validator(action,payload)

		if err != nil {
			log.Error("Failed to validate Role json", "err", err)
		}

		// Perform type assertion
		roleWrapper, ok := result.(RoleWrapper)
		if !ok {
			log.Fatal("Result is not of type api_v2024.Role")
		}

		role := roleWrapper.Role

		// create the role
		roleObj, _, err := apiClient.V2024.RolesAPI.CreateRole(ctx).Role(role).Execute()
		if err != nil {
			log.Error("Failed to create role", "err", err)
		}
		log.Info("Transform created successfully", "role", roleObj.Id)
	} else {
		log.Error("Unknown action", "action", action)
	}
}