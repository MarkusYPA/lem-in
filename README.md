# lem-in

Lem-in moves ants across an ant farm in a most efficient way. A text file specifies the rooms and tunnels connecting them. No upper limit exists to how many ants can move during a turn, but each ant can move only once. Only the start and end rooms are allowed more than one ant at a time. Any tunnel can be used only once per turn. 


## How to run 

Lem-in requires a text file for information about the farm and the ants.
\
\
Try with example00.txt in the testcases folder:
```bash
 go run . testcases/example00.txt 
```

The program prints out it's input and the moves for the fewest amount of turns to move all the ants through the farm.


## Lem-in - How to run the visualizer with Docker

Lem-in comes with a visualizer program in the viz/ folder. The visualizer requiers GraphViz to be installed so it comes with a dockerfile to run it in a compatible docker image. Two script files are provided to make dockerizing the visualizer easy. Make sure Docker is installed and ready before running them.

\
Run the first one to create a docker image:
```bash
./dockerize_visualizer.sh
```

\
The other one runs a docker container with a texfile from the testcases/ folder in port 8080:
```bash
./visualize_test.sh example00
```

Navigate to [localhost:8080](localhost:8080) with a web browser to see an animation of the ants moving through the farm.


### Cleaning up after dockerizing

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
Removing all unused images, containers, networks, and volumes:\
(Warning: -a will remove ALL images not associated with a running container)
```bash
docker system prune -a -f
```

