package govm

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/c4milo/unpackit"
)

const (
	goVMDir             string = ".govm"
	defaultRegistryPath string = "https://golang.org/dl/"
	goVMDownloadUrl     string = "https://github.com/TaceyWong/govm/releases/latest/download/"
	goVMTagsApi         string = "https://raw.githubusercontent.com/TaceyWong/govm/go-tags/tags.json"
)

// Command ...
type Command interface {
	ListVersions()
	ListRemoteVersions()
	CurrentVersion() string
	Uninstall(version string)
	Install(version string)
	Use(version string)
	Upgrade()
	Helper
}

// GoVM struct
type GoVM struct {
	homeDir       string
	installDir    string
	versionsDir   string
	currentDir    string
	currentBinDir string
	currentGoDir  string
	downloadsDir  string
	Command
}

// Helper ...
type Helper interface {
	getArch() string
	existsVersion(version string) bool
	cleanVersionDir(version string)
	mkdirs(version string)
	getVersionDir(version string) string
	downloadAndExtract(version string)
	changeSymblinkGoBin(version string)
	changeSymblinkGo(version string)
	getLatestVersion() string
	getGithubTags(repo string) (result []string)
}

var gvm GoVM
var githubTags map[string][]string

// NewVM instance
func NewGoVm() GoVM {
	gvm.homeDir = os.Getenv("HOME")
	gvm.installDir = filepath.Join(gvm.homeDir, goVMDir)
	gvm.versionsDir = filepath.Join(gvm.installDir, "versions")
	gvm.currentDir = filepath.Join(gvm.installDir, "current")
	gvm.currentBinDir = filepath.Join(gvm.installDir, "current", "bin")
	gvm.currentGoDir = filepath.Join(gvm.installDir, "current", "go")
	gvm.downloadsDir = filepath.Join(gvm.installDir, "downloads")

	return gvm
}

func (g *GoVM) getArch() string {
	return runtime.GOOS + "-" + runtime.GOARCH
}

// ListVersions that are installed by dir ls
// highlight the version that is currently symbolic linked
func (g *GoVM) ListVersions() error {
	entries, err := os.ReadDir(g.versionsDir)
	CheckError(err, "[Error]: List versions failed")
	files := make([]fs.FileInfo, 0, len(entries))
	for _, entry := range entries {
		info, err := entry.Info()
		CheckError(err, "[Error]: List versions failed")
		files = append(files, info)
	}

	cv := g.CurrentVersion()

	versionsSemantic := make([]*semver.Version, 0)

	for _, f := range files {
		if v, err := semver.NewVersion(f.Name()); err == nil {
			versionsSemantic = append(versionsSemantic, v)
		}
	}

	// sort semantic versions
	sort.Sort(semver.Collection(versionsSemantic))

	for _, versionSemantic := range versionsSemantic {
		version := versionSemantic.String()
		// 1.8.0 -> 1.8
		reMajorVersion, _ := regexp.Compile("[0-9]+.[0-9]+.0")
		if reMajorVersion.MatchString(version) {
			version = strings.Split(version, ".")[0] + "." + strings.Split(version, ".")[1]
		}

		if version == cv {
			version = cv + "*"
			Successln(version)
		} else {
			log.Println(version)
		}
	}

	// print rc and beta versions in the end
	for _, f := range files {
		rcVersion := f.Name()
		r, _ := regexp.Compile("beta.*|rc.*")
		matches := r.FindAllString(rcVersion, -1)
		if len(matches) == 1 {
			if rcVersion == cv {
				rcVersion = cv + "*"
				Successln(rcVersion)
			} else {
				log.Println(rcVersion)
			}
		}
	}

	if cv != "" {
		log.Println()
		log.Printf("current: %s", cv)
	}
	return nil
}

// ListRemoteVersions that are installed by dir ls
func (g *GoVM) ListRemoteVersions(print bool) map[string][]string {
	log.Println("[Info]: Fetching remote versions")
	tags := g.getGithubTags("golang/go")

	var versions []string
	for _, tag := range tags {
		versions = append(versions, strings.ReplaceAll(tag, "go", ""))
	}

	return g.getGroupedVersion(versions, print)
}

