// Code generated. DO NOT EDIT.

package {{ .package }}

import (
	"testing"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	watchtools "k8s.io/client-go/tools/watch"
	"github.tools.sap/cs-devops/kubernetes-testing/framework"
	{{- $imports := list }}
    {{- range .groupVersions }}
    {{- $alias := printf "_%s" .name }}
    {{- $imports = append $imports (printf "%s \"%s\"" $alias .import) }}
    {{- end }}
	{{- range .resources }}
	{{- $alias := printf "_%s%s" (default "core" .group) .version | replace "-" "" | replace "." "" }}
	{{- $imports = append $imports (printf "%s \"%s\"" $alias .import) }}
	{{- end }}
	{{- range (uniq $imports) }}
	{{ . }}
	{{- end }}
)

type Environment struct {
	framework.Environment
	{{- range .groupVersions }}
	{{- if .client }}
	{{ .name }}Client *_{{ .name }}.Clientset
	{{- end }}
	{{- end }}
}

func NewEnvironment() *Environment {
	{{- range .groupVersions }}
	{{- if .client }}
	{{ .name }}Client := _{{ .name }}.NewSimpleClientset()
	{{- end }}
	{{- end }}
	groupVersions := []*framework.GroupVersion{
		{{- range .groupVersions }}
		{{- if .client }}
		framework.NewGroupVersion(_{{ .name }}.AddToScheme, {{ .name }}Client),
		{{- else }}
		framework.NewGroupVersion(_{{ .name }}.AddToScheme, nil),
		{{- end }}
		{{- end }}
	}
	return &Environment{
		Environment: framework.NewEnvironment(groupVersions),
		{{- range .groupVersions }}
		{{- if .client }}
		{{ .name }}Client: {{ .name }}Client,
		{{- end }}
		{{- end }}
	}
}

type Must struct {
	env          *Environment
	errorHandler func(error)
}

func (must *Must) handleError(err error) {
	if err != nil {
		must.errorHandler(err)
	}
}

func (env *Environment) Must(errorHandler func(error)) *Must {
	return &Must{env: env, errorHandler: errorHandler}
}

func (env *Environment) MustError(t *testing.T) *Must {
	return env.Must(func(err error) { t.Error(err) })
}

func (env *Environment) MustFatal(t *testing.T) *Must {
	return env.Must(func(err error) { t.Fatal(err) })
}

// Client accessors
{{- range .groupVersions }}
{{- if .client }}

func (env *Environment) {{ .name | camelcase }}Client() *_{{ .name }}.Clientset {
	return env.{{ .name }}Client
}
{{- end }}
{{- end }}

{{- range .resources }}
{{- $alias := printf "_%s%s" (default "core" .group) .version | replace "-" "" | replace "." "" }}
{{- $type := printf "%s.%s" $alias .kind }}
{{- $gvk := printf "schema.GroupVersionKind{Group: \"%s\", Version: \"%s\", Kind: \"%s\"}" .group .version .kind }}
{{- $namespaceParameter := ternary "namespace string, " "" .namespaced }}
{{- $namespaceArgument := ternary "namespace, " "" .namespaced }}
{{- $namespaceValue := ternary "namespace" "\"\"" .namespaced }}

// Typed methods for {{ default "core" .group }}/{{ .version }} {{ .kind }}

func (env *Environment) Load{{ .singular }}FromFile(path string) *{{ $type }} {
	return env.LoadObjectFromFile(path).(*{{ $type }})
}

func (env *Environment) Add{{ .singular }}(obj *{{ $type }}) {
	env.AddObject(obj)
}

func (env *Environment) Add{{ .singular }}FromFile(path string) {
	env.AddObjectFromFile(path)
}

func (env *Environment) Add{{ .plural }}FromFiles(paths ...string) {
	env.AddObjectsFromFiles(paths...)
}

func (env *Environment) With{{ .singular }}(obj *{{ $type }}) *Environment {
	return env.WithObject(obj).(*Environment)
}

