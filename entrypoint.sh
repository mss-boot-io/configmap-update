#!/bin/sh -l

export cluster_url=$1
export token=$2
export name=$3
export namespace=$4
export files=$5
export data=$6

/configmap-update --kubeconfig=$7