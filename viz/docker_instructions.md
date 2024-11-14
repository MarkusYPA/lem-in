## Lem-In - How to dockerize the visualizer

\
Building with the tag "leminvisualizer":
```bash
docker build -t leminvisualizer .
```

\
Run a container with port 8080 and pipe input:
```bash
cat test01out.txt | docker run -i -p 8080:8080 leminvisualizer
```


## Cleaning up afterwards

\
Stop all containers:
```bash
docker stop $(docker ps -a -q)
```


\
Cleaning up stopped containers, unused networks, and dangling images:
```bash
docker system prune -f
```

\
Removing unused volumes:
```bash
docker volume prune -f
```

\
Removing all unused images, containers, networks, and volumes:\
(Warning: -a will remove ALL images not associated with a running container)
```bash
docker system prune -a -f
```
