package main

import (
	"encoding/xml"
	"fmt"
	liberr "github.com/konveyor/controller/pkg/error"
	"github.com/konveyor/tackle2-addon/command"
	"github.com/konveyor/tackle2-addon/repository"
	"github.com/konveyor/tackle2-hub/api"
	"github.com/mundra-ankur/tackle2-addon-jkube/pom"
	"os"
	"path"
	"strings"
)

type Jkube struct {
	application *api.Application
	repository  repository.Repository
}

func (r *Jkube) Run() (err error) {
	output := r.output()
	cmd := command.Command{Path: "/usr/bin/rm"}
	cmd.Options.Add("-rf", output)
	err = cmd.Run()
	if err != nil {
		return
	}
	err = os.MkdirAll(output, 0777)
	if err != nil {
		err = liberr.Wrap(err, "path", output)
		return
	}
	addon.Activity("[Jkube] created: %s.", output)

	// Fetch the repository.
	err = r.repository.Fetch()
	if err != nil {
		return err
	}

	// Add the jkube plugin to the pom.xml
	groupId, artifactId, err := r.addJkubePlugin()
	if err != nil {
		return err
	}

	// Build the maven project
	err = r.buildMvnProject()
	if err != nil {
		return err
	}

	// Copy the resources to the output directory
	err = r.commitResources(groupId, artifactId)
	if err != nil {
		return err
	}
	return
}

// output returns output directory.
func (r *Jkube) output() string {
	return path.Join(SourceDir, "k8sResources")
}

func (r *Jkube) addJkubePlugin() (groupId string, artifactId string, err error) {
	pomXml := path.Join(SourceDir, "pom.xml") // Path to the pom.xml file
	parsedPom, err := pom.Parse(pomXml)       // Parse the pom.xml file
	if err != nil {
		fmt.Printf("Error parsing pom.xml %s", err)
		return
	}

	jkubePlugin := pom.Plugin{
		GroupID:    "org.eclipse.jkube",
		ArtifactID: "kubernetes-maven-plugin",
		Version:    "1.9.1",
	}

	*parsedPom.Build.Plugins = append(*parsedPom.Build.Plugins, jkubePlugin)

	// Marshal the pom back to xml
	output, err := xml.MarshalIndent(parsedPom, "  ", "    ")
	if err != nil {
		fmt.Printf("Error marshalling pom.xml %s", err)
		return
	}

	// Write the xml to a file
	err = os.Chdir(SourceDir)
	err = os.WriteFile("pom.xml", output, 0644)
	if err != nil {
		fmt.Printf("Error writing pom.xml %s", err)
		return
	}

	return parsedPom.GroupID, parsedPom.ArtifactID, nil
}

func (r *Jkube) buildMvnProject() (err error) {
	// Build the maven project
	cmd := command.Command{
		Path:    "./mvnw",
		Options: []string{"package", "-Dmaven.test.skip", "--quiet"},
		Dir:     SourceDir,
	}

	err = cmd.Run()
	if err != nil {
		fmt.Printf("Error building maven project %s", err)
		return
	}

	// Run mvn k8s:build and mvn k8s:resource
	cmd = command.Command{
		Path:    "./mvnw",
		Options: []string{"k8s:build", "-Djkube.build.strategy=jib", "k8s:resource"},
		Dir:     SourceDir,
	}

	err = cmd.Run()
	if err != nil {
		fmt.Printf("Error running mvn k8s:build and k8s:resource %s", err)
		return
	}

	return
}

func (r *Jkube) commitResources(groupId string, artifactId string) (err error) {
	// Copy the k8s resources to the output directory
	cmd := command.Command{
		Path: "/usr/bin/cp",
		Options: []string{
			"-r",
			path.Join(SourceDir, "target", "classes", "META-INF", "jkube/"), ".",
			r.output(),
		},
		Dir: SourceDir,
	}

	err = cmd.Run()
	if err != nil {
		fmt.Printf("Error copying resources %s", err)
		return
	}

	// Copy the Dockerfile to the output directory
	group := strings.ToLower(strings.Split(groupId, ".")[1])
	artifactId = strings.ToLower(artifactId)
	cmd = command.Command{
		Path: "/usr/bin/cp",
		Options: []string{
			path.Join(SourceDir, "target", "docker", group, artifactId, "latest", "build", "Dockerfile"),
			path.Join(r.output(), "Dockerfile")},
		Dir: SourceDir,
	}

	err = cmd.Run()
	if err != nil {
		fmt.Printf("Error copying Dockerfile %s", err)
	}

	cmd = command.Command{
		Path:    "/usr/bin/git",
		Options: []string{"config", "--global", "user.email", "tackle@konveyor.org"},
		Dir:     SourceDir,
	}

	err = cmd.Run()
	if err != nil {
		fmt.Printf("Error setting git config %s", err)
	}

	cmd = command.Command{
		Path:    "/usr/bin/git",
		Options: []string{"config", "--global", "user.name", "tackle"},
		Dir:     SourceDir,
	}

	err = cmd.Run()
	if err != nil {
		fmt.Printf("Error setting git config %s", err)
	}

	cmd = command.Command{
		Path:    "/usr/bin/git",
		Options: []string{"add", path.Base(r.output())},
		Dir:     SourceDir,
	}

	err = cmd.Run()
	if err != nil {
		fmt.Printf("Error adding k8s resources to git %s", err)
		return
	}

	cmd = command.Command{
		Path:    "/usr/bin/git",
		Options: []string{"commit", "-m", "Add k8s resources"},
		Dir:     SourceDir,
	}

	err = cmd.Run()
	if err != nil {
		fmt.Printf("Error committing k8s resources %s", err)
		return
	}

	cmd = command.Command{
		Path:    "/usr/bin/git",
		Options: []string{"push"},
		Dir:     SourceDir,
	}

	err = cmd.Run()
	if err != nil {
		fmt.Printf("Error pushing k8s resources %s", err)
	}
	return
}
