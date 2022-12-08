# annotation-based-deployment-scaler

- Simple controller PoC

- Reacts to Deployment events in the Default Namespace

- Sets their replica counts to 5 if not already 5 when
  `.metadata.annotations["foo"] = "bar"`
