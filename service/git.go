package service

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"
	"tag/config"
	"tag/print"
	"tag/utils"
)

type GitSvc struct{}

var gitSvc *GitSvc

func NewGitSvc() *GitSvc {
	if gitSvc == nil {
		gitSvc = new(GitSvc)
	}
	return gitSvc
}

const (
	TagNumMax   = 30         // 最大保留tag数量
	VersionFile = ".version" // 最新tag存储路径
)

/**
git fetch --tags
git tag --list
git push --tags
git status --short
*/

func (p *GitSvc) Release() error {
	// 检查是否git项目根目录
	if !p.IsGitProject() {
		print.Error("请切换到Git项目根目录重试")
		return nil
	}
	// 检查代码是否需要提交
	undo, err := p.run("git status --short")
	if err != nil {
		print.Error("非Git项目目录 " + err.Error())
		return err
	}
	if undo != "" {
		print.Warn("请先提交最新改动：")
		print.Print(undo)
		return nil
	}

	// 更新全部tag
	p.run("git fetch --tags")

	// 暂存已存在tags
	exists := make([]string, 0)

	// 获取全部tag
	res, _ := p.run("git tag --list")
	if res != "" {
		buf := strings.Split(res, "\n")
		offset := 0
		if len(buf) > TagNumMax {
			offset = len(buf) - TagNumMax
		}
		if offset > 0 {
			expired := buf[:offset]
			for _, tag := range expired {
				p.deleteTag(tag)
			}
		}
		exists = append(exists, buf...)
	}

	// 打tag
	if err = p.newTag(exists); err == nil {
		print.Info("success!")
	}

	return nil
}

func (p *GitSvc) DeleteTag() error {
	tag := config.NewMemConfig().GetDef("tag", "")
	tag = strings.Replace(tag, " ", "", -1)
	return p.deleteTag(tag)
}

// 删除tag 1.删除本地tag； 2.删除远端tag
func (p *GitSvc) deleteTag(tag string) error {
	if tag == "" {
		return nil
	}
	print.Info("delete tag: " + tag)

	p.run("git tag -d " + tag)
	p.run("git push origin :refs/tags/" + tag)
	return nil
}

// 创建新tag 1.tag命令是否已存在 2.更新tag文件
func (p *GitSvc) newTag(exists []string) error {
	var err error

	tag := config.NewMemConfig().GetDef("tag", "")
	if tag == "" {
		tag, err = p.genTag(exists)
		if err != nil {
			return err
		}
	}
	print.Print("release tag:" + tag)

	// 提交 .version文件
	p.run("git add " + VersionFile)
	p.run("git commit -m " + tag)

	_, err = p.run(fmt.Sprintf(`git tag -a "%s" -m "%s"`, tag, config.NewMemConfig().GetDef("comment", "release")))
	if err != nil {
		print.Error("create tag failed:" + err.Error())
		return err
	}

	_, err = p.run("git push origin " + tag)
	if err != nil {
		print.Error("push tags failed:" + err.Error())
		return err
	}

	branch, err := p.run("git symbolic-ref --short -q HEAD")
	if err != nil {
		print.Error("get branch failed:" + err.Error())
		return err
	}
	p.run("git push origin " + branch)

	return nil
}

// 创建tag字符串
func (p *GitSvc) genTag(exists []string) (string, error) {
	oldTag := ""
	data, err := os.ReadFile(VersionFile)
	if err != nil {
		oldTag = "v0.0.0"
	} else {
		oldTag = string(data)
		if oldTag == "" {
			oldTag = "v0.0.0"
		}
	}

	vnum := strings.TrimLeft(strings.TrimSpace(oldTag), "v")
	buf := strings.Split(vnum, ".")

	// 如果添加version_ext，则不更新版本号
	ext := config.NewMemConfig().GetDef("version_ext", "")
	if ext != "" {
		tag := fmt.Sprintf("v%s.%s.%s.%s", buf[0], buf[1], buf[2], ext)
		if utils.StringInSlice(exists, tag) {
			print.Errorf("tag exists:%s", tag)
			return "", errors.New("tag exists")
		}
		p.saveTag(VersionFile, tag)
		return tag, nil
	}

	// 循环生成，排重
	version_type := config.NewMemConfig().GetIntDef("version_type", 3)
	v1 := utils.ToInt(buf[0])
	v2 := utils.ToInt(buf[1])
	v3 := utils.ToInt(buf[2])

	switch version_type {
	case 1:
		v1 += 1
		v2 += 1
		v3 += 1
	case 2:
		v2 += 1
		v3 += 1
	case 3:
		v3 += 1
	}

	tag := ""
	for {
		tag = fmt.Sprintf("v%d.%d.%d", v1, v2, v3)
		if utils.StringInSlice(exists, tag) {
			v3 += 1
			continue
		}
		break
	}
	p.saveTag(VersionFile, tag)
	return tag, nil
}

func (p *GitSvc) saveTag(filename, tag string) error {
	if err := os.WriteFile(filename, []byte(tag), 0666); err != nil {
		print.Errorf("update .version failed:%s", err.Error())
		return err
	}
	return nil
}

// 执行命令,安全起见，禁止外部访问
func (p *GitSvc) run(command string) (string, error) {
	var err error
	var out string
	if runtime.GOOS == "windows" {
		out, err = utils.PowerShell(command + "\n")
	} else {
		out, err = utils.ExecShell(command)
	}
	if err != nil {
		return out, err
	}
	return out, nil
}

// 判断是否git仓库目录
func (p *GitSvc) IsGitProject() bool {
	f := ".git/config"
	_, err := os.Stat(f)
	if err == nil {
		return true
	}

	if _, ok := err.(*os.PathError); ok {
		return false
	}
	return false
}

func (p *GitSvc) Status() (string, error) {
	c := "git status"
	out, err := utils.ExecShell(c)
	if err != nil {
		print.Error(err.Error())
		return out, err
	}
	return out, nil
}
