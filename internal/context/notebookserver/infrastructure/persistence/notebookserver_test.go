package notebookserver

import (
	"context"
	"os"
	"testing"

	"github.com/onsi/gomega"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	mongooptions "go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/Bio-OS/bioos/internal/context/notebookserver/application/query"
	"github.com/Bio-OS/bioos/internal/context/notebookserver/domain"
	"github.com/Bio-OS/bioos/internal/context/notebookserver/infrastructure/persistence/mongo"
	"github.com/Bio-OS/bioos/internal/context/notebookserver/infrastructure/persistence/sql"
	applog "github.com/Bio-OS/bioos/pkg/log"
	"github.com/Bio-OS/bioos/pkg/notebook"
)

func TestMain(m *testing.M) {
	applog.RegisterLogger(&applog.Options{
		Level: "fatal",
	})
	os.Exit(m.Run())
}

func TestSQLiteMemory(t *testing.T) {
	g := gomega.NewWithT(t)
	ctx := context.TODO()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	g.Expect(err).ToNot(gomega.HaveOccurred())

	repo, err := sql.NewRepository(ctx, db)
	g.Expect(err).ToNot(gomega.HaveOccurred())
	testRepository(ctx, g, repo)

	read, err := sql.NewReadModel(ctx, db)
	g.Expect(err).ToNot(gomega.HaveOccurred())
	testReadModel(ctx, g, repo, read)
}

func TestMongoDB(t *testing.T) {
	uri := os.Getenv("MONGO_URI")
	if len(uri) == 0 {
		t.Logf("set env MONGO_URI to enable TestMongoDB. e.g. mongodb://user:passwd@localhost:27017")
		return
	}

	g := gomega.NewWithT(t)
	ctx := context.TODO()

	client, err := mongodriver.Connect(ctx, mongooptions.Client().ApplyURI(uri))
	g.Expect(err).ToNot(gomega.HaveOccurred())
	t.Logf("connected %s and begin test", uri)
	db := client.Database("go-unit-test")

	repo, err := mongo.NewRepository(ctx, db)
	g.Expect(err).ToNot(gomega.HaveOccurred())
	testRepository(ctx, g, repo)

	read, err := mongo.NewReadModel(ctx, db)
	g.Expect(err).ToNot(gomega.HaveOccurred())
	testReadModel(ctx, g, repo, read)
}

func testRepository(ctx context.Context, g *gomega.WithT, repo domain.Repository) {
	srv := newNotebookServer()

	g.Expect(repo.Save(ctx, srv)).ToNot(gomega.HaveOccurred())

	got, err := repo.Get(ctx, srv.ID)
	g.Expect(err).ToNot(gomega.HaveOccurred())
	g.Expect(got).ToNot(gomega.BeNil())
	g.Expect(got.ID).To(gomega.Equal(srv.ID))
	g.Expect(got.WorkspaceID).To(gomega.Equal(srv.WorkspaceID))
	g.Expect(got.Settings).To(gomega.Equal(srv.Settings))

	srv.Settings.DockerImage = "new content"
	g.Expect(repo.Save(ctx, srv)).ToNot(gomega.HaveOccurred())
	got, err = repo.Get(ctx, srv.ID)
	g.Expect(err).ToNot(gomega.HaveOccurred())
	g.Expect(got).ToNot(gomega.BeNil())
	g.Expect(got.Settings).To(gomega.Equal(srv.Settings))
	g.Expect(got.Volumes).To(gomega.Equal(srv.Volumes))

	g.Expect(repo.Delete(ctx, srv)).ToNot(gomega.HaveOccurred())

	got, err = repo.Get(ctx, "no-exist-id")
	g.Expect(err).ToNot(gomega.HaveOccurred())
	g.Expect(got).To(gomega.BeNil())
}

func testReadModel(ctx context.Context, g *gomega.WithT, repo domain.Repository, read query.ReadModel) {
	// prepare data for reading
	srv := newNotebookServer()
	g.Expect(repo.Save(ctx, srv)).ToNot(gomega.HaveOccurred())
	defer func() {
		g.Expect(repo.Delete(ctx, srv)).ToNot(gomega.HaveOccurred())
	}()

	// test get
	got, err := read.GetSettingsByID(ctx, "workspace-1", "id-1")
	g.Expect(err).ToNot(gomega.HaveOccurred())
	g.Expect(got).ToNot(gomega.BeNil())
	g.Expect(got.ID).To(gomega.Equal(srv.ID))
	g.Expect(got.WorkspaceID).To(gomega.Equal(srv.WorkspaceID))
	g.Expect(got.Image).To(gomega.Equal(srv.Settings.DockerImage))
	g.Expect(got.ResourceSize).To(gomega.Equal(srv.Settings.ResourceSize))

	list, err := read.ListSettingsByWorkspace(ctx, "no-exist-workspace")
	g.Expect(err).ToNot(gomega.HaveOccurred())
	g.Expect(list).To(gomega.BeNil())
	list, err = read.ListSettingsByWorkspace(ctx, "workspace-1")
	g.Expect(err).ToNot(gomega.HaveOccurred())
	g.Expect(list).To(gomega.HaveLen(1))
}

func newNotebookServer() *domain.NotebookServer {
	return &domain.NotebookServer{
		ID:          "id-1",
		WorkspaceID: "workspace-1",
		Settings: domain.Settings{
			DockerImage: "image-1",
			ResourceSize: notebook.ResourceSize{
				CPU:    1,
				Memory: 1024,
				Disk:   1024,
			},
		},
		Volumes: []domain.Volume{
			{
				Name:              "v1",
				Type:              "nfs",
				Source:            "/mnt/nfs",
				MountRelativePath: "test",
			},
		},
	}
}
