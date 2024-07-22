# Example: Aggregating GPU Workload running on OpenShift AI Jupyter Notebook

The following instructions show step by step how to use SusQL to aggregate energy data
consumed by a GPU utilizing Jupyter notebook running on OpenShift AI.

## Prerequisites
The following are assumed to be installed and available.
- OpenShift Cluster (verified on version 4.14)
- OpenShift AI (verified on version 2.11)
- Kepler (verified with community version 0.13.)
- SusQL (verified with community version 0.0.22)
- Command line access to the OpenShift Cluster (e.g., logged in with `oc` command)
To use GPU functionality seamlessly within the cluster, a GPU and necessary software must be available:
- NVIDIA GPU Operator (verified with version 24.3)
- Node Feature Discovery Operator (verified with version 4.14)
- Kernel Module Management Operator (verified with version 2.1)

##  Create a Jupyter Notebook 

Any code that runs in a Jupyter notebook can be aggregated by SusQL. The following sample code
demonstrates the use of GPU resources. This is also a good test case to verify that GPU is
configured correctly.

```
pip install pycaret[full]

import torch
import time

if torch.cuda.is_available():
    device = torch.device('cuda')
else:
    device = torch.device('cpu')

matrix_size = 16384

x = torch.randn(matrix_size, matrix_size)
y = torch.randn(matrix_size, matrix_size)

x_gpu = x.to(device)
y_gpu = y.to(device)
torch.cuda.synchronize()

for i in range(10):
    start = time.time()
    result_gpu = torch.matmul(x_gpu, y_gpu)
    print("Run time using device",result_gpu.device,"is","{:.7f}".format(time.time() - start))
```

A Jupyter Notebook can be created and run through the following steps:

- Log into OpenShift
- Open the OpenShift AI Console by clicking the app tile, and then clicking "Red Hat OpenShift AI".
  - If this is your first time, or if you want a refresher, you can scroll down to the "Get oriented with learning resources" section 
    and follow the "Creating a Jupyter notebook" tutorial. Otherwise continue with the instructions below.
- Click "Applications" on the menu on the left hand side of the window. Click "Enabled". Then click "Launch Application" on the newly displayed tile.
  - The first time you will need to create a notebook server. These instructions were verified with the PyTorch image using CUDA v11.8, Python v3.9, and PyTorch v2.2.
- Within the Jupyter notebook server, click the "+" sign and click on the "Python 3.9" tile. Copy the code above into cells in the notebook.
  (It is probably best to put the `pip` command in its own cell to run once before running the rest of the Python code.)
- Verify that the code runs and uses the GPU.


## Attach a SusQL label to the Jupyter Notebook server:

Although the OpenShift Web Console can be used to set labels on existing workloads, this is also easy to do from the command line:

The following command removes a SusQL label on the Jupyter Notebook Server pod, in case one happens to be defined.
```
$ oc label pod $(oc get po -n rhods-notebooks | grep jupyter | head -1 |  cut -f 1 -d" ") -n rhods-notebooks "susql.label/1-"
pod/jupyter-nb-kube-3aadmin-0 unlabeled
```

Next, this command sets the label `Susql.label/1` to `openshiftaij` for the Jupyter notebook server running in namespace rhods-notebooks.
```
$ oc label pod $(oc get po -n rhods-notebooks | grep jupyter | head -1 |  cut -f 1 -d" ") -n rhods-notebooks "susql.label/1=openshiftaij"
pod/jupyter-nb-kube-3aadmin-0 labeled
```

And, finally, this command can verify that the label has been set
```
$ oc describe pod $(oc get po -n rhods-notebooks | grep jupyter | head -1 |  cut -f 1 -d" ") -n rhods-notebooks | grep -i susql
                  susql.label/1=openshiftaij
```

## Create the SusQL LabelGroup

First create a LabelGroup definition file called `openshiftaij.yaml` as follows:
```
---
apiVersion: susql.ibm.com/v1
kind: LabelGroup
metadata:
    name: openshiftaij
    namespace: rhods-notebooks
spec:
    labels:
        - openshiftaij
---
```
And apply the file:
```
$ oc apply -f openshiftaij.yaml
labelgroup.susql.ibm.com/openshiftaij created
```

## Visualize 

- Go back to the Jupyter Notebook Server, and run the workload.
 (Possibly change the number of times the program loops to a very large number to make the energy usage obvious.)
- Then on the OpenShift Web GUI, click "Observe" on the left hand side menu, then click "Metrics" and enter the following search
query, and click "Run queries":
```
susql_total_energy_joules{susql_label_1="openshiftaij"}
```

If you have cloned the GitHub `susql-operator` repository, you could also run the `test/susqltop` command to view energy aggregation from the command line.

```
$ test/susqltop
NameSpace           LabelGroup                              Labels                                                    TotalEnergy (J)
rhods-notebooks     openshiftaij                            ["openshiftaij"]                                          17963.00
```
