"""
SusQL testing
"""

# %% Modules

import os
import sys
import json
import yaml
import math
import time
import threading
import subprocess

from datetime import datetime
from cloud import *

# %% Functions

ACTIVE_CLEANING = True

def clean_appwrappers():
    """
    Checks for completed AppWrappers and deletes them if found
    """
    while ACTIVE_CLEANING:
        appwrappers = bash("kubectl get appwrappers --all-namespaces --no-headers")

        for appwrapper in appwrappers:
            namespace, job_name, _ = appwrapper.split()

            pods = bash("kubectl -n %s get pods --selector appwrapper.mcad.ibm.com=%s --no-headers" % (namespace, job_name))

            if len(pods) == 0:
                continue

            num_completed = 0

            for pod in pods:
                if "Completed" in pod:
                    num_completed += 1

            if num_completed == len(pods):
                bash("kubectl -n %s delete appwrapper %s" % (namespace, job_name))

# %% Main function

def main():
    """
    Main program
    """

    # %% Multiple jobs

    ACTIVE_CLEANING = True
    cleaning = threading.Thread(target = clean_appwrappers, daemon = True)
    cleaning.start()

    projects = []

    # Set 1
    projects.append(("job-1", (2, 1), ( 60,  90), "namespace-1", { "project-name": "model-1", "experiment-name": "exp-1", "group-name": "grp-1" }))

    # Set 2
    projects.append(("job-2", (1, 1), ( 30,   0), "namespace-1", { "project-name": "model-1", "experiment-name": "exp-2", "group-name": "grp-2" }))
    projects.append(("job-3", (2, 1), (120,   0), "namespace-1", { "project-name": "model-1", "experiment-name": "exp-2", "group-name": "grp-2" }))
    projects.append(("job-4", (3, 1), ( 60,  45), "namespace-1", { "project-name": "model-1", "experiment-name": "exp-2", "group-name": "grp-2" }))

    # Set 3
    projects.append(("job-5", (3, 1), ( 45, 120), "namespace-1", { "project-name": "model-1", "experiment-name": "exp-3", "group-name": "grp-1" }))

    # Set 4
    projects.append(("job-6", (8, 1), ( 30,   0), "namespace-1", { "project-name": "model-1", "experiment-name": "exp-4", "group-name": "grp-2" }))
    projects.append(("job-7", (4, 1), ( 45,   0), "namespace-1", { "project-name": "model-1", "experiment-name": "exp-4", "group-name": "grp-2" }))

    for job_name, (num_pods, num_gpus_per_pod), (run_time, wait_time), namespace, custom_labels in projects:
        print("Starting job '%s' in namespace '%s': %s" % (job_name, namespace, datetime.now().strftime("%Y-%m-%d %H:%M:%S")))

        create_podgang(job_name = job_name, \
                       namespace = namespace, \
                       priority = "default-priority", \
                       container_image = "ubuntu", \
                       num_pods = num_pods, \
                       num_cpus_per_pod = "500m", \
                       num_gpus_per_pod = num_gpus_per_pod, \
                       total_memory_per_pod = "1Gi", \
                       shell_commands = ["sleep %d" % (run_time)], \
                       custom_labels = custom_labels, \
                       apply_yaml = True)

        time.sleep(wait_time)

    # Clean appwrappers
    while True:
        appwrappers = bash("kubectl get appwrappers --all-namespaces --no-headers")

        if len(appwrappers) == 0:
            break

    ACTIVE_CLEANING = False

    print("Finished: %s" % (datetime.now().strftime("%Y-%m-%d %H:%M:%S")))

    # %% Test job

    create_podgang(job_name = "long-job-1", \
                   namespace = "org-1", \
                   priority = "default-priority", \
                   container_image = "ubuntu", \
                   num_pods = 4, \
                   num_cpus_per_pod = "500m", \
                   num_gpus_per_pod = 2, \
                   total_memory_per_pod = "1Gi", \
                   shell_commands = ["sleep infinity"], \
                   custom_labels = { "organization-name": "org-1", "group-name": "grp-1", "project-name": "model-1" }, \
                   apply_yaml = True)

    create_podgang(job_name = "long-job-2", \
                   namespace = "org-2", \
                   priority = "default-priority", \
                   container_image = "ubuntu", \
                   num_pods = 4, \
                   num_cpus_per_pod = "500m", \
                   num_gpus_per_pod = 2, \
                   total_memory_per_pod = "1Gi", \
                   shell_commands = ["sleep infinity"], \
                   custom_labels = { "organization-name": "org-2", "group-name": "grp-1", "project-name": "model-2" }, \
                   apply_yaml = True)

    create_podgang(job_name = "long-job-3", \
                   namespace = "org-1", \
                   priority = "default-priority", \
                   container_image = "ubuntu", \
                   num_pods = 4, \
                   num_cpus_per_pod = "500m", \
                   num_gpus_per_pod = 2, \
                   total_memory_per_pod = "1Gi", \
                   shell_commands = ["sleep infinity"], \
                   custom_labels = { "organization-name": "org-1", "group-name": "grp-1", "project-name": "model-2" }, \
                   apply_yaml = True)

    # %% SusQL operator development

    create_job(job_name = "job-1", \
               namespace = "default", \
               priority = "default-priority", \
               container_image = "ubuntu", \
               num_pods = 4, \
               num_cpus_per_pod = "0m", \
               num_gpus_per_pod = 8, \
               total_memory_per_pod = "0Gi", \
               shell_commands = ["sleep infinity"], \
               custom_labels = { "experiment": "exp-1", "group": "grp-1", "project": "model-1" }, \
               quota_labels = { "quota-tree": "A" }, \
               apply_yaml = True)

    # %%

    create_job(job_name = "job-2", \
               namespace = "default", \
               priority = "default-priority", \
               container_image = "ubuntu", \
               num_pods = 4, \
               num_cpus_per_pod = "0m", \
               num_gpus_per_pod = 8, \
               total_memory_per_pod = "0Gi", \
               shell_commands = ["sleep infinity"], \
               custom_labels = { "experiment": "exp-1", "group": "grp-1", "project": "model-2" }, \
               quota_labels = { "quota-tree": "A" }, \
               apply_yaml = True)

    create_job(job_name = "job-3", \
               namespace = "default", \
               priority = "default-priority", \
               container_image = "ubuntu", \
               num_pods = 4, \
               num_cpus_per_pod = "0m", \
               num_gpus_per_pod = 4, \
               total_memory_per_pod = "0Gi", \
               shell_commands = ["sleep infinity"], \
               custom_labels = { "experiment": "exp-1", "group": "grp-1", "project": "model-1" }, \
               quota_labels = { "quota-tree": "A" }, \
               apply_yaml = True)

    # %% Clean up

    #bash("oc -n org-1 delete appwrappers --all")
    #bash("oc -n org-2 delete appwrappers --all")
    #bash("oc -n namespace-1 delete appwrappers --all")

    bash("kubectl -n default delete appwrappers --all", stdout = True)

# %% Main program

if __name__ == "__main__":
    try:
        main()
    except KeyboardInterrupt:
        try:
            quit()
        except SystemExit:
            quit()

# %% End of program
