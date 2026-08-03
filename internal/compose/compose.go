package compose

import (
	"bytes"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/yunkhngn/dforge/internal/fsutil"
	"github.com/yunkhngn/dforge/internal/services"
	"gopkg.in/yaml.v3"
)

type File struct {
	path string
	root *yaml.Node // document node
}

func Load(path string) (*File, error) {
	var doc yaml.Node
	if data, err := os.ReadFile(path); err == nil {
		if err := yaml.Unmarshal(data, &doc); err != nil {
			return nil, err
		}
	} else {
		// empty skeleton: mapping with "services"
		doc = yaml.Node{Kind: yaml.DocumentNode, Content: []*yaml.Node{
			{Kind: yaml.MappingNode},
		}}
	}
	if len(doc.Content) == 0 {
		doc.Content = []*yaml.Node{{Kind: yaml.MappingNode}}
	}
	return &File{path: path, root: &doc}, nil
}

func (f *File) topMap() *yaml.Node { return f.root.Content[0] }

// mapValue returns the value node for key in a mapping, creating it if missing.
func mapValue(m *yaml.Node, key string, create func() *yaml.Node) *yaml.Node {
	for i := 0; i+1 < len(m.Content); i += 2 {
		if m.Content[i].Value == key {
			return m.Content[i+1]
		}
	}
	if create == nil {
		return nil
	}
	k := &yaml.Node{Kind: yaml.ScalarNode, Value: key}
	v := create()
	m.Content = append(m.Content, k, v)
	return v
}

func scalar(v string) *yaml.Node { return &yaml.Node{Kind: yaml.ScalarNode, Value: v} }

func seq(values []string) *yaml.Node {
	n := &yaml.Node{Kind: yaml.SequenceNode}
	for _, v := range values {
		n.Content = append(n.Content, scalar(v))
	}
	return n
}

func (f *File) HasService(name string) bool {
	svcs := mapValue(f.topMap(), "services", nil)
	if svcs == nil {
		return false
	}
	for i := 0; i+1 < len(svcs.Content); i += 2 {
		if svcs.Content[i].Value == name {
			return true
		}
	}
	return false
}

func (f *File) AddService(svc services.Service) error {
	if f.HasService(svc.Name) {
		return fmt.Errorf("service %q already exists", svc.Name)
	}
	svcs := mapValue(f.topMap(), "services", func() *yaml.Node { return &yaml.Node{Kind: yaml.MappingNode} })

	node := &yaml.Node{Kind: yaml.MappingNode}
	add := func(k string, v *yaml.Node) { node.Content = append(node.Content, scalar(k), v) }
	add("image", scalar(svc.Image))
	add("restart", scalar("unless-stopped"))
	if len(svc.Ports) > 0 {
		add("ports", seq(svc.Ports))
	}
	if len(svc.Env) > 0 {
		envNode := &yaml.Node{Kind: yaml.MappingNode}
		keys := make([]string, 0, len(svc.Env))
		for k := range svc.Env {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			envNode.Content = append(envNode.Content, scalar(k), scalar(svc.Env[k]))
		}
		add("environment", envNode)
	}
	if len(svc.Volumes) > 0 {
		add("volumes", seq(svc.Volumes))
	}
	if svc.Healthcheck != nil {
		hcNode := &yaml.Node{Kind: yaml.MappingNode}
		if len(svc.Healthcheck.Test) > 0 {
			hcNode.Content = append(hcNode.Content, scalar("test"), seq(svc.Healthcheck.Test))
		}
		if svc.Healthcheck.Interval != "" {
			hcNode.Content = append(hcNode.Content, scalar("interval"), scalar(svc.Healthcheck.Interval))
		}
		if svc.Healthcheck.Timeout != "" {
			hcNode.Content = append(hcNode.Content, scalar("timeout"), scalar(svc.Healthcheck.Timeout))
		}
		if svc.Healthcheck.Retries > 0 {
			hcNode.Content = append(hcNode.Content, scalar("retries"), scalar(fmt.Sprintf("%d", svc.Healthcheck.Retries)))
		}
		add("healthcheck", hcNode)
	}
	svcs.Content = append(svcs.Content, scalar(svc.Name), node)

	// register named volumes at top level
	for _, v := range svc.Volumes {
		if name, ok := namedVolume(v); ok {
			f.ensureTopVolume(name)
		}
	}
	return nil
}

func namedVolume(mapping string) (string, bool) {
	parts := strings.SplitN(mapping, ":", 2)
	if len(parts) == 2 && !strings.HasPrefix(parts[0], ".") && !strings.HasPrefix(parts[0], "/") {
		return parts[0], true
	}
	return "", false
}

func (f *File) ensureTopVolume(name string) {
	vols := mapValue(f.topMap(), "volumes", func() *yaml.Node { return &yaml.Node{Kind: yaml.MappingNode} })
	for i := 0; i+1 < len(vols.Content); i += 2 {
		if vols.Content[i].Value == name {
			return
		}
	}
	vols.Content = append(vols.Content, scalar(name), &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!null"})
}

func (f *File) RemoveService(name string) error {
	svcs := mapValue(f.topMap(), "services", nil)
	if svcs == nil {
		return fmt.Errorf("no services defined")
	}
	for i := 0; i+1 < len(svcs.Content); i += 2 {
		if svcs.Content[i].Value == name {
			svcs.Content = append(svcs.Content[:i], svcs.Content[i+2:]...)
			f.pruneOrphanVolumes()
			return nil
		}
	}
	return fmt.Errorf("service %q not found", name)
}

// pruneOrphanVolumes removes top-level named volumes no longer referenced by any service.
func (f *File) pruneOrphanVolumes() {
	vols := mapValue(f.topMap(), "volumes", nil)
	if vols == nil {
		return
	}
	used := map[string]bool{}
	svcs := mapValue(f.topMap(), "services", nil)
	if svcs != nil {
		for i := 1; i < len(svcs.Content); i += 2 {
			if v := mapValue(svcs.Content[i], "volumes", nil); v != nil {
				for _, item := range v.Content {
					if name, ok := namedVolume(item.Value); ok {
						used[name] = true
					}
				}
			}
		}
	}
	kept := vols.Content[:0]
	for i := 0; i+1 < len(vols.Content); i += 2 {
		if used[vols.Content[i].Value] {
			kept = append(kept, vols.Content[i], vols.Content[i+1])
		}
	}
	vols.Content = kept
}

func (f *File) Bytes() ([]byte, error) {
	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)
	if err := enc.Encode(f.root); err != nil {
		return nil, err
	}
	enc.Close()
	return buf.Bytes(), nil
}

func (f *File) Save() error {
	data, err := f.Bytes()
	if err != nil {
		return err
	}
	return fsutil.WriteWithBackup(f.path, data)
}
