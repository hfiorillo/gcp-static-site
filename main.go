package main

import (
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/storage"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Create a GCP resource (Storage Bucket)
		bucket, err := storage.NewBucket(ctx, "my-bucket", &storage.BucketArgs{
			Website: storage.BucketWebsiteArgs{
				MainPageSuffix: pulumi.String("index.html"),
				NotFoundPage:   pulumi.String("404.html"),
			},
			UniformBucketLevelAccess: pulumi.Bool(true),
			Location:                 pulumi.String("EU"),
		})
		if err != nil {
			return err
		}

		_, err = storage.NewBucketIAMBinding(ctx, "my-bucket-IAMbinding", &storage.BucketIAMBindingArgs{
			Bucket: bucket.Name,
			Role:   pulumi.String("roles/storage.objectViewer"),
			Members: pulumi.StringArray{
				pulumi.String("allUsers"),
			},
		})
		if err != nil {
			return err
		}

		bucketObject, err := storage.NewBucketObject(ctx, "index.html", &storage.BucketObjectArgs{
			Bucket:      bucket.Name,
			ContentType: pulumi.String("text/html"),
			Source:      pulumi.NewFileAsset("src/index.html"),
		})

		bucketEndpoint := pulumi.Sprintf("http://storage.googleapis.com/%s/%s", bucket.Name, bucketObject.Name)
		if err != nil {
			return err
		}

		bucketObject404, err := storage.NewBucketObject(ctx, "404.html", &storage.BucketObjectArgs{
			Bucket:      bucket.Name,
			ContentType: pulumi.String("text/html"),
			Source:      pulumi.NewFileAsset("src/404.html"),
		})

		bucketEndpoint404 := pulumi.Sprintf("http://storage.googleapis.com/%s/%s", bucket.Name, bucketObject404.Name)
		if err != nil {
			return err
		}

		bucketObjectStatic, err := storage.NewBucketObject(ctx, "styles.css", &storage.BucketObjectArgs{
			Bucket:      bucket.Name,
			ContentType: pulumi.String("text/css"),
			Source:      pulumi.NewFileAsset("src/styles.css"),
		})

		bucketEndpointStatic := pulumi.Sprintf("http://storage.googleapis.com/%s/%s", bucket.Name, bucketObjectStatic.Name)
		if err != nil {
			return err
		}

		// Export the DNS name of the bucket
		ctx.Export("bucketName", bucket.Url)
		ctx.Export("bucketEndpointIndex", bucketEndpoint)
		ctx.Export("bucketEndpoint404", bucketEndpoint404)
		ctx.Export("bucketEndpointStatic", bucketEndpointStatic)
		return nil
	})
}
