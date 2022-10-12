package pom

import (
	"encoding/xml"
	"io"
	"os"
)

func Parse(path string) (*Pom, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	b, _ := io.ReadAll(file)
	var project Pom

	err = xml.Unmarshal(b, &project)
	if err != nil {
		return nil, err
	}
	return &project, nil
}

type Pom struct {
	XMLName                xml.Name                `xml:"project,omitempty"`
	ModelVersion           string                  `xml:"modelVersion,omitempty"`
	Parent                 *Parent                 `xml:"parent"`
	GroupID                string                  `xml:"groupId,omitempty"`
	ArtifactID             string                  `xml:"artifactId,omitempty"`
	Version                string                  `xml:"version,omitempty"`
	Packaging              string                  `xml:"packaging,omitempty"`
	Name                   string                  `xml:"name,omitempty"`
	Description            string                  `xml:"description,omitempty"`
	URL                    string                  `xml:"url,omitempty"`
	InceptionYear          string                  `xml:"inceptionYear,omitempty"`
	Organization           *Organization           `xml:"organization,omitempty"`
	Licenses               *[]License              `xml:"licenses>license,omitempty"`
	Developers             *[]Developer            `xml:"developers>developer,omitempty"`
	Contributors           *[]Contributor          `xml:"contributors>contributor,omitempty"`
	MailingLists           *[]MailingList          `xml:"mailingLists>mailingList,omitempty"`
	Prerequisites          *Prerequisites          `xml:"prerequisites,omitempty"`
	Modules                []string                `xml:"modules>module"`
	SCM                    *Scm                    `xml:"scm"`
	IssueManagement        *IssueManagement        `xml:"issueManagement"`
	CIManagement           *CIManagement           `xml:"ciManagement"`
	DistributionManagement *DistributionManagement `xml:"distributionManagement"`
	DependencyManagement   *DependencyManagement   `xml:"dependencyManagement"`
	Dependencies           *[]Dependency           `xml:"dependencies>dependency"`
	Repositories           *[]Repository           `xml:"repositories>repository"`
	PluginRepositories     *[]PluginRepository     `xml:"pluginRepositories>pluginRepository"`
	Build                  *Build                  `xml:"build"`
	Reporting              *Reporting              `xml:"reporting"`
	Profiles               *[]Profile              `xml:"profiles>profile"`
	Properties             *Properties             `xml:"properties"`
}

type Properties struct {
	Entries map[string]string
}

func (p *Properties) UnmarshalXML(d *xml.Decoder, start xml.StartElement) (err error) {
	type entry struct {
		XMLName xml.Name
		Key     string `xml:"name,attr"`
		Value   string `xml:",chardata"`
	}
	e := entry{}
	p.Entries = map[string]string{}
	for err = d.Decode(&e); err == nil; err = d.Decode(&e) {
		e.Key = e.XMLName.Local
		p.Entries[e.Key] = e.Value
	}
	if err != nil && err != io.EOF {
		return err
	}

	return nil
}

// MarshalXML marshals Properties into XML.
func (p *Properties) MarshalXML(e *xml.Encoder, start xml.StartElement) error {

	tokens := []xml.Token{start}

	for key, value := range p.Entries {
		t := xml.StartElement{Name: xml.Name{Local: key}}
		tokens = append(tokens, t, xml.CharData(value), xml.EndElement{Name: t.Name})
	}

	tokens = append(tokens, xml.EndElement{Name: start.Name})

	for _, t := range tokens {
		err := e.EncodeToken(t)
		if err != nil {
			return err
		}
	}

	// flush to ensure tokens are written
	return e.Flush()
}

type Parent struct {
	GroupID      string `xml:"groupId,omitempty"`
	ArtifactID   string `xml:"artifactId,omitempty"`
	Version      string `xml:"version,omitempty"`
	RelativePath string `xml:"relativePath,omitempty"`
}

type Organization struct {
	Name string `xml:"name,omitempty"`
	URL  string `xml:"url,omitempty"`
}

type License struct {
	Name         string `xml:"name,omitempty"`
	URL          string `xml:"url,omitempty"`
	Distribution string `xml:"distribution,omitempty"`
	Comments     string `xml:"comments,omitempty"`
}

type Developer struct {
	ID              string      `xml:"id,omitempty"`
	Name            string      `xml:"name,omitempty"`
	Email           string      `xml:"email,omitempty"`
	URL             string      `xml:"url,omitempty"`
	Organization    string      `xml:"organization,omitempty"`
	OrganizationURL string      `xml:"organizationUrl,omitempty"`
	Roles           []string    `xml:"roles>role"`
	Timezone        string      `xml:"timezone,omitempty"`
	Properties      *Properties `xml:"properties"`
}

