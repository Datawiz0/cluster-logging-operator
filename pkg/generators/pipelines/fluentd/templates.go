package fluentd

const outputLabelMatchTemplate = `{{define "outputLabelMatch" -}}
<match {{.Tags}}>
    @type relabel
    @label {{labelName .Source}}
</match>
{{- end}}`

const sourceLabelCopyTemplate = `{{define "sourceLabelCopy" -}}
<label {{labelName .Source}}>
    <match **>
		@type copy
{{ range $index, $target := .TargetLabels }}
        <store>
            @type relabel
            @label {{$target}}
        </store>
{{ end }}
    </match>
</label>
{{- end}}`

const outputLabelConfTemplate = `{{define "outputLabelConf" -}}
<label {{.ReLabelName}}>
	<match {{.RetryTag}}>
		@type copy
{{include "store_elasticsearch" . "prefix_as_retry" | indent 4}}
	</match>
	<match **>
		@type copy
{{include "store_elasticsearch" . "include_retry_tag"| indent 4}}
	</match>
</label>
{{- end}}`

const storeElasticsearchTemplate = `{{- define "store_elasticsearch" -}}
<store>
	@type elasticsearch
	@id {{.StoreID }}
	host {{.Host}}
	port {{.Port}}
	{{- if .Target.Certificates }}
	scheme https
	ssl_version TLSv1_2
	{{ else }}
	scheme http
	{{ end -}}
	target_index_key viaq_index_name
	id_key viaq_msg_id
	remove_keys viaq_index_name
	user fluentd
	password changeme
	{{ if .Target.Certificates }}
	client_key {{ .SecretPath "key"}}
	client_cert {{ .SecretPath "cert"}}
	ca_file {{ .SecretPath "cacert"}}
	{{ end -}}
	type_name com.redhat.viaq.common
	{{if .Hints.Has "include_retry_tag" -}}
	retry_tag {{.RetryTag}}
	{{end -}}
	write_operation create
	reload_connections "#{ENV['ES_RELOAD_CONNECTIONS'] || 'true'}"
	# https://github.com/uken/fluent-plugin-elasticsearch#reload-after
	reload_after "#{ENV['ES_RELOAD_AFTER'] || '200'}"
	# https://github.com/uken/fluent-plugin-elasticsearch#sniffer-class-name
	sniffer_class_name "#{ENV['ES_SNIFFER_CLASS_NAME'] || 'Fluent::Plugin::ElasticsearchSimpleSniffer'}"
	reload_on_failure false
	# 2 ^ 31
	request_timeout 2147483648
	<buffer>
		@type file
		path '{{.BufferPath}}'
		flush_interval "#{ENV['ES_FLUSH_INTERVAL'] || '1s'}"
		flush_thread_count "#{ENV['ES_FLUSH_THREAD_COUNT'] || 2}"
		flush_at_shutdown "#{ENV['FLUSH_AT_SHUTDOWN'] || 'false'}"
		retry_max_interval "#{ENV['ES_RETRY_WAIT'] || '300'}"
		retry_forever true
		queued_chunks_limit_size "#{ENV['BUFFER_QUEUE_LIMIT'] || '32' }"
		chunk_limit_size "#{ENV['BUFFER_SIZE_LIMIT'] || '8m' }"
		overflow_action "#{ENV['BUFFER_QUEUE_FULL_ACTION'] || 'block'}"
	</buffer>
</store>
{{- end}}`
