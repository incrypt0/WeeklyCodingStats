# WeeklyCodingStats

This is a weekly coding stats
📊   ascii bar chart generator which uses the WAKATIME API to keep track of your coding stats.
I have seen lot of people have this awesome bar chart with their weekly development breakdown.So I thought why not just try it out.
I just implemented it from scratch in Go (There are lots of ready to use code for implementing this like [wakatime-metrics by Athul](https://github.com/athul/wakatime-metrics)) so that I can learn more about API's and Authorization in http requests and stuff like that (And also cus Iam bored af 😐).

The Program uses GitHub Gists API to update a Gist.

The ASCII bar chart will look like this 👇


![image](asciibar.png)

To set this up yourself you must have go installed and also should have wakatime IDE plugins tracking your coding activity:
Then Follow these steps :

1) Create a Gist and Get the Gist ID 
2) Clone this Repo
3) Set 
