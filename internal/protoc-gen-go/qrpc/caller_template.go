package qrpc

const callerTemplate = `
{{- "" -}}
{{- $contextPkg := .ContextPkg -}}
{{- $qrpcPkg := .QrpcPkg -}}
{{- $uuidPkg := .UUIDPkg -}}

{{- range .Services }}
{{- $service := .Name -}}
{{- $clientInt := ($service | sprintf "Caller%sClient") -}}
{{- $clientType := ($service | sprintf "caller%sClient") -}}
{{- $serverInt := ($service | sprintf "Caller%sServer") }}

// Caller API description

// {{ $clientInt }} is the caller client API for {{ $service }} service.
{{- if .IsDeprecated }}
// Deprecated: Do not use
{{- end }}
type {{ $clientInt }} interface {
	{{- range (.ForwardMethods | concat .BidirectionalMethods) }}
    {{ .Comment }}
	{{- if .IsDeprecated }}
	// Deprecated: Do not use.
	{{- end }}
	{{ .Name }}(ctx {{ $contextPkg }}.Context, in *{{ .InType }}) (string, error)
	{{- end }}
}

type {{ $clientType }} struct {
	cc *{{ $qrpcPkg }}.ClientConn
}

{{- if .IsDeprecated }}
// Deprecated: Do not use.
{{- end }}
func New{{ $clientInt }}(cc *{{ $qrpcPkg }}.ClientConn) {{ $clientInt }} {
	cc.SetService("qrpc.test.api.{{ $service }}.in")
	return &{{ $clientType }}{cc}
}

{{- range (.ForwardMethods | concat .BidirectionalMethods) }}

{{ .Comment }}
{{- if .IsDeprecated }}
// Deprecated: Do not use.
{{- end }}
func (c *{{ $clientType }}) {{ .Name }}(ctx {{ $contextPkg }}.Context, in *{{ .InType }}) (string, error) {
	data, err := proto.Marshal(in)
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

{{- if (.BackwardMethods | concat .BidirectionalMethods) }}

// {{ $serverInt }} is the caller server API for {{ $service }} service.
// It should be used to handle responses from backward methods and bidirectional methods.
{{- if .IsDeprecated }}
// Deprecated: Do not use.
{{- end }}
type {{ $serverInt }} interface {
	{{- range (.BackwardMethods | concat .BidirectionalMethods) }}
	{{ .Comment }}
	{{- if .IsDeprecated }}
	// Deprecated: Do not use.
	{{- end }}
	{{ .Name }}(ctx {{ $contextPkg }}.Context, reqID string, out *{{ .OutType }}) error
	{{- end }}
}

{{- if .IsDeprecated }}
// Deprecated: Do not use.
{{- end }}
func Register{{ $serverInt }}(s *{{ $qrpcPkg }}.Server, srv {{ $serverInt }}) {
	s.RegisterService(&_Caller{{ $service }}_serviceDesc, srv)
}

{{- range (.BackwardMethods | concat .BidirectionalMethods) }}

func _Caller{{ $service }}_{{ .Name }}_Handler(srv interface{}, ctx {{ $contextPkg }}.Context, reqID string, msg []byte) error {
	out := new({{ .OutType }})

	if err := proto.Unmarshal(msg, out); err != nil {
		return err
	}

	return srv.({{ $serverInt }}).{{ .Name }}(ctx, reqID, out)
}
{{- end }}

var _Caller{{ $service }}_serviceDesc = {{ $qrpcPkg }}.ServiceDesc{
	ServiceName: "qrpc.test.api.{{ $service }}.out",
	HandlerType: (*{{ $serverInt }})(nil),
	Methods: []{{ $qrpcPkg }}.MethodDesc{
		{{- range (.BackwardMethods | concat .BidirectionalMethods) }}
		{
			MethodName: "{{ .Name }}",
			Handler: _Caller{{ $service }}_{{ .Name }}_Handler,
		},
		{{- end }}
	},
}
{{- end }}
{{- end }}
`
