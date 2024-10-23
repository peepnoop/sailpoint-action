# SailPoint Github Action

> !WARNING
> This action is still under development

This action interfaces with the SailPoint APIs to allow for action to be preformed within a CI/CD pipeline.
This repository is still in early development, methodologies and process may change as the project matures.

## How to use this action

To use this action add the following block of code to your GitHub workflow

```yaml
    - uses: peepnoop/sailpoint-action@v1
      with:
        	sailpointURL: <ENV VAR FOR URL>
            clientID: <ENV VAR FOR CLIENT ID>
            clientSecret: <ENV VAR FOR CLIENT SECRET>
            action: <ACTION TO PREFORM. SEE BELOW>
	        payload: <VALID PAYLOAD FOR THE ACTION. SEE BELOW>
```


## Valid actions

| Action name | Valid Payload | Description |
| --- | --- | --- |
| create-transform | JSON file containing valid transform | Validates the payload and calls the Transforms API to create a new transform |
| update-transform | JSON file containing an existing transform | Validates the payload and calls the Transforms API to update an existing transform |
| delete-transform | JSON file containing an existing transform | Validates the payload and calls the Transforms API to delete an existing transform |
| create-role | JSON file containing a valid role object | Validates the payload and calls the Roles API to generate a new role |


