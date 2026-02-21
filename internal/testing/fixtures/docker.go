package fixtures

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/containerd/errdefs"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/google/uuid"
	"github.com/testcontainers/testcontainers-go"
	tcexec "github.com/testcontainers/testcontainers-go/exec"
	"github.com/testcontainers/testcontainers-go/modules/mariadb"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	networkName        = "forgejo"
	mariaDBImage       = "mariadb:lts"
	forgejoImage       = "codeberg.org/forgejo/forgejo:11"
	forgejoUserName    = "tfadmin"
	forgejoUserEmail   = "tfadmin@localhost"
	forgejoTokenScopes = "write:organization,write:repository,write:user,write:admin"
)

var (
	dbContainer      *DBContainer
	forgejoContainer *ForgejoContainer
)

type TestContainers struct {
	DBContainer      *DBContainer
	ForgejoContainer *ForgejoContainer
}

func GetTestContainers(ctx context.Context) (*TestContainers, error) {
	db, err := getDBContainer(ctx)
	if err != nil {
		return nil, err
	}

	forgejo, err := getForgejoContainer(ctx)
	if err != nil {
		return nil, err
	}

	return &TestContainers{DBContainer: db, ForgejoContainer: forgejo}, nil
}

type DBContainer struct {
	c *mariadb.MariaDBContainer
}

func getDBContainer(ctx context.Context) (*DBContainer, error) {
	if dbContainer != nil {
		return dbContainer, nil
	}

	err := ensureDockerTestNetwork(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting docker network: %w", err)
	}
	containerName := "forgejo_db"

	containerRequest := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Name:     containerName,
			Networks: []string{networkName},
			ConfigModifier: func(config *container.Config) {
				config.Hostname = containerName
			},
		},
		Reuse: true,
	}

	c, err := mariadb.Run(
		ctx,
		mariaDBImage,
		mariadb.WithDatabase("forgejo"),
		mariadb.WithUsername("forgejo"),
		mariadb.WithPassword("password"),
		testcontainers.CustomizeRequest(containerRequest),
		testcontainers.WithCmd(
			"--transaction-isolation=READ-COMMITTED",
			"--binlog-format=ROW",
		),
	)
	if err != nil {
		return nil, fmt.Errorf("error starting db container: %w", err)
	}

	dbContainer = &DBContainer{c: c}

	return dbContainer, nil
}

type ForgejoContainer struct {
	c     *testcontainers.DockerContainer
	token string
}

func (f *ForgejoContainer) GetHost(ctx context.Context) (string, error) {
	host, err := f.c.Host(ctx)
	if err != nil {
		return "", err
	}

	port, err := f.c.MappedPort(ctx, "3000/tcp")
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("http://%s:%d", host, port.Int()), nil
}

func (f *ForgejoContainer) GetAPIToken(ctx context.Context) (string, error) {
	if f.token != "" {
		return f.token, nil
	}

	err := f.ensureAdminUser(ctx)
	if err != nil {
		return "", fmt.Errorf("error creating admin user: %w", err)
	}

	token, err := f.createAPIToken(ctx)
	if err != nil {
		return "", fmt.Errorf("error creating API token: %w", err)
	}

	return strings.TrimSpace(token), nil
}

func (f *ForgejoContainer) ensureAdminUser(ctx context.Context) error {
	users, err := f.executeAdminCommand(ctx, []string{"user", "list", "--admin"})
	if err != nil {
		return err
	}

	if strings.Contains(users, forgejoUserEmail) {
		return nil
	}

	_, err = f.executeAdminCommand(
		ctx,
		[]string{
			"user",
			"create",
			"--username",
			forgejoUserName,
			"--email",
			forgejoUserEmail,
			"--random-password",
			"--admin",
		},
	)

	if err != nil {
		return err
	}

	return nil
}

func (f *ForgejoContainer) createAPIToken(ctx context.Context) (string, error) {
	token, err := f.executeAdminCommand(
		ctx,
		[]string{
			"user",
			"generate-access-token",
			"--username",
			forgejoUserName,
			"--scopes",
			forgejoTokenScopes,
			"--token-name",
			uuid.NewString(),
			"--raw",
		},
	)

	if err != nil {
		return "", err
	}

	return token, nil
}

func (f *ForgejoContainer) executeAdminCommand(ctx context.Context, cmd []string) (string, error) {
	cmd = append([]string{"/usr/local/bin/forgejo", "admin"}, cmd...)

	exitCode, output, err := f.c.Exec(ctx, cmd, tcexec.WithUser("git"))

	if err != nil {
		return "", err
	}

	var stdout, stderr bytes.Buffer
	_, err = stdcopy.StdCopy(&stdout, &stderr, output)
	if err != nil {
		return "", fmt.Errorf("error reading output: %w", err)
	}

	if exitCode != 0 {
		return "", errors.New(stderr.String())
	}

	return stdout.String(), nil
}

func getForgejoContainer(ctx context.Context) (*ForgejoContainer, error) {
	if forgejoContainer != nil {
		return forgejoContainer, nil
	}

	err := ensureDockerTestNetwork(ctx)
	if err != nil {
		return nil, err
	}

	containerName := "forgejo"

	containerRequest := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Name:     containerName,
			Networks: []string{networkName},
			ConfigModifier: func(config *container.Config) {
				config.Hostname = containerName
			},
		},
		Reuse: true,
	}

	_, f, _, _ := runtime.Caller(0)
	appIniFile := filepath.Join(filepath.Dir(f), "app.ini")

	c, err := testcontainers.Run(
		ctx,
		forgejoImage,
		testcontainers.CustomizeRequest(containerRequest),
		testcontainers.WithFiles(testcontainers.ContainerFile{
			HostFilePath:      appIniFile,
			ContainerFilePath: "/data/gitea/conf/app.ini",
			FileMode:          0o644,
		}),
		testcontainers.WithWaitStrategy(
			wait.ForLog("Starting new Web server"),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("error starting forgejo container: %w", err)
	}

	forgejoContainer = &ForgejoContainer{c: c}

	return forgejoContainer, nil
}

func ensureDockerTestNetwork(ctx context.Context) error {
	docker, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return fmt.Errorf("error creating docker client: %w", err)
	}
	defer docker.Close()

	_, err = docker.NetworkInspect(ctx, networkName, network.InspectOptions{})

	// Network already exists, nothing to do
	if err == nil {
		return nil
	}

	// Anything else than not found is unexpected, return it
	if !errdefs.IsNotFound(err) {
		return fmt.Errorf("error checking for docker network: %w", err)
	}

	_, err = docker.NetworkCreate(ctx, networkName, network.CreateOptions{
		Driver: "bridge",
	})
	if err != nil {
		return fmt.Errorf("error creating docker network: %w", err)
	}

	return nil
}
