## sprots-news-storage

### Purpose
This is an (incomplete) implementation on a golang technical task.

### Building & Running
`vendor` is attached so you shouldn't need to build or pull deps
#### Locally
```bash
$ go build -o api cmd/api/main.go
```
```bash
$ go run cmd/api/main.go 
```
or 
```bash
./api
```

#### With Docker
The only dependency of the api is a mongoDB instance.  
Bring one up with:
```bash
$ docker-compose up -d mongo
```
With a mongo instance running, build and start the api docker container with:
```bash
$ docker-compose up -d api
```
When all services are up and running, run:
```bash
curl -s localhost:8082/v1/articles
```
to get a list of all stored articles.  

Similary, hit:
```bash
$ curl -s localhost:8082/v1/article/{ID}
```
to get a single news article.

*Note, for this to work, the news fetcher needs run first so that newly fetched articles are stored in the db.  
The default news articles fetch interval is 15seconds, which is very aggressive, but it's set so that you don't  
wait 15mins for new articles to be stored into the db.  

Yes, the api could preload db data by getting and storing the news once when the app starts and then on regular intervals to avoid the problem of waiting. 
This feature got dropped because of time limitations :).

#### Missing Features
Unfortunately, I didn't have the time to focus more than a day on this (as per your suggestion), so the task is incomplete.
Here's what is missing or could be improved:  

* Adding support for extending an article's metadata by retrieving more details about the article from the  
https://www.brentfordfc.com/api/incrowd/getnewsarticleinformation?id={ARTICLE_ID} endpoint. This doesn't involve anything non-trivial.  
The way the app is structured it can be easily extended to support this.
* Adding support for pagination. At the moment there's no limit on how many articles are returned in the /v1/articles response.  
This is a perfect use case for pagination which can be implemented using mongodb cursors underneath.


#### Tests & ITs
The repo includes a unit and integration tests. Obviously this was done to the extend of time availability.  
More tests can easily be added but there should be enough to cover the most major cases.  
*Mocks for the mongodb repository were generated with mockgen.  
*Integration tests require a docker container of mongodb to be running. You can facilitate that with:  
```bash
$ docker-compose up -d mongo
```
then to run all tests:
```bash
$ go test ./...
```

#### Personal remarks about the implementation
* Ideally, I'd move the news articles periodic sync module into a different module with its own main and hence have it become a separate go app  
that is solely responsible for fetching news articles periodicallly & updating the database.
* Secondly, I'd improve the costly and unnecessary BulkInsert operation. Since the articles don't change frequently,  
we can easily store a view in memory and on every sync interval performe a set operation for find only the new ones that need adding.  
This logic won't work though if the news articles are frequently updated and therefore need to be upserted.