func (g *GoVM) getGroupedVersion(versions []string, print bool) map[string][]string {
	groupedVersions := make(map[string][]string)
	for _, version := range versions {
		parts := strings.Split(version, ".")
		if len(parts) > 1 {
			majorVersion := fmt.Sprintf("%s.%s", parts[0], parts[1])
			r, _ := regexp.Compile("beta.*|rc.*")
			matches := r.FindAllString(majorVersion, -1)
			if len(matches) == 1 {
				majorVersion = strings.Split(version, matches[0])[0]
			}
			groupedVersions[majorVersion] = append(groupedVersions[majorVersion], version)
		}
	}

	// groupedVersionKeys := []string{"1", "1.1", "1.2", ..., "1.17"}
	groupedVersionKeys := make([]string, 0, len(groupedVersions))
	for groupedVersionKey := range groupedVersions {
		groupedVersionKeys = append(groupedVersionKeys, groupedVersionKey)
	}

	versionsSemantic := make([]*semver.Version, 0)
	for _, r := range groupedVersionKeys {
		if v, err := semver.NewVersion(r); err == nil {
			versionsSemantic = append(versionsSemantic, v)
		}
	}

	// sort semantic versions
	sort.Sort(semver.Collection(versionsSemantic))

	// match 1.0.0 or 2.0.0
	reTopVersion, _ := regexp.Compile("[0-9]+.0.0")

	for _, versionSemantic := range versionsSemantic {
		maxPerLine := 0
		strKey := versionSemantic.String()
		lookupKey := ""
		versionParts := strings.Split(strKey, ".")

		// prepare lookup key for the grouped version map.
		// 1.0.0 -> 1.0, 1.1.1 -> 1.1
		lookupKey = versionParts[0] + "." + versionParts[1]
		// On match 1.0.0, print 1. On match 2.0.0 print 2
		if reTopVersion.MatchString(strKey) {
			if print {
				Major(versionParts[0])
			}
			g.print("\t", print)
		} else {
			if print {
				Major(lookupKey)
			}
			g.print("\t", print)
		}

		groupedVersionsSemantic := make([]*semver.Version, 0)
		for _, r := range groupedVersions[lookupKey] {
			if v, err := semver.NewVersion(r); err == nil {
				groupedVersionsSemantic = append(groupedVersionsSemantic, v)
			}

		}
		// sort semantic versions
		sort.Sort(semver.Collection(groupedVersionsSemantic))

		for _, gvSemantic := range groupedVersionsSemantic {
			maxPerLine++
			if maxPerLine == 6 {
				maxPerLine = 0
				g.print("\n\t", print)
			}
			g.print(gvSemantic.String()+"  ", print)
		}

		// print rc and beta versions in the end
		for _, rcVersion := range groupedVersions[lookupKey] {
			r, _ := regexp.Compile("beta.*|rc.*")
			matches := r.FindAllString(rcVersion, -1)
			if len(matches) == 1 {
				g.print(rcVersion+"  ", print)
				maxPerLine++
				if maxPerLine == 6 {
					maxPerLine = 0
					g.print("\n\t", print)
				}
			}
		}
		g.print("\n", print)
		g.print("\n", print)
	}
	return groupedVersions
}

func (g *GoVM) print(message string, shouldPrint bool) {
	if shouldPrint {
		fmt.Print(message)
	}
}

