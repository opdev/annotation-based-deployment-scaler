# annotation-based-deployment-scaler

## controllers/deployment_controller.go
- Simple controller PoC

- Reacts to Deployment events in the Default Namespace

- Sets their replica counts to 5 if not already 5 when
  `.metadata.annotations["foo"] = "bar"`

## controllers/filtered_deployment_controller.go

- Simple Predicate PoC

- Reacts to Deployment create/update events in the default namespace, but
  only if they have the `.metadata.annotations["scaleme"] = "please"`