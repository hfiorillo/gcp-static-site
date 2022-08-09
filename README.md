# Getting set up

1. Run pulumi new gcp-go on an empty directory to create the necessary scaffolding to get your project started.
2. Run `gcloud init` and follow the steps to configure your gcloud cli.
3. Run `gcloud auth application-default login` to authenticate using the default credentials. https://cloud.google.com/docs/authentication/production#auth-cloud-implicit-go
4. Run `pulumi up` from the root of your project.
