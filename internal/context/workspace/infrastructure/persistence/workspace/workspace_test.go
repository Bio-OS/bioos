package main

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/onsi/gomega"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	mongooptions "go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	query "github.com/Bio-OS/bioos/internal/context/workspace/application/query/workspace"
	"github.com/Bio-OS/bioos/internal/context/workspace/domain/workspace"
	"github.com/Bio-OS/bioos/internal/context/workspace/infrastructure/persistence/workspace/mongo"
	workspacesql "github.com/Bio-OS/bioos/internal/context/workspace/infrastructure/persistence/workspace/sql"
	applog "github.com/Bio-OS/bioos/pkg/log"
	"github.com/Bio-OS/bioos/pkg/utils"
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

	repo, err := workspacesql.NewWorkspaceRepository(ctx, db)
	g.Expect(err).ToNot(gomega.HaveOccurred())
	testRepository(ctx, g, repo)

	read, err := workspacesql.NewWorkspaceReadModel(ctx, db)
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

	repo, err := mongo.NewWorkspaceRepository(ctx, db)
	g.Expect(err).ToNot(gomega.HaveOccurred())
	testRepository(ctx, g, repo)

	read, err := mongo.NewWorkspaceReadModel(ctx, db)
	g.Expect(err).ToNot(gomega.HaveOccurred())
	testReadModel(ctx, g, repo, read)
}

func testRepository(ctx context.Context, g *gomega.WithT, repo workspace.Repository) {
	ws := &workspace.Workspace{
		ID:          "id",
		Name:        "name",
		Description: "desc",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Storage: workspace.Storage{
			NFS: &workspace.NFSStorage{
				MountPath: "/mount/path",
			},
		},
	}

	g.Expect(repo.Save(ctx, ws)).ToNot(gomega.HaveOccurred())

	w, err := repo.Get(ctx, ws.ID)
	g.Expect(err).ToNot(gomega.HaveOccurred())
	g.Expect(w.Name).To(gomega.Equal(ws.Name))
	g.Expect(w.Description).To(gomega.Equal(ws.Description))
	g.Expect(w.Storage.NFS).ToNot(gomega.BeNil())
	g.Expect(w.Storage.NFS.MountPath).To(gomega.Equal("/mount/path"))

	ws.Description = "new desc"
	g.Expect(repo.Save(ctx, ws)).ToNot(gomega.HaveOccurred())
	w, err = repo.Get(ctx, ws.ID)
	g.Expect(err).ToNot(gomega.HaveOccurred())
	g.Expect(w.Description).To(gomega.Equal(ws.Description))

	g.Expect(repo.Delete(ctx, ws)).ToNot(gomega.HaveOccurred())
}

func testReadModel(ctx context.Context, g *gomega.WithT, repo workspace.Repository, read query.WorkspaceReadModel) {
	// prepare data for reading
	ws := &workspace.Workspace{
		Name:        "name",
		Description: "desc",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Storage: workspace.Storage{
			NFS: &workspace.NFSStorage{
				MountPath: "/mount/path",
			},
		},
	}
	for i := 0; i < 20; i++ {
		ws.ID = fmt.Sprintf("id-%d", i)
		ws.Name = fmt.Sprintf("name-%d", i)
		g.Expect(repo.Save(ctx, ws)).ToNot(gomega.HaveOccurred())
	}

	cases := []struct {
		filter *query.ListWorkspacesFilter
		page   *utils.Pagination
		length int
	}{
		{nil, utils.NewPagination(20, 1), 20},
		{&query.ListWorkspacesFilter{SearchWord: "name-1"}, utils.NewPagination(20, 1), 11},
		{&query.ListWorkspacesFilter{SearchWord: "_"}, utils.NewPagination(20, 1), 0},
		{&query.ListWorkspacesFilter{IDs: []string{"id-1"}}, utils.NewPagination(20, 1), 1},
	}
	for _, c := range cases {
		list, err := read.ListWorkspaces(ctx, *c.page, c.filter)
		g.Expect(err).ToNot(gomega.HaveOccurred())
		g.Expect(list).To(gomega.HaveLen(c.length))
		if len(list) > 0 {
			g.Expect(list[0].Storage.NFS).ToNot(gomega.BeNil())
		}

		count, err := read.CountWorkspaces(ctx, c.filter)
		g.Expect(list).To(gomega.HaveLen(c.length))
		g.Expect(count).To(gomega.Equal(c.length))
	}
}
