package main

import (
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/storage"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func createBucketObject(ctx *pulumi.Context, bucket *storage.Bucket, filePath string, contentType string) error {
	bucketObject, err := storage.NewBucketObject(ctx, filePath, &storage.BucketObjectArgs{
		Bucket:      bucket.Name,
		ContentType: pulumi.String(contentType),
		Source:      pulumi.NewFileAsset("src/" + filePath),
	})

	bucketEndpoint := pulumi.Sprintf("http://storage.googleapis.com/%s/%s", bucket.Name, bucketObject.Name)
	if err != nil {
		return err
	}

	ctx.Export("bucketEndpoint"+filePath, bucketEndpoint)
	return nil
}

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
			ForceDestroy:             pulumi.Bool(true),
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

		createBucketObject(ctx, bucket, "index.html", "text/html")
		createBucketObject(ctx, bucket, "404.html", "text/html")
		createBucketObject(ctx, bucket, "styles.css", "text/css")

		// Export the DNS name of the bucket
		ctx.Export("bucketName", bucket.Url)
		return nil
	})
}
