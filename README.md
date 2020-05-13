# WeeklyCodingStats

This is a weekly coding stats
📊   ascii bar chart generator which uses the WAKATIME API to keep track of your coding stats.
I have seen lot of people have this awesome bar chart with their weekly development breakdown.So I thought why not just try to implement it myself.so that I can learn more about API's and Authorization in http requests and stuff like that (And also cus Iam bored af 😐).
There are some ready to use tools like [waka-box github action](https://github.com/marketplace/actions/waka-box) to do this.
 

The Program uses GitHub Gists API to update a Gist after getting your coding stats for the last 7 days from WAKATIME API.

The ASCII bar chart will look like this 👇


![image](asciibar.png)

<!-- To set this up yourself you must have go installed and also should have wakatime IDE plugins tracking your coding activity:
Then Follow these steps :

1) Create a Gist and Get the Gist ID 
2) Clone this Repo
3) Set the environment variables
   1) GIST_TOKEN - Generate a personal access token with permissions to gists from Github settings
   2) GIST_ID - The id of the gist into which we are updating the coding stats into.
   3) WAKATIME_API_KEY - If you are using the api its required to set this environemnt variable
   4) WAKATIME_EMBED_URL - If you are not using the API instead using the EMBED URL set this env variable

After setting these env variables

If you are going to use API then

`go run main.go -api`

Or if you are usng the json embed url you got from the share section of Wakatime Dashboard

`go run main.go` -->