type Contributor struct {
	Name            string      `xml:"name,omitempty"`
	Email           string      `xml:"email,omitempty"`
	URL             string      `xml:"url,omitempty"`
	Organization    string      `xml:"organization,omitempty"`
	OrganizationURL string      `xml:"organizationUrl,omitempty"`
	Roles           []string    `xml:"roles>role"`
	Timezone        string      `xml:"timezone,omitempty"`
	Properties      *Properties `xml:"properties"`
}

type MailingList struct {
	Name          string   `xml:"name,omitempty"`
	Subscribe     string   `xml:"subscribe,omitempty"`
	Unsubscribe   string   `xml:"unsubscribe,omitempty"`
	Post          string   `xml:"post,omitempty"`
	Archive       string   `xml:"archive,omitempty"`
	OtherArchives []string `xml:"otherArchives>otherArchive"`
}

type Prerequisites struct {
	Maven string `xml:"maven,omitempty"`
}

type Scm struct {
	Connection          string `xml:"connection,omitempty"`
	DeveloperConnection string `xml:"developerConnection,omitempty"`
	Tag                 string `xml:"tag,omitempty"`
	URL                 string `xml:"url,omitempty"`
}

type IssueManagement struct {
	System string `xml:"system,omitempty"`
	URL    string `xml:"url,omitempty"`
}

type CIManagement struct {
	System    string      `xml:"system,omitempty"`
	URL       string      `xml:"url,omitempty"`
	Notifiers *[]Notifier `xml:"notifiers>notifier"`
}

type Notifier struct {
	Type          string      `xml:"type,omitempty"`
	SendOnError   bool        `xml:"sendOnError,omitempty"`
	SendOnFailure bool        `xml:"sendOnFailure,omitempty"`
	SendOnSuccess bool        `xml:"sendOnSuccess,omitempty"`
	SendOnWarning bool        `xml:"sendOnWarning,omitempty"`
	Address       string      `xml:"address,omitempty"`
	Configuration *Properties `xml:"configuration"`
}

type DistributionManagement struct {
	Repository         *Repository `xml:"repository"`
	SnapshotRepository *Repository `xml:"snapshotRepository"`
	Site               *Site       `xml:"site"`
	DownloadURL        string      `xml:"downloadUrl,omitempty"`
	Relocation         *Relocation `xml:"relocation"`
	Status             string      `xml:"status,omitempty"`
}

type Site struct {
	ID   string `xml:"id,omitempty"`
	Name string `xml:"name,omitempty"`
	URL  string `xml:"url,omitempty"`
}

type Relocation struct {
	GroupID    string `xml:"groupId,omitempty"`
	ArtifactID string `xml:"artifactId,omitempty"`
	Version    string `xml:"version,omitempty"`
	Message    string `xml:"message,omitempty"`
}

type DependencyManagement struct {
	Dependencies *[]Dependency `xml:"dependencies>dependency,omitempty"`
}

type Dependency struct {
	GroupID    string       `xml:"groupId,omitempty"`
	ArtifactID string       `xml:"artifactId,omitempty"`
	Version    string       `xml:"version,omitempty"`
	Type       string       `xml:"type,omitempty"`
	Classifier string       `xml:"classifier,omitempty"`
	Scope      string       `xml:"scope,omitempty"`
	SystemPath string       `xml:"systemPath,omitempty"`
	Exclusions *[]Exclusion `xml:"exclusions>exclusion"`
	Optional   string       `xml:"optional,omitempty"`
}

type Exclusion struct {
	ArtifactID string `xml:"artifactId,omitempty"`
	GroupID    string `xml:"groupId,omitempty"`
}

type Repository struct {
	UniqueVersion bool              `xml:"uniqueVersion,omitempty"`
	Releases      *RepositoryPolicy `xml:"releases"`
	Snapshots     *RepositoryPolicy `xml:"snapshots"`
	ID            string            `xml:"id,omitempty"`
	Name          string            `xml:"name,omitempty"`
	URL           string            `xml:"url,omitempty"`
	Layout        string            `xml:"layout,omitempty"`
}

type RepositoryPolicy struct {
	Enabled        string `xml:"enabled,omitempty"`
	UpdatePolicy   string `xml:"updatePolicy,omitempty"`
	ChecksumPolicy string `xml:"checksumPolicy,omitempty"`
}

type PluginRepository struct {
	Releases  *RepositoryPolicy `xml:"releases"`
	Snapshots *RepositoryPolicy `xml:"snapshots"`
	ID        string            `xml:"id,omitempty"`
	Name      string            `xml:"name,omitempty"`
	URL       string            `xml:"url,omitempty"`
	Layout    string            `xml:"layout,omitempty"`
}

