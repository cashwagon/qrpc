package qrpc

const handlerTemplate = `
{{- "" -}}
{{- $contextPkg := .ContextPkg -}}
{{- $qrpcPkg := .QrpcPkg -}}
{{- $uuidPkg := .UUIDPkg -}}

{{- range .Services }}
{{- $service := .Name -}}
{{- $clientInt := ($service | sprintf "Handler%sClient") -}}
{{- $clientType := ($service | sprintf "handler%sClient") -}}
{{- $serverInt := ($service | sprintf "Handler%sServer") }}

// Handler API description

{{- if (.BackwardMethods | concat .BidirectionalMethods) }}

// {{ $clientInt }} is the handler client API for {{ $service }} service.
// It should be used to send response to the caller.
{{- if .IsDeprecated }}
// Deprecated: Do not use.
{{- end }}
type {{ $clientInt }} interface {
	{{- range .BidirectionalMethods }}
	{{ .Comment }}
	{{- if .IsDeprecated }}
	// Deprecated: Do not use.
	{{- end }}
	{{ .Name }}(ctx {{ $contextPkg }}.Context, reqID string, out *{{ .OutType }}) error
	{{- end }}

	{{- range .BackwardMethods }}
	{{ .Comment }}
	{{- if .IsDeprecated }}
	// Deprecated: Do not use.
	{{- end }}
	{{ .Name }}(ctx {{ $contextPkg }}.Context, out *{{ .OutType }}) (string, error)
	{{- end }}
}

type {{ $clientType }} struct {
	cc *{{ $qrpcPkg }}.ClientConn
}

{{- if .IsDeprecated }}
// Deprecated: Do not use.
{{- end }}
func New{{ $clientInt }}(cc *{{ $qrpcPkg }}.ClientConn) {{ $clientInt }} {
	cc.SetService("qrpc.test.api.{{ $service }}.out")
	return &{{ $clientType }}{cc}
}

{{- range .BidirectionalMethods }}

{{ .Comment }}
{{- if .IsDeprecated }}
// Deprecated: Do not use.
{{- end }}
func (c *{{ $clientType }}) {{ .Name }}(ctx {{ $contextPkg }}.Context, reqID string, out *{{ .OutType }}) error {
	data, err := proto.Marshal(out)
	if err != nil {
		return err
	}

	return c.cc.Invoke(ctx, {{ $qrpcPkg }}.Message{
		Method:    "{{ .Name }}",
		RequestID: reqID,
		Data:      data,
	})
}
{{- end }}

{{- range .BackwardMethods }}

{{ .Comment }}
{{- if .IsDeprecated }}
// Deprecated: Do not use.
{{- end }}
func (c *{{ $clientType }}) {{ .Name }}(ctx {{ $contextPkg }}.Context, out *{{ .OutType }}) (string, error) {
	data, err := proto.Marshal(out)
	if err != nil {
		return "", err
	}

	msg := {{ $qrpcPkg }}.Message{
		Method:    "{{ .Name }}",
		RequestID: {{ $uuidPkg }}.New().String(),
		Data:      data,
	}

	return msg.RequestID, c.cc.Invoke(ctx, msg)
}
{{- end }}
{{- end }}

// {{ $serverInt }} is the handler server API for {{ $service }} service.
{{- if .IsDeprecated }}
// Deprecated: Do not use.
{{- end }}
type {{ $serverInt }} interface {
	{{- range (.ForwardMethods | concat .BidirectionalMethods) }}
	{{ .Comment }}
	{{- if .IsDeprecated }}
	// Deprecated: Do not use.
	{{- end }}
	{{ .Name }}(ctx {{ $contextPkg }}.Context, reqID string, in *{{ .InType }}) error
	{{- end }}
}

{{- if .IsDeprecated }}
// Deprecated: Do not use.
{{- end }}
func Register{{ $serverInt }}(s *{{ $qrpcPkg }}.Server, srv {{ $serverInt }}) {
	s.RegisterService(&_Handler{{ $service }}_serviceDesc, srv)
}

{{- range (.ForwardMethods | concat .BidirectionalMethods) }}

func _Handler{{ $service }}_{{ .Name }}_Handler(srv interface{}, ctx {{ $contextPkg }}.Context, reqID string, msg []byte) error {
	in := new({{ .InType }})

	if err := proto.Unmarshal(msg, in); err != nil {
		return err
	}

	return srv.({{ $serverInt }}).{{ .Name }}(ctx, reqID, in)
}
{{- end }}

var _Handler{{ $service }}_serviceDesc = {{ $qrpcPkg }}.ServiceDesc{
	ServiceName: "qrpc.test.api.{{ $service }}.in",
	HandlerType: (*{{ $serverInt }})(nil),
	Methods: []{{ $qrpcPkg }}.MethodDesc{
		{{- range (.ForwardMethods | concat .BidirectionalMethods) }}
		{
			MethodName: "{{ .Name }}",
			Handler: _Handler{{ $service }}_{{ .Name }}_Handler,
		},
		{{- end }}
	},
}
{{- end }}
`
