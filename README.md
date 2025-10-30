# To get the backend working:

- You will need a spotify dev account along with the tokens that get generated
- When running the program for the first time it will generate a file: ./userdata/data.json
- data.json -> v1 will need to be updated with the client and secret (token is automatically retrieved)
- data.json -> v2 is much trickier as it uses the private spotify APIs.

## How to get v2 spotify apis:

- Get a http scraper such as [fiddler] (https://www.telerik.com/fiddler).
- Have the scraper running and then open up spotify (free or premium, doesn't matter)
- Search for an artist and open their profile
- Any spotify request should have the headers, but /pathfinder/v2/query contains the data returned to the spotify client
- The headers of interest are: client-token and authorization
- Copy paste this data into the v2 section of data.json. It has a max expiry of 1 hour

## How to run backend

go run .

If running for first time you can update the v2 token information with:
`settokens`

Otherwise you can update the data.json file directly.

To get an artists information:
`getartist Muse`

This will pull all album information for Muse. Live albums will be filtered out. Note this will pull all albums, and sometimes spotify marks data incorrectly and you may end up with more albums than expected. Same goes for rereleases of albums

To get an artists information based on set albums:
`get`

This uses the information located at ./artistdata/\_input.json. This will get the entire artist catalogue from spotify then attempt to match based on the information.

If an exact match can be found then that album is saved. Otherwise it will try to grab based on releaseyear and a partial match of the album name. Note that spotify is not accurate when it comes to release years and it can be off (December releases especially bad). Also spotify may have a different album name compared to the original release. Rereleases can also have their release dates changed to when it was rereleased rather than the actual release date. See Johnny Cash as an example of terrible data.

If you need to regrab the data, then update the \_input.json file and change processed=true to false. Then delete the retrieved artist data in the ./artistdata directory

# To display the FE

It's best to download the echarts.js manually and place it in the /display directory. There is the option to download from their site, but you will need to comment out the code

## To run frontend

`npx http-server -c-1`

Data is hardcoded to the artist's name. This was to get a clean 1080p screenshot when the browser was in full screenmode (f11)
Simply update artistName to be the one you want. If the data is not there, then a blank white screen will be displayed

## Tips on how to get artist info

Spotify does not distinguish albums. They are either EPs or albums, so compilations and rereleases get caught up together. To avoid this, use some AI tool (chatGPT, grok, etc) to get the artist's studio albums. This is far more consistent than relying off spotify exclusively. I have included an example prompt I use at: example prompt.txt
