package yamlOps

import (
	"gopkg.in/yaml.v3"
)

type SwapperFunc func(value string) (replacement string, err error)

type Walker struct {
	swapperFunc SwapperFunc
}

func WalkFile(inB []byte, swapper SwapperFunc) (outB []byte, err error) {
	walker := Walker{
		swapperFunc: swapper,
	}
	root := new(yaml.Node)
	if err = yaml.Unmarshal(inB, root); err == nil {
		if err = walker.walkYaml(root); err == nil {
			return yaml.Marshal(root)
		}
	}
	return
}

func (c *Walker) walkYaml(node *yaml.Node) (err error) {
	switch node.Kind {
	case yaml.ScalarNode:
		err = c.callSwapper(node)
	case yaml.DocumentNode:
		return c.walkDocument(node)
	case yaml.MappingNode:
		return c.walkMapping(node)
	case yaml.SequenceNode:
		return c.walkSequence(node)
	}
	return
}

func (c *Walker) walkDocument(node *yaml.Node) error {
	if len(node.Content) > 0 {
		return c.walkYaml(node.Content[0])
	}
	return nil
}

func (c *Walker) walkMapping(node *yaml.Node) (err error) {
	for i := 1; err == nil && i < len(node.Content); i += 2 {
		err = c.walkYaml(node.Content[i])
	}
	return
}

func (c *Walker) walkSequence(node *yaml.Node) (err error) {
	for _, v := range node.Content {
		err = c.callSwapper(v)
		if err != nil {
			return
		}
	}
	return
}

func (c *Walker) callSwapper(node *yaml.Node) (err error) {
	var replacement string
	replacement, err = c.swapperFunc(node.Value)
	if err == nil && node.Value != replacement {
		node.Value = replacement
	}
	return
}