func (g *GoVM) existsVersion(version string) bool {
	path := filepath.Join(g.versionsDir, version, "go")
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

// CurrentVersion get current version from symb link
func (g *GoVM) CurrentVersion() string {

	fp, err := filepath.EvalSymlinks(g.currentBinDir)
	if err != nil {
		return ""
	}

	version := strings.ReplaceAll(fp, "/go/bin", "")
	version = strings.ReplaceAll(version, g.versionsDir, "")
	version = strings.ReplaceAll(version, "/", "")
	return version
}

// Uninstall the given version of go
func (g *GoVM) Uninstall(version string) {
	if version == "" {
		log.Fatal("[Error] No version provided")
	}
	if g.CurrentVersion() == version {
		Errorf("[Error] Version: %s you are trying to remove is your current version. Please use a different version first before uninstalling the current version\n", version)
		os.Exit(1)
	}
	if !g.existsVersion(version) {
		Errorf("[Error] Version: %s you are trying to remove is not installed\n", version)
		os.Exit(1)
	}
	g.cleanVersionDir(version)
	Successf("[Success] Version: %s uninstalled\n", version)
}

func (g *GoVM) cleanVersionDir(version string) {
	_ = os.RemoveAll(g.getVersionDir(version))
}

func (g *GoVM) cleanDownloadsDir() {
	_ = os.RemoveAll(g.downloadsDir)
}

// Install the given version of go
func (g *GoVM) Install(version string) {
	if version == "" {
		log.Fatal("[Error] No version provided")
	}
	version = g.judgeVersion(version)
	g.mkdirs(version)
	if g.existsVersion(version) {
		Infof("[Info] Version: %s exists \n", version)
		return
	}

	Infof("[Info] Downloading version: %s \n", version)
	g.downloadAndExtract(version)
	g.cleanDownloadsDir()
	Successf("[Success] Downloaded version: %s\n", version)
}

func (g *GoVM) judgeVersion(version string) string {
	judgedVersion := ""
	rcBetaOk := false
	reRcOrBeta, _ := regexp.Compile("beta.*|rc.*")
	// check if version string ends with x

	if strings.HasSuffix(version, "x") {
		judgedVersion = version[:len(version)-1]
	}

	if strings.HasSuffix(version, ".x") {
		judgedVersion = version[:len(version)-2]
	}
	if strings.HasSuffix(version, "@latest") {
		judgedVersion = version[:len(version)-7]
	}
	if strings.HasSuffix(version, "@dev-latest") {
		judgedVersion = version[:len(version)-11]
		rcBetaOk = true
	}

	if version == "latest" || version == "dev-latest" {
		groupedVersions := g.ListRemoteVersions(false) // donot print
		groupedVersionKeys := make([]string, 0, len(groupedVersions))
		for groupedVersionKey := range groupedVersions {
			groupedVersionKeys = append(groupedVersionKeys, groupedVersionKey)
		}
		versionsSemantic := make([]*semver.Version, 0)
		for _, r := range groupedVersionKeys {
			if v, err := semver.NewVersion(r); err == nil {
				versionsSemantic = append(versionsSemantic, v)
			}
		}

		// sort semantic versions
		sort.Sort(semver.Collection(versionsSemantic))
		// loop in reverse
		for i := len(versionsSemantic) - 1; i >= 0; i-- {
			judgedVersions := groupedVersions[versionsSemantic[i].Original()]
			// get last element
			if version == "dev-latest" {
				return judgedVersions[len(judgedVersions)-1]
			}

			// loop in reverse
			for j := len(judgedVersions) - 1; j >= 0; j-- {
				matches := reRcOrBeta.FindAllString(judgedVersions[j], -1)
				if len(matches) == 0 {
					return judgedVersions[j]
				}
			}
		}

		latest := versionsSemantic[len(versionsSemantic)-1].String()
		return g.judgeVersion(latest)
	}

	if judgedVersion != "" {
		groupedVersions := g.ListRemoteVersions(false) // donot print
		// check if judgedVersion is in the groupedVersions
		if _, ok := groupedVersions[judgedVersion]; ok {
			// get last item in the groupedVersions excluding rc and beta
			// loop in reverse groupedVersions
			for i := len(groupedVersions[judgedVersion]) - 1; i >= 0; i-- {
				matches := reRcOrBeta.FindAllString(groupedVersions[judgedVersion][i], -1)
				if len(matches) == 0 {
					return groupedVersions[judgedVersion][i]
				}
			}
			if rcBetaOk {
				// return last element including beta and rc if present
				return groupedVersions[judgedVersion][len(groupedVersions[judgedVersion])-1]
			}
		}
	}

	return version
}

// Use a version
func (g *GoVM) Use(version string) {
	version = g.judgeVersion(version)
	if g.CurrentVersion() == version {
		Infof("[Info] Version: %s is already your current version \n", version)
		return
	}
	Infof("[Info] Changing go version to: %s \n", version)
	g.changeSymblinkGoBin(version)
	g.changeSymblinkGo(version)
	Successf("[Success] Changed go version to: %s\n", version)
}

// Upgrade of GoBrew
func (g *GoVM) Upgrade(currentVersion string) {
	if "v"+currentVersion == g.getLatestVersion() {
		Infoln("[INFO] your version is already newest")
		return
	}

	mkdirTemp, _ := os.MkdirTemp("", "govm")
	tmpFile := filepath.Join(mkdirTemp, "gobrew")
	url := goVMDownloadUrl + "govm-" + g.getArch()
	if err := DownloadWithProgress(url, "gobrew", mkdirTemp); err != nil {
		Errorln("[Error] Download GoBrew failed:", err)
		return
	}

	source, err := os.Open(tmpFile)
	if err != nil {
		Errorln("[Error] Cannot open file", err)
		return
	}
	defer func(source *os.File) {
		_ = source.Close()
	}(source)

	goBrewFile := filepath.Join(g.installDir, "/bin/gobrew")
	destination, err := os.Create(goBrewFile)
	if err != nil {
		Errorf("[Error] Cannot open file: %s", err)
		return
	}
	defer func(destination *os.File) {
		_ = destination.Close()
	}(destination)

	if _, err = io.Copy(destination, source); err != nil {
		Errorf("[Error] Cannot copy file: %s", err)
		return
	}

	if err = os.Chmod(goBrewFile, 0755); err != nil {
		Errorf("[Error] Cannot set file as executable: %s", err)
		return
	}

	if err = os.Remove(tmpFile); err != nil {
		Errorf("[Error] Cannot remove tmp file: %s", err)
		return
	}

	Infoln("Upgrade successful")
}

func (g *GoVM) mkdirs(version string) {
	_ = os.MkdirAll(g.installDir, os.ModePerm)
	_ = os.MkdirAll(g.currentDir, os.ModePerm)
	_ = os.MkdirAll(g.versionsDir, os.ModePerm)
	_ = os.MkdirAll(g.getVersionDir(version), os.ModePerm)
	_ = os.MkdirAll(g.downloadsDir, os.ModePerm)
}

func (g *GoVM) getVersionDir(version string) string {
	return filepath.Join(g.versionsDir, version)
}
func (g *GoVM) downloadAndExtract(version string) {
	tarName := "go" + version + "." + g.getArch() + ".tar.gz"

	registryPath := defaultRegistryPath
	if p := os.Getenv("GOBREW_REGISTRY"); p != "" {
		registryPath = p
	}
	downloadURL := registryPath + tarName
	Infof("[Info] Downloading from: %s \n", downloadURL)

	dstDownloadDir := filepath.Join(g.downloadsDir)
	Infof("[Info] Downloading to: %s \n", dstDownloadDir)
	err := DownloadWithProgress(downloadURL, tarName, dstDownloadDir)

	if err != nil {
		g.cleanVersionDir(version)
		Infof("[Info]: Downloading version failed: %s \n", err)
		Errorf("[Error]: Please check connectivity to url: %s\n", downloadURL)
		os.Exit(1)
	}

	srcTar := filepath.Join(g.downloadsDir, tarName)
	dstDir := g.getVersionDir(version)

	Infof("[Info] Extracting from: %s \n", srcTar)
	Infof("[Info] Extracting to: %s \n", dstDir)

	err = g.ExtractTarGz(srcTar, dstDir)
	if err != nil {
		// clean up dir
		g.cleanVersionDir(version)
		Infof("[Info]: Untar failed: %s \n", err)
		Errorf("[Error]: Please check if version exists from url: %s\n", downloadURL)
		os.Exit(1)
	}
	Infof("[Success] Untar to %s\n", g.getVersionDir(version))
}

func (g *GoVM) ExtractTarGz(srcTar string, dstDir string) error {
	//#nosec G304
	file, err := os.Open(srcTar)
	if err != nil {
		return err
	}
	_, err = unpackit.Unpack(file, dstDir)
	if err != nil {
		return err
	}

	return nil
}

func (g *GoVM) changeSymblinkGoBin(version string) {
	goBinDst := filepath.Join(g.versionsDir, version, "/go/bin")
	_ = os.RemoveAll(g.currentBinDir)

	if err := os.Symlink(goBinDst, g.currentBinDir); err != nil {
		Errorf("[Error]: symbolic link failed: %s\n", err)
		os.Exit(1)
	}
}

func (g *GoVM) changeSymblinkGo(version string) {
	_ = os.RemoveAll(g.currentGoDir)
	versionGoDir := filepath.Join(g.versionsDir, version, "go")

	if err := os.Symlink(versionGoDir, g.currentGoDir); err != nil {
		Errorf("[Error]: symbolic link failed: %s\n", err)
		os.Exit(1)
	}
}

func (g *GoVM) getLatestVersion() string {
	tags := g.getGithubTags("kevincobain2000/gobrew")

	if len(tags) == 0 {
		return ""
	}

	return tags[len(tags)-1]
}

func (g *GoVM) getGithubTags(repo string) (result []string) {
	if len(githubTags[repo]) > 0 {
		return githubTags[repo]
	}

	githubTags = make(map[string][]string, 0)
	client := &http.Client{}
	url := "https://api.github.com/repos/TaceyWong/govm/git/refs/tags"
	if repo == "golang/go" {
		url = goVMTagsApi
	}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		Errorf("[Error] Cannot create request: %s", err)
		return
	}

	request.Header.Set("User-Agent", "govm")

	response, err := client.Do(request)
	if err != nil {
		Errorf("[Error] Cannot get response: %s", err)
		return
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)

	data, err := io.ReadAll(response.Body)
	if err != nil {
		Errorf("[Error] Cannot read response: %s", err)
		return
	}

	type Tag struct {
		Ref string
	}
	var tags []Tag

	if err := json.Unmarshal(data, &tags); err != nil {
		Errorf("[Error] Rate limit exceed")
		os.Exit(2)
	}

	for _, tag := range tags {
		t := strings.ReplaceAll(tag.Ref, "refs/tags/", "")
		if strings.HasPrefix(t, "v") || strings.HasPrefix(t, "go") {
			result = append(result, t)
		}
	}

	githubTags[repo] = result
	return result
}



