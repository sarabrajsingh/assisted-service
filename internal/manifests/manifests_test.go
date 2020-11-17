package manifests_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"errors"

	"github.com/go-openapi/strfmt"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/openshift/assisted-service/internal/common"
	"github.com/openshift/assisted-service/internal/manifests"
	"github.com/openshift/assisted-service/models"
	"github.com/openshift/assisted-service/pkg/filemiddleware"
	"github.com/openshift/assisted-service/pkg/s3wrapper"
	operations "github.com/openshift/assisted-service/restapi/operations/manifests"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
)

func TestValidator(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "manifests_test")
}

func getTestLog() logrus.FieldLogger {
	l := logrus.New()
	l.SetOutput(ioutil.Discard)
	return l
}

var _ = Describe("ClusterManifestTests", func() {
	var (
		manifestsAPI  *manifests.Manifests
		db            *gorm.DB
		ctx           = context.Background()
		ctrl          *gomock.Controller
		mockS3Client  *s3wrapper.MockAPI
		dbName        = "cluster_manifest"
		content       = "aGVsbG8gd29ybGQhCg=="
		fileName      = "99-test.yaml"
		validFolder   = "openshift"
		defaultFolder = "manifests"
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		db = common.PrepareTestDB(dbName)
		mockS3Client = s3wrapper.NewMockAPI(ctrl)
		manifestsAPI = manifests.NewManifestsAPI(db, getTestLog(), mockS3Client)
	})

	AfterEach(func() {
		ctrl.Finish()

	})

	registerCluster := func() *common.Cluster {
		clusterID := strfmt.UUID(uuid.New().String())
		cluster := common.Cluster{
			Cluster: models.Cluster{
				ID: &clusterID,
			},
		}
		Expect(db.Create(&cluster).Error).ShouldNot(HaveOccurred())
		return &cluster
	}

	addManifestToCluster := func(clusterID *strfmt.UUID, content, fileName, folderName string) {
		response := manifestsAPI.CreateClusterManifest(ctx, operations.CreateClusterManifestParams{
			ClusterID: *clusterID,
			CreateManifestParams: &models.CreateManifestParams{
				Content:  &content,
				FileName: &fileName,
			},
		})
		Expect(response).Should(BeAssignableToTypeOf(operations.NewCreateClusterManifestCreated()))
	}

	getObjectName := func(clusterID *strfmt.UUID, folderName, fileName string) string {
		return fmt.Sprintf("%s/manifests/%s/%s", *clusterID, folderName, fileName)
	}

	mockObjectExists := func(exists bool) {
		mockS3Client.EXPECT().DoesObjectExist(ctx, gomock.Any()).Return(exists, nil).Times(1)
	}

	mockUpload := func(times int) {
		mockS3Client.EXPECT().Upload(ctx, gomock.Any(), gomock.Any()).Return(nil).Times(times)
	}

	mockListByPrefix := func(clusterID *strfmt.UUID, files []string) {
		prefix := fmt.Sprintf("%s/manifests", *clusterID)
		mockS3Client.EXPECT().ListObjectsByPrefix(ctx, prefix).Return(files, nil).Times(1)
	}

	Context("CreateClusterManifest", func() {
		It("creates manifest successfully with default folder", func() {
			clusterID := registerCluster().ID
			mockUpload(1)
			response := manifestsAPI.CreateClusterManifest(ctx, operations.CreateClusterManifestParams{
				ClusterID: *clusterID,
				CreateManifestParams: &models.CreateManifestParams{
					Content:  &content,
					FileName: &fileName,
				},
			})
			Expect(response).Should(BeAssignableToTypeOf(operations.NewCreateClusterManifestCreated()))
			responsePayload := response.(*operations.CreateClusterManifestCreated)
			Expect(responsePayload.Payload).ShouldNot(BeNil())
			Expect(responsePayload.Payload.FileName).To(Equal(fileName))
			Expect(responsePayload.Payload.Folder).To(Equal(defaultFolder))
		})

		It("creates manifest successfully with 'openshift' folder", func() {
			clusterID := registerCluster().ID
			mockUpload(1)
			response := manifestsAPI.CreateClusterManifest(ctx, operations.CreateClusterManifestParams{
				ClusterID: *clusterID,
				CreateManifestParams: &models.CreateManifestParams{
					Content:  &content,
					FileName: &fileName,
					Folder:   &validFolder,
				},
			})
			Expect(response).Should(BeAssignableToTypeOf(operations.NewCreateClusterManifestCreated()))
			responsePayload := response.(*operations.CreateClusterManifestCreated)
			Expect(responsePayload.Payload).ShouldNot(BeNil())
			Expect(responsePayload.Payload.FileName).To(Equal(fileName))
			Expect(responsePayload.Payload.Folder).To(Equal(validFolder))
		})

		It("override an existing manifest", func() {
			clusterID := registerCluster().ID
			mockUpload(2)
			response := manifestsAPI.CreateClusterManifest(ctx, operations.CreateClusterManifestParams{
				ClusterID: *clusterID,
				CreateManifestParams: &models.CreateManifestParams{
					Content:  &content,
					FileName: &fileName,
				},
			})
			Expect(response).Should(BeAssignableToTypeOf(operations.NewCreateClusterManifestCreated()))
			responsePayload := response.(*operations.CreateClusterManifestCreated)
			Expect(responsePayload.Payload).ShouldNot(BeNil())
			Expect(responsePayload.Payload.FileName).To(Equal(fileName))
			Expect(responsePayload.Payload.Folder).To(Equal(defaultFolder))

			response = manifestsAPI.CreateClusterManifest(ctx, operations.CreateClusterManifestParams{
				ClusterID: *clusterID,
				CreateManifestParams: &models.CreateManifestParams{
					Content:  &content,
					FileName: &fileName,
				},
			})
			Expect(response).Should(BeAssignableToTypeOf(operations.NewCreateClusterManifestCreated()))
			responsePayload = response.(*operations.CreateClusterManifestCreated)
			Expect(responsePayload.Payload).ShouldNot(BeNil())
			Expect(responsePayload.Payload.FileName).To(Equal(fileName))
			Expect(responsePayload.Payload.Folder).To(Equal(defaultFolder))
		})

		It("cluster doesn't exist", func() {
			response := manifestsAPI.CreateClusterManifest(ctx, operations.CreateClusterManifestParams{
				ClusterID: strfmt.UUID(uuid.New().String()),
				CreateManifestParams: &models.CreateManifestParams{
					Content:  &content,
					FileName: &fileName,
				},
			})

			Expect(response).Should(BeAssignableToTypeOf(common.NewApiError(http.StatusNotFound, errors.New(""))))
			err := response.(*common.ApiErrorResponse)
			Expect(err.StatusCode()).To(Equal(int32(http.StatusNotFound)))
		})

		It("fails due to non-base64 file content", func() {
			clusterID := registerCluster().ID
			invalidContent := "not base64 content"
			response := manifestsAPI.CreateClusterManifest(ctx, operations.CreateClusterManifestParams{
				ClusterID: *clusterID,
				CreateManifestParams: &models.CreateManifestParams{
					Content:  &invalidContent,
					FileName: &fileName,
				},
			})
			Expect(response).Should(BeAssignableToTypeOf(common.NewApiError(http.StatusBadRequest, errors.New(""))))
			err := response.(*common.ApiErrorResponse)
			Expect(err.StatusCode()).To(Equal(int32(http.StatusBadRequest)))
			Expect(err.Error()).To(ContainSubstring("failed to base64-decode cluster manifest content"))
		})
	})

	Context("ListClusterManifests", func() {
		It("lists manifest from different folders", func() {
			clusterID := registerCluster().ID
			mockUpload(3)
			files := []string{
				getObjectName(clusterID, validFolder, "file-1.yaml"),
				getObjectName(clusterID, validFolder, "file-2.yaml"),
				getObjectName(clusterID, defaultFolder, "file-3.yaml"),
			}
			mockListByPrefix(clusterID, files)
			addManifestToCluster(clusterID, content, "file-1.yaml", validFolder)
			addManifestToCluster(clusterID, content, "file-2.yaml", validFolder)
			addManifestToCluster(clusterID, content, "file-1.yaml", defaultFolder)
			response := manifestsAPI.ListClusterManifests(ctx, operations.ListClusterManifestsParams{
				ClusterID: *clusterID,
			})
			Expect(response).Should(BeAssignableToTypeOf(operations.NewListClusterManifestsOK()))
			responsePayload := response.(*operations.ListClusterManifestsOK)
			Expect(responsePayload.Payload).ShouldNot(BeNil())
			Expect(len(responsePayload.Payload)).To(Equal(3))
		})

		It("list manifests for new cluster", func() {
			clusterID := registerCluster().ID
			mockListByPrefix(clusterID, []string{})
			response := manifestsAPI.ListClusterManifests(ctx, operations.ListClusterManifestsParams{
				ClusterID: *clusterID,
			})
			Expect(response).Should(BeAssignableToTypeOf(operations.NewListClusterManifestsOK()))
			responsePayload := response.(*operations.ListClusterManifestsOK)
			Expect(responsePayload.Payload).ShouldNot(BeNil())
			Expect(len(responsePayload.Payload)).To(Equal(0))
		})

		It("cluster doesn't exist", func() {
			response := manifestsAPI.ListClusterManifests(ctx, operations.ListClusterManifestsParams{
				ClusterID: strfmt.UUID(uuid.New().String()),
			})

			Expect(response).Should(BeAssignableToTypeOf(common.NewApiError(http.StatusNotFound, errors.New(""))))
			err := response.(*common.ApiErrorResponse)
			Expect(err.StatusCode()).To(Equal(int32(http.StatusNotFound)))
		})
	})

	Context("DeleteClusterManifest", func() {
		It("deletes manifest from default folder", func() {
			clusterID := registerCluster().ID
			mockUpload(1)
			mockObjectExists(true)
			mockS3Client.EXPECT().DeleteObject(ctx, getObjectName(clusterID, defaultFolder, "file-1.yaml")).Return(nil)
			addManifestToCluster(clusterID, content, "file-1.yaml", defaultFolder)

			response := manifestsAPI.DeleteClusterManifest(ctx, operations.DeleteClusterManifestParams{
				ClusterID: *clusterID,
				FileName:  "file-1.yaml",
			})
			Expect(response).Should(BeAssignableToTypeOf(operations.NewDeleteClusterManifestOK()))
		})

		It("deletes manifest from a different folder", func() {
			clusterID := registerCluster().ID
			mockUpload(1)
			mockObjectExists(true)
			mockS3Client.EXPECT().DeleteObject(ctx, getObjectName(clusterID, validFolder, "file-1.yaml")).Return(nil)
			addManifestToCluster(clusterID, content, "file-1.yaml", validFolder)

			response := manifestsAPI.DeleteClusterManifest(ctx, operations.DeleteClusterManifestParams{
				ClusterID: *clusterID,
				FileName:  "file-1.yaml",
				Folder:    &validFolder,
			})
			Expect(response).Should(BeAssignableToTypeOf(operations.NewDeleteClusterManifestOK()))
		})

		It("deletes missing manifest", func() {
			clusterID := registerCluster().ID
			mockObjectExists(false)

			response := manifestsAPI.DeleteClusterManifest(ctx, operations.DeleteClusterManifestParams{
				ClusterID: *clusterID,
				FileName:  "file-1.yaml",
			})
			Expect(response).Should(BeAssignableToTypeOf(operations.NewDeleteClusterManifestOK()))
		})

		It("cluster doesn't exist", func() {
			response := manifestsAPI.DeleteClusterManifest(ctx, operations.DeleteClusterManifestParams{
				ClusterID: strfmt.UUID(uuid.New().String()),
				FileName:  "file-1.yaml",
			})

			Expect(response).Should(BeAssignableToTypeOf(common.NewApiError(http.StatusNotFound, errors.New(""))))
			err := response.(*common.ApiErrorResponse)
			Expect(err.StatusCode()).To(Equal(int32(http.StatusNotFound)))
		})
	})

	Context("DownloadClusterManifest", func() {
		It("downloads manifest from different folder", func() {
			clusterID := registerCluster().ID
			mockUpload(1)
			mockObjectExists(true)
			mockS3Client.EXPECT().Download(ctx, gomock.Any()).Return(VoidReadCloser{}, int64(0), nil)
			addManifestToCluster(clusterID, content, "file-1.yaml", defaultFolder)

			response := manifestsAPI.DownloadClusterManifest(ctx, operations.DownloadClusterManifestParams{
				ClusterID: *clusterID,
				FileName:  "file-1.yaml",
			})
			Expect(response).Should(BeAssignableToTypeOf(filemiddleware.NewResponder(nil, "", int64(0))))
		})

		It("downloads missing manifest", func() {
			clusterID := registerCluster().ID
			mockObjectExists(false)

			response := manifestsAPI.DownloadClusterManifest(ctx, operations.DownloadClusterManifestParams{
				ClusterID: *clusterID,
				FileName:  "file-1.yaml",
			})
			Expect(response).Should(BeAssignableToTypeOf(common.NewApiError(http.StatusNotFound, errors.New(""))))
		})

		It("cluster doesn't exist", func() {
			response := manifestsAPI.DownloadClusterManifest(ctx, operations.DownloadClusterManifestParams{
				ClusterID: strfmt.UUID(uuid.New().String()),
			})

			Expect(response).Should(BeAssignableToTypeOf(common.NewApiError(http.StatusNotFound, errors.New(""))))
			err := response.(*common.ApiErrorResponse)
			Expect(err.StatusCode()).To(Equal(int32(http.StatusNotFound)))
		})
	})
})

type VoidReadCloser struct {
}

func (VoidReadCloser) Read(p []byte) (int, error) {
	return 0, nil
}

func (VoidReadCloser) Close() error {
	return nil
}
