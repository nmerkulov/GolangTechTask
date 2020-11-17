Migrations done by https://github.com/golang-migrate/migrate#cli-usage

Database - postgresql, docker-compose file attached

Requests

1. `GET /stream/{stream_id}`

    That returns a single stream along with the buffs associated with that stream.
    
2. `POST /stream`
   
    ```bash
     curl -X POST -d '{"name":"test1", "buff_ids": [1,2]}' localhost:8080/streams -D -
    ```
   
   
unfortunately i won't have much time this week to complete pagination, testing and things like logging or config management,
but i'll simply explain how it should be implemmented


1) Pagination - there are few options available, - offsetting and "search_after"
Using offseting is fairly straightforward but is inefficient. For small amount of data (thousands, not millions) it
probably won't be an issue. But when it comes to more-or-less big amount of data - it will be slow.
Reason - postgress will literally iterate over output to find out what is offset. from 0 to 500000 element

Another approach - search after. In pagination response we return some specific value, lets say - ID.
And to get next page we need to pass this search-after-id parameter. This allows us to use database more efficient

2) Config management. I'm not a big fan of "config" packages and things like that.
I prefer, first of all, dependency injection principle. By saing that, i mean lets say GStore should know only
PG connecction string, it doesn't require any Config struct. Thats one story 
Another thing is..environment variables are much more easy to manage. I prefer to follow 12 factor apps recomendations
https://12factor.net/
And in such context there are no need in thihngs like config management. You can literally parse anythging you need from environment variable

Of course i understand that there are a lot of complicated things in industry. Like...get configuration from things like
etcd/zookeeper/vault. Or you want to adopt spf13/   Viper. Thats also okay to me, i'm flexible

3) Structured logging. Well. Logging is never easy. I'll share this link https://dave.cheney.net/2015/11/05/lets-talk-about-logging and say
that i aggree with mr. Cheney. And will add few bullets from my own
* You can start with standard log.Logger. It is not best solution but easy. Disadvantages - as far as i remember it is not buffered so can
lead to efficiency drops (because of IO)
* You can use whatever things you like - structlog, uber.Zap - whatever. Just don't blow stdout with EVERYTHING if Debug set to False
* In modern world logging quite often should be traceable. So things like opentracing/jaeger are very important to use
* And, of course, metrics. It is inderctly related to logging, so i want to highlight it here - app need measurements to
be successfully used in producton


4) As for layout - i'm flexible in a way that different teams composing their applications. So i didn't assign new package to postgres code, let me explain why

It is about 100 LoC. Thats it. only one small file. So it felt like okay. Bt what is actually important - low coupling and high cohesion.
These two guys are extremely important when building software. And all those tricks with layouts are exists to fit these
requirements. Layouts can be different - layered with storage in the heart of app, hex/onion/clean with domain in the heart -
i can do whatever you want until we support low coupling and high cohesion


5) Testing
I agree that test i wrote isn't idiomatic. The most idiomatic go-way of unit-testig is table-tests or closure-tests (which is pretty much the same but 
looks better https://medium.com/@cep21/closure-driven-tests-an-alternative-style-to-table-driven-tests-in-go-628a41497e5e)