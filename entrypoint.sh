#!/bin/sh -l

export name=$1
export namespace=$2
export files=$3
export data=$4

aws eks --region $5 update-kubeconfig --name $6

/configmap-update