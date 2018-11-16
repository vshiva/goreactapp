# GoReactApp - Go and React JS micro service that serves static content and also hosts REST Apis

## Application Stack

    - Go 1.12
    - ReactJS
    - Kaniko for Container Builder
    - Argo WF for CI & CD

## Building this App

    - argo submit ci.ayml

## Deploying this App

### Assumptions

- Initial Version of the application is deployed outside of this.
- Service is fixed.
- Deployment has a TAG to indicate it is blue or green.
- Parameters
    -- Overlay

### Deployment Steps

#### Step 1

Query the current version from the service tag and find out if we need to do blue or green deployment

#### Step 2

Do data migration or any presteps

#### Step 3

Use the overlay and generate the artifact with prefix that is green or blue

#### Step 4

Apply the artifacts that was generated earlier. At this step we will have a newer version of the app deployed

#### Step 5

Run the Test against the new service end point

#### Step 6

    Update the public service LB with new end points ( Simply another K8S Service )

#### Step 7

Scale down the old deployment to 0

#### Step 8

Handle Any Post deployment activies here