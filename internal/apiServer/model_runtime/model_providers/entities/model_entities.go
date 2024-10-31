package entities

type ModelType string

const (
	LLM            ModelType = "llm"
	TEXT_EMBEDDING ModelType = "text-embedding"
	RERANK         ModelType = "rerank"
	SPEECH2TEXT    ModelType = "speech2text"
	MODERATION     ModelType = "moderation"
	TTS            ModelType = "tts"
	TEXT2IMG       ModelType = "text2img"
)

type ModelFeature string

const (
	TOOL_CALL        ModelFeature = "tool-call"
	MULTI_TOOL_CALL  ModelFeature = "multi-tool-call"
	AGENT_THOUGHT    ModelFeature = "agent-thought"
	VISION           ModelFeature = "vision"
	STREAM_TOOL_CALL ModelFeature = "stream-tool-call"
)

type FetchFrom string

const (
	PREDEFINED_MODEL_FROM   FetchFrom = "predefined-model"
	CUSTOMIZABLE_MODEL_FROM FetchFrom = "customizable-model"
)

type ModelPropertyKey string

const (
	MODE                      ModelPropertyKey = "mode"
	CONTEXT_SIZE              ModelPropertyKey = "context_size"
	MAX_CHUNKS                ModelPropertyKey = "max_chunks"
	FILE_UPLOAD_LIMIT         ModelPropertyKey = "file_upload_limit"
	SUPPORTED_FILE_EXTENSIONS ModelPropertyKey = "supported_file_extensions"
	MAX_CHARACTERS_PER_CHUNK  ModelPropertyKey = "max_characters_per_chunk"
	DEFAULT_VOICE             ModelPropertyKey = "default_voice"
	VOICES                    ModelPropertyKey = "voices"
	WORD_LIMIT                ModelPropertyKey = "word_limit"
	AUDIO_TYPE                ModelPropertyKey = "audio_type"
	MAX_WORKERS               ModelPropertyKey = "max_workers"
)

type ProviderModel struct {
	Model           string                           `json:"model"`
	Label           I18nObject                       `json:"label"`
	ModelType       ModelType                        `json:"model_type"`
	Features        []ModelFeature                   `json:"features"`
	FetchFrom       FetchFrom                        `json:"fetch_from"`
	ModelProperties map[ModelPropertyKey]interface{} `json:"model_properties"`
	Deprecated      bool                             `json:"deprecated"`
	ModelConfig     interface{}                      `json:"model_config"`
}
