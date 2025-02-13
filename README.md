# metal-device-plugin
A simple k8s-device-plugin that mimics the NVIDIA plugin but for the Apple Metal specification

# Building the plugin for the appropriate microprocessor architecture
Run the plugin as follows (For e.g. set the following environment variables as appropriate. Set GOOS=linux, GOARCH=arm64 as the docker container on M1/M2/M3 would be aarch64. Use go build -o metal-device-plugin to build the plugin for the appropriate microprocessor architecture. Remember to set GO111MODULE to on if you're relying on traditional go build -o.

# Validating the plugin
1. Setup a simple minikube cluster after starting Docker locally on your workstation
2. minikube start


# Copy the metal-device-plugin (built above) using minikube scp 
3. minikube scp metal-device-plugin docker@$(minikube ip):/home/docker/ 

# Check for presence of the metal-gpu socket file (inside the container).
docker@minikube:~$ sudo ls -ld /var/lib/kubelet/device-plugins/metal-gpu.sock

# Run the metal-device-plugin locally from inside the docker container (Run as sudo if required)
4. minikube ssh
5. ./metal-device-plugin

# If all went well, you will see the following output
Metal Device Plugin running...

# Now run your kubectl commands as follows
6. kubecelt describe node minikube

You should see a line like this:

quest1.io/gpu      0           0

# Now run the following command to look for the GPU (relying on the metal interface)
7. kubectl describe node minikube > minikube.before.allocation.txt
Depending on your MacBook configuration, you will find something like this:

Listed under Capacity:<BR>
quest1.io/gpu:      1

Listed under Allocatable:<BR>
quest1.io/gpu:      1

Listed under Allocated resources:<BR>
quest1.io/gpu      0           0

# Now run the following command to allocate the GPU (relying on the metal interface)
8. kubectl describe node minikube | grep -i allocatable -A 10 > minikube.after.allocation.txt 

And voila, you should see something like this<BR>
Allocatable:<BR>
  cpu:                16<BR>
  ephemeral-storage:  61202244Ki<BR>
  hugepages-1Gi:      0<BR>
  hugepages-2Mi:      0<BR>
  hugepages-32Mi:     0<BR>
  hugepages-64Ki:     0<BR>
  memory:             8028596Ki<BR>
  pods:               110<BR>
  quest1.io/gpu:      1<BR>
System Info:<BR>
--<BR>
  Normal  NodeAllocatableEnforced  56m                kubelet          Updated Node Allocatable limit across pods<BR>
  Normal  Starting                 56m                kubelet          Starting kubelet.<BR>
  Normal  NodeAllocatableEnforced  56m                kubelet          Updated Node Allocatable limit across pods<BR>
  Normal  NodeHasSufficientMemory  56m                kubelet          Node minikube status is now: NodeHasSufficientMemory<BR>
  Normal  NodeHasNoDiskPressure    56m                kubelet          Node minikube status is now: NodeHasNoDiskPressure<BR>
  Normal  NodeHasSufficientPID     56m                kubelet          Node minikube status is now: NodeHasSufficientPID<BR>
  Normal  RegisteredNode           56m                node-controller  Node minikube event: Registered Node minikube in Controller<BR>


