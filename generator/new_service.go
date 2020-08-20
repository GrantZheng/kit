package generator

import (
	"fmt"
	"os/exec"
	"path"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/kujtimiihoxha/kit/fs"
	"github.com/kujtimiihoxha/kit/utils"
	"github.com/spf13/viper"
)

// NewService implements Gen and is used to create a new service.
type NewService struct {
	BaseGenerator
	name          string
	interfaceName string
	destPath      string
	filePath      string
}

// NewNewService returns a initialized and ready generator.
//
// The name parameter is the name of the service that will be created
// this name should be without the `Service` suffix
func NewNewService(name string) Gen {
	gs := &NewService{
		name:          name,
		interfaceName: utils.ToCamelCase(name + "Service"),
		destPath:      fmt.Sprintf(viper.GetString("gk_service_path_format"), utils.ToLowerSnakeCase(name)),
	}
	gs.filePath = path.Join(gs.destPath, viper.GetString("gk_service_file_name"))
	gs.srcFile = jen.NewFilePath(strings.Replace(gs.destPath, "\\", "/", -1))
	gs.InitPg()
	gs.fs = fs.Get()
	return gs
}

// Generate will run the generator.
func (g *NewService) Generate() error {
	g.CreateFolderStructure(g.destPath)
	err := g.genModule()
	if err != nil {
		println(err.Error())
		return err
	}

	comments := []string{
		"Add your methods here",
		"e.x: Foo(ctx context.Context,s string)(rs string, err error)",
	}
	partial := NewPartialGenerator(nil)
	partial.appendMultilineComment(comments)
	g.code.Raw().Commentf("%s describes the service.", g.interfaceName).Line()
	g.code.appendInterface(
		g.interfaceName,
		[]jen.Code{partial.Raw()},
	)
	return g.fs.WriteFile(g.filePath, g.srcFile.GoString(), false)
}

func (g *NewService) genModule() error {
	prjName := utils.ToLowerSnakeCase(g.name)
	exist, _ := g.fs.Exists(prjName + "/go.mod")
	if exist {
		return nil
	}

	moduleName := prjName
	if utils.IsInGopath() {
		moduleName = ""
	}

	if viper.GetString("n_s_module") != "" {
		moduleName = viper.GetString("n_s_module")
		moduleNameSlice := strings.Split(moduleName, "/")
		moduleNameSlice[len(moduleNameSlice)-1] = utils.ToLowerSnakeCase(moduleNameSlice[len(moduleNameSlice)-1])
		moduleName = strings.Join(moduleNameSlice, "/")
	}
	cmdStr := "cd " + prjName + " && go mod init " + moduleName
	cmd := exec.Command("sh", "-c", cmdStr)
	_, err := cmd.Output()
	return err
}