func (env *Environment) With{{ .singular }}FromFile(path string) *Environment {
	return env.WithObjectFromFile(path).(*Environment)
}

func (env *Environment) With{{ .plural }}FromFiles(paths ...string) *Environment {
	return env.WithObjectsFromFiles(paths...).(*Environment)
}

func (env *Environment) Assert{{ .singular }}(obj *{{ $type }}) error {
	return env.AssertObject(obj)
}

func (env *Environment) Assert{{ .singular }}FromFile(path string) error {
	return env.AssertObjectFromFile(path)
}

func (env *Environment) Assert{{ .singular }}Count({{ $namespaceParameter }}labelSelector string, count int) error {
	return env.AssertObjectCount({{ $gvk }}, {{ $namespaceValue }}, labelSelector, count)
}

func (env *Environment) Get{{ .singular }}({{ $namespaceParameter }}name string) (*{{ $type }}, error) {
	retobj, err := env.GetObject({{ $gvk }}, {{ $namespaceValue }}, name)
	if err != nil {
		return nil, err
	}
	return retobj.(*{{ $type }}), nil
}

func (env *Environment) List{{ .plural }}({{ $namespaceParameter }}labelSelector string) ([]*{{ $type }}, error) {
	retobjs, err := env.ListObjects({{ $gvk }}, {{ $namespaceValue }}, labelSelector)
	if err != nil {
		return nil, err
	}
	typedretobjs := make([]*{{ $type }}, len(retobjs))
	for i, retobj := range retobjs {
		typedretobjs[i] = retobj.(*{{ $type }})
	}
	return typedretobjs, nil
}

func (env *Environment) Create{{ .singular }}(obj *{{ $type }}) (*{{ $type }}, error) {
	retobj, err := env.CreateObject(obj)
	if err != nil {
		return nil, err
	}
	return retobj.(*{{ $type }}), nil
}

func (env *Environment) Create{{ .singular }}FromFile(path string) (*{{ $type }}, error) {
	retobj, err := env.CreateObjectFromFile(path)
	if err != nil {
		return nil, err
	}
	return retobj.(*{{ $type }}), nil
}

func (env *Environment) Update{{ .singular }}(obj *{{ $type }}) (*{{ $type }}, error) {
	retobj, err := env.UpdateObject(obj)
	if err != nil {
		return nil, err
	}
	return retobj.(*{{ $type }}), nil
}

func (env *Environment) Update{{ .singular }}FromFile(path string) (*{{ $type }}, error) {
	retobj, err := env.UpdateObjectFromFile(path)
	if err != nil {
		return nil, err
	}
	return retobj.(*{{ $type }}), nil
}

func (env *Environment) Patch{{ .singular }}({{ $namespaceParameter }}name string, patchType types.PatchType, patch []byte) (*{{ $type }}, error) {
	retobj, err := env.PatchObject({{ $gvk }}, {{ $namespaceValue }}, name, patchType, patch)
	if err != nil {
		return nil, err
	}
	return retobj.(*{{ $type }}), nil
}

func (env *Environment) Label{{ .singular }}({{ $namespaceParameter }}name string, key string, value string) (*{{ $type }}, error) {
	retobj, err := env.LabelObject({{ $gvk }}, {{ $namespaceValue }}, name, key, value)
	if err != nil {
		return nil, err
	}
	return retobj.(*{{ $type }}), nil
}

func (env *Environment) Unlabel{{ .singular }}({{ $namespaceParameter }}name string, key string) (*{{ $type }}, error) {
	retobj, err := env.UnlabelObject({{ $gvk }}, {{ $namespaceValue }}, name, key)
	if err != nil {
		return nil, err
	}
	return retobj.(*{{ $type }}), nil
}

func (env *Environment) Delete{{ .singular }}({{ $namespaceParameter }}name string) error {
	return env.DeleteObject({{ $gvk }}, {{ $namespaceValue }}, name)
}

func (env *Environment) WaitFor{{ .singular }}(obj *{{ $type }}, conditions ...watchtools.ConditionFunc) (*{{ $type }}, error) {
	retobj, err := env.WaitForObject(obj, conditions...)
	if err != nil {
		return nil, err
	}
	return retobj.(*{{ $type }}), nil
}

