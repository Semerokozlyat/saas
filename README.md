## General information and purpose

Screenshot-as-a-Service (SaaS) application.

Makes website screenshots according to URL provided, stores image data in local storage (it is Redis currently), returns screenshots by demand.
 
## Requirements

- linux (preferably)
- golang 1.12
- docker
- docker-compose
- internet connection (redis image is to be downloaded)
 
## Installation

1. Clone repository folder into `$GOPATH/src/saas`
2. Change current directory to `$GOPATH/src/saas`
3. Modify docker-compose.yaml to customize:

 - listening port, 8000 by default
 - local path, mounted as volume here. I.e. `/home/eugene/Pictures/screens` is the directory on local machine 
 where screenshots will be saved after they requested. 
 ```
     volumes:
       - /home/eugene/Pictures/screens:/screens 
```

4. Run the command:

`` docker-compose up --build -d ``

5. Make sure service is started:

```
[]$ docker-compose ps
       Name                      Command               State           Ports         
-------------------------------------------------------------------------------------
saas_redis_1          docker-entrypoint.sh redis ...   Up      0.0.0.0:6379->6379/tcp
saas_saas-service_1   /opt/saas/saas                   Up      0.0.0.0:8000->8000/tcp
 ```

6. When service is not required anymore:

`` docker-compose down --volumes ``

## Workflow (how to use the service)

Install, launch, make sure `[ ]$ curl -k -X GET localhost:8000/status` response code is 200.

Send request to make screenshot. Comma-separated URLs are supported in this way:
```
[ ]$ curl -k -X GET localhost:8000/make_screenshot?url=https://stopgame.ru,https://drom.ru

Request accepted. Screenshot filename will be: 1572519469.0.png
Request accepted. Screenshot filename will be: 1572519469.1.png

```

*What happens in app*: handler sends request to "put" channel, screenshot is made by headless chrome and saved to both
directory and Redis storage as array of bytes.

Send request to receive link to your stored screenshot:
```
[ ]$ curl -k -X GET localhost:8000/get_screenshot?file_name=1572519469.1.png
Screenshot is uploaded to share.
 You may download it by using this link: http://localhost:8000/screens/1572519469.1.png, its size is 3 kB
```

*What happens in app*: handler sends request to "get" channel, screenshot is retrieved from Redis (no matter if 
it is already removed from directory), saved to directory, direct link is created

## Debug

Watch docker container logs:

```
[ ]$ docker logs -f saas_saas-service_1
2019/10/30 23:54:29 Creating Redis storage
2019/10/30 23:54:29 Redis storage connection: tcp/redis:6379 db 1
2019/10/30 23:54:29 Redis storage created
2019/10/30 23:54:29 Service started
...
```

## Not-yet solved problems

Even though `--disable-gpu` parameter is used in headless chrome commandline, on my PC it crashes from time to time 
 and returns white pictures instead of website screenshots.
 
## Possible improvements

 - Add more intelligent request parsing, change methods to POST - to provide request parameters in body
 - Move to k8s for better scalability (current storage scheme supports HA by means of several k8s pods, storage logic is 
 separated from functional logic, any kind of storage could be used due to common interface method)
 - Replace headless chrome with something more stable (crashes too often)
 - Make File Server robust and secure