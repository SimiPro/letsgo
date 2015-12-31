package images

/*
import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	storageApi "google.golang.org/api/storage/v1"
	"google.golang.org/cloud/storage"
)

// bucket is a local cache of the app's default bucket name.
//var bucket string // or: var bucket = "<your-app-id>.appspot.com"

func createCloudStorageService() *GoogleCloudStorageHandler {

	client, err := google.DefaultClient(context.Background(), scope)
	if err != nil {
		fmt.Printf("Unable to get default client: %v", err)
	}
	service, err := storageApi.New(client)
	if err != nil {
		fmt.Printf("Unable to create storage service: %v", err)
	}

	// If the bucket already exists and the user has access, warn the user, but don't try to create it.
	bucketName := "simisDefaultBucket"
	var bucketli *storage.BucketHandle
	bucketli, err = service.Buckets.Get(bucketName).Do()
	if err == nil {
		fmt.Printf("Bucket %s already exists - skipping buckets.insert call.", bucketName)
	} else {
		// Create a bucket.
		bucketli, err := service.Buckets.Insert("lestgo-1145", &storageApi.Bucket{Name: bucketName}).Do()
		if err == nil {
			fmt.Printf("Created bucket %v at location %v\n\n", bucketli.Name, bucketli.SelfLink)
		} else {
			fatalf(service, "Failed creating bucket %s: %v", bucketName, err)
		}

	}

	storageService := &GoogleCloudStorageHandler{
		bucket: bucketli,
		client: client,
		w:      os.Stdout,
		ctx:    context.Context,
	}

	return storageService
}

// holds information needed to run the various functions.
type GoogleCloudStorageHandler struct {
	bucket *storage.BucketHandle
	client *storage.Client
	w      io.Writer
	ctx    context.Context
}

// createFile creates a file in Google Cloud Storage.
func (cloudStorage *GoogleCloudStorageHandler) createFile(fileName string) {
	fmt.Fprintf(cloudStorage.w, "Creating file /%v/%v\n", cloudStorage.bucket, fileName)

	wc := cloudStorage.bucket.Object(fileName).NewWriter(cloudStorage.ctx)
	wc.ContentType = "text/plain"
	wc.Metadata = map[string]string{
		"x-goog-meta-foo": "foo",
		"x-goog-meta-bar": "bar",
	}

	if _, err := wc.Write([]byte("abcde\n")); err != nil {
		fmt.Fprintf(cloudStorage.w, "createFile: unable to write data to bucket %q, file %q: %v", bucket, fileName, err)
		return
	}
	if _, err := wc.Write([]byte(strings.Repeat("f", 1024*4) + "\n")); err != nil {
		fmt.Fprintf(cloudStorage.w, "createFile: unable to write data to bucket %q, file %q: %v", bucket, fileName, err)
		return
	}
	if err := wc.Close(); err != nil {
		fmt.Fprintf(cloudStorage.w, "createFile: unable to close bucket %q, file %q: %v", bucket, fileName, err)
		return
	}
}

// readFile reads the named file in Google Cloud Storage.
func (cloudStorage *GoogleCloudStorageHandler) readFile(fileName string) {
	fmt.Fprintf(cloudStorage.w, cloudStorage.w, "\nAbbreviated file content (first line and last 1K):\n")

	rc, err := cloudStorage.bucket.Object(fileName).NewReader(cloudStorage.ctx)
	if err != nil {
		cloudStorage.errorf("readFile: unable to open file from bucket %q, file %q: %v", bucket, fileName, err)
		return
	}
	defer rc.Close()
	slurp, err := ioutil.ReadAll(rc)
	if err != nil {
		cloudStorage.errorf("readFile: unable to read data from bucket %q, file %q: %v", bucket, fileName, err)
		return
	}

	fmt.Fprintf(cloudStorage.w, "%s\n", bytes.SplitN(slurp, []byte("\n"), 2)[0])
	if len(slurp) > 1024 {
		fmt.Fprintf(cloudStorage.w, "...%s\n", slurp[len(slurp)-1024:])
	} else {
		fmt.Fprintf(cloudStorage.w, "%s\n", slurp)
	}
}

// copyFile copies a file in Google Cloud Storage.
func (cloudStorage *GoogleCloudStorageHandler) copyFile(fileName string) {
	copyName := fileName + "-copy"
	fmt.Fprintf(cloudStorage.w, "Copying file /%v/%v to /%v/%v:\n", bucket, fileName, bucket, copyName)

	obj, err := cloudStorage.client.CopyObject(cloudStorage.ctx, bucket, fileName, bucket, copyName, nil)
	if err != nil {
		fmt.Fprintf(cloudStorage.w, "copyFile: unable to copy /%v/%v to bucket %q, file %q: %v", bucket, fileName, bucket, copyName, err)
		return
	}

	cloudStorage.dumpStats(obj)
}

func (cloudStorage *GoogleCloudStorageHandler) dumpStats(obj *storage.ObjectAttrs) {
	fmt.Fprintf(cloudStorage.w, "(filename: /%v/%v, ", obj.Bucket, obj.Name)
	fmt.Fprintf(cloudStorage.w, "ContentType: %q, ", obj.ContentType)
	fmt.Fprintf(cloudStorage.w, "ACL: %#v, ", obj.ACL)
	fmt.Fprintf(cloudStorage.w, "Owner: %v, ", obj.Owner)
	fmt.Fprintf(cloudStorage.w, "ContentEncoding: %q, ", obj.ContentEncoding)
	fmt.Fprintf(cloudStorage.w, "Size: %v, ", obj.Size)
	fmt.Fprintf(cloudStorage.w, "MD5: %q, ", obj.MD5)
	fmt.Fprintf(cloudStorage.w, "CRC32C: %q, ", obj.CRC32C)
	fmt.Fprintf(cloudStorage.w, "Metadata: %#v, ", obj.Metadata)
	fmt.Fprintf(cloudStorage.w, "MediaLink: %q, ", obj.MediaLink)
	fmt.Fprintf(cloudStorage.w, "StorageClass: %q, ", obj.StorageClass)
	if !obj.DeletecloudStorage.IsZero() {
		fmt.Fprintf(cloudStorage.w, "Deleted: %v, ", obj.Deleted)
	}
	fmt.Fprintf(cloudStorage.w, "Updated: %v)\n", obj.Updated)
}

// statFile reads the stats of the named file in Google Cloud Storage.
func (cloudStorage *GoogleCloudStorageHandler) statFile(fileName string) {
	io.WriteString(cloudStorage.w, "\nFile stat:\n")

	obj, err := cloudStorage.bucket.Object(fileName).Attrs(cloudStorage.ctx)
	if err != nil {
		cloudStorage.errorf("statFile: unable to stat file from bucket %q, file %q: %v", bucket, fileName, err)
		return
	}

	cloudStorage.dumpStats(obj)
}

// createListFiles creates files that will be used by listBucket.
func (cloudStorage *GoogleCloudStorageHandler) createListFiles() {
	io.WriteString(dcloudStorage.w, "\nCreating more files for listbucket...\n")
	for _, n := range []string{"foo1", "foo2", "bar", "bar/1", "bar/2", "boo/"} {
		cloudStorage.createFile(n)
	}
}

// listBucket lists the contents of a bucket in Google Cloud Storage.
func (cloudStorage *GoogleCloudStorageHandler) listBucket() {
	io.WriteString(cloudStorage.w, "\nListbucket result:\n")

	query := &storage.Query{Prefix: "foo"}
	for query != nil {
		objs, err := cloudStorage.bucket.List(cloudStorage.ctx, query)
		if err != nil {
			cloudStorage.errorf("listBucket: unable to list bucket %q: %v", bucket, err)
			return
		}
		query = objs.Next

		for _, obj := range objs.Results {
			cloudStorage.dumpStats(obj)
		}
	}
}

func (cloudStorage *GoogleCloudStorageHandler) listDir(name, indent string) {
	query := &storage.Query{Prefix: name, Delimiter: "/"}
	for query != nil {
		objs, err := cloudStorage.bucket.List(cloudStorage.ctx, query)
		if err != nil {
			cloudStorage.errorf("listBucketDirMode: unable to list bucket %q: %v", bucket, err)
			return
		}
		query = objs.Next

		for _, obj := range objs.Results {
			fmt.Fprint(cloudStorage.w, indent)
			cloudStorage.dumpStats(obj)
		}
		for _, dir := range objs.Prefixes {
			fmt.Fprintf(cloudStorage.w, "%v(directory: /%v/%v)\n", indent, bucket, dir)
			cloudStorage.listDir(dir, indent+"  ")
		}
	}
}

// listBucketDirMode lists the contents of a bucket in dir mode in Google Cloud Storage.
func (cloudStorage *GoogleCloudStorageHandler) listBucketDirMode() {
	io.WriteString(cloudStorage.w, "\nListbucket directory mode result:\n")
	cloudStorage.listDir("b", "")
}

// dumpDefaultACL prints out the default object ACL for this bucket.
func (cloudStorage *GoogleCloudStorageHandler) dumpDefaultACL() {
	acl, err := cloudStorage.bucket.ACL().List(cloudStorage.ctx)
	if err != nil {
		cloudStorage.errorf("defaultACL: unable to list default object ACL for bucket %q: %v", bucket, err)
		return
	}
	for _, v := range acl {
		fmt.Fprintf(cloudStorage.w, "Scope: %q, Permission: %q\n", v.Entity, v.Role)
	}
}

// defaultACL displays the default object ACL for this bucket.
func (cloudStorage *GoogleCloudStorageHandler) defaultACL() {
	io.WriteString(cloudStorage.w, "\nDefault object ACL:\n")
	cloudStorage.dumpDefaultACL()
}

// putDefaultACLRule adds the "allUsers" default object ACL rule for this bucket.
func (cloudStorage *GoogleCloudStorageHandler) putDefaultACLRule() {
	io.WriteString(cloudStorage.w, "\nPut Default object ACL Rule:\n")
	err := cloudStorage.bucket.DefaultObjectACL().Set(cloudStorage.ctx, storage.AllUsers, storage.RoleReader)
	if err != nil {
		cloudStorage.errorf("putDefaultACLRule: unable to save default object ACL rule for bucket %q: %v", bucket, err)
		return
	}
	cloudStorage.dumpDefaultACL()
}

// deleteDefaultACLRule deleted the "allUsers" default object ACL rule for this bucket.
func (cloudStorage *GoogleCloudStorageHandler) deleteDefaultACLRule() {
	io.WriteString(cloudStorage.w, "\nDelete Default object ACL Rule:\n")
	err := cloudStorage.bucket.DefaultObjectACL().Delete(cloudStorage.ctx, storage.AllUsers)
	if err != nil {
		cloudStorage.errorf("deleteDefaultACLRule: unable to delete default object ACL rule for bucket %q: %v", bucket, err)
		return
	}
	cloudStorage.dumpDefaultACL()
}

// dumpBucketACL prints out the bucket ACL.
func (cloudStorage *GoogleCloudStorageHandler) dumpBucketACL() {
	acl, err := cloudStorage.bucket.ACL().List(cloudStorage.ctx)
	if err != nil {
		cloudStorage.errorf("dumpBucketACL: unable to list bucket ACL for bucket %q: %v", bucket, err)
		return
	}
	for _, v := range acl {
		fmt.Fprintf(cloudStorage.w, "Scope: %q, Permission: %q\n", v.Entity, v.Role)
	}
}

// bucketACL displays the bucket ACL for this bucket.
func (cloudStorage *GoogleCloudStorageHandler) bucketACL() {
	io.WriteString(cloudStorage.w, "\nBucket ACL:\n")
	cloudStorage.dumpBucketACL()
}

// putBucketACLRule adds the "allUsers" bucket ACL rule for this bucket.
func (cloudStorage *GoogleCloudStorageHandler) putBucketACLRule() {
	io.WriteString(cloudStorage.w, "\nPut Bucket ACL Rule:\n")
	err := cloudStorage.bucket.ACL().Set(cloudStorage.ctx, storage.AllUsers, storage.RoleReader)
	if err != nil {
		cloudStorage.errorf("putBucketACLRule: unable to save bucket ACL rule for bucket %q: %v", bucket, err)
		return
	}
	cloudStorage.dumpBucketACL()
}

// deleteBucketACLRule deleted the "allUsers" bucket ACL rule for this bucket.
func (cloudStorage *GoogleCloudStorageHandler) deleteBucketACLRule() {
	io.WriteString(cloudStorage.w, "\nDelete Bucket ACL Rule:\n")
	err := cloudStorage.bucket.ACL().Delete(cloudStorage.ctx, storage.AllUsers)
	if err != nil {
		cloudStorage.errorf("deleteBucketACLRule: unable to delete bucket ACL rule for bucket %q: %v", bucket, err)
		return
	}
	cloudStorage.dumpBucketACL()
}

// dumpACL prints out the ACL of the named file.
func (cloudStorage *GoogleCloudStorageHandler) dumpACL(fileName string) {
	acl, err := cloudStorage.bucket.Object(fileName).ACL().List(cloudStorage.ctx)
	if err != nil {
		cloudStorage.errorf("dumpACL: unable to list file ACL for bucket %q, file %q: %v", bucket, fileName, err)
		return
	}
	for _, v := range acl {
		fmt.Fprintf(cloudStorage.w, "Scope: %q, Permission: %q\n", v.Entity, v.Role)
	}
}

// acl displays the ACL for the named file.
func (cloudStorage *GoogleCloudStorageHandler) acl(fileName string) {
	fmt.Fprintf(cloudStorage.w, "\nACL for file %v:\n", fileName)
	cloudStorage.dumpACL(fileName)
}

// putACLRule adds the "allUsers" ACL rule for the named file.
func (cloudStorage *GoogleCloudStorageHandler) putACLRule(fileName string) {
	fmt.Fprintf(cloudStorage.w, "\nPut ACL rule for file %v:\n", fileName)
	err := cloudStorage.bucket.Object(fileName).ACL().Set(cloudStorage.ctx, storage.AllUsers, storage.RoleReader)
	if err != nil {
		cloudStorage.errorf("putACLRule: unable to save ACL rule for bucket %q, file %q: %v", bucket, fileName, err)
		return
	}
	cloudStorage.dumpACL(fileName)
}

// deleteACLRule deleted the "allUsers" ACL rule for the named file.
func (cloudStorage *GoogleCloudStorageHandler) deleteACLRule(fileName string) {
	fmt.Fprintf(cloudStorage.w, "\nDelete ACL rule for file %v:\n", fileName)
	err := cloudStorage.bucket.Object(fileName).ACL().Delete(cloudStorage.ctx, storage.AllUsers)
	if err != nil {
		cloudStorage.errorf("deleteACLRule: unable to delete ACL rule for bucket %q, file %q: %v", bucket, fileName, err)
		return
	}
	cloudStorage.dumpACL(fileName)
}

// deleteFiles deletes all the temporary files from a bucket created by this demo.
func (cloudStorage *GoogleCloudStorageHandler) deleteFiles() {
	io.WriteString(cloudStorage.w, "\nDeleting files...\n")
	for _, v := range cloudStorage.cleanUp {
		fmt.Fprintf(cloudStorage.w, "Deleting file %v\n", v)
		if err := cloudStorage.bucket.Object(v).Delete(cloudStorage.ctx); err != nil {
			cloudStorage.errorf("deleteFiles: unable to delete bucket %q, file %q: %v", bucket, v, err)
			return
		}
	}
}
*/