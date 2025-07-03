
# CloudLib - S3 File Retrieval Service

## Overview
A serverless API that retrieves text files from an S3 bucket via Lambda. Architecture:

`API Gateway (GET /text)  ->  Lambda Function  ->  S3 Bucket`

HTTP Endpoint: 

`GET /text?file={file-name}`

## Deployment

### 1. Build Lambda Function

```bash
make build-CloudLibFunction
sam build --use-container
```

### 2. Deploy with SAM

If you have aws API keys set locally you can deploy directly. 

```bash
sam deploy --guided
```

If you are like me and prefer to use `cloudshell` upload an archive of the project to `cloudshell` and then `build` + `deploy` directly in the aws shell. 

A successful deploy will print the endpoint - `https://${PubApi}.execute-api.${AWS::Region}.amazonaws.com/Prod/text`

### 3. Test functionality 

First you need to create a file in `S3`. Then simply run:

``` bash
curl `https://${PubApi}.execute-api.${AWS::Region}.amazonaws.com/Prod/text?file={file-name}`
```





