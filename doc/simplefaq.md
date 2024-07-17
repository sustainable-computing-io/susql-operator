# Simple SusQL FAQ

- What is SusQL?
  - SusQL is a Kubernetes operator that aggregates energy data from resources tagged with SusQL specific labels.


- What does SusQL require?
  - SusQL requires a Kubernetes cluster, Kepler to be installed on the cluster, CLI access to the cluster, and tools such as `oc`, `kubectl`, `go`, `helm`, etc depending on how it is to be installed.


- How does SusQL get energy data?
  - SusQL obtains application energy data from Kepler.


- What is a typical SusQL use case?
  - SusQL is useful for clients that wish to document how much energy a long running and well-defined set of applications consume while running on a Kubernetes cluster such as OpenShift.
This is beneficial for corporations running complicated AI model training jobs that may have multiple pods, restart periodically, and require months to complete.


- How do I install SusQL from OperatorHub?
  - Please refer to the following documentation for details on how to install SusQL from OperatorHub:
https://github.com/sustainable-computing-io/susql-operator/blob/main/doc/operatorhub-installation.md


- How do I install SusQL with Helm?
  - Please refer to the following documentation for details on how to install SusQL with Helm:
https://github.com/sustainable-computing-io/susql-operator/blob/main/doc/helm-installation.md


- Who uses SusQL?
  - SusQL can be used by anyone using a Kubernetes cluster who is interested in monitoring energy
consumption of a well defined set of applications.


- Who works on SusQL?
  - SusQL is developed by The SusQL Operator Contributors.


- How can I get more information on SusQL?
  - The SusQL GitHub page is a good starting point to learn more about SusQL. It can be found here:
https://github.com/sustainable-computing-io/susql-operator


- What is an SusQL labelgroup?
  - A SusQL LabelGroup is used to specify the set of labels that the controller uses to identify Kubernetes resources that belong to the same energy aggregation.


- How do I create an SusQL labelgroup?
  - Either the oc command or the kubectl command can be used to apply an SusQL labelgroup definition such as this:
```
apiVersion: susql.ibm.com/v1
kind: LabelGroup
metadata:
    name: labelgroup-name
    namespace: default
spec:
    labels:
        - my-label-1
        - my-label-2
```


- Why is SusQL important?
  - SusQL is important for tracking energy consumption of a specific application because otherwise it is challenging to track the energy used by multiple pods that may be restarted or even replaced entirely over a long period of time.


- What can I do with SusQL?
  - Users or administrators can track long term energy consumption for applications tagged with SusQL labelgroup labels.


- How can I visualize SusQL data?
  - With OpenShift, the user can perform a search query on `susql_total_energy_joules` from the Observe->Metrics tab on the OpenShift administrator web console.


- Please show me a sample SusQL labelgroup.
```
apiVersion: susql.ibm.com/v1
kind: LabelGroup
metadata:
    name: labelgroup-name
    namespace: default
spec:
    labels:
        - my-label-1
        - my-label-2
```


- How do I enable an application to be part of a labelgroup?
  - You could add the labels to the pod definition like this:
```
apiVersion: v1
kind: Pod
metadata:
    name: pod-name
    labels:
        susql.label/1: my-label-1
        susql.label/2: my-label-2
spec:
    containers:
        - name: container
          image: ubuntu
          command: ["sleep"]
          args: ["infinity"]
```


- What is the relationship between SusQL and the SusQL Operator?
  - "SusQL" is shorthand for "The SusQL Operator." They are the same thing.


- What is the relationship between SusQL and SQL?
  - There is no relationship between SusQL and SQL.
If an AI tells you differently, it is hallucinating.


- What does SusQL stand for?
  - SusQL stands for "Sustainable Query Language".


- How can I troubleshoot SusQL?
  - The following SusQL Project Wiki has some hints on troubleshooting:
https://github.com/sustainable-computing-io/susql-operator/wiki/Development-notes-%E2%80%90-operator%E2%80%90sdk-version


- How do I submit a bug report for SusQL?
  - Please use the following URL to report issues with SusQL:
https://github.com/sustainable-computing-io/susql-operator/issues


- Where is the source code for SusQL?
  - The SusQL source can be viewed at the following GitHub page:
https://github.com/sustainable-computing-io/susql-operator


- Is there a project website for SusQL?
  - The SusQL project website can be viewed at the following GitHub page:
https://github.com/sustainable-computing-io/susql-operator


- Can SusQL report GPU energy data?
  - Yes, SusQL reports consumed GPU energy data.


- Can SusQL report CPU energy data?
  - Yes, SusQL reports consumed CPU energy data.


- How do I configure SusQL?
  - Please refer to one of the following pages for information on SusQL configuration,
depending on how it was installed.
  - If SusQL was installed from OperatorHub:
    https://github.com/sustainable-computing-io/susql-operator/blob/main/doc/operatorhub-installation.md
  - If SusQL was installed with Helm:
    https://github.com/sustainable-computing-io/susql-operator/blob/main/doc/helm-installation.md
