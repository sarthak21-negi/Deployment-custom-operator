# Description 
A custom Kubernetes Operator built using Golang, Kubebuilder, and controller-runtime that automatically scales Kubernetes Deployments based on configurable time schedules.

The operator watches custom resources and dynamically updates Deployment replica counts during active and inactive hours to optimize cluster resource utilization.

---

# Features

- Custom Kubernetes Operator built with Kubebuilder
- Time-based autoscaling for Kubernetes Deployments
- Custom Resource Definition (CRD) for scaling schedules
- Declarative scaling configuration
- Automatic scale-up and scale-down
- Reconciliation loop using controller-runtime
- Status updates for active schedules
- RBAC-enabled controller permissions
- Tested locally using Kind cluster

---

# Tech Stack

- Golang
- Kubernetes
- Kubebuilder
- controller-runtime
- Kind
- Docker

---

# Project Structure

```bash
deployment-custom-operator/
├── api/
│   └── v1alpha1/
├── cmd/
├── config/
│   ├── crd/
│   ├── rbac/
│   ├── manager/
│   └── samples/
├── internal/
│   └── controller/
├── Makefile
└── README.md
```

---

# Architecture

The operator follows the Kubernetes Operator pattern:

1. User creates a custom DeploymentCustomOperator resource
2. Controller watches the resource
3. Reconcile loop checks current time
4. Controller compares desired state vs actual state
5. Deployment replicas are updated automatically

---

# Custom Resource Example

```yaml
apiVersion: dpscaler.sarthak.dev/v1alpha1
kind: DeploymentCustomOperator

metadata:
  name: deploymentcustomoperator-sample
  namespace: default

spec:
  targets:
    - name: nginx-demo
      namespace: default
      replicas: 5

  schedule:
    startHour: 9
    endHour: 18
    defaultReplicas: 1
```

---

# How It Works

- During active hours (9 AM – 6 PM), the Deployment scales to 5 replicas
- Outside active hours, the Deployment scales back to 1 replica
- The controller continuously monitors the cluster state and reconciles resources every minute

---

# Prerequisites

Make sure the following are installed:

- Go >= 1.24
- Docker
- kubectl
- Kind
- Kubebuilder

---

# Setup Instructions

## 1. Clone Repository

```bash
git clone https://github.com/sarthak21-negi/deployment-custom-operator.git

cd deployment-custom-operator
```

---

## 2. Create Kind Cluster

```bash
kind create cluster --name dpscaler
```

---

## 3. Install CRD

```bash
make install
```

---

## 4. Run Controller

```bash
make run
```

---

## 5. Deploy Test Application

Create deployment:

```yaml
apiVersion: apps/v1
kind: Deployment

metadata:
  name: nginx-demo

spec:
  replicas: 1

  selector:
    matchLabels:
      app: nginx-demo

  template:
    metadata:
      labels:
        app: nginx-demo

    spec:
      containers:
        - name: nginx
          image: nginx:latest
          ports:
            - containerPort: 80
```

Apply deployment:

```bash
kubectl apply -f nginx-deployment.yaml
```

---

## 6. Apply Custom Resource

```bash
kubectl apply -f config/samples/dpscaler_v1alpha1_deploymentcustomoperator.yaml
```

---

# Verify Scaling

Watch deployment replicas:

```bash
kubectl get deployment nginx-demo -w
```

---

# Controller Logic

The controller:

- Watches DeploymentCustomOperator resources
- Reads scaling schedules
- Checks current system hour
- Fetches target Deployments
- Updates replica counts dynamically
- Updates CRD status fields

---

# Status Fields

The operator updates status information:

```yaml
status:
  active: true
  lastScaleTime: "2026-05-11T10:00:00Z"
```

---

# Testing The Controller

Use the make test command to run tests

```bash
make test
```

<img width="1903" height="246" alt="make test" src="https://github.com/user-attachments/assets/9fc1c4a6-ba88-407b-b994-92b3ae47c85b" />

---

# ScreenShots

Operator in action

<img width="1920" height="599" alt="Screenshot make run" src="https://github.com/user-attachments/assets/8c10b537-27b9-49cb-9429-a939e11cd3c5" />

---

<img width="1301" height="112" alt="Screenshot (16)" src="https://github.com/user-attachments/assets/3ba29cd9-d40b-4dce-8d98-ef2f641f69c0" />

