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

type ValidatorResult interface {
}
func (t api_v2024.Transform) ValidatorResult() {}
type ValidatorString string
func (s ValidatorString) ValidatorResult() {}

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
		return transform,nil
	}
	return nil,fmt.Errorf("Unknown action: %s", action)
}

func main() {
	// sailpointURL := os.Getenv("INPUT_SAILPOINT_URL")
	// clientID := os.Getenv("INPUT_CLIENT_ID")
	// clientSecret := os.Getenv("INPUT_CLIENT_SECRET")
	// action := os.Getenv("INPUT_ACTION")
	// payload := os.Getenv("INPUT_PAYLOAD")


	// TESTING VALUES REMOVE BEFORE PUSHING

	// set up the context for the api
	ctx := context.TODO()
	// this is building a custom configuration object. This could be replaced with using the default ENV configuration
	configuration := sailpoint.NewConfiguration(config(sailpointURL,clientID,clientSecret))
	apiClient := sailpoint.NewAPIClient(configuration)
	configuration.HTTPClient.RetryMax = 10

	// get the list of sources for validation if necessary
	sources, _, err := apiClient.V2024.SourcesAPI.ListSources(ctx).Execute()
	source_info := make(map[string]string)

	if err != nil {
		fmt.Println(err)
	}

	// Loop sources and get the id for each source
	for _, source := range sources {
		// append the id and name to a map of source ids and names
		source_info[source.GetId()] = source.Name
		
	}

	// DEBUG
	fmt.Println(source_info)
	// END DEBUG


	// Transforms
	if action == "create-transform" {

		// pass the payload to the validate function to ensure the transform is valid
		transform,err := validator(action,payload)

		if transform != nil {
			fmt.Println("Transform is invalid")
			fmt.Println(transform)
		}

		// actually create the transform
		transformObj, resp, err := apiClient.V2024.TransformsAPI.CreateTransform(ctx).Transform(transform).Execute()
		if err != nil {
			fmt.Println(resp,err)
		}
		fmt.Print("Transform created successfully")
		fmt.Print(transformObj.Id)



	}

	// Sources

	// Identity Profiles
}