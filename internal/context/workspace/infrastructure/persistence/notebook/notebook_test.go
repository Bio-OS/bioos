package notebook

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/onsi/gomega"

	query "github.com/Bio-OS/bioos/internal/context/workspace/application/query/notebook"
	"github.com/Bio-OS/bioos/internal/context/workspace/domain/notebook"
	"github.com/Bio-OS/bioos/internal/context/workspace/infrastructure/persistence/notebook/filesystem"
	applog "github.com/Bio-OS/bioos/pkg/log"
)

func TestMain(m *testing.M) {
	applog.RegisterLogger(&applog.Options{
		Level: "fatal",
	})
	os.Exit(m.Run())
}

func TestFileSystem(t *testing.T) {
	g := gomega.NewWithT(t)
	ctx := context.TODO()

	tempdir := os.TempDir()
	repo, err := filesystem.NewRepository(tempdir)
	g.Expect(err).ToNot(gomega.HaveOccurred())
	testRepository(ctx, g, repo)

	read, err := filesystem.NewReadModel(tempdir)
	g.Expect(err).ToNot(gomega.HaveOccurred())
	testReadModel(ctx, g, repo, read)
}

func testRepository(ctx context.Context, g *gomega.WithT, repo notebook.Repository) {
	nb := &notebook.Notebook{
		Name:      "nb-1",
		Namespace: "workspace-1",
		Content:   []byte("abcd"),
	}

	g.Expect(repo.Save(ctx, nb)).ToNot(gomega.HaveOccurred())

	got, err := repo.Get(ctx, nb.Path())
	g.Expect(err).ToNot(gomega.HaveOccurred())
	g.Expect(got.Name).To(gomega.Equal(nb.Name))
	g.Expect(got.Namespace).To(gomega.Equal(nb.Namespace))
	g.Expect(got.Content).To(gomega.Equal(nb.Content))

	nb.Content = []byte("new content")
	g.Expect(repo.Save(ctx, nb)).ToNot(gomega.HaveOccurred())
	got, err = repo.Get(ctx, nb.Path())
	g.Expect(err).ToNot(gomega.HaveOccurred())
	g.Expect(got.Content).To(gomega.Equal(nb.Content))

	g.Expect(repo.Delete(ctx, nb)).ToNot(gomega.HaveOccurred())
}

func testReadModel(ctx context.Context, g *gomega.WithT, repo notebook.Repository, read query.ReadModel) {
	// prepare data for reading
	content := []byte("abcd")
	nb := &notebook.Notebook{
		Namespace: "workspace-1",
		Content:   content,
	}
	count := 6
	for i := 0; i < count; i++ {
		nb.Name = fmt.Sprintf("name-%d", i)
		g.Expect(repo.Save(ctx, nb)).ToNot(gomega.HaveOccurred())
	}
	defer func() {
		for i := 0; i < count; i++ {
			nb.Name = fmt.Sprintf("name-%d", i)
			g.Expect(repo.Delete(ctx, nb)).ToNot(gomega.HaveOccurred())
		}
	}()

	// test get
	got, err := read.Get(ctx, "workspace-1", "name-0")
	g.Expect(err).ToNot(gomega.HaveOccurred())
	g.Expect(got).ToNot(gomega.BeNil())
	g.Expect(got.Content).To(gomega.Equal(content))

	// test list
	cases := []struct {
		workspaceID string
		length      int
	}{
		{"workspace-1", count},
	}
	for _, c := range cases {
		list, err := read.ListByWorkspace(ctx, c.workspaceID)
		g.Expect(err).ToNot(gomega.HaveOccurred())
		g.Expect(list).To(gomega.HaveLen(c.length))
		for _, n := range list {
			g.Expect(n.WorkspaceID).To(gomega.Equal(c.workspaceID))
			g.Expect(n.Size).To(gomega.BeNumerically("==", len(content)))
		}
	}
}
