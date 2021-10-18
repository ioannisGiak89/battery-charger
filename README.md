# Battery Charger

A microservice to control our batteries so that we can keep the grid low-carbon.

## Requirements 

1. Go 1.16

## Installation

Open a terminal and navigate to the project's directory. Make sure the env is set up(see [Environment Variables](#Environment Variables)) and run:

```bash
    go install -v ./...
```

To execute the app run

```bash
    arenko
```

## Run Locally

Open a terminal and navigate to the project's directory. Make sure the env is set up(see [Environment Variables](#Environment Variables)) and run:

```bash
    go run .
```

## Running Tests

To run tests, make sure the env is set up(see [Environment Variables](#Environment Variables)) and run: the following command

```bash
   go test ./...
```

## Assumptions

1) I based the decision for setting the charging rate on the index field of the intensity data. 
    If the value of the index is low or very low I charge the batteries. If it's high or very high I discharge them. Otherwise, 
    I just wait for the next value.
2) Although the National Grid API documentation says that the intensity endpoint returns the Intensity data for the current half hour, I noticed 
    that the `from` and `to` fields refer to the last half hour. Here is an example where I log the time of my machine as well as the response. 
    Both times are in UTC:
    
    ```.env
    UTC: 2021-07-27 10:37:09.664434007 +0000 UTC
    Response: {Data:[{From:2021-07-27T10:00Z To:2021-07-27T10:30Z Intensity:{Forecast:245 Actual:245 Index:moderate}}]}
    Carbon Intensity is moderate. Let's wait..
    
    UTC: 2021-07-27 11:07:09.974287103 +0000 UTC
    Response: {Data:[{From:2021-07-27T10:30Z To:2021-07-27T11:00Z Intensity:{Forecast:240 Actual:253 Index:moderate}}]}
    Carbon Intensity is moderate. Let's wait..
    ```
    
    Then, I looked at the `GET /intensity/{from}` endpoint. Both `from` and `to` fields were more accurate, however the actual intensity was 0:
    
    ```.env
    UTC: 2021-07-27 10:04:05.990849284 +0000 UTC
    Response: {Data:[{From:2021-07-27T10:00Z To:2021-07-27T10:30Z Intensity:{Forecast:244 Actual:0 Index:moderate}}]}
    Carbon Intensity is moderate. Let's wait..
    
    UTC: 2021-07-27 10:34:08.119501983 +0000 UTC
    Response: {Data:[{From:2021-07-27T10:30Z To:2021-07-27T11:00Z Intensity:{Forecast:243 Actual:0 Index:moderate}}]}
    Carbon Intensity is moderate. Let's wait..  
    ```
   
    At this point I made the assumption that the first endpoint is accurate based on the National Grid docs
3) I decided to swallow some errors, so I don't interrupt the execution of the app. In a production environment I would
    create an alarm, so we can monitor these errors.

## Environment Variables

To run this project, you will need to add the following environment variables to your `.env` file. 

`NATIONAL_GRID_BASE_URL` default value: https://api.carbonintensity.org.uk/

`INTENSITY_ENDPOINT` default value: intensity

`INTENSITY_CHECK_INTERVAL_SECONDS` default value: 1800

*Note:* A `.env` file is included in the zip that I submitted. If it is missing please create one and populate the above values.