func (env *Environment) WaitFor{{ .singular }}FromFile(path string, conditions ...watchtools.ConditionFunc) (*{{ $type }}, error) {
	retobj, err := env.WaitForObjectFromFile(path, conditions...)
	if err != nil {
		return nil, err
	}
	return retobj.(*{{ $type }}), nil
}

func (must *Must) Assert{{ .singular }}(obj *{{ $type }}) {
	err := must.env.Assert{{ .singular }}(obj)
	must.handleError(err)
}

func (must *Must) Assert{{ .singular }}FromFile(path string) {
	err := must.env.Assert{{ .singular }}FromFile(path)
	must.handleError(err)
}

func (must *Must) Assert{{ .singular }}Count({{ $namespaceParameter }}labelSelector string, count int) {
	err := must.env.Assert{{ .singular }}Count({{ $namespaceArgument }}labelSelector, count)
	must.handleError(err)
}

func (must *Must) Get{{ .singular }}({{ $namespaceParameter }}name string) *{{ $type }} {
	retobj, err := must.env.Get{{ .singular }}({{ $namespaceArgument }}name)
	must.handleError(err)
	return retobj
}

func (must *Must) List{{ .plural }}({{ $namespaceParameter }}labelSelector string) []*{{ $type }} {
	retobjs, err := must.env.List{{ .plural }}({{ $namespaceArgument }}labelSelector)
	must.handleError(err)
	return retobjs
}

func (must *Must) Create{{ .singular }}(obj *{{ $type }}) *{{ $type }} {
	retobj, err := must.env.Create{{ .singular }}(obj)
	must.handleError(err)
	return retobj
}

func (must *Must) Create{{ .singular }}FromFile(path string) *{{ $type }} {
	retobj, err := must.env.Create{{ .singular }}FromFile(path)
	must.handleError(err)
	return retobj
}

func (must *Must) Update{{ .singular }}(obj *{{ $type }}) *{{ $type }} {
	retobj, err := must.env.Update{{ .singular }}(obj)
	must.handleError(err)
	return retobj
}

func (must *Must) Update{{ .singular }}FromFile(path string) *{{ $type }} {
	retobj, err := must.env.Update{{ .singular }}FromFile(path)
	must.handleError(err)
	return retobj
}

func (must *Must) Patch{{ .singular }}({{ $namespaceParameter }}name string, patchType types.PatchType, patch []byte) *{{ $type }} {
	retobj, err := must.env.Patch{{ .singular }}({{ $namespaceArgument }}name, patchType, patch)
	must.handleError(err)
	return retobj
}

func (must *Must) Label{{ .singular }}({{ $namespaceParameter }}name string, key string, value string) *{{ $type }} {
	retobj, err := must.env.Label{{ .singular }}({{ $namespaceArgument }}name, key, value)
	must.handleError(err)
	return retobj
}

func (must *Must) Unlabel{{ .singular }}({{ $namespaceParameter }}name string, key string) *{{ $type }} {
	retobj, err := must.env.Unlabel{{ .singular }}({{ $namespaceArgument }}name, key)
	must.handleError(err)
	return retobj
}

func (must *Must) Delete{{ .singular }}({{ $namespaceParameter }}name string) {
	err := must.env.Delete{{ .singular }}({{ $namespaceArgument }}name)
	must.handleError(err)
}

func (must *Must) WaitFor{{ .singular }}(obj *{{ $type }}, conditions ...watchtools.ConditionFunc) *{{ $type }} {
	retobj, err := must.env.WaitFor{{ .singular }}(obj, conditions...)
	must.handleError(err)
	return retobj
}

func (must *Must) WaitFor{{ .singular }}FromFile(path string, conditions ...watchtools.ConditionFunc) *{{ $type }} {
	retobj, err := must.env.WaitFor{{ .singular }}FromFile(path, conditions...)
	must.handleError(err)
	return retobj
}
{{- end }}