type BuildBase struct {
	DefaultGoal      string            `xml:"defaultGoal,omitempty"`
	Resources        *[]Resource       `xml:"resources>resource"`
	TestResources    *[]Resource       `xml:"testResources>testResource"`
	Directory        string            `xml:"directory,omitempty"`
	FinalName        string            `xml:"finalName,omitempty"`
	Filters          []string          `xml:"filters>filter"`
	PluginManagement *PluginManagement `xml:"pluginManagement"`
	Plugins          *[]Plugin         `xml:"plugins>plugin"`
}

type Build struct {
	SourceDirectory       string       `xml:"sourceDirectory,omitempty"`
	ScriptSourceDirectory string       `xml:"scriptSourceDirectory,omitempty"`
	TestSourceDirectory   string       `xml:"testSourceDirectory,omitempty"`
	OutputDirectory       string       `xml:"outputDirectory,omitempty"`
	TestOutputDirectory   string       `xml:"testOutputDirectory,omitempty"`
	Extensions            *[]Extension `xml:"extensions>extension"`
	BuildBase
}

type Extension struct {
	GroupID    string `xml:"groupId,omitempty"`
	ArtifactID string `xml:"artifactId,omitempty"`
	Version    string `xml:"version,omitempty"`
}

type Resource struct {
	TargetPath string   `xml:"targetPath,omitempty"`
	Filtering  string   `xml:"filtering,omitempty"`
	Directory  string   `xml:"directory,omitempty"`
	Includes   []string `xml:"includes>include"`
	Excludes   []string `xml:"excludes>exclude"`
}

type PluginManagement struct {
	Plugins *[]Plugin `xml:"plugins>plugin"`
}

type Plugin struct {
	GroupID      string             `xml:"groupId,omitempty"`
	ArtifactID   string             `xml:"artifactId,omitempty"`
	Version      string             `xml:"version,omitempty"`
	Extensions   string             `xml:"extensions,omitempty"`
	Executions   *[]PluginExecution `xml:"executions>execution"`
	Dependencies *[]Dependency      `xml:"dependencies>dependency"`
	Inherited    string             `xml:"inherited,omitempty"`
}

type PluginExecution struct {
	ID        string   `xml:"id,omitempty"`
	Phase     string   `xml:"phase,omitempty"`
	Goals     []string `xml:"goals>goal"`
	Inherited string   `xml:"inherited,omitempty"`
}

type Reporting struct {
	ExcludeDefaults string             `xml:"excludeDefaults,omitempty"`
	OutputDirectory string             `xml:"outputDirectory,omitempty"`
	Plugins         *[]ReportingPlugin `xml:"plugins>plugin"`
}

type ReportingPlugin struct {
	GroupID    string       `xml:"groupId,omitempty"`
	ArtifactID string       `xml:"artifactId,omitempty"`
	Version    string       `xml:"version,omitempty"`
	Inherited  string       `xml:"inherited,omitempty"`
	ReportSets *[]ReportSet `xml:"reportSets>reportSet"`
}

type ReportSet struct {
	ID        string   `xml:"id,omitempty"`
	Reports   []string `xml:"reports>report"`
	Inherited string   `xml:"inherited,omitempty"`
}

type Profile struct {
	ID                     string                  `xml:"id,omitempty"`
	Activation             *Activation             `xml:"activation"`
	Build                  *BuildBase              `xml:"build"`
	Modules                []string                `xml:"modules>module"`
	DistributionManagement *DistributionManagement `xml:"distributionManagement"`
	Properties             *Properties             `xml:"properties"`
	DependencyManagement   *DependencyManagement   `xml:"dependencyManagement"`
	Dependencies           *[]Dependency           `xml:"dependencies>dependency"`
	Repositories           *[]Repository           `xml:"repositories>repository"`
	PluginRepositories     *[]PluginRepository     `xml:"pluginRepositories>pluginRepository"`
	Reporting              *Reporting              `xml:"reporting"`
}

type Activation struct {
	ActiveByDefault bool                `xml:"activeByDefault,omitempty"`
	JDK             string              `xml:"jdk,omitempty"`
	OS              *ActivationOS       `xml:"os"`
	Property        *ActivationProperty `xml:"property"`
	File            *ActivationFile     `xml:"file"`
}

type ActivationOS struct {
	Name    string `xml:"name,omitempty"`
	Family  string `xml:"family,omitempty"`
	Arch    string `xml:"arch,omitempty"`
	Version string `xml:"version,omitempty"`
}

type ActivationProperty struct {
	Name  string `xml:"name,omitempty"`
	Value string `xml:"value,omitempty"`
}

type ActivationFile struct {
	Missing string `xml:"missing,omitempty"`
	Exists  string `xml:"exists,omitempty"`
}
